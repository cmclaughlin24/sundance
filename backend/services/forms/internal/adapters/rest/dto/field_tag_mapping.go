package dto

import (
	"sundance/backend/services/forms/internal/core/domain"
	"time"
)

type upsertFieldTagMappingRequest struct {
	TagVersionID string `json:"tagVersionId"`
	Priority     int    `json:"priority"`
}

type FieldTagMappingResponse struct {
	ID           domain.FieldTagMappingID `json:"id"`
	FieldID      domain.FieldID           `json:"fieldId"`
	TagVersionID domain.TagVersionID      `json:"tagVersionId"`
	Priority     int                      `json:"priority"`
	CreatedAt    time.Time                `json:"createdAt"`
	UpdatedAt    time.Time                `json:"updatedAt"`
}

func requestToFieldTagMappingConfigs(dtos []upsertFieldTagMappingRequest) []domain.FieldTagMappingConfig {
	configs := make([]domain.FieldTagMappingConfig, 0, len(dtos))

	for _, dto := range dtos {
		configs = append(configs, domain.FieldTagMappingConfig{
			TagVersionID: domain.TagVersionID(dto.TagVersionID),
			Priority:     dto.Priority,
		})
	}

	return configs
}

func fieldTagMappingsToResponses(tags []*domain.FieldTagMapping) []*FieldTagMappingResponse {
	dtos := make([]*FieldTagMappingResponse, 0, len(tags))

	for _, tag := range tags {
		dtos = append(dtos, &FieldTagMappingResponse{
			ID:           tag.ID,
			FieldID:      tag.FieldID,
			TagVersionID: tag.TagVersionID,
			Priority:     tag.Priority,
			CreatedAt:    tag.CreatedAt,
			UpdatedAt:    tag.UpdatedAt,
		})
	}

	return dtos
}
