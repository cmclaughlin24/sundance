package ports

import (
	"context"

	"github.com/cmclaughlin24/sundance/backend/pkg/database"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
)

type Repository struct {
	Database database.Database
	Forms    FormsRepository
	Versions VersionRepository
}

type FormsRepository interface {
	Find(context.Context, *FormFilters) ([]*domain.Form, error)
	FindByID(context.Context, domain.FormID) (*domain.Form, error)
	Upsert(context.Context, *domain.Form) (*domain.Form, error)
}

type VersionRepository interface {
	Find(context.Context, domain.FormID) ([]*domain.Version, error)
	FindByID(context.Context, domain.FormID, domain.VersionID) (*domain.Version, error)
	FindNextVersionNumber(context.Context, domain.FormID) (int, error)
	Upsert(context.Context, *domain.Version) (*domain.Version, error)
}
