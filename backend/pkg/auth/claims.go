package auth

import (
	"context"
	"errors"
	"log/slog"
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
		err := errors.New("failed to get claims from context; claims not found or of wrong type")
		slog.ErrorContext(ctx, "error", err)
		panic(err)
	}

	return claims
}
