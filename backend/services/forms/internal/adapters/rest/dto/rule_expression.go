package dto

import (
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports/commands"
)

type RuleExpressionRequest struct {
	FieldKey         string  `json:"fieldKey" validate:"required"`
	Operator         string  `json:"operator" validate:"required"`
	Value            any     `json:"value"`
	JoinWithPrevious *string `json:"joinWithPrevious"`
	Position         float32 `json:"position" validate:"gte=0"`
}

type RuleExpressionResponse struct {
	FieldKey         string  `json:"fieldKey"`
	Operator         string  `json:"operator"`
	Value            any     `json:"value"`
	JoinWithPrevious *string `json:"joinWithPrevious"`
	Position         float32 `json:"position"`
}

func requestsToRuleExpressionData(dtos []*RuleExpressionRequest) []*commands.RuleExpressionData {
	expressions := make([]*commands.RuleExpressionData, 0, len(dtos))

	for _, dto := range dtos {
		expressions = append(expressions, &commands.RuleExpressionData{
			FieldKey:         dto.FieldKey,
			Operator:         dto.Operator,
			Value:            dto.Value,
			JoinWithPrevious: dto.JoinWithPrevious,
			Position:         dto.Position,
		})
	}

	return expressions
}

func ruleExpressionsToResponse(expressions []*domain.RuleExpression) []*RuleExpressionResponse {
	dtos := make([]*RuleExpressionResponse, 0, len(expressions))
	for _, exp := range expressions {
		dtos = append(dtos, &RuleExpressionResponse{
			FieldKey:         string(exp.FieldKey),
			Operator:         string(exp.Operator),
			Value:            exp.Value,
			JoinWithPrevious: (*string)(exp.JoinWithPrevious),
			Position:         exp.GetPosition(),
		})
	}

	return dtos
}
