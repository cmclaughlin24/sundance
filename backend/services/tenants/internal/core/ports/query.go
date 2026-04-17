package ports

import (
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/validate"
)

type ListDataSourceQuery struct {
	TenantID domain.TenantID `validate:"required"`
}

func NewListDataSourceQuery(tenantId domain.TenantID) (*ListDataSourceQuery, error) {
	query := &ListDataSourceQuery{
		TenantID: tenantId,
	}

	if err := validate.ValidateStruct(query); err != nil {
		return nil, err
	}

	return query, nil
}
