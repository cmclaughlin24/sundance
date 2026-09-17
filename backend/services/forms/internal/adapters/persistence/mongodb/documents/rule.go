package documents

import (
	"sundance/backend/services/forms/internal/core/domain"
)

type ruleDocument struct {
	ID          string                    `bson:"_id"`
	Type        string                    `bson:"type"`
	Expressions []*ruleExpressionDocument `bson:"expressions"`
}

func RulesToDocuments(rules map[domain.RuleType]*domain.Rule) []*ruleDocument {
	documents := make([]*ruleDocument, 0, len(rules))
	for _, r := range rules {
		documents = append(documents, toRuleDocument(r))
	}
	return documents
}

func toRuleDocument(r *domain.Rule) *ruleDocument {
	expressions := r.GetExpressions()
	documents := make([]*ruleExpressionDocument, 0, len(expressions))

	for _, e := range expressions {
		documents = append(documents, toRuleExpressionDocument(e))
	}

	return &ruleDocument{
		ID:          string(r.ID),
		Type:        string(r.Type),
		Expressions: documents,
	}
}

func documentsToRules(documents []*ruleDocument) ([]*domain.Rule, error) {
	rules := make([]*domain.Rule, 0, len(documents))
	for _, d := range documents {
		rule, err := fromRuleDocument(d)
		if err != nil {
			return nil, err
		}

		rules = append(rules, rule)
	}

	return rules, nil
}

func fromRuleDocument(doc *ruleDocument) (*domain.Rule, error) {
	r := domain.HydrateRule(
		domain.RuleID(doc.ID),
		domain.RuleType(doc.Type),
	)

	expressions := make([]*domain.RuleExpression, 0, len(doc.Expressions))
	for _, e := range doc.Expressions {
		expressions = append(expressions, fromRuleExpressionDocument(e))
	}

	if err := r.AddExpressions(expressions...); err != nil {
		return nil, err
	}

	return r, nil
}

type ruleExprSource struct {
	Type string `bson:"type"`
	Key  string `bson:"key"`
}

type ruleExpressionDocument struct {
	Source           ruleExprSource `bson:"source"`
	Operator         string         `bson:"operator"`
	Value            any            `bson:"value"`
	JoinWithPrevious *string        `bson:"join_with_previous"`
	Position         float32        `bson:"position"`
}

func toRuleExpressionDocument(e *domain.RuleExpression) *ruleExpressionDocument {
	return &ruleExpressionDocument{
		Source: ruleExprSource{
			Type: string(e.Source.Type),
			Key:  e.Source.Key,
		},
		Operator:         string(e.Operator),
		Value:            e.Value,
		JoinWithPrevious: (*string)(e.JoinWithPrevious),
		Position:         e.GetPosition(),
	}
}

func fromRuleExpressionDocument(e *ruleExpressionDocument) *domain.RuleExpression {
	return domain.HydrateRuleExpression(
		domain.RuleExprSource{
			Type: domain.RuleExprSourceType(e.Source.Type),
			Key:  e.Source.Key,
		},
		domain.ExprOperator(e.Operator),
		e.Value,
		(*domain.JoinOperator)(e.JoinWithPrevious),
		e.Position,
	)
}
