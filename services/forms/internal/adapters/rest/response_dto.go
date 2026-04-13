package rest

import (
	"time"

	"github.com/cmclaughlin24/sundance/forms/internal/core/domain"
)

type formResponseDto struct {
	ID          domain.FormID `json:"id"`
	TenantID    string        `json:"tenantId"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
}

func formToResponseDto(form *domain.Form) *formResponseDto {
	if form == nil {
		return nil
	}

	return &formResponseDto{
		ID:          form.ID,
		TenantID:    form.TenantID,
		Name:        form.Name,
		Description: form.Description,
		CreatedAt:   form.CreatedAt,
		UpdatedAt:   form.UpdatedAt,
	}
}

type versionResponseDto struct {
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
	Pages         []*pageResponseDto   `json:"pages"`
}

func versionToResponseDto(version *domain.Version) *versionResponseDto {
	if version == nil {
		return nil
	}

	pages := make([]*pageResponseDto, 0, len(version.Pages))
	for _, p := range version.Pages {
		pages = append(pages, pageToResponseDto(p))
	}

	return &versionResponseDto{
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

type pageResponseDto struct {
	ID         domain.PageID                 `json:"id"`
	Key        string                        `json:"key"`
	Name       string                        `json:"name"`
	Position   int                           `json:"position"`
	Sections   []*sectionResponseDto         `json:"sections"`
	Conditions []*conditionalRuleResponseDto `json:"conditions"`
}

func pageToResponseDto(page *domain.Page) *pageResponseDto {
	if page == nil {
		return nil
	}

	sections := make([]*sectionResponseDto, 0, len(page.Sections))
	for _, s := range page.Sections {
		sections = append(sections, sectionToResponseDto(s))
	}

	conditions := conditionalRulesToResponseDtos(page.Conditions...)

	return &pageResponseDto{
		ID:         page.ID,
		Key:        page.Key,
		Name:       page.Name,
		Position:   page.Position,
		Sections:   sections,
		Conditions: conditions,
	}
}

type sectionResponseDto struct {
	ID         domain.SectionID              `json:"id"`
	Key        string                        `json:"key"`
	Name       string                        `json:"name"`
	Position   int                           `json:"position"`
	Fields     []*fieldResponseDto           `json:"fields"`
	Conditions []*conditionalRuleResponseDto `json:"conditions"`
}

func sectionToResponseDto(section *domain.Section) *sectionResponseDto {
	if section == nil {
		return nil
	}

	fields := make([]*fieldResponseDto, 0, len(section.Fields))
	for _, f := range section.Fields {
		fields = append(fields, fieldToResponseDto(f))
	}

	conditions := conditionalRulesToResponseDtos(section.Conditions...)

	return &sectionResponseDto{
		ID:         section.ID,
		Key:        section.Key,
		Name:       section.Name,
		Position:   section.Position,
		Fields:     fields,
		Conditions: conditions,
	}
}

type fieldResponseDto struct {
	ID         domain.FieldID                `json:"id"`
	Key        string                        `json:"key"`
	Name       string                        `json:"name"`
	FieldType  string                        `json:"fieldType"`
	Position   int                           `json:"position"`
	Conditions []*conditionalRuleResponseDto `json:"conditions"`
}

func fieldToResponseDto(field *domain.Field) *fieldResponseDto {
	if field == nil {
		return nil
	}

	conditions := conditionalRulesToResponseDtos(field.Conditions...)

	return &fieldResponseDto{
		ID:         field.ID,
		Key:        field.Key,
		Name:       field.Name,
		FieldType:  field.FieldType,
		Position:   field.Position,
		Conditions: conditions,
	}
}

type conditionalRuleResponseDto struct {
	ID domain.ConditionalRuleID `json:"id"`
}

func conditionalRulesToResponseDtos(rules ...*domain.ConditionalRule) []*conditionalRuleResponseDto {
	conditions := make([]*conditionalRuleResponseDto, 0, len(rules))
	for _, c := range rules {
		conditions = append(conditions, &conditionalRuleResponseDto{
			ID: c.ID,
		})
	}

	return conditions
}
