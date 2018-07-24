package restapi

import (
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/projectodd/kwsk/models"
	"github.com/projectodd/kwsk/restapi/operations"
	"github.com/projectodd/kwsk/restapi/operations/triggers"

	knative "github.com/knative/serving/pkg/client/clientset/versioned"
)

func configureTriggers(api *operations.KwskAPI, knativeClient *knative.Clientset) {
	api.TriggersDeleteTriggerHandler = triggers.DeleteTriggerHandlerFunc(func(params triggers.DeleteTriggerParams, principal *models.Principal) middleware.Responder {
		return middleware.NotImplemented("operation triggers.DeleteTrigger has not yet been implemented")
	})
	api.TriggersFireTriggerHandler = triggers.FireTriggerHandlerFunc(func(params triggers.FireTriggerParams, principal *models.Principal) middleware.Responder {
		return middleware.NotImplemented("operation triggers.FireTrigger has not yet been implemented")
	})
	api.TriggersGetAllTriggersHandler = triggers.GetAllTriggersHandlerFunc(func(params triggers.GetAllTriggersParams, principal *models.Principal) middleware.Responder {
		return triggers.NewGetAllTriggersOK()
	})
	api.TriggersGetTriggerByNameHandler = triggers.GetTriggerByNameHandlerFunc(func(params triggers.GetTriggerByNameParams, principal *models.Principal) middleware.Responder {
		return middleware.NotImplemented("operation triggers.GetTriggerByName has not yet been implemented")
	})
	api.TriggersUpdateTriggerHandler = triggers.UpdateTriggerHandlerFunc(func(params triggers.UpdateTriggerParams, principal *models.Principal) middleware.Responder {
		return middleware.NotImplemented("operation triggers.UpdateTrigger has not yet been implemented")
	})
}
