package rest

import "github.com/cmclaughlin24/sundance/submissions/internal/core"

type handlers struct {
	app *core.Application
}

func newHandlers(app *core.Application) *handlers {
	return &handlers{
		app: app,
	}
}
