package dto

import "github.com/cmclaughlin24/sundance/forms/internal/core/domain"

type FieldDto struct {
	Key        string `json:"key"`
	Name       string `json:"name"`
	FieldType  string `json:"fieldType"`
	Position   int    `json:"position"`
	Attributes any    `json:"attributes"`
}

type FieldResponseDto struct {
	ID         domain.FieldID                `json:"id"`
	Key        string                        `json:"key"`
	Name       string                        `json:"name"`
	FieldType  string                        `json:"fieldType"`
	Position   int                           `json:"position"`
	Conditions []*ConditionalRuleResponseDto `json:"conditions"`
}

func DtoToField(dto FieldDto) (*domain.Field, error) {
	attributes, err := attributesFromDto(domain.FieldType(dto.FieldType), dto.Attributes)

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

func FieldToResponseDto(field *domain.Field) *FieldResponseDto {
	if field == nil {
		return nil
	}

	conditions := ConditionalRulesToResponseDtos(field.Conditions...)

	return &FieldResponseDto{
		ID:         field.ID,
		Key:        field.Key,
		Name:       field.Name,
		FieldType:  string(field.FieldType),
		Position:   field.Position,
		Conditions: conditions,
	}
}
