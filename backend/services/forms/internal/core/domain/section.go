package domain

import (
	"errors"
	"maps"
	"slices"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/validate"
	"github.com/google/uuid"
)

type SectionID string

var (
	ErrInvalidSection = errors.New("invalid section")
)

type Section struct {
	ID     SectionID
	Key    string `validate:"required,notblank"`
	Name   string `validate:"required,notblank"`
	fields map[int]*Field
	withPosition
	withRules
}

func NewSection(key, name string, position int) (*Section, error) {
	if !isValidPosition(position) {
		return nil, ErrInvalidPosition
	}

	s := &Section{
		ID:     SectionID(uuid.NewString()),
		Key:    key,
		Name:   name,
		fields: make(map[int]*Field),
		withPosition: withPosition{
			position: position,
		},
	}

	if err := validate.ValidateStruct(s); err != nil {
		return nil, err
	}

	return s, nil
}

func HydrateSection(id SectionID, key, name string, position int) *Section {
	return &Section{
		ID:     id,
		Key:    key,
		Name:   name,
		fields: make(map[int]*Field),
		withPosition: withPosition{
			position: position,
		},
	}
}

func (s *Section) GetFields() map[int]*Field {
	return s.fields
}

func (s *Section) GetFieldsSlice() []*Field {
	positions := slices.Sorted(maps.Keys(s.fields))
	fields := make([]*Field, 0, len(s.fields))

	for _, p := range positions {
		fields = append(fields, s.fields[p])
	}

	return fields
}

func (s *Section) SetFields(fields ...*Field) error {
	if s == nil {
		return ErrInvalidSection
	}

	if s.fields == nil {
		s.fields = make(map[int]*Field)
	}

	for _, field := range fields {
		position := field.GetPosition()
		_, exists := s.fields[position]

		if exists {
			return ErrDuplicatePosition
		}

		s.fields[position] = field
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
