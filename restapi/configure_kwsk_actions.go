package restapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	middleware "github.com/go-openapi/runtime/middleware"

	models "github.com/projectodd/kwsk/models"
	"github.com/projectodd/kwsk/restapi/operations"
	"github.com/projectodd/kwsk/restapi/operations/actions"

	corev1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"

	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	knative "github.com/knative/serving/pkg/client/clientset/versioned"
)

func configureActions(api *operations.KwskAPI, knativeClient *knative.Clientset, cache cache.Store) {
	api.ActionsDeleteActionHandler = actions.DeleteActionHandlerFunc(deleteActionFunc(knativeClient))

	api.ActionsGetActionByNameHandler = actions.GetActionByNameHandlerFunc(getActionByNameFunc(knativeClient))

	api.ActionsGetAllActionsHandler = actions.GetAllActionsHandlerFunc(getAllActionsFunc(knativeClient))

	api.ActionsInvokeActionHandler = actions.InvokeActionHandlerFunc(invokeActionFunc(knativeClient, cache))

	api.ActionsUpdateActionHandler = actions.UpdateActionHandlerFunc(updateActionFunc(knativeClient))
}

func deleteActionFunc(knativeClient *knative.Clientset) actions.DeleteActionHandlerFunc {
	return func(params actions.DeleteActionParams, principal *models.Principal) middleware.Responder {
		serviceName := sanitizeObjectName(params.ActionName)
		namespace := namespaceOrDefault(params.Namespace)
		err := knativeClient.ServingV1alpha1().Services(namespace).Delete(serviceName, &metav1.DeleteOptions{})
		if err != nil {
			errorMessage := errorMessageFromErr(err)
			if k8sErrors.IsNotFound(err) {
				return actions.NewDeleteActionNotFound().WithPayload(errorMessage)
			}
			return actions.NewDeleteActionInternalServerError().WithPayload(errorMessage)
		}

		// Wait for the action to be deleted
		deleteTimeout := 1 * time.Minute
		err = wait.PollImmediate(1*time.Second, deleteTimeout, func() (bool, error) {
			_, serviceErr := knativeClient.ServingV1alpha1().Services(namespace).Get(serviceName, metav1.GetOptions{})
			_, routeErr := knativeClient.ServingV1alpha1().Routes(namespace).Get(serviceName, metav1.GetOptions{})
			if serviceErr != nil && k8sErrors.IsNotFound(serviceErr) && routeErr != nil && k8sErrors.IsNotFound(routeErr) {
				// Service and Route are gone, so that's good enough
				return true, nil
			}
			if serviceErr != nil {
				return false, err
			}
			return false, nil
		})

		if err != nil {
			fmt.Printf("Error waiting on action to delete: %s\n", err)
			return actions.NewDeleteActionInternalServerError()
		}

		return actions.NewDeleteActionOK()
	}
}

