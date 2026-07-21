package dto

import (
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports/commands"
)

type ElementRequest struct {
	ID         *string                          `json:"id,omitempty" validate:"omitempty,uuidv7"`
	Key        string                           `json:"key" validate:"required,max=25"`
	Name       string                           `json:"name" validate:"required,max=250"`
	Type       string                           `json:"type" validate:"required"`
	Position   float32                          `json:"position" validate:"gte=0,lte=50"`
	Attributes any                              `json:"attributes" validate:"required" swaggertype:"object"`
	Tags       []upsertElementTagMappingRequest `json:"tags" validate:"dive"`
	Rules      []RuleRequest                    `json:"rules" validate:"dive"`
}

type ElementResponse struct {
	ID         domain.ElementID              `json:"id"`
	Key        string                        `json:"key"`
	Name       string                        `json:"name"`
	Type       string                        `json:"type"`
	Position   float32                       `json:"position"`
	Attributes any                           `json:"attributes" swaggertype:"object"`
	Tags       []*ElementTagMappingResponse  `json:"tags"`
	Rules      []*RuleResponse               `json:"rules"`
}

func RequestToElementData(dto ElementRequest) (commands.ElementData, error) {
	attributes, err := attributesFromRequest(domain.ElementType(dto.Type), dto.Attributes)
	if err != nil {
		return commands.ElementData{}, err
	}

	rules := RequestsToRuleData(dto.Rules)
	tags := requestToElementTagMappingData(dto.Tags)

	return commands.ElementData{
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

func ElementToResponse(element *domain.Element) *ElementResponse {
	if element == nil {
		return nil
	}

	attr := elementAttributesToResponse(element.Attributes)
	tags := elementTagMappingsToResponses(element.GetTags())
	rules := RulesToResponse(element.GetRules())

	return &ElementResponse{
		ID:         element.ID,
		Key:        element.Key,
		Name:       element.Name,
		Type:       string(element.Type),
		Position:   element.GetPosition(),
		Attributes: attr,
		Tags:       tags,
		Rules:      rules,
	}
}
