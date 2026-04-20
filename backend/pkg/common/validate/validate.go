package validate

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var v *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

func IsValidationErr(err error) bool {
	return errors.As(err, &validator.ValidationErrors{})
}

func ValidateStruct(s any) error {
	return v.Struct(s)
}
