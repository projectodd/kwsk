package restapi

import (
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/projectodd/kwsk/restapi/operations"
	"github.com/projectodd/kwsk/restapi/operations/rules"

	knative "github.com/knative/serving/pkg/client/clientset/versioned"
)

func configureRules(api *operations.KwskAPI, knativeClient *knative.Clientset) {
	api.RulesDeleteRuleHandler = rules.DeleteRuleHandlerFunc(func(params rules.DeleteRuleParams) middleware.Responder {
		return middleware.NotImplemented("operation rules.DeleteRule has not yet been implemented")
	})
	api.RulesGetAllRulesHandler = rules.GetAllRulesHandlerFunc(func(params rules.GetAllRulesParams) middleware.Responder {
		return rules.NewGetAllRulesOK()
	})
	api.RulesGetRuleByNameHandler = rules.GetRuleByNameHandlerFunc(func(params rules.GetRuleByNameParams) middleware.Responder {
		return middleware.NotImplemented("operation rules.GetRuleByName has not yet been implemented")
	})
	api.RulesSetStateHandler = rules.SetStateHandlerFunc(func(params rules.SetStateParams) middleware.Responder {
		return middleware.NotImplemented("operation rules.SetState has not yet been implemented")
	})
	api.RulesUpdateRuleHandler = rules.UpdateRuleHandlerFunc(func(params rules.UpdateRuleParams) middleware.Responder {
		return middleware.NotImplemented("operation rules.UpdateRule has not yet been implemented")
	})
}
