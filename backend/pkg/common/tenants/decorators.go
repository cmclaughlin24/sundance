package tenants

import (
	"net/http"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/httputil"
)

type TenantHandlerFunc func(http.ResponseWriter, *http.Request, string)

func WithTenant(fn TenantHandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID, err := TenantFromContext(r.Context())

		if err != nil {
			httputil.SendErrorResponse(w, err)
			return
		}

		fn(w, r, tenantID)
	})
}
