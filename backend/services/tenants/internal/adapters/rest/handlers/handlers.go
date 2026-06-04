package handlers

import (
	"errors"
	"net/http"

	"sundance/backend/pkg/common/httputil"
	"sundance/backend/services/tenants/internal/adapters/rest/dto"
	"sundance/backend/services/tenants/internal/core"
	"sundance/backend/services/tenants/internal/core/domain"
)

type result[T any] struct {
	data T
	err  error
}

type Handlers struct {
	app *core.Application
}

func NewHandlers(app *core.Application) *Handlers {
	return &Handlers{
		app: app,
	}
}

func (h *Handlers) sendErrorResponse(w http.ResponseWriter, err error) {
	switch {
	case isBadRequest(err):
		httputil.SendJSONResponse(w, http.StatusBadRequest, httputil.APIErrorResponse{
			Message:    "Bad Request",
			Error:      err.Error(),
			StatusCode: http.StatusBadRequest,
		})

	default:
		httputil.SendErrorResponse(w, err)
	}
}

func isBadRequest(err error) bool {
	return errors.Is(err, dto.ErrDataSourceAttrParse) ||
		errors.Is(err, domain.ErrInvalidSourceType) ||
		errors.Is(err, domain.ErrInvalidSourceTypeAttributes) ||
		errors.Is(err, domain.ErrMissingRequiredKeys)

}
