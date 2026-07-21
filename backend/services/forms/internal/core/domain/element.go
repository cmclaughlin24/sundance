package domain

import (
	"errors"
	"fmt"
	"slices"

	"sundance/backend/pkg/common/validate"
)

type ElementID string

type ElementType string

const (
	ElementTypeText     ElementType = "text"
	ElementTypeNumber   ElementType = "number"
	ElementTypeSelect   ElementType = "select"
	ElementTypeCheckbox ElementType = "checkbox"
	ElementTypeDate     ElementType = "date"
)

var (
	ErrInvalidElement             = errors.New("invalid element")
	ErrInvalidElementType         = errors.New("invalid element type")
	ErrInvalidElementAttributes   = errors.New("invalid element attributes for type")
	ErrDuplicateElementTagMapping = errors.New("duplicate element tag mapping for same tag version")
)

type Element struct {
	ID         ElementID
	Key        string `validate:"required,nowhitespace"`
	Name       string `validate:"required"`
	Type       ElementType
	Attributes ElementAttributes
	tags       []*ElementTagMapping
	withPosition
	withRules
}

func NewElement(key, name string, elementType ElementType, attr ElementAttributes, position float32) (*Element, error) {
	e := &Element{
		ID:         ElementID(NewID()),
		Key:        key,
		Name:       name,
		Type:       elementType,
		Attributes: attr,
		withPosition: withPosition{
			position: position,
		},
	}

	if err := e.validate(); err != nil {
		return nil, err
	}

	return e, nil
}

func HydrateElement(
	id ElementID,
	key,
	name string,
	elementType ElementType,
	attr ElementAttributes,
	position float32,
	tags []*ElementTagMapping,
) *Element {
	return &Element{
		ID:         id,
		Key:        key,
		Name:       name,
		Type:       elementType,
		Attributes: attr,
		tags:       tags,
		withPosition: withPosition{
			position: position,
		},
	}
}

func (e *Element) Update(key, name string, elementType ElementType, attr ElementAttributes, position float32) error {
	if e == nil {
		return ErrInvalidElement
	}

	cpy := *e
	cpy.Key = key
	cpy.Name = name
	cpy.Type = elementType
	cpy.Attributes = attr
	cpy.position = position

	if err := cpy.validate(); err != nil {
		return err
	}

	*e = cpy

	return nil
}

func (e *Element) GetTags() []*ElementTagMapping {
	return e.tags
}

func (e *Element) AddTags(mappings ...ElementTagMappingConfig) error {
	for _, tag := range mappings {
		idx := slices.IndexFunc(e.tags, func(etm *ElementTagMapping) bool {
			return etm.TagVersionID == tag.TagVersionID
		})

		if idx != -1 {
			return fmt.Errorf("%w: tagVersion=%s", ErrDuplicateElementTagMapping, tag.TagVersionID)
		}

		etm, err := NewElementTagMapping(e.ID, tag.TagVersionID, tag.Priority)
		if err != nil {
			return err
		}

		e.tags = append(e.tags, etm)
	}
	return nil
}

func (e *Element) ReplaceTags(mappings ...ElementTagMappingConfig) error {
	tags := make([]*ElementTagMapping, 0, len(mappings))

	for _, tag := range mappings {
		idx := slices.IndexFunc(tags, func(etm *ElementTagMapping) bool {
			return etm.TagVersionID == tag.TagVersionID
		})

		if idx != -1 {
			return fmt.Errorf("%w: tagVersion=%s", ErrDuplicateElementTagMapping, tag.TagVersionID)
		}

		etm, err := NewElementTagMapping(e.ID, tag.TagVersionID, tag.Priority)
		if err != nil {
			return err
		}

		tags = append(tags, etm)
	}

	e.tags = tags

	return nil
}

func (e *Element) validate() error {
	if !isValidPosition(e.position) {
		return ErrInvalidPosition
	}

	if !isValidElementType(e.Type) {
		return ErrInvalidElementType
	}

	if !isValidElementAttributes(e.Type, e.Attributes) {
		return ErrInvalidElementAttributes
	}

	if err := validate.ValidateStruct(e); err != nil {
		return err
	}

	return nil
}

var isValidElementType = validate.NewTypeValidator([]ElementType{
	ElementTypeText,
	ElementTypeNumber,
	ElementTypeSelect,
	ElementTypeCheckbox,
	ElementTypeDate,
})

func isValidElementAttributes(elementType ElementType, attr ElementAttributes) bool {
	switch attr.(type) {
	case *TextElementAttributes:
		return elementType == ElementTypeText
	case *NumberElementAttributes:
		return elementType == ElementTypeNumber
	case *SelectElementAttributes:
		return elementType == ElementTypeSelect
	case *CheckboxElementAttributes:
		return elementType == ElementTypeCheckbox
	case *DateElementAttributes:
		return elementType == ElementTypeDate
	default:
		return false
	}
}
