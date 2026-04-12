package ports

import (
	"context"

	"github.com/cmclaughlin24/sundance/common/database"
	"github.com/cmclaughlin24/sundance/forms/internal/core/domain"
)

type Repository struct {
	Database database.Database
	Forms    FormsRepository
}

type FormsRepository interface {
	Find(context.Context) ([]*domain.Form, error)
	FindById(context.Context, domain.FormID) (*domain.Form, error)
	Create(context.Context, *domain.Form) (*domain.Form, error)
	Update(context.Context, *domain.Form) (*domain.Form, error)
	GetVersions(context.Context, domain.FormID) ([]*domain.Version, error)
	GetVersion(context.Context, domain.FormID, domain.VersionID) (*domain.Version, error)
	CreateVersion(context.Context, *domain.Version) (*domain.Version, error)
	UpdateVersion(context.Context, *domain.Version) (*domain.Version, error)
}
