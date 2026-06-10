package authenticators

import (
	"sundance/backend/pkg/auth"

	"github.com/golang-jwt/jwt/v5"
)

type PingFedClaims struct {
	jwt.RegisteredClaims
}

type PingFedValidator struct {
	audience string
	issuer   string
}

func NewPingFedValidator(opts ...AuthenticatorOption) (*PingFedValidator, error) {
	settings := make(autenticatorSettings)

	for _, opt := range opts {
		opt(settings)
	}

	audience, ok := settings[settingsKeyAudience]
	if !ok {
		return nil, ErrMissingAuthenticatorSettings(settingsKeyAudience)
	}

	issuer, ok := settings[settingsKeyIssuer]
	if !ok {
		return nil, ErrMissingAuthenticatorSettings(settingsKeyIssuer)
	}

	_, ok = settings[settingsKeyJWK]
	if !ok {
		return nil, ErrMissingAuthenticatorSettings(settingsKeyJWK)
	}

	return &PingFedValidator{
		audience: audience,
		issuer:   issuer,
	}, nil
}

func (v *PingFedValidator) Validate(t string) (auth.Claims, error) {
	return nil, nil
}
