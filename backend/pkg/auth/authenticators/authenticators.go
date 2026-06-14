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

type authenticatorSettings map[settingsKey]string

type AuthenticatorOption func(authenticatorSettings)

func WithAudience(audience string) AuthenticatorOption {
	return func(as authenticatorSettings) {
		as[settingsKeyAudience] = audience
	}
}

func WithIssuer(issuer string) AuthenticatorOption {
	return func(as authenticatorSettings) {
		as[settingsKeyIssuer] = issuer
	}
}

func WithJWK(jwk string) AuthenticatorOption {
	return func(as authenticatorSettings) {
		as[settingsKeyJWK] = jwk
	}
}
