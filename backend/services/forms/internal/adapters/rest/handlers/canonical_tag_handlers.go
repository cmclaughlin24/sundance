package handlers

import (
	"net/http"
	"sundance/backend/pkg/common/httputil"
	"sundance/backend/services/forms/internal/adapters/rest/dto"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

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

func (h *Handlers) GetCanonicalTag(w http.ResponseWriter, r *http.Request) {
	// tenantID := httputil.TenantFromContext(r.Context())
	// query := ports.NewFindFormsQuery(tenantID)
	resultChan := make(chan result[*domain.CanonicalTag], 1)

	go func() {
		defer close(resultChan)
		tag, err := h.app.API.CanonicalTags.FindById(r.Context())
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

func (h *Handlers) CreateCanonicalTag(w http.ResponseWriter, r *http.Request) {
	var body dto.UpsertCanonicalTagRequest
	if err := httputil.ReadValidateJSONPayload(r, &body); err != nil {
		h.app.Logger.WarnContext(r.Context(), "invalid request body", "error", err)
		h.sendErrorResponse(w, err)
		return
	}

	// tenantID := httputil.TenantFromContext(r.Context())
	resultChan := make(chan result[*domain.CanonicalTag], 1)
	// command := ports.NewCreateFormCommand(tenantID, body.Name, body.Description)

	go func() {
		defer close(resultChan)
		tag, err := h.app.API.CanonicalTags.Create(r.Context(), nil)
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

func (h *Handlers) DeleteCanonicalTag(w http.ResponseWriter, r *http.Request) {
	// tenantID := httputil.TenantFromContext(r.Context())
	// formID := h.getFormIDPathValue(r)
	// command := ports.NewRemoveFormCommand(tenantID, formID)
	resultChan := make(chan result[any], 1)

	go func() {
		defer close(resultChan)
		err := h.app.API.CanonicalTags.Delete(r.Context(), nil)
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
