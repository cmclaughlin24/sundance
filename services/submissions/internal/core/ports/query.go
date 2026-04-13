package ports

type FindByIdQuery[T any] struct {
	ID       T
	TenantID string
}
