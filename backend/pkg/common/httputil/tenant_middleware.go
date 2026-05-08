package httputil

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
)

const tenantIDKey contextKey = "tenantID"

var ErrMissingTenantID = errors.New("tenant header is required")

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

func NewTenantMiddleware(header string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get(header)

			if id == "" {
				SendErrorResponse(w, ErrMissingTenantID)
				return
			}

			ctx := SetTenantContext(r.Context(), id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
