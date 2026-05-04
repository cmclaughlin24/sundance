package rest

import (
	"net/http"

	"github.com/cmclaughlin24/sundance/backend/pkg/auth"
	"github.com/cmclaughlin24/sundance/backend/pkg/auth/authenticators"
	"github.com/cmclaughlin24/sundance/backend/pkg/common/tenants"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRoutes(app *core.Application) http.Handler {
	h := newHandlers(app)
	mux := chi.NewRouter()
	placeholderAuthenticator := authenticators.NewPlaceholderAuthenticator("placholder") // TODO: Remove for a proper authentication implementation.

	mux.Use(middleware.RequestID)
	mux.Use(middleware.Logger)
	mux.Use(tenants.NewMiddleware("X-Tenant-ID"))
	mux.Use(auth.NewMiddleware(placeholderAuthenticator))

	mux.Route("/api/v1", func(routes chi.Router) {
		routes.Route("/forms", func(formsRoutes chi.Router) {
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

	return mux
}
