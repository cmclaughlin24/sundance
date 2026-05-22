package ports

import (
	"sundance/backend/pkg/common/validate"
	"sundance/backend/services/tenants/internal/core/domain"
)

type ListDataSourceQuery struct {
	// TODO: Add pagination support through embedded struct.
	TenantID domain.TenantID `validate:"required"`
}

func NewListDataSourceQuery(tenantID domain.TenantID) *ListDataSourceQuery {
	return &ListDataSourceQuery{
		TenantID: tenantID,
	}
}

func (q *ListDataSourceQuery) Validate() error {
	return validate.ValidateStruct(q)
}

type FindDataSourceByIDQuery struct {
	TenantID domain.TenantID     `validate:"required"`
	ID       domain.DataSourceID `validate:"required"`
}

func NewFindDataSourceByIDQuery(tenantID domain.TenantID, sourceID domain.DataSourceID) *FindDataSourceByIDQuery {
	return &FindDataSourceByIDQuery{
		TenantID: tenantID,
		ID:       sourceID,
	}
}

func (q *FindDataSourceByIDQuery) Validate() error {
	return validate.ValidateStruct(q)
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

func (q *GetDataSourceLookupsQuery) Validate() error {
	return validate.ValidateStruct(q)
}

type FindDataSourceJobsQuery struct {
	Take       int `validate:"min=0"`
	RetryLimit int `validate:"min=0"`
}

func NewFindDataSourceJobsQuery(take int, retryLimit int) *FindDataSourceJobsQuery {
	return &FindDataSourceJobsQuery{Take: take, RetryLimit: retryLimit}
}

func (q *FindDataSourceJobsQuery) Validate() error {
	return validate.ValidateStruct(q)
}
