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
		// TODO: Throw an ApiResponseError that the request timed out.
	case res := <-resultChan:
		if res.err != nil {
			// TODO: Throw an ApiResponseError that there was an error.
			return
		}
		// TODO: Return a successful response with the data.
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
		// TODO: Throw an ApiResponseError that the request timed out.
	case res := <-resultChan:
		if res.err != nil {
			// TODO: Throw an ApiResponseError that there was an error.
			return
		}
		// TODO: Return a successful response with the data.
	}
}

func (h *handlers) createTenant(w http.ResponseWriter, r *http.Request) {
	resultChan := make(chan result[*domain.Tenant], 1)

	var dto tenantDto
	if err := common.ReadJsonPayload(r, &dto); err != nil {
		return
	}

	go func() {
		defer close(resultChan)
		command := ports.NewCreateTenantCommand(dto.Name, dto.Description)
		tenant, err := h.app.Services.Tenants.Create(r.Context(), command)
		resultChan <- result[*domain.Tenant]{tenant, err}
	}()

	select {
	case <-r.Context().Done():
		// TODO: Throw an ApiResponseError that the request timed out.
	case res := <-resultChan:
		if res.err != nil {
			// TODO: Throw an ApiResponseError that there was an error.
			return
		}
		// TODO: Return a successful response with the data.
	}
}

func (h *handlers) updateTenant(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("tenantId")
	resultChan := make(chan result[*domain.Tenant], 1)

	var dto tenantDto
	if err := common.ReadJsonPayload(r, &dto); err != nil {
		return
	}

	go func() {
		defer close(resultChan)
		command := ports.NewUpdateTenantCommand(domain.TenantID(id), dto.Name, dto.Description)
		tenant, err := h.app.Services.Tenants.Update(r.Context(), command)
		resultChan <- result[*domain.Tenant]{tenant, err}
	}()

	select {
	case <-r.Context().Done():
		// TODO: Throw an ApiResponseError that the request timed out.
	case res := <-resultChan:
		if res.err != nil {
			// TODO: Throw an ApiResponseError that there was an error.
			return
		}
		// TODO: Return a successful response with the data.
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
		// TODO: Throw an ApiResponseError that the request timed out.
	case res := <-resultChan:
		if res.err != nil {
			// TODO: Throw an ApiResponseError that there was an error.
			return
		}
		// TODO: Return a successful response with the data.
	}
}

func (h *handlers) getDataSources(w http.ResponseWriter, r *http.Request) {
	tenantId := r.PathValue("tenantId")
	resultChan := make(chan result[[]*domain.DataSource], 1)

	go func() {
		defer close(resultChan)
		sources, err := h.app.Services.DataSources.Find(r.Context(), ports.NewListDataSourceQuery(domain.TenantID(tenantId)))
		resultChan <- result[[]*domain.DataSource]{sources, err}
	}()

	select {
	case <-r.Context().Done():
		// TODO: Throw an ApiResponseError that the request timed out.
	case res := <-resultChan:
		if res.err != nil {
			// TODO: Throw an ApiResponseError that there was an error.
			return
		}
		// TODO: Return a successful response with the data.
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
		// TODO: Throw an ApiResponseError that the request timed out.
	case res := <-resultChan:
		if res.err != nil {
			// TODO: Throw an ApiResponseError that there was an error.
			return
		}
		// TODO: Return a successful response with the data.
	}
}

func (h *handlers) createDataSource(w http.ResponseWriter, r *http.Request) {
	tenantId := r.PathValue("tenantId")
	resultChan := make(chan result[*domain.DataSource], 1)

	var dto dataSourceDto
	if err := common.ReadJsonPayload(r, &dto); err != nil {
		return
	}

	go func() {
		defer close(resultChan)
		command := ports.NewCreateDataSourceCommand(domain.TenantID(tenantId), dto.Type, dto.Attributes)
		tenant, err := h.app.Services.DataSources.Create(r.Context(), command)
		resultChan <- result[*domain.DataSource]{tenant, err}
	}()

	select {
	case <-r.Context().Done():
		// TODO: Throw an ApiResponseError that the request timed out.
	case res := <-resultChan:
		if res.err != nil {
			// TODO: Throw an ApiResponseError that there was an error.
			return
		}
		// TODO: Return a successful response with the data.
	}
}

func (h *handlers) updateDataSource(w http.ResponseWriter, r *http.Request) {
	tenantId := r.PathValue("tenantId")
	sourceId := r.PathValue("dataSourceId")
	resultChan := make(chan result[*domain.DataSource], 1)

	var dto dataSourceDto
	if err := common.ReadJsonPayload(r, &dto); err != nil {
		return
	}

	go func() {
		defer close(resultChan)
		command := ports.NewUpdateDataSourceCommand(domain.DataSourceID(sourceId), domain.TenantID(tenantId), dto.Type, dto.Attributes)
		tenant, err := h.app.Services.DataSources.Update(r.Context(), command)
		resultChan <- result[*domain.DataSource]{tenant, err}
	}()

	select {
	case <-r.Context().Done():
		// TODO: Throw an ApiResponseError that the request timed out.
	case res := <-resultChan:
		if res.err != nil {
			// TODO: Throw an ApiResponseError that there was an error.
			return
		}
		// TODO: Return a successful response with the data.
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
		// TODO: Throw an ApiResponseError that the request timed out.
	case res := <-resultChan:
		if res.err != nil {
			// TODO: Throw an ApiResponseError that there was an error.
			return
		}
		// TODO: Return a successful response with the data.
	}
}
