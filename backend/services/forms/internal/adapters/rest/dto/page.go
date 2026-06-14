package dto

import (
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports/commands"
)

type PageRequest struct {
	ID       *string          `json:"id,omitempty" validate:"omitempty,uuidv7"`
	Key      string           `json:"key" validate:"required,max=50"`
	Name     string           `json:"name" validate:"required,max=75"`
	Position float32          `json:"position" validate:"gte=0,lte=10"`
	Sections []SectionRequest `json:"sections" validate:"dive"`
	Rules    []RuleRequest    `json:"rules" validate:"dive"`
}

type PageResponse struct {
	ID       domain.PageID      `json:"id"`
	Key      string             `json:"key"`
	Name     string             `json:"name"`
	Position float32            `json:"position"`
	Sections []*SectionResponse `json:"sections"`
	Rules    []*RuleResponse    `json:"rules"`
}

func RequestToPages(dto UpsertFormVersionRequest) ([]commands.PageData, error) {
	pages := make([]commands.PageData, 0, len(dto.Pages))

	for _, p := range dto.Pages {
		page, err := RequestToPageData(p)

		if err != nil {
			return nil, err
		}

		pages = append(pages, page)
	}

	return pages, nil
}

func RequestToPageData(dto PageRequest) (commands.PageData, error) {
	sections := make([]commands.SectionData, 0, len(dto.Sections))

	for _, s := range dto.Sections {
		section, err := RequestToSectionData(s)

		if err != nil {
			return commands.PageData{}, err
		}

		sections = append(sections, section)
	}

	rules := RequestsToRuleData(dto.Rules)

	return commands.PageData{
		ID:       dto.ID,
		Key:      dto.Key,
		Name:     dto.Name,
		Position: dto.Position,
		SectionsData: sections,
		Rules:    rules,
	}, nil
}

func PageToResponse(page *domain.Page) *PageResponse {
	if page == nil {
		return nil
	}

	sections := page.GetSections()
	dtos := make([]*SectionResponse, 0, len(sections))

	for _, s := range sections {
		dtos = append(dtos, SectionToResponse(s))
	}

	rules := RulesToResponse(page.GetRules())

	return &PageResponse{
		ID:       page.ID,
		Key:      page.Key,
		Name:     page.Name,
		Position: page.GetPosition(),
		Sections: dtos,
		Rules:    rules,
	}
}
