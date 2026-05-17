package strategies

import (
	"log/slog"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/stratreg"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/ports"
)

type strategyOptions struct {
	logger *slog.Logger
}

func Bootstrap(opts ...func(*strategyOptions)) *ports.Strategies {
	var so strategyOptions
	for _, opt := range opts {
		opt(&so)
	}

	fieldValidatorStrategies := stratreg.New[domain.FieldType, ports.FieldValidatorStrategy]().
		Set(domain.FieldTypeText, NewTextFieldValidatorStrategy(so.logger)).
		Set(domain.FieldTypeNumber, NewNumberFieldValidatorStrategy(so.logger)).
		Set(domain.FieldTypeSelect, NewSelectFieldValidatorStrategy(so.logger)).
		Set(domain.FieldTypeCheckbox, NewCheckboxFieldValidatorStrategy(so.logger)).
		Set(domain.FieldTypeDate, NewDateFieldValidatorStrategy(so.logger))

	return &ports.Strategies{
		FieldValidator: fieldValidatorStrategies,
	}
}

func WithLogger(logger *slog.Logger) func(*strategyOptions) {
	return func(so *strategyOptions) {
		so.logger = logger
	}
}
