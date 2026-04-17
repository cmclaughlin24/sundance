package validate

import (
	"github.com/go-playground/validator/v10"
)

var v *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

func IsValidationErr(err error) bool {
	_, ok := err.(validator.ValidationErrors)

	if !ok {
		return false
	}

	return  true
}

func ValidateStruct(s any) error {
	if err := v.Struct(s); err != nil {
		return err
	}

	return nil
}
