package domain

import (
	"errors"
)

type SectionID string

var (
	ErrInvalidSection = errors.New("invalid section")
)

type Section struct {
	ID         SectionID
	Key        string
	Name       string
	Position   int
	Fields     map[int]*Field
	Conditions []*ConditionalRule
}

func NewSection(id SectionID, key, name string, position int) (*Section, error) {
	s := &Section{
		ID:       id,
		Key:      key,
		Name:     name,
		Position: position,
		Fields:   make(map[int]*Field),
	}

	// TODO: Implement domain specific validation.

	return s, nil
}

func (s *Section) SetFields(fields ...*Field) error {
	if s == nil {
		return ErrInvalidSection
	}

	if s.Fields == nil {
		s.Fields = make(map[int]*Field)
	}

	for _, field := range fields {
		_, exists := s.Fields[field.Position]

		if exists {
			return ErrDuplicatePosition
		}

		s.Fields[field.Position] = field
	}

	return nil
}

func (s *Section) UpdateFields(fields ...*Field) error {
	if s == nil {
		return ErrInvalidSection
	}

	s.Fields = make(map[int]*Field)

	return s.SetFields(fields...)
}
