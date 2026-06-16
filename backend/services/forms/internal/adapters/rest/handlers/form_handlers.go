package handlers

import (
	"net/http"
	"sundance/backend/pkg/auth"
	"sundance/backend/pkg/common/httputil"
	"sundance/backend/services/forms/internal/adapters/rest/dto"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
	"sundance/backend/services/forms/internal/core/ports/commands"

	"github.com/go-chi/chi/v5"
)

// @summary		Get all forms
// @tags		Forms
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param 		X-Request-ID header string false "Client-supplied request trace ID (generated if absent)"
// @param 		X-Correlation-ID header string false "Client-supplied correlation ID for tracing"
// @param 		X-Request-Date header string false "Client-supplied request date in ISO 8601 format" Format(date)
// @success		200 {array} dto.FormResponse
// @failure		500 {object} httputil.APIErrorResponse
// @security 	BearerAuth
// @Router		/forms [get]
func (h *Handlers) GetForms(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
	query := ports.NewFindFormsQuery(tenantID)
	resultChan := make(chan result[[]*domain.Form], 1)

	go func() {
		defer close(resultChan)
		forms, err := h.app.API.Forms.Find(r.Context(), query)
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
// @param 		X-Request-ID header string false "Client-supplied request trace ID (generated if absent)"
// @param 		X-Correlation-ID header string false "Client-supplied correlation ID for tracing"
// @param 		X-Request-Date header string false "Client-supplied request date in ISO 8601 format" Format(date)
// @param		formId path string true "Form ID"
// @success		200 {object} dto.FormResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @security 	BearerAuth
// @Router		/forms/{formId} [get]
func (h *Handlers) GetForm(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
	formID := h.getFormIDPathValue(r)
	query := ports.NewFindByIDQuery(tenantID, formID)
	resultChan := make(chan result[*domain.Form], 1)

	go func() {
		defer close(resultChan)
		form, err := h.app.API.Forms.FindByID(r.Context(), query)
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
// @param 		X-Request-ID header string false "Client-supplied request trace ID (generated if absent)"
// @param 		X-Correlation-ID header string false "Client-supplied correlation ID for tracing"
// @param 		X-Request-Date header string false "Client-supplied request date in ISO 8601 format" Format(date)
// @param		body body dto.UpsertFormRequest true "Create Form"
// @success		201 {object} httputil.APIResponse[dto.FormResponse]
// @failure		400 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @security 	BearerAuth
// @Router		/forms [post]
func (h *Handlers) CreateForm(w http.ResponseWriter, r *http.Request) {
	var body dto.UpsertFormRequest
	if err := httputil.ReadValidateJSONPayload(r, &body); err != nil {
		h.app.Logger.WarnContext(r.Context(), "invalid request body", "error", err)
		h.sendErrorResponse(w, err)
		return
	}

	tenantID := httputil.TenantFromContext(r.Context())
	resultChan := make(chan result[*domain.Form], 1)
	command := commands.NewCreateFormCommand(tenantID, body.Name, body.Description)

	go func() {
		defer close(resultChan)
		form, err := h.app.API.Forms.Create(r.Context(), command)
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
// @param 		X-Request-ID header string false "Client-supplied request trace ID (generated if absent)"
// @param 		X-Correlation-ID header string false "Client-supplied correlation ID for tracing"
// @param 		X-Request-Date header string false "Client-supplied request date in ISO 8601 format" Format(date)
// @param		formId path string true "Form ID"
// @param		body body dto.UpsertFormRequest true "Update Form"
// @success		200 {object} httputil.APIResponse[dto.FormResponse]
// @failure		400 {object} httputil.APIErrorResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @security 	BearerAuth
// @Router		/forms/{formId} [put]
func (h *Handlers) UpdateForm(w http.ResponseWriter, r *http.Request) {
	var body dto.UpsertFormRequest
	if err := httputil.ReadValidateJSONPayload(r, &body); err != nil {
		h.app.Logger.WarnContext(r.Context(), "invalid request body", "error", err)
		h.sendErrorResponse(w, err)
		return
	}

	tenantID := httputil.TenantFromContext(r.Context())
	formID := h.getFormIDPathValue(r)
	resultChan := make(chan result[*domain.Form], 1)
	command := commands.NewUpdateFormCommand(tenantID, formID, body.Name, body.Description)

	go func() {
		defer close(resultChan)
		form, err := h.app.API.Forms.Update(r.Context(), command)
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
// @param 		X-Request-ID header string false "Client-supplied request trace ID (generated if absent)"
// @param 		X-Correlation-ID header string false "Client-supplied correlation ID for tracing"
// @param 		X-Request-Date header string false "Client-supplied request date in ISO 8601 format" Format(date)
// @param		formId path string true "Form ID"
// @success		204
// @failure		400 {object} httputil.APIErrorResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @security 	BearerAuth
// @Router		/forms/{formId} [delete]
func (h *Handlers) DeleteForm(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
	formID := h.getFormIDPathValue(r)
	command := commands.NewDeleteCommand(tenantID, formID)
	resultChan := make(chan result[any], 1)

	go func() {
		defer close(resultChan)
		err := h.app.API.Forms.Delete(r.Context(), command)
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
// @tags		Forms
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param 		X-Request-ID header string false "Client-supplied request trace ID (generated if absent)"
// @param 		X-Correlation-ID header string false "Client-supplied correlation ID for tracing"
// @param 		X-Request-Date header string false "Client-supplied request date in ISO 8601 format" Format(date)
// @param		formId path string true "Form ID"
// @success		200 {array} dto.FormVersionResponse
// @failure		500 {object} httputil.APIErrorResponse
// @security 	BearerAuth
// @Router		/forms/{formId}/versions [get]
func (h *Handlers) GetFormVersions(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
	formID := h.getFormIDPathValue(r)
	query := ports.NewFindFormVersionsQuery(tenantID, formID)
	resultChan := make(chan result[[]*domain.FormVersion], 1)

	go func() {
		defer close(resultChan)
		versions, err := h.app.API.Forms.FindVersions(r.Context(), query)
		resultChan <- result[[]*domain.FormVersion]{versions, err}
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

		dtos := make([]*dto.FormVersionResponse, 0, len(res.data))
		for _, v := range res.data {
			dtos = append(dtos, dto.FormVersionToResponse(v))
		}

		httputil.SendJSONResponse(w, http.StatusOK, dtos)
	}
}

// @summary		Get a version by ID
// @tags		Forms
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param 		X-Request-ID header string false "Client-supplied request trace ID (generated if absent)"
// @param 		X-Correlation-ID header string false "Client-supplied correlation ID for tracing"
// @param 		X-Request-Date header string false "Client-supplied request date in ISO 8601 format" Format(date)
// @param		formId path string true "Form ID"
// @param		versionId path string true "Version ID"
// @success		200 {object} dto.FormVersionResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @security 	BearerAuth
// @Router		/forms/{formId}/versions/{versionId} [get]
func (h *Handlers) GetFormVersion(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
	formID := h.getFormIDPathValue(r)
	versionID := h.getVersionIDPathValue(r)
	query := ports.NewFindFormVersionByIDQuery(tenantID, formID, versionID)
	resultChan := make(chan result[*domain.FormVersion], 1)

	go func() {
		defer close(resultChan)
		version, err := h.app.API.Forms.FindVersion(r.Context(), query)
		resultChan <- result[*domain.FormVersion]{version, err}
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

		httputil.SendJSONResponse(w, http.StatusOK, dto.FormVersionToResponse(res.data))
	}
}

// @summary		Create a version
// @description	The pages field defines the structure of the form version including sections, fields, and validation rules.
// @tags		Forms
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param 		X-Request-ID header string false "Client-supplied request trace ID (generated if absent)"
// @param 		X-Correlation-ID header string false "Client-supplied correlation ID for tracing"
// @param 		X-Request-Date header string false "Client-supplied request date in ISO 8601 format" Format(date)
// @param		formId path string true "Form ID"
// @param		body body dto.UpsertFormVersionRequest true "Create Version"
// @success		201 {object} httputil.APIResponse[dto.FormVersionResponse]
// @failure		400 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @security 	BearerAuth
// @Router		/forms/{formId}/versions [post]
func (h *Handlers) CreateFormVersion(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
	formID := h.getFormIDPathValue(r)

	var body dto.UpsertFormVersionRequest
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

	resultChan := make(chan result[*domain.FormVersion], 1)
	command := commands.NewCreateFormVersionCommand(tenantID, formID, pages)

	go func() {
		defer close(resultChan)
		version, err := h.app.API.Forms.CreateVersion(r.Context(), command)
		resultChan <- result[*domain.FormVersion]{version, err}
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

		httputil.SendJSONResponse(w, http.StatusCreated, httputil.APIResponse[*dto.FormVersionResponse]{
			Message: "Successfully created!",
			Data:    dto.FormVersionToResponse(res.data),
		})
	}
}

// @summary		Update a version
// @description	Only draft versions can be updated. Published or retired versions are locked.
// @tags		Forms
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param 		X-Request-ID header string false "Client-supplied request trace ID (generated if absent)"
// @param 		X-Correlation-ID header string false "Client-supplied correlation ID for tracing"
// @param 		X-Request-Date header string false "Client-supplied request date in ISO 8601 format" Format(date)
// @param		formId path string true "Form ID"
// @param		versionId path string true "Version ID"
// @param		body body dto.UpsertFormVersionRequest true "Update Version"
// @success		200 {object} httputil.APIResponse[dto.FormVersionResponse]
// @failure		400 {object} httputil.APIErrorResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @security 	BearerAuth
// @Router		/forms/{formId}/versions/{versionId} [put]
func (h *Handlers) UpdateFormVersion(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
	formID := h.getFormIDPathValue(r)
	versionID := h.getVersionIDPathValue(r)

	var body dto.UpsertFormVersionRequest
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

	resultChan := make(chan result[*domain.FormVersion], 1)
	command := commands.NewUpdateFormVersionCommand(tenantID, versionID, formID, pages)

	go func() {
		defer close(resultChan)
		version, err := h.app.API.Forms.UpdateVersion(r.Context(), command)
		resultChan <- result[*domain.FormVersion]{version, err}
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

		httputil.SendJSONResponse(w, http.StatusOK, httputil.APIResponse[*dto.FormVersionResponse]{
			Message: "Successfully updated!",
			Data:    dto.FormVersionToResponse(res.data),
		})
	}
}

// @summary		Publish a version
// @description	Transitions a draft version to published status. Only one version per form can be published at a time.
// @tags		Forms
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param 		X-Request-ID header string false "Client-supplied request trace ID (generated if absent)"
// @param 		X-Correlation-ID header string false "Client-supplied correlation ID for tracing"
// @param 		X-Request-Date header string false "Client-supplied request date in ISO 8601 format" Format(date)
// @param		formId path string true "Form ID"
// @param		versionId path string true "Version ID"
// @success		200 {object} httputil.APIResponse[dto.FormVersionResponse]
// @failure		400 {object} httputil.APIErrorResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @security 	BearerAuth
// @Router		/forms/{formId}/versions/{versionId}/publish [post]
func (h *Handlers) PublishFormVersion(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
	formID := h.getFormIDPathValue(r)
	versionID := h.getVersionIDPathValue(r)
	claims := auth.GetClaimsFromContext(r.Context())
	sub, _ := claims.GetSubject()
	resultChan := make(chan result[*domain.FormVersion], 1)
	command := commands.NewPublishFormVersionCommand(tenantID, formID, versionID, sub)

	go func() {
		defer close(resultChan)
		version, err := h.app.API.Forms.PublishVersion(r.Context(), command)
		resultChan <- result[*domain.FormVersion]{version, err}
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
			Data:    dto.FormVersionToResponse(res.data),
		})
	}
}

// @summary		Retire a version
// @description	Transitions a published version to retired status, making it no longer active.
// @tags		Forms
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param 		X-Request-ID header string false "Client-supplied request trace ID (generated if absent)"
// @param 		X-Correlation-ID header string false "Client-supplied correlation ID for tracing"
// @param 		X-Request-Date header string false "Client-supplied request date in ISO 8601 format" Format(date)
// @param		formId path string true "Form ID"
// @param		versionId path string true "Version ID"
// @success		200 {object} httputil.APIResponse[dto.FormVersionResponse]
// @failure		400 {object} httputil.APIErrorResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @security 	BearerAuth
// @Router		/forms/{formId}/versions/{versionId}/retire [post]
func (h *Handlers) RetireFormVersion(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
	formID := h.getFormIDPathValue(r)
	versionID := h.getVersionIDPathValue(r)
	resultChan := make(chan result[*domain.FormVersion], 1)
	claims := auth.GetClaimsFromContext(r.Context())
	sub, _ := claims.GetSubject()
	command := commands.NewRetireFormVersionCommand(tenantID, formID, versionID, sub)

	go func() {
		defer close(resultChan)
		version, err := h.app.API.Forms.RetireVersion(r.Context(), command)
		resultChan <- result[*domain.FormVersion]{version, err}
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
			Data:    dto.FormVersionToResponse(res.data),
		})
	}
}

func (h *Handlers) getFormIDPathValue(r *http.Request) domain.FormID {
	id := chi.URLParam(r, "formId")
	return domain.FormID(id)
}

func (h *Handlers) getVersionIDPathValue(r *http.Request) domain.FormVersionID {
	id := chi.URLParam(r, "versionId")
	return domain.FormVersionID(id)
}
