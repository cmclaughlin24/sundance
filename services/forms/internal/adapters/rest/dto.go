package rest

import (
	"time"

	"github.com/cmclaughlin24/sundance/forms/internal/core/domain"
)

type upsertFormDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type versionDto struct {
	ID            domain.VersionID     `json:"id"`
	FormID        domain.FormID        `json:"formId"`
	Version       int                  `json:"version"`
	Status        domain.VersionStatus `json:"status"`
	PublishedByID string               `json:"publishedById"`
	PublishedAt   time.Time            `json:"publishedAt"`
	RetiredByID   string               `json:"retiredById"`
	RetiredAt     time.Time            `json:"retiredAt"`
	CreatedAt     time.Time            `json:"createdAt"`
	UpdatedAt     time.Time            `json:"updatedAt"`
	Pages         []*pageDto           `json:"pages"`
}

func versionToDto(version *domain.Version) *versionDto {
	if version == nil {
		return nil
	}

	pages := make([]*pageDto, 0, len(version.Pages))
	for _, p := range version.Pages {
		pages = append(pages, pageToDto(p))
	}

	return &versionDto{
		ID:            version.ID,
		FormID:        version.FormID,
		Version:       version.Version,
		Status:        version.Status,
		PublishedByID: version.PublishedByID,
		PublishedAt:   version.PublishedAt,
		RetiredByID:   version.RetiredByID,
		RetiredAt:     version.RetiredAt,
		CreatedAt:     version.CreatedAt,
		UpdatedAt:     version.UpdatedAt,
		Pages:         pages,
	}
}

type pageDto struct {
	ID         domain.PageID         `json:"id"`
	Key        string                `json:"key"`
	Name       string                `json:"name"`
	Position   int                   `json:"position"`
	Sections   []*sectionDto         `json:"sections"`
	Conditions []*conditionalRuleDto `json:"conditions"`
}

func pageToDto(page *domain.Page) *pageDto {
	if page == nil {
		return nil
	}

	sections := make([]*sectionDto, 0, len(page.Sections))
	for _, s := range page.Sections {
		sections = append(sections, sectionToDto(s))
	}

	conditions := conditionalRulesToDtos(page.Conditions...)

	return &pageDto{
		ID:         page.ID,
		Key:        page.Key,
		Name:       page.Name,
		Position:   page.Position,
		Sections:   sections,
		Conditions: conditions,
	}
}

type sectionDto struct {
	ID         domain.SectionID      `json:"id"`
	Key        string                `json:"key"`
	Name       string                `json:"name"`
	Position   int                   `json:"position"`
	Fields     []*fieldDto           `json:"fields"`
	Conditions []*conditionalRuleDto `json:"conditions"`
}

func sectionToDto(section *domain.Section) *sectionDto {
	if section == nil {
		return nil
	}

	fields := make([]*fieldDto, 0, len(section.Fields))
	for _, f := range section.Fields {
		fields = append(fields, fieldToDto(f))
	}

	conditions := conditionalRulesToDtos(section.Conditions...)

	return &sectionDto{
		ID:         section.ID,
		Key:        section.Key,
		Name:       section.Name,
		Position:   section.Position,
		Fields:     fields,
		Conditions: conditions,
	}
}

type fieldDto struct {
	ID         domain.FieldID        `json:"id"`
	Key        string                `json:"key"`
	Name       string                `json:"name"`
	FieldType  string                `json:"fieldType"`
	Position   int                   `json:"position"`
	Conditions []*conditionalRuleDto `json:"conditions"`
}

func fieldToDto(field *domain.Field) *fieldDto {
	if field == nil {
		return nil
	}

	conditions := conditionalRulesToDtos(field.Conditions...)

	return &fieldDto{
		ID:         field.ID,
		Key:        field.Key,
		Name:       field.Name,
		FieldType:  field.FieldType,
		Position:   field.Position,
		Conditions: conditions,
	}
}

type conditionalRuleDto struct {
	ID domain.ConditionalRuleID `json:"id"`
}

func conditionalRulesToDtos(rules ...*domain.ConditionalRule) []*conditionalRuleDto {
	conditions := make([]*conditionalRuleDto, 0, len(rules))
	for _, c := range rules {
		conditions = append(conditions, &conditionalRuleDto{
			ID: c.ID,
		})
	}

	return conditions
}
