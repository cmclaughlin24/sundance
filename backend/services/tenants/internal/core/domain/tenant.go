package domain

import (
	"errors"
	"time"

	"sundance/backend/pkg/common/validate"
)

var ErrInvalidTenant = errors.New("invalid tenant")

type TenantID string

type Tenant struct {
	ID          TenantID
	Name        string `validate:"required"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewTenant(name, description string) (*Tenant, error) {
	t := &Tenant{
		ID:          TenantID(NewID()),
		Name:        name,
		Description: description,
		CreatedAt:   Now(),
	}

	if err := validate.ValidateStruct(t); err != nil {
		return nil, err
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

	cpy := *t
	cpy.Name = name
	cpy.Description = description

	if err := validate.ValidateStruct(cpy); err != nil {
		return err
	}

	*t = cpy
	t.UpdatedAt = Now()

	return nil
}
