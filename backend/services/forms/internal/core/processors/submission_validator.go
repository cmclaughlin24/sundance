package processors

import (
	"context"
	"fmt"
	"log/slog"
	"sundance/backend/services/forms/internal/core/ports"
	"sundance/backend/services/forms/internal/core/strategies"
)

type submissionValidator struct {
	elementValidatorStrategies ports.ElementValidatorRegistry
}

func newSubmissionValidator(strats *ports.Strategies) *submissionValidator {
	return &submissionValidator{
		elementValidatorStrategies: strats.ElementValidator,
	}
}

func (v *submissionValidator) validate(ctx context.Context, logger *slog.Logger, resolved []resolveElement) error {
	for _, re := range resolved {
		element := re.element
		value := re.value

		elementValidator, err := v.elementValidatorStrategies.Get(element.Type)
		if err != nil {
			logger.ErrorContext(ctx, "failed to validate elements; missing element validation strategy", "element_id", element.ID, "element_type", element.Type)
			return err
		}

		if value == nil {
			if re.required {
				logger.WarnContext(ctx, "element validation failed; required element missing", "element_id", element.ID, "element_key", element.Key)
				return fmt.Errorf("%w; id=%s key=%s", strategies.ErrElementRequired, element.ID, element.Key)
			}

			continue
		}

		if err = elementValidator.Validate(ctx, *element, *value); err != nil {
			return err
		}
	}

	return nil
}
