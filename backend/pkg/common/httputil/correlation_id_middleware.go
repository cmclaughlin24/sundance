package httputil

import (
	"context"
	"log/slog"
	"net/http"
)

const correlationIDKey contextKey = "correlationID"

func SetCorrelationIDContext(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, correlationIDKey, correlationID)
}

func CorrelationIDFromContext(ctx context.Context) string {
	correlationID, ok := ctx.Value(correlationIDKey).(string)

	if !ok || correlationID == "" {
		slog.DebugContext(ctx, "correlation ID not found in context, returning empty string")
		return ""
	}

	return correlationID
}

func NewCorrelationIDMiddleware(header string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get(header)
			ctx := r.Context()

			if id != "" {
				ctx = SetCorrelationIDContext(ctx, id)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
