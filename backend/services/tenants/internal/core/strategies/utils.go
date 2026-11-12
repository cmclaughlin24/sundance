package strategies

import (
	"reflect"
	"strings"
	"time"
	"unicode"

	"sundance/backend/services/tenants/internal/core/domain"
	"sundance/backend/services/tenants/internal/core/ports"
)

func reflectAttributes[T domain.DataSourceAttributes]() []ports.AttributeDescriptor {
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

		// NOTE: "-" denotes when a field should be ignored.
		desc := field.Tag.Get("desc")
		if desc == "-" {
			continue
		}

		*out = append(*out, ports.AttributeDescriptor{
			Type:        goTypeToSchemaType(field.Type),
			Name:        toLowerCamelCase(field.Name),
			Description: desc,
			Required:    field.Type.Kind() != reflect.Pointer,
		})
	}
}

func toLowerCamelCase(s string) string {
	// NOTE: Lowercase the leading run of capitals so that acronyms are handled
	// (e.g. "URL" -> "url", "URLPath" -> "urlPath", "IntervalHours" -> "intervalHours").
	n := 0
	for n < len(s) && unicode.IsUpper(rune(s[n])) {
		n++
	}

	if n > 1 && n < len(s) {
		n--
	}

	return strings.ToLower(s[:n]) + s[n:]
}

func goTypeToSchemaType(t reflect.Type) string {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	// NOTE: time.Time is a struct, but is serialized as an RFC 3339 string.
	if t == reflect.TypeFor[time.Time]() {
		return "string"
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
