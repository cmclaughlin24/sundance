package dto

import (
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports/commands"
)

type RuleExprSourceRequest struct {
	Type string `json:"type" validate:"required"`
	Key  string `json:"key" validate:"required"`
}

type RuleExpressionRequest struct {
	Source           RuleExprSourceRequest `json:"source" validate:"dive"`
	Operator         string                `json:"operator" validate:"required"`
	Value            any                   `json:"value"`
	JoinWithPrevious *string               `json:"joinWithPrevious"`
	Position         float32               `json:"position" validate:"gte=0"`
}

type RuleExprSourceResponse struct {
	Type string `json:"type"`
	Key  string `json:"key"`
}

type RuleExpressionResponse struct {
	Source           RuleExprSourceResponse `json:"source"`
	Operator         string                 `json:"operator"`
	Value            any                    `json:"value"`
	JoinWithPrevious *string                `json:"joinWithPrevious"`
	Position         float32                `json:"position"`
}

func requestsToRuleExpressionData(dtos []*RuleExpressionRequest) []*commands.RuleExpressionData {
	expressions := make([]*commands.RuleExpressionData, 0, len(dtos))

	for _, dto := range dtos {
		expressions = append(expressions, &commands.RuleExpressionData{
			Source: commands.RuleExprSourceData{
				Type: dto.Source.Type,
				Key:  dto.Source.Key,
			},
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
			Source: RuleExprSourceResponse{
				Type: string(exp.Source.Type),
				Key:  exp.Source.Key,
			},
			Operator:         string(exp.Operator),
			Value:            exp.Value,
			JoinWithPrevious: (*string)(exp.JoinWithPrevious),
			Position:         exp.GetPosition(),
		})
	}

	return dtos
}
