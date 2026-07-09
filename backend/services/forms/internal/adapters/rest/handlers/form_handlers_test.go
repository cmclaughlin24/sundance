package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"sundance/backend/pkg/auth"
	"sundance/backend/pkg/common"
	"sundance/backend/pkg/common/httputil"
	"sundance/backend/services/forms/internal/adapters/rest/dto"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
	"sundance/backend/services/forms/internal/core/ports/commands"

	"github.com/go-chi/chi/v5"
)

func Test_handlers_GetForms(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context, ports.FindFormsQuery) ([]*domain.Form, error)
		statusCode int
		count      int
	}{
		{
			"should yield OK if the request is successful with results",
			func(ctx context.Context, query ports.FindFormsQuery) ([]*domain.Form, error) {
				return []*domain.Form{
					{ID: domain.FormID("1"), TenantID: "tenant-1", Name: "Master Sword"},
					{ID: domain.FormID("2"), TenantID: "tenant-1", Name: "Hylian Shield"},
				}, nil
			},
			http.StatusOK,
			2,
		},
		{
			"should yield OK if the request is successful with empty results",
			func(ctx context.Context, query ports.FindFormsQuery) ([]*domain.Form, error) {
				return []*domain.Form{}, nil
			},
			http.StatusOK,
			0,
		},
		{
			"should yield INTERNAL SERVER ERROR if the request fails",
			func(ctx context.Context, query ports.FindFormsQuery) ([]*domain.Form, error) {
				return nil, errors.New("internal error")
			},
			http.StatusInternalServerError,
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := &ports.API{Forms: &mockFormsService{findFn: tt.fn}}
			h := newTestHandlers(s)
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/forms", nil)
			ctx := httputil.SetTenantContext(req.Context(), "tenant-1")
			req = req.WithContext(ctx)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.GetForms)

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
				t.Errorf("expected %d forms but got %d", tt.count, len(body))
			}
		})
	}
}

