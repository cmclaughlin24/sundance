package ports

type ListDataSourceQuery struct {
	// TODO: Add pagination support through embedded struct.
}

func NewListDataSourceQuery() *ListDataSourceQuery {
	return &ListDataSourceQuery{}
}
