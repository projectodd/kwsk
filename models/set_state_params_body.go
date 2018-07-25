// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// SetStateParamsBody set state params body
// swagger:model setStateParamsBody
type SetStateParamsBody struct {

	// status
	// Required: true
	// Enum: [inactive active]
	Status *string `json:"status"`
}

// Validate validates this set state params body
func (m *SetStateParamsBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var setStateParamsBodyTypeStatusPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["inactive","active"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		setStateParamsBodyTypeStatusPropEnum = append(setStateParamsBodyTypeStatusPropEnum, v)
	}
}

const (

	// SetStateParamsBodyStatusInactive captures enum value "inactive"
	SetStateParamsBodyStatusInactive string = "inactive"

	// SetStateParamsBodyStatusActive captures enum value "active"
	SetStateParamsBodyStatusActive string = "active"
)

// prop value enum
func (m *SetStateParamsBody) validateStatusEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, setStateParamsBodyTypeStatusPropEnum); err != nil {
		return err
	}
	return nil
}

func (m *SetStateParamsBody) validateStatus(formats strfmt.Registry) error {

	if err := validate.Required("status", "body", m.Status); err != nil {
		return err
	}

	// value enum
	if err := m.validateStatusEnum("status", "body", *m.Status); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *SetStateParamsBody) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SetStateParamsBody) UnmarshalBinary(b []byte) error {
	var res SetStateParamsBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}