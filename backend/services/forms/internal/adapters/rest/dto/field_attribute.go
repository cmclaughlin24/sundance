package dto

import (
	"encoding/json"
	"errors"

	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
)

type attributeParser func([]byte) (domain.FieldAttributes, error)

var attributeParsers = map[domain.FieldType]attributeParser{
	domain.FieldTypeText: func(data []byte) (domain.FieldAttributes, error) {
		return parseAttributes[domain.TextFieldAttributes](data)
	},
	domain.FieldTypeNumber: func(data []byte) (domain.FieldAttributes, error) {
		return parseAttributes[domain.NumberFieldAttributes](data)
	},
	domain.FieldTypeCheckbox: func(data []byte) (domain.FieldAttributes, error) {
		return parseAttributes[domain.CheckboxFieldAttributes](data)
	},
	domain.FieldTypeSelect: func(data []byte) (domain.FieldAttributes, error) {
		return parseAttributes[domain.SelectFieldAttributes](data)
	},
	domain.FieldTypeDate: func(data []byte) (domain.FieldAttributes, error) {
		return parseAttributes[domain.DateFieldAttributes](data)
	},
}

func attributesFromRequest(fieldType domain.FieldType, raw any) (domain.FieldAttributes, error) {
	if fieldType == "" {
		return nil, errors.New("field type is required")
	}

	attrBytes, err := json.Marshal(raw)

	if err != nil {
		return nil, err
	}

	parser, ok := attributeParsers[fieldType]

	if !ok {
		return nil, errors.New("unsupported field type")
	}

	return parser(attrBytes)
}

func parseAttributes[T any](data []byte) (domain.FieldAttributes, error) {
	var attributes T

	if err := json.Unmarshal(data, &attributes); err != nil {
		return nil, err
	}

	return attributes, nil
}
