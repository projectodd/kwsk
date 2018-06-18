// Code generated by go-swagger; DO NOT EDIT.

package packages

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/projectodd/kwsk/models"
)

// UpdatePackageOKCode is the HTTP code returned for type UpdatePackageOK
const UpdatePackageOKCode int = 200

/*UpdatePackageOK Updated Item

swagger:response updatePackageOK
*/
type UpdatePackageOK struct {

	/*
	  In: Body
	*/
	Payload *models.ItemID `json:"body,omitempty"`
}

// NewUpdatePackageOK creates UpdatePackageOK with default headers values
func NewUpdatePackageOK() *UpdatePackageOK {

	return &UpdatePackageOK{}
}

// WithPayload adds the payload to the update package o k response
func (o *UpdatePackageOK) WithPayload(payload *models.ItemID) *UpdatePackageOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update package o k response
func (o *UpdatePackageOK) SetPayload(payload *models.ItemID) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdatePackageOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdatePackageBadRequestCode is the HTTP code returned for type UpdatePackageBadRequest
const UpdatePackageBadRequestCode int = 400

/*UpdatePackageBadRequest Bad request

swagger:response updatePackageBadRequest
*/
type UpdatePackageBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorMessage `json:"body,omitempty"`
}

// NewUpdatePackageBadRequest creates UpdatePackageBadRequest with default headers values
func NewUpdatePackageBadRequest() *UpdatePackageBadRequest {

	return &UpdatePackageBadRequest{}
}

// WithPayload adds the payload to the update package bad request response
func (o *UpdatePackageBadRequest) WithPayload(payload *models.ErrorMessage) *UpdatePackageBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update package bad request response
func (o *UpdatePackageBadRequest) SetPayload(payload *models.ErrorMessage) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdatePackageBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdatePackageUnauthorizedCode is the HTTP code returned for type UpdatePackageUnauthorized
const UpdatePackageUnauthorizedCode int = 401

/*UpdatePackageUnauthorized Unauthorized request

swagger:response updatePackageUnauthorized
*/
type UpdatePackageUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorMessage `json:"body,omitempty"`
}

// NewUpdatePackageUnauthorized creates UpdatePackageUnauthorized with default headers values
func NewUpdatePackageUnauthorized() *UpdatePackageUnauthorized {

	return &UpdatePackageUnauthorized{}
}

// WithPayload adds the payload to the update package unauthorized response
func (o *UpdatePackageUnauthorized) WithPayload(payload *models.ErrorMessage) *UpdatePackageUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update package unauthorized response
func (o *UpdatePackageUnauthorized) SetPayload(payload *models.ErrorMessage) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdatePackageUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdatePackageConflictCode is the HTTP code returned for type UpdatePackageConflict
const UpdatePackageConflictCode int = 409

/*UpdatePackageConflict Conflicting item already exists

swagger:response updatePackageConflict
*/
type UpdatePackageConflict struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorMessage `json:"body,omitempty"`
}

// NewUpdatePackageConflict creates UpdatePackageConflict with default headers values
func NewUpdatePackageConflict() *UpdatePackageConflict {

	return &UpdatePackageConflict{}
}

// WithPayload adds the payload to the update package conflict response
func (o *UpdatePackageConflict) WithPayload(payload *models.ErrorMessage) *UpdatePackageConflict {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update package conflict response
func (o *UpdatePackageConflict) SetPayload(payload *models.ErrorMessage) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdatePackageConflict) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(409)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdatePackageRequestEntityTooLargeCode is the HTTP code returned for type UpdatePackageRequestEntityTooLarge
const UpdatePackageRequestEntityTooLargeCode int = 413

/*UpdatePackageRequestEntityTooLarge Request entity too large

swagger:response updatePackageRequestEntityTooLarge
*/
type UpdatePackageRequestEntityTooLarge struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorMessage `json:"body,omitempty"`
}

// NewUpdatePackageRequestEntityTooLarge creates UpdatePackageRequestEntityTooLarge with default headers values
func NewUpdatePackageRequestEntityTooLarge() *UpdatePackageRequestEntityTooLarge {

	return &UpdatePackageRequestEntityTooLarge{}
}

// WithPayload adds the payload to the update package request entity too large response
func (o *UpdatePackageRequestEntityTooLarge) WithPayload(payload *models.ErrorMessage) *UpdatePackageRequestEntityTooLarge {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update package request entity too large response
func (o *UpdatePackageRequestEntityTooLarge) SetPayload(payload *models.ErrorMessage) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdatePackageRequestEntityTooLarge) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(413)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdatePackageInternalServerErrorCode is the HTTP code returned for type UpdatePackageInternalServerError
const UpdatePackageInternalServerErrorCode int = 500

/*UpdatePackageInternalServerError Server error

swagger:response updatePackageInternalServerError
*/
type UpdatePackageInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorMessage `json:"body,omitempty"`
}

// NewUpdatePackageInternalServerError creates UpdatePackageInternalServerError with default headers values
func NewUpdatePackageInternalServerError() *UpdatePackageInternalServerError {

	return &UpdatePackageInternalServerError{}
}

// WithPayload adds the payload to the update package internal server error response
func (o *UpdatePackageInternalServerError) WithPayload(payload *models.ErrorMessage) *UpdatePackageInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update package internal server error response
func (o *UpdatePackageInternalServerError) SetPayload(payload *models.ErrorMessage) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdatePackageInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
