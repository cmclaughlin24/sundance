package tenants

import (
	"net/http"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/httputil"
)

func NewMiddleware(header string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get(header)

			if id == "" {
				httputil.SendErrorResponse(w, ErrMissingTenantID)
				return
			}

			ctx := SetTenantContext(r.Context(), id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
