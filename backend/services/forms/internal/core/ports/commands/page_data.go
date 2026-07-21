package commands

import "sundance/backend/services/forms/internal/core/domain"

type RuleExpressionData struct {
	FieldKey         string
	Operator         string
	Value            any
	JoinWithPrevious *string
	Position         float32
}

type RuleData struct {
	ID          *string
	Type        string
	Expressions []*RuleExpressionData
}

type ElementTagMappingData struct {
	TagVersionID string
	Priority     int
}

type ElementData struct {
	ID         *string
	Key        string
	Name       string
	Type       string
	Position   float32
	Attributes domain.ElementAttributes
	Tags       []ElementTagMappingData
	Rules      []RuleData
}

type SectionData struct {
	ID           *string
	Key          string
	Name         string
	Position     float32
	ElementsData []ElementData
	Rules        []RuleData
}

type PageData struct {
	ID           *string
	Key          string
	Name         string
	Position     float32
	SectionsData []SectionData
	Rules        []RuleData
}
