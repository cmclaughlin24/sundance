package domain

import (
	"cmp"
	"errors"
	"fmt"
	"slices"

	"sundance/backend/pkg/common/validate"
)

type SectionID string

var (
	ErrInvalidSection       = errors.New("invalid section")
	ErrInvalidSectionRule   = errors.New("invalid section rule type")
	ErrSectionInstanceCount = errors.New("section instance count out of range")

	sectionRuleTypes = map[RuleType]bool{
		RuleTypeVisible: true,
	}
)

type Section struct {
	ID       SectionID
	Key      string `validate:"required,nowhitespace"`
	Name     string `validate:"required"`
	repeat   *SectionRepeat
	elements PositionElements[*Element]
	withPosition
	withRules
}

type SectionRepeat struct {
	MinInstances int     `validate:"min=0"`
	MaxInstances *int    `validate:"omitempty,min=1,gtefield=MinInstances"`
	AddLabel     *string `validate:"omitempty,min=1,max=100"`
}

func NewSection(key, name string, position float32, repeat *SectionRepeat) (*Section, error) {
	if !isValidPosition(position) {
		return nil, ErrInvalidPosition
	}

	s := &Section{
		ID:       SectionID(NewID()),
		Key:      key,
		Name:     name,
		repeat:   cloneSectionRepeat(repeat),
		elements: make(PositionElements[*Element], 0),
		withPosition: withPosition{
			position: position,
		},
	}

	if err := s.validate(); err != nil {
		return nil, err
	}

	return s, nil
}

func HydrateSection(id SectionID, key, name string, position float32, repeat *SectionRepeat) *Section {
	return &Section{
		ID:       id,
		Key:      key,
		Name:     name,
		repeat:   cloneSectionRepeat(repeat),
		elements: make(PositionElements[*Element], 0),
		withPosition: withPosition{
			position: position,
		},
	}
}

func (s *Section) Update(key, name string, position float32, repeat *SectionRepeat) error {
	if s == nil {
		return ErrInvalidSection
	}

	if !isValidPosition(position) {
		return ErrInvalidPosition
	}

	cpy := *s
	cpy.Key = key
	cpy.Name = name
	cpy.repeat = cloneSectionRepeat(repeat)
	cpy.position = position

	if err := cpy.validate(); err != nil {
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

func (s *Section) GetRepeat() *SectionRepeat {
	if s == nil {
		return nil
	}

	return cloneSectionRepeat(s.repeat)
}

func (s *Section) IsRepeatable() bool {
	return s != nil && s.repeat != nil
}

func (s *Section) ValidateInstanceCount(n int) error {
	if s == nil {
		return ErrInvalidSection
	}

	if !s.IsRepeatable() {
		if n > 1 {
			return fmt.Errorf("%w; section=%s max=1 got=%d", ErrSectionInstanceCount, s.Key, n)
		}
		return nil
	}

	if n < s.repeat.MinInstances {
		return fmt.Errorf("%w; section=%s min=%d got=%d", ErrSectionInstanceCount, s.Key, s.repeat.MinInstances, n)
	}

	if s.repeat.MaxInstances != nil && n > *s.repeat.MaxInstances {
		return fmt.Errorf("%w; section=%s max=%d got=%d", ErrSectionInstanceCount, s.Key, *s.repeat.MaxInstances, n)
	}

	return nil
}

func (s *Section) FieldValues(instance SectionInstance) map[string]any {
	fields := make(map[string]any)

	for _, element := range s.GetElements() {
		values := instance.Values(element.ID)

		switch {
		case len(values) == 0:
			continue
		case values[0].ItemIndex == nil:
			fields[element.Key] = values[0].Value
		default:
			slices.SortFunc(values, func(a, b *SubmissionValue) int {
				return cmp.Compare(*a.ItemIndex, *b.ItemIndex)
			})

			list := make([]any, 0, len(values))
			for _, v := range values {
				list = append(list, v.Value)
			}
			fields[element.Key] = list
		}
	}

	return fields
}

func (s *Section) validate() error {
	if err := validate.ValidateStruct(s); err != nil {
		return err
	}

	return s.repeat.validate()
}

func (sr *SectionRepeat) validate() error {
	if sr == nil {
		return nil
	}
	return validate.ValidateStruct(sr)
}

func cloneSectionRepeat(sr *SectionRepeat) *SectionRepeat {
	if sr == nil {
		return nil
	}

	cpy := *sr
	if sr.MaxInstances != nil {
		cpy.MaxInstances = new(*sr.MaxInstances)
	}
	if sr.AddLabel != nil {
		cpy.AddLabel = new(*sr.AddLabel)
	}

	return &cpy
}
