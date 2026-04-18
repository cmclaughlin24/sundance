package domain

import (
	"errors"

	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/types"
)

type PageID string

var (
	ErrInvalidPage = errors.New("invalid page")
)

type Page struct {
	ID         PageID
	Key        string
	Name       string
	Position   int
	Sections   map[int]*Section
	Conditions []*ConditionalRule
}

func NewPage(id PageID, key, name string, position int) (*Page, error) {
	p := &Page{
		ID:       id,
		Key:      key,
		Name:     name,
		Position: position,
		Sections: make(map[int]*Section),
	}

	// TODO: Implement domain specific validation.

	return p, nil
}

func (p *Page) SetSections(sections ...*Section) error {
	if p == nil {
		return ErrInvalidPage
	}

	if p.Sections == nil {
		p.Sections = make(map[int]*Section)
	}

	for _, section := range sections {
		_, exists := p.Sections[section.Position]

		if exists {
			return types.ErrDuplicatePosition
		}

		p.Sections[section.Position] = section
	}

	return nil
}

func (p *Page) UpdateSections(section ...*Section) error {
	if p == nil {
		return ErrInvalidPage
	}

	p.Sections = make(map[int]*Section)

	return p.SetSections(section...)
}
