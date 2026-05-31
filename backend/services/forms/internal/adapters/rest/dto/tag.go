package dto

import (
	"sundance/backend/services/forms/internal/core/domain"
	"time"
)

type CreateTagRequest struct {
	Key         string `json:"key"`
	DisplayName string `json:"displayName"`
}

type UpdateTagRequest struct {
	DisplayName string `json:"displayName"`
}

type TagResponse struct {
	ID          domain.TagID `json:"id"`
	Key         string       `json:"key"`
	DisplayName string       `json:"displayName"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}

func TagToResponse(ct *domain.Tag) *TagResponse {
	return &TagResponse{
		ID:          ct.ID,
		Key:         ct.Key,
		DisplayName: ct.DisplayName,
		CreatedAt:   ct.CreatedAt,
		UpdatedAt:   ct.UpdatedAt,
	}
}
