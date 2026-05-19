package strategies

import (
	"fmt"

	"sundance/backend/services/forms/internal/core/domain"
)

func checkValueRequired[T comparable](attr domain.FieldAttributes, value any) (T, error) {
	var zero T

	switch val := value.(type) {
	case nil:
		if attr.GetIsRequired() {
			return zero, fmt.Errorf("")
		}

		return zero, nil
	case T:
		if attr.GetIsRequired() && val == zero {
			return zero, fmt.Errorf("")
		}

		return val, nil
	default:
		return zero, fmt.Errorf("")
	}
}
