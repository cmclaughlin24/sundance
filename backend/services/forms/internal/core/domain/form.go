package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrInvalidForm = errors.New("invalid form")

type FormID string

type Form struct {
	ID          FormID
	TenantID    string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewForm(tenantID, name, description string) (*Form, error) {
	f := &Form{
		ID:          FormID(uuid.NewString()),
		TenantID:    tenantID,
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
	}

	// TODO: Implement domain specific validation.

	return f, nil
}

func HydrateForm(id FormID, tenantID, name, description string, createdAt, updatedAt time.Time) *Form {
	return &Form{
		ID:          id,
		TenantID:    tenantID,
		Name:        name,
		Description: description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func (f *Form) Update(name, description string) error {
	if f == nil {
		return ErrInvalidForm
	}

	f.Name = name
	f.Description = description
	f.UpdatedAt = time.Now()

	// TODO: Implement domain specific validation.

	return nil
}
