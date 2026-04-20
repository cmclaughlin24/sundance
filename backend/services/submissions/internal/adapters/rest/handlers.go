package rest

import (
	"net/http"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/httputil"
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
			httputil.SendErrorResponse(w, res.err)
			return
		}

		dtos := make([]*dto.SubmissionResponse, 0, len(res.data))
		for _, submission := range res.data {
			dtos = append(dtos, dto.SubmissionToResponse(submission))
		}

		httputil.SendJsonResponse(w, http.StatusOK, dtos)
	}
}

func (h *handlers) getSubmissionByReferenceID(w http.ResponseWriter, r *http.Request) {
	referenceID := h.getReferenceIdPathValue(r)
	resultChan := make(chan result[*domain.Submission], 1)
	query := ports.NewFindByIdQuery(referenceID)

	go func() {
		defer close(resultChan)
		submission, err := h.app.Services.Submissions.FindByReferenceId(r.Context(), query)
		resultChan <- result[*domain.Submission]{submission, err}
	}()

	select {
	case <-r.Context().Done():
		return
	case res := <-resultChan:
		if res.err != nil {
			httputil.SendErrorResponse(w, res.err)
			return
		}

		httputil.SendJsonResponse(w, http.StatusOK, dto.SubmissionToResponse(res.data))
	}
}

func (h *handlers) createSubmission(w http.ResponseWriter, r *http.Request) {}

func (h *handlers) getSubmissionAttempts(w http.ResponseWriter, r *http.Request) {}

func (h *handlers) getSubmissionStatus(w http.ResponseWriter, r *http.Request) {}

func (h *handlers) replaySubmission(w http.ResponseWriter, r *http.Request) {}

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
