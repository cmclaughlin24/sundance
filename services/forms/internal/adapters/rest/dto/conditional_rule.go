package dto

import "github.com/cmclaughlin24/sundance/forms/internal/core/domain"

type ConditionalRuleResponseDto struct {
	ID domain.ConditionalRuleID `json:"id"`
}

func ConditionalRulesToResponseDtos(rules ...*domain.ConditionalRule) []*ConditionalRuleResponseDto {
	conditions := make([]*ConditionalRuleResponseDto, 0, len(rules))
	for _, c := range rules {
		conditions = append(conditions, &ConditionalRuleResponseDto{
			ID: c.ID,
		})
	}

	return conditions
}
