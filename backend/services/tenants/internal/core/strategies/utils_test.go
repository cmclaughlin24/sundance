package strategies

import (
	"reflect"
	"testing"
	"time"
)

func TestToLowerCamelCase(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"should yield an empty string", "", ""},
		{"should lowercase the first letter", "IntervalHours", "intervalHours"},
		{"should lowercase a single letter", "A", "a"},
		{"should lowercase an acronym", "URL", "url"},
		{"should lowercase a leading acronym", "URLPath", "urlPath"},
		{"should leave a lower camel case string unchanged", "valueField", "valueField"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act.
			got := toLowerCamelCase(tt.in)

			// Assert.
			if got != tt.want {
				t.Errorf("expected %q but got %q", tt.want, got)
			}
		})
	}
}

func TestGoTypeToSchemaType(t *testing.T) {
	tests := []struct {
		name string
		in   reflect.Type
		want string
	}{
		{"should yield string for a string", reflect.TypeFor[string](), "string"},
		{"should yield integer for an int", reflect.TypeFor[int](), "integer"},
		{"should yield integer for an int64", reflect.TypeFor[int64](), "integer"},
		{"should yield number for a float64", reflect.TypeFor[float64](), "number"},
		{"should yield boolean for a bool", reflect.TypeFor[bool](), "boolean"},
		{"should yield array for a slice", reflect.TypeFor[[]string](), "array"},
		{"should yield array for an array", reflect.TypeFor[[2]int](), "array"},
		{"should yield object for a map", reflect.TypeFor[map[string]string](), "object"},
		{"should yield object for a struct", reflect.TypeFor[struct{ Name string }](), "object"},
		{"should yield string for a time", reflect.TypeFor[time.Time](), "string"},
		{"should yield string for a time pointer", reflect.TypeFor[*time.Time](), "string"},
		{"should yield the element type for a pointer", reflect.TypeFor[*int](), "integer"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act.
			got := goTypeToSchemaType(tt.in)

			// Assert.
			if got != tt.want {
				t.Errorf("expected %q but got %q", tt.want, got)
			}
		})
	}
}
