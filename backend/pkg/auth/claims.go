package auth

import (
	"context"
	"log/slog"
	"os"
)

type ClaimsContextKey string

const ClaimsKey ClaimsContextKey = "claims"

type Claims interface {
	GetSubject() string
}

func SetClaimsContext(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, ClaimsKey, claims)
}

func GetClaimsFromContext(ctx context.Context) Claims {
	claims, ok := ctx.Value(ClaimsKey).(Claims)

	if !ok {
		slog.ErrorContext(ctx, "failed to get claims from context; claims not found or of wrong type")
		os.Exit(1)
	}

	return claims
}
