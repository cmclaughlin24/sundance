package mongodb

import (
	"time"

	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
)

type formDocument struct {
	ID          string    `bson:"_id"`
	TenantID    string    `bson:"tenant_id"`
	Name        string    `bson:"name"`
	Description string    `bson:"description"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}

func toFormDocument(f *domain.Form) *formDocument {
	return &formDocument{
		ID:          string(f.ID),
		TenantID:    f.TenantID,
		Name:        f.Name,
		Description: f.Description,
		CreatedAt:   f.CreatedAt,
		UpdatedAt:   f.UpdatedAt,
	}
}

func fromFormDocument(f *formDocument) *domain.Form {
	return domain.HydrateForm(
		domain.FormID(f.ID),
		f.TenantID,
		f.Name,
		f.Description,
		f.CreatedAt,
		f.UpdatedAt,
	)
}
