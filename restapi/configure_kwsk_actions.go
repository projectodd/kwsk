package restapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	models "github.com/projectodd/kwsk/models"
	"github.com/projectodd/kwsk/restapi/operations"
	"github.com/projectodd/kwsk/restapi/operations/actions"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	build "github.com/knative/build/pkg/apis/build/v1alpha1"
	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	knative "github.com/knative/serving/pkg/client/clientset/versioned"
)

func configureActions(api *operations.KwskAPI, knativeClient *knative.Clientset) {
	api.ActionsDeleteActionHandler = actions.DeleteActionHandlerFunc(deleteActionFunc(knativeClient))

	api.ActionsGetActionByNameHandler = actions.GetActionByNameHandlerFunc(getActionByNameFunc(knativeClient))

	api.ActionsGetAllActionsHandler = actions.GetAllActionsHandlerFunc(getAllActionsFunc(knativeClient))

	api.ActionsInvokeActionHandler = actions.InvokeActionHandlerFunc(invokeActionFunc(knativeClient))

	api.ActionsUpdateActionHandler = actions.UpdateActionHandlerFunc(updateActionFunc(knativeClient))
}

func deleteActionFunc(knativeClient *knative.Clientset) actions.DeleteActionHandlerFunc {
	return func(params actions.DeleteActionParams) middleware.Responder {
		namespace := namespaceOrDefault(params.Namespace)
		err := knativeClient.ServingV1alpha1().Configurations(namespace).Delete(params.ActionName, &metav1.DeleteOptions{})
		if err != nil {
			msg := err.Error()
			errorMessage := &models.ErrorMessage{
				Error: &msg,
			}
			if errors.IsNotFound(err) {
				return actions.NewDeleteActionNotFound().WithPayload(errorMessage)
			}
			return actions.NewDeleteActionInternalServerError().WithPayload(errorMessage)
		}

		err = knativeClient.ServingV1alpha1().Routes(namespace).Delete(params.ActionName, &metav1.DeleteOptions{})
		if err != nil {
			msg := err.Error()
			errorMessage := &models.ErrorMessage{
				Error: &msg,
			}
			if errors.IsNotFound(err) {
				return actions.NewDeleteActionNotFound().WithPayload(errorMessage)
			}
			return actions.NewDeleteActionInternalServerError().WithPayload(errorMessage)
		}

		return actions.NewDeleteActionOK()
	}
}

func updateActionFunc(knativeClient *knative.Clientset) actions.UpdateActionHandlerFunc {
	return func(params actions.UpdateActionParams) middleware.Responder {
		name := params.ActionName
		namespace := namespaceOrDefault(params.Namespace)

		annotations := make(map[string]string)
		annotations["kwsk_action_version"] = params.Action.Version
		annotations["kwsk_action_kind"] = *params.Action.Exec.Kind
		annotations["kwsk_action_code"] = params.Action.Exec.Code

		config := &v1alpha1.Configuration{
			ObjectMeta: metav1.ObjectMeta{
				Name:        name,
				Namespace:   namespace,
				Annotations: annotations,
			},
			Spec: v1alpha1.ConfigurationSpec{
				RevisionTemplate: v1alpha1.RevisionTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{},
					Spec:       v1alpha1.RevisionSpec{},
				},
			},
		}

		image := params.Action.Exec.Image
		if image == "" {
			// TODO: Map the kind of the action to an image instead of
			// just assuming everything is node8
			image = "openwhisk/action-nodejs-v8"
			// TODO: This is just a dummy, placeholder BuildSpec. We
			// don't actually use it for anything.
			buildSpec := &build.BuildSpec{
				Steps: []corev1.Container{
					corev1.Container{
						Image:   image,
						Command: []string{"/bin/bash"},
						Args:    []string{"-c", "echo 'hi'"},
					},
				},
			}
			config.Spec.Build = buildSpec
		}
		container := corev1.Container{
			Image: image,
		}
		config.Spec.RevisionTemplate.Spec.Container = container

		fmt.Printf("Creating configuration %+v\n", config)
		createdConfig, err := knativeClient.ServingV1alpha1().Configurations(namespace).Create(config)
		if err != nil {
			msg := err.Error()
			errorMessage := &models.ErrorMessage{
				Error: &msg,
			}
			return actions.NewUpdateActionInternalServerError().WithPayload(errorMessage)
		}

		route := &v1alpha1.Route{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Spec: v1alpha1.RouteSpec{
				Traffic: []v1alpha1.TrafficTarget{
					v1alpha1.TrafficTarget{
						ConfigurationName: name,
						Percent:           100,
					},
				},
			},
		}
		_, err = knativeClient.ServingV1alpha1().Routes(namespace).Create(route)
		id := createdConfig.ObjectMeta.SelfLink
		itemId := &models.ItemID{
			ID: &id,
		}
		return actions.NewUpdateActionOK().WithPayload(itemId)
	}
}

