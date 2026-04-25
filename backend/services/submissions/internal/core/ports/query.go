package ports

type FindByIDQuery[T any] struct {
	TenantID string `validate:"required"`
	ID       T      `validate:"required"`
}

func NewFindByIDQuery[T any](tenantID string, id T) *FindByIDQuery[T] {
	query := &FindByIDQuery[T]{
		TenantID: tenantID,
		ID:       id,
	}

	return query
}
