package restapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/projectodd/kwsk/restapi/operations"
	"github.com/projectodd/kwsk/restapi/operations/activations"

	knative "github.com/knative/serving/pkg/client/clientset/versioned"
)

func configureActivations(api *operations.KwskAPI, knativeClient *knative.Clientset) {
	api.ActivationsGetNamespacesNamespaceActivationsActivationidLogsHandler = activations.GetNamespacesNamespaceActivationsActivationidLogsHandlerFunc(func(params activations.GetNamespacesNamespaceActivationsActivationidLogsParams) middleware.Responder {
		return middleware.NotImplemented("operation activations.GetNamespacesNamespaceActivationsActivationidLogs has not yet been implemented")
	})
	api.ActivationsGetNamespacesNamespaceActivationsActivationidResultHandler = activations.GetNamespacesNamespaceActivationsActivationidResultHandlerFunc(func(params activations.GetNamespacesNamespaceActivationsActivationidResultParams) middleware.Responder {
		return middleware.NotImplemented("operation activations.GetNamespacesNamespaceActivationsActivationidResult has not yet been implemented")
	})
	api.ActivationsGetActivationByIDHandler = activations.GetActivationByIDHandlerFunc(func(params activations.GetActivationByIDParams) middleware.Responder {
		return middleware.NotImplemented("operation activations.GetActivationByID has not yet been implemented")
	})
	api.ActivationsGetActivationsHandler = activations.GetActivationsHandlerFunc(func(params activations.GetActivationsParams) middleware.Responder {
		return middleware.NotImplemented("operation activations.GetActivations has not yet been implemented")
	})
}
