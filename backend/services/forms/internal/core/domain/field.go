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
	ErrInvalidField             = errors.New("invalid field")
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

func NewField(key, name string, fieldType FieldType, attr FieldAttributes, position float32) (*Field, error) {
	f := &Field{
		ID:         FieldID(NewID()),
		Key:        key,
		Name:       name,
		Type:       fieldType,
		Attributes: attr,
		withPosition: withPosition{
			position: position,
		},
	}

	if err := f.validate(); err != nil {
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

func (f *Field) Update(key, name string, fieldType FieldType, attr FieldAttributes, position float32) error {
	if f == nil {
		return ErrInvalidField
	}

	cpy := *f
	cpy.Key = key
	cpy.Name = name
	cpy.Type = fieldType
	cpy.Attributes = attr
	cpy.position = position

	if err := cpy.validate(); err != nil {
		return err
	}

	*f = cpy

	return nil
}

func (f *Field) GetTags() []*FieldTagMapping {
	return f.tags
}

func (f *Field) AddTags(mappings ...FieldTagMappingConfig) error {
	for _, tag := range mappings {
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

func (f *Field) ReplaceTags(mappings ...FieldTagMappingConfig) error {
	tags := make([]*FieldTagMapping, 0, len(mappings))

	for _, tag := range mappings {
		idx := slices.IndexFunc(tags, func(ftm *FieldTagMapping) bool {
			return ftm.TagVersionID == tag.TagVersionID
		})

		if idx != -1 {
			return fmt.Errorf("%w: tagVersion=%s", ErrDuplicateFieldTagMapping, tag.TagVersionID)
		}

		ftm, err := NewFieldTagMapping(f.ID, tag.TagVersionID, tag.Priority)
		if err != nil {
			return err
		}

		tags = append(tags, ftm)
	}

	f.tags = tags

	return nil
}

func (f *Field) validate() error {
	if !isValidPosition(f.position) {
		return ErrInvalidPosition
	}

	if !isValidFieldType(f.Type) {
		return ErrInvalidFieldType
	}

	if !isValidFieldAttributes(f.Type, f.Attributes) {
		return ErrInvalidFieldAttributes
	}

	if err := validate.ValidateStruct(f); err != nil {
		return err
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
