package rest

import (
	"errors"
	"net/http"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/httputil"
	"github.com/cmclaughlin24/sundance/backend/pkg/common/tenants"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/adapters/rest/dto"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
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

func (h *handlers) getTenants(w http.ResponseWriter, r *http.Request) {
	resultChan := make(chan result[[]*domain.Tenant], 1)

	go func() {
		defer close(resultChan)
		tenants, err := h.app.Services.Tenants.Find(r.Context())
		resultChan <- result[[]*domain.Tenant]{tenants, err}
	}()

	select {
	case <-r.Context().Done():
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

		httputil.SendJsonResponse(w, http.StatusOK, dtos)
	}
}

func (h *handlers) getTenant(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantIDPathValue(r)
	resultChan := make(chan result[*domain.Tenant], 1)

	go func() {
		defer close(resultChan)
		tenant, err := h.app.Services.Tenants.FindById(r.Context(), tenantID)
		resultChan <- result[*domain.Tenant]{tenant, err}
	}()

	select {
	case <-r.Context().Done():
		return
	case res := <-resultChan:
		if res.err != nil {
			h.sendErrorResponse(w, res.err)
			return
		}

		httputil.SendJsonResponse(w, http.StatusOK, dto.TenantToResponse(res.data))
	}
}

func (h *handlers) createTenant(w http.ResponseWriter, r *http.Request) {
	resultChan := make(chan result[*domain.Tenant], 1)

	var body dto.TenantRequest
	if err := httputil.ReadJsonPayload(r, &body); err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	command := ports.NewCreateTenantCommand(body.Name, body.Description)

	go func() {
		defer close(resultChan)
		tenant, err := h.app.Services.Tenants.Create(r.Context(), command)
		resultChan <- result[*domain.Tenant]{tenant, err}
	}()

	select {
	case <-r.Context().Done():
		return
	case res := <-resultChan:
		if res.err != nil {
			h.sendErrorResponse(w, res.err)
			return
		}

		httputil.SendJsonResponse(w, http.StatusCreated, httputil.ApiResponse[dto.TenantResponse]{
			Message: "Successfully created!",
			Data:    *dto.TenantToResponse(res.data),
		})
	}
}

func (h *handlers) updateTenant(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantIDPathValue(r)
	resultChan := make(chan result[*domain.Tenant], 1)

	var body dto.TenantRequest
	if err := httputil.ReadJsonPayload(r, &body); err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	command := ports.NewUpdateTenantCommand(tenantID, body.Name, body.Description)

	go func() {
		defer close(resultChan)
		tenant, err := h.app.Services.Tenants.Update(r.Context(), command)
		resultChan <- result[*domain.Tenant]{tenant, err}
	}()

	select {
	case <-r.Context().Done():
		return
	case res := <-resultChan:
		if res.err != nil {
			h.sendErrorResponse(w, res.err)
			return
		}

		httputil.SendJsonResponse(w, http.StatusOK, httputil.ApiResponse[dto.TenantResponse]{
			Message: "Successfully updated!",
			Data:    *dto.TenantToResponse(res.data),
		})
	}
}

func (h *handlers) removeTenant(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantIDPathValue(r)
	resultChan := make(chan result[any], 1)

	go func() {
		defer close(resultChan)
		err := h.app.Services.Tenants.Remove(r.Context(), tenantID)
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

func (h *handlers) getDataSources(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantFromContext(r)
	if err != nil {
		httputil.SendErrorResponse(w, err)
		return
	}

	query := ports.NewListDataSourceQuery(tenantID)
	resultChan := make(chan result[[]*domain.DataSource], 1)

	go func() {
		defer close(resultChan)
		sources, err := h.app.Services.DataSources.Find(r.Context(), query)
		resultChan <- result[[]*domain.DataSource]{sources, err}
	}()

	select {
	case <-r.Context().Done():
		return
	case res := <-resultChan:
		if res.err != nil {
			h.sendErrorResponse(w, res.err)
			return
		}

		dtos := make([]*dto.DataSourceResponse, 0, len(res.data))
		for _, source := range res.data {
			dtos = append(dtos, dto.DataSourceToResponse(source))
		}

		httputil.SendJsonResponse(w, http.StatusOK, dtos)
	}
}

func (h *handlers) getDataSource(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantFromContext(r)
	if err != nil {
		httputil.SendErrorResponse(w, err)
		return
	}

	sourceID := chi.URLParam(r, "dataSourceId")
	resultChan := make(chan result[*domain.DataSource], 1)
	query := ports.NewFindDataSourceByID(tenantID, domain.DataSourceID(sourceID))

	go func() {
		defer close(resultChan)
		source, err := h.app.Services.DataSources.FindById(r.Context(), query)
		resultChan <- result[*domain.DataSource]{source, err}
	}()

	select {
	case <-r.Context().Done():
		return
	case res := <-resultChan:
		if res.err != nil {
			h.sendErrorResponse(w, res.err)
			return
		}

		httputil.SendJsonResponse(w, http.StatusOK, dto.DataSourceToResponse(res.data))
	}
}

func (h *handlers) createDataSource(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantFromContext(r)
	if err != nil {
		httputil.SendErrorResponse(w, err)
		return
	}

	resultChan := make(chan result[*domain.DataSource], 1)

	var body dto.DataSourceRequest
	if err := httputil.ReadJsonPayload(r, &body); err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	attributes, err := dto.RequestToDataSourceAttributes(body.Type, body.Attributes)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	command := ports.NewCreateDataSourceCommand(
		tenantID,
		body.Name,
		body.Description,
		body.Type,
		attributes,
	)

	go func() {
		defer close(resultChan)
		tenant, err := h.app.Services.DataSources.Create(r.Context(), command)
		resultChan <- result[*domain.DataSource]{tenant, err}
	}()

	select {
	case <-r.Context().Done():
		return
	case res := <-resultChan:
		if res.err != nil {
			h.sendErrorResponse(w, res.err)
			return
		}

		httputil.SendJsonResponse(w, http.StatusCreated, httputil.ApiResponse[dto.DataSourceResponse]{
			Message: "Successfully created!",
			Data:    *dto.DataSourceToResponse(res.data),
		})
	}
}

func (h *handlers) updateDataSource(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantFromContext(r)
	if err != nil {
		httputil.SendErrorResponse(w, err)
		return
	}

	sourceID := chi.URLParam(r, "dataSourceId")
	resultChan := make(chan result[*domain.DataSource], 1)

	var body dto.DataSourceRequest
	if err := httputil.ReadJsonPayload(r, &body); err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	attributes, err := dto.RequestToDataSourceAttributes(body.Type, body.Attributes)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	command := ports.NewUpdateDataSourceCommand(
		tenantID,
		domain.DataSourceID(sourceID),
		body.Name,
		body.Description,
		body.Type,
		attributes,
	)

	go func() {
		defer close(resultChan)
		tenant, err := h.app.Services.DataSources.Update(r.Context(), command)
		resultChan <- result[*domain.DataSource]{tenant, err}
	}()

	select {
	case <-r.Context().Done():
		return
	case res := <-resultChan:
		if res.err != nil {
			h.sendErrorResponse(w, res.err)
			return
		}

		httputil.SendJsonResponse(w, http.StatusOK, httputil.ApiResponse[dto.DataSourceResponse]{
			Message: "Successfully updated!",
			Data:    *dto.DataSourceToResponse(res.data),
		})
	}
}

func (h *handlers) removeDataSource(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantFromContext(r)
	if err != nil {
		httputil.SendErrorResponse(w, err)
		return
	}

	sourceID := chi.URLParam(r, "dataSourceId")
	resultChan := make(chan result[any], 1)
	command := ports.NewRemoveDataSourceCommand(tenantID, domain.DataSourceID(sourceID))

	go func() {
		defer close(resultChan)
		err := h.app.Services.DataSources.Remove(r.Context(), command)
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

func (h *handlers) getDataSourceLookup(w http.ResponseWriter, r *http.Request) {
	tenantID, err := h.getTenantFromContext(r)
	if err != nil {
		httputil.SendErrorResponse(w, err)
		return
	}

	sourceID := chi.URLParam(r, "dataSourceId")
	resultChan := make(chan result[[]*domain.DataSourceLookup], 1)
	command := ports.NewGetDataSourceLookupsCommand(tenantID, domain.DataSourceID(sourceID))

	go func() {
		defer close(resultChan)
		lookups, err := h.app.Services.DataSources.Lookup(r.Context(), command)
		resultChan <- result[[]*domain.DataSourceLookup]{lookups, err}
	}()

	select {
	case <-r.Context().Done():
		return
	case res := <-resultChan:
		if res.err != nil {
			h.sendErrorResponse(w, res.err)
			return
		}

		dtos := make([]*dto.DataSourceLookupResponse, 0, len(res.data))
		for _, lookup := range res.data {
			dtos = append(dtos, dto.DataSourceLookupToResponse(lookup))
		}

		httputil.SendJsonResponse(w, http.StatusOK, dtos)
	}
}

func (h *handlers) getTenantFromContext(r *http.Request) (domain.TenantID, error) {
	tenantID, err := tenants.TenantFromContext(r.Context())

	if err != nil {
		return "", err
	}

	return domain.TenantID(tenantID), nil
}

func (h *handlers) getTenantIDPathValue(r *http.Request) domain.TenantID {
	id := chi.URLParam(r, "tenantId")
	return domain.TenantID(id)
}

func (h *handlers) sendErrorResponse(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidSourceTypeAttributes):
		httputil.SendJsonResponse(w, http.StatusBadRequest, httputil.ApiErrorResponse{
			Message:    "Bad Request",
			Error:      err.Error(),
			StatusCode: http.StatusBadRequest,
		})

	default:
		httputil.SendErrorResponse(w, err)
	}
}
