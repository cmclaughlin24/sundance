package rest

import (
	"fmt"
	"net/http"
	"regexp"

	"sundance/backend/pkg/auth"
	"sundance/backend/pkg/auth/authenticators"
	"sundance/backend/pkg/common/httputil"
	"sundance/backend/services/forms/internal/adapters/rest/docs"
	"sundance/backend/services/forms/internal/core"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRoutes(app *core.Application, host string) http.Handler {
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

	re := regexp.MustCompile(`https?://`)
	docs.SwaggerInfo.Host = re.ReplaceAllString(host, "")
	docs.SwaggerInfo.Title = "Form Builder SaaS | Forms Service"
	docs.SwaggerInfo.Description = "The Forms Service manages form definitions and their versioned schemas for the Form Builder SaaS platform. It provides CRUD operations for forms and versions, where versions define the structure of a form (pages, sections, fields, and validation rules) and follow a lifecycle of draft, published, and retired statuses."
	docs.SwaggerInfo.Version = "1.0.0"
	docs.SwaggerInfo.BasePath = "/api/v1"

	mux.Route("/swagger", func(r chi.Router) {
		r.Get("/*", httpSwagger.Handler(httpSwagger.URL(fmt.Sprintf("%s/swagger/doc.json", host))))
	})

	return mux
}
