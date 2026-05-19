package domain

import (
	"errors"

	"sundance/backend/pkg/common/validate"
)

type ExprOperator string

type JoinOperator string

const (
	ExprOperatorEquals  ExprOperator = "equal"
	ExprOperatorNEquals ExprOperator = "nequal"
	ExprOperatorLT      ExprOperator = "lt"
	ExprOperatorGT      ExprOperator = "gt"
	ExprOperatorLTE     ExprOperator = "lte"
	ExprOperatorGTE     ExprOperator = "gte"

	JoinOperatorAnd JoinOperator = "and"
	JoinOperatorOr  JoinOperator = "or"
)

var (
	ErrInvalidExprOperator = errors.New("invalid expression operator")
	ErrInvalidJoinOperator = errors.New("invalid join operator")
)

type RuleExpression struct {
	FieldKey         string
	Operator         ExprOperator
	Value            any
	JoinWithPrevious *JoinOperator
	withPosition
}

func NewRuleExpression(
	fieldID string,
	operator ExprOperator,
	value any,
	joinWithPrevious *JoinOperator,
	position float32,
) (*RuleExpression, error) {
	if !isValidExprOperator(operator) {
		return nil, ErrInvalidExprOperator
	}

	if joinWithPrevious != nil && !isValidJoinOperator(*joinWithPrevious) {
		return nil, ErrInvalidJoinOperator
	}

	return &RuleExpression{
		FieldKey:         fieldID,
		Operator:         operator,
		Value:            value,
		JoinWithPrevious: joinWithPrevious,
		withPosition: withPosition{
			position: position,
		},
	}, nil
}

func HydrateRuleExpression(
	fieldID string,
	operator ExprOperator,
	value any,
	joinWithPrevious *JoinOperator,
	position float32,
) *RuleExpression {
	return &RuleExpression{
		FieldKey:         fieldID,
		Operator:         operator,
		Value:            value,
		JoinWithPrevious: joinWithPrevious,
		withPosition: withPosition{
			position: position,
		},
	}
}

var isValidJoinOperator = validate.NewTypeValidator([]JoinOperator{
	JoinOperatorAnd,
	JoinOperatorOr,
})

var isValidExprOperator = validate.NewTypeValidator([]ExprOperator{
	ExprOperatorEquals,
	ExprOperatorNEquals,
	ExprOperatorLT,
	ExprOperatorGT,
	ExprOperatorLTE,
	ExprOperatorGTE,
})
