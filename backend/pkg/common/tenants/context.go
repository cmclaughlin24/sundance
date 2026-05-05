package tenants

import (
	"context"
	"errors"
	"log/slog"
	"os"
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
		slog.ErrorContext(ctx, "failed to get tenant from context; tenant not found or of wrong type")
		os.Exit(1)
	}

	return tenantID
}
