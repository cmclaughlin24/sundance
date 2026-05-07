package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cmclaughlin24/sundance/backend/pkg/common"
	"github.com/cmclaughlin24/sundance/backend/pkg/common/tenants"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/adapters/rest/dto"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/ports"
	"github.com/go-chi/chi/v5"
)

func newTestHandlers(services *ports.Services) *handlers {
	var buf bytes.Buffer

	app := &core.Application{
		Services: services,
		Logger:   slog.New(slog.NewTextHandler(&buf, nil)),
	}

	return newHandlers(app)
}

func Test_handlers_getTenants(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context) ([]*domain.Tenant, error)
		statusCode int
		count        int
	}{
		{
			"should yield OK if the request is successful with results",
			func(ctx context.Context) ([]*domain.Tenant, error) {
				return []*domain.Tenant{
					{ID: domain.TenantID("1"), Name: "Tenant One"},
					{ID: domain.TenantID("2"), Name: "Tenant Two"},
				}, nil
			},
			http.StatusOK,
			2,
		},
		{
			"should yield OK if the request is successful with empty results",
			func(ctx context.Context) ([]*domain.Tenant, error) {
				return []*domain.Tenant{}, nil
			},
			http.StatusOK,
			0,
		},
		{
			"should yield INTERNAL SERVER ERROR if the request fails",
			func(ctx context.Context) ([]*domain.Tenant, error) {
				return nil, errors.New("internal error")
			},
			http.StatusInternalServerError,
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ports.Services{Tenants: &mockTenantsService{findFn: tt.fn}}
			h := newTestHandlers(s)
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/tenants", nil)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.getTenants)

			// Act.
			handler.ServeHTTP(rr, req)

			// Assert.
			resp := rr.Result()

			if resp.StatusCode != tt.statusCode {
				t.Errorf("expected status code %d but got %d", tt.statusCode, resp.StatusCode)
			}

			if tt.statusCode < 200 || tt.statusCode >= 300 {
				return
			}

			var body []map[string]any
			if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}

			if len(body) != tt.count {
				t.Errorf("expected %d tenants but got %d", tt.count, len(body))
			}
		})
	}
}

func Test_handlers_getTenant(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context, domain.TenantID) (*domain.Tenant, error)
		statusCode int
		id         string
	}{
		{
			"should yield OK if the request is successful",
			func(ctx context.Context, id domain.TenantID) (*domain.Tenant, error) {
				return &domain.Tenant{
					ID:   id,
					Name: "Tenant 1",
				}, nil
			},
			http.StatusOK,
			"1",
		},
		{
			"should yield NOT FOUND if the resource is not found",
			func(ctx context.Context, id domain.TenantID) (*domain.Tenant, error) {
				return nil, common.ErrNotFound
			},
			http.StatusNotFound,
			"1",
		},
		{
			"should yield INTERNAL SERVER ERROR if the request fails",
			func(ctx context.Context, id domain.TenantID) (*domain.Tenant, error) {
				return nil, errors.New("internal error")
			},
			http.StatusInternalServerError,
			"1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ports.Services{Tenants: &mockTenantsService{findByIdFn: tt.fn}}
			h := newTestHandlers(s)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("tenantId", tt.id)
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/tenants/"+tt.id, nil)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.getTenant)

			// Act.
			handler.ServeHTTP(rr, req)

			// Assert.
			resp := rr.Result()

			if resp.StatusCode != tt.statusCode {
				t.Errorf("expected status code %d but got %d", tt.statusCode, resp.StatusCode)
			}

			if tt.statusCode < 200 || tt.statusCode >= 300 {
				return
			}

			var body map[string]any
			if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}

			if body["id"] != tt.id {
				t.Errorf("expected id %s but got %s", tt.id, body["id"])
			}
		})
	}
}
func Test_handlers_createTenant(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context, *ports.CreateTenantCommand) (*domain.Tenant, error)
		statusCode int
		body       dto.TenantRequest
	}{
		{
			"should yield CREATED if the request is successful",
			func(ctx context.Context, command *ports.CreateTenantCommand) (*domain.Tenant, error) {
				return domain.NewTenant(command.Name, command.Description)
			},
			http.StatusCreated,
			dto.TenantRequest{Name: "Tenant 1", Description: "Most ACDs have a small white star on their forehead, known as a Bentley Mark."},
		},
		{
			"should yield BAD REQUEST if the request body is invalid",
			func(ctx context.Context, command *ports.CreateTenantCommand) (*domain.Tenant, error) {
				return nil, nil
			},
			http.StatusBadRequest,
			dto.TenantRequest{Name: "", Description: ""},
		},
		{
			"should yield INTERNAL SERVER ERROR if the request fails",
			func(ctx context.Context, command *ports.CreateTenantCommand) (*domain.Tenant, error) {
				return nil, errors.New("internal error")
			},
			http.StatusInternalServerError,
			dto.TenantRequest{Name: "Tenant 2", Description: "ACDs are very vocal, communicating with distinct grunts, whines, and demanding barks."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ports.Services{Tenants: &mockTenantsService{createFn: tt.fn}}
			h := newTestHandlers(s)
			body, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/tenants", bytes.NewReader(body))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.createTenant)

			// Act.
			handler.ServeHTTP(rr, req)

			// Assert.
			resp := rr.Result()

			if resp.StatusCode != tt.statusCode {
				t.Errorf("expected status code %d but got %d", tt.statusCode, resp.StatusCode)
			}

			if tt.statusCode < 200 || tt.statusCode >= 300 {
				return
			}
		})
	}
}

