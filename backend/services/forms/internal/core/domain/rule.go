package domain

import "errors"

type RuleID string

type RuleType string

const (
	RuleTypeVisible  RuleType = "visible"
	RuleTypeRequired RuleType = "required"
	RuleTypeReadOnly RuleType = "readonly"
)

var (
	ErrDuplicateRuleType = errors.New("duplicate rule type")
)

type Rule struct {
	ID         RuleID
	Type       RuleType
	Expression string
}

func NewRule(id RuleID, ruleType RuleType, expression string) (*Rule, error) {
	cr := &Rule{
		ID:         id,
		Type:       ruleType,
		Expression: expression,
	}

	// TODO: Implement domain validation.

	return cr, nil
}

type baseWithRules struct {
	Rules map[RuleType]*Rule
}

func (b *baseWithRules) SetRules(rules ...*Rule) error {
	if b.Rules == nil {
		b.Rules = make(map[RuleType]*Rule)
	}

	for _, rule := range rules {
		_, exists := b.Rules[rule.Type]

		if exists {
			return ErrDuplicateRuleType
		}

		b.Rules[rule.Type] = rule
	}

	return nil
}
