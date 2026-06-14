package dto

import (
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports/commands"
)

type RuleRequest struct {
	ID          *string                  `json:"id,omitempty" validate:"omitempty,uuidv7"`
	Type        string                   `json:"type" validate:"required"`
	Expressions []*RuleExpressionRequest `json:"expressions" validate:"dive"`
}

type RuleResponse struct {
	ID          domain.RuleID             `json:"id"`
	Type        domain.RuleType           `json:"type"`
	Expressions []*RuleExpressionResponse `json:"expressions"`
}

func RequestsToRuleData(dtos []RuleRequest) []commands.RuleData {
	rules := make([]commands.RuleData, 0, len(dtos))

	for _, dto := range dtos {
		expressions := requestsToRuleExpressionData(dto.Expressions)

		rules = append(rules, commands.RuleData{
			ID:          dto.ID,
			Type:        dto.Type,
			Expressions: expressions,
		})
	}

	return rules
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
