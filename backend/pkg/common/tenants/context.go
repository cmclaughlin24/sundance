package tenants

import (
	"context"
	"errors"
	"log/slog"
)

type contextKey string

const tenantIDKey contextKey = "tenantID"

var ErrMissingTenantID = errors.New("X-Tenant-ID header is required")

func SetTenantContext(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, tenantIDKey, tenantID)
}

func TenantFromContext(ctx context.Context) string {
	tenantID, ok := ctx.Value(tenantIDKey).(string)

	if !ok || tenantID == "" {
		err := errors.New("failed to get tenant from context; tenant not found or of wrong type")
		slog.ErrorContext(ctx, "error", err)
		panic(err)
	}

	return tenantID
}
