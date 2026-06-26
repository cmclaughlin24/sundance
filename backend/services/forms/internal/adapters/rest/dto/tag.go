package dto

import (
	"sundance/backend/services/forms/internal/core/domain"
	"time"
)

type CreateTagRequest struct {
	Key           string  `json:"key"`
	DisplayName   string  `json:"displayName"`
	NodeType      string  `json:"nodeType"`
	PrimitiveType *string `json:"primitiveType"`
	IsCollection  bool    `json:"isCollection"`
}

type UpdateTagRequest struct {
	DisplayName string `json:"displayName"`
}

type TagResponse struct {
	ID            domain.TagID            `json:"id"`
	Key           string                  `json:"key"`
	DisplayName   string                  `json:"displayName"`
	NodeType      domain.TagNodeType      `json:"nodeType"`
	PrimitiveType *domain.TagPrimitiveType `json:"primitiveType"`
	IsCollection  bool                    `json:"isCollection"`
	CreatedAt     time.Time               `json:"createdAt"`
	UpdatedAt     time.Time               `json:"updatedAt"`
}

func TagToResponse(ct *domain.Tag) TagResponse {
	return TagResponse{
		ID:            ct.ID,
		Key:           ct.Key,
		DisplayName:   ct.DisplayName,
		NodeType:      ct.NodeType,
		PrimitiveType: ct.PrimitiveType,
		IsCollection:  ct.IsCollection,
		CreatedAt:     ct.CreatedAt,
		UpdatedAt:     ct.UpdatedAt,
	}
}
