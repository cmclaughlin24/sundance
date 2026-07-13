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

type resolveField struct {
	field    *domain.Field
	value    *domain.SubmissionFieldValue
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

func (r *submissionResolver) resolve(ctx context.Context, s *domain.Submission, fv *domain.FormVersion) ([]resolveField, error) {
	if fv.Status == domain.FormVersionStatusDraft {
		r.logger.WarnContext(ctx, "skipping submission job; invalid status", "submission_id", s.ID, "form_id", s.FormID, "version_id", fv.ID, "version_status", fv.Status)
		return nil, domain.ErrInvalidVersionStatus
	}

	evalCtx := make(ports.RuleEvaluationContext, len(s.Values))
	for _, field := range fv.FlatFields() {
		if val, ok := s.GetFieldValue(field.ID); ok {
			evalCtx[field.Key] = val.Value
		}
	}

	resolved := make([]resolveField, 0)

pageLoop:
	for _, page := range fv.GetPages() {
		visible, err := r.shouldValidate(ctx, page, evalCtx)

		if err != nil {
			return nil, err
		}

		if !visible {
			continue pageLoop
		}

	sectionLoop:
		for _, section := range page.GetSections() {
			visible, err := r.shouldValidate(ctx, section, evalCtx)

			if err != nil {
				return nil, err
			}

			if !visible {
				continue sectionLoop
			}

		fieldLoop:
			for _, field := range section.GetFields() {
				visible, err := r.shouldValidate(ctx, field, evalCtx)

				if err != nil {
					return nil, err
				}

				if !visible {
					continue fieldLoop
				}

				required, err := r.isRequired(ctx, field, evalCtx)
				if err != nil {
					return nil, err
				}

				if required == nil {
					req := field.Attributes.GetIsRequired()
					required = &req
				}

				val, _ := s.GetFieldValue(field.ID)
				resolved = append(resolved, resolveField{
					field:    field,
					value:    val,
					required: *required,
				})
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
