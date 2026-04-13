package rest

import (
	"net/http"

	"github.com/cmclaughlin24/sundance/submissions/internal/core"
)

func NewRoutes(app *core.Application) http.Handler {
	_ = newHandlers(app)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/submissions", nil)
	mux.HandleFunc("POST /api/v1/submissions", nil)
	mux.HandleFunc("GET /api/v1/submissions/{referenceId}", nil)
	mux.HandleFunc("GET /api/v1/submissions/{referenceId}/attempts", nil)
	mux.HandleFunc("GET /api/v1/submissions/{referenceId}/status", nil)
	mux.HandleFunc("POST /api/v1/submissions/{referenceId}/replay", nil)

	return mux
}
