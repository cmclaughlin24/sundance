package commands

import "sundance/backend/services/forms/internal/core/domain"

type RuleExprSourceData struct {
	Type string
	Key  string
}

type RuleExpressionData struct {
	Source           RuleExprSourceData
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
	TagVersionID   string
	Priority       int
	HasStaticValue bool
	StaticValue    any
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
	Repeat       *SectionRepeatData
	ElementsData []ElementData
	Rules        []RuleData
}

type SectionRepeatData struct {
	MinInstances int
	MaxInstances *int
}

func (s *SectionRepeatData) ToDomain() *domain.SectionRepeat {
	if s == nil {
		return nil
	}

	return &domain.SectionRepeat{
		MinInstances: s.MinInstances,
		MaxInstances: s.MaxInstances,
	}
}

type PageData struct {
	ID           *string
	Key          string
	Name         string
	Position     float32
	SectionsData []SectionData
	Rules        []RuleData
}
