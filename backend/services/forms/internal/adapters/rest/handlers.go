package rest

import (
	"errors"
	"net/http"

	"github.com/cmclaughlin24/sundance/backend/pkg/auth"
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
	tenantID, err := h.getTenantFromContext(r)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	query := ports.NewFindFormsQuery(tenantID)
	resultChan := make(chan result[[]*domain.Form], 1)

	go func() {
		defer close(resultChan)
		forms, err := h.app.Services.Forms.Find(r.Context(), query)
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

		httputil.SendJSONResponse(w, http.StatusOK, dtos)
	}
}

func (h *handlers) getForm(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantFromContext(r)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

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
		return
	case res := <-resultChan:
		if res.err != nil {
			h.sendErrorResponse(w, res.err)
			return
		}

		httputil.SendJSONResponse(w, http.StatusOK, dto.FormToResponse(res.data))
	}
}

func (h *handlers) createForm(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantFromContext(r)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	var body dto.UpsertFormRequest
	if err := httputil.ReadValidateJSONPayload(r, &body); err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	resultChan := make(chan result[*domain.Form], 1)
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

		httputil.SendJSONResponse(w, http.StatusCreated, httputil.APIResponse[*dto.FormResponse]{
			Message: "Successfully created!",
			Data:    dto.FormToResponse(res.data),
		})
	}
}

func (h *handlers) updateForm(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantFromContext(r)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	formID := h.getFormIDPathValue(r)

	var body dto.UpsertFormRequest
	if err := httputil.ReadValidateJSONPayload(r, &body); err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	resultChan := make(chan result[*domain.Form], 1)
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

		httputil.SendJSONResponse(w, http.StatusOK, httputil.APIResponse[*dto.FormResponse]{
			Message: "Successfully updated!",
			Data:    dto.FormToResponse(res.data),
		})
	}
}

func (h *handlers) deleteForm(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantFromContext(r)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

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
		return
	case res := <-resultChan:
		if res.err != nil {
			h.sendErrorResponse(w, res.err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *handlers) getVersions(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantFromContext(r)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

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

func (h *handlers) getVersion(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantFromContext(r)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

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
		return
	case res := <-resultChan:
		if res.err != nil {
			h.sendErrorResponse(w, res.err)
			return
		}

		httputil.SendJSONResponse(w, http.StatusOK, dto.VersionToResponse(res.data))
	}
}

func (h *handlers) createVersion(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantFromContext(r)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	formID := h.getFormIDPathValue(r)

	var body dto.UpsertVersionRequest
	if err := httputil.ReadValidateJSONPayload(r, &body); err != nil {
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

func (h *handlers) updateVersion(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantFromContext(r)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	formID := h.getFormIDPathValue(r)
	versionID := h.getVersionIDPathValue(r)

	var body dto.UpsertVersionRequest
	if err := httputil.ReadValidateJSONPayload(r, &body); err != nil {
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

func (h *handlers) publishVersion(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantFromContext(r)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

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

func (h *handlers) retireVersion(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantFromContext(r)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

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

func (h *handlers) getTenantFromContext(r *http.Request) (string, error) {
	tenantID, err := tenants.TenantFromContext(r.Context())

	if err != nil {
		return "", err
	}

	return tenantID, nil
}

func (h *handlers) getFormIDPathValue(r *http.Request) domain.FormID {
	id := chi.URLParam(r, "formId")
	return domain.FormID(id)
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
