package auth

import "context"

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
	}

	return claims
}
