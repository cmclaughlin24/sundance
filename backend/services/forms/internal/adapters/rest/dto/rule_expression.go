package dto

import "sundance/backend/services/forms/internal/core/domain"

type RuleExpressionRequest struct {
	FieldID          string  `json:"fieldId" validate:"required"`
	Operator         string  `json:"operator" validate:"required"`
	Value            any     `json:"value"`
	JoinWithPrevious *string `json:"joinWithPrevious"`
	Position         float32 `json:"position" validate:"gte=0"`
}

type RuleExpressionResponse struct {
	FieldID          string  `json:"fieldId"`
	Operator         string  `json:"operator"`
	Value            any     `json:"value"`
	JoinWithPrevious *string `json:"joinWithPrevious"`
	Position         float32 `json:"position"`
}

func requestToRuleExpression(dto *RuleExpressionRequest) (*domain.RuleExpression, error) {
	return domain.NewRuleExpression(
		domain.FieldID(dto.FieldID),
		domain.ExpressionOperator(dto.Operator),
		dto.Value,
		(*domain.JoinOperator)(dto.JoinWithPrevious),
		dto.Position,
	)
}

func requestsToRuleExpressions(dtos []*RuleExpressionRequest) ([]*domain.RuleExpression, error) {
	expressions := make([]*domain.RuleExpression, 0, len(dtos))
	for _, dto := range dtos {
		e, err := requestToRuleExpression(dto)
		if err != nil {
			return nil, err
		}

		expressions = append(expressions, e)
	}

	return expressions, nil
}

func ruleExpressionsToResponse(expressions []*domain.RuleExpression) []*RuleExpressionResponse {
	dtos := make([]*RuleExpressionResponse, 0, len(expressions))
	for _, exp := range expressions {
		dtos = append(dtos, &RuleExpressionResponse{
			FieldID:          string(exp.FieldID),
			Operator:         string(exp.Operator),
			Value:            exp.Value,
			JoinWithPrevious: (*string)(exp.JoinWithPrevious),
			Position:         exp.GetPosition(),
		})
	}

	return dtos
}
