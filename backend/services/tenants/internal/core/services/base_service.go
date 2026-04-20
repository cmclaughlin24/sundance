package services

import (
	"context"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/tenants"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
)

type baseService struct {
}

func (s baseService) getTenantFromContext(ctx context.Context) (domain.TenantID, error) {
	// FIXME: Change where the decision to throw an error happens.
	tenantID, err := tenants.TenantFromContext(ctx)

	if err != nil {
		return "", err
	}

	return domain.TenantID(tenantID), nil
}
