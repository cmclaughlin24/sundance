package domain

import (
	"errors"

	"sundance/backend/pkg/common/validate"
)

type ExprOperator string

type JoinOperator string

type RuleExprSourceType string

const (
	ExprOperatorEquals  ExprOperator = "equal"
	ExprOperatorNEquals ExprOperator = "nequal"
	ExprOperatorLT      ExprOperator = "lt"
	ExprOperatorGT      ExprOperator = "gt"
	ExprOperatorLTE     ExprOperator = "lte"
	ExprOperatorGTE     ExprOperator = "gte"

	JoinOperatorAnd JoinOperator = "and"
	JoinOperatorOr  JoinOperator = "or"

	RuleExprSourceTypeField     RuleExprSourceType = "field"
	RuleExprSourceTypeUserClaim RuleExprSourceType = "user_claim"
)

var (
	ErrInvalidExprOperator       = errors.New("invalid expression operator")
	ErrInvalidJoinOperator       = errors.New("invalid join operator")
	ErrInvalidRuleExprSourceType = errors.New("invalid rule expression source type")
)

type RuleExprSource struct {
	Type RuleExprSourceType
	Key  string
}

type RuleExpression struct {
	Source           RuleExprSource
	Operator         ExprOperator
	Value            any
	JoinWithPrevious *JoinOperator
	withPosition
}

func NewRuleExpression(
	source RuleExprSource,
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

	if !isValidExprSourceType(source.Type) {
		return nil, ErrInvalidRuleExprSourceType
	}

	return &RuleExpression{
		Source:           source,
		Operator:         operator,
		Value:            value,
		JoinWithPrevious: joinWithPrevious,
		withPosition: withPosition{
			position: position,
		},
	}, nil
}

func HydrateRuleExpression(
	source RuleExprSource,
	operator ExprOperator,
	value any,
	joinWithPrevious *JoinOperator,
	position float32,
) *RuleExpression {
	return &RuleExpression{
		Source:           source,
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

var isValidExprSourceType = validate.NewTypeValidator([]RuleExprSourceType{
	RuleExprSourceTypeField,
	RuleExprSourceTypeUserClaim,
})
