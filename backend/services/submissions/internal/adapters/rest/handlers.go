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
	tenantID := tenants.TenantFromContext(r.Context())
	query := ports.NewFindSubmissionsQuery(tenantID)
	resultChan := make(chan result[[]*domain.Submission], 1)

	go func() {
		defer close(resultChan)
		submissions, err := h.app.Services.Submissions.Find(r.Context(), query)
		resultChan <- result[[]*domain.Submission]{submissions, err}
	}()

	select {
	case <-r.Context().Done():
		h.app.Logger.WarnContext(r.Context(), "context cancellation")
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
	tenantID := tenants.TenantFromContext(r.Context())
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
		h.app.Logger.WarnContext(r.Context(), "context cancellation")
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
	var body dto.SubmissionRequest
	if err := httputil.ReadValidateJSONPayload(r, &body); err != nil {
		h.app.Logger.WarnContext(r.Context(), "invalid request body", "error", err)
		h.sendErrorResponse(w, err)
		return
	}

	tenantID := tenants.TenantFromContext(r.Context())
	command := ports.NewCreateSubmissionCommand(
		tenantID,
		body.FormID,
		body.VersionID,
		domain.IdempotencyID(""),
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
		h.app.Logger.WarnContext(r.Context(), "context cancellation")
		return
	case res := <-resultChan:
		if res.err != nil {
			h.sendErrorResponse(w, res.err)
			return
		}

		httputil.SendJSONResponse(w, http.StatusAccepted, httputil.APIResponse[*dto.SubmissionResponse]{
			Message: "Accepted!",
			Data:    dto.SubmissionToResponse(res.data),
		})
	}
}

func (h *handlers) getSubmissionStatus(w http.ResponseWriter, r *http.Request) {
	tenantID := tenants.TenantFromContext(r.Context())
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
		h.app.Logger.WarnContext(r.Context(), "context cancellation")
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
	tenantID := tenants.TenantFromContext(r.Context())
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
		h.app.Logger.WarnContext(r.Context(), "context cancellation")
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
