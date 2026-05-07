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

	"github.com/cmclaughlin24/sundance/backend/pkg/auth"
	"github.com/cmclaughlin24/sundance/backend/pkg/common"
	"github.com/cmclaughlin24/sundance/backend/pkg/common/tenants"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/adapters/rest/dto"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/ports"
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

func Test_handlers_getForms(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context, *ports.FindFormsQuery) ([]*domain.Form, error)
		statusCode int
		count      int
	}{
		{
			"should yield OK if the request is successful with results",
			func(ctx context.Context, query *ports.FindFormsQuery) ([]*domain.Form, error) {
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
			func(ctx context.Context, query *ports.FindFormsQuery) ([]*domain.Form, error) {
				return []*domain.Form{}, nil
			},
			http.StatusOK,
			0,
		},
		{
			"should yield INTERNAL SERVER ERROR if the request fails",
			func(ctx context.Context, query *ports.FindFormsQuery) ([]*domain.Form, error) {
				return nil, errors.New("internal error")
			},
			http.StatusInternalServerError,
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := &ports.Services{Forms: &mockFormsService{findFn: tt.fn}}
			h := newTestHandlers(s)
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/forms", nil)
			ctx := tenants.SetTenantContext(req.Context(), "tenant-1")
			req = req.WithContext(ctx)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.getForms)

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

func Test_handlers_getForm(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context, *ports.FindFormsByIDQuery) (*domain.Form, error)
		statusCode int
		id         string
	}{
		{
			"should yield OK if the request is successful",
			func(ctx context.Context, query *ports.FindFormsByIDQuery) (*domain.Form, error) {
				return &domain.Form{
					ID:       query.FormID,
					TenantID: query.TenantID,
					Name:     "Ocarina of Time",
				}, nil
			},
			http.StatusOK,
			"1",
		},
		{
			"should yield NOT FOUND if the resource is not found",
			func(ctx context.Context, query *ports.FindFormsByIDQuery) (*domain.Form, error) {
				return nil, common.ErrNotFound
			},
			http.StatusNotFound,
			"1",
		},
		{
			"should yield INTERNAL SERVER ERROR if the request fails",
			func(ctx context.Context, query *ports.FindFormsByIDQuery) (*domain.Form, error) {
				return nil, errors.New("internal error")
			},
			http.StatusInternalServerError,
			"1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := &ports.Services{Forms: &mockFormsService{findByIDFn: tt.fn}}
			h := newTestHandlers(s)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("formId", tt.id)
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/forms/"+tt.id, nil)
			ctx := tenants.SetTenantContext(req.Context(), "tenant-1")
			req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.getForm)

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

func Test_handlers_createForm(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context, *ports.CreateFormCommand) (*domain.Form, error)
		statusCode int
		body       dto.UpsertFormRequest
	}{
		{
			"should yield CREATED if the request is successful",
			func(ctx context.Context, command *ports.CreateFormCommand) (*domain.Form, error) {
				return domain.NewForm(command.TenantID, command.Name, command.Description)
			},
			http.StatusCreated,
			dto.UpsertFormRequest{Name: "Temple of Time", Description: "The Master Sword is a legendary blade that evil can never touch."},
		},
		{
			"should yield BAD REQUEST if the request body is invalid",
			func(ctx context.Context, command *ports.CreateFormCommand) (*domain.Form, error) {
				return nil, nil
			},
			http.StatusBadRequest,
			dto.UpsertFormRequest{Name: "", Description: ""},
		},
		{
			"should yield INTERNAL SERVER ERROR if the request fails",
			func(ctx context.Context, command *ports.CreateFormCommand) (*domain.Form, error) {
				return nil, errors.New("internal error")
			},
			http.StatusInternalServerError,
			dto.UpsertFormRequest{Name: "Great Deku Tree", Description: "The Triforce consists of three golden triangles representing Power, Wisdom, and Courage."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := &ports.Services{Forms: &mockFormsService{createFn: tt.fn}}
			h := newTestHandlers(s)
			body, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/forms", bytes.NewReader(body))
			ctx := tenants.SetTenantContext(req.Context(), "tenant-1")
			req = req.WithContext(ctx)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.createForm)

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

func Test_handlers_updateForm(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context, *ports.UpdateFormCommand) (*domain.Form, error)
		statusCode int
		id         string
		body       dto.UpsertFormRequest
	}{
		{
			"should yield OK if the request is successful",
			func(ctx context.Context, command *ports.UpdateFormCommand) (*domain.Form, error) {
				return domain.NewForm(command.TenantID, command.Name, command.Description)
			},
			http.StatusOK,
			"1",
			dto.UpsertFormRequest{Name: "Lon Lon Ranch", Description: "Link has been reincarnated across countless generations to battle evil."},
		},
		{
			"should yield BAD REQUEST if the request body is invalid",
			func(ctx context.Context, command *ports.UpdateFormCommand) (*domain.Form, error) {
				return nil, nil
			},
			http.StatusBadRequest,
			"1",
			dto.UpsertFormRequest{Name: "", Description: ""},
		},
		{
			"should yield INTERNAL SERVER ERROR if the request fails",
			func(ctx context.Context, command *ports.UpdateFormCommand) (*domain.Form, error) {
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
			s := &ports.Services{Forms: &mockFormsService{updateFn: tt.fn}}
			h := newTestHandlers(s)
			body, _ := json.Marshal(tt.body)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("formId", tt.id)
			req, _ := http.NewRequest(http.MethodPut, "/api/v1/forms/"+tt.id, bytes.NewReader(body))
			ctx := tenants.SetTenantContext(req.Context(), "tenant-1")
			req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.updateForm)

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

func Test_handlers_deleteForm(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context, *ports.RemoveFormCommand) error
		statusCode int
		id         string
	}{
		{
			"should yield NO CONTENT if the request is successful",
			func(ctx context.Context, command *ports.RemoveFormCommand) error {
				return nil
			},
			http.StatusNoContent,
			"1",
		},
		{
			"should yield NOT FOUND if the resource is not found",
			func(ctx context.Context, command *ports.RemoveFormCommand) error {
				return common.ErrNotFound
			},
			http.StatusNotFound,
			"1",
		},
		{
			"should yield INTERNAL SERVER ERROR if the request fails",
			func(ctx context.Context, command *ports.RemoveFormCommand) error {
				return errors.New("internal error")
			},
			http.StatusInternalServerError,
			"1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := &ports.Services{Forms: &mockFormsService{deleteFn: tt.fn}}
			h := newTestHandlers(s)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("formId", tt.id)
			req, _ := http.NewRequest(http.MethodDelete, "/api/v1/forms/"+tt.id, nil)
			ctx := tenants.SetTenantContext(req.Context(), "tenant-1")
			req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.deleteForm)

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

func Test_handlers_getVersions(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context, *ports.FindVersionsQuery) ([]*domain.Version, error)
		statusCode int
		count      int
	}{
		{
			"should yield OK if the request is successful with results",
			func(ctx context.Context, query *ports.FindVersionsQuery) ([]*domain.Version, error) {
				v1, _ := domain.NewVersion("form-1", 1, domain.VersionStatusDraft)
				v2, _ := domain.NewVersion("form-1", 2, domain.VersionStatusActive)
				return []*domain.Version{v1, v2}, nil
			},
			http.StatusOK,
			2,
		},
		{
			"should yield OK if the request is successful with empty results",
			func(ctx context.Context, query *ports.FindVersionsQuery) ([]*domain.Version, error) {
				return []*domain.Version{}, nil
			},
			http.StatusOK,
			0,
		},
		{
			"should yield INTERNAL SERVER ERROR if the request fails",
			func(ctx context.Context, query *ports.FindVersionsQuery) ([]*domain.Version, error) {
				return nil, errors.New("internal error")
			},
			http.StatusInternalServerError,
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := &ports.Services{Forms: &mockFormsService{findVersionsFn: tt.fn}}
			h := newTestHandlers(s)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("formId", "form-1")
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/forms/form-1/versions", nil)
			ctx := tenants.SetTenantContext(req.Context(), "tenant-1")
			req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.getVersions)

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

func Test_handlers_getVersion(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context, *ports.FindVersionByIDQuery) (*domain.Version, error)
		statusCode int
		versionId  string
	}{
		{
			"should yield OK if the request is successful",
			func(ctx context.Context, query *ports.FindVersionByIDQuery) (*domain.Version, error) {
				v, _ := domain.NewVersion(query.FormID, 1, domain.VersionStatusDraft)
				return v, nil
			},
			http.StatusOK,
			"v-1",
		},
		{
			"should yield NOT FOUND if the resource is not found",
			func(ctx context.Context, query *ports.FindVersionByIDQuery) (*domain.Version, error) {
				return nil, common.ErrNotFound
			},
			http.StatusNotFound,
			"v-1",
		},
		{
			"should yield INTERNAL SERVER ERROR if the request fails",
			func(ctx context.Context, query *ports.FindVersionByIDQuery) (*domain.Version, error) {
				return nil, errors.New("internal error")
			},
			http.StatusInternalServerError,
			"v-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := &ports.Services{Forms: &mockFormsService{findVersionFn: tt.fn}}
			h := newTestHandlers(s)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("formId", "form-1")
			rctx.URLParams.Add("versionId", tt.versionId)
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/forms/form-1/versions/"+tt.versionId, nil)
			ctx := tenants.SetTenantContext(req.Context(), "tenant-1")
			req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.getVersion)

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

func Test_handlers_createVersion(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context, *ports.CreateVersionCommand) (*domain.Version, error)
		statusCode int
		body       dto.UpsertVersionRequest
	}{
		{
			"should yield CREATED if the request is successful",
			func(ctx context.Context, command *ports.CreateVersionCommand) (*domain.Version, error) {
				v, _ := domain.NewVersion(command.FormID, 1, domain.VersionStatusDraft)
				return v, nil
			},
			http.StatusCreated,
			dto.UpsertVersionRequest{
				Pages: []dto.PageRequest{
					{Key: "p1", Name: "Hyrule Field", Position: 0, Sections: []dto.SectionRequest{}, Rules: []dto.RuleRequest{}},
				},
			},
		},
		{
			"should yield INTERNAL SERVER ERROR if the request fails",
			func(ctx context.Context, command *ports.CreateVersionCommand) (*domain.Version, error) {
				return nil, errors.New("internal error")
			},
			http.StatusInternalServerError,
			dto.UpsertVersionRequest{
				Pages: []dto.PageRequest{
					{Key: "p1", Name: "Hyrule Field", Position: 0, Sections: []dto.SectionRequest{}, Rules: []dto.RuleRequest{}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := &ports.Services{Forms: &mockFormsService{createVersionFn: tt.fn}}
			h := newTestHandlers(s)
			body, _ := json.Marshal(tt.body)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("formId", "form-1")
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/forms/form-1/versions", bytes.NewReader(body))
			ctx := tenants.SetTenantContext(req.Context(), "tenant-1")
			req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.createVersion)

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

func Test_handlers_updateVersion(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context, *ports.UpdateVersionCommand) (*domain.Version, error)
		statusCode int
		body       dto.UpsertVersionRequest
	}{
		{
			"should yield OK if the request is successful",
			func(ctx context.Context, command *ports.UpdateVersionCommand) (*domain.Version, error) {
				v, _ := domain.NewVersion(command.FormID, 1, domain.VersionStatusDraft)
				return v, nil
			},
			http.StatusOK,
			dto.UpsertVersionRequest{
				Pages: []dto.PageRequest{
					{Key: "p1", Name: "Hyrule Field", Position: 0, Sections: []dto.SectionRequest{}, Rules: []dto.RuleRequest{}},
				},
			},
		},
		{
			"should yield INTERNAL SERVER ERROR if the request fails",
			func(ctx context.Context, command *ports.UpdateVersionCommand) (*domain.Version, error) {
				return nil, errors.New("internal error")
			},
			http.StatusInternalServerError,
			dto.UpsertVersionRequest{
				Pages: []dto.PageRequest{
					{Key: "p1", Name: "Hyrule Field", Position: 0, Sections: []dto.SectionRequest{}, Rules: []dto.RuleRequest{}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := &ports.Services{Forms: &mockFormsService{updateVersionFn: tt.fn}}
			h := newTestHandlers(s)
			body, _ := json.Marshal(tt.body)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("formId", "form-1")
			rctx.URLParams.Add("versionId", "v-1")
			req, _ := http.NewRequest(http.MethodPut, "/api/v1/forms/form-1/versions/v-1", bytes.NewReader(body))
			ctx := tenants.SetTenantContext(req.Context(), "tenant-1")
			req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.updateVersion)

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

func Test_handlers_publishVersion(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context, *ports.PublishVersionCommand) (*domain.Version, error)
		statusCode int
	}{
		{
			"should yield OK if the request is successful",
			func(ctx context.Context, command *ports.PublishVersionCommand) (*domain.Version, error) {
				v, _ := domain.NewVersion(command.FormID, 1, domain.VersionStatusActive)
				return v, nil
			},
			http.StatusOK,
		},
		{
			"should yield INTERNAL SERVER ERROR if the request fails",
			func(ctx context.Context, command *ports.PublishVersionCommand) (*domain.Version, error) {
				return nil, errors.New("internal error")
			},
			http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := &ports.Services{Forms: &mockFormsService{publishVersionFn: tt.fn}}
			h := newTestHandlers(s)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("formId", "form-1")
			rctx.URLParams.Add("versionId", "v-1")
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/forms/form-1/versions/v-1/publish", nil)
			ctx := tenants.SetTenantContext(req.Context(), "tenant-1")
			ctx = auth.SetClaimsContext(ctx, &mockClaims{subject: "user-1"})
			req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.publishVersion)

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

func Test_handlers_retireVersion(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(context.Context, *ports.RetireVersionCommand) (*domain.Version, error)
		statusCode int
	}{
		{
			"should yield OK if the request is successful",
			func(ctx context.Context, command *ports.RetireVersionCommand) (*domain.Version, error) {
				v, _ := domain.NewVersion(command.FormID, 1, domain.VersionStatusRetired)
				return v, nil
			},
			http.StatusOK,
		},
		{
			"should yield INTERNAL SERVER ERROR if the request fails",
			func(ctx context.Context, command *ports.RetireVersionCommand) (*domain.Version, error) {
				return nil, errors.New("internal error")
			},
			http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange.
			s := &ports.Services{Forms: &mockFormsService{retireVersionFn: tt.fn}}
			h := newTestHandlers(s)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("formId", "form-1")
			rctx.URLParams.Add("versionId", "v-1")
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/forms/form-1/versions/v-1/retire", nil)
			ctx := tenants.SetTenantContext(req.Context(), "tenant-1")
			ctx = auth.SetClaimsContext(ctx, &mockClaims{subject: "user-1"})
			req = req.WithContext(context.WithValue(ctx, chi.RouteCtxKey, rctx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.retireVersion)

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

func Test_isBadRequest(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			"should yield true when err is ErrVersionLocked",
			domain.ErrVersionLocked,
			true,
		},
		{
			"should yield true when err is ErrInvalidVersion",
			domain.ErrInvalidVersion,
			true,
		},
		{
			"should yield true when err is ErrInvalidVersionStatus",
			domain.ErrInvalidVersionStatus,
			true,
		},
		{
			"should yield true when err is ErrDuplicateVersion",
			domain.ErrDuplicateVersion,
			true,
		},
		{
			"should yield true when err is ErrInvalidPosition",
			domain.ErrInvalidPosition,
			true,
		},
		{
			"should yield true when err is ErrDuplicatePosition",
			domain.ErrDuplicatePosition,
			true,
		},
		{
			"should yield true when err is ErrInvalidRuleType",
			domain.ErrInvalidRuleType,
			true,
		},
		{
			"should yield true when err is ErrDuplicateRuleType",
			domain.ErrDuplicateRuleType,
			true,
		},
		{
			"should yield true when err is ErrPublishedByRequired",
			domain.ErrPublishedByRequired,
			true,
		},
		{
			"should yield true when err is ErrRetiredByRequired",
			domain.ErrRetiredByRequired,
			true,
		},
		{
			"should yield true when err is ErrInvalidFieldType",
			domain.ErrInvalidFieldType,
			true,
		},
		{
			"should yield true when err is ErrInvalidFieldAttributes",
			domain.ErrInvalidFieldAttributes,
			true,
		},
		{
			"should yield true when err is ErrInvalidForm",
			domain.ErrInvalidForm,
			true,
		},
		{
			"should yield true when err is ErrFormHasActiveVersion",
			domain.ErrFormHasActiveVersion,
			true,
		},
		{
			"should yield true when err is ErrInvalidPage",
			domain.ErrInvalidPage,
			true,
		},
		{
			"should yield true when err is ErrInvalidSection",
			domain.ErrInvalidSection,
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