func Test_handlers_updateTenant(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context, *ports.UpdateTenantCommand) (*domain.Tenant, error)
		statusCode int
		id         string
		body       dto.TenantRequest
	}{
		{
			"should yield OK if the request is successful",
			func(ctx context.Context, command *ports.UpdateTenantCommand) (*domain.Tenant, error) {
				return domain.NewTenant(command.Name, command.Description)
			},
			http.StatusOK,
			"1",
			dto.TenantRequest{Name: "Tenant 1", Description: "ACDs have thick double coats."},
		},
		{
			"should yield BAD REQUEST if the request body is invalid",
			func(ctx context.Context, command *ports.UpdateTenantCommand) (*domain.Tenant, error) {
				return nil, nil
			},
			http.StatusBadRequest,
			"1",
			dto.TenantRequest{Name: "", Description: ""},
		},
		{
			"should yield INTERNAL SERVER ERROR if the request fails",
			func(ctx context.Context, command *ports.UpdateTenantCommand) (*domain.Tenant, error) {
				return nil, errors.New("internal error")
			},
			http.StatusInternalServerError,
			"2",
			dto.TenantRequest{Name: "Tenant 2", Description: "ACDs are born entirely white and turn either red or blue as time goes on."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ports.Services{Tenants: &mockTenantsService{updateFn: tt.fn}}
			h := newTestHandlers(s)
			body, _ := json.Marshal(tt.body)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("tenantId", tt.id)
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/tenants/"+tt.id, bytes.NewReader(body))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.updateTenant)

			// Act.
			handler.ServeHTTP(rr, req)

			// Assert.
			resp := rr.Result()

			if resp.StatusCode != tt.statusCode {
				t.Errorf("expected status code %d but got %d", tt.statusCode, resp.StatusCode)
			}

			if tt.statusCode < 200 || tt.statusCode >= 300 {
				return
			}
		})
	}
}

func Test_handlers_deleteTenant(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context, domain.TenantID) error
		statusCode int
		id         string
	}{
		{
			"should yield NOT CONTENT if the request is successful",
			func(ctx context.Context, id domain.TenantID) error {
				return nil
			},
			http.StatusNoContent,
			"1",
		},
		{
			"should yield NOT FOUND if the resource is not found",
			func(ctx context.Context, id domain.TenantID) error {
				return common.ErrNotFound
			},
			http.StatusNotFound,
			"1",
		},
		{
			"should yield INTERNAL SERVER ERROR if the request fails",
			func(ctx context.Context, id domain.TenantID) error {
				return errors.New("internal error")
			},
			http.StatusInternalServerError,
			"1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ports.Services{Tenants: &mockTenantsService{deleteFn: tt.fn}}
			h := newTestHandlers(s)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("tenantId", tt.id)
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/tenants/"+tt.id, nil)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.deleteTenant)

			// Act.
			handler.ServeHTTP(rr, req)

			// Assert.
			resp := rr.Result()

			if resp.StatusCode != tt.statusCode {
				t.Errorf("expected status code %d but got %d", tt.statusCode, resp.StatusCode)
			}
		})
	}
}

