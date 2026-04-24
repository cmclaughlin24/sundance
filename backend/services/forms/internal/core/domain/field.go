package domain

import (
	"errors"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/validate"
	"github.com/google/uuid"
)

type FieldID string

type FieldType string

const (
	FieldTypeText     FieldType = "text"
	FieldTypeNumber   FieldType = "number"
	FieldTypeSelect   FieldType = "select"
	FieldTypeCheckbox FieldType = "checkbox"
	FieldTypeDate     FieldType = "date"
)

var (
	ErrInvalidFieldType       = errors.New("invalid field type")
	ErrInvalidFieldAttributes = errors.New("invalid field attributes for type")
)

type Field struct {
	ID         FieldID
	Key        string
	Name       string
	Type       FieldType
	Attributes FieldAttributes
	Position   int
	baseWithRules
}

func NewField(key, name string, fieldType FieldType, attributes FieldAttributes, position int) (*Field, error) {
	if !isValidFieldType(fieldType) {
		return nil, ErrInvalidFieldType
	}

	if !isValidFieldAttributes(fieldType, attributes) {
		return nil, ErrInvalidFieldAttributes
	}

	return &Field{
		ID:         FieldID(uuid.NewString()),
		Key:        key,
		Name:       name,
		Type:       fieldType,
		Attributes: attributes,
		Position:   position,
	}, nil
}

func HydrateField(id FieldID, key, name string, fieldType FieldType, attr FieldAttributes, position int) *Field {
	return &Field{
		ID:         id,
		Key:        key,
		Name:       name,
		Type:       fieldType,
		Attributes: attr,
		Position:   position,
	}
}

var isValidFieldType = validate.NewTypeValidator([]FieldType{
	FieldTypeText,
	FieldTypeNumber,
	FieldTypeSelect,
	FieldTypeCheckbox,
	FieldTypeDate,
})

func isValidFieldAttributes(fieldType FieldType, attr FieldAttributes) bool {
	switch attr.(type) {
	case TextFieldAttributes:
		return fieldType == FieldTypeText
	case NumberFieldAttributes:
		return fieldType == FieldTypeNumber
	case SelectFieldAttributes:
		return fieldType == FieldTypeSelect
	case CheckboxFieldAttributes:
		return fieldType == FieldTypeCheckbox
	case DateFieldAttributes:
		return fieldType == FieldTypeDate
	default:
		return false
	}
}
