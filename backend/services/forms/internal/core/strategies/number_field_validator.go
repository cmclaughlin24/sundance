package strategies

import (
	"context"
	"fmt"
	"log/slog"

	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

type NumberFieldValidatorStrategy struct {
	logger *slog.Logger
}

func NewNumberFieldValidatorStrategy(logger *slog.Logger) ports.FieldValidatorStrategy {
	return &NumberFieldValidatorStrategy{
		logger: logger,
	}
}

func (s *NumberFieldValidatorStrategy) Validate(ctx context.Context, field domain.Field, fv domain.SubmissionFieldValue) error {
	attr, err := domain.GetFieldAttributes[domain.NumberFieldAttributes](field.Attributes)
	if err != nil {
		s.logger.ErrorContext(ctx, "strategy attributes mismatch", "field_id", field.ID, "field_type", field.Type, "error", err)
		return err
	}

	val, err := checkValueRequired[float64](attr, fv.Value)
	if err != nil {
		return newValidationErr(field.Key, err)
	}

	if fv.Value == nil {
		return nil
	}

	if attr.Min != nil && val < *attr.Min {
		return newValidationErr(field.Key, fmt.Errorf("min value"))
	}

	if attr.Max != nil && val > *attr.Max {
		return newValidationErr(field.Key, fmt.Errorf("max value"))
	}

	return nil
}
