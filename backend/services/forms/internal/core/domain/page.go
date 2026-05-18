package domain

import (
	"errors"
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
	Name     string `validate:"required"`
	sections PositionElements[*Section]
	withPosition
	withRules
}

func NewPage(key, name string, position float32) (*Page, error) {
	if !isValidPosition(position) {
		return nil, ErrInvalidPosition
	}

	p := &Page{
		ID:       PageID(NewID()),
		Key:      key,
		Name:     name,
		sections: make(PositionElements[*Section], 0),
		withPosition: withPosition{
			position: position,
		},
	}

	if err := validate.ValidateStruct(p); err != nil {
		return nil, err
	}

	return p, nil
}

func HydratePage(id PageID, key, name string, position float32) *Page {
	return &Page{
		ID:       id,
		Key:      key,
		Name:     name,
		sections: make(PositionElements[*Section], 0),
		withPosition: withPosition{
			position: position,
		},
	}
}

func (p *Page) GetSections() []*Section {
	return p.sections
}

func (p *Page) AddSections(sections ...*Section) error {
	if p == nil {
		return ErrInvalidPage
	}

	cpy := slices.Clone(p.sections)
	cpy = append(cpy, sections...)

	if ok := hasUniqueElements(cpy); !ok {
		return ErrDuplicatePosition
	}

	sortElements(cpy)
	p.sections = cpy

	return nil
}

func (p *Page) ReplaceSections(section ...*Section) error {
	if p == nil {
		return ErrInvalidPage
	}

	old := p.sections
	p.sections = make(PositionElements[*Section], 0)

	if err := p.AddSections(section...); err != nil {
		p.sections = old
		return err
	}

	return nil
}
