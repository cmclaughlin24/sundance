package ports

type FindSubmissionsQuery struct {
	TenantID string `validate:"required"`
}

func NewFindSubmissionsQuery(tenantID string) *FindSubmissionsQuery {
	return &FindSubmissionsQuery{
		TenantID: tenantID,
	}
}

type FindSubmissionByIDQuery[T any] struct {
	TenantID string `validate:"required"`
	ID       T      `validate:"required"`
}

func NewFindSubmissionByIDQuery[T any](tenantID string, id T) *FindSubmissionByIDQuery[T] {
	query := &FindSubmissionByIDQuery[T]{
		TenantID: tenantID,
		ID:       id,
	}

	return query
}

type FindSubmissionsFilter struct {
	TenantID string 
}