func updateActionFunc(knativeClient *knative.Clientset) actions.UpdateActionHandlerFunc {
	return func(params actions.UpdateActionParams, principal *models.Principal) middleware.Responder {
		name := params.ActionName
		serviceName := sanitizeObjectName(name)
		namespace := namespaceOrDefault(params.Namespace)
		kind := params.Action.Exec.Kind
		version := params.Action.Version
		if version == "" {
			version = "0.0.1"
		}

		annotations := make(map[string]string)
		annotations[KwskName] = name
		annotations[KwskVersion] = version

		for _, kv := range params.Action.Annotations {
			key := fmt.Sprintf("kwsk_action_anno_%s", kv.Key)
			value := fmt.Sprintf("%s", kv.Value)
			annotations[key] = value
		}

		var image string
		if params.Action.Exec != nil {
			image = params.Action.Exec.Image
			annotations[KwskKind] = kind
			annotations[KwskImage] = image
		}

		configSpec := v1alpha1.ConfigurationSpec{
			RevisionTemplate: v1alpha1.RevisionTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{},
				Spec:       v1alpha1.RevisionSpec{},
			},
		}

		if image == "" {
			prefix := imagePrefix()
			tag := imageTag()
			switch kind {
			case "nodejs:default", "nodejs:6":
				image = "kwsk-nodejs6action"
			case "nodejs:8":
				image = "kwsk-action-nodejs-v8"
			case "python:default", "python:2":
				image = "kwsk-python2action"
			case "python:3":
				image = "kwsk-python3action"
			case "java:default", "java":
				image = "kwsk-java8action"
			}
			image = fmt.Sprintf("%s/%s:%s", prefix, image, tag)
		}

		actionParamsMap := map[string]interface{}{}
		for _, kv := range params.Action.Parameters {
			actionParamsMap[kv.Key] = kv.Value
		}
		actionParamsJson, err := json.Marshal(actionParamsMap)
		if err != nil {
			fmt.Println("Error marshaling action parameters: ", err)
			return actions.NewUpdateActionInternalServerError()
		}

		containerEnv := []corev1.EnvVar{
			corev1.EnvVar{
				Name:  "KWSK_ACTION_CODE",
				Value: params.Action.Exec.Code,
			},
			corev1.EnvVar{
				Name:  "KWSK_ACTION_PARAMS",
				Value: string(actionParamsJson),
			},
		}
		container := corev1.Container{
			Image: image,
			Env:   containerEnv,
		}
		configSpec.RevisionTemplate.Spec.Container = container

		service := &v1alpha1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:        serviceName,
				Namespace:   namespace,
				Annotations: annotations,
			},
			Spec: v1alpha1.ServiceSpec{
				RunLatest: &v1alpha1.RunLatestType{
					Configuration: configSpec,
				},
			},
		}

		dbg := fmt.Sprintf("Creating service %+v\n", service)
		fmt.Printf("%.2000s\n", dbg)
		service, err = knativeClient.ServingV1alpha1().Services(namespace).Create(service)
		if err != nil {
			fmt.Println("Error updating action: ", err)
			return actions.NewUpdateActionInternalServerError().WithPayload(errorMessageFromErr(err))
		}

		action, err := getActionByName(knativeClient, name, namespace, false)
		if err != nil {
			fmt.Println("Error retrieving updated action: ", err)
			return actions.NewUpdateActionInternalServerError().WithPayload(errorMessageFromErr(err))
		}
		return actions.NewUpdateActionOK().WithPayload(action)
	}
}

func serviceToAction(service *v1alpha1.Service) *models.Action {
	objectMeta := service.ObjectMeta
	name := objectMeta.Annotations[KwskName]
	if name == "" {
		name = objectMeta.Name
	}
	version := objectMeta.Annotations[KwskVersion]
	kind := objectMeta.Annotations[KwskKind]
	image := objectMeta.Annotations[KwskImage]

	var code string
	var actionParams map[string]interface{}
	configurationSpec := service.Spec.RunLatest.Configuration
	for _, env := range configurationSpec.RevisionTemplate.Spec.Container.Env {
		if env.Name == "KWSK_ACTION_CODE" {
			code = env.Value
		}
		if env.Name == "KWSK_ACTION_PARAMS" {
			err := json.Unmarshal([]byte(env.Value), &actionParams)
			if err != nil {
				fmt.Println("Failed to unmarshal action parameters:", err)
			}
		}
	}

	var params []*models.KeyValue
	for key, value := range actionParams {
		param := &models.KeyValue{
			Key:   key,
			Value: value,
		}
		params = append(params, param)
	}

	var annotations []*models.KeyValue
	for key, value := range objectMeta.Annotations {
		if strings.HasPrefix(key, "kwsk_action_anno_") {
			annotation := &models.KeyValue{
				Key:   strings.TrimPrefix(key, "kwsk_action_anno_"),
				Value: value,
			}
			annotations = append(annotations, annotation)
		}
	}

	publish := false
	return &models.Action{
		Name:      &name,
		Namespace: &objectMeta.Namespace,
		Version:   &version,
		Exec: &models.ActionExec{
			Image: image,
			Kind:  kind,
			Code:  code,
		},
		Parameters:  params,
		Annotations: annotations,
		Publish:     &publish,
	}
}

