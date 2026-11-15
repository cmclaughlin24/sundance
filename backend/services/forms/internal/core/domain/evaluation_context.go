package domain

import (
	"maps"
)

type EvaluationContext map[RuleExprSourceType]map[string]any

// BaseEvalContext layers non-repeatable field values over the submission's eval context.
// Fields inside repeatable sections are only visible to rules evaluated within their instance.
func BaseEvalContext(s *Submission, fv *FormVersion) EvaluationContext {
	fields := make(map[string]any)

	for _, page := range fv.GetPages() {
		for _, section := range page.GetSections() {
			if section.IsRepeatable() {
				continue
			}

			for _, instance := range s.GetSection(section.ID).Instances() {
				maps.Copy(fields, section.FieldValues(instance))
			}
		}
	}

	return s.GetEvalContext().WithFields(fields)
}

// WithFields returns a copy of ctx with fields merged over its field source.
// The input context is never mutated.
func (ctx EvaluationContext) WithFields(fields map[string]any) EvaluationContext {
	out := maps.Clone(ctx)
	if out == nil {
		out = make(EvaluationContext)
	}

	merged := maps.Clone(ctx[RuleExprSourceTypeField])
	if merged == nil {
		merged = make(map[string]any, len(fields))
	}

	maps.Copy(merged, fields)
	out[RuleExprSourceTypeField] = merged

	return out
}
