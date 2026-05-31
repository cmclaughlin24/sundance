package dto

import (
	"sundance/backend/services/forms/internal/core/domain"
	"time"
)

type UpsertTagVersionRequest struct {
	Type string `json:"required"`
}

type TagVersionResponse struct {
	ID        domain.TagVersionID `json:"id"`
	TagID     domain.TagID        `json:"tagId"`
	Version   int                 `json:"version"`
	Type      domain.TagType      `json:"type"`
	Status    domain.TagStatus    `json:"status"`
	CreatedAt time.Time           `json:"createdAt"`
	RetiredAt time.Time           `json:"retiredAt"`
}

func TagVersionToResponse(t *domain.TagVersion) TagVersionResponse {
	return TagVersionResponse{
		ID:        t.ID,
		TagID:     t.TagID,
		Version:   t.Version,
		Type:      t.Type,
		Status:    t.Status,
		CreatedAt: t.CreatedAt,
		RetiredAt: t.RetiredAt,
	}
}
