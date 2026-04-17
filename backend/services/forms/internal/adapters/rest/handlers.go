package rest

import (
	"net/http"

	"github.com/cmclaughlin24/sundance/backend/pkg/common"
	"github.com/cmclaughlin24/sundance/backend/pkg/common/validate"
	"github.com/cmclaughlin24/sundance/forms/internal/adapters/rest/dto"
	"github.com/cmclaughlin24/sundance/forms/internal/core"
	"github.com/cmclaughlin24/sundance/forms/internal/core/domain"
	"github.com/cmclaughlin24/sundance/forms/internal/core/ports"
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

		common.SendJsonResponse(w, http.StatusOK, dtos)
	}
}

func (h *handlers) getForm(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromContext(r.Context())

	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	formID := h.getFormIdPathValue(r)
	query := ports.NewFindByIDQuery(formID, tenantID)
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

		common.SendJsonResponse(w, http.StatusOK, dto.FormToResponse(res.data))
	}
}

func (h *handlers) createForm(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromContext(r.Context())
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	resultChan := make(chan result[*domain.Form], 1)

	var body dto.UpsertFormRequest
	if err := common.ReadJsonPayload(r, &body); err != nil {
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

		common.SendJsonResponse(w, http.StatusCreated, common.ApiResponse[*dto.FormResponse]{
			Message: "Successfully created!",
			Data:    dto.FormToResponse(res.data),
		})
	}
}

func (h *handlers) updateForm(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromContext(r.Context())
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	formID := h.getFormIdPathValue(r)
	resultChan := make(chan result[*domain.Form], 1)

	var body dto.UpsertFormRequest
	if err := common.ReadJsonPayload(r, &body); err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	command := ports.NewUpdateFormCommand(formID, tenantID, body.Name, body.Description)

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

		common.SendJsonResponse(w, http.StatusOK, common.ApiResponse[*dto.FormResponse]{
			Message: "Successfully updated!",
			Data:    dto.FormToResponse(res.data),
		})
	}
}

func (h *handlers) getVersions(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromContext(r.Context())
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	formID := h.getFormIdPathValue(r)
	query := ports.NewFindVersionsQuery(formID, tenantID)
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

		common.SendJsonResponse(w, http.StatusOK, dtos)
	}
}

func (h *handlers) getVersion(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromContext(r.Context())
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	formID := h.getFormIdPathValue(r)
	versionID := h.getVersionIdPathValue(r)
	query := ports.NewFindVersionByIDQuery(formID, tenantID, versionID)
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

		common.SendJsonResponse(w, http.StatusOK, dto.VersionToResponse(res.data))
	}
}

func (h *handlers) createVersion(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromContext(r.Context())
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	formID := h.getFormIdPathValue(r)
	resultChan := make(chan result[*domain.Version], 1)

	var body dto.CreateVersionDto
	if err := common.ReadJsonPayload(r, &body); err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	command := ports.NewCreateVersionCommand(formID, tenantID)

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

		common.SendJsonResponse(w, http.StatusCreated, common.ApiResponse[*dto.VersionResponseDto]{
			Message: "Successfully created!",
			Data:    dto.VersionToResponse(res.data),
		})
	}
}

func (h *handlers) updateVersion(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromContext(r.Context())
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	formID := h.getFormIdPathValue(r)
	versionID := h.getVersionIdPathValue(r)
	resultChan := make(chan result[*domain.Version], 1)

	var body dto.UpdateVersionDto
	if err := common.ReadJsonPayload(r, &body); err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	pages, err := dto.RequestToPages(body)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	command := ports.NewUpdateVersionCommand(versionID, formID, tenantID, pages)

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

		common.SendJsonResponse(w, http.StatusOK, common.ApiResponse[*dto.VersionResponseDto]{
			Message: "Successfully updated!",
			Data:    dto.VersionToResponse(res.data),
		})
	}
}

func (h *handlers) publishVersion(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromContext(r.Context())
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	formID := h.getFormIdPathValue(r)
	versionID := h.getVersionIdPathValue(r)
	resultChan := make(chan result[*domain.Version], 1)
	// FIXME: Remove temporary placeholder for user ID.
	command := ports.NewPublishVersionCommand(formID, tenantID, versionID, "placeholder")

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

		common.SendJsonResponse(w, http.StatusOK, common.ApiResponse[any]{
			Message: "Successfully published!",
			Data:    nil,
		})
	}
}

func (h *handlers) retireVersion(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromContext(r.Context())
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	formId := h.getFormIdPathValue(r)
	versionId := h.getVersionIdPathValue(r)
	resultChan := make(chan result[*domain.Version], 1)
	// FIXME: Remove temporary placeholder for user ID.
	command := ports.NewRetireVersionCommand(formId, tenantID, versionId, "placeholder")

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

		common.SendJsonResponse(w, http.StatusOK, common.ApiResponse[any]{
			Message: "Successfully retired!",
			Data:    nil,
		})
	}
}

func (h *handlers) getFormIdPathValue(r *http.Request) domain.FormID {
	id := r.PathValue("formId")
	return domain.FormID(id)
}

func (h *handlers) getVersionIdPathValue(r *http.Request) domain.VersionID {
	id := r.PathValue("versionId")
	return domain.VersionID(id)
}

func (h *handlers) sendErrorResponse(w http.ResponseWriter, err error) {
	switch {
	case validate.IsValidationErr(err):
		common.SendJsonResponse(w, http.StatusBadRequest, common.ApiErrorResponse{
			Message:    "Bad Request",
			Error:      err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	default:
		common.SendErrorResponse(w, err)
	}

}
