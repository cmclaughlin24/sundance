package rest

import (
	"net/http"

	"github.com/cmclaughlin24/sundance/common"
	"github.com/cmclaughlin24/sundance/tenants/internal/core"
	"github.com/cmclaughlin24/sundance/tenants/internal/core/domain"
	"github.com/cmclaughlin24/sundance/tenants/internal/core/ports"
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
			common.SendErrorResponse(w, res.err)
			return
		}

		dtos := make([]*tenantResponseDto, 0, len(res.data))
		for _, tenant := range res.data {
			dtos = append(dtos, tenantToResponseDto(tenant))
		}

		common.SendJsonResponse(w, http.StatusOK, dtos)
	}
}

func (h *handlers) getTenant(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("tenantId")
	resultChan := make(chan result[*domain.Tenant], 1)

	go func() {
		defer close(resultChan)
		tenant, err := h.app.Services.Tenants.FindById(r.Context(), domain.TenantID(id))
		resultChan <- result[*domain.Tenant]{tenant, err}
	}()

	select {
	case <-r.Context().Done():
		return
	case res := <-resultChan:
		if res.err != nil {
			common.SendErrorResponse(w, res.err)
			return
		}

		common.SendJsonResponse(w, http.StatusOK, tenantToResponseDto(res.data))
	}
}

func (h *handlers) createTenant(w http.ResponseWriter, r *http.Request) {
	resultChan := make(chan result[*domain.Tenant], 1)

	var dto upsertTenantDto
	if err := common.ReadJsonPayload(r, &dto); err != nil {
		return
	}

	command, err := ports.NewCreateTenantCommand(dto.Name, dto.Description)
	if err != nil {
		common.SendErrorResponse(w, err)
		return
	}

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
			common.SendErrorResponse(w, res.err)
			return
		}

		common.SendJsonResponse(w, http.StatusCreated, common.ApiResponse[tenantResponseDto]{
			Message: "Successfully created!",
			Data:    *tenantToResponseDto(res.data),
		})
	}
}

func (h *handlers) updateTenant(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("tenantId")
	resultChan := make(chan result[*domain.Tenant], 1)

	var dto upsertTenantDto
	if err := common.ReadJsonPayload(r, &dto); err != nil {
		return
	}

	command, err := ports.NewUpdateTenantCommand(domain.TenantID(id), dto.Name, dto.Description)
	if err != nil {
		common.SendErrorResponse(w, err)
		return
	}

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
			common.SendErrorResponse(w, res.err)
			return
		}

		common.SendJsonResponse(w, http.StatusOK, common.ApiResponse[tenantResponseDto]{
			Message: "Success updated!",
			Data:    *tenantToResponseDto(res.data),
		})
	}
}

