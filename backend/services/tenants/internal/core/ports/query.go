package ports

import "github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"

type ListDataSourceQuery struct {
	// TODO: Add pagination support through embedded struct.
	TenantID domain.TenantID `validate:"required"`
}

func NewListDataSourceQuery(tenantID domain.TenantID) *ListDataSourceQuery {
	return &ListDataSourceQuery{
		TenantID: tenantID,
	}
}

type FindDataSourceByIDQuery struct {
	TenantID domain.TenantID     `validate:"required"`
	ID       domain.DataSourceID `validate:"required"`
}

func NewFindDataSourceByID(tenantID domain.TenantID, sourceID domain.DataSourceID) *FindDataSourceByIDQuery {
	return &FindDataSourceByIDQuery{
		TenantID: tenantID,
		ID:       sourceID,
	}
}

type GetDataSourceLookupsQuery struct {
	TenantID domain.TenantID     `validate:"required"`
	ID       domain.DataSourceID `validate:"required"`
}

func NewGetDataSourceLookupsQuery(tenantID domain.TenantID, sourceID domain.DataSourceID) *GetDataSourceLookupsQuery {
	return &GetDataSourceLookupsQuery{
		TenantID: tenantID,
		ID:       sourceID,
	}
}
