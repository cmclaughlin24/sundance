package authenticators

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

type PingFedClaims struct {
	jwt.RegisteredClaims
}

type PingFedTokenValidator struct {
	audience string
	issuer   string
	jwk      keyfunc.Keyfunc
}

func NewPingFedTokenValidator(opts ...AuthenticatorOption) (*PingFedTokenValidator, error) {
	settings := make(authenticatorSettings)

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

	jwkURI, ok := settings[settingsKeyJWK]
	if !ok {
		return nil, ErrMissingAuthenticatorSettings(settingsKeyJWK)
	}

	k, err := keyfunc.NewDefaultCtx(context.Background(), []string{jwkURI})
	if err != nil {
		return nil, fmt.Errorf("failed to initialise JWKs keyfunc: %w", err)
	}

	return &PingFedTokenValidator{
		audience: audience,
		issuer:   issuer,
		jwk:      k,
	}, nil
}

func (v *PingFedTokenValidator) Validate(ctx context.Context, t string) (*PingFedClaims, error) {
	claims := &PingFedClaims{}

	token, err := jwt.ParseWithClaims(
		t,
		claims,
		v.jwk.Keyfunc,
		jwt.WithAudience(v.audience),
		jwt.WithExpirationRequired(),
		jwt.WithIssuer(v.issuer),
		jwt.WithIssuedAt(),
		jwt.WithLeeway(5*time.Second),
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token is not valid")
	}

	return claims, nil
}
