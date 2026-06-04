package domain

import (
	"errors"
	"fmt"
	"slices"

	"sundance/backend/pkg/common/validate"
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
	ErrInvalidFieldType         = errors.New("invalid field type")
	ErrInvalidFieldAttributes   = errors.New("invalid field attributes for type")
	ErrDuplicateFieldTagMapping = errors.New("duplicate field tag mapping for same tag version")
)

type Field struct {
	ID         FieldID
	Key        string `validate:"required,nowhitespace"`
	Name       string `validate:"required"`
	Type       FieldType
	Attributes FieldAttributes
	tags       []*FieldTagMapping
	withPosition
	withRules
}

func NewField(key, name string, fieldType FieldType, attributes FieldAttributes, position float32) (*Field, error) {
	if !isValidPosition(position) {
		return nil, ErrInvalidPosition
	}

	if !isValidFieldType(fieldType) {
		return nil, ErrInvalidFieldType
	}

	if !isValidFieldAttributes(fieldType, attributes) {
		return nil, ErrInvalidFieldAttributes
	}

	f := &Field{
		ID:         FieldID(NewID()),
		Key:        key,
		Name:       name,
		Type:       fieldType,
		Attributes: attributes,
		withPosition: withPosition{
			position: position,
		},
	}

	if err := validate.ValidateStruct(f); err != nil {
		return nil, err
	}

	return f, nil
}

func HydrateField(
	id FieldID,
	key,
	name string,
	fieldType FieldType,
	attr FieldAttributes,
	position float32,
	tags []*FieldTagMapping,
) *Field {
	return &Field{
		ID:         id,
		Key:        key,
		Name:       name,
		Type:       fieldType,
		Attributes: attr,
		tags:       tags,
		withPosition: withPosition{
			position: position,
		},
	}
}

func (f Field) GetTags() []*FieldTagMapping {
	return f.tags
}

func (f *Field) AddTags(tags ...FieldTagMappingConfig) error {
	for _, tag := range tags {
		idx := slices.IndexFunc(f.tags, func(ftm *FieldTagMapping) bool {
			return ftm.TagVersionID == tag.TagVersionID
		})

		if idx != -1 {
			return fmt.Errorf("%w: tagVersion=%s", ErrDuplicateFieldTagMapping, tag.TagVersionID)
		}

		ftm, err := NewFieldTagMapping(f.ID, tag.TagVersionID, tag.Priority)
		if err != nil {
			return err
		}

		f.tags = append(f.tags, ftm)
	}
	return nil
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
	case *TextFieldAttributes:
		return fieldType == FieldTypeText
	case *NumberFieldAttributes:
		return fieldType == FieldTypeNumber
	case *SelectFieldAttributes:
		return fieldType == FieldTypeSelect
	case *CheckboxFieldAttributes:
		return fieldType == FieldTypeCheckbox
	case *DateFieldAttributes:
		return fieldType == FieldTypeDate
	default:
		return false
	}
}
