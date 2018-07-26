package restapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/projectodd/kwsk/models"
	"github.com/projectodd/kwsk/restapi/operations"
	"github.com/projectodd/kwsk/restapi/operations/activations"
	"k8s.io/client-go/tools/cache"

	knative "github.com/knative/serving/pkg/client/clientset/versioned"
)

func configureActivations(api *operations.KwskAPI, knativeClient *knative.Clientset, cache cache.Store) {
	api.ActivationsGetNamespacesNamespaceActivationsActivationidLogsHandler = activations.GetNamespacesNamespaceActivationsActivationidLogsHandlerFunc(func(params activations.GetNamespacesNamespaceActivationsActivationidLogsParams, principal *models.Principal) middleware.Responder {
		return middleware.NotImplemented("operation activations.GetNamespacesNamespaceActivationsActivationidLogs has not yet been implemented")
	})
	api.ActivationsGetNamespacesNamespaceActivationsActivationidResultHandler = activations.GetNamespacesNamespaceActivationsActivationidResultHandlerFunc(func(params activations.GetNamespacesNamespaceActivationsActivationidResultParams, principal *models.Principal) middleware.Responder {
		return middleware.NotImplemented("operation activations.GetNamespacesNamespaceActivationsActivationidResult has not yet been implemented")
	})

	api.ActivationsGetActivationByIDHandler = activations.GetActivationByIDHandlerFunc(getActivationByIDFunc(cache))

	api.ActivationsGetActivationsHandler = activations.GetActivationsHandlerFunc(getActivationsFunc(cache))
}

func getActivationByIDFunc(cache cache.Store) activations.GetActivationByIDHandlerFunc {
	return func(params activations.GetActivationByIDParams, principal *models.Principal) middleware.Responder {
		activationId := params.Activationid
		obj, exists, err := cache.GetByKey(activationId)
		if err != nil {
			msg := fmt.Sprintf("Error retrieving activation record: %s\n", err)
			errorMessage := &models.ErrorMessage{
				Error: &msg,
			}
			return activations.NewGetActivationByIDInternalServerError().WithPayload(errorMessage)
		}
		if !exists {
			return activations.NewGetActivationByIDNotFound()
		}
		if activation, ok := obj.(models.Activation); ok {
			fmt.Printf("Got Activation: %+v\n", activation)
			return activations.NewGetActivationByIDOK().WithPayload(&activation)
		}
		return activations.NewGetActivationByIDInternalServerError()
	}
}

func getActivationsFunc(cache cache.Store) activations.GetActivationsHandlerFunc {
	return func(params activations.GetActivationsParams, principal *models.Principal) middleware.Responder {
		objs := cache.List()
		var payload = make([]*models.Activation, len(objs))
		for i, obj := range objs {
			if activation, ok := obj.(models.Activation); ok {
				payload[i] = &activation
			}
		}
		return activations.NewGetActivationsOK().WithPayload(payload)
	}
}
