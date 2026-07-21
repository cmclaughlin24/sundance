package strategies

import (
	"context"
	"fmt"
	"log/slog"

	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

type NumberElementValidatorStrategy struct {
	logger *slog.Logger
}

func NewNumberElementValidatorStrategy(logger *slog.Logger) ports.ElementValidatorStrategy {
	return &NumberElementValidatorStrategy{
		logger: logger,
	}
}

func (s *NumberElementValidatorStrategy) Validate(ctx context.Context, element domain.Element, sv domain.SubmissionValue) error {
	attr, err := domain.GetElementAttributes[*domain.NumberElementAttributes](element.Attributes)
	if err != nil {
		s.logger.ErrorContext(ctx, "strategy attributes mismatch", "element_id", element.ID, "element_type", element.Type, "error", err)
		return err
	}

	val, err := checkValueRequired[float64](attr, sv.Value)
	if err != nil {
		return newValidationErr(element.Key, err)
	}

	if sv.Value == nil {
		return nil
	}

	if attr.Min != nil && val < *attr.Min {
		return newValidationErr(element.Key, fmt.Errorf("min value"))
	}

	if attr.Max != nil && val > *attr.Max {
		return newValidationErr(element.Key, fmt.Errorf("max value"))
	}

	return nil
}
