package domain

import "time"

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
