package domain

import (
	"time"
)

type TenantID string

type Tenant struct {
	ID          TenantID
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewTenant(id TenantID, name, description string) (*Tenant, error) {
	return &Tenant{
		ID:          id,
		Name:        name,
		Description: description,
	}, nil
}
