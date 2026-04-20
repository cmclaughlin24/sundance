package httputil

import (
	"context"
	"errors"
	"net/http"
)

type contextKey string

const tenantIDKey contextKey = "tenantID"

var ErrMissingTenantID = errors.New("X-Tenant-ID header is required")

func TenantMiddleware(header string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get(header)

			if id == "" {
				SendErrorResponse(w, ErrMissingTenantID)
				return
			}

			ctx := withTenant(r.Context(), id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func withTenant(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, tenantIDKey, tenantID)
}

func TenantFromContext(ctx context.Context) (string, error) {
	tenantID, ok := ctx.Value(tenantIDKey).(string)

	if !ok || tenantID == "" {
		return "", ErrMissingTenantID
	}

	return tenantID, nil
}
