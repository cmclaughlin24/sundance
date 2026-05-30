package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sundance/backend/pkg/common"
	"sundance/backend/pkg/common/httputil"
	"sundance/backend/services/tenants/internal/adapters/rest/dto"
	"sundance/backend/services/tenants/internal/core/domain"
	"sundance/backend/services/tenants/internal/core/ports"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_handlers_GetDataSources(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context, *ports.ListDataSourceQuery) ([]*domain.DataSource, error)
		statusCode int
		count      int
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
			ctx := httputil.SetTenantContext(req.Context(), "tenant-1")
			req = req.WithContext(ctx)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.GetDataSources)

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

func Test_handlers_GetDataSource(t *testing.T) {
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
			ctx := httputil.SetTenantContext(req.Context(), "tenant-1")
			req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.GetDataSource)

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

func Test_handlers_CreateDataSource(t *testing.T) {
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
			ctx := httputil.SetTenantContext(req.Context(), "tenant-1")
			req = req.WithContext(ctx)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.CreateDataSource)

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

func Test_handlers_UpdateDataSource(t *testing.T) {
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
			ctx := httputil.SetTenantContext(req.Context(), "tenant-1")
			req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.UpdateDataSource)

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

func Test_handlers_DeleteDataSource(t *testing.T) {
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
			ctx := httputil.SetTenantContext(req.Context(), "tenant-1")
			req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.DeleteDataSource)

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

func Test_handlers_GetLookups(t *testing.T) {
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
			ctx := httputil.SetTenantContext(req.Context(), "tenant-1")
			req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.GetLookups)

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
