package httputil

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

const requestDateKey contextKey = "requestDate"

var (
	ErrInvalidRequestDateFormat = errors.New("invalid request date format, must be in RFC3339 format")
)

func SetRequestDateContext(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, correlationIDKey, correlationID)
}

func RequestDateFromContext(ctx context.Context) string {
	requestDate, ok := ctx.Value(requestDateKey).(string)

	if !ok || requestDate == "" {
		slog.DebugContext(ctx, "request date not found in context, returning empty string")
		return ""
	}

	return requestDate
}

func NewRequestDateMiddleware(header string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			date := r.Header.Get(header)
			ctx := r.Context()

			if date != "" {
				if _, err := time.Parse(time.RFC3339, date); err != nil {
					SendErrorResponse(w, ErrInvalidRequestDateFormat)
					return
				}

				ctx = SetRequestDateContext(ctx, date)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
