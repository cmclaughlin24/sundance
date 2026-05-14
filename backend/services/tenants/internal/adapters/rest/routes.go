package rest

import (
	"net/http"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/httputil"
	_ "github.com/cmclaughlin24/sundance/backend/services/tenants/internal/adapters/rest/docs"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"
	"github.com/swaggo/http-swagger/v2"
)

// @title 			Form Builder SaaS | Tenants Service
// @version 		1.0.0
// @host			localhost:80
// @BasePath 		/api/v1

func NewRoutes(app *core.Application) http.Handler {
	h := newHandlers(app)
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(httplog.RequestLogger(app.Logger, &httplog.Options{
		Schema: httplog.SchemaOTEL,
	}))

	mux.Route("/api/v1", func(routes chi.Router) {
		routes.Route("/tenants", func(tenantsRoutes chi.Router) {
			tenantsRoutes.Get("/", h.getTenants)
			tenantsRoutes.Post("/", h.createTenant)

			tenantsRoutes.Route("/{tenantId}", func(tenantRoutes chi.Router) {
				tenantRoutes.Get("/", h.getTenant)
				tenantRoutes.Put("/", h.updateTenant)
				tenantRoutes.Delete("/", h.deleteTenant)
			})
		})

		routes.Route("/data-sources", func(dataSourcesRoutes chi.Router) {
			dataSourcesRoutes.Use(httputil.NewTenantMiddleware("X-Tenant-ID"))

			dataSourcesRoutes.Get("/", h.getDataSources)
			dataSourcesRoutes.Post("/", h.createDataSource)

			dataSourcesRoutes.Route("/{dataSourceId}", func(dataSourceRoutes chi.Router) {
				dataSourceRoutes.Get("/", h.getDataSource)
				dataSourceRoutes.Put("/", h.updateDataSource)
				dataSourceRoutes.Delete("/", h.deleteDataSource)
				dataSourceRoutes.Get("/look-ups", h.getLookups)
			})
		})
	})

	mux.Route("/swagger", func(r chi.Router) {
		r.Get("/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8082/swagger/doc.json")))
	})

	return mux
}
