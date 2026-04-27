package domain

import (
	"errors"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/validate"
	"github.com/google/uuid"
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
	Expression string `validate:"required,nowhitespace"`
}

func NewRule(ruleType RuleType, expression string) (*Rule, error) {
	if !isValidRuleType(ruleType) {
		return nil, ErrInvalidRuleType
	}

	r := &Rule{
		ID:         RuleID(uuid.NewString()),
		Type:       ruleType,
		Expression: expression,
	}

	if err := validate.ValidateStruct(r); err != nil {
		return nil, err
	}

	return r, nil
}

func HydrateRule(id RuleID, ruleType RuleType, expression string) *Rule {
	return &Rule{
		ID:         id,
		Type:       ruleType,
		Expression: expression,
	}
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
