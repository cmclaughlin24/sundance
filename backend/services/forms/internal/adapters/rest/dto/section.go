package dto

import (
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports/commands"
)

type SectionRequest struct {
	ID       *string               `json:"id,omitempty" validate:"omitempty,uuidv7"`
	Key      string                `json:"key" validate:"required,max=25"`
	Name     string                `json:"name" validate:"required,max=75"`
	Position float32               `json:"position" validate:"gte=0,lte=10"`
	Repeat   *SectionRepeatRequest `json:"repeat,omitempty"`
	Elements []ElementRequest      `json:"elements" validate:"dive"`
	Rules    []RuleRequest         `json:"rules" validate:"dive"`
}

type SectionRepeatRequest struct {
	MinInstances int  `json:"minInstances,omitempty" validate:"omitempty,min=0"`
	MaxInstances *int `json:"maxInstances,omitempty" validate:"omitempty,min=1,gtefield=MinInstances"`
}

type SectionResponse struct {
	ID       domain.SectionID       `json:"id"`
	Key      string                 `json:"key"`
	Name     string                 `json:"name"`
	Position float32                `json:"position"`
	Repeat   *SectionRepeatResponse `json:"repeat"`
	Elements []*ElementResponse     `json:"elements"`
	Rules    []*RuleResponse        `json:"rules"`
}

type SectionRepeatResponse struct {
	MinInstances int  `json:"minInstances"`
	MaxInstances *int `json:"maxInstances"`
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

	var repeat *commands.SectionRepeatData
	if dto.Repeat != nil {
		repeat = &commands.SectionRepeatData{
			MinInstances: dto.Repeat.MinInstances,
			MaxInstances: dto.Repeat.MaxInstances,
		}
	}

	return commands.SectionData{
		ID:           dto.ID,
		Key:          dto.Key,
		Name:         dto.Name,
		Position:     dto.Position,
		Repeat:       repeat,
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

	var repeat *SectionRepeatResponse
	if section.IsRepeatable() {
		r := section.GetRepeat()
		repeat = &SectionRepeatResponse{
			MinInstances: r.MinInstances,
			MaxInstances: r.MaxInstances,
		}
	}

	return &SectionResponse{
		ID:       section.ID,
		Key:      section.Key,
		Name:     section.Name,
		Position: section.GetPosition(),
		Repeat:   repeat,
		Elements: dtos,
		Rules:    rules,
	}
}
