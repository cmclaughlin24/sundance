package dto

import (
	"time"

	"sundance/backend/services/forms/internal/core/domain"
)

type UpsertFormVersionRequest struct {
	Pages []PageRequest `json:"pages" validate:"dive"`
}

type FormVersionResponse struct {
	ID            domain.FormVersionID     `json:"id"`
	FormID        domain.FormID        `json:"formId"`
	Version       int                  `json:"version"`
	Status        domain.FormVersionStatus `json:"status"`
	PublishedByID string               `json:"publishedById"`
	PublishedAt   time.Time            `json:"publishedAt"`
	RetiredByID   string               `json:"retiredById"`
	RetiredAt     time.Time            `json:"retiredAt"`
	CreatedAt     time.Time            `json:"createdAt"`
	UpdatedAt     time.Time            `json:"updatedAt"`
	Pages         []*PageResponse      `json:"pages"`
}

func FormVersionToResponse(version *domain.FormVersion) *FormVersionResponse {
	if version == nil {
		return nil
	}

	pages := version.GetPages()
	dtos := make([]*PageResponse, 0, len(pages))

	for _, p := range pages {
		dtos = append(dtos, PageToResponse(p))
	}

	return &FormVersionResponse{
		ID:            version.ID,
		FormID:        version.FormID,
		Version:       version.Version,
		Status:        version.Status,
		PublishedByID: version.PublishedBy,
		PublishedAt:   version.PublishedAt,
		RetiredByID:   version.RetiredBy,
		RetiredAt:     version.RetiredAt,
		CreatedAt:     version.CreatedAt,
		UpdatedAt:     version.UpdatedAt,
		Pages:         dtos,
	}
}
