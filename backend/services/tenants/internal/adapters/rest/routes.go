package rest

import (
	"fmt"
	"net/http"
	"regexp"

	"sundance/backend/pkg/common/httputil"
	"sundance/backend/services/tenants/internal/adapters/rest/docs"
	"sundance/backend/services/tenants/internal/core"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRoutes(app *core.Application, host string) http.Handler {
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

	re := regexp.MustCompile(`https?://`)
	docs.SwaggerInfo.Host = re.ReplaceAllString(host, "")
	docs.SwaggerInfo.Title = "Form Builder SaaS | Tenants Service"
	docs.SwaggerInfo.Description = "The Tenants Service manages multi-tenant configurations and their associated data sources for the Form Builder SaaS platform. It provides CRUD operations for tenants and data sources, where data sources supply lookup key-value pairs (e.g., for populating dropdowns) via static data, scheduled external fetches, or on-demand webhook calls."
	docs.SwaggerInfo.Version = "1.0.0"
	docs.SwaggerInfo.BasePath = "/api/v1"

	mux.Route("/swagger", func(r chi.Router) {
		r.Get("/*", httpSwagger.Handler(httpSwagger.URL(fmt.Sprintf("%s/swagger/doc.json", host))))
	})

	return mux
}
