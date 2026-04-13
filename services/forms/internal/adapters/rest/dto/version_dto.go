package dto

import (
	"time"

	"github.com/cmclaughlin24/sundance/forms/internal/core/domain"
)

type CreateVersionDto struct{}

type UpdateVersionDto struct {
	Pages []PageDto `json:"pages"`
}

type VersionResponseDto struct {
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
	Pages         []*PageResponseDto   `json:"pages"`
}

func VersionToResponseDto(version *domain.Version) *VersionResponseDto {
	if version == nil {
		return nil
	}

	pages := make([]*PageResponseDto, 0, len(version.Pages))
	for _, p := range version.Pages {
		pages = append(pages, PageToResponseDto(p))
	}

	return &VersionResponseDto{
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
