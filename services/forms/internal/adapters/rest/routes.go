package rest

import (
	"net/http"

	"github.com/cmclaughlin24/sundance/forms/internal/core"
)

func NewRoutes(app *core.Application) http.Handler {
	h := newHandlers(app)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/forms", h.getForms)
	mux.HandleFunc("POST /api/v1/forms", h.createForm)
	mux.HandleFunc("GET /api/v1/forms/{formId}", h.getForm)
	mux.HandleFunc("PUT /api/v1/forms/{formId}", h.updateForm)

	mux.HandleFunc("GET /api/v1/forms/{formId}/versions", h.getVersions)
	mux.HandleFunc("POST /api/v1/forms/{formId}/versions", h.createVersion)
	mux.HandleFunc("GET /api/v1/forms/{formId}/versions/{versionId}", h.getVersion)
	mux.HandleFunc("PUT /api/v1/forms/{formId}/versions/{versionId}", h.updateVersion)
	mux.HandleFunc("POST /api/v1/forms/{formId}/versions/{versionId}/publish", h.publishVersion)
	mux.HandleFunc("POST /api/v1/forms/{formId}/versions/{versionId}/retire", h.retireVersion)

	return tenantMiddleware(mux)
}
