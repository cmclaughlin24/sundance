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
	return ctx.Value(ClaimsKey).(Claims)
}
