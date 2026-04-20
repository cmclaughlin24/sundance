package tenants

import (
	"context"
	"errors"
)

type contextKey string

const tenantIDKey contextKey = "tenantID"

var ErrMissingTenantID = errors.New("X-Tenant-ID header is required")

func WithTenant(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, tenantIDKey, tenantID)
}

func TenantFromContext(ctx context.Context) (string, error) {
	tenantID, ok := ctx.Value(tenantIDKey).(string)

	if !ok || tenantID == "" {
		return "", ErrMissingTenantID
	}

	return tenantID, nil
}
