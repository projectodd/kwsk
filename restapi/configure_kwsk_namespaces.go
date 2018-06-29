package restapi

import (
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/projectodd/kwsk/restapi/operations"
	"github.com/projectodd/kwsk/restapi/operations/namespaces"

	knative "github.com/knative/serving/pkg/client/clientset/versioned"
)

func configureNamespaces(api *operations.KwskAPI, knativeClient *knative.Clientset) {
	api.NamespacesGetAllNamespacesHandler = namespaces.GetAllNamespacesHandlerFunc(func(params namespaces.GetAllNamespacesParams) middleware.Responder {
		return middleware.NotImplemented("operation namespaces.GetAllNamespaces has not yet been implemented")
	})
}
