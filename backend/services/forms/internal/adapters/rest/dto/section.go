package dto

import "github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"

type SectionRequest struct {
	Key      string         `json:"key"`
	Name     string         `json:"name"`
	Position int            `json:"position"`
	Fields   []FieldRequest `json:"fields"`
	Rules    []RuleRequest  `json:"rules"`
}

type SectionResponse struct {
	ID       domain.SectionID `json:"id"`
	Key      string           `json:"key"`
	Name     string           `json:"name"`
	Position int              `json:"position"`
	Fields   []*FieldResponse `json:"fields"`
	Rules    []*RuleResponse  `json:"rules"`
}

func RequestToSection(dto SectionRequest) (*domain.Section, error) {
	fields := make([]*domain.Field, 0, len(dto.Fields))

	for _, f := range dto.Fields {
		field, err := RequestToField(f)

		if err != nil {
			return nil, err
		}

		fields = append(fields, field)
	}

	rules, err := RequestsToRules(dto.Rules)
	if err != nil {
		return nil, err
	}

	section, err := domain.NewSection(dto.Key, dto.Name, dto.Position)
	if err != nil {
		return nil, err
	}

	if err := section.SetFields(fields...); err != nil {
		return nil, err
	}

	if err := section.SetRules(rules...); err != nil {
		return nil, err
	}

	return section, nil
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

	rules := RuleToResponse(section.GetRules())

	return &SectionResponse{
		ID:       section.ID,
		Key:      section.Key,
		Name:     section.Name,
		Position: section.Position,
		Fields:   dtos,
		Rules:    rules,
	}
}
