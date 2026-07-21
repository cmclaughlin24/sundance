package dto

import (
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports/commands"
)

type SectionRequest struct {
	ID       *string          `json:"id,omitempty" validate:"omitempty,uuidv7"`
	Key      string           `json:"key" validate:"required,max=25"`
	Name     string           `json:"name" validate:"required,max=75"`
	Position float32          `json:"position" validate:"gte=0,lte=10"`
	Elements []ElementRequest `json:"elements" validate:"dive"`
	Rules    []RuleRequest    `json:"rules" validate:"dive"`
}

type SectionResponse struct {
	ID       domain.SectionID   `json:"id"`
	Key      string             `json:"key"`
	Name     string             `json:"name"`
	Position float32            `json:"position"`
	Elements []*ElementResponse `json:"elements"`
	Rules    []*RuleResponse    `json:"rules"`
}

func RequestToSectionData(dto SectionRequest) (commands.SectionData, error) {
	elements := make([]commands.ElementData, 0, len(dto.Elements))

	for _, e := range dto.Elements {
		element, err := RequestToElementData(e)
		if err != nil {
			return commands.SectionData{}, err
		}

		elements = append(elements, element)
	}

	rules := RequestsToRuleData(dto.Rules)

	return commands.SectionData{
		ID:           dto.ID,
		Key:          dto.Key,
		Name:         dto.Name,
		Position:     dto.Position,
		ElementsData: elements,
		Rules:        rules,
	}, nil
}

func SectionToResponse(section *domain.Section) *SectionResponse {
	if section == nil {
		return nil
	}

	elements := section.GetElements()
	dtos := make([]*ElementResponse, 0, len(elements))

	for _, e := range elements {
		dtos = append(dtos, ElementToResponse(e))
	}

	rules := RulesToResponse(section.GetRules())

	return &SectionResponse{
		ID:       section.ID,
		Key:      section.Key,
		Name:     section.Name,
		Position: section.GetPosition(),
		Elements: dtos,
		Rules:    rules,
	}
}
