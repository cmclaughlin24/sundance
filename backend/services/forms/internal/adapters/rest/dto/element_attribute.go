package dto

import (
	"encoding/json"
	"errors"
	"fmt"

	"sundance/backend/pkg/common/stratreg"
	"sundance/backend/services/forms/internal/core/domain"
)

var (
	ErrElementAttrParse = errors.New("failed to deserialize element attributes")
)

type attributeParser func([]byte) (domain.ElementAttributes, error)

var attributeParsers = stratreg.New[domain.ElementType, attributeParser]().
	Set(domain.ElementTypeText, func(data []byte) (domain.ElementAttributes, error) {
		return parseAttributes[*domain.TextElementAttributes](data)
	}).
	Set(domain.ElementTypeNumber, func(data []byte) (domain.ElementAttributes, error) {
		return parseAttributes[*domain.NumberElementAttributes](data)
	}).
	Set(domain.ElementTypeCheckbox, func(data []byte) (domain.ElementAttributes, error) {
		return parseAttributes[*domain.CheckboxElementAttributes](data)
	}).
	Set(domain.ElementTypeSelect, func(data []byte) (domain.ElementAttributes, error) {
		return parseAttributes[*domain.SelectElementAttributes](data)
	}).
	Set(domain.ElementTypeDate, func(data []byte) (domain.ElementAttributes, error) {
		return parseAttributes[*domain.DateElementAttributes](data)
	})

func attributesFromRequest(elementType domain.ElementType, raw any) (domain.ElementAttributes, error) {
	if elementType == "" {
		return nil, errors.New("element type is required")
	}

	attrBytes, err := json.Marshal(raw)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrElementAttrParse, err)
	}

	strategy, err := attributeParsers.Get(elementType)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrElementAttrParse, err)
	}

	return strategy(attrBytes)
}

func parseAttributes[T domain.ElementAttributes](data []byte) (domain.ElementAttributes, error) {
	var attributes T

	if err := json.Unmarshal(data, &attributes); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrElementAttrParse, err)
	}

	return attributes, nil
}

type baseElementAttributeResponse struct {
	IsRequired bool `json:"isRequired"`
	IsReadOnly bool `json:"isReadOnly"`
}

type textElementAttributesResponse struct {
	baseElementAttributeResponse
	MinLength   *int   `json:"minLength"`
	MaxLength   *int   `json:"maxLength"`
	Pattern     string `json:"pattern"`
	Placeholder string `json:"placeholder"`
}

type numberElementAttributesResponse struct {
	baseElementAttributeResponse
	Min  *float64 `json:"min"`
	Max  *float64 `json:"max"`
	Step *float64 `json:"step"`
}

type selectElementAttributesResponse struct {
	baseElementAttributeResponse
	Data          []any                 `json:"data"`
	DataSourceRef *domain.DataSourceRef `json:"dataSourceRef"`
	Multiple      bool                  `json:"multiple"`
	MinSelected   *int                  `json:"minSelected"`
	MaxSelected   *int                  `json:"maxSelected"`
}

type checkboxElementAttributesResponse struct {
	baseElementAttributeResponse
	IsCheckedByDefault bool                  `json:"isCheckedByDefault"`
	Data               []any                 `json:"data"`
	DataSourceRef      *domain.DataSourceRef `json:"dataSourceRef"`
}

type dateElementAttributesResponse struct {
	baseElementAttributeResponse
	MinDate *string `json:"minDate"`
	MaxDate *string `json:"maxDate"`
}

func elementAttributesToResponse(attr domain.ElementAttributes) any {
	base := baseElementAttributeResponse{
		IsRequired: attr.GetIsRequired(),
		IsReadOnly: attr.GetIsReadOnly(),
	}

	switch t := attr.(type) {
	case *domain.TextElementAttributes:
		return textElementAttributesResponse{
			baseElementAttributeResponse: base,
			MinLength:                    t.MinLength,
			MaxLength:                    t.MaxLength,
			Pattern:                      t.Pattern,
			Placeholder:                  t.Placeholder,
		}
	case *domain.NumberElementAttributes:
		return numberElementAttributesResponse{
			baseElementAttributeResponse: base,
			Min:                          t.Min,
			Max:                          t.Max,
			Step:                         t.Step,
		}
	case *domain.SelectElementAttributes:
		return selectElementAttributesResponse{
			baseElementAttributeResponse: base,
			Data:                         t.Data,
			DataSourceRef:                t.DataSourceRef,
			Multiple:                     t.Multiple,
			MinSelected:                  t.MinSelected,
			MaxSelected:                  t.MaxSelected,
		}
	case *domain.CheckboxElementAttributes:
		return checkboxElementAttributesResponse{
			baseElementAttributeResponse: base,
			IsCheckedByDefault:           t.IsCheckedByDefault,
			Data:                         t.Data,
			DataSourceRef:                t.DataSourceRef,
		}
	case *domain.DateElementAttributes:
		return dateElementAttributesResponse{
			baseElementAttributeResponse: base,
			MinDate:                      t.MinDate,
			MaxDate:                      t.MaxDate,
		}
	default:
		return attr
	}
}
