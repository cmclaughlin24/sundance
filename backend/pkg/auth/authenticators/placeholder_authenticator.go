package authenticators

import (
	"net/http"

	"github.com/cmclaughlin24/sundance/backend/pkg/auth"
)

type PlaceholderClaims struct {
	subject string
}

func (c PlaceholderClaims) GetSubject() string {
	return c.subject
}

func NewPlaceholderAuthenticator(subject string) auth.Authenticator {
	return func(r *http.Request) (auth.Claims, error) {
		return &PlaceholderClaims{subject}, nil
	}
}
