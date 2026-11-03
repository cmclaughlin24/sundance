package strategies

import (
	"context"
	"log/slog"

	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

type CheckboxElementDefinitionStrategy struct {
	logger *slog.Logger
}

func NewCheckboxElementDefinitionStrategy(logger *slog.Logger) ports.ElementDefinitionStrategy {
	return &CheckboxElementDefinitionStrategy{
		logger: logger,
	}
}

func (s *CheckboxElementDefinitionStrategy) Description() string {
	return "Checkbox input supporting boolean toggle state or multi-option selection."
}

func (s *CheckboxElementDefinitionStrategy) Attributes() []ports.AttributeDescriptor {
	return reflectAttributes[*domain.CheckboxElementAttributes]()
}

func (s *CheckboxElementDefinitionStrategy) Validate(ctx context.Context, element domain.Element, sv domain.SubmissionValue) error {
	_, err := domain.GetElementAttributes[*domain.CheckboxElementAttributes](element.Attributes)
	if err != nil {
		s.logger.ErrorContext(ctx, "strategy attributes mismatch", "element_id", element.ID, "element_type", element.Type, "error", err)
		return err
	}

	// TODO: Implement validation.

	return nil
}
