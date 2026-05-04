package rest

import (
	"net/http"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/tenants"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRoutes(app *core.Application) http.Handler {
	h := newHandlers(app)
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.Logger)
	mux.Use(tenants.TenantMiddleware("X-Tenant-ID"))

	mux.Route("/api/v1", func(routes chi.Router) {
		routes.Route("/forms", func(formsRoutes chi.Router) {
			formsRoutes.Get("/", tenants.WithTenant(h.getForms))
			formsRoutes.Post("/", tenants.WithTenant(h.createForm))

			formsRoutes.Route("/{formId}", func(formRoutes chi.Router) {
				formRoutes.Get("/", tenants.WithTenant(h.getForm))
				formRoutes.Put("/", tenants.WithTenant(h.updateForm))
				formRoutes.Delete("/", tenants.WithTenant(h.deleteForm))

				formRoutes.Route("/versions", func(versionsRoutes chi.Router) {
					versionsRoutes.Get("/", tenants.WithTenant(h.getVersions))
					versionsRoutes.Post("/", tenants.WithTenant(h.createVersion))

					versionsRoutes.Route("/{versionId}", func(versionRoutes chi.Router) {
						versionRoutes.Get("/", tenants.WithTenant(h.getVersion))
						versionRoutes.Put("/", tenants.WithTenant(h.updateVersion))
						versionRoutes.Post("/publish", tenants.WithTenant(h.publishVersion))
						versionRoutes.Post("/retire", tenants.WithTenant(h.retireVersion))
					})
				})
			})
		})
	})

	return mux
}
