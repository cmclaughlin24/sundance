package handlers

import (
	"net/http"
	"sundance/backend/pkg/common/httputil"
	"sundance/backend/services/tenants/internal/adapters/rest/dto"
	"sundance/backend/services/tenants/internal/core/domain"
	"sundance/backend/services/tenants/internal/core/ports/commands"

	"github.com/go-chi/chi/v5"
)

// @summary		Get all tenants
// @tags		Tenants
// @accept		json
// @produce		json
// @param 		X-Request-ID header string false "Client-supplied request trace ID (generated if absent)"
// @param 		X-Correlation-ID header string false "Client-supplied correlation ID for tracing"
// @param 		X-Request-Date header string false "Client-supplied request date in ISO 8601 format" Format(date)
// @success		200 {array} dto.TenantResponse
// @failure		500 {object} httputil.APIErrorResponse
// @security 	BearerAuth
// @Router		/tenants [get]
func (h *Handlers) GetTenants(w http.ResponseWriter, r *http.Request) {
	resultChan := make(chan result[[]*domain.Tenant], 1)

	go func() {
		defer close(resultChan)
		tenants, err := h.app.API.Tenants.Find(r.Context())
		resultChan <- result[[]*domain.Tenant]{tenants, err}
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

		dtos := make([]*dto.TenantResponse, 0, len(res.data))
		for _, tenant := range res.data {
			dtos = append(dtos, dto.TenantToResponse(tenant))
		}

		httputil.SendJSONResponse(w, http.StatusOK, dtos)
	}
}

// @summary		Get a tenant by ID
// @tags		Tenants
// @accept		json
// @produce		json
// @param		id path string true "Tenant ID"
// @param 		X-Request-ID header string false "Client-supplied request trace ID (generated if absent)"
// @param 		X-Correlation-ID header string false "Client-supplied correlation ID for tracing"
// @param 		X-Request-Date header string false "Client-supplied request date in ISO 8601 format" Format(date)
// @success		200 {object} dto.TenantResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @security 	BearerAuth
// @Router		/tenants/{id} [get]
func (h *Handlers) GetTenant(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantIDPathValue(r)
	resultChan := make(chan result[*domain.Tenant], 1)

	go func() {
		defer close(resultChan)
		tenant, err := h.app.API.Tenants.FindByID(r.Context(), tenantID)
		resultChan <- result[*domain.Tenant]{tenant, err}
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

		httputil.SendJSONResponse(w, http.StatusOK, dto.TenantToResponse(res.data))
	}
}

// @summary		Create a tenant
// @tags		Tenants
// @accept		json
// @produce		json
// @param		body body dto.TenantRequest true "Create Tenant"
// @param 		X-Request-ID header string false "Client-supplied request trace ID (generated if absent)"
// @param 		X-Correlation-ID header string false "Client-supplied correlation ID for tracing"
// @param 		X-Request-Date header string false "Client-supplied request date in ISO 8601 format" Format(date)
// @success		201 {object} httputil.APIResponse[dto.TenantResponse]
// @failure		400 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @security 	BearerAuth
// @Router		/tenants [post]
func (h *Handlers) CreateTenant(w http.ResponseWriter, r *http.Request) {
	var body dto.TenantRequest
	if err := httputil.ReadValidateJSONPayload(r, &body); err != nil {
		h.app.Logger.WarnContext(r.Context(), "invalid request body", "error", err)
		h.sendErrorResponse(w, err)
		return
	}

	resultChan := make(chan result[*domain.Tenant], 1)
	command := commands.NewCreateTenantCommand(body.Name, body.Description)

	go func() {
		defer close(resultChan)
		tenant, err := h.app.API.Tenants.Create(r.Context(), command)
		resultChan <- result[*domain.Tenant]{tenant, err}
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

		httputil.SendJSONResponse(w, http.StatusCreated, httputil.APIResponse[*dto.TenantResponse]{
			Message: "Successfully created!",
			Data:    dto.TenantToResponse(res.data),
		})
	}
}

// @summary		Update a tenant
// @tags		Tenants
// @accept		json
// @produce		json
// @param		id path string true "Tenant ID"
// @param		body body dto.TenantRequest true "Update Tenant"
// @param 		X-Request-ID header string false "Client-supplied request trace ID (generated if absent)"
// @param 		X-Correlation-ID header string false "Client-supplied correlation ID for tracing"
// @param 		X-Request-Date header string false "Client-supplied request date in ISO 8601 format" Format(date)
// @success		200 {object} httputil.APIResponse[dto.TenantResponse]
// @failure		400 {object} httputil.APIErrorResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @security 	BearerAuth
// @Router		/tenants/{id} [put]
func (h *Handlers) UpdateTenant(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantIDPathValue(r)

	var body dto.TenantRequest
	if err := httputil.ReadValidateJSONPayload(r, &body); err != nil {
		h.app.Logger.WarnContext(r.Context(), "invalid request body", "error", err)
		h.sendErrorResponse(w, err)
		return
	}

	resultChan := make(chan result[*domain.Tenant], 1)
	command := commands.NewUpdateTenantCommand(tenantID, body.Name, body.Description)

	go func() {
		defer close(resultChan)
		tenant, err := h.app.API.Tenants.Update(r.Context(), command)
		resultChan <- result[*domain.Tenant]{tenant, err}
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

		httputil.SendJSONResponse(w, http.StatusOK, httputil.APIResponse[*dto.TenantResponse]{
			Message: "Successfully updated!",
			Data:    dto.TenantToResponse(res.data),
		})
	}
}

// @summary		Delete a tenant
// @description	All data sources belonging to the tenant will also be deleted.
// @tags		Tenants
// @accept		json
// @produce		json
// @param		id path string true "Tenant ID"
// @param 		X-Request-ID header string false "Client-supplied request trace ID (generated if absent)"
// @param 		X-Correlation-ID header string false "Client-supplied correlation ID for tracing"
// @param 		X-Request-Date header string false "Client-supplied request date in ISO 8601 format" Format(date)
// @success		204
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @security 	BearerAuth
// @Router		/tenants/{id} [delete]
func (h *Handlers) DeleteTenant(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantIDPathValue(r)
	resultChan := make(chan result[any], 1)

	go func() {
		defer close(resultChan)
		err := h.app.API.Tenants.Delete(r.Context(), tenantID)
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

func (h *Handlers) getTenantIDPathValue(r *http.Request) domain.TenantID {
	id := chi.URLParam(r, "tenantId")
	return domain.TenantID(id)
}
