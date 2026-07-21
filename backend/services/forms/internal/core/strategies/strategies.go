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

	elementValidatorStrategies := stratreg.New[domain.ElementType, ports.ElementValidatorStrategy]().
		Set(domain.ElementTypeText, NewTextElementValidatorStrategy(so.logger)).
		Set(domain.ElementTypeNumber, NewNumberElementValidatorStrategy(so.logger)).
		Set(domain.ElementTypeSelect, NewSelectElementValidatorStrategy(so.logger)).
		Set(domain.ElementTypeCheckbox, NewCheckboxElementValidatorStrategy(so.logger)).
		Set(domain.ElementTypeDate, NewDateElementValidatorStrategy(so.logger))

	return &ports.Strategies{
		ElementValidator: elementValidatorStrategies,
	}
}

func WithLogger(logger *slog.Logger) func(*strategyOptions) {
	return func(so *strategyOptions) {
		so.logger = logger
	}
}
