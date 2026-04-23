package dto

import "github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"

type FieldRequest struct {
	Key        string        `json:"key"`
	Name       string        `json:"name"`
	Type       string        `json:"type"`
	Position   int           `json:"position"`
	Attributes any           `json:"attributes"`
	Rules      []RuleRequest `json:"rules"`
}

type FieldResponse struct {
	ID         domain.FieldID  `json:"id"`
	Key        string          `json:"key"`
	Name       string          `json:"name"`
	Type       string          `json:"type"`
	Position   int             `json:"position"`
	Attributes any             `json:"attributes"`
	Rules      []*RuleResponse `json:"rules"`
}

func RequestToField(dto FieldRequest) (*domain.Field, error) {
	attributes, err := attributesFromRequest(domain.FieldType(dto.Type), dto.Attributes)

	if err != nil {
		return nil, err
	}

	rules, err := RequestsToRules(dto.Rules)
	if err != nil {
		return nil, err
	}

	f, err := domain.NewField(
		"",
		dto.Key,
		dto.Name,
		domain.FieldType(dto.Type),
		attributes,
		dto.Position,
	)

	if err != nil {
		return nil, err
	}

	if err := f.SetRules(rules...); err != nil {
		return nil, err
	}

	return f, nil
}

func FieldToResponse(field *domain.Field) *FieldResponse {
	if field == nil {
		return nil
	}

	rules := RuleToResponse(field.GetRules())

	return &FieldResponse{
		ID:         field.ID,
		Key:        field.Key,
		Name:       field.Name,
		Type:       string(field.Type),
		Position:   field.Position,
		Attributes: field.Attributes,
		Rules:      rules,
	}
}
