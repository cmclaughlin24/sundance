package services

import (
	"context"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/tenants"
)

type baseService struct {
}

func (s baseService) getTenantFromContext(ctx context.Context) (string, error) {
	// FIXME: Change where the decision to throw an error happens.
	tenantID, err := tenants.TenantFromContext(ctx)

	if err != nil {
		return "", err
	}

	return tenantID, nil
}
