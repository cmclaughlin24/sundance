package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/cmclaughlin24/sundance/backend/pkg/common"
)

const tenantIDHeader = "X-Tenant-ID"

type contextKey string

const tenantIDKey contextKey = "tenantID"

var ErrMissingTenantID = errors.New("X-Tenant-ID header is required")

func tenantMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID := r.Header.Get(tenantIDHeader)

		if tenantID == "" {
			common.SendErrorResponse(w, ErrMissingTenantID)
			return
		}

		ctx := withTenantID(r.Context(), tenantID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func withTenantID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, tenantIDKey, tenantID)
}

func tenantIDFromContext(ctx context.Context) (string, error) {
	tenantID, ok := ctx.Value(tenantIDKey).(string)

	if !ok || tenantID == "" {
		return "", ErrMissingTenantID
	}

	return tenantID, nil
}
