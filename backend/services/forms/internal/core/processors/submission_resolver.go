package processors

import (
	"context"
	"log/slog"
	"slices"

	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

type ruleByTypeGetter interface {
	GetRuleByType(domain.RuleType) *domain.Rule
}

type resolveElement struct {
	element  *domain.Element
	value    *domain.SubmissionValue
	path     []int // [instanceIndex?, itemIndex?]; nil for scalar values in non-repeatable sections
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

	base := domain.BaseEvalContext(s, fv)
	resolved := make([]resolveElement, 0)

pageLoop:
	for _, page := range fv.GetPages() {
		visible, err := r.shouldValidate(ctx, page, base)
		if err != nil {
			return nil, err
		}

		if !visible {
			continue pageLoop
		}

	sectionLoop:
		for _, section := range page.GetSections() {
			visible, err := r.shouldValidate(ctx, section, base)
			if err != nil {
				return nil, err
			}

			if !visible {
				continue sectionLoop
			}

			instances := s.GetSection(section.ID).Instances()
			if err := section.ValidateInstanceCount(len(instances)); err != nil {
				r.logger.WarnContext(ctx, "section instance count invalid", "submission_id", s.ID, "section_id", section.ID, "section_key", section.Key, "error", err)
				return nil, err
			}

			if !section.IsRepeatable() {
				var instance domain.SectionInstance
				if len(instances) == 1 {
					instance = instances[0]
				}

				res, err := r.resolveSection(ctx, section, instance, nil, base)
				if err != nil {
					return nil, err
				}

				resolved = append(resolved, res...)
				continue sectionLoop
			}

			for idx, instance := range instances {
				res, err := r.resolveSection(ctx, section, instance, []int{idx}, base.WithFields(section.FieldValues(instance)))
				if err != nil {
					return nil, err
				}

				resolved = append(resolved, res...)
			}
		}
	}

	return resolved, nil
}

func (r *submissionResolver) resolveSection(
	ctx context.Context,
	section *domain.Section,
	instance domain.SectionInstance,
	prefix []int,
	evalCtx ports.RuleEvaluationContext,
) ([]resolveElement, error) {
	var resolved []resolveElement

	for _, element := range section.GetElements() {
		res, err := r.resolveElement(ctx, element, instance.Values(element.ID), prefix, evalCtx)
		if err != nil {
			return nil, err
		}

		resolved = append(resolved, res...)
	}

	return resolved, nil
}

func (r *submissionResolver) resolveElement(
	ctx context.Context,
	element *domain.Element,
	values []*domain.SubmissionValue,
	prefix []int,
	evalCtx ports.RuleEvaluationContext,
) ([]resolveElement, error) {
	visible, err := r.shouldValidate(ctx, element, evalCtx)
	if err != nil {
		return nil, err
	}

	if !visible {
		return nil, nil
	}

	required, err := r.isRequired(ctx, element, evalCtx)
	if err != nil {
		return nil, err
	}

	if required == nil {
		req := element.Attributes.GetIsRequired()
		required = &req
	}

	if len(values) == 0 {
		return []resolveElement{{
			element:  element,
			path:     slices.Clone(prefix),
			required: *required,
		}}, nil
	}

	resolved := make([]resolveElement, 0, len(values))
	for _, v := range values {
		resolved = append(resolved, resolveElement{
			element:  element,
			value:    v,
			path:     valuePath(prefix, v),
			required: *required,
		})
	}

	return resolved, nil
}

func (r *submissionResolver) isRequired(ctx context.Context, rg ruleByTypeGetter, evalCtx ports.RuleEvaluationContext) (*bool, error) {
	rule := rg.GetRuleByType(domain.RuleTypeRequired)
	if rule == nil {
		return nil, nil
	}

	result, err := r.evaluator.Evaluate(ctx, rule, evalCtx)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *submissionResolver) shouldValidate(ctx context.Context, rg ruleByTypeGetter, evalCtx ports.RuleEvaluationContext) (bool, error) {
	rule := rg.GetRuleByType(domain.RuleTypeVisible)
	if rule == nil {
		return true, nil
	}

	return r.evaluator.Evaluate(ctx, rule, evalCtx)
}

func valuePath(prefix []int, v *domain.SubmissionValue) []int {
	path := slices.Clone(prefix)
	if v.ItemIndex != nil {
		path = append(path, *v.ItemIndex)
	}

	return path
}
