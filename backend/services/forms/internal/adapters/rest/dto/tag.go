package dto

import (
	"sundance/backend/services/forms/internal/core/domain"
	"time"
)

type CreateTagRequest struct {
	KeyPath       string  `json:"keyPath"`
	DisplayName   string  `json:"displayName"`
	NodeType      string  `json:"nodeType"`
	PrimitiveType *string `json:"primitiveType"`
	IsCollection  bool    `json:"isCollection"`
}

type UpdateTagRequest struct {
	DisplayName string `json:"displayName"`
}

type TagResponse struct {
	ID            domain.TagID             `json:"id"`
	KeyPath       string                   `json:"keyPath"`
	DisplayName   string                   `json:"displayName"`
	NodeType      domain.TagNodeType       `json:"nodeType"`
	PrimitiveType *domain.TagPrimitiveType `json:"primitiveType"`
	CreatedAt     time.Time                `json:"createdAt"`
	UpdatedAt     time.Time                `json:"updatedAt"`
}

func TagToResponse(ct *domain.Tag) TagResponse {
	return TagResponse{
		ID:            ct.ID,
		KeyPath:       ct.KeyPath,
		DisplayName:   ct.DisplayName,
		NodeType:      ct.NodeType,
		PrimitiveType: ct.PrimitiveType,
		CreatedAt:     ct.CreatedAt,
		UpdatedAt:     ct.UpdatedAt,
	}
}
