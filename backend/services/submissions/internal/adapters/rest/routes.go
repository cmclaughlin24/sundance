package rest

import (
	"net/http"

	"github.com/cmclaughlin24/sundance/backend/pkg/auth"
	"github.com/cmclaughlin24/sundance/backend/pkg/auth/authenticators"
	"github.com/cmclaughlin24/sundance/backend/pkg/common/httputil"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/adapters/rest/docs"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core"
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
		routes.Use(httputil.NewTenantMiddleware("X-Tenant-ID"))
		routes.Use(auth.NewMiddleware(placeholderAuthenticator))

		routes.Route("/submissions", func(submissionsRoutes chi.Router) {
			submissionsRoutes.Get("/", h.getSubmissions)
			submissionsRoutes.With(httputil.IdempotencyMiddleware).Post("/", h.createSubmission)

			submissionsRoutes.Route("/{submissionId}", func(submissionRoutes chi.Router) {
				submissionRoutes.Post("/replay", h.replaySubmission)
			})

			submissionsRoutes.Route("/by-reference/{referenceId}", func(submissionRoutes chi.Router) {
				submissionRoutes.Get("/", h.getSubmissionByReferenceID)
				submissionRoutes.Get("/status", h.getSubmissionStatus)
			})
		})
	})

	docs.SwaggerInfo.Host = "localhost:8081"
	docs.SwaggerInfo.Title = "Form Builder SaaS | Submissions Service"
	docs.SwaggerInfo.Description = "The Submissions Service handles form submission intake and lifecycle tracking for the Form Builder SaaS platform. It provides idempotent submission creation with asynchronous processing, status tracking via reference IDs, and the ability to replay submission events for reprocessing by downstream consumers."
	docs.SwaggerInfo.Version = "1.0.0"
	docs.SwaggerInfo.BasePath = "/api/v1"

	mux.Route("/swagger", func(r chi.Router) {
		r.Get("/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8081/swagger/doc.json")))
	})

	return mux
}
