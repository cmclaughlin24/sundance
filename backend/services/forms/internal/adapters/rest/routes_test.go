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
		{"/api/v1/forms/", "GET"},
		{"/api/v1/forms/", "POST"},
		{"/api/v1/forms/{formId}/", "GET"},
		{"/api/v1/forms/{formId}/", "PUT"},
		{"/api/v1/forms/{formId}/", "DELETE"},
		{"/api/v1/forms/{formId}/versions/", "GET"},
		{"/api/v1/forms/{formId}/versions/", "POST"},
		{"/api/v1/forms/{formId}/versions/{versionId}/", "GET"},
		{"/api/v1/forms/{formId}/versions/{versionId}/", "PUT"},
		{"/api/v1/forms/{formId}/versions/{versionId}/publish", "POST"},
		{"/api/v1/forms/{formId}/versions/{versionId}/retire", "POST"},
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
