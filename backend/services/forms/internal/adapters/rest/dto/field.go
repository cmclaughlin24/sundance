package dto

import (
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports/commands"
)

type FieldRequest struct {
	ID         *string                        `json:"id,omitempty" validate:"omitempty,uuidv7"`
	Key        string                         `json:"key" validate:"required,max=25"`
	Name       string                         `json:"name" validate:"required,max=250"`
	Type       string                         `json:"type" validate:"required"`
	Position   float32                        `json:"position" validate:"gte=0,lte=50"`
	Attributes any                            `json:"attributes" validate:"required" swaggertype:"object"`
	Tags       []upsertFieldTagMappingRequest `json:"tags" validate:"dive"`
	Rules      []RuleRequest                  `json:"rules" validate:"dive"`
}

type FieldResponse struct {
	ID         domain.FieldID             `json:"id"`
	Key        string                     `json:"key"`
	Name       string                     `json:"name"`
	Type       string                     `json:"type"`
	Position   float32                    `json:"position"`
	Attributes any                        `json:"attributes" swaggertype:"object"`
	Tags       []*FieldTagMappingResponse `json:"tags"`
	Rules      []*RuleResponse            `json:"rules"`
}

func RequestToFieldData(dto FieldRequest) (commands.FieldData, error) {
	attributes, err := attributesFromRequest(domain.FieldType(dto.Type), dto.Attributes)
	if err != nil {
		return commands.FieldData{}, err
	}

	rules := RequestsToRuleData(dto.Rules)
	tags := requestToFieldTagMappingData(dto.Tags)

	return commands.FieldData{
		ID:         dto.ID,
		Key:        dto.Key,
		Name:       dto.Name,
		Type:       dto.Type,
		Position:   dto.Position,
		Attributes: attributes,
		Tags:       tags,
		Rules:      rules,
	}, nil
}

func FieldToResponse(field *domain.Field) *FieldResponse {
	if field == nil {
		return nil
	}

	attr := fieldAttributesToResponse(field.Attributes)
	tags := fieldTagMappingsToResponses(field.GetTags())
	rules := RulesToResponse(field.GetRules())

	return &FieldResponse{
		ID:         field.ID,
		Key:        field.Key,
		Name:       field.Name,
		Type:       string(field.Type),
		Position:   field.GetPosition(),
		Attributes: attr,
		Tags:       tags,
		Rules:      rules,
	}
}
