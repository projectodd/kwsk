package restapi

import (
	"fmt"

	middleware "github.com/go-openapi/runtime/middleware"

	models "github.com/projectodd/kwsk/models"
	"github.com/projectodd/kwsk/restapi/operations"
	"github.com/projectodd/kwsk/restapi/operations/actions"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/knative/serving/pkg/apis/serving/v1alpha1"
	knative "github.com/knative/serving/pkg/client/clientset/versioned"
)

func configureActions(api *operations.KwskAPI, knativeClient *knative.Clientset) {
	api.ActionsDeleteActionHandler = actions.DeleteActionHandlerFunc(func(params actions.DeleteActionParams) middleware.Responder {
		return middleware.NotImplemented("operation actions.DeleteAction has not yet been implemented")
	})

	api.ActionsGetActionByNameHandler = actions.GetActionByNameHandlerFunc(func(params actions.GetActionByNameParams) middleware.Responder {
		return actions.NewGetActionByNameOK()
	})

	api.ActionsGetAllActionsHandler = actions.GetAllActionsHandlerFunc(getAllActionsFunc(knativeClient))

	api.ActionsInvokeActionHandler = actions.InvokeActionHandlerFunc(func(params actions.InvokeActionParams) middleware.Responder {
		return middleware.NotImplemented("operation actions.InvokeAction has not yet been implemented")
	})

	api.ActionsUpdateActionHandler = actions.UpdateActionHandlerFunc(updateActionFunc(knativeClient))
}

func updateActionFunc(knativeClient *knative.Clientset) actions.UpdateActionHandlerFunc {
	return func(params actions.UpdateActionParams) middleware.Responder {
		namespace := namespaceOrDefault(params.Namespace)
		container := corev1.Container{
			Name: "action",
		}
		config := &v1alpha1.Configuration{
			ObjectMeta: metav1.ObjectMeta{
				Name:      params.ActionName,
				Namespace: namespace,
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
		_, err := knativeClient.ServingV1alpha1().Configurations(namespace).Create(config)
		if err != nil {
			msg := err.Error()
			errorMessage := &models.ErrorMessage{
				Error: &msg,
			}
			return actions.NewUpdateActionInternalServerError().WithPayload(errorMessage)
		}
		return actions.NewUpdateActionOK()
	}
}

func getAllActionsFunc(knativeClient *knative.Clientset) actions.GetAllActionsHandlerFunc {
	return func(params actions.GetAllActionsParams) middleware.Responder {
		namespace := namespaceOrDefault(params.Namespace)
		configs, err := knativeClient.ServingV1alpha1().Configurations(namespace).List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		var payload = make([]*models.EntityBrief, len(configs.Items))
		for i, item := range configs.Items {
			name := item.ObjectMeta.Name
			namespace := item.ObjectMeta.Namespace
			payload[i] = &models.EntityBrief{
				Name:      &name,
				Namespace: &namespace,
			}
		}
		return actions.NewGetAllActionsOK().WithPayload(payload)
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
