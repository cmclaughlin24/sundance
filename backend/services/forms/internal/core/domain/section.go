package domain

import (
	"errors"
	"fmt"
	"slices"

	"sundance/backend/pkg/common/validate"
)

type SectionID string

var (
	ErrInvalidSection     = errors.New("invalid section")
	ErrInvalidSectionRule = errors.New("invalid section rule type")

	sectionRuleTypes = map[RuleType]bool{
		RuleTypeVisible: true,
	}
)

type Section struct {
	ID     SectionID
	Key    string `validate:"required,nowhitespace"`
	Name   string `validate:"required"`
	fields PositionElements[*Field]
	withPosition
	withRules
}

func NewSection(key, name string, position float32) (*Section, error) {
	if !isValidPosition(position) {
		return nil, ErrInvalidPosition
	}

	s := &Section{
		ID:     SectionID(NewID()),
		Key:    key,
		Name:   name,
		fields: make(PositionElements[*Field], 0),
		withPosition: withPosition{
			position: position,
		},
	}

	if err := validate.ValidateStruct(s); err != nil {
		return nil, err
	}

	return s, nil
}

func HydrateSection(id SectionID, key, name string, position float32) *Section {
	return &Section{
		ID:     id,
		Key:    key,
		Name:   name,
		fields: make(PositionElements[*Field], 0),
		withPosition: withPosition{
			position: position,
		},
	}
}

func (s *Section) GetFields() PositionElements[*Field] {
	return s.fields
}

func (s *Section) AddFields(fields ...*Field) error {
	if s == nil {
		return ErrInvalidSection
	}

	cpy := slices.Clone(s.fields)
	cpy = append(cpy, fields...)

	if ok := hasUniqueElements(cpy); !ok {
		return ErrDuplicatePosition
	}

	sortElements(cpy)
	s.fields = cpy

	return nil
}

func (s *Section) ReplaceFields(fields ...*Field) error {
	if s == nil {
		return ErrInvalidSection
	}

	old := s.fields
	s.fields = make(PositionElements[*Field], 0)

	if err := s.AddFields(fields...); err != nil {
		s.fields = old
		return err
	}

	return nil
}

func (s *Section) SetRules(rules ...*Rule) error {
	for _, rule := range rules {
		if allow, ok := sectionRuleTypes[rule.Type]; !allow || !ok {
			return fmt.Errorf("rule type %s not supported; %w", rule.Type, ErrInvalidSectionRule)
		}
	}

	return s.withRules.SetRules(rules...)
}
