// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Activation activation
// swagger:model Activation
type Activation struct {

	// Id of the activation
	// Required: true
	ActivationID *string `json:"activationId"`

	// Time when the activation completed
	// Required: true
	End *string `json:"end"`

	// Logs generated by the activation
	// Required: true
	Logs *string `json:"logs"`

	// Name of the item
	// Required: true
	Name *string `json:"name"`

	// Namespace of the associated item
	// Required: true
	Namespace *string `json:"namespace"`

	// Whether to publish the item or not
	// Required: true
	Publish *bool `json:"publish"`

	// result
	// Required: true
	Result *ActivationResult `json:"result"`

	// Time when the activation began
	// Required: true
	Start *string `json:"start"`

	// The subject that activated the item
	// Required: true
	Subject *string `json:"subject"`

	// Semantic version of the item
	// Required: true
	Version *string `json:"version"`
}

// Validate validates this activation
func (m *Activation) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateActivationID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateEnd(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLogs(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNamespace(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePublish(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateResult(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStart(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSubject(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateVersion(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Activation) validateActivationID(formats strfmt.Registry) error {

	if err := validate.Required("activationId", "body", m.ActivationID); err != nil {
		return err
	}

	return nil
}

func (m *Activation) validateEnd(formats strfmt.Registry) error {

	if err := validate.Required("end", "body", m.End); err != nil {
		return err
	}

	return nil
}

func (m *Activation) validateLogs(formats strfmt.Registry) error {

	if err := validate.Required("logs", "body", m.Logs); err != nil {
		return err
	}

	return nil
}

func (m *Activation) validateName(formats strfmt.Registry) error {

	if err := validate.Required("name", "body", m.Name); err != nil {
		return err
	}

	return nil
}

func (m *Activation) validateNamespace(formats strfmt.Registry) error {

	if err := validate.Required("namespace", "body", m.Namespace); err != nil {
		return err
	}

	return nil
}

func (m *Activation) validatePublish(formats strfmt.Registry) error {

	if err := validate.Required("publish", "body", m.Publish); err != nil {
		return err
	}

	return nil
}

func (m *Activation) validateResult(formats strfmt.Registry) error {

	if err := validate.Required("result", "body", m.Result); err != nil {
		return err
	}

	if m.Result != nil {
		if err := m.Result.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("result")
			}
			return err
		}
	}

	return nil
}

func (m *Activation) validateStart(formats strfmt.Registry) error {

	if err := validate.Required("start", "body", m.Start); err != nil {
		return err
	}

	return nil
}

func (m *Activation) validateSubject(formats strfmt.Registry) error {

	if err := validate.Required("subject", "body", m.Subject); err != nil {
		return err
	}

	return nil
}

func (m *Activation) validateVersion(formats strfmt.Registry) error {

	if err := validate.Required("version", "body", m.Version); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Activation) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Activation) UnmarshalBinary(b []byte) error {
	var res Activation
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
