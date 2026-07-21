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
	IsReadOnly bool
	IsRequired bool
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
	MinLength   *int
	MaxLength   *int
	Pattern     string
	Placeholder string
}

type NumberElementAttributes struct {
	BaseElementAttributes
	Min  *float64
	Max  *float64
	Step *float64
}

type SelectElementAttributes struct {
	BaseElementAttributes
	Data          []any
	DataSourceRef *DataSourceRef
	Multiple      bool
	MinSelected   *int
	MaxSelected   *int
}

func (a *SelectElementAttributes) GetReferencedKeys() []string {
	return getReferencedKeys(a.DataSourceRef)
}

type CheckboxElementAttributes struct {
	BaseElementAttributes
	IsCheckedByDefault bool
}

type DateElementAttributes struct {
	BaseElementAttributes
	MinDate *string
	MaxDate *string
}

func GetElementAttributes[T ElementAttributes](attr ElementAttributes) (T, error) {
	switch t := attr.(type) {
	case T:
		return t, nil
	default:
		return *new(T), ErrElementAttributeMismatch
	}
}
