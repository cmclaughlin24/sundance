package strategies

import (
	"context"
	"log/slog"

	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

type DateFieldValidatorStrategy struct {
	logger *slog.Logger
}

func NewDateFieldValidatorStrategy(logger *slog.Logger) ports.FieldValidatorStrategy {
	return &DateFieldValidatorStrategy{
		logger: logger,
	}
}

func (s *DateFieldValidatorStrategy) Validate(ctx context.Context, field domain.Field, fv domain.SubmissionFieldValue) error {
	attr, err := domain.GetFieldAttributes[domain.DateFieldAttributes](field.Attributes)
	if err != nil {
		s.logger.ErrorContext(ctx, "strategy attributes mismatch", "field_id", field.ID, "field_type", field.Type, "error", err)
		return err
	}

	_, err = checkValueRequired[string](attr, fv.Value)
	if err != nil {
		return err
	}

	if fv.Value == nil {
		return nil
	}

	// TODO: Implement validation.

	return nil
}
