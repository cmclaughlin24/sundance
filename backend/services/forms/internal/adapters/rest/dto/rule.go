package dto

import "github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"

type RuleRequest struct {
	Type       string `json:"type"`
	Expression string `json:"expression"`
}

type RuleResponse struct {
	ID         domain.RuleID   `json:"id"`
	Type       domain.RuleType `json:"type"`
	Expression string          `json:"expression"`
}

func RequestToRule(dto RuleRequest) (*domain.Rule, error) {
	return domain.NewRule(
		domain.RuleType(dto.Type),
		dto.Expression,
	)
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

func RuleToResponse(rules map[domain.RuleType]*domain.Rule) []*RuleResponse {
	dtos := make([]*RuleResponse, 0, len(rules))
	for _, r := range rules {
		dtos = append(dtos, &RuleResponse{
			ID:         r.ID,
			Type:       r.Type,
			Expression: r.Expression,
		})
	}

	return dtos
}
