package processors

import (
	"context"
	"log/slog"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

type ruleByTypeGetter interface {
	GetRuleByType(domain.RuleType) *domain.Rule
}

type resolveElement struct {
	element  *domain.Element
	value    *domain.SubmissionValue
	required bool
}

type submissionResolver struct {
	logger    *slog.Logger
	evaluator ports.RuleEvaluator
}

func newSubmissionResolver(logger *slog.Logger, evaluator ports.RuleEvaluator) *submissionResolver {
	return &submissionResolver{
		logger:    logger,
		evaluator: evaluator,
	}
}

func (r *submissionResolver) resolve(ctx context.Context, s *domain.Submission, fv *domain.FormVersion) ([]resolveElement, error) {
	if fv.Status == domain.FormVersionStatusDraft {
		r.logger.WarnContext(ctx, "skipping submission job; invalid status", "submission_id", s.ID, "form_id", s.FormID, "version_id", fv.ID, "version_status", fv.Status)
		return nil, domain.ErrInvalidVersionStatus
	}

	evalContext := s.GetEvalContext()
	fieldCtx := make(map[string]any, len(s.Values))

	for _, element := range fv.FlatElements() {
		// FIXME: Can result in a bug where only the first found element is checked in rule expressions.
		if val, ok := s.GetValue(element.ID); ok {
			fieldCtx[element.Key] = val.Value
		}
	}
	evalContext[domain.RuleExprSourceTypeField] = fieldCtx

	resolved := make([]resolveElement, 0)

pageLoop:
	for _, page := range fv.GetPages() {
		visible, err := r.shouldValidate(ctx, page, evalContext)

		if err != nil {
			return nil, err
		}

		if !visible {
			continue pageLoop
		}

	sectionLoop:
		for _, section := range page.GetSections() {
			visible, err := r.shouldValidate(ctx, section, evalContext)

			if err != nil {
				return nil, err
			}

			if !visible {
				continue sectionLoop
			}

		elementLoop:
			for _, element := range section.GetElements() {
				visible, err := r.shouldValidate(ctx, element, evalContext)

				if err != nil {
					return nil, err
				}

				if !visible {
					continue elementLoop
				}

				required, err := r.isRequired(ctx, element, evalContext)
				if err != nil {
					return nil, err
				}

				if required == nil {
					req := element.Attributes.GetIsRequired()
					required = &req
				}

				values := s.GetValues(element.ID)
				if len(values) == 0 {
					resolved = append(resolved, resolveElement{
						element:  element,
						value:    nil,
						required: *required,
					})
					continue elementLoop
				}

				for _, v := range values {
					resolved = append(resolved, resolveElement{
						element:  element,
						value:    v,
						required: *required,
					})
				}
			}
		}
	}

	return resolved, nil
}

func (s *submissionResolver) isRequired(ctx context.Context, rg ruleByTypeGetter, evalCtx ports.RuleEvaluationContext) (*bool, error) {
	rule := rg.GetRuleByType(domain.RuleTypeRequired)
	if rule == nil {
		return nil, nil
	}

	result, err := s.evaluator.Evaluate(ctx, rule, evalCtx)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *submissionResolver) shouldValidate(ctx context.Context, rg ruleByTypeGetter, evalCtx ports.RuleEvaluationContext) (bool, error) {
	rule := rg.GetRuleByType(domain.RuleTypeVisible)

	if rule == nil {
		return true, nil
	}

	return s.evaluator.Evaluate(ctx, rule, evalCtx)
}