func getActionByName(knativeClient *knative.Clientset, name string, namespace string, includeCode bool) (*models.Action, error) {
	serviceName := sanitizeObjectName(name)
	namespace = namespaceOrDefault(namespace)
	service, err := knativeClient.ServingV1alpha1().Services(namespace).Get(serviceName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	action := serviceToAction(service)
	if !includeCode {
		action.Exec.Code = ""
	}
	return action, nil
}

func getActionByNameFunc(knativeClient *knative.Clientset) actions.GetActionByNameHandlerFunc {
	return func(params actions.GetActionByNameParams, principal *models.Principal) middleware.Responder {
		var code bool
		if params.Code != nil {
			code = *params.Code
		}
		action, err := getActionByName(knativeClient, params.ActionName, params.Namespace, code)
		if err != nil {
			errorMessage := errorMessageFromErr(err)
			if k8sErrors.IsNotFound(err) {
				return actions.NewGetActionByNameNotFound().WithPayload(errorMessage)
			}
			return actions.NewGetActionByNameInternalServerError().WithPayload(errorMessage)
		}
		return actions.NewGetActionByNameOK().WithPayload(action)
	}
}

func getAllActionsFunc(knativeClient *knative.Clientset) actions.GetAllActionsHandlerFunc {
	return func(params actions.GetAllActionsParams, principal *models.Principal) middleware.Responder {
		namespace := namespaceOrDefault(params.Namespace)
		services, err := knativeClient.ServingV1alpha1().Services(namespace).List(metav1.ListOptions{})
		if err != nil {
			return actions.NewGetAllActionsInternalServerError().WithPayload(errorMessageFromErr(err))
		}
		var payload = make([]*models.Action, len(services.Items))
		for i, service := range services.Items {
			payload[i] = serviceToAction(&service)
		}
		return actions.NewGetAllActionsOK().WithPayload(payload)
	}
}

type ActionInitMessage struct {
	Value ActionInitValue `json:"value,omitempty"`
}

type ActionInitValue struct {
	Main string `json:"main,omitempty"`
	Code string `json:"code,omitempty"`
}

type ActionRunMessage struct {
	Value interface{} `json:"value"`
}

func withRoutesReady(knativeClient *knative.Clientset, service *v1alpha1.Service) (*v1alpha1.Service, error) {
	// Wait for the service routes to be ready
	readyTimeout := 5 * time.Minute
	if !serviceRoutesReady(service) {
		err := wait.Poll(1*time.Second, readyTimeout, func() (bool, error) {
			newService, err := knativeClient.ServingV1alpha1().Services(service.Namespace).Get(service.Name, metav1.GetOptions{})
			if err != nil {
				if k8sErrors.IsNotFound(err) {
					// If not found then keep trying, assuming it will
					// be found later
					return false, nil
				}
				fmt.Println("Error waiting for service route readiness: ", err)
				return false, err
			}
			service = newService
			return serviceRoutesReady(service), nil
		})
		if err != nil {
			fmt.Printf("Error waiting on service to become ready: %s\n", err)
			return service, err
		}
	}
	return service, nil
}

func serviceRoutesReady(service *v1alpha1.Service) bool {
	if c := service.Status.GetCondition(v1alpha1.ServiceConditionRoutesReady); c != nil {
		return c.Status == corev1.ConditionTrue
	}
	return false
}

func invokeActionFunc(knativeClient *knative.Clientset, cache cache.Store) actions.InvokeActionHandlerFunc {
	return func(params actions.InvokeActionParams, principal *models.Principal) middleware.Responder {
		serviceName := sanitizeObjectName(params.ActionName)
		namespace := namespaceOrDefault(params.Namespace)
		blocking := params.Blocking != nil && *params.Blocking == "true"
		result := params.Result != nil && *params.Result == "true"

		service, err := knativeClient.ServingV1alpha1().Services(namespace).Get(serviceName, metav1.GetOptions{})
		if err != nil {
			errorMessage := errorMessageFromErr(err)
			if k8sErrors.IsNotFound(err) {
				return actions.NewInvokeActionNotFound().WithPayload(errorMessage)
			}
			return actions.NewInvokeActionInternalServerError().WithPayload(errorMessage)
		}
		service, err = withRoutesReady(knativeClient, service)
		if err != nil {
			fmt.Printf("Error waiting for action readiness: %s\n", err)
			msg := "Error waiting for action readiness\n"
			errorMessage := &models.ErrorMessage{
				Error: &msg,
			}
			return actions.NewInvokeActionInternalServerError().WithPayload(errorMessage)
		}
		actionHost := service.Status.Domain
		istioHostAndPort := istioHostAndPort()

		return runAction(istioHostAndPort, actionHost, service.Name, namespace, params.Payload, blocking, result, cache)
	}
}

func runAction(istioHostAndPort string, actionHost string, name string, namespace string, params interface{}, blocking bool, result bool, cache cache.Store) middleware.Responder {

	start := time.Now().UnixNano() / 1000000 // milliseconds since epoch

	var resStatus int
	var resBody []byte
	var err error

	// Wait for the action to be ready
	readyTimeout := 5 * time.Minute
	err = wait.PollImmediate(1*time.Second, readyTimeout, func() (bool, error) {
		resStatus, resBody, err = actionRequest(istioHostAndPort, actionHost, "", params)
		if err != nil {
			fmt.Printf("Action not yet ready: %s\n", err)
			return false, err
		}
		return resStatus != http.StatusNotFound && resStatus != http.StatusServiceUnavailable, nil
	})

	if err != nil {
		fmt.Printf("Error waiting on action to become ready: %s\n", err)
		return actions.NewInvokeActionInternalServerError().WithPayload(errorMessageFromErr(err))
	}

	if resStatus != http.StatusOK {
		msg := fmt.Sprintf("Error invoking action. Status: %d, Message: %s\n", resStatus, resBody)
		errorMessage := &models.ErrorMessage{
			Error: &msg,
		}
		fmt.Println(msg)
		return actions.NewInvokeActionInternalServerError().WithPayload(errorMessage)
	}

	var resultJson interface{}
	err = json.Unmarshal(resBody, &resultJson)
	if err != nil {
		msg := fmt.Sprintf("Action invocation result was not valid JSON. Result: %s\n", resStatus, resBody)
		errorMessage := &models.ErrorMessage{
			Error: &msg,
		}
		fmt.Println(msg)
		return actions.NewInvokeActionInternalServerError().WithPayload(errorMessage)
	}
	activationResult := &models.ActivationResult{
		Result:  resultJson,
		Success: true,
		Status:  "success",
	}

	activationId, err := newActivationId()
	if err != nil {
		fmt.Printf("Error generating activationId: %s\n", err)
		return actions.NewInvokeActionInternalServerError().WithPayload(errorMessageFromErr(err))
	}
	logs := []string{}
	annotations := []*models.KeyValue{}
	end := time.Now().UnixNano() / 1000000 // milliseconds since epoch
	duration := end - start
	activation := models.Activation{
		ActivationID: &activationId,
		Name:         &name,
		Namespace:    &namespace,
		Response:     activationResult,
		Start:        &start,
		End:          end,
		Duration:     duration,
		Logs:         logs,
		Annotations:  annotations,
	}

	err = cache.Add(activation)
	if err != nil {
		msg := fmt.Sprintf("Error storing activation record: %s\n", err)
		errorMessage := &models.ErrorMessage{
			Error: &msg,
		}
		fmt.Println(msg)
		return actions.NewInvokeActionInternalServerError().WithPayload(errorMessage)
	}

	if result {
		fmt.Printf("Returning Activation result: %+v\n", resultJson)
		return actions.NewInvokeActionOK().WithPayload(resultJson)
	} else if blocking {
		fmt.Printf("Returning Activation: %+v\n", activation)
		return actions.NewInvokeActionOK().WithPayload(activation)
	} else {
		payload := &models.ActivationID{
			ActivationID: &activationId,
		}
		fmt.Printf("Returning ActivationId: %s\n", activationId)
		return actions.NewInvokeActionAccepted().WithPayload(payload)
	}
}

func actionRequest(istioHostAndPort string, actionHost string, path string, requestBody interface{}) (int, []byte, error) {
	url := fmt.Sprintf("http://%s/%s", istioHostAndPort, path)
	fmt.Printf("Sending POST to url %s with host %s\n", url, actionHost)

	body, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Printf("Error marshaling action request body: %s\n", err)
		return 500, nil, err
	}
	fmt.Printf("Request Body: %s\n", body)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Error creating http request for action: %s\n", err)
		return 500, nil, err
	}

	req.Host = actionHost
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error calling action http endpoint: %s\n", err)
		return 500, nil, err
	}

	defer res.Body.Close()
	resBody, _ := ioutil.ReadAll(res.Body)
	fmt.Printf("Response: %+v\n", res)
	fmt.Printf("Response Body: %s\n", string(resBody))

	return res.StatusCode, resBody, nil
}
