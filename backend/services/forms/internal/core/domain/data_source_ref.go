package domain

import (
	"errors"
	"sundance/backend/pkg/common/validate"
)

type BindingSourceType string

const (
	BindingSourceTypeField BindingSourceType = "field"
	BindSourceTypeStatic   BindingSourceType = "static"
)

var (
	ErrInvalidBindingSourceType = errors.New("invalid binding source type")
)

type BindingSource struct {
	Type  BindingSourceType
	Key   string
	Value any
}

func NewBindingSourceType(sourceType BindingSourceType, key string, value any) (BindingSource, error) {
	if !isValidBindingSourceType(sourceType) {
		return BindingSource{}, ErrInvalidBindingSourceType
	}

	bs := BindingSource{
		Type:  sourceType,
		Key:   key,
		Value: value,
	}

	return bs, nil
}

var isValidBindingSourceType = validate.NewTypeValidator([]BindingSourceType{
	BindingSourceTypeField,
	BindSourceTypeStatic,
})

type DataSourceRef struct {
	DataSourceID string
	Bindings     map[string]BindingSource
}
