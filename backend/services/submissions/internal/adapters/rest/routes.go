package rest

import (
	"net/http"

	"github.com/cmclaughlin24/sundance/submissions/internal/core"
	"github.com/go-chi/chi/v5"
)

func NewRoutes(app *core.Application) http.Handler {
	h := newHandlers(app)
	mux := chi.NewRouter()

	mux.Route("/api/v1", func(routes chi.Router) {
		routes.Route("/submissions", func(submissionsRoutes chi.Router) {
			submissionsRoutes.Get("/", h.getSubmissions)
			submissionsRoutes.Post("/", h.createSubmission)

			submissionsRoutes.Route("/{referenceId}", func(submissionRoutes chi.Router) {
				submissionRoutes.Get("/", h.getSubmissionByReferenceID)
				submissionRoutes.Get("/attempts", h.getSubmissionAttempts)
				submissionRoutes.Get("/status", h.getSubmissionStatus)
				submissionRoutes.Post("/replay", h.replaySubmission)
			})
		})
	})

	return mux
}
