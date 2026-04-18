package ports

import (
	"context"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/database"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/domain"
)

type Repository struct {
	Database    database.Database
	Submissions SubmissionsRepository
}

type SubmissionsRepository interface {
	Find(context.Context) ([]*domain.Submission, error)
	FindById(context.Context, domain.SubmissionID) (*domain.Submission, error)
	FindByReferenceId(context.Context, domain.ReferenceID) (*domain.Submission, error)
}
