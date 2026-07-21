package dto

import (
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports/commands"
	"time"
)

type upsertElementTagMappingRequest struct {
	TagVersionID string `json:"tagVersionId"`
	Priority     int    `json:"priority"`
}

type ElementTagMappingResponse struct {
	ID           domain.ElementTagMappingID `json:"id"`
	ElementID    domain.ElementID           `json:"elementId"`
	TagVersionID domain.TagVersionID        `json:"tagVersionId"`
	Priority     int                        `json:"priority"`
	CreatedAt    time.Time                  `json:"createdAt"`
	UpdatedAt    time.Time                  `json:"updatedAt"`
}

func requestToElementTagMappingData(dtos []upsertElementTagMappingRequest) []commands.ElementTagMappingData {
	configs := make([]commands.ElementTagMappingData, 0, len(dtos))

	for _, dto := range dtos {
		configs = append(configs, commands.ElementTagMappingData{
			TagVersionID: dto.TagVersionID,
			Priority:     dto.Priority,
		})
	}

	return configs
}

func elementTagMappingsToResponses(tags []*domain.ElementTagMapping) []*ElementTagMappingResponse {
	dtos := make([]*ElementTagMappingResponse, 0, len(tags))

	for _, tag := range tags {
		dtos = append(dtos, &ElementTagMappingResponse{
			ID:           tag.ID,
			ElementID:    tag.ElementID,
			TagVersionID: tag.TagVersionID,
			Priority:     tag.Priority,
			CreatedAt:    tag.CreatedAt,
			UpdatedAt:    tag.UpdatedAt,
		})
	}

	return dtos
}
