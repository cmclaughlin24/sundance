package evaluators

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"sundance/backend/pkg/common/stratreg"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"

	"github.com/expr-lang/expr"
)

var (
	ErrInvalidExpression       = errors.New("invalid expression")
	ErrInvalidExpressionOutput = errors.New("invalid expression output")
)

type ExprRuleEvaluator struct {
	logger *slog.Logger
}

func NewExprRuleEvaluator(logger *slog.Logger) *ExprRuleEvaluator {
	return &ExprRuleEvaluator{logger}
}

func (e *ExprRuleEvaluator) Evaluate(ctx context.Context, r *domain.Rule, evalCtx ports.RuleEvaluationContext) (bool, error) {
	stmt, err := e.statement(ctx, r)
	if err != nil {
		return false, err
	}

	program, err := expr.Compile(stmt, expr.AsBool())
	if err != nil {
		e.logger.WarnContext(ctx, "expression compilation failed", "statement", stmt, "error", err)
		return false, fmt.Errorf("%w: %w", ErrInvalidExpression, err)
	}

	output, err := expr.Run(program, evalCtx)
	if err != nil {
		e.logger.ErrorContext(ctx, "expression execution failed", "statement", stmt, "error", err)
		return false, err
	}

	result, ok := output.(bool)
	if !ok {
		e.logger.WarnContext(ctx, "expression output type mismatch", "statement", stmt)
		return false, ErrInvalidExpressionOutput
	}

	return result, nil
}

func (e *ExprRuleEvaluator) statement(ctx context.Context, r *domain.Rule) (string, error) {
	stmt := ""

	for _, re := range r.GetExpressions() {
		statementFn, err := exprRegistry.Get(re.Operator)
		if err != nil {
			e.logger.WarnContext(ctx, "invalid expression operator", "operator", re.Operator)
			return "", domain.ErrInvalidExprOperator
		}

		join, err := joinOperator(re)
		if err != nil {
			e.logger.WarnContext(ctx, "invalid join operator", "operator", *re.JoinWithPrevious)
			return "", domain.ErrInvalidJoinOperator
		}

		stmt = stmt + join + statementFn(re)
	}

	return stmt, nil
}

func joinOperator(re *domain.RuleExpression) (string, error) {
	if re.JoinWithPrevious == nil {
		return "", nil
	}

	operator := ""

	switch *re.JoinWithPrevious {
	case domain.JoinOperatorAnd:
		operator = "&&"
	case domain.JoinOperatorOr:
		operator = "||"
	default:
		return "", domain.ErrInvalidJoinOperator
	}

	return fmt.Sprintf(" %s ", operator), nil
}

type statementFn = func(*domain.RuleExpression) string

var exprRegistry = stratreg.New[domain.ExprOperator, statementFn]().
	Set(domain.ExprOperatorEquals, newDefaultStatementFn("==")).
	Set(domain.ExprOperatorNEquals, newDefaultStatementFn("!=")).
	Set(domain.ExprOperatorLT, newDefaultStatementFn("<")).
	Set(domain.ExprOperatorGT, newDefaultStatementFn(">")).
	Set(domain.ExprOperatorLTE, newDefaultStatementFn("<=")).
	Set(domain.ExprOperatorGTE, newDefaultStatementFn(">="))

func newDefaultStatementFn(operator string) statementFn {
	return func(re *domain.RuleExpression) string {
		return fmt.Sprintf("%s %s %v", re.FieldKey, operator, re.Value)
	}
}
