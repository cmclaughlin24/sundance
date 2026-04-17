package rest

import (
	"net/http"

	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core"
	"github.com/go-chi/chi/v5"
)

func NewRoutes(app *core.Application) http.Handler {
	h := newHandlers(app)
	mux := chi.NewRouter()

	mux.Route("/api/v1", func(routes chi.Router) {
		routes.Route("/tenants", func(tenantsRoutes chi.Router) {
			tenantsRoutes.Get("/", h.getTenants)
			tenantsRoutes.Post("/", h.createTenant)

			tenantsRoutes.Route("/{tenantId}", func(tenantRoutes chi.Router) {
				tenantRoutes.Get("/", h.getTenant)
				tenantRoutes.Put("/", h.updateTenant)
				tenantRoutes.Delete("/", h.removeTenant)

				tenantRoutes.Route("/data-sources", func(dataSourcesRoutes chi.Router) {
					dataSourcesRoutes.Get("/", h.getDataSources)
					dataSourcesRoutes.Post("/", h.createDataSource)

					dataSourcesRoutes.Route("/{dataSourceId}", func(dataSourceRoutes chi.Router) {
						dataSourceRoutes.Get("/", h.getDataSource)
						dataSourceRoutes.Put("/", h.updateDataSource)
						dataSourceRoutes.Put("/", h.removeDataSource)
						dataSourceRoutes.Get("/look-ups", h.getDataSourceLookup)
					})
				})

			})
		})

	})

	return mux
}
