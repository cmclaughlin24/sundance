package strategies

import (
	"errors"
	"fmt"

	"sundance/backend/services/forms/internal/core/domain"
)

var (
	ErrElementRequired  = errors.New("element required")
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
