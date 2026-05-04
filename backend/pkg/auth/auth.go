package auth

import (
	"net/http"
	"reflect"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/httputil"
)

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

				valOfClaims := reflect.ValueOf(claims)

				if valOfClaims.IsValid() && !valOfClaims.IsZero() {
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
