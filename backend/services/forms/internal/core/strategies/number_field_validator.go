package strategies

import (
	"context"
	"log/slog"

	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/ports"
)

type NumberFieldValidatorStrategy struct {
	logger *slog.Logger
}

func NewNumberFieldValidatorStrategy(logger *slog.Logger) ports.FieldValidatorStrategy {
	return &NumberFieldValidatorStrategy{
		logger: logger,
	}
}

func (s *NumberFieldValidatorStrategy) Validate(ctx context.Context, field domain.Field, value domain.SubmissionFieldValue) error {
	_, err := domain.GetFieldAttributes[domain.NumberFieldAttributes](field.Attributes)
	if err != nil {
		s.logger.ErrorContext(ctx, "strategy attributes mismatch", "field_id", field.ID, "field_type", field.Type, "error", err)
		return err
	}

	// TODO: Implement validation.

	return nil
}
