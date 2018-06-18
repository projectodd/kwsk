// Code generated by go-swagger; DO NOT EDIT.

package triggers

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/projectodd/kwsk/models"
)

// UpdateTriggerOKCode is the HTTP code returned for type UpdateTriggerOK
const UpdateTriggerOKCode int = 200

/*UpdateTriggerOK Updated Item

swagger:response updateTriggerOK
*/
type UpdateTriggerOK struct {

	/*
	  In: Body
	*/
	Payload *models.ItemID `json:"body,omitempty"`
}

// NewUpdateTriggerOK creates UpdateTriggerOK with default headers values
func NewUpdateTriggerOK() *UpdateTriggerOK {

	return &UpdateTriggerOK{}
}

// WithPayload adds the payload to the update trigger o k response
func (o *UpdateTriggerOK) WithPayload(payload *models.ItemID) *UpdateTriggerOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update trigger o k response
func (o *UpdateTriggerOK) SetPayload(payload *models.ItemID) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateTriggerOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateTriggerBadRequestCode is the HTTP code returned for type UpdateTriggerBadRequest
const UpdateTriggerBadRequestCode int = 400

/*UpdateTriggerBadRequest Bad request

swagger:response updateTriggerBadRequest
*/
type UpdateTriggerBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorMessage `json:"body,omitempty"`
}

// NewUpdateTriggerBadRequest creates UpdateTriggerBadRequest with default headers values
func NewUpdateTriggerBadRequest() *UpdateTriggerBadRequest {

	return &UpdateTriggerBadRequest{}
}

// WithPayload adds the payload to the update trigger bad request response
func (o *UpdateTriggerBadRequest) WithPayload(payload *models.ErrorMessage) *UpdateTriggerBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update trigger bad request response
func (o *UpdateTriggerBadRequest) SetPayload(payload *models.ErrorMessage) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateTriggerBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateTriggerUnauthorizedCode is the HTTP code returned for type UpdateTriggerUnauthorized
const UpdateTriggerUnauthorizedCode int = 401

/*UpdateTriggerUnauthorized Unauthorized request

swagger:response updateTriggerUnauthorized
*/
type UpdateTriggerUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorMessage `json:"body,omitempty"`
}

// NewUpdateTriggerUnauthorized creates UpdateTriggerUnauthorized with default headers values
func NewUpdateTriggerUnauthorized() *UpdateTriggerUnauthorized {

	return &UpdateTriggerUnauthorized{}
}

// WithPayload adds the payload to the update trigger unauthorized response
func (o *UpdateTriggerUnauthorized) WithPayload(payload *models.ErrorMessage) *UpdateTriggerUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update trigger unauthorized response
func (o *UpdateTriggerUnauthorized) SetPayload(payload *models.ErrorMessage) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateTriggerUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateTriggerConflictCode is the HTTP code returned for type UpdateTriggerConflict
const UpdateTriggerConflictCode int = 409

/*UpdateTriggerConflict Conflicting item already exists

swagger:response updateTriggerConflict
*/
type UpdateTriggerConflict struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorMessage `json:"body,omitempty"`
}

// NewUpdateTriggerConflict creates UpdateTriggerConflict with default headers values
func NewUpdateTriggerConflict() *UpdateTriggerConflict {

	return &UpdateTriggerConflict{}
}

// WithPayload adds the payload to the update trigger conflict response
func (o *UpdateTriggerConflict) WithPayload(payload *models.ErrorMessage) *UpdateTriggerConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update trigger conflict response
func (o *UpdateTriggerConflict) SetPayload(payload *models.ErrorMessage) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateTriggerConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(409)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateTriggerRequestEntityTooLargeCode is the HTTP code returned for type UpdateTriggerRequestEntityTooLarge
const UpdateTriggerRequestEntityTooLargeCode int = 413

/*UpdateTriggerRequestEntityTooLarge Request entity too large

swagger:response updateTriggerRequestEntityTooLarge
*/
type UpdateTriggerRequestEntityTooLarge struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorMessage `json:"body,omitempty"`
}

// NewUpdateTriggerRequestEntityTooLarge creates UpdateTriggerRequestEntityTooLarge with default headers values
func NewUpdateTriggerRequestEntityTooLarge() *UpdateTriggerRequestEntityTooLarge {

	return &UpdateTriggerRequestEntityTooLarge{}
}

// WithPayload adds the payload to the update trigger request entity too large response
func (o *UpdateTriggerRequestEntityTooLarge) WithPayload(payload *models.ErrorMessage) *UpdateTriggerRequestEntityTooLarge {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update trigger request entity too large response
func (o *UpdateTriggerRequestEntityTooLarge) SetPayload(payload *models.ErrorMessage) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateTriggerRequestEntityTooLarge) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(413)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateTriggerInternalServerErrorCode is the HTTP code returned for type UpdateTriggerInternalServerError
const UpdateTriggerInternalServerErrorCode int = 500

/*UpdateTriggerInternalServerError Server error

swagger:response updateTriggerInternalServerError
*/
type UpdateTriggerInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorMessage `json:"body,omitempty"`
}

// NewUpdateTriggerInternalServerError creates UpdateTriggerInternalServerError with default headers values
func NewUpdateTriggerInternalServerError() *UpdateTriggerInternalServerError {

	return &UpdateTriggerInternalServerError{}
}

// WithPayload adds the payload to the update trigger internal server error response
func (o *UpdateTriggerInternalServerError) WithPayload(payload *models.ErrorMessage) *UpdateTriggerInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update trigger internal server error response
func (o *UpdateTriggerInternalServerError) SetPayload(payload *models.ErrorMessage) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateTriggerInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
