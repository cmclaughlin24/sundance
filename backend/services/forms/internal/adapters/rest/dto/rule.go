package dto

import "sundance/backend/services/forms/internal/core/domain"

type RuleRequest struct {
	ID          *string                  `json:"id,omitempty" validate:"uuidv7"`
	Type        string                   `json:"type" validate:"required"`
	Expressions []*RuleExpressionRequest `json:"expressions" validate:"dive"`
}

type RuleResponse struct {
	ID          domain.RuleID             `json:"id"`
	Type        domain.RuleType           `json:"type"`
	Expressions []*RuleExpressionResponse `json:"expressions"`
}

func RequestToRule(dto RuleRequest) (*domain.Rule, error) {
	expressions, err := requestsToRuleExpressions(dto.Expressions)
	if err != nil {
		return nil, err
	}

	r, err := domain.NewRule(domain.RuleType(dto.Type))
	if err != nil {
		return nil, err
	}

	if err := r.AddExpressions(expressions...); err != nil {
		return nil, err
	}

	return r, nil
}

func RequestsToRules(dtos []RuleRequest) ([]*domain.Rule, error) {
	rules := make([]*domain.Rule, 0, len(dtos))

	for _, dto := range dtos {
		r, err := RequestToRule(dto)

		if err != nil {
			return nil, err
		}

		rules = append(rules, r)
	}

	return rules, nil
}

func RulesToResponse(rules map[domain.RuleType]*domain.Rule) []*RuleResponse {
	dtos := make([]*RuleResponse, 0, len(rules))

	for _, r := range rules {
		expressions := ruleExpressionsToResponse(r.GetExpressions())

		dtos = append(dtos, &RuleResponse{
			ID:          r.ID,
			Type:        r.Type,
			Expressions: expressions,
		})
	}

	return dtos
}
