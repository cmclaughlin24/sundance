package rest

import (
	"net/http"

	"github.com/cmclaughlin24/sundance/backend/pkg/auth"
	"github.com/cmclaughlin24/sundance/backend/pkg/auth/authenticators"
	"github.com/cmclaughlin24/sundance/backend/pkg/common/httputil"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/adapters/rest/docs"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRoutes(app *core.Application) http.Handler {
	h := newHandlers(app)
	mux := chi.NewRouter()
	placeholderAuthenticator := authenticators.NewPlaceholderAuthenticator("placeholder") // TODO: Remove for a proper authentication implementation.

	mux.Use(middleware.RequestID)
	mux.Use(httplog.RequestLogger(app.Logger, &httplog.Options{
		Schema: httplog.SchemaOTEL,
	}))

	mux.Route("/api/v1", func(routes chi.Router) {
		routes.Use(auth.NewMiddleware(placeholderAuthenticator))

		routes.Route("/forms", func(formsRoutes chi.Router) {
			formsRoutes.Use(httputil.NewTenantMiddleware("X-Tenant-ID"))

			formsRoutes.Get("/", h.getForms)
			formsRoutes.Post("/", h.createForm)

			formsRoutes.Route("/{formId}", func(formRoutes chi.Router) {
				formRoutes.Get("/", h.getForm)
				formRoutes.Put("/", h.updateForm)
				formRoutes.Delete("/", h.deleteForm)

				formRoutes.Route("/versions", func(versionsRoutes chi.Router) {
					versionsRoutes.Get("/", h.getVersions)
					versionsRoutes.Post("/", h.createVersion)

					versionsRoutes.Route("/{versionId}", func(versionRoutes chi.Router) {
						versionRoutes.Get("/", h.getVersion)
						versionRoutes.Put("/", h.updateVersion)
						versionRoutes.Post("/publish", h.publishVersion)
						versionRoutes.Post("/retire", h.retireVersion)
					})
				})
			})
		})
	})

	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.Title = "Form Builder SaaS | Forms Service"
	docs.SwaggerInfo.Description = ""
	docs.SwaggerInfo.Version = "1.0.0"
	docs.SwaggerInfo.BasePath = "/api/v1"

	mux.Route("/swagger", func(r chi.Router) {
		r.Get("/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/swagger/doc.json")))
	})

	return mux
}
