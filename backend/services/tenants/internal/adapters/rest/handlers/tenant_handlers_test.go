package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sundance/backend/pkg/common"
	"sundance/backend/services/tenants/internal/adapters/rest/dto"
	"sundance/backend/services/tenants/internal/core/domain"
	"sundance/backend/services/tenants/internal/core/ports"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_handlers_GetTenants(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context) ([]*domain.Tenant, error)
		statusCode int
		count      int
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
			s := &ports.API{Tenants: &mockTenantsService{findFn: tt.fn}}
			h := newTestHandlers(s)
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/tenants", nil)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.GetTenants)

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

func Test_handlers_GetTenant(t *testing.T) {
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
			s := &ports.API{Tenants: &mockTenantsService{findByIdFn: tt.fn}}
			h := newTestHandlers(s)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("tenantId", tt.id)
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/tenants/"+tt.id, nil)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.GetTenant)

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
func Test_handlers_CreateTenant(t *testing.T) {
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
			s := &ports.API{Tenants: &mockTenantsService{createFn: tt.fn}}
			h := newTestHandlers(s)
			body, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/tenants", bytes.NewReader(body))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.CreateTenant)

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

func Test_handlers_UpdateTenant(t *testing.T) {
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
			s := &ports.API{Tenants: &mockTenantsService{updateFn: tt.fn}}
			h := newTestHandlers(s)
			body, _ := json.Marshal(tt.body)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("tenantId", tt.id)
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/tenants/"+tt.id, bytes.NewReader(body))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.UpdateTenant)

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

func Test_handlers_DeleteTenant(t *testing.T) {
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
			s := &ports.API{Tenants: &mockTenantsService{deleteFn: tt.fn}}
			h := newTestHandlers(s)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("tenantId", tt.id)
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/tenants/"+tt.id, nil)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.DeleteTenant)

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
