package rest

import (
	"errors"
	"net/http"

	"github.com/cmclaughlin24/sundance/backend/pkg/common"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/adapters/rest/dto"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/validate"
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

		common.SendJsonResponse(w, http.StatusOK, dtos)
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

		common.SendJsonResponse(w, http.StatusOK, dto.TenantToResponse(res.data))
	}
}

func (h *handlers) createTenant(w http.ResponseWriter, r *http.Request) {
	resultChan := make(chan result[*domain.Tenant], 1)

	var body dto.TenantRequest
	if err := common.ReadJsonPayload(r, &body); err != nil {
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

		common.SendJsonResponse(w, http.StatusCreated, common.ApiResponse[dto.TenantResponse]{
			Message: "Successfully created!",
			Data:    *dto.TenantToResponse(res.data),
		})
	}
}

func (h *handlers) updateTenant(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantIDPathValue(r)
	resultChan := make(chan result[*domain.Tenant], 1)

	var body dto.TenantRequest
	if err := common.ReadJsonPayload(r, &body); err != nil {
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

		common.SendJsonResponse(w, http.StatusOK, common.ApiResponse[dto.TenantResponse]{
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
	tenantID := h.getTenantIDPathValue(r)
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

		common.SendJsonResponse(w, http.StatusOK, dtos)
	}
}

func (h *handlers) getDataSource(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantIDPathValue(r)
	sourceID := r.PathValue("dataSourceId")
	resultChan := make(chan result[*domain.DataSource], 1)

	go func() {
		defer close(resultChan)
		source, err := h.app.Services.DataSources.FindById(r.Context(), tenantID, domain.DataSourceID(sourceID))
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

		common.SendJsonResponse(w, http.StatusOK, dto.DataSourceToResponse(res.data))
	}
}

func (h *handlers) createDataSource(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantIDPathValue(r)
	resultChan := make(chan result[*domain.DataSource], 1)

	var body dto.DataSourceRequest
	if err := common.ReadJsonPayload(r, &body); err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	attributes, err := dto.RequestToDataSourceAttributes(body.Type, body.Attributes)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	command := ports.NewCreateDataSourceCommand(tenantID, body.Type, attributes)

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

		common.SendJsonResponse(w, http.StatusCreated, common.ApiResponse[dto.DataSourceResponse]{
			Message: "Successfully created!",
			Data:    *dto.DataSourceToResponse(res.data),
		})
	}
}

func (h *handlers) updateDataSource(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantIDPathValue(r)
	sourceID := r.PathValue("dataSourceId")
	resultChan := make(chan result[*domain.DataSource], 1)

	var body dto.DataSourceRequest
	if err := common.ReadJsonPayload(r, &body); err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	attributes, err := dto.RequestToDataSourceAttributes(body.Type, body.Attributes)
	if err != nil {
		h.sendErrorResponse(w, err)
		return
	}

	command := ports.NewUpdateDataSourceCommand(domain.DataSourceID(sourceID), tenantID, body.Type, attributes)

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

		common.SendJsonResponse(w, http.StatusOK, common.ApiResponse[dto.DataSourceResponse]{
			Message: "Successfully updated!",
			Data:    *dto.DataSourceToResponse(res.data),
		})
	}
}

func (h *handlers) removeDataSource(w http.ResponseWriter, r *http.Request) {
	tenantID := h.getTenantIDPathValue(r)
	sourceID := r.PathValue("dataSourceId")
	resultChan := make(chan result[any], 1)

	go func() {
		defer close(resultChan)
		err := h.app.Services.DataSources.Remove(r.Context(), tenantID, domain.DataSourceID(sourceID))
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
	tenantID := h.getTenantIDPathValue(r)
	sourceID := r.PathValue("dataSourceId")
	resultChan := make(chan result[[]*domain.DataSourceLookup], 1)

	go func() {
		defer close(resultChan)
		lookups, err := h.app.Services.DataSources.Lookup(r.Context(), tenantID, domain.DataSourceID(sourceID))
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

		common.SendJsonResponse(w, http.StatusOK, dtos)
	}
}

func (h *handlers) sendErrorResponse(w http.ResponseWriter, err error) {
	switch {
	case validate.IsValidationErr(err) || errors.Is(err, dto.ErrDataSourceAttrParse):
		common.SendJsonResponse(w, http.StatusBadRequest, common.ApiErrorResponse{
			Message:    "Bad Request",
			Error:      err.Error(),
			StatusCode: http.StatusBadRequest,
		})
	default:
		common.SendErrorResponse(w, err)
	}

}

func (h *handlers) getTenantIDPathValue(r *http.Request) domain.TenantID {
	id := r.PathValue("tenantId")
	return domain.TenantID(id)
}
