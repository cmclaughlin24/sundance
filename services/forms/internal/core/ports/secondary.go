package ports

import (
	"context"

	"github.com/cmclaughlin24/sundance/forms/internal/core/domain"
)

type Repository struct {
	Database Database
	Forms    FormsRepository
}

type Database interface {
	Close() error
	BeginTx(context.Context) (context.Context, error)
	GetTx(context.Context) (any, error)
	CommitTx(context.Context) error
	RollbackTx(context.Context) error
}

type FormsRepository interface {
	Find(context.Context) ([]*domain.Form, error)
	FindById(context.Context, domain.FormID) (*domain.Form, error)
	Create(context.Context, *domain.Form) (*domain.Form, error)
	Update(context.Context, *domain.Form) (*domain.Form, error)
	GetVersions(context.Context, domain.FormID) ([]*domain.Version, error)
	GetVersion(context.Context, domain.FormID, domain.VersionID) ([]*domain.Version, error)
	CreateVersion(context.Context, *CreateVersionCommand) (*domain.Version, error)
	UpdateVersion(context.Context, *UpdateVersionCommand) (*domain.Version, error)
}
