package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cmclaughlin24/sundance/backend/pkg/auth"
	"github.com/cmclaughlin24/sundance/backend/pkg/common/httputil"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/adapters/rest/dto"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/ports"
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

// @summary		Get all forms
// @tags		Forms
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @success		200 {array} dto.FormResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/forms [get]
func (h *handlers) getForms(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
	query := ports.NewFindFormsQuery(tenantID)
	resultChan := make(chan result[[]*domain.Form], 1)

	go func() {
		defer close(resultChan)
		forms, err := h.app.Services.Forms.Find(r.Context(), query)
		resultChan <- result[[]*domain.Form]{forms, err}
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

		dtos := make([]*dto.FormResponse, 0, len(res.data))
		for _, form := range res.data {
			dtos = append(dtos, dto.FormToResponse(form))
		}

		httputil.SendJSONResponse(w, http.StatusOK, dtos)
	}
}

// @summary		Get a form by ID
// @tags		Forms
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param		formId path string true "Form ID"
// @success		200 {object} dto.FormResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/forms/{formId} [get]
func (h *handlers) getForm(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
	formID := h.getFormIDPathValue(r)
	query := ports.NewFindFormsByIDQuery(tenantID, formID)
	resultChan := make(chan result[*domain.Form], 1)

	go func() {
		defer close(resultChan)
		form, err := h.app.Services.Forms.FindByID(r.Context(), query)
		resultChan <- result[*domain.Form]{form, err}
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

		httputil.SendJSONResponse(w, http.StatusOK, dto.FormToResponse(res.data))
	}
}

// @summary		Create a form
// @tags		Forms
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param		body body dto.UpsertFormRequest true "Create Form"
// @success		201 {object} httputil.APIResponse[dto.FormResponse]
// @failure		400 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/forms [post]
func (h *handlers) createForm(w http.ResponseWriter, r *http.Request) {
	var body dto.UpsertFormRequest
	if err := httputil.ReadValidateJSONPayload(r, &body); err != nil {
		h.app.Logger.WarnContext(r.Context(), "invalid request body", "error", err)
		h.sendErrorResponse(w, err)
		return
	}

	tenantID := httputil.TenantFromContext(r.Context())
	resultChan := make(chan result[*domain.Form], 1)
	command := ports.NewCreateFormCommand(tenantID, body.Name, body.Description)

	go func() {
		defer close(resultChan)
		form, err := h.app.Services.Forms.Create(r.Context(), command)
		resultChan <- result[*domain.Form]{form, err}
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

		httputil.SendJSONResponse(w, http.StatusCreated, httputil.APIResponse[*dto.FormResponse]{
			Message: "Successfully created!",
			Data:    dto.FormToResponse(res.data),
		})
	}
}

// @summary		Update a form
// @tags		Forms
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param		formId path string true "Form ID"
// @param		body body dto.UpsertFormRequest true "Update Form"
// @success		200 {object} httputil.APIResponse[dto.FormResponse]
// @failure		400 {object} httputil.APIErrorResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/forms/{formId} [put]
func (h *handlers) updateForm(w http.ResponseWriter, r *http.Request) {
	var body dto.UpsertFormRequest
	if err := httputil.ReadValidateJSONPayload(r, &body); err != nil {
		h.app.Logger.WarnContext(r.Context(), "invalid request body", "error", err)
		h.sendErrorResponse(w, err)
		return
	}

	tenantID := httputil.TenantFromContext(r.Context())
	formID := h.getFormIDPathValue(r)
	resultChan := make(chan result[*domain.Form], 1)
	command := ports.NewUpdateFormCommand(tenantID, formID, body.Name, body.Description)

	go func() {
		defer close(resultChan)
		form, err := h.app.Services.Forms.Update(r.Context(), command)
		resultChan <- result[*domain.Form]{form, err}
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

		httputil.SendJSONResponse(w, http.StatusOK, httputil.APIResponse[*dto.FormResponse]{
			Message: "Successfully updated!",
			Data:    dto.FormToResponse(res.data),
		})
	}
}

