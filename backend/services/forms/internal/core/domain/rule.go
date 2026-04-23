package domain

import (
	"errors"

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
	ErrInvalidRuleType   = errors.New("invalid rule type")
	ErrDuplicateRuleType = errors.New("duplicate rule type")
)

type Rule struct {
	ID         RuleID
	Type       RuleType
	Expression string
}

func NewRule(id RuleID, ruleType RuleType, expression string) (*Rule, error) {
	if !isValidRuleType(ruleType) {
		return nil, ErrInvalidRuleType
	}

	return &Rule{
		ID:         id,
		Type:       ruleType,
		Expression: expression,
	}, nil

}

var isValidRuleType = validate.NewTypeValidator([]RuleType{
	RuleTypeVisible,
	RuleTypeRequired,
	RuleTypeReadOnly,
})

type baseWithRules struct {
	rules map[RuleType]*Rule
}

func (b *baseWithRules) GetRules() map[RuleType]*Rule {
	return b.rules
}

func (b *baseWithRules) GetRule(ruleType RuleType) *Rule {
	r, ok := b.rules[ruleType]

	if !ok {
		return nil
	}

	return r
}

func (b *baseWithRules) SetRules(rules ...*Rule) error {
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
