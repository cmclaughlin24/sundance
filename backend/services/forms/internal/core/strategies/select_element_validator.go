package strategies

import (
	"context"
	"log/slog"

	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

type SelectElementValidatorStrategy struct {
	logger *slog.Logger
}

func NewSelectElementValidatorStrategy(logger *slog.Logger) ports.ElementValidatorStrategy {
	return &SelectElementValidatorStrategy{
		logger: logger,
	}
}

func (s *SelectElementValidatorStrategy) Validate(ctx context.Context, element domain.Element, sv domain.SubmissionValue) error {
	_, err := domain.GetElementAttributes[*domain.SelectElementAttributes](element.Attributes)
	if err != nil {
		s.logger.ErrorContext(ctx, "strategy attributes mismatch", "element_id", element.ID, "element_type", element.Type, "error", err)
		return err
	}

	// TODO: Implement validation.

	return nil
}
