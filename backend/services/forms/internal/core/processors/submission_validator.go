package processors

import (
	"context"
	"fmt"
	"log/slog"
	"sundance/backend/services/forms/internal/core/ports"
	"sundance/backend/services/forms/internal/core/strategies"
)

type submissionValidator struct {
	fieldValidatorStrategies ports.FieldValidatorRegistry
}

func newSubmissionValidator(strats *ports.Strategies) *submissionValidator {
	return &submissionValidator{
		fieldValidatorStrategies: strats.FieldValidator,
	}
}

func (v *submissionValidator) validate(ctx context.Context, logger *slog.Logger, resolved []resolveField) error {
	for _, rf := range resolved {
		field := rf.field
		value := rf.value

		fieldValidator, err := v.fieldValidatorStrategies.Get(field.Type)
		if err != nil {
			logger.ErrorContext(ctx, "failed to validate fields; missing field validation strategy", "field_id", field.ID, "field_type", field.Type)
			return err
		}

		if value == nil {
			if rf.required {
				logger.WarnContext(ctx, "field validation failed; required field missing", "field_id", field.ID, "field_key", field.Key)
				return fmt.Errorf("%w; id=%s key=%s", strategies.ErrFieldRequired, field.ID, field.Key)
			}

			continue
		}

		if err = fieldValidator.Validate(ctx, *field, *value); err != nil {
			return err
		}
	}

	return nil
}
