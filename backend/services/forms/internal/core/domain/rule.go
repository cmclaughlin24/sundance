package domain

import (
	"errors"
	"slices"

	"sundance/backend/pkg/common/validate"
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
	expressions PositionElements[*RuleExpression]
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

func (r *Rule) GetExpressions() PositionElements[*RuleExpression] {
	return r.expressions
}

func (r *Rule) AddExpressions(expressions ...*RuleExpression) error {
	if r == nil {
		return ErrInvalidRule
	}

	cpy := slices.Clone(r.expressions)
	cpy = append(cpy, expressions...)

	if ok := hasUniqueElements(cpy); !ok {
		return ErrDuplicatePosition
	}

	sortElements(cpy)
	r.expressions = cpy

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

func (b *withRules) GetRuleByType(ruleType RuleType) *Rule {
	r, ok := b.rules[ruleType]

	if !ok {
		return nil
	}

	return r
}

func (b *withRules) GetRule(ruleID RuleID) *Rule {
	for _, rule := range b.rules {
		if ruleID == rule.ID {
			return rule
		}
	}

	return nil
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

func (b *withRules) ReplaceRules(rules ...*Rule) error {
	cpy := *b
	cpy.rules = make(map[RuleType]*Rule)

	if err := cpy.SetRules(rules...); err != nil {
		return err
	}

	*b = cpy

	return nil
}
