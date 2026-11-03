package strategies

import (
	"context"
	"log/slog"

	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

type DateElementDefinitionStrategy struct {
	logger *slog.Logger
}

func NewDateElementDefinitionStrategy(logger *slog.Logger) ports.ElementDefinitionStrategy {
	return &DateElementDefinitionStrategy{
		logger: logger,
	}
}

func (s *DateElementDefinitionStrategy) Description() string {
	return "Date selection input supporting minimum and maximum date boundaries."
}

func (s *DateElementDefinitionStrategy) Attributes() []ports.AttributeDescriptor {
	return reflectAttributes[*domain.DateElementAttributes]()
}

func (s *DateElementDefinitionStrategy) Validate(ctx context.Context, element domain.Element, sv domain.SubmissionValue) error {
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