func Test_handlers_GetForm(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context, ports.FindByIDQuery[domain.FormID]) (*domain.Form, error)
		statusCode int
		id         string
	}{
		{
			"should yield OK if the request is successful",
			func(ctx context.Context, query ports.FindByIDQuery[domain.FormID]) (*domain.Form, error) {
				return &domain.Form{
					ID:       query.ID,
					TenantID: query.TenantID,
					Name:     "Ocarina of Time",
				}, nil
			},
			http.StatusOK,
			"1",
		},
		{
			"should yield NOT FOUND if the resource is not found",
			func(ctx context.Context, query ports.FindByIDQuery[domain.FormID]) (*domain.Form, error) {
				return nil, common.ErrNotFound
			},
			http.StatusNotFound,
			"1",
		},
		{
			"should yield INTERNAL SERVER ERROR if the request fails",
			func(ctx context.Context, query ports.FindByIDQuery[domain.FormID]) (*domain.Form, error) {
				return nil, errors.New("internal error")
			},
			http.StatusInternalServerError,
			"1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := &ports.API{Forms: &mockFormsService{findByIDFn: tt.fn}}
			h := newTestHandlers(s)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("formId", tt.id)
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/forms/"+tt.id, nil)
			ctx := httputil.SetTenantContext(req.Context(), "tenant-1")
			req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.GetForm)

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

func Test_handlers_CreateForm(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context, commands.CreateFormCommand) (*domain.Form, error)
		statusCode int
		body       dto.UpsertFormRequest
	}{
		{
			"should yield CREATED if the request is successful",
			func(ctx context.Context, command commands.CreateFormCommand) (*domain.Form, error) {
				return domain.NewForm(command.TenantID, command.Name, command.Description)
			},
			http.StatusCreated,
			dto.UpsertFormRequest{Name: "Temple of Time", Description: "The Master Sword is a legendary blade that evil can never touch."},
		},
		{
			"should yield BAD REQUEST if the request body is invalid",
			func(ctx context.Context, command commands.CreateFormCommand) (*domain.Form, error) {
				return nil, nil
			},
			http.StatusBadRequest,
			dto.UpsertFormRequest{Name: "", Description: ""},
		},
		{
			"should yield INTERNAL SERVER ERROR if the request fails",
			func(ctx context.Context, command commands.CreateFormCommand) (*domain.Form, error) {
				return nil, errors.New("internal error")
			},
			http.StatusInternalServerError,
			dto.UpsertFormRequest{Name: "Great Deku Tree", Description: "The Triforce consists of three golden triangles representing Power, Wisdom, and Courage."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := &ports.API{Forms: &mockFormsService{createFn: tt.fn}}
			h := newTestHandlers(s)
			body, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/forms", bytes.NewReader(body))
			ctx := httputil.SetTenantContext(req.Context(), "tenant-1")
			req = req.WithContext(ctx)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.CreateForm)

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

func Test_handlers_UpdateForm(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context, commands.UpdateFormCommand) (*domain.Form, error)
		statusCode int
		id         string
		body       dto.UpsertFormRequest
	}{
		{
			"should yield OK if the request is successful",
			func(ctx context.Context, command commands.UpdateFormCommand) (*domain.Form, error) {
				return domain.NewForm(command.TenantID, command.Name, command.Description)
			},
			http.StatusOK,
			"1",
			dto.UpsertFormRequest{Name: "Lon Lon Ranch", Description: "Link has been reincarnated across countless generations to battle evil."},
		},
		{
			"should yield BAD REQUEST if the request body is invalid",
			func(ctx context.Context, command commands.UpdateFormCommand) (*domain.Form, error) {
				return nil, nil
			},
			http.StatusBadRequest,
			"1",
			dto.UpsertFormRequest{Name: "", Description: ""},
		},
		{
			"should yield INTERNAL SERVER ERROR if the request fails",
			func(ctx context.Context, command commands.UpdateFormCommand) (*domain.Form, error) {
				return nil, errors.New("internal error")
			},
			http.StatusInternalServerError,
			"2",
			dto.UpsertFormRequest{Name: "Gerudo Valley", Description: "Ganondorf is the sole male Gerudo born every hundred years, destined to bear the Triforce of Power."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := &ports.API{Forms: &mockFormsService{updateFn: tt.fn}}
			h := newTestHandlers(s)
			body, _ := json.Marshal(tt.body)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("formId", tt.id)
			req, _ := http.NewRequest(http.MethodPut, "/api/v1/forms/"+tt.id, bytes.NewReader(body))
			ctx := httputil.SetTenantContext(req.Context(), "tenant-1")
			req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.UpdateForm)

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

func Test_handlers_DeleteForm(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context, commands.DeleteCommand[domain.FormID]) error
		statusCode int
		id         string
	}{
		{
			"should yield NO CONTENT if the request is successful",
			func(ctx context.Context, command commands.DeleteCommand[domain.FormID]) error {
				return nil
			},
			http.StatusNoContent,
			"1",
		},
		{
			"should yield NOT FOUND if the resource is not found",
			func(ctx context.Context, command commands.DeleteCommand[domain.FormID]) error {
				return common.ErrNotFound
			},
			http.StatusNotFound,
			"1",
		},
		{
			"should yield INTERNAL SERVER ERROR if the request fails",
			func(ctx context.Context, command commands.DeleteCommand[domain.FormID]) error {
				return errors.New("internal error")
			},
			http.StatusInternalServerError,
			"1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := &ports.API{Forms: &mockFormsService{deleteFn: tt.fn}}
			h := newTestHandlers(s)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("formId", tt.id)
			req, _ := http.NewRequest(http.MethodDelete, "/api/v1/forms/"+tt.id, nil)
			ctx := httputil.SetTenantContext(req.Context(), "tenant-1")
			req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.DeleteForm)

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

func Test_handlers_GetFormVersions(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context, ports.FindFormVersionsQuery) ([]*domain.FormVersion, error)
		statusCode int
		count      int
	}{
		{
			"should yield OK if the request is successful with results",
			func(ctx context.Context, query ports.FindFormVersionsQuery) ([]*domain.FormVersion, error) {
			v1, _ := domain.NewFormVersion("form-1", 1, domain.FormVersionStatusDraft, nil)
			v2, _ := domain.NewFormVersion("form-1", 2, domain.FormVersionStatusActive, nil)
				return []*domain.FormVersion{v1, v2}, nil
			},
			http.StatusOK,
			2,
		},
		{
			"should yield OK if the request is successful with empty results",
			func(ctx context.Context, query ports.FindFormVersionsQuery) ([]*domain.FormVersion, error) {
				return []*domain.FormVersion{}, nil
			},
			http.StatusOK,
			0,
		},
		{
			"should yield INTERNAL SERVER ERROR if the request fails",
			func(ctx context.Context, query ports.FindFormVersionsQuery) ([]*domain.FormVersion, error) {
				return nil, errors.New("internal error")
			},
			http.StatusInternalServerError,
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := &ports.API{Forms: &mockFormsService{findVersionsFn: tt.fn}}
			h := newTestHandlers(s)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("formId", "form-1")
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/forms/form-1/versions", nil)
			ctx := httputil.SetTenantContext(req.Context(), "tenant-1")
			req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.GetFormVersions)

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
				t.Errorf("expected %d versions but got %d", tt.count, len(body))
			}
		})
	}
}