func Test_handlers_getDataSources(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context, *ports.ListDataSourceQuery) ([]*domain.DataSource, error)
		statusCode int
		count        int
	}{
		{
			"should yield OK if the request is successful with results",
			func(ctx context.Context, query *ports.ListDataSourceQuery) ([]*domain.DataSource, error) {
				return []*domain.DataSource{
					{ID: domain.DataSourceID("1"), TenantID: "tenant-1", Name: "Source One"},
					{ID: domain.DataSourceID("2"), TenantID: "tenant-1", Name: "Source Two"},
				}, nil
			},
			http.StatusOK,
			2,
		},
		{
			"should yield OK if the request is successful with empty results",
			func(ctx context.Context, query *ports.ListDataSourceQuery) ([]*domain.DataSource, error) {
				return []*domain.DataSource{}, nil
			},
			http.StatusOK,
			0,
		},
		{
			"should yield INTERNAL SERVER ERROR if the request fails",
			func(ctx context.Context, query *ports.ListDataSourceQuery) ([]*domain.DataSource, error) {
				return nil, errors.New("internal error")
			},
			http.StatusInternalServerError,
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := &ports.Services{DataSources: &mockDataSourcesService{findFn: tt.fn}}
			h := newTestHandlers(s)
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/data-sources", nil)
			ctx := tenants.SetTenantContext(req.Context(), "tenant-1")
			req = req.WithContext(ctx)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.getDataSources)

			// Act.
			handler.ServeHTTP(rr, req)

			// Assert.
			resp := rr.Result()

			if resp.StatusCode != tt.statusCode {
				t.Errorf("expected status code %d but got %d", tt.statusCode, resp.StatusCode)
			}

			if tt.statusCode < 200 || tt.statusCode >= 300 {
				return
			}

			var body []map[string]any
			if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}

			if len(body) != tt.count {
				t.Errorf("expected %d data sources but got %d", tt.count, len(body))
			}
		})
	}
}

func Test_handlers_getDataSource(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context, *ports.FindDataSourceByIDQuery) (*domain.DataSource, error)
		statusCode int
		id         string
	}{
		{
			"should yield OK if the request is successful",
			func(ctx context.Context, query *ports.FindDataSourceByIDQuery) (*domain.DataSource, error) {
				return &domain.DataSource{
					ID:       query.ID,
					TenantID: query.TenantID,
					Name:     "Source One",
				}, nil
			},
			http.StatusOK,
			"1",
		},
		{
			"should yield NOT FOUND if the resource is not found",
			func(ctx context.Context, query *ports.FindDataSourceByIDQuery) (*domain.DataSource, error) {
				return nil, common.ErrNotFound
			},
			http.StatusNotFound,
			"1",
		},
		{
			"should yield INTERNAL SERVER ERROR if the request fails",
			func(ctx context.Context, query *ports.FindDataSourceByIDQuery) (*domain.DataSource, error) {
				return nil, errors.New("internal error")
			},
			http.StatusInternalServerError,
			"1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := &ports.Services{DataSources: &mockDataSourcesService{findByIdFn: tt.fn}}
			h := newTestHandlers(s)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("dataSourceId", tt.id)
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/data-sources/"+tt.id, nil)
			ctx := tenants.SetTenantContext(req.Context(), "tenant-1")
			req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.getDataSource)

			// Act.
			handler.ServeHTTP(rr, req)

			// Assert.
			resp := rr.Result()

			if resp.StatusCode != tt.statusCode {
				t.Errorf("expected status code %d but got %d", tt.statusCode, resp.StatusCode)
			}

			if tt.statusCode < 200 || tt.statusCode >= 300 {
				return
			}

			var body map[string]any
			if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}

			if body["id"] != tt.id {
				t.Errorf("expected id %s but got %s", tt.id, body["id"])
			}
		})
	}
}

