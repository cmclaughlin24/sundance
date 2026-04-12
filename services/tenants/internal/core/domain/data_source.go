package domain

type DataSourceID string

type DataSourceType string

const (
	DataSourceTypeStatic    DataSourceType = "static"
	DataSourceTypeScheduled DataSourceType = "scheduled"
	DataSourceTypeQuery     DataSourceType = "query"
)

type DataSource struct {
	ID         DataSourceID
	Type       DataSourceType
	Attributes DataSourceAttributes
}
