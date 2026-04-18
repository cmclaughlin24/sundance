package dto

import "github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"

type FieldRequest struct {
	Key        string `json:"key"`
	Name       string `json:"name"`
	FieldType  string `json:"fieldType"`
	Position   int    `json:"position"`
	Attributes any    `json:"attributes"`
}

type FieldResponse struct {
	ID         domain.FieldID             `json:"id"`
	Key        string                     `json:"key"`
	Name       string                     `json:"name"`
	FieldType  string                     `json:"fieldType"`
	Position   int                        `json:"position"`
	Conditions []*ConditionalRuleResponse `json:"conditions"`
}

func RequestToField(dto FieldRequest) (*domain.Field, error) {
	attributes, err := attributesFromRequest(domain.FieldType(dto.FieldType), dto.Attributes)

	if err != nil {
		return nil, err
	}

	return domain.NewField(
		"",
		dto.Key,
		dto.Name,
		domain.FieldType(dto.FieldType),
		attributes,
		dto.Position,
	)
}

func FieldToResponse(field *domain.Field) *FieldResponse {
	if field == nil {
		return nil
	}

	conditions := ConditionalRulesToResponse(field.Conditions...)

	return &FieldResponse{
		ID:         field.ID,
		Key:        field.Key,
		Name:       field.Name,
		FieldType:  string(field.FieldType),
		Position:   field.Position,
		Conditions: conditions,
	}
}