func Test_handlers_createDataSource(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context, *ports.CreateDataSourceCommand) (*domain.DataSource, error)
		statusCode int
		body       dto.DataSourceRequest
	}{
		{
			"should yield CREATED if the request is successful",
			func(ctx context.Context, command *ports.CreateDataSourceCommand) (*domain.DataSource, error) {
				return &domain.DataSource{
					ID:       "ds-1",
					TenantID: command.TenantID,
					Name:     command.Name,
					Type:     command.Type,
				}, nil
			},
			http.StatusCreated,
			dto.DataSourceRequest{
				Name:        "Source One",
				Description: "Australian Cattle Dogs were bred by crossing Dingoes with Collies and other herding dogs.",
				Type:        domain.DataSourceTypeStatic,
				Attributes:  map[string]any{"data": []any{}},
			},
		},
		{
			"should yield BAD REQUEST if the request body is invalid",
			func(ctx context.Context, command *ports.CreateDataSourceCommand) (*domain.DataSource, error) {
				return nil, nil
			},
			http.StatusBadRequest,
			dto.DataSourceRequest{Name: "", Description: "", Type: "", Attributes: nil},
		},
		{
			"should yield INTERNAL SERVER ERROR if the request fails",
			func(ctx context.Context, command *ports.CreateDataSourceCommand) (*domain.DataSource, error) {
				return nil, errors.New("internal error")
			},
			http.StatusInternalServerError,
			dto.DataSourceRequest{
				Name:        "Source Two",
				Description: "ACDs are known to be one of the most intelligent dog breeds in the world.",
				Type:        domain.DataSourceTypeStatic,
				Attributes:  map[string]any{"data": []any{}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := &ports.Services{DataSources: &mockDataSourcesService{createFn: tt.fn}}
			h := newTestHandlers(s)
			body, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/data-sources", bytes.NewReader(body))
			ctx := tenants.SetTenantContext(req.Context(), "tenant-1")
			req = req.WithContext(ctx)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.createDataSource)

			// Act.
			handler.ServeHTTP(rr, req)

			// Assert.
			resp := rr.Result()

			if resp.StatusCode != tt.statusCode {
				t.Errorf("expected status code %d but got %d", tt.statusCode, resp.StatusCode)
			}

			if tt.statusCode < 200 || tt.statusCode >= 300 {
				return
			}
		})
	}
}

func Test_handlers_updateDataSource(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context, *ports.UpdateDataSourceCommand) (*domain.DataSource, error)
		statusCode int
		id         string
		body       dto.DataSourceRequest
	}{
		{
			"should yield OK if the request is successful",
			func(ctx context.Context, command *ports.UpdateDataSourceCommand) (*domain.DataSource, error) {
				return &domain.DataSource{
					ID:       command.ID,
					TenantID: command.TenantID,
					Name:     command.Name,
					Type:     command.Type,
				}, nil
			},
			http.StatusOK,
			"ds-1",
			dto.DataSourceRequest{
				Name:        "Source One",
				Description: "Blue Heelers can have either a blue or red speckled coat pattern.",
				Type:        domain.DataSourceTypeStatic,
				Attributes:  map[string]any{"data": []any{}},
			},
		},
		{
			"should yield BAD REQUEST if the request body is invalid",
			func(ctx context.Context, command *ports.UpdateDataSourceCommand) (*domain.DataSource, error) {
				return nil, nil
			},
			http.StatusBadRequest,
			"ds-1",
			dto.DataSourceRequest{Name: "", Description: "", Type: "", Attributes: nil},
		},
		{
			"should yield INTERNAL SERVER ERROR if the request fails",
			func(ctx context.Context, command *ports.UpdateDataSourceCommand) (*domain.DataSource, error) {
				return nil, errors.New("internal error")
			},
			http.StatusInternalServerError,
			"ds-2",
			dto.DataSourceRequest{
				Name:        "Source Two",
				Description: "An ACD named Bluey holds the record for the oldest dog ever at 29 years and 5 months.",
				Type:        domain.DataSourceTypeStatic,
				Attributes:  map[string]any{"data": []any{}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := &ports.Services{DataSources: &mockDataSourcesService{updateFn: tt.fn}}
			h := newTestHandlers(s)
			body, _ := json.Marshal(tt.body)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("dataSourceId", tt.id)
			req, _ := http.NewRequest(http.MethodPut, "/api/v1/data-sources/"+tt.id, bytes.NewReader(body))
			ctx := tenants.SetTenantContext(req.Context(), "tenant-1")
			req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.updateDataSource)

			// Act.
			handler.ServeHTTP(rr, req)

			// Assert.
			resp := rr.Result()

			if resp.StatusCode != tt.statusCode {
				t.Errorf("expected status code %d but got %d", tt.statusCode, resp.StatusCode)
			}

			if tt.statusCode < 200 || tt.statusCode >= 300 {
				return
			}
		})
	}
}

