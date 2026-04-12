package ports

import (
	"context"

	"github.com/cmclaughlin24/sundance/forms/internal/core/domain"
)

type Services struct {
	Forms FormsService
}

type FormsService interface {
	Find(context.Context) ([]*domain.Form, error)
	FindById(context.Context, domain.FormID) (*domain.Form, error)
	Create(context.Context, *CreateFormCommand) (*domain.Form, error)
	Update(context.Context, *UpdateFormCommand) (*domain.Form, error)
	CreateVersion(context.Context, *CreateVersionCommand) (*domain.Version, error)
	UpdateVersion(context.Context, *UpdateVersionCommand) (*domain.Version, error)
	PublishVersion(context.Context, *PublishVersionCommand) (*domain.Version, error)
	RetireVersion(context.Context, *RetireVersionCommand) (*domain.Version, error)
}
