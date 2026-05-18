package evaluators

import (
	"context"

	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/ports"
)

type ExprRuleEvaluator struct{}

func (e *ExprRuleEvaluator) Evaluate(context.Context, *domain.Rule, ports.RuleEvaluationContext) (bool, error) {
	return true, nil
}
