package domain

import (
	"errors"
	"slices"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/validate"
)

type RuleID string

type RuleType string

const (
	RuleTypeVisible  RuleType = "visible"
	RuleTypeRequired RuleType = "required"
	RuleTypeReadOnly RuleType = "readonly"
)

var (
	ErrInvalidRule       = errors.New("invalid rule")
	ErrInvalidRuleType   = errors.New("invalid rule type")
	ErrDuplicateRuleType = errors.New("duplicate rule type")
)

type Rule struct {
	ID          RuleID
	Type        RuleType
	expressions []*RuleExpression
}

func NewRule(ruleType RuleType) (*Rule, error) {
	if !isValidRuleType(ruleType) {
		return nil, ErrInvalidRuleType
	}

	r := &Rule{
		ID:   RuleID(NewID()),
		Type: ruleType,
	}

	if err := validate.ValidateStruct(r); err != nil {
		return nil, err
	}

	return r, nil
}

func HydrateRule(id RuleID, ruleType RuleType) *Rule {
	return &Rule{
		ID:   id,
		Type: ruleType,
	}
}

func (r *Rule) GetExpressions() []*RuleExpression {
	return r.expressions
}

func (r *Rule) SetExpressions(expressions ...*RuleExpression) error {
	if r == nil {
		return ErrInvalidRule
	}

	seen := make(map[float32]struct{}, len(expressions))
	for _, exp := range expressions {
		position := exp.GetPosition()

		if _, exists := seen[position]; exists {
			return ErrDuplicatePosition
		}

		seen[position] = struct{}{}
	}

	r.expressions = slices.Clone(expressions)
	slices.SortFunc(r.expressions, func(a, b *RuleExpression) int {
		if a.GetPosition() < b.GetPosition() {
			return -1
		}
		if a.GetPosition() > b.GetPosition() {
			return 1
		}
		return 0
	})

	return nil
}

var isValidRuleType = validate.NewTypeValidator([]RuleType{
	RuleTypeVisible,
	RuleTypeRequired,
	RuleTypeReadOnly,
})

type withRules struct {
	rules map[RuleType]*Rule
}

func (b *withRules) GetRules() map[RuleType]*Rule {
	return b.rules
}

func (b *withRules) GetRule(ruleType RuleType) *Rule {
	r, ok := b.rules[ruleType]

	if !ok {
		return nil
	}

	return r
}

func (b *withRules) SetRules(rules ...*Rule) error {
	if b.rules == nil {
		b.rules = make(map[RuleType]*Rule)
	}

	for _, rule := range rules {
		_, exists := b.rules[rule.Type]

		if exists {
			return ErrDuplicateRuleType
		}

		b.rules[rule.Type] = rule
	}

	return nil
}
