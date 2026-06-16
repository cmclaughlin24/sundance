package rest

import (
	"fmt"
	"net/http"
	"regexp"

	"sundance/backend/pkg/auth"
	"sundance/backend/pkg/auth/authenticators"
	"sundance/backend/pkg/common/httputil"
	"sundance/backend/services/forms/internal/adapters/rest/docs"
	"sundance/backend/services/forms/internal/adapters/rest/handlers"
	"sundance/backend/services/forms/internal/core"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v3"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type ServerOptions struct {
	AllowedOrigins []string         `json:"allowedOrigins"`
	Auth           auth.AuthOptions `json:"security"`
}

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer JWT. Format: "Bearer <token>"

func NewRoutes(app *core.Application, host string, options ServerOptions) http.Handler {
	h := handlers.NewHandlers(app)
	mux := chi.NewRouter()

	placeHolderTokenValidator := &authenticators.PlaceholderTokenValidator{}
	bearerAuthenticator := authenticators.NewBearerAuthenticator(placeHolderTokenValidator)

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: options.AllowedOrigins,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))
	mux.Use(middleware.RequestID)
	mux.Use(httplog.RequestLogger(app.Logger, &httplog.Options{
		Schema: httplog.SchemaOTEL,
	}))
	mux.Use(httputil.NewCorrelationIDMiddleware("X-Correlation-ID"))
	mux.Use(httputil.NewRequestDateMiddleware("X-Request-Date"))

	mux.Route("/api/v1", func(routes chi.Router) {
		routes.Use(httputil.NewTenantMiddleware("X-Tenant-ID"))
		routes.Use(auth.NewMiddleware(bearerAuthenticator))

		routes.Route("/forms", func(formsRoutes chi.Router) {
			formsRoutes.Get("/", h.GetForms)
			formsRoutes.Post("/", h.CreateForm)

			formsRoutes.Route("/{formId}", func(formRoutes chi.Router) {
				formRoutes.Get("/", h.GetForm)
				formRoutes.Put("/", h.UpdateForm)
				formRoutes.Delete("/", h.DeleteForm)

				formRoutes.Route("/versions", func(versionsRoutes chi.Router) {
					versionsRoutes.Get("/", h.GetFormVersions)
					versionsRoutes.Post("/", h.CreateFormVersion)

					versionsRoutes.Route("/{versionId}", func(versionRoutes chi.Router) {
						versionRoutes.Get("/", h.GetFormVersion)
						versionRoutes.Put("/", h.UpdateFormVersion)
						versionRoutes.Post("/publish", h.PublishFormVersion)
						versionRoutes.Post("/retire", h.RetireFormVersion)
					})
				})
			})
		})

		routes.Route("/submissions", func(submissionsRoutes chi.Router) {
			submissionsRoutes.Get("/", h.GetSubmissions)
			submissionsRoutes.With(httputil.IdempotencyMiddleware).Post("/", h.CreateSubmission)

			submissionsRoutes.Route("/{submissionId}", func(submissionRoutes chi.Router) {
				submissionRoutes.Post("/replay", h.ReplaySubmission)
			})

			submissionsRoutes.Route("/by-reference/{referenceId}", func(submissionRoutes chi.Router) {
				submissionRoutes.Get("/", h.GetSubmissionByReferenceID)
				submissionRoutes.Get("/status", h.GetSubmissionStatus)
			})
		})

		routes.Route("/tags", func(tagsRoutes chi.Router) {
			tagsRoutes.Get("/", h.GetTags)
			tagsRoutes.Post("/", h.CreateTag)

			tagsRoutes.Route("/{tagId}", func(tagRoutes chi.Router) {
				tagRoutes.Get("/", h.GetTag)
				tagRoutes.Put("/", h.UpdateTag)
				tagRoutes.Delete("/", h.DeleteTag)

				tagRoutes.Route("/versions", func(versionsRoutes chi.Router) {
					versionsRoutes.Get("/", h.GetTagVersions)
					versionsRoutes.Post("/", h.CreateTagVersion)

					versionsRoutes.Route("/{versionId}", func(versionRoutes chi.Router) {
						versionRoutes.Get("/", h.GetTagVersion)
						versionRoutes.Post("/deprecate", h.DeprecateTagVersion)
						versionRoutes.Post("/publish", h.PublishTagVersion)
						versionRoutes.Post("/retire", h.RetireTagVersion)
					})
				})
			})
		})
	})

	re := regexp.MustCompile(`https?://`)
	docs.SwaggerInfo.Host = re.ReplaceAllString(host, "")
	docs.SwaggerInfo.Title = "Forms Hub SaaS | Forms Service"
	docs.SwaggerInfo.Description = "The Forms Service is the system of record for form definitions, submissions, and tags in the Form Builder SaaS platform. **Forms** are composed of versioned schemas (pages, sections, fields, and validation rules) that follow a draft → published → retired lifecycle. **Submissions** are accepted asynchronously against a published form version, deduplicated via an idempotency key, and queryable by reference ID. **Tags** are tenant-scoped, stable identifiers used to attach semantic meaning to form fields and preserve historical associations across schema changes."
	docs.SwaggerInfo.Version = "1.0.0"
	docs.SwaggerInfo.BasePath = "/api/v1"

	mux.Route("/swagger", func(r chi.Router) {
		r.Get("/*", httpSwagger.Handler(httpSwagger.URL(fmt.Sprintf("%s/swagger/doc.json", host))))
	})

	return mux
}