func getActionByNameFunc(knativeClient *knative.Clientset) actions.GetActionByNameHandlerFunc {
	return func(params actions.GetActionByNameParams) middleware.Responder {
		namespace := namespaceOrDefault(params.Namespace)
		config, err := knativeClient.ServingV1alpha1().Configurations(namespace).Get(params.ActionName, metav1.GetOptions{})
		if err != nil {
			msg := err.Error()
			errorMessage := &models.ErrorMessage{
				Error: &msg,
			}
			if errors.IsNotFound(err) {
				return actions.NewGetActionByNameNotFound().WithPayload(errorMessage)
			}
			return actions.NewGetActionByNameInternalServerError().WithPayload(errorMessage)
		}
		objectMeta := config.ObjectMeta
		kind := objectMeta.Annotations["kwsk_action_kind"]
		version := objectMeta.Annotations["kwsk_action_version"]
		code := objectMeta.Annotations["kwsk_action_code"]
		payload := &models.Action{
			Name:      &objectMeta.Name,
			Namespace: &objectMeta.Namespace,
			Version:   &version,
			Exec: &models.ActionExec{
				Image: config.Spec.RevisionTemplate.Spec.Container.Image,
				Kind:  &kind,
				Code:  code,
			},
		}
		return actions.NewGetActionByNameOK().WithPayload(payload)
	}
}

