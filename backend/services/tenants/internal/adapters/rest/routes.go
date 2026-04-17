package rest

import (
	"net/http"

	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core"
)

func NewRoutes(app *core.Application) http.Handler {
	h := newHandlers(app)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/tenants", h.getTenants)
	mux.HandleFunc("POST /api/v1/tenants", h.createTenant)
	mux.HandleFunc("GET /api/v1/tenants/{tenantId}", h.getTenant)
	mux.HandleFunc("PUT /api/v1/tenants/{tenantId}", h.updateTenant)
	mux.HandleFunc("DELETE /api/v1/tenants/{tenantId}", h.removeTenant)

	mux.HandleFunc("GET /api/v1/tenants/{tenantId}/data-sources", h.getDataSources)
	mux.HandleFunc("POST /api/v1/tenants/{tenantId}/data-sources", h.createDataSource)
	mux.HandleFunc("GET /api/v1/tenants/{tenantId}/data-sources/{dataSourceId}", h.getDataSource)
	mux.HandleFunc("PUT /api/v1/tenants/{tenantId}/data-sources/{dataSourceId}", h.updateDataSource)
	mux.HandleFunc("DELETE /api/v1/tenants/{tenantId}/data-sources/{dataSourceId}", h.removeDataSource)
	mux.HandleFunc("GET /api/v1/tenants/{tenantId}/data-sources/{dataSourceId}/look-ups", h.getDataSourceLookup)

	return mux
}
