package rest

import (
	"errors"
	"net/http"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/httputil"
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

// @summary		Get all tenants
// @tags		Tenants
// @accept		json
// @produce		json
// @success		200 {array} dto.TenantResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/tenants [get]
func (h *handlers) getTenants(w http.ResponseWriter, r *http.Request) {
	resultChan := make(chan result[[]*domain.Tenant], 1)

	go func() {
		defer close(resultChan)
		tenants, err := h.app.Services.Tenants.Find(r.Context())
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
// @success		200 {object} dto.TenantResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/tenants/{id} [get]
func (h *handlers) getTenant(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantIDPathValue(r)
	resultChan := make(chan result[*domain.Tenant], 1)

	go func() {
		defer close(resultChan)
		tenant, err := h.app.Services.Tenants.FindByID(r.Context(), tenantID)
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
// @success		201 {object} httputil.APIResponse[dto.TenantResponse]
// @failure		400 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/tenants [post]
func (h *handlers) createTenant(w http.ResponseWriter, r *http.Request) {
	var body dto.TenantRequest
	if err := httputil.ReadValidateJSONPayload(r, &body); err != nil {
		h.app.Logger.WarnContext(r.Context(), "invalid request body", "error", err)
		h.sendErrorResponse(w, err)
		return
	}

	resultChan := make(chan result[*domain.Tenant], 1)
	command := ports.NewCreateTenantCommand(body.Name, body.Description)

	go func() {
		defer close(resultChan)
		tenant, err := h.app.Services.Tenants.Create(r.Context(), command)
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
// @success		200 {object} httputil.APIResponse[dto.TenantResponse]
// @failure		400 {object} httputil.APIErrorResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/tenants/{id} [put]
func (h *handlers) updateTenant(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantIDPathValue(r)

	var body dto.TenantRequest
	if err := httputil.ReadValidateJSONPayload(r, &body); err != nil {
		h.app.Logger.WarnContext(r.Context(), "invalid request body", "error", err)
		h.sendErrorResponse(w, err)
		return
	}

	resultChan := make(chan result[*domain.Tenant], 1)
	command := ports.NewUpdateTenantCommand(tenantID, body.Name, body.Description)

	go func() {
		defer close(resultChan)
		tenant, err := h.app.Services.Tenants.Update(r.Context(), command)
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
// @success		204
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/tenants/{id} [delete]
func (h *handlers) deleteTenant(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantIDPathValue(r)
	resultChan := make(chan result[any], 1)

	go func() {
		defer close(resultChan)
		err := h.app.Services.Tenants.Delete(r.Context(), tenantID)
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

// @summary		Get all data sources
// @tags		Data Sources
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @success		200 {array} dto.DataSourceResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/data-sources [get]
func (h *handlers) getDataSources(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantFromContext(r)
	query := ports.NewListDataSourceQuery(tenantID)
	resultChan := make(chan result[[]*domain.DataSource], 1)

	go func() {
		defer close(resultChan)
		sources, err := h.app.Services.DataSources.Find(r.Context(), query)
		resultChan <- result[[]*domain.DataSource]{sources, err}
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

		dtos := make([]*dto.DataSourceResponse, 0, len(res.data))
		for _, source := range res.data {
			dtos = append(dtos, dto.DataSourceToResponse(source))
		}

		httputil.SendJSONResponse(w, http.StatusOK, dtos)
	}
}

// @summary		Get a data source by ID
// @tags		Data Sources
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param		id path string true "Data Source ID"
// @success		200 {object} dto.DataSourceResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/data-sources/{id} [get]
func (h *handlers) getDataSource(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantFromContext(r)
	sourceID := chi.URLParam(r, "dataSourceId")
	resultChan := make(chan result[*domain.DataSource], 1)
	query := ports.NewFindDataSourceByIDQuery(tenantID, domain.DataSourceID(sourceID))

	go func() {
		defer close(resultChan)
		source, err := h.app.Services.DataSources.FindByID(r.Context(), query)
		resultChan <- result[*domain.DataSource]{source, err}
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

		httputil.SendJSONResponse(w, http.StatusOK, dto.DataSourceToResponse(res.data))
	}
}

// @summary		Create a data source
// @description	The attributes field is validated against the schema for the specified data source type.
// @tags		Data Sources
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param		body body dto.DataSourceRequest true "Create Data Source"
// @success		201 {object} httputil.APIResponse[dto.DataSourceResponse]
// @failure		400 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/data-sources [post]
func (h *handlers) createDataSource(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantFromContext(r)

	var body dto.DataSourceRequest
	if err := httputil.ReadValidateJSONPayload(r, &body); err != nil {
		h.app.Logger.WarnContext(r.Context(), "invalid request body", "error", err)
		h.sendErrorResponse(w, err)
		return
	}

	attributes, err := dto.RequestToDataSourceAttributes(body.Type, body.Attributes)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	resultChan := make(chan result[*domain.DataSource], 1)
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
		h.app.Logger.WarnContext(r.Context(), "context cancellation")
		return
	case res := <-resultChan:
		if res.err != nil {
			h.sendErrorResponse(w, res.err)
			return
		}

		httputil.SendJSONResponse(w, http.StatusCreated, httputil.APIResponse[*dto.DataSourceResponse]{
			Message: "Successfully created!",
			Data:    dto.DataSourceToResponse(res.data),
		})
	}
}

// @summary		Update a data source
// @description	The attributes field is validated against the schema for the specified data source type.
// @tags		Data Sources
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param		id path string true "Data Source ID"
// @param		body body dto.DataSourceRequest true "Update Data Source"
// @success		200 {object} httputil.APIResponse[dto.DataSourceResponse]
// @failure		400 {object} httputil.APIErrorResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/data-sources/{id} [put]
func (h *handlers) updateDataSource(w http.ResponseWriter, r *http.Request) {
	var body dto.DataSourceRequest
	if err := httputil.ReadValidateJSONPayload(r, &body); err != nil {
		h.app.Logger.WarnContext(r.Context(), "invalid request body", "error", err)
		h.sendErrorResponse(w, err)
		return
	}

	tenantID := h.getTenantFromContext(r)
	sourceID := chi.URLParam(r, "dataSourceId")
	attributes, err := dto.RequestToDataSourceAttributes(body.Type, body.Attributes)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	resultChan := make(chan result[*domain.DataSource], 1)
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
		h.app.Logger.WarnContext(r.Context(), "context cancellation")
		return
	case res := <-resultChan:
		if res.err != nil {
			h.sendErrorResponse(w, res.err)
			return
		}

		httputil.SendJSONResponse(w, http.StatusOK, httputil.APIResponse[*dto.DataSourceResponse]{
			Message: "Successfully updated!",
			Data:    dto.DataSourceToResponse(res.data),
		})
	}
}

// @summary		Delete a data source
// @tags		Data Sources
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param		id path string true "Data Source ID"
// @success		204
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/data-sources/{id} [delete]
func (h *handlers) deleteDataSource(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantFromContext(r)
	sourceID := chi.URLParam(r, "dataSourceId")
	resultChan := make(chan result[any], 1)
	command := ports.NewRemoveDataSourceCommand(tenantID, domain.DataSourceID(sourceID))

	go func() {
		defer close(resultChan)
		err := h.app.Services.DataSources.Delete(r.Context(), command)
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

// @summary		Get data source look-ups
// @description	Returns key-value pairs suitable for populating selection inputs from the specified data source.
// @tags		Data Sources
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @param		id path string true "Data Source ID"
// @success		200 {array} dto.LookupResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/data-sources/{id}/look-ups [get]
func (h *handlers) getLookups(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantFromContext(r)
	sourceID := chi.URLParam(r, "dataSourceId")
	resultChan := make(chan result[[]*domain.Lookup], 1)
	command := ports.NewGetDataSourceLookupsQuery(tenantID, domain.DataSourceID(sourceID))

	go func() {
		defer close(resultChan)
		lookups, err := h.app.Services.DataSources.Lookup(r.Context(), command)
		resultChan <- result[[]*domain.Lookup]{lookups, err}
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

		dtos := dto.LookupsToResponse(res.data)
		httputil.SendJSONResponse(w, http.StatusOK, dtos)
	}
}

func (h *handlers) getTenantFromContext(r *http.Request) domain.TenantID {
	tenantID := httputil.TenantFromContext(r.Context())
	return domain.TenantID(tenantID)
}

func (h *handlers) getTenantIDPathValue(r *http.Request) domain.TenantID {
	id := chi.URLParam(r, "tenantId")
	return domain.TenantID(id)
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
	return errors.Is(err, dto.ErrDataSourceAttrParse) ||
		errors.Is(err, domain.ErrInvalidSourceType) ||
		errors.Is(err, domain.ErrInvalidSourceTypeAttributes)

}
