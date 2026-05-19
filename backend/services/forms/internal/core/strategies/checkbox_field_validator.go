package strategies

import (
	"context"
	"log/slog"

	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

type CheckboxFieldValidatorStrategy struct {
	logger *slog.Logger
}

func NewCheckboxFieldValidatorStrategy(logger *slog.Logger) ports.FieldValidatorStrategy {
	return &CheckboxFieldValidatorStrategy{
		logger: logger,
	}
}

func (s *CheckboxFieldValidatorStrategy) Validate(ctx context.Context, field domain.Field, fv domain.SubmissionFieldValue) error {
	_, err := domain.GetFieldAttributes[domain.CheckboxFieldAttributes](field.Attributes)
	if err != nil {
		s.logger.ErrorContext(ctx, "strategy attributes mismatch", "field_id", field.ID, "field_type", field.Type, "error", err)
		return err
	}

	// TODO: Implement validation.

	return nil
}
