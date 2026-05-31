package handlers

import (
	"net/http"
	"sundance/backend/pkg/common/httputil"
	"sundance/backend/services/forms/internal/adapters/rest/dto"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"

	"github.com/go-chi/chi/v5"
)

// @summary		Get all canonical tags
// @tags		Canonical Tags
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @success		200 {array} dto.CanonicalTagResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/canonical-tags [get]
func (h *Handlers) GetCanonicalTags(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
	query := ports.NewCanonicalTagsQuery(tenantID)
	resultChan := make(chan result[[]*domain.CanonicalTag], 1)

	go func() {
		defer close(resultChan)
		tags, err := h.app.API.CanonicalTags.Find(r.Context(), query)
		resultChan <- result[[]*domain.CanonicalTag]{tags, err}
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

		dtos := make([]*dto.CanonicalTagResponse, 0, len(res.data))
		for _, tag := range res.data {
			dtos = append(dtos, dto.CanonicalTagToResponse(tag))
		}

		httputil.SendJSONResponse(w, http.StatusOK, dtos)
	}
}

// @summary		Get a canonical tag by ID
// @tags		Canonical Tags
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param		tagId path string true "Tag ID"
// @success		200 {object} dto.CanonicalTagResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/canonical-tags/{tagId} [get]
func (h *Handlers) GetCanonicalTag(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
	tagID := h.getCanonicalTagIDPathValue(r)
	query := ports.NewFindByIDQuery(tenantID, tagID)
	resultChan := make(chan result[*domain.CanonicalTag], 1)

	go func() {
		defer close(resultChan)
		tag, err := h.app.API.CanonicalTags.FindById(r.Context(), query)
		resultChan <- result[*domain.CanonicalTag]{tag, err}
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

		httputil.SendJSONResponse(w, http.StatusOK, dto.CanonicalTagToResponse(res.data))
	}
}

// @summary		Create a canonical tag
// @description	Creates a new canonical tag for the tenant. The key must be unique within the tenant and is used as the stable identifier for tag-to-field associations.
// @tags		Canonical Tags
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param		body body dto.CreateCanonicalTagRequest true "Create Canonical Tag"
// @success		201 {object} httputil.APIResponse[dto.CanonicalTagResponse]
// @failure		400 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/canonical-tags [post]
func (h *Handlers) CreateCanonicalTag(w http.ResponseWriter, r *http.Request) {
	var body dto.CreateCanonicalTagRequest
	if err := httputil.ReadValidateJSONPayload(r, &body); err != nil {
		h.app.Logger.WarnContext(r.Context(), "invalid request body", "error", err)
		h.sendErrorResponse(w, err)
		return
	}

	tenantID := httputil.TenantFromContext(r.Context())
	resultChan := make(chan result[*domain.CanonicalTag], 1)
	command := ports.NewCreateCanonicalTagCommand(tenantID, body.Key, body.DisplayName)

	go func() {
		defer close(resultChan)
		tag, err := h.app.API.CanonicalTags.Create(r.Context(), command)
		resultChan <- result[*domain.CanonicalTag]{tag, err}
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

		httputil.SendJSONResponse(w, http.StatusCreated, httputil.APIResponse[*dto.CanonicalTagResponse]{
			Message: "Successfully created!",
			Data:    dto.CanonicalTagToResponse(res.data),
		})
	}
}

// @summary		Update a canonical tag
// @description	Updates the display name of a canonical tag. The key is immutable to preserve historical associations.
// @tags		Canonical Tags
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param		tagId path string true "Tag ID"
// @param		body body dto.UpdateCanonicalTagRequest true "Update Canonical Tag"
// @success		200 {object} httputil.APIResponse[dto.CanonicalTagResponse]
// @failure		400 {object} httputil.APIErrorResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/canonical-tags/{tagId} [put]
func (h *Handlers) UpdateCanonicalTag(w http.ResponseWriter, r *http.Request) {
	var body dto.UpdateCanonicalTagRequest
	if err := httputil.ReadValidateJSONPayload(r, &body); err != nil {
		h.app.Logger.WarnContext(r.Context(), "invalid request body", "error", err)
		h.sendErrorResponse(w, err)
		return
	}

	tenantID := httputil.TenantFromContext(r.Context())
	tagID := h.getCanonicalTagIDPathValue(r)
	resultChan := make(chan result[*domain.CanonicalTag], 1)
	command := ports.NewUpdateCanonicalTagCommand(tenantID, tagID, body.DisplayName)

	go func() {
		defer close(resultChan)
		tag, err := h.app.API.CanonicalTags.Update(r.Context(), command)
		resultChan <- result[*domain.CanonicalTag]{tag, err}
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

		httputil.SendJSONResponse(w, http.StatusOK, httputil.APIResponse[*dto.CanonicalTagResponse]{
			Message: "Successfully updated!",
			Data:    dto.CanonicalTagToResponse(res.data),
		})
	}
}

// @summary		Delete a canonical tag
// @description	Deletes a canonical tag. Deletion is only permitted when a tag version has never been associated with any form field, in order to preserve historical context.
// @tags		Canonical Tags
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param		tagId path string true "Tag ID"
// @success		204
// @failure		400 {object} httputil.APIErrorResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		409 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/canonical-tags/{tagId} [delete]
func (h *Handlers) DeleteCanonicalTag(w http.ResponseWriter, r *http.Request) {
	tenantID := httputil.TenantFromContext(r.Context())
	tagID := h.getCanonicalTagIDPathValue(r)
	command := ports.NewDeleteCommand(tenantID, tagID)
	resultChan := make(chan result[any], 1)

	go func() {
		defer close(resultChan)
		err := h.app.API.CanonicalTags.Delete(r.Context(), command)
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

func (h *Handlers) GetCanonicalTagVersions(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) GetCanonicalTagVersion(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) CreateCanonicalTagVersion(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) UpdateCanonicalTagVersion(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) PublishCanonicalTagVersion(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) DeprecateCanonicalTagVersion(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) RetireCanonicalTagVersion(w http.ResponseWriter, r *http.Request) {}

func (h *Handlers) getCanonicalTagIDPathValue(r *http.Request) domain.CanonicalTagID {
	id := chi.URLParam(r, "tagId")
	return domain.CanonicalTagID(id)
}
