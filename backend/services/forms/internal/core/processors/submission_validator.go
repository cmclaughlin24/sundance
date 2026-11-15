package processors

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
	"sundance/backend/services/forms/internal/core/strategies"

	"golang.org/x/sync/errgroup"
)

var (
	ErrSubmissionFormMismatch = errors.New("submission form does not match form version")
	ErrUnknownSection         = errors.New("submission references a section not in the form version")
	ErrUnknownElement         = errors.New("submission references an element not in the section")
)

type submissionValidator struct {
	elementDefinitionStrategies ports.ElementDefinitionRegistry
}

func newSubmissionValidator(strats *ports.Strategies) *submissionValidator {
	return &submissionValidator{
		elementDefinitionStrategies: strats.ElementDefinition,
	}
}

func (v *submissionValidator) validateStructure(ctx context.Context, logger *slog.Logger, s *domain.Submission, fv *domain.FormVersion) error {
	if s.FormID != fv.FormID {
		logger.WarnContext(ctx, "structure validation failed; form mismatch", "version_form_id", fv.FormID)
		return fmt.Errorf("%w; submission_form=%s version_form=%s", ErrSubmissionFormMismatch, s.FormID, fv.FormID)
	}

	sections := make(map[domain.SectionID]*domain.Section)
	for _, page := range fv.GetPages() {
		for _, section := range page.GetSections() {
			sections[section.ID] = section
		}
	}

	for _, resp := range s.Sections {
		section, ok := sections[resp.SectionID]
		if !ok {
			logger.WarnContext(ctx, "structure validation failed; unknown section", "section_id", resp.SectionID)
			return fmt.Errorf("%w; section=%s", ErrUnknownSection, resp.SectionID)
		}

		for _, instance := range resp.Instances() {
			for _, val := range instance {
				if section.GetElement(val.ElementID) == nil {
					logger.WarnContext(ctx, "structure validation failed; unknown element", "section_id", section.ID, "element_id", val.ElementID)
					return fmt.Errorf("%w; section=%s element=%s", ErrUnknownElement, section.Key, val.ElementID)
				}
			}
		}
	}

	return nil
}

func (v *submissionValidator) validateResolved(ctx context.Context, logger *slog.Logger, resolved []resolveElement) error {
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
