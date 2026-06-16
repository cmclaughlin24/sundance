package auth

import (
	"net/http"

	"sundance/backend/pkg/common/httputil"
)

type OAuth2 struct {
	Audience string `json:"audience" env:"AUDIENCE"`
	Issuer   string `json:"issuer" env:"ISSUER"`
	JWK      string `json:"jwk" env:"JWK"`
}

type AuthOptions struct {
	OAuth2 OAuth2 `json:"oauth2" envPrefix:"OAUTH2_"`
}

type Authenticator = func(*http.Request) (Claims, error)

func NewMiddleware(authenticators ...Authenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, authenticator := range authenticators {
				claims, err := authenticator(r)

				if err != nil {
					// TODO: Add a debug log message indicating the authenticator method failed.
					continue
				}

				if claims != nil {
					r := r.WithContext(SetClaimsContext(r.Context(), claims))
					next.ServeHTTP(w, r)
					return
				}
			}

			httputil.SendJSONResponse(w, http.StatusUnauthorized, httputil.APIErrorResponse{
				Message:    "Unauthorized",
				Error:      "Unauthorized",
				StatusCode: http.StatusUnauthorized,
			})
		})
	}
}
