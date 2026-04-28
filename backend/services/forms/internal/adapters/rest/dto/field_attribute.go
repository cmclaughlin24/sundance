package dto

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/stratreg"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
)

var (
	ErrFieldAttrParse = errors.New("failed to deserialize field attributes")
)

type attributeParser func([]byte) (domain.FieldAttributes, error)

var attributeParsers = stratreg.New[domain.FieldType, attributeParser]().
	Set(domain.FieldTypeText, func(data []byte) (domain.FieldAttributes, error) {
		return parseAttributes[domain.TextFieldAttributes](data)
	}).
	Set(domain.FieldTypeNumber, func(data []byte) (domain.FieldAttributes, error) {
		return parseAttributes[domain.NumberFieldAttributes](data)
	}).
	Set(domain.FieldTypeCheckbox, func(data []byte) (domain.FieldAttributes, error) {
		return parseAttributes[domain.CheckboxFieldAttributes](data)
	}).
	Set(domain.FieldTypeSelect, func(data []byte) (domain.FieldAttributes, error) {
		return parseAttributes[domain.SelectFieldAttributes](data)
	}).
	Set(domain.FieldTypeDate, func(data []byte) (domain.FieldAttributes, error) {
		return parseAttributes[domain.DateFieldAttributes](data)
	})

func attributesFromRequest(fieldType domain.FieldType, raw any) (domain.FieldAttributes, error) {
	if fieldType == "" {
		return nil, errors.New("field type is required")
	}

	attrBytes, err := json.Marshal(raw)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFieldAttrParse, err)
	}

	strategy, err := attributeParsers.Get(fieldType)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFieldAttrParse, err)
	}

	return strategy(attrBytes)
}

func parseAttributes[T domain.FieldAttributes](data []byte) (domain.FieldAttributes, error) {
	var attributes T

	if err := json.Unmarshal(data, &attributes); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFieldAttrParse, err)
	}

	return attributes, nil
}

type baseFieldAttributeResponse struct {
	IsRequired bool `json:"isRequired"`
	IsReadOnly bool `json:"isReadOnly"`
}

type textFieldAttributesResponse struct {
	baseFieldAttributeResponse
	MinLength   *int   `json:"minLength"`
	MaxLength   *int   `json:"maxLength"`
	Pattern     string `json:"pattern"`
	Placeholder string `json:"placeholder"`
}

type numberFieldAttributesResponse struct {
	baseFieldAttributeResponse
	Min  *float64 `json:"min"`
	Max  *float64 `json:"max"`
	Step *float64 `json:"step"`
}

type selectFieldAttributesResponse struct {
	baseFieldAttributeResponse
	Multiple    bool `json:"multiple"`
	MinSelected *int `json:"minSelected"`
	MaxSelected *int `json:"maxSelected"`
}

type checkboxFieldAttributesResponse struct {
	baseFieldAttributeResponse
	IsCheckedByDefault bool `json:"isCheckedByDefault"`
}

type dateFieldAttributesResponse struct {
	baseFieldAttributeResponse
	MinDate *string `json:"minDate"`
	MaxDate *string `json:"maxDate"`
}

func fieldAttributesToResponse(attr domain.FieldAttributes) any {
	base := baseFieldAttributeResponse{
		IsRequired: attr.GetIsRequired(),
		IsReadOnly: attr.GetIsReadOnly(),
	}

	switch t := attr.(type) {
	case domain.TextFieldAttributes:
		return textFieldAttributesResponse{
			baseFieldAttributeResponse: base,
			MinLength:                  t.MinLength,
			MaxLength:                  t.MaxLength,
			Pattern:                    t.Pattern,
			Placeholder:                t.Placeholder,
		}
	case domain.NumberFieldAttributes:
		return numberFieldAttributesResponse{
			baseFieldAttributeResponse: base,
			Min:                        t.Min,
			Max:                        t.Max,
			Step:                       t.Step,
		}
	case domain.SelectFieldAttributes:
		return selectFieldAttributesResponse{
			baseFieldAttributeResponse: base,
			Multiple:                   t.Multiple,
			MinSelected:                t.MinSelected,
			MaxSelected:                t.MaxSelected,
		}
	case domain.CheckboxFieldAttributes:
		return checkboxFieldAttributesResponse{
			baseFieldAttributeResponse: base,
			IsCheckedByDefault:         t.IsCheckedByDefault,
		}
	case domain.DateFieldAttributes:
		return dateFieldAttributesResponse{
			baseFieldAttributeResponse: base,
			MinDate:                    t.MinDate,
			MaxDate:                    t.MaxDate,
		}
	default:
		return attr
	}
}
