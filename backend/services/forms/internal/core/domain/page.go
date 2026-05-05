package domain

import (
	"errors"
	"maps"
	"slices"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/validate"
)

type PageID string

var (
	ErrInvalidPage = errors.New("invalid page")
)

type Page struct {
	ID       PageID
	Key      string `validate:"required,nowhitespace"`
	Name     string `validate:"required,nowhitespace"`
	sections map[int]*Section
	withPosition
	withRules
}

func NewPage(key, name string, position int) (*Page, error) {
	if !isValidPosition(position) {
		return nil, ErrInvalidPosition
	}

	p := &Page{
		ID:       PageID(NewID()),
		Key:      key,
		Name:     name,
		sections: make(map[int]*Section),
		withPosition: withPosition{
			position: position,
		},
	}

	if err := validate.ValidateStruct(p); err != nil {
		return nil, err
	}

	return p, nil
}

func HydratePage(id PageID, key, name string, position int) *Page {
	return &Page{
		ID:       id,
		Key:      key,
		Name:     name,
		sections: make(map[int]*Section),
		withPosition: withPosition{
			position: position,
		},
	}
}

func (p *Page) GetSections() map[int]*Section {
	return p.sections
}

func (p *Page) GetSectionsSlice() []*Section {
	positions := slices.Sorted(maps.Keys(p.sections))
	sections := make([]*Section, 0, len(p.sections))

	for _, position := range positions {
		sections = append(sections, p.sections[position])
	}

	return sections
}

func (p *Page) SetSections(sections ...*Section) error {
	if p == nil {
		return ErrInvalidPage
	}

	if p.sections == nil {
		p.sections = make(map[int]*Section)
	}

	for _, section := range sections {
		position := section.GetPosition()
		_, exists := p.sections[position]

		if exists {
			return ErrDuplicatePosition
		}

		p.sections[position] = section
	}

	return nil
}

func (p *Page) ReplaceSections(section ...*Section) error {
	if p == nil {
		return ErrInvalidPage
	}

	old := p.sections
	p.sections = make(map[int]*Section)

	if err := p.SetSections(section...); err != nil {
		p.sections = old
		return err
	}

	return nil
}
