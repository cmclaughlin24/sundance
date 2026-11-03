package strategies

import (
	"log/slog"

	"sundance/backend/pkg/common/stratreg"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

type strategyOptions struct {
	logger *slog.Logger
}

func Bootstrap(opts ...func(*strategyOptions)) *ports.Strategies {
	var so strategyOptions
	for _, opt := range opts {
		opt(&so)
	}

	elementDefinitionStrategies := stratreg.New[domain.ElementType, ports.ElementDefinitionStrategy]().
		Set(domain.ElementTypeText, NewTextElementDefinitionStrategy(so.logger)).
		Set(domain.ElementTypeNumber, NewNumberElementDefinitionStrategy(so.logger)).
		Set(domain.ElementTypeSelect, NewSelectElementDefinitionStrategy(so.logger)).
		Set(domain.ElementTypeCheckbox, NewCheckboxElementDefinitionStrategy(so.logger)).
		Set(domain.ElementTypeDate, NewDateElementDefinitionStrategy(so.logger))

	return &ports.Strategies{
		ElementDefinition: elementDefinitionStrategies,
	}
}

func WithLogger(logger *slog.Logger) func(*strategyOptions) {
	return func(so *strategyOptions) {
		so.logger = logger
	}
}
