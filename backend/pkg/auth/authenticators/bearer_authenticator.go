package authenticators

import (
	"errors"
	"net/http"
	"strings"
	"sundance/backend/pkg/auth"
)

func NewBearerAuthenticator[T auth.Claims](validator TokenValidator[T]) auth.Authenticator {
	return func(r *http.Request) (auth.Claims, error) {
		token := strings.Split(r.Header.Get("Authorization"), " ")

		if len(token) != 2 {
			return nil, errors.New("authorization header is not a tuple with type and token")
		}

		if token[0] != "Bearer" {
			return nil, errors.New("authorization header is not type \"Bearer\"")
		}

		return validator.Validate(r.Context(), token[1])
	}
}