// @summary		Delete a form
// @description	All versions belonging to the form will also be deleted.
// @tags		Forms
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param		formId path string true "Form ID"
// @success		204
// @failure		400 {object} httputil.APIErrorResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/forms/{formId} [delete]
func (h *handlers) deleteForm(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
	formID := h.getFormIDPathValue(r)
	command := ports.NewRemoveFormCommand(tenantID, formID)
	resultChan := make(chan result[any], 1)

	go func() {
		defer close(resultChan)
		err := h.app.Services.Forms.Delete(r.Context(), command)
		resultChan <- result[any]{nil, err}
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

		w.WriteHeader(http.StatusNoContent)
	}
}

// @summary		Get all versions
// @tags		Versions
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param		formId path string true "Form ID"
// @success		200 {array} dto.VersionResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/forms/{formId}/versions [get]
func (h *handlers) getVersions(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
	formID := h.getFormIDPathValue(r)
	query := ports.NewFindVersionsQuery(tenantID, formID)
	resultChan := make(chan result[[]*domain.Version], 1)

	go func() {
		defer close(resultChan)
		versions, err := h.app.Services.Forms.FindVersions(r.Context(), query)
		resultChan <- result[[]*domain.Version]{versions, err}
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

		dtos := make([]*dto.VersionResponse, 0, len(res.data))
		for _, v := range res.data {
			dtos = append(dtos, dto.VersionToResponse(v))
		}

		httputil.SendJSONResponse(w, http.StatusOK, dtos)
	}
}

// @summary		Get a version by ID
// @tags		Versions
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param		formId path string true "Form ID"
// @param		versionId path string true "Version ID"
// @success		200 {object} dto.VersionResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/forms/{formId}/versions/{versionId} [get]
func (h *handlers) getVersion(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
	formID := h.getFormIDPathValue(r)
	versionID := h.getVersionIDPathValue(r)
	query := ports.NewFindVersionByIDQuery(tenantID, formID, versionID)
	resultChan := make(chan result[*domain.Version], 1)

	go func() {
		defer close(resultChan)
		version, err := h.app.Services.Forms.FindVersion(r.Context(), query)
		resultChan <- result[*domain.Version]{version, err}
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

		httputil.SendJSONResponse(w, http.StatusOK, dto.VersionToResponse(res.data))
	}
}

// @summary		Create a version
// @description	The pages field defines the structure of the form version including sections, fields, and validation rules.
// @tags		Versions
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param		formId path string true "Form ID"
// @param		body body dto.UpsertVersionRequest true "Create Version"
// @success		201 {object} httputil.APIResponse[dto.VersionResponse]
// @failure		400 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/forms/{formId}/versions [post]
func (h *handlers) createVersion(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
	formID := h.getFormIDPathValue(r)

	var body dto.UpsertVersionRequest
	if err := httputil.ReadValidateJSONPayload(r, &body); err != nil {
		h.app.Logger.WarnContext(r.Context(), "invalid request body", "error", err)
		h.sendErrorResponse(w, err)
		return
	}

	pages, err := dto.RequestToPages(body)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	resultChan := make(chan result[*domain.Version], 1)
	command := ports.NewCreateVersionCommand(tenantID, formID, pages)

	go func() {
		defer close(resultChan)
		version, err := h.app.Services.Forms.CreateVersion(r.Context(), command)
		resultChan <- result[*domain.Version]{version, err}
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

		httputil.SendJSONResponse(w, http.StatusCreated, httputil.APIResponse[*dto.VersionResponse]{
			Message: "Successfully created!",
			Data:    dto.VersionToResponse(res.data),
		})
	}
}