func Test_handlers_deleteDataSource(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context, *ports.RemoveDataSourceCommand) error
		statusCode int
		id         string
	}{
		{
			"should yield NO CONTENT if the request is successful",
			func(ctx context.Context, command *ports.RemoveDataSourceCommand) error {
				return nil
			},
			http.StatusNoContent,
			"ds-1",
		},
		{
			"should yield NOT FOUND if the resource is not found",
			func(ctx context.Context, command *ports.RemoveDataSourceCommand) error {
				return common.ErrNotFound
			},
			http.StatusNotFound,
			"ds-1",
		},
		{
			"should yield INTERNAL SERVER ERROR if the request fails",
			func(ctx context.Context, command *ports.RemoveDataSourceCommand) error {
				return errors.New("internal error")
			},
			http.StatusInternalServerError,
			"ds-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := &ports.Services{DataSources: &mockDataSourcesService{deleteFn: tt.fn}}
			h := newTestHandlers(s)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("dataSourceId", tt.id)
			req, _ := http.NewRequest(http.MethodDelete, "/api/v1/data-sources/"+tt.id, nil)
			ctx := tenants.SetTenantContext(req.Context(), "tenant-1")
			req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.deleteDataSource)

			// Act.
			handler.ServeHTTP(rr, req)

			// Assert.
			resp := rr.Result()

			if resp.StatusCode != tt.statusCode {
				t.Errorf("expected status code %d but got %d", tt.statusCode, resp.StatusCode)
			}
		})
	}
}

func Test_handlers_getLookups(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context, *ports.GetDataSourceLookupsQuery) ([]*domain.Lookup, error)
		statusCode int
		id         string
		len        int
	}{
		{
			"should yield OK if the request is successful with results",
			func(ctx context.Context, query *ports.GetDataSourceLookupsQuery) ([]*domain.Lookup, error) {
				return []*domain.Lookup{
					{Value: "blue", Label: "Blue Heeler"},
					{Value: "red", Label: "Red Heeler"},
				}, nil
			},
			http.StatusOK,
			"ds-1",
			2,
		},
		{
			"should yield OK if the request is successful with empty results",
			func(ctx context.Context, query *ports.GetDataSourceLookupsQuery) ([]*domain.Lookup, error) {
				return []*domain.Lookup{}, nil
			},
			http.StatusOK,
			"ds-1",
			0,
		},
		{
			"should yield INTERNAL SERVER ERROR if the request fails",
			func(ctx context.Context, query *ports.GetDataSourceLookupsQuery) ([]*domain.Lookup, error) {
				return nil, errors.New("internal error")
			},
			http.StatusInternalServerError,
			"ds-1",
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := &ports.Services{DataSources: &mockDataSourcesService{lookupFn: tt.fn}}
			h := newTestHandlers(s)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("dataSourceId", tt.id)
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/data-sources/"+tt.id+"/look-ups", nil)
			ctx := tenants.SetTenantContext(req.Context(), "tenant-1")
			req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.getLookups)

			// Act.
			handler.ServeHTTP(rr, req)

			// Assert.
			resp := rr.Result()

			if resp.StatusCode != tt.statusCode {
				t.Errorf("expected status code %d but got %d", tt.statusCode, resp.StatusCode)
			}

			if tt.statusCode < 200 || tt.statusCode >= 300 {
				return
			}

			var body []map[string]any
			if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}

			if len(body) != tt.len {
				t.Errorf("expected %d lookups but got %d", tt.len, len(body))
			}
		})
	}
}

func Test_isBadRequest(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			"should yield true when err is ErrDataSourceAttrParse",
			dto.ErrDataSourceAttrParse,
			true,
		},
		{
			"should yield true when err is ErrInvalidSourceType",
			domain.ErrInvalidSourceType,
			true,
		},
		{
			"should yield true when err is ErrInvalidSourceTypeAttributes",
			domain.ErrInvalidSourceTypeAttributes,
			true,
		},
		{
			"should yield false otherwise",
			errors.New("unknown error"),
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isBadRequest(tt.err)

			if got != tt.want {
				t.Errorf("expected %v but got %v", tt.want, got)
			}
		})
	}
}
