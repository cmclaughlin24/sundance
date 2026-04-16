package rest

import (
	"time"

	"github.com/cmclaughlin24/sundance/tenants/internal/core/domain"
)

type tenantResponse struct {
	ID          domain.TenantID `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

func tenantToResponse(tenant *domain.Tenant) *tenantResponse {
	return &tenantResponse{
		ID:          tenant.ID,
		Name:        tenant.Name,
		Description: tenant.Description,
		CreatedAt:   tenant.CreatedAt,
		UpdatedAt:   tenant.UpdatedAt,
	}
}

type dataSourceResponse struct {
	ID         domain.DataSourceID   `json:"id"`
	TenantID   domain.TenantID       `json:"tenantId"`
	Type       domain.DataSourceType `json:"type"`
	Attributes any                   `json:"attributes"`
	CreatedAt  time.Time             `json:"createdAt"`
	UpdatedAt  time.Time             `json:"updatedAt"`
}

func dataSourceToResponseDto(source *domain.DataSource) *dataSourceResponse {
	return &dataSourceResponse{
		ID:         source.ID,
		TenantID:   source.TenantID,
		Type:       source.Type,
		Attributes: source.Attributes,
		CreatedAt:  source.CreatedAt,
		UpdatedAt:  source.UpdatedAt,
	}
}

type dataSourceLookupResponse struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

func dataSourceLookupToResponse(lookup *domain.DataSourceLookup) *dataSourceLookupResponse {
	return &dataSourceLookupResponse{
		Code:        lookup.Code,
		Description: lookup.Description,
	}
}
