package documents

import (
	"sundance/backend/services/forms/internal/core/domain"
	"time"
)

type FormDocument struct {
	ID          string    `bson:"_id"`
	TenantID    string    `bson:"tenant_id"`
	Name        string    `bson:"name"`
	Description string    `bson:"description"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}

func ToFormDocument(f *domain.Form) *FormDocument {
	return &FormDocument{
		ID:          string(f.ID),
		TenantID:    f.TenantID,
		Name:        f.Name,
		Description: f.Description,
		CreatedAt:   f.CreatedAt,
		UpdatedAt:   f.UpdatedAt,
	}
}

func FromFormDocument(f *FormDocument) *domain.Form {
	return domain.HydrateForm(
		domain.FormID(f.ID),
		f.TenantID,
		f.Name,
		f.Description,
		f.CreatedAt,
		f.UpdatedAt,
	)
}
