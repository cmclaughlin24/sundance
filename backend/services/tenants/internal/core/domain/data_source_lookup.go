package domain

type DataSourceLookup struct {
	Code        string
	Description string
}

func NewDataSourceLookup(code, description string) *DataSourceLookup {
	return &DataSourceLookup{
		Code:        code,
		Description: description,
	}
}
