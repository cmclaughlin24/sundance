package domain

import (
	"errors"
	"sundance/backend/pkg/common/validate"
)

type BindingSourceType string

const (
	BindingSourceTypeField  BindingSourceType = "field"
	BindingSourceTypeStatic BindingSourceType = "static"
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
	BindingSourceTypeStatic,
})

type DataSourceRef struct {
	DataSourceID string
	Bindings     map[string]BindingSource
}

func getReferencedKeys(ds *DataSourceRef) []string {
	keys := make([]string, 0)

	if ds == nil {
		return keys
	}

	for _, bs := range ds.Bindings {
		if bs.Type != BindingSourceTypeField {
			continue
		}

		keys = append(keys, bs.Key)
	}

	return keys
}
