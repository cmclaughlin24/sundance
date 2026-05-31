package handlers

import (
	"net/http"
	"sundance/backend/pkg/common/httputil"
	"sundance/backend/services/forms/internal/adapters/rest/dto"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"

	"github.com/go-chi/chi/v5"
)

// @summary		Get all tags
// @tags		Tags
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @success		200 {array} dto.TagResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/tags [get]
func (h *Handlers) GetTags(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
	query := ports.NewTagsQuery(tenantID)
	resultChan := make(chan result[[]*domain.Tag], 1)

	go func() {
		defer close(resultChan)
		tags, err := h.app.API.Tags.Find(r.Context(), query)
		resultChan <- result[[]*domain.Tag]{tags, err}
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

		dtos := make([]dto.TagResponse, 0, len(res.data))
		for _, tag := range res.data {
			dtos = append(dtos, dto.TagToResponse(tag))
		}

		httputil.SendJSONResponse(w, http.StatusOK, dtos)
	}
}

// @summary		Get a tag by ID
// @tags		Tags
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param		tagId path string true "Tag ID"
// @success		200 {object} dto.TagResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/tags/{tagId} [get]
func (h *Handlers) GetTag(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
	tagID := h.getTagIDPathValue(r)
	query := ports.NewFindByIDQuery(tenantID, tagID)
	resultChan := make(chan result[*domain.Tag], 1)

	go func() {
		defer close(resultChan)
		tag, err := h.app.API.Tags.FindById(r.Context(), query)
		resultChan <- result[*domain.Tag]{tag, err}
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

		httputil.SendJSONResponse(w, http.StatusOK, dto.TagToResponse(res.data))
	}
}

// @summary		Create a tag
// @description	Creates a new tag for the tenant. The key must be unique within the tenant and is used as the stable identifier for tag-to-field associations.
// @tags		Tags
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param		body body dto.CreateTagRequest true "Create Tag"
// @success		201 {object} httputil.APIResponse[dto.TagResponse]
// @failure		400 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/tags [post]
func (h *Handlers) CreateTag(w http.ResponseWriter, r *http.Request) {
	var body dto.CreateTagRequest
	if err := httputil.ReadValidateJSONPayload(r, &body); err != nil {
		h.app.Logger.WarnContext(r.Context(), "invalid request body", "error", err)
		h.sendErrorResponse(w, err)
		return
	}

	tenantID := httputil.TenantFromContext(r.Context())
	resultChan := make(chan result[*domain.Tag], 1)
	command := ports.NewCreateTagCommand(tenantID, body.Key, body.DisplayName)

	go func() {
		defer close(resultChan)
		tag, err := h.app.API.Tags.Create(r.Context(), command)
		resultChan <- result[*domain.Tag]{tag, err}
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

		httputil.SendJSONResponse(w, http.StatusCreated, httputil.APIResponse[dto.TagResponse]{
			Message: "Successfully created!",
			Data:    dto.TagToResponse(res.data),
		})
	}
}

// @summary		Update a tag
// @description	Updates the display name of a tag. The key is immutable to preserve historical associations.
// @tags		Tags
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param		tagId path string true "Tag ID"
// @param		body body dto.UpdateTagRequest true "Update Tag"
// @success		200 {object} httputil.APIResponse[dto.TagResponse]
// @failure		400 {object} httputil.APIErrorResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/tags/{tagId} [put]
func (h *Handlers) UpdateTag(w http.ResponseWriter, r *http.Request) {
	var body dto.UpdateTagRequest
	if err := httputil.ReadValidateJSONPayload(r, &body); err != nil {
		h.app.Logger.WarnContext(r.Context(), "invalid request body", "error", err)
		h.sendErrorResponse(w, err)
		return
	}

	tenantID := httputil.TenantFromContext(r.Context())
	tagID := h.getTagIDPathValue(r)
	resultChan := make(chan result[*domain.Tag], 1)
	command := ports.NewUpdateTagCommand(tenantID, tagID, body.DisplayName)

	go func() {
		defer close(resultChan)
		tag, err := h.app.API.Tags.Update(r.Context(), command)
		resultChan <- result[*domain.Tag]{tag, err}
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

		httputil.SendJSONResponse(w, http.StatusOK, httputil.APIResponse[dto.TagResponse]{
			Message: "Successfully updated!",
			Data:    dto.TagToResponse(res.data),
		})
	}
}

// @summary		Delete a tag
// @description	Deletes a tag. Deletion is only permitted when a tag version has never been associated with any form field, in order to preserve historical context.
// @tags		Tags
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param		tagId path string true "Tag ID"
// @success		204
// @failure		400 {object} httputil.APIErrorResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		409 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/tags/{tagId} [delete]
func (h *Handlers) DeleteTag(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
	tagID := h.getTagIDPathValue(r)
	command := ports.NewDeleteCommand(tenantID, tagID)
	resultChan := make(chan result[any], 1)

	go func() {
		defer close(resultChan)
		err := h.app.API.Tags.Delete(r.Context(), command)
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

// @summary		Get all tag versions
// @tags		Tags
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param		tagId path string true "Tag ID"
// @success		200 {array} dto.TagVersionResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/tags/{tagId}/versions [get]
func (h *Handlers) GetTagVersions(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
	tagID := h.getTagIDPathValue(r)
	query := ports.NewFindTagVersionsQuery(tenantID, tagID)
	resultChan := make(chan result[[]*domain.TagVersion], 1)

	go func() {
		defer close(resultChan)
		tags, err := h.app.API.Tags.FindVersions(r.Context(), query)
		resultChan <- result[[]*domain.TagVersion]{tags, err}
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

		dtos := make([]dto.TagVersionResponse, 0, len(res.data))
		for _, version := range res.data {
			dtos = append(dtos, dto.TagVersionToResponse(version))
		}

		httputil.SendJSONResponse(w, http.StatusOK, dtos)
	}
}

// @summary		Get a tag version by ID
// @tags		Tags
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param		tagId path string true "Tag ID"
// @param		versionId path string true "Version ID"
// @success		200 {object} dto.TagVersionResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/tags/{tagId}/versions/{versionId} [get]
func (h *Handlers) GetTagVersion(w http.ResponseWriter, r *http.Request) {
}

// @summary		Create a tag version
// @description	Creates a new draft version for the tag. Versions follow a draft → active → deprecated → retired lifecycle and define the data type associated with the tag.
// @tags		Tags
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param		tagId path string true "Tag ID"
// @param		body body dto.UpsertTagVersionRequest true "Create Tag Version"
// @success		201 {object} httputil.APIResponse[dto.TagVersionResponse]
// @failure		400 {object} httputil.APIErrorResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/tags/{tagId}/versions [post]
func (h *Handlers) CreateTagVersion(w http.ResponseWriter, r *http.Request) {
}

func (h *Handlers) UpdateTagVersion(w http.ResponseWriter, r *http.Request) {
}

func (h *Handlers) PublishTagVersion(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) DeprecateTagVersion(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) RetireTagVersion(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) getTagIDPathValue(r *http.Request) domain.TagID {
	id := chi.URLParam(r, "tagId")
	return domain.TagID(id)
}

func (h *Handlers) getTagVersionIDPathValue(r *http.Request) domain.TagVersionID {
	id := chi.URLParam(r, "versionId")
	return domain.TagVersionID(id)
}