func getAllActionsFunc(knativeClient *knative.Clientset) actions.GetAllActionsHandlerFunc {
	return func(params actions.GetAllActionsParams) middleware.Responder {
		namespace := namespaceOrDefault(params.Namespace)
		configs, err := knativeClient.ServingV1alpha1().Configurations(namespace).List(metav1.ListOptions{})
		if err != nil {
			msg := err.Error()
			errorMessage := &models.ErrorMessage{
				Error: &msg,
			}
			return actions.NewGetAllActionsInternalServerError().WithPayload(errorMessage)
		}
		var payload = make([]*models.EntityBrief, len(configs.Items))
		for i, item := range configs.Items {
			name := item.ObjectMeta.Name
			namespace := item.ObjectMeta.Namespace
			version := item.ObjectMeta.Annotations["kwsk_action_version"]
			payload[i] = &models.EntityBrief{
				Name:      &name,
				Namespace: &namespace,
				Version:   &version,
			}
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
	Value interface{} `json:"value,omitempty"`
}

func invokeActionFunc(knativeClient *knative.Clientset) actions.InvokeActionHandlerFunc {
	return func(params actions.InvokeActionParams) middleware.Responder {
		namespace := namespaceOrDefault(params.Namespace)
		route, err := knativeClient.ServingV1alpha1().Routes(namespace).Get(params.ActionName, metav1.GetOptions{})
		if err != nil {
			errorMessage := errorMessageFromErr(err)
			if errors.IsNotFound(err) {
				return actions.NewInvokeActionNotFound().WithPayload(errorMessage)
			}
			return actions.NewInvokeActionInternalServerError().WithPayload(errorMessage)
		}

		config, err := knativeClient.ServingV1alpha1().Configurations(namespace).Get(params.ActionName, metav1.GetOptions{})
		if err != nil {
			errorMessage := errorMessageFromErr(err)
			if errors.IsNotFound(err) {
				return actions.NewInvokeActionNotFound().WithPayload(errorMessage)
			}
			return actions.NewInvokeActionInternalServerError().WithPayload(errorMessage)
		}
		annotations := config.Annotations

		actionHost := route.Status.Domain

		// If we're running in-cluster this needs to be an internal
		// hostname. If we're running outside the cluster, this needs
		// to be the exposed route and/or nodeport. For now, don't
		// worry about magic and expect it to be explicitly configured
		// via a flag.
		//
		// host := "istio-ingress.istio-system.svc.cluster.local"
		istioHostAndPort := kwskFlags.Istio
		if istioHostAndPort == "" {
			panic("Istio host and port must be provided via --istio flag to invoke actions")
		}

		// TODO: Don't init the action every time it's invoked
		errResponder := initAction(istioHostAndPort, actionHost, annotations["kwsk_action_code"])
		if errResponder != nil {
			return errResponder
		}
		return runAction(istioHostAndPort, actionHost, config.Name, namespace, params.Payload)
	}
}

func initAction(istioHostAndPort string, actionHost string, actionCode string) middleware.Responder {
	initBody := &ActionInitMessage{
		Value: ActionInitValue{
			Main: "main",
			Code: actionCode,
		},
	}
	resStatus, resBody, err := actionRequest(istioHostAndPort, actionHost, "init", initBody)
	if err != nil {
		return actions.NewInvokeActionInternalServerError().WithPayload(errorMessageFromErr(err))
	}

	if resStatus != http.StatusOK {
		msg := fmt.Sprintf("Error initializating action. Status: %d, Message: %s\n", resStatus, resBody)
		errorMessage := &models.ErrorMessage{
			Error: &msg,
		}
		return actions.NewInvokeActionInternalServerError().WithPayload(errorMessage)
	}

	return nil
}

func runAction(istioHostAndPort string, actionHost string, name string, namespace string, payload interface{}) middleware.Responder {

	runBody := &ActionRunMessage{
		Value: payload,
	}
	resStatus, resBody, err := actionRequest(istioHostAndPort, actionHost, "run", runBody)
	if err != nil {
		return actions.NewInvokeActionInternalServerError().WithPayload(errorMessageFromErr(err))
	}

	if resStatus != http.StatusOK {
		msg := fmt.Sprintf("Error invoking action. Status: %d, Message: %s\n", resStatus, resBody)
		errorMessage := &models.ErrorMessage{
			Error: &msg,
		}
		return actions.NewInvokeActionInternalServerError().WithPayload(errorMessage)
	}

	var resultJson interface{}
	err = json.Unmarshal(resBody, &resultJson)
	if err != nil {
		msg := fmt.Sprintf("Action invocation result was not valid JSON. Result: %s\n", resStatus, resBody)
		errorMessage := &models.ErrorMessage{
			Error: &msg,
		}
		return actions.NewInvokeActionInternalServerError().WithPayload(errorMessage)
	}
	activationResult := &models.ActivationResult{
		Value: resultJson,
	}

	activationId := "dummyactivationid"
	logs := ""
	activation := &models.Activation{
		ActivationID: &activationId,
		Name:         &name,
		Namespace:    &namespace,
		Result:       activationResult,
		Logs:         &logs,
	}
	return actions.NewInvokeActionOK().WithPayload(activation)
}

func actionRequest(istioHostAndPort string, actionHost string, path string, requestBody interface{}) (int, []byte, error) {
	url := fmt.Sprintf("http://%s/%s", istioHostAndPort, path)
	fmt.Printf("Sending POST to url %s\n", url)

	body, err := json.Marshal(requestBody)
	if err != nil {
		return 500, nil, err
	}
	fmt.Printf("Request Body: %s\n", body)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return 500, nil, err
	}

	req.Host = actionHost
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 500, nil, err
	}

	defer res.Body.Close()
	resBody, _ := ioutil.ReadAll(res.Body)
	fmt.Printf("Response Body: %s\n", string(resBody))

	return res.StatusCode, resBody, nil
}

func namespaceOrDefault(namespace string) string {
	if namespace == "_" {
		// TODO: In OpenWhisk land, the "_" namespace means the
		// default namespace of the authenticated user. Because we're
		// not dealing with any auth yet, just hardcode this to the
		// "default" namespace for now.
		namespace = "default"
	}
	return namespace
}

func errorMessageFromErr(err error) *models.ErrorMessage {
	msg := err.Error()
	return &models.ErrorMessage{
		Error: &msg,
	}
}
