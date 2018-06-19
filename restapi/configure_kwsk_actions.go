package restapi

import (
	"fmt"

	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/projectodd/kwsk/restapi/operations"
	"github.com/projectodd/kwsk/restapi/operations/actions"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	knative "github.com/knative/serving/pkg/client/clientset/versioned"
)

func configureActions(api *operations.KwskAPI, knativeClient *knative.Clientset) {
	api.ActionsDeleteActionHandler = actions.DeleteActionHandlerFunc(func(params actions.DeleteActionParams) middleware.Responder {
		return middleware.NotImplemented("operation actions.DeleteAction has not yet been implemented")
	})

	api.ActionsGetActionByNameHandler = actions.GetActionByNameHandlerFunc(func(params actions.GetActionByNameParams) middleware.Responder {
		return actions.NewGetActionByNameOK()
	})

	api.ActionsGetAllActionsHandler = actions.GetAllActionsHandlerFunc(actionsGetAllActionsFunc(knativeClient))

	api.ActionsInvokeActionHandler = actions.InvokeActionHandlerFunc(func(params actions.InvokeActionParams) middleware.Responder {
		return middleware.NotImplemented("operation actions.InvokeAction has not yet been implemented")
	})

	api.ActionsUpdateActionHandler = actions.UpdateActionHandlerFunc(func(params actions.UpdateActionParams) middleware.Responder {
		return actions.NewUpdateActionOK()
	})
}

func actionsGetAllActionsFunc(knativeClient *knative.Clientset) actions.GetAllActionsHandlerFunc {
	return func(params actions.GetAllActionsParams) middleware.Responder {
		// TODO: This is just stubbed in here to show an example of
		// fetching knative CRDs. The namespace should definitely not
		// be hardcoded, for example.
		configs, err := knativeClient.ServingV1alpha1().Configurations("default").List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d configs in the cluster\n", len(configs.Items))
		return actions.NewGetAllActionsOK()
	}
}
