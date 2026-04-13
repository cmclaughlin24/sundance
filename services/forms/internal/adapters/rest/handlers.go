package rest

import (
	"net/http"

	"github.com/cmclaughlin24/sundance/common"
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
			common.SendErrorResponse(w, res.err)
			return
		}

		dtos := make([]*formResponseDto, 0, len(res.data))
		for _, form := range res.data {
			dtos = append(dtos, formToResponseDto(form))
		}

		common.SendJsonResponse(w, http.StatusOK, res.data)
	}
}

func (h *handlers) getForm(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromContext(r.Context())

	if err != nil {
		common.SendErrorResponse(w, err)
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
			common.SendErrorResponse(w, res.err)
			return
		}

		common.SendJsonResponse(w, http.StatusOK, formToResponseDto(res.data))
	}
}

func (h *handlers) createForm(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromContext(r.Context())
	if err != nil {
		common.SendErrorResponse(w, err)
		return
	}

	resultChan := make(chan result[*domain.Form], 1)

	var dto upsertFormDto
	if err := common.ReadJsonPayload(r, &dto); err != nil {
		return
	}

	command, err := ports.NewCreateFormCommand(tenantID, dto.Name, dto.Description)
	if err != nil {
		common.SendErrorResponse(w, err)
		return
	}

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
			common.SendErrorResponse(w, res.err)
			return
		}

		common.SendJsonResponse(w, http.StatusCreated, common.ApiResponse[*formResponseDto]{
			Message: "Successfully created!",
			Data:    formToResponseDto(res.data),
		})
	}
}

func (h *handlers) updateForm(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromContext(r.Context())
	if err != nil {
		common.SendErrorResponse(w, err)
		return
	}

	formID := h.getFormIdPathValue(r)
	resultChan := make(chan result[*domain.Form], 1)

	var dto upsertFormDto
	if err := common.ReadJsonPayload(r, &dto); err != nil {
		return
	}

	command, err := ports.NewUpdateFormCommand(formID, tenantID, dto.Name, dto.Description)
	if err != nil {
		common.SendErrorResponse(w, err)
		return
	}

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
			common.SendErrorResponse(w, res.err)
			return
		}

		common.SendJsonResponse(w, http.StatusOK, common.ApiResponse[*formResponseDto]{
			Message: "Successfully updated!",
			Data:    formToResponseDto(res.data),
		})
	}
}

func (h *handlers) getVersions(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromContext(r.Context())
	if err != nil {
		common.SendErrorResponse(w, err)
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
			common.SendErrorResponse(w, res.err)
			return
		}

		dtos := make([]*versionResponseDto, 0, len(res.data))
		for _, v := range res.data {
			dtos = append(dtos, versionToResponseDto(v))
		}

		common.SendJsonResponse(w, http.StatusOK, dtos)
	}
}

func (h *handlers) getVersion(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromContext(r.Context())
	if err != nil {
		common.SendErrorResponse(w, err)
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
			common.SendErrorResponse(w, res.err)
			return
		}

		common.SendJsonResponse(w, http.StatusOK, versionToResponseDto(res.data))
	}
}

func (h *handlers) createVersion(w http.ResponseWriter, r *http.Request) {}

func (h *handlers) updateVersion(w http.ResponseWriter, r *http.Request) {}

func (h *handlers) removeVersion(w http.ResponseWriter, r *http.Request) {}

func (h *handlers) publishVersion(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromContext(r.Context())
	if err != nil {
		common.SendErrorResponse(w, err)
		return
	}

	formID := h.getFormIdPathValue(r)
	versionID := h.getVersionIdPathValue(r)
	resultChan := make(chan result[*domain.Version], 1)

	command, err := ports.NewPublishVersionCommand(formID, tenantID, versionID, "")
	if err != nil {
		common.SendErrorResponse(w, err)
		return
	}

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
			common.SendErrorResponse(w, res.err)
			return
		}

		// TODO: Send response.
	}
}

func (h *handlers) retireVersion(w http.ResponseWriter, r *http.Request) {
	tenantID, err := tenantIDFromContext(r.Context())
	if err != nil {
		common.SendErrorResponse(w, err)
		return
	}

	formId := h.getFormIdPathValue(r)
	versionId := h.getVersionIdPathValue(r)
	resultChan := make(chan result[*domain.Version], 1)

	command, err := ports.NewRetireVersionCommand(formId, tenantID, versionId, "")
	if err != nil {
		common.SendErrorResponse(w, err)
		return
	}

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
			common.SendErrorResponse(w, res.err)
			return
		}

		// TODO: Send response.
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
