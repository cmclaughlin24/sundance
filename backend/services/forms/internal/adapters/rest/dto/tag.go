package dto

import (
	"sundance/backend/services/forms/internal/core/domain"
	"time"
)

type CreateTagRequest struct {
	Key          string `json:"key"`
	DisplayName  string `json:"displayName"`
	ValueKind    string `json:"valueKind"`
	IsCollection bool   `json:"isCollection"`
}

type UpdateTagRequest struct {
	DisplayName string `json:"displayName"`
}

type TagResponse struct {
	ID           domain.TagID        `json:"id"`
	Key          string              `json:"key"`
	DisplayName  string              `json:"displayName"`
	ValueKind    domain.TagValueKind `json:"valueKind"`
	IsCollection bool                `json:"isCollection"`
	CreatedAt    time.Time           `json:"createdAt"`
	UpdatedAt    time.Time           `json:"updatedAt"`
}

func TagToResponse(ct *domain.Tag) TagResponse {
	return TagResponse{
		ID:           ct.ID,
		Key:          ct.Key,
		DisplayName:  ct.DisplayName,
		ValueKind:    ct.ValueKind,
		IsCollection: ct.IsCollection,
		CreatedAt:    ct.CreatedAt,
		UpdatedAt:    ct.UpdatedAt,
	}
}
