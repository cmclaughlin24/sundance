package ports

type FindByIdQuery[T any] struct {
	ID T `validate:"required"`
}

func NewFindByIdQuery[T any](id T) *FindByIdQuery[T] {
	query := &FindByIdQuery[T]{
		ID: id,
	}

	return query
}
