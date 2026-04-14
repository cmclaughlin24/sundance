package rest

import (
	"net/http"

	"github.com/cmclaughlin24/sundance/submissions/internal/core"
)

func NewRoutes(app *core.Application) http.Handler {
	h := newHandlers(app)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/submissions", h.getSubmissions)
	mux.HandleFunc("POST /api/v1/submissions", nil)
	mux.HandleFunc("GET /api/v1/submissions/{referenceId}", h.getSubmissionByReferenceID)
	mux.HandleFunc("GET /api/v1/submissions/{referenceId}/attempts", nil)
	mux.HandleFunc("GET /api/v1/submissions/{referenceId}/status", nil)
	mux.HandleFunc("POST /api/v1/submissions/{referenceId}/replay", nil)

	return mux
}
