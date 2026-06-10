package authenticators

import (
	"context"
	"fmt"
	"sundance/backend/pkg/auth"
)

type TokenValidator[T auth.Claims] interface {
	Validate(context.Context, string) (T, error)
}

type settingsKey string

const (
	settingsKeyAudience settingsKey = "audience"
	settingsKeyIssuer   settingsKey = "issuer"
	settingsKeyJWK      settingsKey = "jwk"
)

var (
	ErrMissingAuthenticatorSettings = func(key settingsKey) error {
		return fmt.Errorf("missing authenticator setting: %s", key)
	}
)

type autenticatorSettings map[settingsKey]string

type AuthenticatorOption func(autenticatorSettings)

func WithAudience(audience string) AuthenticatorOption {
	return func(as autenticatorSettings) {
		as[settingsKeyAudience] = audience
	}
}

func WithIssuer(issuer string) AuthenticatorOption {
	return func(as autenticatorSettings) {
		as[settingsKeyIssuer] = issuer
	}
}

func WithJWK(jwk string) AuthenticatorOption {
	return func(as autenticatorSettings) {
		as[settingsKeyJWK] = jwk
	}
}
