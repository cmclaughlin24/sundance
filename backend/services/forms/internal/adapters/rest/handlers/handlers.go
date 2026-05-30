package handlers

import (
	"errors"
	"net/http"

	"sundance/backend/pkg/common/httputil"
	"sundance/backend/services/forms/internal/core"
	"sundance/backend/services/forms/internal/core/domain"
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
	return errors.Is(err, domain.ErrVersionLocked) ||
		errors.Is(err, domain.ErrInvalidVersion) ||
		errors.Is(err, domain.ErrInvalidVersionStatus) ||
		errors.Is(err, domain.ErrDuplicateVersion) ||
		errors.Is(err, domain.ErrInvalidPosition) ||
		errors.Is(err, domain.ErrDuplicatePosition) ||
		errors.Is(err, domain.ErrInvalidRuleType) ||
		errors.Is(err, domain.ErrDuplicateRuleType) ||
		errors.Is(err, domain.ErrPublishedByRequired) ||
		errors.Is(err, domain.ErrRetiredByRequired) ||
		errors.Is(err, domain.ErrInvalidFieldType) ||
		errors.Is(err, domain.ErrInvalidFieldAttributes) ||
		errors.Is(err, domain.ErrInvalidForm) ||
		errors.Is(err, domain.ErrFormHasActiveVersion) ||
		errors.Is(err, domain.ErrInvalidPage) ||
		errors.Is(err, domain.ErrInvalidSection) ||
		errors.Is(err, domain.ErrInvalidExprOperator) ||
		errors.Is(err, domain.ErrInvalidJoinOperator)
}
