package dto

import (
	"sundance/backend/services/forms/internal/core/domain"
	"time"
)

type UpsertTagVersionRequest struct {
	Type domain.TagType `json:"type" validate:"required"`
}

type TagVersionResponse struct {
	ID           domain.TagVersionID `json:"id"`
	TagID        domain.TagID        `json:"tagId"`
	Version      int                 `json:"version"`
	Type         domain.TagType      `json:"type"`
	Status       domain.TagStatus    `json:"status"`
	CreatedAt    time.Time           `json:"createdAt"`
	DeprecatedAt time.Time           `json:"deprecatedAt"`
	PublishedAt  time.Time           `json:"publishedAt"`
	RetiredAt    time.Time           `json:"retiredAt"`
}

func TagVersionToResponse(tv *domain.TagVersion) TagVersionResponse {
	return TagVersionResponse{
		ID:           tv.ID,
		TagID:        tv.TagID,
		Version:      tv.Version,
		Type:         tv.Type,
		Status:       tv.Status,
		CreatedAt:    tv.CreatedAt,
		DeprecatedAt: tv.DeprecatedAt,
		PublishedAt:  tv.PublishedAt,
		RetiredAt:    tv.RetiredAt,
	}
}
