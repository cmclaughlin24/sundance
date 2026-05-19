package strategies

import (
	"errors"
	"fmt"

	"sundance/backend/services/forms/internal/core/domain"
)

var (
	ErrFieldRequired   = errors.New("field required")
	ErrFieldValidation = errors.New("field validation")
	ErrFieldTypeValue  = errors.New("value does not match expected for field type")
)

func checkValueRequired[T comparable](attr domain.FieldAttributes, value any) (T, error) {
	var zero T

	switch val := value.(type) {
	case nil:
		if attr.GetIsRequired() {
			return zero, ErrFieldRequired
		}

		return zero, nil
	case T:
		return val, nil
	default:
		return zero, ErrFieldTypeValue
	}
}

func newValidationErr(field string, err error) error {
	return fmt.Errorf("%w: '%s' failed on %w", ErrFieldValidation, field, err)
}
