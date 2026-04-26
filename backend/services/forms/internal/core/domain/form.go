package domain

import (
	"errors"
	"time"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/validate"
	"github.com/google/uuid"
)

var (
	ErrInvalidForm          = errors.New("invalid form")
	ErrFormHasActiveVersion = errors.New("form has at least one active version")
)

type FormID string

type Form struct {
	ID          FormID
	TenantID    string `validate:"required,notblank"`
	Name        string `validate:"required,notblank"`
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
		CreatedAt:   Now(),
	}

	if err := validate.ValidateStruct(f); err != nil {
		return nil, err
	}

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
	f.UpdatedAt = Now()

	if err := validate.ValidateStruct(f); err != nil {
		return err
	}

	return nil
}
