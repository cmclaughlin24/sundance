package ports

type FindByIdQuery[T any] struct {
	TenantID string `validate:"required"`
	ID       T      `validate:"required"`
}

func NewFindByIdQuery[T any](tenantID string, id T) *FindByIdQuery[T] {
	query := &FindByIdQuery[T]{
		TenantID: tenantID,
		ID:       id,
	}

	return query
}
