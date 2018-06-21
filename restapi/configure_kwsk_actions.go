package restapi

import (
	"fmt"

	middleware "github.com/go-openapi/runtime/middleware"

	models "github.com/projectodd/kwsk/models"
	"github.com/projectodd/kwsk/restapi/operations"
	"github.com/projectodd/kwsk/restapi/operations/actions"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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
		return actions.NewDeleteActionOK()
	}
}

func updateActionFunc(knativeClient *knative.Clientset) actions.UpdateActionHandlerFunc {
	return func(params actions.UpdateActionParams) middleware.Responder {
		namespace := namespaceOrDefault(params.Namespace)

		annotations := make(map[string]string)
		annotations["kwsk_action_version"] = params.Action.Version
		annotations["kwsk_action_kind"] = *params.Action.Exec.Kind
		annotations["kwsk_action_code"] = params.Action.Exec.Code

		container := corev1.Container{
			Image: params.Action.Exec.Image,
		}
		config := &v1alpha1.Configuration{
			ObjectMeta: metav1.ObjectMeta{
				Name:        params.ActionName,
				Namespace:   namespace,
				Annotations: annotations,
			},
			Spec: v1alpha1.ConfigurationSpec{
				Generation: 0,
				RevisionTemplate: v1alpha1.RevisionTemplateSpec{
					Spec: v1alpha1.RevisionSpec{
						Container: container,
					},
				},
			},
		}
		fmt.Printf("Creating action %+v\n", config)
		createdConfig, err := knativeClient.ServingV1alpha1().Configurations(namespace).Create(config)
		if err != nil {
			msg := err.Error()
			errorMessage := &models.ErrorMessage{
				Error: &msg,
			}
			return actions.NewUpdateActionInternalServerError().WithPayload(errorMessage)
		}
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

func invokeActionFunc(knativeClient *knative.Clientset) actions.InvokeActionHandlerFunc {
	return func(params actions.InvokeActionParams) middleware.Responder {
		activationId := "fake-activation"
		activation := &models.Activation{
			ActivationID: &activationId,
		}
		return actions.NewInvokeActionOK().WithPayload(activation)
	}
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
