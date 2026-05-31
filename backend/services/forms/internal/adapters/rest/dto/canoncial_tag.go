package dto

import (
	"sundance/backend/services/forms/internal/core/domain"
	"time"
)

type CreateCanonicalTagRequest struct {
	Key         string `json:"key"`
	DisplayName string `json:"displayName"`
}

type UpdateCanonicalTagRequest struct {
	DisplayName string `json:"displayName"`
}

type CanonicalTagResponse struct {
	ID          domain.CanonicalTagID `json:"id"`
	Key         string                `json:"key"`
	DisplayName string                `json:"displayName"`
	CreatedAt   time.Time             `json:"createdAt"`
	UpdatedAt   time.Time             `json:"updatedAt"`
}

func CanonicalTagToResponse(ct *domain.CanonicalTag) *CanonicalTagResponse {
	return &CanonicalTagResponse{
		ID:          ct.ID,
		Key:         ct.Key,
		DisplayName: ct.DisplayName,
		CreatedAt:   ct.CreatedAt,
		UpdatedAt:   ct.UpdatedAt,
	}
}
