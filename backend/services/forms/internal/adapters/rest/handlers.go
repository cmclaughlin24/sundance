package rest

import (
	"net/http"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/httputil"
	"github.com/cmclaughlin24/sundance/backend/pkg/common/tenants"
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

func (h *handlers) getForms(w http.ResponseWriter, r *http.Request) {
	resultChan := make(chan result[[]*domain.Form], 1)

	go func() {
		defer close(resultChan)
		forms, err := h.app.Services.Forms.Find(r.Context())
		resultChan <- result[[]*domain.Form]{forms, err}
	}()

	select {
	case <-r.Context().Done():
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

		httputil.SendJsonResponse(w, http.StatusOK, dtos)
	}
}

func (h *handlers) getForm(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantFromContext(r)
	if err != nil {
		httputil.SendErrorResponse(w, err)
		return
	}

	formID := h.getFormIdPathValue(r)
	query := ports.NewFindByIDQuery(tenantID, formID)
	resultChan := make(chan result[*domain.Form], 1)

	go func() {
		defer close(resultChan)
		form, err := h.app.Services.Forms.FindById(r.Context(), query)
		resultChan <- result[*domain.Form]{form, err}
	}()

	select {
	case <-r.Context().Done():
		return
	case res := <-resultChan:
		if res.err != nil {
			h.sendErrorResponse(w, res.err)
			return
		}

		httputil.SendJsonResponse(w, http.StatusOK, dto.FormToResponse(res.data))
	}
}

func (h *handlers) createForm(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantFromContext(r)
	if err != nil {
		httputil.SendErrorResponse(w, err)
		return
	}

	resultChan := make(chan result[*domain.Form], 1)

	var body dto.UpsertFormRequest
	if err := httputil.ReadJsonPayload(r, &body); err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	command := ports.NewCreateFormCommand(tenantID, body.Name, body.Description)

	go func() {
		defer close(resultChan)
		form, err := h.app.Services.Forms.Create(r.Context(), command)
		resultChan <- result[*domain.Form]{form, err}
	}()

	select {
	case <-r.Context().Done():
		return
	case res := <-resultChan:
		if res.err != nil {
			h.sendErrorResponse(w, res.err)
			return
		}

		httputil.SendJsonResponse(w, http.StatusCreated, httputil.ApiResponse[*dto.FormResponse]{
			Message: "Successfully created!",
			Data:    dto.FormToResponse(res.data),
		})
	}
}

func (h *handlers) updateForm(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantFromContext(r)
	if err != nil {
		httputil.SendErrorResponse(w, err)
		return
	}

	formID := h.getFormIdPathValue(r)
	resultChan := make(chan result[*domain.Form], 1)

	var body dto.UpsertFormRequest
	if err := httputil.ReadJsonPayload(r, &body); err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	command := ports.NewUpdateFormCommand(tenantID, formID, body.Name, body.Description)

	go func() {
		defer close(resultChan)
		form, err := h.app.Services.Forms.Update(r.Context(), command)
		resultChan <- result[*domain.Form]{form, err}
	}()

	select {
	case <-r.Context().Done():
		return
	case res := <-resultChan:
		if res.err != nil {
			h.sendErrorResponse(w, res.err)
			return
		}

		httputil.SendJsonResponse(w, http.StatusOK, httputil.ApiResponse[*dto.FormResponse]{
			Message: "Successfully updated!",
			Data:    dto.FormToResponse(res.data),
		})
	}
}

func (h *handlers) getVersions(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantFromContext(r)
	if err != nil {
		httputil.SendErrorResponse(w, err)
		return
	}

	formID := h.getFormIdPathValue(r)
	query := ports.NewFindVersionsQuery(tenantID, formID)
	resultChan := make(chan result[[]*domain.Version], 1)

	go func() {
		defer close(resultChan)
		versions, err := h.app.Services.Forms.FindVersions(r.Context(), query)
		resultChan <- result[[]*domain.Version]{versions, err}
	}()

	select {
	case <-r.Context().Done():
		return
	case res := <-resultChan:
		if res.err != nil {
			h.sendErrorResponse(w, res.err)
			return
		}

		dtos := make([]*dto.VersionResponseDto, 0, len(res.data))
		for _, v := range res.data {
			dtos = append(dtos, dto.VersionToResponse(v))
		}

		httputil.SendJsonResponse(w, http.StatusOK, dtos)
	}
}

func (h *handlers) getVersion(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantFromContext(r)
	if err != nil {
		httputil.SendErrorResponse(w, err)
		return
	}

	formID := h.getFormIdPathValue(r)
	versionID := h.getVersionIdPathValue(r)
	query := ports.NewFindVersionByIDQuery(tenantID, formID, versionID)
	resultChan := make(chan result[*domain.Version], 1)

	go func() {
		defer close(resultChan)
		version, err := h.app.Services.Forms.FindVersion(r.Context(), query)
		resultChan <- result[*domain.Version]{version, err}
	}()

	select {
	case <-r.Context().Done():
		return
	case res := <-resultChan:
		if res.err != nil {
			h.sendErrorResponse(w, res.err)
			return
		}

		httputil.SendJsonResponse(w, http.StatusOK, dto.VersionToResponse(res.data))
	}
}

func (h *handlers) createVersion(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantFromContext(r)
	if err != nil {
		httputil.SendErrorResponse(w, err)
		return
	}

	formID := h.getFormIdPathValue(r)
	resultChan := make(chan result[*domain.Version], 1)

	var body dto.CreateVersionDto
	if err := httputil.ReadJsonPayload(r, &body); err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	command := ports.NewCreateVersionCommand(tenantID, formID)

	go func() {
		defer close(resultChan)
		version, err := h.app.Services.Forms.CreateVersion(r.Context(), command)
		resultChan <- result[*domain.Version]{version, err}
	}()

	select {
	case <-r.Context().Done():
		return
	case res := <-resultChan:
		if res.err != nil {
			h.sendErrorResponse(w, res.err)
			return
		}

		httputil.SendJsonResponse(w, http.StatusCreated, httputil.ApiResponse[*dto.VersionResponseDto]{
			Message: "Successfully created!",
			Data:    dto.VersionToResponse(res.data),
		})
	}
}

func (h *handlers) updateVersion(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantFromContext(r)
	if err != nil {
		httputil.SendErrorResponse(w, err)
		return
	}

	formID := h.getFormIdPathValue(r)
	versionID := h.getVersionIdPathValue(r)
	resultChan := make(chan result[*domain.Version], 1)

	var body dto.UpdateVersionDto
	if err := httputil.ReadJsonPayload(r, &body); err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	pages, err := dto.RequestToPages(body)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	command := ports.NewUpdateVersionCommand(tenantID, versionID, formID, pages)

	go func() {
		defer close(resultChan)
		version, err := h.app.Services.Forms.UpdateVersion(r.Context(), command)
		resultChan <- result[*domain.Version]{version, err}
	}()

	select {
	case <-r.Context().Done():
		return
	case res := <-resultChan:
		if res.err != nil {
			h.sendErrorResponse(w, res.err)
			return
		}

		httputil.SendJsonResponse(w, http.StatusOK, httputil.ApiResponse[*dto.VersionResponseDto]{
			Message: "Successfully updated!",
			Data:    dto.VersionToResponse(res.data),
		})
	}
}

func (h *handlers) publishVersion(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantFromContext(r)
	if err != nil {
		httputil.SendErrorResponse(w, err)
		return
	}

	formID := h.getFormIdPathValue(r)
	versionID := h.getVersionIdPathValue(r)
	resultChan := make(chan result[*domain.Version], 1)
	// FIXME: Remove temporary placeholder for user ID.
	command := ports.NewPublishVersionCommand(tenantID, formID, versionID, "placeholder")

	go func() {
		defer close(resultChan)
		version, err := h.app.Services.Forms.PublishVersion(r.Context(), command)
		resultChan <- result[*domain.Version]{version, err}
	}()

	select {
	case <-r.Context().Done():
		return
	case res := <-resultChan:
		if res.err != nil {
			h.sendErrorResponse(w, res.err)
			return
		}

		httputil.SendJsonResponse(w, http.StatusOK, httputil.ApiResponse[any]{
			Message: "Successfully published!",
			Data:    nil,
		})
	}
}

func (h *handlers) retireVersion(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantFromContext(r)
	if err != nil {
		httputil.SendErrorResponse(w, err)
		return
	}

	formId := h.getFormIdPathValue(r)
	versionId := h.getVersionIdPathValue(r)
	resultChan := make(chan result[*domain.Version], 1)
	// FIXME: Remove temporary placeholder for user ID.
	command := ports.NewRetireVersionCommand(tenantID, formId, versionId, "placeholder")

	go func() {
		defer close(resultChan)
		version, err := h.app.Services.Forms.RetireVersion(r.Context(), command)
		resultChan <- result[*domain.Version]{version, err}
	}()

	select {
	case <-r.Context().Done():
		return
	case res := <-resultChan:
		if res.err != nil {
			h.sendErrorResponse(w, res.err)
			return
		}

		httputil.SendJsonResponse(w, http.StatusOK, httputil.ApiResponse[any]{
			Message: "Successfully retired!",
			Data:    nil,
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

func (h *handlers) getFormIdPathValue(r *http.Request) domain.FormID {
	id := chi.URLParam(r, "formId")
	return domain.FormID(id)
}

func (h *handlers) getVersionIdPathValue(r *http.Request) domain.VersionID {
	id := chi.URLParam(r, "versionId")
	return domain.VersionID(id)
}

func (h *handlers) sendErrorResponse(w http.ResponseWriter, err error) {
	switch {
	default:
		httputil.SendErrorResponse(w, err)
	}
}
