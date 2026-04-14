package rest

import (
	"time"

	"github.com/cmclaughlin24/sundance/tenants/internal/core/domain"
)

type tenantResponseDto struct {
	ID          domain.TenantID `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

func tenantToResponseDto(tenant *domain.Tenant) *tenantResponseDto {
	return &tenantResponseDto{
		ID:          tenant.ID,
		Name:        tenant.Name,
		Description: tenant.Description,
		CreatedAt:   tenant.CreatedAt,
		UpdatedAt:   tenant.UpdatedAt,
	}
}

type dataSourceResponseDto struct {
	ID         domain.DataSourceID   `json:"id"`
	TenantID   domain.TenantID       `json:"tenantId"`
	Type       domain.DataSourceType `json:"type"`
	Attributes any                   `json:"attributes"`
	CreatedAt  time.Time             `json:"createdAt"`
	UpdatedAt  time.Time             `json:"updatedAt"`
}

func dataSourceToResponseDto(source *domain.DataSource) *dataSourceResponseDto {
	return &dataSourceResponseDto{
		ID:         source.ID,
		TenantID:   source.TenantID,
		Type:       source.Type,
		Attributes: source.Attributes,
		CreatedAt:  source.CreatedAt,
		UpdatedAt:  source.UpdatedAt,
	}
}

type dataSourceLookupResponseDto struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

func dataSourceLookupToResponseDto(lookup *domain.DataSourceLookup) *dataSourceLookupResponseDto {
	return &dataSourceLookupResponseDto{
		Code:        lookup.Code,
		Description: lookup.Description,
	}
}
