package dto

import (
	"time"

	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
)

type TenantResponse struct {
	ID          domain.TenantID `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

func TenantToResponse(tenant *domain.Tenant) *TenantResponse {
	return &TenantResponse{
		ID:          tenant.ID,
		Name:        tenant.Name,
		Description: tenant.Description,
		CreatedAt:   tenant.CreatedAt,
		UpdatedAt:   tenant.UpdatedAt,
	}
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
	return &DataSourceResponse{
		ID:          source.ID,
		TenantID:    source.TenantID,
		Name:        source.Name,
		Description: source.Description,
		Type:        source.Type,
		Attributes:  source.Attributes,
		CreatedAt:   source.CreatedAt,
		UpdatedAt:   source.UpdatedAt,
	}
}

type DataSourceLookupResponse struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

func DataSourceLookupToResponse(lookup *domain.DataSourceLookup) *DataSourceLookupResponse {
	return &DataSourceLookupResponse{
		Code:        lookup.Code,
		Description: lookup.Description,
	}
}
