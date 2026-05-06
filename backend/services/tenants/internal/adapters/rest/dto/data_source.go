package dto

import (
	"time"

	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
)

type DataSourceRequest struct {
	Name        string                `json:"name" validate:"required,max=75"`
	Description string                `json:"description" validate:"max=500"`
	Type        domain.DataSourceType `json:"type" validate:"required"`
	Attributes  any                   `json:"attributes" validate:"required"`
}

type DataSourceResponse struct {
	ID          domain.DataSourceID   `json:"id"`
	TenantID    domain.TenantID       `json:"tenantId"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Type        domain.DataSourceType `json:"type"`
	Attributes  any                   `json:"attributes"`
	CreatedAt   time.Time             `json:"createdAt"`
	UpdatedAt   time.Time             `json:"updatedAt"`
}

func DataSourceToResponse(source *domain.DataSource) *DataSourceResponse {
	attr := dataSourceAttributesToResponse(source.Attributes)

	return &DataSourceResponse{
		ID:          source.ID,
		TenantID:    source.TenantID,
		Name:        source.Name,
		Description: source.Description,
		Type:        source.Type,
		Attributes:  attr,
		CreatedAt:   source.CreatedAt,
		UpdatedAt:   source.UpdatedAt,
	}
}
