package domain

import (
	"errors"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/validate"
)

type ExpressionOperator string

type JoinOperator string

const (
	JoinOperatorAnd JoinOperator = "and"
	JoinOperatorOr  JoinOperator = "or"
)

var (
	ErrInvalidJoinOperator = errors.New("invalid join operator")
)

type RuleExpression struct {
	FieldID          FieldID
	Operator         ExpressionOperator
	Value            any
	JoinWithPrevious *JoinOperator
	withPosition
}

func NewRuleExpression(
	fieldID FieldID,
	operator ExpressionOperator,
	value any,
	joinWithPrevious *JoinOperator,
	position float32,
) (*RuleExpression, error) {
	// TODO: Implement validation for operators.

	if joinWithPrevious != nil && !isValidJoinOperator(*joinWithPrevious) {
		return nil, ErrInvalidJoinOperator
	}

	return &RuleExpression{
		FieldID:          fieldID,
		Operator:         operator,
		Value:            value,
		JoinWithPrevious: joinWithPrevious,
		withPosition: withPosition{
			position: position,
		},
	}, nil
}

func HydrateRuleExpression(
	fieldID FieldID,
	operator ExpressionOperator,
	value any,
	joinWithPrevious *JoinOperator,
	position float32,
) *RuleExpression {
	return &RuleExpression{
		FieldID:          fieldID,
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
