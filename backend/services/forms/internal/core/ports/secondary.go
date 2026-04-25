package ports

import (
	"context"

	"github.com/cmclaughlin24/sundance/backend/pkg/database"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
)

type Repository struct {
	Database database.Database
	Forms    FormsRepository
}

type FormsRepository interface {
	Find(context.Context, *FormFilters) ([]*domain.Form, error)
	FindByID(context.Context, domain.FormID) (*domain.Form, error)
	Upsert(context.Context, *domain.Form) (*domain.Form, error)
	FindVersions(context.Context, domain.FormID) ([]*domain.Version, error)
	FindVersion(context.Context, domain.FormID, domain.VersionID) (*domain.Version, error)
	FindNextVersionNumber(context.Context, domain.FormID) (int, error)
	UpsertVersion(context.Context, *domain.Version) (*domain.Version, error)
}