func Test_handlers_GetFormVersion(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context, ports.FindFormVersionByIDQuery) (*domain.FormVersion, error)
		statusCode int
		versionId  string
	}{
		{
			"should yield OK if the request is successful",
			func(ctx context.Context, query ports.FindFormVersionByIDQuery) (*domain.FormVersion, error) {
				v, _ := domain.NewFormVersion(query.ID, 1, domain.FormVersionStatusDraft, nil)
				return v, nil
			},
			http.StatusOK,
			"v-1",
		},
		{
			"should yield NOT FOUND if the resource is not found",
			func(ctx context.Context, query ports.FindFormVersionByIDQuery) (*domain.FormVersion, error) {
				return nil, common.ErrNotFound
			},
			http.StatusNotFound,
			"v-1",
		},
		{
			"should yield INTERNAL SERVER ERROR if the request fails",
			func(ctx context.Context, query ports.FindFormVersionByIDQuery) (*domain.FormVersion, error) {
				return nil, errors.New("internal error")
			},
			http.StatusInternalServerError,
			"v-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := &ports.API{Forms: &mockFormsService{findVersionFn: tt.fn}}
			h := newTestHandlers(s)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("formId", "form-1")
			rctx.URLParams.Add("versionId", tt.versionId)
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/forms/form-1/versions/"+tt.versionId, nil)
			ctx := httputil.SetTenantContext(req.Context(), "tenant-1")
			req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.GetFormVersion)

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

func Test_handlers_CreateFormVersion(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context, *commands.CreateFormVersionCommand) (*domain.FormVersion, error)
		statusCode int
		body       dto.UpsertFormVersionRequest
	}{
		{
			"should yield CREATED if the request is successful",
			func(ctx context.Context, command *commands.CreateFormVersionCommand) (*domain.FormVersion, error) {
			v, _ := domain.NewFormVersion(command.FormID, 1, domain.FormVersionStatusDraft, nil)
			return v, nil
			},
			http.StatusCreated,
			dto.UpsertFormVersionRequest{
				Pages: []dto.PageRequest{
					{Key: "p1", Name: "Hyrule Field", Position: 0, Sections: []dto.SectionRequest{}, Rules: []dto.RuleRequest{}},
				},
			},
		},
		{
			"should yield INTERNAL SERVER ERROR if the request fails",
			func(ctx context.Context, command *commands.CreateFormVersionCommand) (*domain.FormVersion, error) {
				return nil, errors.New("internal error")
			},
			http.StatusInternalServerError,
			dto.UpsertFormVersionRequest{
				Pages: []dto.PageRequest{
					{Key: "p1", Name: "Hyrule Field", Position: 0, Sections: []dto.SectionRequest{}, Rules: []dto.RuleRequest{}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := &ports.API{Forms: &mockFormsService{createVersionFn: tt.fn}}
			h := newTestHandlers(s)
			body, _ := json.Marshal(tt.body)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("formId", "form-1")
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/forms/form-1/versions", bytes.NewReader(body))
			ctx := httputil.SetTenantContext(req.Context(), "tenant-1")
			req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.CreateFormVersion)

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

func Test_handlers_UpdateFormVersion(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context, *commands.UpdateFormVersionCommand) (*domain.FormVersion, error)
		statusCode int
		body       dto.UpsertFormVersionRequest
	}{
		{
			"should yield OK if the request is successful",
			func(ctx context.Context, command *commands.UpdateFormVersionCommand) (*domain.FormVersion, error) {
			v, _ := domain.NewFormVersion(command.FormID, 1, domain.FormVersionStatusDraft, nil)
			return v, nil
			},
			http.StatusOK,
			dto.UpsertFormVersionRequest{
				Pages: []dto.PageRequest{
					{Key: "p1", Name: "Hyrule Field", Position: 0, Sections: []dto.SectionRequest{}, Rules: []dto.RuleRequest{}},
				},
			},
		},
		{
			"should yield INTERNAL SERVER ERROR if the request fails",
			func(ctx context.Context, command *commands.UpdateFormVersionCommand) (*domain.FormVersion, error) {
				return nil, errors.New("internal error")
			},
			http.StatusInternalServerError,
			dto.UpsertFormVersionRequest{
				Pages: []dto.PageRequest{
					{Key: "p1", Name: "Hyrule Field", Position: 0, Sections: []dto.SectionRequest{}, Rules: []dto.RuleRequest{}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := &ports.API{Forms: &mockFormsService{updateVersionFn: tt.fn}}
			h := newTestHandlers(s)
			body, _ := json.Marshal(tt.body)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("formId", "form-1")
			rctx.URLParams.Add("versionId", "v-1")
			req, _ := http.NewRequest(http.MethodPut, "/api/v1/forms/form-1/versions/v-1", bytes.NewReader(body))
			ctx := httputil.SetTenantContext(req.Context(), "tenant-1")
			req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.UpdateFormVersion)

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

func Test_handlers_PublishFormVersion(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context, commands.PublishFormVersionCommand) (*domain.FormVersion, error)
		statusCode int
	}{
		{
			"should yield OK if the request is successful",
			func(ctx context.Context, command commands.PublishFormVersionCommand) (*domain.FormVersion, error) {
				v, _ := domain.NewFormVersion(command.FormID, 1, domain.FormVersionStatusActive, nil)
				return v, nil
			},
			http.StatusOK,
		},
		{
			"should yield INTERNAL SERVER ERROR if the request fails",
			func(ctx context.Context, command commands.PublishFormVersionCommand) (*domain.FormVersion, error) {
				return nil, errors.New("internal error")
			},
			http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := &ports.API{Forms: &mockFormsService{publishVersionFn: tt.fn}}
			h := newTestHandlers(s)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("formId", "form-1")
			rctx.URLParams.Add("versionId", "v-1")
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/forms/form-1/versions/v-1/publish", nil)
			ctx := httputil.SetTenantContext(req.Context(), "tenant-1")
			ctx = auth.SetClaimsContext(ctx, &mockClaims{subject: "user-1"})
			req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.PublishFormVersion)

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

func Test_handlers_RetireFormVersion(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context, commands.RetireFormVersionCommand) (*domain.FormVersion, error)
		statusCode int
	}{
		{
			"should yield OK if the request is successful",
			func(ctx context.Context, command commands.RetireFormVersionCommand) (*domain.FormVersion, error) {
				v, _ := domain.NewFormVersion(command.FormID, 1, domain.FormVersionStatusRetired, nil)
				return v, nil
			},
			http.StatusOK,
		},
		{
			"should yield INTERNAL SERVER ERROR if the request fails",
			func(ctx context.Context, command commands.RetireFormVersionCommand) (*domain.FormVersion, error) {
				return nil, errors.New("internal error")
			},
			http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := &ports.API{Forms: &mockFormsService{retireVersionFn: tt.fn}}
			h := newTestHandlers(s)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("formId", "form-1")
			rctx.URLParams.Add("versionId", "v-1")
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/forms/form-1/versions/v-1/retire", nil)
			ctx := httputil.SetTenantContext(req.Context(), "tenant-1")
			ctx = auth.SetClaimsContext(ctx, &mockClaims{subject: "user-1"})
			req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.RetireFormVersion)

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