// @summary		Update a version
// @description	Only draft versions can be updated. Published or retired versions are locked.
// @tags		Versions
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param		formId path string true "Form ID"
// @param		versionId path string true "Version ID"
// @param		body body dto.UpsertVersionRequest true "Update Version"
// @success		200 {object} httputil.APIResponse[dto.VersionResponse]
// @failure		400 {object} httputil.APIErrorResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/forms/{formId}/versions/{versionId} [put]
func (h *handlers) updateVersion(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
	formID := h.getFormIDPathValue(r)
	versionID := h.getVersionIDPathValue(r)

	var body dto.UpsertVersionRequest
	if err := httputil.ReadValidateJSONPayload(r, &body); err != nil {
		h.app.Logger.WarnContext(r.Context(), "invalid request body", "error", err)
		h.sendErrorResponse(w, err)
		return
	}

	pages, err := dto.RequestToPages(body)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	resultChan := make(chan result[*domain.Version], 1)
	command := ports.NewUpdateVersionCommand(tenantID, versionID, formID, pages)

	go func() {
		defer close(resultChan)
		version, err := h.app.Services.Forms.UpdateVersion(r.Context(), command)
		resultChan <- result[*domain.Version]{version, err}
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

		httputil.SendJSONResponse(w, http.StatusOK, httputil.APIResponse[*dto.VersionResponse]{
			Message: "Successfully updated!",
			Data:    dto.VersionToResponse(res.data),
		})
	}
}

// @summary		Publish a version
// @description	Transitions a draft version to published status. Only one version per form can be published at a time.
// @tags		Versions
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param		formId path string true "Form ID"
// @param		versionId path string true "Version ID"
// @success		200 {object} httputil.APIResponse[dto.VersionResponse]
// @failure		400 {object} httputil.APIErrorResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/forms/{formId}/versions/{versionId}/publish [post]
func (h *handlers) publishVersion(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
	formID := h.getFormIDPathValue(r)
	versionID := h.getVersionIDPathValue(r)
	claims := auth.GetClaimsFromContext(r.Context())
	resultChan := make(chan result[*domain.Version], 1)
	command := ports.NewPublishVersionCommand(tenantID, formID, versionID, claims.GetSubject())

	go func() {
		defer close(resultChan)
		version, err := h.app.Services.Forms.PublishVersion(r.Context(), command)
		resultChan <- result[*domain.Version]{version, err}
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

		httputil.SendJSONResponse(w, http.StatusOK, httputil.APIResponse[any]{
			Message: "Successfully published!",
			Data:    dto.VersionToResponse(res.data),
		})
	}
}

// @summary		Retire a version
// @description	Transitions a published version to retired status, making it no longer active.
// @tags		Versions
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param		formId path string true "Form ID"
// @param		versionId path string true "Version ID"
// @success		200 {object} httputil.APIResponse[dto.VersionResponse]
// @failure		400 {object} httputil.APIErrorResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/forms/{formId}/versions/{versionId}/retire [post]
func (h *handlers) retireVersion(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
	formID := h.getFormIDPathValue(r)
	versionID := h.getVersionIDPathValue(r)
	resultChan := make(chan result[*domain.Version], 1)
	claims := auth.GetClaimsFromContext(r.Context())
	command := ports.NewRetireVersionCommand(tenantID, formID, versionID, claims.GetSubject())

	go func() {
		defer close(resultChan)
		version, err := h.app.Services.Forms.RetireVersion(r.Context(), command)
		resultChan <- result[*domain.Version]{version, err}
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

		httputil.SendJSONResponse(w, http.StatusOK, httputil.APIResponse[any]{
			Message: "Successfully retired!",
			Data:    dto.VersionToResponse(res.data),
		})
	}
}

