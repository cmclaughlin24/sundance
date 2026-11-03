package processors

import (
	"context"
	"fmt"
	"log/slog"
	"sundance/backend/services/forms/internal/core/ports"
	"sundance/backend/services/forms/internal/core/strategies"

	"golang.org/x/sync/errgroup"
)

type submissionValidator struct {
	elementDefinitionStrategies ports.ElementDefinitionRegistry
}

func newSubmissionValidator(strats *ports.Strategies) *submissionValidator {
	return &submissionValidator{
		elementDefinitionStrategies: strats.ElementDefinition,
	}
}

func (v *submissionValidator) validate(ctx context.Context, logger *slog.Logger, resolved []resolveElement) error {
	g, gCtx := errgroup.WithContext(ctx)

	for _, re := range resolved {
		g.Go(func() error {
			element := re.element
			value := re.value

			elementDefinition, err := v.elementDefinitionStrategies.Get(element.Type)
			if err != nil {
				logger.ErrorContext(ctx, "failed to validate elements; missing element definition strategy", "element_id", element.ID, "element_type", element.Type)
				return err
			}

			if value == nil {
				if re.required {
					logger.WarnContext(ctx, "element validation failed; required element missing", "element_id", element.ID, "element_key", element.Key)
					return fmt.Errorf("%w; id=%s key=%s", strategies.ErrElementRequired, element.ID, element.Key)
				}

				return nil
			}

			return elementDefinition.Validate(gCtx, *element, *value)
		})
	}

	return g.Wait()
}
