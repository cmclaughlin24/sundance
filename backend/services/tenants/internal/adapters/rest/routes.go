package rest

import (
	"fmt"
	"net/http"
	"regexp"

	"sundance/backend/pkg/auth"
	"sundance/backend/pkg/auth/authenticators"
	"sundance/backend/pkg/common/httputil"
	"sundance/backend/services/tenants/internal/adapters/rest/docs"
	"sundance/backend/services/tenants/internal/adapters/rest/handlers"
	"sundance/backend/services/tenants/internal/core"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer JWT. Format: "Bearer <token>"

func NewRoutes(app *core.Application, host string, _ auth.AuthOptions) http.Handler {
	h := handlers.NewHandlers(app)
	mux := chi.NewRouter()

	placeHolderTokenValidator := &authenticators.PlaceholderTokenValidator{}
	bearerAuthenticator := authenticators.NewBearerAuthenticator(placeHolderTokenValidator)

	mux.Use(middleware.RequestID)
	mux.Use(httplog.RequestLogger(app.Logger, &httplog.Options{
		Schema: httplog.SchemaOTEL,
	}))
	mux.Use(httputil.NewCorrelationIDMiddleware("X-Correlation-ID"))
	mux.Use(httputil.NewRequestDateMiddleware("X-Request-Date"))

	mux.Route("/api/v1", func(routes chi.Router) {
		routes.Use(auth.NewMiddleware(bearerAuthenticator))

		routes.Route("/tenants", func(tenantsRoutes chi.Router) {
			tenantsRoutes.Get("/", h.GetTenants)
			tenantsRoutes.Post("/", h.CreateTenant)

			tenantsRoutes.Route("/{tenantId}", func(tenantRoutes chi.Router) {
				tenantRoutes.Get("/", h.GetTenant)
				tenantRoutes.Put("/", h.UpdateTenant)
				tenantRoutes.Delete("/", h.DeleteTenant)
			})
		})

		routes.Route("/data-sources", func(dataSourcesRoutes chi.Router) {
			dataSourcesRoutes.Use(httputil.NewTenantMiddleware("X-Tenant-ID"))

			dataSourcesRoutes.Get("/", h.GetDataSources)
			dataSourcesRoutes.Post("/", h.CreateDataSource)

			dataSourcesRoutes.Route("/{dataSourceId}", func(dataSourceRoutes chi.Router) {
				dataSourceRoutes.Get("/", h.GetDataSource)
				dataSourceRoutes.Put("/", h.UpdateDataSource)
				dataSourceRoutes.Delete("/", h.DeleteDataSource)
				dataSourceRoutes.Get("/look-ups", h.GetLookups)
			})
		})
	})

	re := regexp.MustCompile(`https?://`)
	docs.SwaggerInfo.Host = re.ReplaceAllString(host, "")
	docs.SwaggerInfo.Title = "Forms Hub SaaS | Tenants Service"
	docs.SwaggerInfo.Description = "The Tenants Service manages multi-tenant configurations and their associated data sources for the Form Builder SaaS platform. It provides CRUD operations for tenants and data sources, where data sources supply lookup key-value pairs (e.g., for populating dropdowns) via static data, scheduled external fetches, or on-demand webhook calls."
	docs.SwaggerInfo.Version = "1.0.0"
	docs.SwaggerInfo.BasePath = "/api/v1"

	mux.Route("/swagger", func(r chi.Router) {
		r.Get("/*", httpSwagger.Handler(httpSwagger.URL(fmt.Sprintf("%s/swagger/doc.json", host))))
	})

	return mux
}