// @summary		Get all Submissions
// @tags		Submissions
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @success		200 {array} dto.SubmissionResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/submissions [get]
func (h *handlers) getSubmissions(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
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

// @summary		Get a submission by reference ID
// @tags		Submissions
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param		referenceId path string true "Reference ID"
// @success		200 {object} dto.SubmissionResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/submissions/by-reference/{referenceId} [get]
func (h *handlers) getSubmissionByReferenceID(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
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

// @summary		Create a submission
// @description	Accepts a form submission for asynchronous processing. An Idempotency-Key header is required to prevent duplicate submissions.
// @tags		Submissions
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param		Idempotency-Key header string true "Idempotency Key"
// @param		body body dto.SubmissionRequest true "Create Submission"
// @success		202 {object} httputil.APIResponse[dto.SubmissionResponse]
// @failure		400 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/submissions [post]
func (h *handlers) createSubmission(w http.ResponseWriter, r *http.Request) {
	var body dto.SubmissionRequest
	if err := httputil.ReadValidateJSONPayload(r, &body); err != nil {
		h.app.Logger.WarnContext(r.Context(), "invalid request body", "error", err)
		h.sendErrorResponse(w, err)
		return
	}

	tenantID := httputil.TenantFromContext(r.Context())
	idempotencyID := httputil.IdempotencyFromContext(r.Context())
	command := ports.NewCreateSubmissionCommand(
		tenantID,
		body.FormID,
		body.VersionID,
		domain.IdempotencyID(idempotencyID),
		body.Payload,
	)
	resultChan := make(chan result[*domain.Submission], 1)

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

// @summary		Get a submission status
// @tags		Submissions
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param		referenceId path string true "Reference ID"
// @success		200 {object} object{status=string}
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/submissions/by-reference/{referenceId}/status [get]
func (h *handlers) getSubmissionStatus(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
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

// @summary		Replay a submission
// @description	Re-publishes the submission event for reprocessing by downstream consumers.
// @tags		Submissions
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param		submissionId path string true "Submission ID"
// @success		202 {object} httputil.APIResponse[dto.SubmissionResponse]
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/submissions/{submissionId}/replay [post]
func (h *handlers) replaySubmission(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
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

		httputil.SendJSONResponse(w, http.StatusAccepted, httputil.APIResponse[*dto.SubmissionResponse]{
			Message: fmt.Sprintf("Successfully replayed submission %s", id),
		})
	}
}

func (h *handlers) getFormIDPathValue(r *http.Request) domain.FormID {
	id := chi.URLParam(r, "formId")
	return domain.FormID(id)
}

func (h *handlers) getReferenceIDPathValue(r *http.Request) domain.ReferenceID {
	id := chi.URLParam(r, "referenceId")
	return domain.ReferenceID(id)
}

func (h *handlers) getVersionIDPathValue(r *http.Request) domain.VersionID {
	id := chi.URLParam(r, "versionId")
	return domain.VersionID(id)
}

func (h *handlers) sendErrorResponse(w http.ResponseWriter, err error) {
	switch {
	case isBadRequest(err):
		httputil.SendJSONResponse(w, http.StatusBadRequest, httputil.APIErrorResponse{
			Message:    "Bad Request",
			Error:      err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	default:
		httputil.SendErrorResponse(w, err)
	}
}

func isBadRequest(err error) bool {
	return errors.Is(err, domain.ErrVersionLocked) ||
		errors.Is(err, domain.ErrInvalidVersion) ||
		errors.Is(err, domain.ErrInvalidVersionStatus) ||
		errors.Is(err, domain.ErrDuplicateVersion) ||
		errors.Is(err, domain.ErrInvalidPosition) ||
		errors.Is(err, domain.ErrDuplicatePosition) ||
		errors.Is(err, domain.ErrInvalidRuleType) ||
		errors.Is(err, domain.ErrDuplicateRuleType) ||
		errors.Is(err, domain.ErrPublishedByRequired) ||
		errors.Is(err, domain.ErrRetiredByRequired) ||
		errors.Is(err, domain.ErrInvalidFieldType) ||
		errors.Is(err, domain.ErrInvalidFieldAttributes) ||
		errors.Is(err, domain.ErrInvalidForm) ||
		errors.Is(err, domain.ErrFormHasActiveVersion) ||
		errors.Is(err, domain.ErrInvalidPage) ||
		errors.Is(err, domain.ErrInvalidSection)
}
