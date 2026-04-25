package domain

import (
	"errors"

	"github.com/google/uuid"
)

type PageID string

var (
	ErrInvalidPage = errors.New("invalid page")
)

type Page struct {
	ID       PageID
	Key      string
	Name     string
	sections map[int]*Section
	withPosition
	withRules
}

func NewPage(key, name string, position int) (*Page, error) {
	p := &Page{
		ID:       PageID(uuid.NewString()),
		Key:      key,
		Name:     name,
		sections: make(map[int]*Section),
		withPosition: withPosition{
			position: position,
		},
	}

	// TODO: Implement domain specific validation.

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

func (p *Page) UpdateSections(section ...*Section) error {
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
