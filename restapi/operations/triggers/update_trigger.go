// Code generated by go-swagger; DO NOT EDIT.

package triggers

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// UpdateTriggerHandlerFunc turns a function with the right signature into a update trigger handler
type UpdateTriggerHandlerFunc func(UpdateTriggerParams) middleware.Responder

// Handle executing the request and returning a response
func (fn UpdateTriggerHandlerFunc) Handle(params UpdateTriggerParams) middleware.Responder {
	return fn(params)
}

// UpdateTriggerHandler interface for that can handle valid update trigger params
type UpdateTriggerHandler interface {
	Handle(UpdateTriggerParams) middleware.Responder
}

// NewUpdateTrigger creates a new http.Handler for the update trigger operation
func NewUpdateTrigger(ctx *middleware.Context, handler UpdateTriggerHandler) *UpdateTrigger {
	return &UpdateTrigger{Context: ctx, Handler: handler}
}

/*UpdateTrigger swagger:route PUT /namespaces/{namespace}/triggers/{triggerName} Triggers updateTrigger

Update a trigger

Update a trigger

*/
type UpdateTrigger struct {
	Context *middleware.Context
	Handler UpdateTriggerHandler
}

func (o *UpdateTrigger) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewUpdateTriggerParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
