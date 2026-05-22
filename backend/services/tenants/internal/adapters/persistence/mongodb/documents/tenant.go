package documents

import (
	"sundance/backend/services/tenants/internal/core/domain"
	"time"
)

type TenantDocument struct {
	ID          string    `bson:"_id"`
	Name        string    `bson:"name"`
	Description string    `bson:"description"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}

func ToTenantDocument(t *domain.Tenant) *TenantDocument {
	return &TenantDocument{
		ID:          string(t.ID),
		Name:        t.Name,
		Description: t.Description,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

func FromTenantDocument(t *TenantDocument) *domain.Tenant {
	return domain.HydrateTenant(
		domain.TenantID(t.ID),
		t.Name,
		t.Description,
		t.CreatedAt,
		t.UpdatedAt,
	)
}
