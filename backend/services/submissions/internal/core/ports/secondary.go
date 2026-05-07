package ports

import (
	"context"

	"github.com/cmclaughlin24/sundance/backend/pkg/database"
	"github.com/cmclaughlin24/sundance/backend/services/submissions/internal/core/domain"
)

type Repository struct {
	Database    database.Database
	Submissions SubmissionsRepository
}

type SubmissionsRepository interface {
	Find(context.Context, *FindSubmissionsFilter) ([]*domain.Submission, error)
	FindByID(context.Context, domain.SubmissionID) (*domain.Submission, error)
	FindByReferenceID(context.Context, domain.ReferenceID) (*domain.Submission, error)
	FindByIdempotencyID(context.Context, domain.IdempotencyID) (*domain.Submission, error)
	Upsert(context.Context, *domain.Submission) (*domain.Submission, error)
}
