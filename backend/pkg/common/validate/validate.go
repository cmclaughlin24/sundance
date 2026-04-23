package validate

import (
	"errors"
	"slices"

	"github.com/go-playground/validator/v10"
)

var v *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

func IsValidationErr(err error) bool {
	return errors.As(err, &validator.ValidationErrors{})
}

func ValidateStruct(s any) error {
	return v.Struct(s)
}

func NewTypeValidator[T comparable](types []T) func(T) bool {
	cpy := slices.Clone(types)

	return func(t T) bool {
		return slices.Contains(cpy, t)
	}
}
