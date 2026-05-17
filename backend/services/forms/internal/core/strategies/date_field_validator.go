package strategies

import (
	"context"
	"log/slog"

	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/ports"
)

type DateFieldValidatorStrategy struct {
	logger *slog.Logger
}

func NewDateFieldValidatorStrategy(logger *slog.Logger) ports.FieldValidatorStrategy {
	return &DateFieldValidatorStrategy{
		logger: logger,
	}
}

func (s *DateFieldValidatorStrategy) Validate(ctx context.Context, field *domain.Field) error {
	_, err := domain.GetFieldAttributes[domain.DateFieldAttributes](field.Attributes)
	if err != nil {
		s.logger.ErrorContext(ctx, "strategy attributes mismatch", "field_id", field.ID, "field_type", field.Type, "error", err)
		return err
	}

	// TODO: Implement validation.

	return nil
}
