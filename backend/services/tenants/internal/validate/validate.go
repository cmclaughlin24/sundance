package validate

import (
	"github.com/go-playground/validator/v10"
)

var v *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

func ValidateStruct(s any) error {
	if err := v.Struct(s); err != nil {
		return err
	}

	return nil
}
