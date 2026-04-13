package rest

import (
	"net/http"

	"github.com/cmclaughlin24/sundance/common"
	"github.com/cmclaughlin24/sundance/forms/internal/core"
	"github.com/cmclaughlin24/sundance/forms/internal/core/domain"
	"github.com/cmclaughlin24/sundance/forms/internal/core/ports"
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

func (h *handlers) getForms(w http.ResponseWriter, r *http.Request) {
	resultChan := make(chan result[[]*domain.Form], 1)

	go func() {
		defer close(resultChan)
		forms, err := h.app.Services.Forms.Find(r.Context())
		resultChan <- result[[]*domain.Form]{forms, err}
	}()

	select {
	case <-r.Context().Done():
		return
	case res := <-resultChan:
		if res.err != nil {
			common.SendErrorResponse(w, res.err)
			return
		}

		// TODO: Send response.
	}
}

func (h *handlers) getForm(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("formId")
	tenantId := r.URL.Query().Get("tenantId")
	query := ports.FindByIdQuery{
		ID:       domain.FormID(id),
		TenantID: tenantId,
	}
	resultChan := make(chan result[*domain.Form], 1)

	go func() {
		defer close(resultChan)
		form, err := h.app.Services.Forms.FindById(r.Context(), query)
		resultChan <- result[*domain.Form]{form, err}
	}()

	select {
	case <-r.Context().Done():
		return
	case res := <-resultChan:
		if res.err != nil {
			common.SendErrorResponse(w, res.err)
			return
		}

		// TODO: Send response.
	}
}

func (h *handlers) createForm(w http.ResponseWriter, r *http.Request) {
	resultChan := make(chan result[*domain.Form], 1)

	var dto upsertFormDto
	if err := common.ReadJsonPayload(r, &dto); err != nil {
		return
	}

	command, err := ports.NewCreateFormCommand(dto.TenantID, dto.Name, dto.Description)
	if err != nil {
		common.SendErrorResponse(w, err)
		return
	}

	go func() {
		defer close(resultChan)
		form, err := h.app.Services.Forms.Create(r.Context(), command)
		resultChan <- result[*domain.Form]{form, err}
	}()

	select {
	case <-r.Context().Done():
		return
	case res := <-resultChan:
		if res.err != nil {
			common.SendErrorResponse(w, res.err)
			return
		}

		// TODO: Send response.
	}
}

func (h *handlers) updateForm(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("formId")
	resultChan := make(chan result[*domain.Form], 1)

	var dto upsertFormDto
	if err := common.ReadJsonPayload(r, &dto); err != nil {
		return
	}

	command, err := ports.NewUpdateFormCommand(domain.FormID(id), dto.TenantID, dto.Name, dto.Description)
	if err != nil {
		common.SendErrorResponse(w, err)
		return
	}

	go func() {
		defer close(resultChan)
		form, err := h.app.Services.Forms.Update(r.Context(), command)
		resultChan <- result[*domain.Form]{form, err}
	}()

	select {
	case <-r.Context().Done():
		return
	case res := <-resultChan:
		if res.err != nil {
			common.SendErrorResponse(w, res.err)
			return
		}

		// TODO: Send response.
	}
}

func (h *handlers) getVersions(w http.ResponseWriter, r *http.Request) {}

func (h *handlers) getVersion(w http.ResponseWriter, r *http.Request) {}

func (h *handlers) createVersion(w http.ResponseWriter, r *http.Request) {}

func (h *handlers) updateVersion(w http.ResponseWriter, r *http.Request) {}

func (h *handlers) removeVersion(w http.ResponseWriter, r *http.Request) {
}

func (h *handlers) publishVersion(w http.ResponseWriter, r *http.Request) {
	formId, versionId := h.getVersionPathValues(r)
	resultChan := make(chan result[*domain.Version], 1)

	command, err := ports.NewPublishVersionCommand(formId, "", versionId, "")
	if err != nil {
		common.SendErrorResponse(w, err)
		return
	}

	go func() {
		defer close(resultChan)
		version, err := h.app.Services.Forms.PublishVersion(r.Context(), command)
		resultChan <- result[*domain.Version]{version, err}
	}()

	select {
	case <-r.Context().Done():
		return
	case res := <-resultChan:
		if res.err != nil {
			common.SendErrorResponse(w, res.err)
			return
		}

		// TODO: Send response.
	}
}

func (h *handlers) retireVersion(w http.ResponseWriter, r *http.Request) {
	formId, versionId := h.getVersionPathValues(r)
	resultChan := make(chan result[*domain.Version], 1)

	command, err := ports.NewRetireVersionCommand(formId, "", versionId, "")
	if err != nil {
		common.SendErrorResponse(w, err)
		return
	}

	go func() {
		defer close(resultChan)
		version, err := h.app.Services.Forms.RetireVersion(r.Context(), command)
		resultChan <- result[*domain.Version]{version, err}
	}()

	select {
	case <-r.Context().Done():
		return
	case res := <-resultChan:
		if res.err != nil {
			common.SendErrorResponse(w, res.err)
			return
		}

		// TODO: Send response.
	}
}

func (h *handlers) getVersionPathValues(r *http.Request) (domain.FormID, domain.VersionID) {
	formId := r.PathValue("formId")
	versionId := r.PathValue("versionId")

	return domain.FormID(formId), domain.VersionID(versionId)
}
