package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrInvalidTenant = errors.New("invalid tenant")

type TenantID string

type Tenant struct {
	ID          TenantID
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewTenant(name, description string) (*Tenant, error) {
	t := &Tenant{
		ID:          TenantID(uuid.NewString()),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
	}

	return t, nil
}

func HydrateTenant(id TenantID, name, description string, createdAt, updatedAt time.Time) *Tenant {
	t := &Tenant{
		ID:          id,
		Name:        name,
		Description: description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	return t
}

func (t *Tenant) Update(name, description string) error {
	if t == nil {
		return ErrInvalidTenant
	}

	t.Name = name
	t.Description = description
	t.UpdatedAt = time.Now()

	return nil
}
