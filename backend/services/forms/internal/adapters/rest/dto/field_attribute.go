package dto

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/strategy"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
)

var (
	ErrFieldAttrParse = errors.New("failed to deserialize field attributes")
)

type attributeParser func([]byte) (domain.FieldAttributes, error)

var attributeParsers = strategy.NewStrategies[domain.FieldType, attributeParser]().
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
