package dto

import "github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"

type ConditionalRuleResponse struct {
	ID domain.ConditionalRuleID `json:"id"`
}

func ConditionalRulesToResponse(rules ...*domain.ConditionalRule) []*ConditionalRuleResponse {
	conditions := make([]*ConditionalRuleResponse, 0, len(rules))
	for _, c := range rules {
		conditions = append(conditions, &ConditionalRuleResponse{
			ID: c.ID,
		})
	}

	return conditions
}
