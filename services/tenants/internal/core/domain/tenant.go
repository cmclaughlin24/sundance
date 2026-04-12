package domain

import "time"

type TenantID string

type Tenant struct {
	ID          TenantID
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DataSources []*DataSource
}

func NewTenant(name, description string) *Tenant {
	return &Tenant{
		Name:        name,
		Description: description,
	}
}
