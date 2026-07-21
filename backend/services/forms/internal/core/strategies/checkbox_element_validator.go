package strategies

import (
	"context"
	"log/slog"

	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

type CheckboxElementValidatorStrategy struct {
	logger *slog.Logger
}

func NewCheckboxElementValidatorStrategy(logger *slog.Logger) ports.ElementValidatorStrategy {
	return &CheckboxElementValidatorStrategy{
		logger: logger,
	}
}

func (s *CheckboxElementValidatorStrategy) Validate(ctx context.Context, element domain.Element, sv domain.SubmissionValue) error {
	_, err := domain.GetElementAttributes[*domain.CheckboxElementAttributes](element.Attributes)
	if err != nil {
		s.logger.ErrorContext(ctx, "strategy attributes mismatch", "element_id", element.ID, "element_type", element.Type, "error", err)
		return err
	}

	// TODO: Implement validation.

	return nil
}
