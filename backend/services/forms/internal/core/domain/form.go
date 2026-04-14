package domain

import (
	"errors"
	"time"
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

func NewForm(id FormID, tenantID, name, description string) (*Form, error) {
	f := &Form{
		ID:          id,
		TenantID:    tenantID,
		Name:        name,
		Description: description,
	}

	// TODO: Implement domain specific validation.

	return f, nil
}

func (f *Form) Update(name, description string) error {
	if f == nil {
		return ErrInvalidForm
	}

	f.Name = name
	f.Description = description
	
	// TODO: Implement domain specific validation.

	return nil
}
