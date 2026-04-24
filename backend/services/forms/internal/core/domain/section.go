package domain

import (
	"errors"

	"github.com/google/uuid"
)

type SectionID string

var (
	ErrInvalidSection = errors.New("invalid section")
)

type Section struct {
	ID       SectionID
	Key      string
	Name     string
	Position int
	fields   map[int]*Field
	baseWithRules
}

func NewSection(key, name string, position int) (*Section, error) {
	s := &Section{
		ID:       SectionID(uuid.NewString()),
		Key:      key,
		Name:     name,
		Position: position,
		fields:   make(map[int]*Field),
	}

	// TODO: Implement domain specific validation.

	return s, nil
}

func HydrateSection(id SectionID, key, name string, position int) *Section {
	return &Section{
		ID:       id,
		Key:      key,
		Name:     name,
		Position: position,
	}
}

func (s *Section) GetFields() map[int]*Field {
	return s.fields
}

func (s *Section) SetFields(fields ...*Field) error {
	if s == nil {
		return ErrInvalidSection
	}

	if s.fields == nil {
		s.fields = make(map[int]*Field)
	}

	for _, field := range fields {
		_, exists := s.fields[field.Position]

		if exists {
			return ErrDuplicatePosition
		}

		s.fields[field.Position] = field
	}

	return nil
}

func (s *Section) UpdateFields(fields ...*Field) error {
	if s == nil {
		return ErrInvalidSection
	}

	old := s.fields
	s.fields = make(map[int]*Field)

	if err := s.SetFields(fields...); err != nil {
		s.fields = old
		return err
	}

	return nil
}
