package domain

import "errors"

var (
	ErrFieldAttributeMismatch = errors.New("field type and attributes mismatch")
)

type FieldAttributes interface {
	GetIsRequired() bool
	SetIsRequired(bool)
	GetIsReadOnly() bool
}

type BaseFieldAttributes struct {
	IsReadOnly bool
	IsRequired bool
}

func (a BaseFieldAttributes) GetIsRequired() bool {
	return a.IsRequired
}

func (a BaseFieldAttributes) GetIsReadOnly() bool {
	return a.IsReadOnly
}

func (a *BaseFieldAttributes) SetIsRequired(required bool) {
	a.IsRequired = required
}

type TextFieldAttributes struct {
	BaseFieldAttributes
	MinLength   *int
	MaxLength   *int
	Pattern     string
	Placeholder string
}

type NumberFieldAttributes struct {
	BaseFieldAttributes
	Min  *float64
	Max  *float64
	Step *float64
}

type SelectFieldAttributes struct {
	BaseFieldAttributes
	Multiple    bool
	MinSelected *int
	MaxSelected *int
}

type CheckboxFieldAttributes struct {
	BaseFieldAttributes
	IsCheckedByDefault bool
}

type DateFieldAttributes struct {
	BaseFieldAttributes
	MinDate *string
	MaxDate *string
}

func GetFieldAttributes[T FieldAttributes](attr FieldAttributes) (T, error) {
	switch t := attr.(type) {
	case T:
		return t, nil
	default:
		return *new(T), ErrFieldAttributeMismatch
	}
}
