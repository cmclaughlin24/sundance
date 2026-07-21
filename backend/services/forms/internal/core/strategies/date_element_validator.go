package strategies

import (
	"context"
	"log/slog"

	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

type DateElementValidatorStrategy struct {
	logger *slog.Logger
}

func NewDateElementValidatorStrategy(logger *slog.Logger) ports.ElementValidatorStrategy {
	return &DateElementValidatorStrategy{
		logger: logger,
	}
}

func (s *DateElementValidatorStrategy) Validate(ctx context.Context, element domain.Element, sv domain.SubmissionValue) error {
	attr, err := domain.GetElementAttributes[*domain.DateElementAttributes](element.Attributes)
	if err != nil {
		s.logger.ErrorContext(ctx, "strategy attributes mismatch", "element_id", element.ID, "element_type", element.Type, "error", err)
		return err
	}

	_, err = checkValueRequired[string](attr, sv.Value)
	if err != nil {
		return newValidationErr(element.Key, err)
	}

	if sv.Value == nil {
		return nil
	}

	// TODO: Implement validation.

	return nil
}
