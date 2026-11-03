package strategies

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

var (
	ErrElementRequired   = errors.New("element required")
	ErrElementValidation = errors.New("element validation")
	ErrElementTypeValue  = errors.New("value does not match expected for element type")
)

func checkValueRequired[T comparable](attr domain.ElementAttributes, value any) (T, error) {
	var zero T

	switch val := value.(type) {
	case nil:
		if attr.GetIsRequired() {
			return zero, ErrElementRequired
		}

		return zero, nil
	case T:
		return val, nil
	default:
		return zero, ErrElementTypeValue
	}
}

func newValidationErr(key string, err error) error {
	return fmt.Errorf("%w: '%s' failed on %w", ErrElementValidation, key, err)
}

func reflectAttributes[T domain.ElementAttributes]() []ports.AttributeDescriptor {
	t := reflect.TypeFor[T]()

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil
	}

	var descriptors []ports.AttributeDescriptor
	extractAttributes(t, &descriptors)
	return descriptors
}

func extractAttributes(t reflect.Type, out *[]ports.AttributeDescriptor) {
	for i := range t.NumField() {
		field := t.Field(i)

		if !field.IsExported() {
			continue
		}

		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			extractAttributes(field.Type, out)
			continue
		}

		*out = append(*out, ports.AttributeDescriptor{
			Type:        goTypeToSchemaType(field.Type),
			Name:        toLowerCamelCase(field.Name),
			Description: field.Tag.Get("desc"),
			Required:    field.Type.Kind() != reflect.Pointer,
		})
	}
}

func toLowerCamelCase(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func goTypeToSchemaType(t reflect.Type) string {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.Slice, reflect.Array:
		return "array"
	case reflect.Struct, reflect.Map:
		return "object"
	default:
		return "string"
	}
}
