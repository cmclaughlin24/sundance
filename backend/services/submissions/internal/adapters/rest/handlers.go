package rest

import (
	"net/http"

	"github.com/cmclaughlin24/sundance/common"
	"github.com/cmclaughlin24/sundance/submissions/internal/core"
	"github.com/cmclaughlin24/sundance/submissions/internal/core/domain"
	"github.com/cmclaughlin24/sundance/submissions/internal/core/ports"
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
			common.SendErrorResponse(w, res.err)
			return
		}

		// TODO: Convert submission domain object to dto.
		common.SendJsonResponse(w, http.StatusOK, res.data)
	}
}

func (h *handlers) getSubmissionByReferenceID(w http.ResponseWriter, r *http.Request) {
	referenceID := h.getReferenceIdPathValue(r)
	resultChan := make(chan result[*domain.Submission], 1)

	query, err := ports.NewFindByIdQuery(referenceID, "")
	if err != nil {
		common.SendErrorResponse(w, err)
		return
	}

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
			common.SendErrorResponse(w, res.err)
			return
		}

		// TODO: Convert submission domain object to dto.
		common.SendJsonResponse(w, http.StatusOK, res.data)
	}
}

func (h *handlers) getReferenceIdPathValue(r *http.Request) domain.ReferenceID {
	id := r.PathValue("referenceId")
	return domain.ReferenceID(id)
}
