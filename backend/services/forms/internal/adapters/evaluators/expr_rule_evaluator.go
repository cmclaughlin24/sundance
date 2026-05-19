package evaluators

import (
	"context"

	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

var exprRegistry = map[domain.ExprOperator]string{
	domain.ExprOperatorEquals:  "==",
	domain.ExprOperatorNEquals: "!=",
	domain.ExprOperatorLT:      "<",
	domain.ExprOperatorGT:      ">",
	domain.ExprOperatorLTE:     "<=",
	domain.ExprOperatorGTE:     ">=",
}

type ExprRuleEvaluator struct{}

func (e *ExprRuleEvaluator) Evaluate(context.Context, *domain.Rule, ports.RuleEvaluationContext) (bool, error) {
	return true, nil
}
