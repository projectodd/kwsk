// Code generated by go-swagger; DO NOT EDIT.

package triggers

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	models "github.com/projectodd/kwsk/models"
)

// FireTriggerHandlerFunc turns a function with the right signature into a fire trigger handler
type FireTriggerHandlerFunc func(FireTriggerParams, *models.Principal) middleware.Responder

// Handle executing the request and returning a response
func (fn FireTriggerHandlerFunc) Handle(params FireTriggerParams, principal *models.Principal) middleware.Responder {
	return fn(params, principal)
}

// FireTriggerHandler interface for that can handle valid fire trigger params
type FireTriggerHandler interface {
	Handle(FireTriggerParams, *models.Principal) middleware.Responder
}

// NewFireTrigger creates a new http.Handler for the fire trigger operation
func NewFireTrigger(ctx *middleware.Context, handler FireTriggerHandler) *FireTrigger {
	return &FireTrigger{Context: ctx, Handler: handler}
}

/*FireTrigger swagger:route POST /namespaces/{namespace}/triggers/{triggerName} Triggers fireTrigger

Fire a trigger

Fire a trigger

*/
type FireTrigger struct {
	Context *middleware.Context
	Handler FireTriggerHandler
}

func (o *FireTrigger) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewFireTriggerParams()

	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal *models.Principal
	if uprinc != nil {
		principal = uprinc.(*models.Principal) // this is really a models.Principal, I promise
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
