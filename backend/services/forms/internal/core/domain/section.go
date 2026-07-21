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
	ID       SectionID
	Key      string `validate:"required,nowhitespace"`
	Name     string `validate:"required"`
	elements PositionElements[*Element]
	withPosition
	withRules
}

func NewSection(key, name string, position float32) (*Section, error) {
	if !isValidPosition(position) {
		return nil, ErrInvalidPosition
	}

	s := &Section{
		ID:       SectionID(NewID()),
		Key:      key,
		Name:     name,
		elements: make(PositionElements[*Element], 0),
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
		ID:       id,
		Key:      key,
		Name:     name,
		elements: make(PositionElements[*Element], 0),
		withPosition: withPosition{
			position: position,
		},
	}
}

func (s *Section) Update(key, name string, position float32) error {
	if s == nil {
		return ErrInvalidSection
	}

	if !isValidPosition(position) {
		return ErrInvalidPosition
	}

	cpy := *s
	cpy.Key = key
	cpy.Name = name
	cpy.position = position

	if err := validate.ValidateStruct(cpy); err != nil {
		return err
	}

	*s = cpy

	return nil
}

func (s *Section) GetElements() PositionElements[*Element] {
	return s.elements
}

func (s *Section) GetElement(elementID ElementID) *Element {
	idx := slices.IndexFunc(s.elements, func(e *Element) bool {
		return e.ID == elementID
	})

	if idx == -1 {
		return nil
	}

	return s.elements[idx]
}

func (s *Section) AddElements(elements ...*Element) error {
	if s == nil {
		return ErrInvalidSection
	}

	cpy := slices.Clone(s.elements)
	cpy = append(cpy, elements...)

	if ok := hasUniqueElements(cpy); !ok {
		return ErrDuplicatePosition
	}

	sortElements(cpy)
	s.elements = cpy

	return nil
}

func (s *Section) ReplaceElements(elements ...*Element) error {
	if s == nil {
		return ErrInvalidSection
	}

	old := s.elements
	s.elements = make(PositionElements[*Element], 0)

	if err := s.AddElements(elements...); err != nil {
		s.elements = old
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
