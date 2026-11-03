package strategies

import (
	"context"
	"log/slog"

	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

type SelectElementDefinitionStrategy struct {
	logger *slog.Logger
}

func NewSelectElementDefinitionStrategy(logger *slog.Logger) ports.ElementDefinitionStrategy {
	return &SelectElementDefinitionStrategy{
		logger: logger,
	}
}

func (s *SelectElementDefinitionStrategy) Description() string {
	return "Dropdown selection input supporting static options or external data source lookups."
}

func (s *SelectElementDefinitionStrategy) Attributes() []ports.AttributeDescriptor {
	return reflectAttributes[*domain.SelectElementAttributes]()
}

func (s *SelectElementDefinitionStrategy) Validate(ctx context.Context, element domain.Element, sv domain.SubmissionValue) error {
	_, err := domain.GetElementAttributes[*domain.SelectElementAttributes](element.Attributes)
	if err != nil {
		s.logger.ErrorContext(ctx, "strategy attributes mismatch", "element_id", element.ID, "element_type", element.Type, "error", err)
		return err
	}

	// TODO: Implement validation.

	return nil
}
