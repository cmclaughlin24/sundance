package ports

import (
	"github.com/cmclaughlin24/sundance/tenants/internal/core/domain"
	"github.com/go-playground/validator/v10"
)

type ListDataSourceQuery struct {
	TenantID domain.TenantID
}

func NewListDataSourceQuery(tenantId domain.TenantID) (*ListDataSourceQuery, error) {
	query := &ListDataSourceQuery{
		TenantID: tenantId,
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(query); err != nil {
		return nil, err
	}

	return query, nil
}
