package rest

import (
	"time"

	"github.com/cmclaughlin24/sundance/tenants/internal/core/domain"
)

type upsertTenantDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type tenantDto struct {
	ID          domain.TenantID `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

func tenantToDto(tenant *domain.Tenant) *tenantDto {
	return &tenantDto{
		ID:          tenant.ID,
		Name:        tenant.Name,
		Description: tenant.Description,
		CreatedAt:   tenant.CreatedAt,
		UpdatedAt:   tenant.UpdatedAt,
	}
}

type upsertDataSourceDto struct {
	Type       domain.DataSourceType `json:"type"`
	Attributes any                   `json:"attributes"`
}

type dataSourceDto struct {
	ID         domain.DataSourceID   `json:"id"`
	TenantID   domain.TenantID       `json:"tenantId"`
	Type       domain.DataSourceType `json:"type"`
	Attributes any                   `json:"attributes"`
	CreatedAt  time.Time             `json:"createdAt"`
	UpdatedAt  time.Time             `json:"updatedAt"`
}

func dataSourceToDto(source *domain.DataSource) *dataSourceDto {
	return &dataSourceDto{
		ID:         source.ID,
		TenantID:   source.TenantID,
		Type:       source.Type,
		Attributes: source.Attributes,
		CreatedAt:  source.CreatedAt,
		UpdatedAt:  source.UpdatedAt,
	}
}

type dataSourceLookupDto struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

func dataSourceLookupToDto(lookup *domain.DataSourceLookup) *dataSourceLookupDto {
	return &dataSourceLookupDto{
		Code:        lookup.Code,
		Description: lookup.Description,
	}
}
