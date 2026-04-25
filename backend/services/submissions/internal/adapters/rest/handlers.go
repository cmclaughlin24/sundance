package rest

import (
	"net/http"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/httputil"
	"github.com/cmclaughlin24/sundance/backend/pkg/common/tenants"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/adapters/rest/dto"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/ports"
	"github.com/go-chi/chi/v5"
)

type result[T any] struct {
	data T
	err  error
}

type handlers struct {
	app *core.Application
}

func newHandlers(app *core.Application) *handlers {
	return &handlers{
		app: app,
	}
}

func (h *handlers) getSubmissions(w http.ResponseWriter, r *http.Request) {
	resultChan := make(chan result[[]*domain.Submission], 1)

	go func() {
		defer close(resultChan)
		submissions, err := h.app.Services.Submissions.Find(r.Context())
		resultChan <- result[[]*domain.Submission]{submissions, err}
	}()

	select {
	case <-r.Context().Done():
		return
	case res := <-resultChan:
		if res.err != nil {
			h.sendErrorResponse(w, res.err)
			return
		}

		dtos := make([]*dto.SubmissionResponse, 0, len(res.data))
		for _, submission := range res.data {
			dtos = append(dtos, dto.SubmissionToResponse(submission))
		}

		httputil.SendJSONResponse(w, http.StatusOK, dtos)
	}
}

func (h *handlers) getSubmissionByReferenceID(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantFromContext(r)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	referenceID := h.getReferenceIdPathValue(r)
	resultChan := make(chan result[*domain.Submission], 1)
	query := ports.NewFindByIDQuery(tenantID, referenceID)

	go func() {
		defer close(resultChan)
		submission, err := h.app.Services.Submissions.FindByReferenceID(r.Context(), query)
		resultChan <- result[*domain.Submission]{submission, err}
	}()

	select {
	case <-r.Context().Done():
		return
	case res := <-resultChan:
		if res.err != nil {
			h.sendErrorResponse(w, res.err)
			return
		}

		httputil.SendJSONResponse(w, http.StatusOK, dto.SubmissionToResponse(res.data))
	}
}

func (h *handlers) createSubmission(w http.ResponseWriter, r *http.Request) {}

func (h *handlers) getSubmissionAttempts(w http.ResponseWriter, r *http.Request) {}

func (h *handlers) getSubmissionStatus(w http.ResponseWriter, r *http.Request) {}

func (h *handlers) replaySubmission(w http.ResponseWriter, r *http.Request) {}

func (h *handlers) getTenantFromContext(r *http.Request) (string, error) {
	tenantID, err := tenants.TenantFromContext(r.Context())

	if err != nil {
		return "", err
	}

	return tenantID, nil
}

func (h *handlers) getReferenceIdPathValue(r *http.Request) domain.ReferenceID {
	id := chi.URLParam(r, "referenceId")
	return domain.ReferenceID(id)
}

func (h *handlers) sendErrorResponse(w http.ResponseWriter, err error) {
	switch {
	default:
		httputil.SendErrorResponse(w, err)
	}
}
