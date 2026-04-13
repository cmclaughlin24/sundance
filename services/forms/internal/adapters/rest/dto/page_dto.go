package dto

import "github.com/cmclaughlin24/sundance/forms/internal/core/domain"

type PageDto struct {
	Key      string       `json:"key"`
	Name     string       `json:"name"`
	Position int          `json:"position"`
	Sections []SectionDto `json:"sections"`
}

type PageResponseDto struct {
	ID         domain.PageID                 `json:"id"`
	Key        string                        `json:"key"`
	Name       string                        `json:"name"`
	Position   int                           `json:"position"`
	Sections   []*SectionResponseDto         `json:"sections"`
	Conditions []*ConditionalRuleResponseDto `json:"conditions"`
}

func DtoToPages(dto UpdateVersionDto) ([]*domain.Page, error) {
	pages := make([]*domain.Page, 0, len(dto.Pages))

	for _, p := range dto.Pages {
		page, err := DtoToPage(p)

		if err != nil {
			return nil, err
		}

		pages = append(pages, page)
	}

	return pages, nil
}

func DtoToPage(dto PageDto) (*domain.Page, error) {
	sections := make([]*domain.Section, 0, len(dto.Sections))

	for _, s := range dto.Sections {
		section, err := DtoToSection(s)

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

func PageToResponseDto(page *domain.Page) *PageResponseDto {
	if page == nil {
		return nil
	}

	sections := make([]*SectionResponseDto, 0, len(page.Sections))
	for _, s := range page.Sections {
		sections = append(sections, SectionToResponseDto(s))
	}

	conditions := ConditionalRulesToResponseDtos(page.Conditions...)

	return &PageResponseDto{
		ID:         page.ID,
		Key:        page.Key,
		Name:       page.Name,
		Position:   page.Position,
		Sections:   sections,
		Conditions: conditions,
	}
}
