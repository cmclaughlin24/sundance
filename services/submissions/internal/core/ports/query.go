package ports

import "github.com/go-playground/validator/v10"

type FindByIdQuery[T any] struct {
	ID       T
	TenantID string
}

func NewFindByIdQuery[T any](id T, tenantID string) (*FindByIdQuery[T], error) {
	query := &FindByIdQuery[T]{
		ID:       id,
		TenantID: tenantID,
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(query); err != nil {
		return nil, err
	}

	return query, nil
}
