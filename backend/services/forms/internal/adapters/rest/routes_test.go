package rest_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/adapters/rest"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core"
	"github.com/go-chi/chi/v5"
)

func TestNewRoutes(t *testing.T) {
	// Arrange.
	routes := []struct {
		route  string
		method string
	}{
		{"/api/v1/forms/", http.MethodGet},
		{"/api/v1/forms/", http.MethodPost},
		{"/api/v1/forms/{formId}/", http.MethodGet},
		{"/api/v1/forms/{formId}/", http.MethodPut},
		{"/api/v1/forms/{formId}/", http.MethodDelete},
		{"/api/v1/forms/{formId}/versions/", http.MethodGet},
		{"/api/v1/forms/{formId}/versions/", http.MethodPost},
		{"/api/v1/forms/{formId}/versions/{versionId}/", http.MethodGet},
		{"/api/v1/forms/{formId}/versions/{versionId}/", http.MethodPut},
		{"/api/v1/forms/{formId}/versions/{versionId}/publish", http.MethodPost},
		{"/api/v1/forms/{formId}/versions/{versionId}/retire", http.MethodPost},
	}
	mux := rest.NewRoutes(&core.Application{})

	for _, r := range routes {
		// Act/Assert.
		if !routeExists(r.method, r.route, mux.(chi.Routes)) {
			t.Errorf("expected route [%s] %s to exist", r.method, r.route)
		}
	}
}

func routeExists(method, route string, routes chi.Routes) bool {
	var found bool

	chi.Walk(routes, func(m, r string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if strings.EqualFold(r, route) && strings.EqualFold(m, method) {
			found = true
		}
		return nil
	})

	return found
}
