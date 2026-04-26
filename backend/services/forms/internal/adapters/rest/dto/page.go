package dto

import (
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
)

type PageRequest struct {
	Key      string           `json:"key"`
	Name     string           `json:"name"`
	Position int              `json:"position"`
	Sections []SectionRequest `json:"sections"`
	Rules    []RuleRequest    `json:"rules"`
}

type PageResponse struct {
	ID       domain.PageID      `json:"id"`
	Key      string             `json:"key"`
	Name     string             `json:"name"`
	Position int                `json:"position"`
	Sections []*SectionResponse `json:"sections"`
	Rules    []*RuleResponse    `json:"rules"`
}

func RequestToPages(dto UpdateVersionRequest) ([]*domain.Page, error) {
	pages := make([]*domain.Page, 0, len(dto.Pages))

	for _, p := range dto.Pages {
		page, err := RequestToPage(p)

		if err != nil {
			return nil, err
		}

		pages = append(pages, page)
	}

	return pages, nil
}

func RequestToPage(dto PageRequest) (*domain.Page, error) {
	sections := make([]*domain.Section, 0, len(dto.Sections))

	for _, s := range dto.Sections {
		section, err := RequestToSection(s)

		if err != nil {
			return nil, err
		}

		sections = append(sections, section)
	}

	rules, err := RequestsToRules(dto.Rules)
	if err != nil {
		return nil, err
	}

	page, err := domain.NewPage(dto.Key, dto.Name, dto.Position)
	if err != nil {
		return nil, err
	}

	if err := page.SetSections(sections...); err != nil {
		return nil, err
	}

	if err := page.SetRules(rules...); err != nil {
		return nil, err
	}

	return page, nil
}

func PageToResponse(page *domain.Page) *PageResponse {
	if page == nil {
		return nil
	}

	sections := page.GetSectionsSlice()
	dtos := make([]*SectionResponse, 0, len(sections))

	for _, s := range sections {
		dtos = append(dtos, SectionToResponse(s))
	}

	rules := RuleToResponse(page.GetRules())

	return &PageResponse{
		ID:       page.ID,
		Key:      page.Key,
		Name:     page.Name,
		Position: page.GetPosition(),
		Sections: dtos,
		Rules:    rules,
	}
}
