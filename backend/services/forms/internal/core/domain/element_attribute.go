package domain

import (
	"errors"
)

var (
	ErrElementAttributeMismatch = errors.New("element type and attributes mismatch")
)

type ElementAttributes interface {
	GetIsRequired() bool
	SetIsRequired(bool)
	GetIsReadOnly() bool
	GetReferencedKeys() []string
}

type BaseElementAttributes struct {
	IsReadOnly bool `desc:"Whether the field is disabled for user editing"`
	IsRequired bool `desc:"Whether the field must be populated before submitting"`
}

func (a BaseElementAttributes) GetIsRequired() bool {
	return a.IsRequired
}

func (a BaseElementAttributes) GetIsReadOnly() bool {
	return a.IsReadOnly
}

func (a *BaseElementAttributes) GetReferencedKeys() []string {
	return make([]string, 0)
}

func (a *BaseElementAttributes) SetIsRequired(required bool) {
	a.IsRequired = required
}

type TextElementAttributes struct {
	BaseElementAttributes
	MinLength   *int   `desc:"Minimum allowed character count"`
	MaxLength   *int   `desc:"Maximum allowed character count"`
	Pattern     string `desc:"Regex pattern the text value must match"`
	Placeholder string `desc:"Placeholder text displayed in the input"`
}

type NumberElementAttributes struct {
	BaseElementAttributes
	Min  *float64 `desc:"Minimum numeric value permitted"`
	Max  *float64 `desc:"Maximum numeric value permitted"`
	Step *float64 `desc:"Permitted step interval between numeric values"`
}

type SelectElementAttributes struct {
	BaseElementAttributes
	Data          []any          `desc:"Static list of key-value lookup items"`
	DataSourceRef *DataSourceRef `desc:"Reference to an external data source in the Tenants service"`
	Multiple      bool           `desc:"Whether multiple selections are permitted"`
	MinSelected   *int           `desc:"Minimum selection count when multiple is true"`
	MaxSelected   *int           `desc:"Maximum selection count when multiple is true"`
}

func (a *SelectElementAttributes) GetReferencedKeys() []string {
	return getReferencedKeys(a.DataSourceRef)
}

type CheckboxElementAttributes struct {
	BaseElementAttributes
	IsCheckedByDefault bool           `desc:"Initial checked state"`
	Data               []any          `desc:"Static list of key-value lookup items"`
	DataSourceRef      *DataSourceRef `desc:"Reference to an external data source in the Tenants service"`
}

type DateElementAttributes struct {
	BaseElementAttributes
	MinDate *string `desc:"Earliest selectable date in ISO 8601 format"`
	MaxDate *string `desc:"Latest selectable date in ISO 8601 format"`
}

func GetElementAttributes[T ElementAttributes](attr ElementAttributes) (T, error) {
	switch t := attr.(type) {
	case T:
		return t, nil
	default:
		return *new(T), ErrElementAttributeMismatch
	}
}
