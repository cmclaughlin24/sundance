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

func NewTypeValidator[T comparable](values []T) func(T) bool {
	return func(t T) bool {
		return slices.Contains(values, t)
	}
}
