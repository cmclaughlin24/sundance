package rest

import (
	"net/http"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/tenants"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"
)

func NewRoutes(app *core.Application) http.Handler {
	h := newHandlers(app)
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(httplog.RequestLogger(app.Logger, &httplog.Options{
		Schema: httplog.SchemaOTEL,
	}))
	mux.Use(tenants.NewMiddleware("X-Tenant-ID"))

	mux.Route("/api/v1", func(routes chi.Router) {
		routes.Route("/submissions", func(submissionsRoutes chi.Router) {
			submissionsRoutes.Get("/", h.getSubmissions)
			submissionsRoutes.Post("/", h.createSubmission)

			submissionsRoutes.Route("/{submissionId}", func(submissionRoutes chi.Router) {
				submissionRoutes.Post("/replay", h.replaySubmission)
			})

			submissionsRoutes.Route("/by-reference/{referenceId}", func(submissionRoutes chi.Router) {
				submissionRoutes.Get("/", h.getSubmissionByReferenceID)
				submissionRoutes.Get("/status", h.getSubmissionStatus)
			})
		})
	})

	return mux
}
