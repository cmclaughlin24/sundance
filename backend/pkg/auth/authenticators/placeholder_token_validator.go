package authenticators

import (
	"context"
)

type PlaceholderClaims struct {
	subject string
}

func (c PlaceholderClaims) GetSubject() (string, error) {
	return c.subject, nil
}

type PlaceholderTokenValidator struct{}

func (v *PlaceholderTokenValidator) Validate(_ context.Context, subject string) (*PlaceholderClaims, error) {
	return &PlaceholderClaims{subject}, nil
}
