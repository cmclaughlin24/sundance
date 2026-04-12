package domain

import "time"

type FormID string

type Form struct {
	ID          FormID
	Name        string
	Description string
	TenantID    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
