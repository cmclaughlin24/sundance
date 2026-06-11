package authenticators

import (
	"net/http"

	"sundance/backend/pkg/auth"
)

type PlaceholderClaims struct {
	subject string
}

func (c PlaceholderClaims) GetSubject() (string, error) {
	return c.subject, nil
}

func NewPlaceholderAuthenticator(subject string) auth.Authenticator {
	return func(r *http.Request) (auth.Claims, error) {
		return &PlaceholderClaims{subject}, nil
	}
}
