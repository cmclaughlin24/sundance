package dto

import "github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"

type PageRequest struct {
	Key      string           `json:"key"`
	Name     string           `json:"name"`
	Position int              `json:"position"`
	Sections []SectionRequest `json:"sections"`
}

type PageResponse struct {
	ID         domain.PageID              `json:"id"`
	Key        string                     `json:"key"`
	Name       string                     `json:"name"`
	Position   int                        `json:"position"`
	Sections   []*SectionResponse         `json:"sections"`
	Conditions []*ConditionalRuleResponse `json:"conditions"`
}

func RequestToPages(dto UpdateVersionDto) ([]*domain.Page, error) {
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

	page, err := domain.NewPage("", dto.Key, dto.Name, dto.Position)
	if err != nil {
		return nil, err
	}

	if err := page.SetSections(sections...); err != nil {
		return nil, err
	}

	return page, nil
}

func PageToResponse(page *domain.Page) *PageResponse {
	if page == nil {
		return nil
	}

	sections := make([]*SectionResponse, 0, len(page.Sections))
	for _, s := range page.Sections {
		sections = append(sections, SectionToResponse(s))
	}

	conditions := ConditionalRulesToResponse(page.Conditions...)

	return &PageResponse{
		ID:         page.ID,
		Key:        page.Key,
		Name:       page.Name,
		Position:   page.Position,
		Sections:   sections,
		Conditions: conditions,
	}
}
