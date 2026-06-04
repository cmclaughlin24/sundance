package handlers

import (
	"encoding/json"
	"net/http"
	"sundance/backend/pkg/common/httputil"
	"sundance/backend/services/tenants/internal/adapters/rest/dto"
	"sundance/backend/services/tenants/internal/core/domain"
	"sundance/backend/services/tenants/internal/core/ports"

	"github.com/go-chi/chi/v5"
)

// @summary		Get all data sources
// @tags		Data Sources
// @accept		json
// @produce		json
// @param		X-Tenant-ID header string true "Tenant ID"
// @success		200 {array} dto.DataSourceResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/data-sources [get]
func (h *Handlers) GetDataSources(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantFromContext(r)
	query := ports.NewListDataSourceQuery(tenantID)
	resultChan := make(chan result[[]*domain.DataSource], 1)

	go func() {
		defer close(resultChan)
		sources, err := h.app.API.DataSources.Find(r.Context(), query)
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
func (h *Handlers) GetDataSource(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantFromContext(r)
	sourceID := chi.URLParam(r, "dataSourceId")
	resultChan := make(chan result[*domain.DataSource], 1)
	query := ports.NewFindDataSourceByIDQuery(tenantID, domain.DataSourceID(sourceID))

	go func() {
		defer close(resultChan)
		source, err := h.app.API.DataSources.FindByID(r.Context(), query)
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
func (h *Handlers) CreateDataSource(w http.ResponseWriter, r *http.Request) {
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
		tenant, err := h.app.API.DataSources.Create(r.Context(), command)
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
func (h *Handlers) UpdateDataSource(w http.ResponseWriter, r *http.Request) {
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
		tenant, err := h.app.API.DataSources.Update(r.Context(), command)
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
func (h *Handlers) DeleteDataSource(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantFromContext(r)
	sourceID := chi.URLParam(r, "dataSourceId")
	resultChan := make(chan result[any], 1)
	command := ports.NewRemoveDataSourceCommand(tenantID, domain.DataSourceID(sourceID))

	go func() {
		defer close(resultChan)
		err := h.app.API.DataSources.Delete(r.Context(), command)
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
// @param		query query object false "Optional parameters for look-up retrieval, such as query parameters for external fetches or payload for webhook calls."
// @success		200 {array} dto.LookupResponse
// @failure		404 {object} httputil.APIErrorResponse
// @failure		500 {object} httputil.APIErrorResponse
// @Router		/data-sources/{id}/look-ups [get]
func (h *Handlers) GetLookups(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantFromContext(r)
	sourceID := chi.URLParam(r, "dataSourceId")
	resultChan := make(chan result[[]*domain.Lookup], 1)

	query, err := h.parseDataSourceLookupQuery(r)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	command := ports.NewGetDataSourceLookupsQuery(tenantID, domain.DataSourceID(sourceID), query)

	go func() {
		defer close(resultChan)
		lookups, err := h.app.API.DataSources.Lookup(r.Context(), command)
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

func (h *Handlers) getTenantFromContext(r *http.Request) domain.TenantID {
	tenantID := httputil.TenantFromContext(r.Context())
	return domain.TenantID(tenantID)
}

func (h *Handlers) parseDataSourceLookupQuery(r *http.Request) (map[string]any, error) {
	query := r.URL.Query().Get("query")
	data := make(map[string]any)

	if query == "" {
		return data, nil
	}

	if err := json.Unmarshal([]byte(query), &data); err != nil {
		return nil, err
	}

	return data, nil
}
