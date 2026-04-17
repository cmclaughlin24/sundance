package ports

import (
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
)

type ListDataSourceQuery struct {
	TenantID domain.TenantID `validate:"required"`
}

func NewListDataSourceQuery(tenantId domain.TenantID) *ListDataSourceQuery {
	return &ListDataSourceQuery{
		TenantID: tenantId,
	}
}
