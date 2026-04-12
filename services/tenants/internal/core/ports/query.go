package ports

import "github.com/cmclaughlin24/sundance/tenants/internal/core/domain"

type ListDataSourceQuery struct {
	TenantID domain.TenantID
}

func NewListDataSourceQuery(tenantId domain.TenantID) ListDataSourceQuery {
	return ListDataSourceQuery{
		TenantID: tenantId,
	}
}
