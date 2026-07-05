package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"sundance/backend/pkg/common/httputil"
	"sundance/backend/services/forms/internal/adapters/rest/dto"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
	"sundance/backend/services/forms/internal/core/ports/commands"

	"github.com/go-chi/chi/v5"
)

var (
	errSubmissionStatus        = errors.New("submission status")
	errSubmissionStatusPending = fmt.Errorf("%w; submission pending", errSubmissionStatus)
)

// @summary		Get all Submissions
// @tags		Submissions
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param 		X-Request-ID header string false "Client-supplied request trace ID (generated if absent)"
// @param 		X-Correlation-ID header string false "Client-supplied correlation ID for tracing"
// @param 		X-Request-Date header string false "Client-supplied request date in ISO 8601 format" Format(date)
// @success		200 {array} dto.SubmissionResponse
// @failure		500 {object} httputil.APIErrorResponse
// @security 	BearerAuth
// @Router		/submissions [get]
func (h *Handlers) GetSubmissions(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
	query := ports.NewFindSubmissionsQuery(tenantID)
	resultChan := make(chan result[[]*domain.Submission], 1)

	go func() {
		defer close(resultChan)
		submissions, err := h.app.API.Submissions.Find(r.Context(), query)
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

// @summary		Get a submission by reference ID
// @tags		Submissions
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param 		X-Request-ID header string false "Client-supplied request trace ID (generated if absent)"
// @param 		X-Correlation-ID header string false "Client-supplied correlation ID for tracing"
// @param 		X-Request-Date header string false "Client-supplied request date in ISO 8601 format" Format(date)
// @param		referenceId path string true "Reference ID"
// @success		200 {object} dto.SubmissionResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @security 	BearerAuth
// @Router		/submissions/by-reference/{referenceId} [get]
func (h *Handlers) GetSubmissionByReferenceID(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
	referenceID := h.getReferenceIDPathValue(r)
	query := ports.NewFindSubmissionByIDQuery(tenantID, referenceID)
	resultChan := make(chan result[*domain.Submission], 1)

	go func() {
		defer close(resultChan)
		submission, err := h.app.API.Submissions.FindByReferenceID(r.Context(), query)
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

// @summary		Create a submission
// @description	Accepts a form submission for asynchronous processing. An Idempotency-Key header is required to prevent duplicate submissions.
// @tags		Submissions
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param 		X-Request-ID header string false "Client-supplied request trace ID (generated if absent)"
// @param 		X-Correlation-ID header string false "Client-supplied correlation ID for tracing"
// @param 		X-Request-Date header string false "Client-supplied request date in ISO 8601 format" Format(date)
// @param		Idempotency-Key header string true "Idempotency Key"
// @param		body body dto.SubmissionRequest true "Create Submission"
// @success		202 {object} httputil.APIResponse[dto.SubmissionResponse]
// @failure		400 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @security 	BearerAuth
// @Router		/submissions [post]
func (h *Handlers) CreateSubmission(w http.ResponseWriter, r *http.Request) {
	var body dto.SubmissionRequest
	if err := httputil.ReadValidateJSONPayload(r, &body); err != nil {
		h.app.Logger.WarnContext(r.Context(), "invalid request body", "error", err)
		h.sendErrorResponse(w, err)
		return
	}

	values := make([]*domain.SubmissionFieldValue, 0, len(body.Values))
	for _, fv := range body.Values {
		values = append(values, domain.NewSubmissionFieldValue(fv.FieldID, fv.Value, fv.CollectionIndex))
	}

	tenantID := httputil.TenantFromContext(r.Context())
	idempotencyID := httputil.IdempotencyFromContext(r.Context())
	command := commands.NewCreateSubmissionCommand(
		tenantID,
		domain.FormID(body.FormID),
		domain.FormVersionID(body.VersionID),
		domain.IdempotencyID(idempotencyID),
		values,
	)
	resultChan := make(chan result[*domain.Submission], 1)

	go func() {
		defer close(resultChan)
		submission, err := h.app.API.Submissions.Create(r.Context(), command)
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

// @summary		Get submission facts
// @description	Returns the canonical fact map for an accepted submission, keyed by tag paths.
// @description	A 400 is returned if the submission status is not `accepted` (e.g. pending, rejected, or failed).
// @tags		Submissions
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param 		X-Request-ID header string false "Client-supplied request trace ID (generated if absent)"
// @param 		X-Correlation-ID header string false "Client-supplied correlation ID for tracing"
// @param 		X-Request-Date header string false "Client-supplied request date in ISO 8601 format" Format(date)
// @param		referenceId path string true "Reference ID"
// @success		200 {object} object "Canonical fact map keyed by tag paths"
// @failure		400 {object} httputil.APIErrorResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @security 	BearerAuth
// @Router		/submissions/by-reference/{referenceId}/facts [get]
func (h *Handlers) GetSubmissionFacts(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
	referenceID := h.getReferenceIDPathValue(r)
	query := ports.NewFindSubmissionByIDQuery(tenantID, referenceID)
	resultChan := make(chan result[domain.FactMap], 1)

	go func() {
		defer close(resultChan)

		submission, err := h.app.API.Submissions.FindByReferenceID(r.Context(), query)
		if err != nil {
			resultChan <- result[domain.FactMap]{nil, err}
			return
		}

		switch submission.Status {
		case domain.SubmissionStatusPending:
			resultChan <- result[domain.FactMap]{nil, errSubmissionStatusPending}
			return
		case domain.SubmissionStatusRejected, domain.SubmissionStatusFailed:
			resultChan <- result[domain.FactMap]{nil, errSubmissionStatus}
			return
		}

		resultChan <- result[domain.FactMap]{submission.ToFactMap(), nil}
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

		httputil.SendJSONResponse(w, http.StatusOK, res.data)
	}
}

// @summary		Get a submission status
// @tags		Submissions
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param 		X-Request-ID header string false "Client-supplied request trace ID (generated if absent)"
// @param 		X-Correlation-ID header string false "Client-supplied correlation ID for tracing"
// @param 		X-Request-Date header string false "Client-supplied request date in ISO 8601 format" Format(date)
// @param		referenceId path string true "Reference ID"
// @success		200 {object} object{status=string}
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @security 	BearerAuth
// @Router		/submissions/by-reference/{referenceId}/status [get]
func (h *Handlers) GetSubmissionStatus(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
	referenceID := h.getReferenceIDPathValue(r)
	query := ports.NewFindSubmissionByIDQuery(tenantID, referenceID)
	resultChan := make(chan result[*domain.Submission], 1)

	go func() {
		defer close(resultChan)
		submission, err := h.app.API.Submissions.FindByReferenceID(r.Context(), query)
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

// @summary		Replay a submission
// @description	Re-publishes the submission event for reprocessing by downstream consumers.
// @tags		Submissions
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param 		X-Request-ID header string false "Client-supplied request trace ID (generated if absent)"
// @param 		X-Correlation-ID header string false "Client-supplied correlation ID for tracing"
// @param 		X-Request-Date header string false "Client-supplied request date in ISO 8601 format" Format(date)
// @param		submissionId path string true "Submission ID"
// @success		202 {object} httputil.APIResponse[dto.SubmissionResponse]
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @security 	BearerAuth
// @Router		/submissions/{submissionId}/replay [post]
func (h *Handlers) ReplaySubmission(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
	id := chi.URLParam(r, "submissionId")
	command := commands.NewReplaySubmissionCommand(
		tenantID,
		domain.SubmissionID(id),
	)
	resultChan := make(chan result[any])

	go func() {
		defer close(resultChan)
		err := h.app.API.Submissions.Replay(r.Context(), command)
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

		httputil.SendJSONResponse(w, http.StatusAccepted, httputil.APIResponse[*dto.SubmissionResponse]{
			Message: fmt.Sprintf("Successfully replayed submission %s", id),
		})
	}
}

func (h *Handlers) getReferenceIDPathValue(r *http.Request) domain.ReferenceID {
	id := chi.URLParam(r, "referenceId")
	return domain.ReferenceID(id)
}
