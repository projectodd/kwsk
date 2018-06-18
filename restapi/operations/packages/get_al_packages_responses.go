// Code generated by go-swagger; DO NOT EDIT.

package packages

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/projectodd/kwsk/models"
)

// GetAlPackagesOKCode is the HTTP code returned for type GetAlPackagesOK
const GetAlPackagesOKCode int = 200

/*GetAlPackagesOK Packages response

swagger:response getAlPackagesOK
*/
type GetAlPackagesOK struct {

	/*
	  In: Body
	*/
	Payload []*models.EntityBrief `json:"body,omitempty"`
}

// NewGetAlPackagesOK creates GetAlPackagesOK with default headers values
func NewGetAlPackagesOK() *GetAlPackagesOK {

	return &GetAlPackagesOK{}
}

// WithPayload adds the payload to the get al packages o k response
func (o *GetAlPackagesOK) WithPayload(payload []*models.EntityBrief) *GetAlPackagesOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get al packages o k response
func (o *GetAlPackagesOK) SetPayload(payload []*models.EntityBrief) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAlPackagesOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		payload = make([]*models.EntityBrief, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

// GetAlPackagesUnauthorizedCode is the HTTP code returned for type GetAlPackagesUnauthorized
const GetAlPackagesUnauthorizedCode int = 401

/*GetAlPackagesUnauthorized Unauthorized request

swagger:response getAlPackagesUnauthorized
*/
type GetAlPackagesUnauthorized struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorMessage `json:"body,omitempty"`
}

// NewGetAlPackagesUnauthorized creates GetAlPackagesUnauthorized with default headers values
func NewGetAlPackagesUnauthorized() *GetAlPackagesUnauthorized {

	return &GetAlPackagesUnauthorized{}
}

// WithPayload adds the payload to the get al packages unauthorized response
func (o *GetAlPackagesUnauthorized) WithPayload(payload *models.ErrorMessage) *GetAlPackagesUnauthorized {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get al packages unauthorized response
func (o *GetAlPackagesUnauthorized) SetPayload(payload *models.ErrorMessage) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAlPackagesUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(401)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetAlPackagesInternalServerErrorCode is the HTTP code returned for type GetAlPackagesInternalServerError
const GetAlPackagesInternalServerErrorCode int = 500

/*GetAlPackagesInternalServerError Server error

swagger:response getAlPackagesInternalServerError
*/
type GetAlPackagesInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorMessage `json:"body,omitempty"`
}

// NewGetAlPackagesInternalServerError creates GetAlPackagesInternalServerError with default headers values
func NewGetAlPackagesInternalServerError() *GetAlPackagesInternalServerError {

	return &GetAlPackagesInternalServerError{}
}

// WithPayload adds the payload to the get al packages internal server error response
func (o *GetAlPackagesInternalServerError) WithPayload(payload *models.ErrorMessage) *GetAlPackagesInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get al packages internal server error response
func (o *GetAlPackagesInternalServerError) SetPayload(payload *models.ErrorMessage) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetAlPackagesInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