func (h *handlers) removeTenant(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("tenantId")
	resultChan := make(chan result[any], 1)

	go func() {
		defer close(resultChan)
		err := h.app.Services.Tenants.Remove(r.Context(), domain.TenantID(id))
		resultChan <- result[any]{nil, err}
	}()

	select {
	case <-r.Context().Done():
		return
	case res := <-resultChan:
		if res.err != nil {
			common.SendErrorResponse(w, res.err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *handlers) getDataSources(w http.ResponseWriter, r *http.Request) {
	tenantId := r.PathValue("tenantId")

	query, err := ports.NewListDataSourceQuery(domain.TenantID(tenantId))
	if err != nil {
		common.SendErrorResponse(w, err)
		return
	}

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
			common.SendErrorResponse(w, res.err)
			return
		}

		dtos := make([]*dataSourceResponseDto, 0, len(res.data))
		for _, source := range res.data {
			dtos = append(dtos, dataSourceToResponseDto(source))
		}

		common.SendJsonResponse(w, http.StatusOK, dtos)
	}
}

func (h *handlers) getDataSource(w http.ResponseWriter, r *http.Request) {
	tenantId := r.PathValue("tenantId")
	sourceId := r.PathValue("dataSourceId")
	resultChan := make(chan result[*domain.DataSource], 1)

	go func() {
		defer close(resultChan)
		source, err := h.app.Services.DataSources.FindById(r.Context(), domain.TenantID(tenantId), domain.DataSourceID(sourceId))
		resultChan <- result[*domain.DataSource]{source, err}
	}()

	select {
	case <-r.Context().Done():
		return
	case res := <-resultChan:
		if res.err != nil {
			common.SendErrorResponse(w, res.err)
			return
		}

		common.SendJsonResponse(w, http.StatusOK, dataSourceToResponseDto(res.data))
	}
}

func (h *handlers) createDataSource(w http.ResponseWriter, r *http.Request) {
	tenantId := r.PathValue("tenantId")
	resultChan := make(chan result[*domain.DataSource], 1)

	var dto upsertDataSourceDto
	if err := common.ReadJsonPayload(r, &dto); err != nil {
		return
	}

	command, err := ports.NewCreateDataSourceCommand(domain.TenantID(tenantId), dto.Type, dto.Attributes)
	if err != nil {
		common.SendErrorResponse(w, err)
		return
	}

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
			common.SendErrorResponse(w, res.err)
			return
		}

		common.SendJsonResponse(w, http.StatusCreated, common.ApiResponse[dataSourceResponseDto]{
			Message: "Successfully created!",
			Data:    *dataSourceToResponseDto(res.data),
		})
	}
}

func (h *handlers) updateDataSource(w http.ResponseWriter, r *http.Request) {
	tenantId := r.PathValue("tenantId")
	sourceId := r.PathValue("dataSourceId")
	resultChan := make(chan result[*domain.DataSource], 1)

	var dto upsertDataSourceDto
	if err := common.ReadJsonPayload(r, &dto); err != nil {
		return
	}

	command, err := ports.NewUpdateDataSourceCommand(domain.DataSourceID(sourceId), domain.TenantID(tenantId), dto.Type, dto.Attributes)
	if err != nil {
		common.SendErrorResponse(w, err)
		return
	}

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
			common.SendErrorResponse(w, res.err)
			return
		}

		common.SendJsonResponse(w, http.StatusOK, common.ApiResponse[dataSourceResponseDto]{
			Message: "Successfully updated!",
			Data:    *dataSourceToResponseDto(res.data),
		})
	}
}

func (h *handlers) removeDataSource(w http.ResponseWriter, r *http.Request) {
	tenantId := r.PathValue("tenantId")
	sourceId := r.PathValue("dataSourceId")
	resultChan := make(chan result[any], 1)

	go func() {
		defer close(resultChan)
		err := h.app.Services.DataSources.Remove(r.Context(), domain.TenantID(tenantId), domain.DataSourceID(sourceId))
		resultChan <- result[any]{nil, err}
	}()

	select {
	case <-r.Context().Done():
		return
	case res := <-resultChan:
		if res.err != nil {
			common.SendErrorResponse(w, res.err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *handlers) getDataSourceLookup(w http.ResponseWriter, r *http.Request) {
	tenantId := r.PathValue("tenantId")
	sourceId := r.PathValue("dataSourceId")
	resultChan := make(chan result[[]*domain.DataSourceLookup], 1)

	go func() {
		defer close(resultChan)
		lookups, err := h.app.Services.DataSources.Lookup(r.Context(), domain.TenantID(tenantId), domain.DataSourceID(sourceId))
		resultChan <- result[[]*domain.DataSourceLookup]{lookups, err}
	}()

	select {
	case <-r.Context().Done():
		return
	case res := <-resultChan:
		if res.err != nil {
			common.SendErrorResponse(w, res.err)
			return
		}

		dtos := make([]*dataSourceLookupResponseDto, 0, len(res.data))
		for _, lookup := range res.data {
			dtos = append(dtos, dataSourceLookupToResponseDto(lookup))
		}

		common.SendJsonResponse(w, http.StatusOK, dtos)
	}
}
