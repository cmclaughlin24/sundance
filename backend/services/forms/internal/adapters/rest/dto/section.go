package dto

import (
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports/commands"
)

type SectionRequest struct {
	ID       *string        `json:"id,omitempty" validate:"omitempty,uuidv7"`
	Key      string         `json:"key" validate:"required,max=25"`
	Name     string         `json:"name" validate:"required,max=75"`
	Position float32        `json:"position" validate:"gte=0,lte=10"`
	Fields   []FieldRequest `json:"fields" validate:"dive"`
	Rules    []RuleRequest  `json:"rules" validate:"dive"`
}

type SectionResponse struct {
	ID       domain.SectionID `json:"id"`
	Key      string           `json:"key"`
	Name     string           `json:"name"`
	Position float32          `json:"position"`
	Fields   []*FieldResponse `json:"fields"`
	Rules    []*RuleResponse  `json:"rules"`
}

func RequestToSectionData(dto SectionRequest) (commands.SectionData, error) {
	fields := make([]commands.FieldData, 0, len(dto.Fields))

	for _, f := range dto.Fields {
		field, err := RequestToFieldData(f)
		if err != nil {
			return commands.SectionData{}, err
		}

		fields = append(fields, field)
	}

	rules := RequestsToRuleData(dto.Rules)

	return commands.SectionData{
		ID:         dto.ID,
		Key:        dto.Key,
		Name:       dto.Name,
		Position:   dto.Position,
		FieldsData: fields,
		Rules:      rules,
	}, nil
}

func SectionToResponse(section *domain.Section) *SectionResponse {
	if section == nil {
		return nil
	}

	fields := section.GetFields()
	dtos := make([]*FieldResponse, 0, len(fields))

	for _, f := range fields {
		dtos = append(dtos, FieldToResponse(f))
	}

	rules := RulesToResponse(section.GetRules())

	return &SectionResponse{
		ID:       section.ID,
		Key:      section.Key,
		Name:     section.Name,
		Position: section.GetPosition(),
		Fields:   dtos,
		Rules:    rules,
	}
}
