package dto

import "github.com/cmclaughlin24/sundance/forms/internal/core/domain"

type SectionDto struct {
	Key      string     `json:"key"`
	Name     string     `json:"name"`
	Position int        `json:"position"`
	Fields   []FieldDto `json:"fields"`
}

type SectionResponseDto struct {
	ID         domain.SectionID              `json:"id"`
	Key        string                        `json:"key"`
	Name       string                        `json:"name"`
	Position   int                           `json:"position"`
	Fields     []*FieldResponseDto           `json:"fields"`
	Conditions []*ConditionalRuleResponseDto `json:"conditions"`
}

func DtoToSection(dto SectionDto) (*domain.Section, error) {
	fields := make([]*domain.Field, 0, len(dto.Fields))

	for _, f := range dto.Fields {
		field, err := DtoToField(f)

		if err != nil {
			return nil, err
		}

		fields = append(fields, field)
	}

	section, err := domain.NewSection("", dto.Key, dto.Name, dto.Position)
	if err != nil {
		return nil, err
	}

	if err := section.SetFields(fields...); err != nil {
		return nil, err
	}

	return section, nil
}

func SectionToResponseDto(section *domain.Section) *SectionResponseDto {
	if section == nil {
		return nil
	}

	fields := make([]*FieldResponseDto, 0, len(section.Fields))
	for _, f := range section.Fields {
		fields = append(fields, FieldToResponseDto(f))
	}

	conditions := ConditionalRulesToResponseDtos(section.Conditions...)

	return &SectionResponseDto{
		ID:         section.ID,
		Key:        section.Key,
		Name:       section.Name,
		Position:   section.Position,
		Fields:     fields,
		Conditions: conditions,
	}
}
