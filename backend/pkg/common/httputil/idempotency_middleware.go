package httputil

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
)

var (
	ErrMissingIdempotencyHeader = errors.New("idempotency header is required")
	ErrIdempotencyContext       = errors.New("idempotency not found or wrong type")
)

const idempotencyIDKey contextKey = "idempotencyID"

func SetIdempotencyContext(ctx context.Context, idempotencyID string) context.Context {
	return context.WithValue(ctx, idempotencyIDKey, idempotencyID)
}

func IdempotencyFromContext(ctx context.Context) string {
	idempotencyID, ok := ctx.Value(idempotencyIDKey).(string)

	if !ok || idempotencyID == "" {
		slog.ErrorContext(ctx, "failed to get idempotency from context", "error", ErrIdempotencyContext)
		panic(ErrIdempotencyContext)
	}

	return idempotencyID
}

func IdempotencyMiddlware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("Idempotency-Key")

		if id == "" {
			SendErrorResponse(w, ErrMissingIdempotencyHeader)
			return
		}

		ctx := SetIdempotencyContext(r.Context(), id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
