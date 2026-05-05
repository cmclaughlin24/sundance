package rest

import (
	"fmt"
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
	tenantID, err := h.getTenantFromContext(r)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	query := ports.NewFindSubmissionsQuery(tenantID)
	resultChan := make(chan result[[]*domain.Submission], 1)

	go func() {
		defer close(resultChan)
		submissions, err := h.app.Services.Submissions.Find(r.Context(), query)
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

	referenceID := h.getReferenceIDPathValue(r)
	query := ports.NewFindSubmissionByIDQuery(tenantID, referenceID)
	resultChan := make(chan result[*domain.Submission], 1)

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

func (h *handlers) createSubmission(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantFromContext(r)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	var body dto.SubmissionRequest
	if err := httputil.ReadValidateJSONPayload(r, &body); err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	command := ports.NewCreateSubmissionCommand(
		tenantID,
		body.FormID,
		body.VersionID,
		body.Payload,
	)
	resultChan := make(chan result[*domain.Submission])

	go func() {
		defer close(resultChan)
		submission, err := h.app.Services.Submissions.Create(r.Context(), command)
		resultChan <- result[*domain.Submission]{data: submission, err: err}
	}()

	select {
	case <-r.Context().Done():
		return
	case res := <-resultChan:
		if res.err != nil {
			h.sendErrorResponse(w, res.err)
			return
		}

		httputil.SendJSONResponse(w, http.StatusCreated, httputil.APIResponse[*dto.SubmissionResponse]{
			Message: "Successfully created!",
			Data:    dto.SubmissionToResponse(res.data),
		})
	}
}

func (h *handlers) getSubmissionStatus(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantFromContext(r)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	referenceID := h.getReferenceIDPathValue(r)
	query := ports.NewFindSubmissionByIDQuery(tenantID, referenceID)
	resultChan := make(chan result[*domain.Submission], 1)

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

		httputil.SendJSONResponse(w, http.StatusOK, struct {
			Status string `json:"status"`
		}{
			Status: string(res.data.Status),
		})
	}
}

func (h *handlers) replaySubmission(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantFromContext(r)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	id := chi.URLParam(r, "submissionId")
	command := ports.NewReplaySubmissionCommand(
		tenantID,
		domain.SubmissionID(id),
	)
	resultChan := make(chan result[any])

	go func() {
		defer close(resultChan)
		err := h.app.Services.Submissions.Replay(r.Context(), command)
		resultChan <- result[any]{data: nil, err: err}
	}()

	select {
	case <-r.Context().Done():
		return
	case res := <-resultChan:
		if res.err != nil {
			h.sendErrorResponse(w, res.err)
			return
		}

		httputil.SendJSONResponse(w, http.StatusCreated, httputil.APIResponse[*dto.SubmissionResponse]{
			Message: fmt.Sprintf("Successfully replayed submission %s", id),
		})
	}
}

func (h *handlers) getTenantFromContext(r *http.Request) (string, error) {
	tenantID, err := tenants.TenantFromContext(r.Context())

	if err != nil {
		return "", err
	}

	return tenantID, nil
}

func (h *handlers) getReferenceIDPathValue(r *http.Request) domain.ReferenceID {
	id := chi.URLParam(r, "referenceId")
	return domain.ReferenceID(id)
}

func (h *handlers) sendErrorResponse(w http.ResponseWriter, err error) {
	switch {
	default:
		httputil.SendErrorResponse(w, err)
	}
}
