package domain

type DataSourceAttributes interface {
	isDataSourceAttributes()
}

type baseDataSourceAttributes struct{}

func (b baseDataSourceAttributes) isDataSourceAttributes() {}

type StaticDataSourceAttributes struct {
	baseDataSourceAttributes
	Data []DataSourceLookup
}

type ScheduledDataSourceAttributes struct {
	baseDataSourceAttributes
	Data []DataSourceLookup
}

type QueryDataSourceType string

const (
	QueryDataSourceTypeREST QueryDataSourceType = "rest"
	QueryDataSourceTypeGRPC QueryDataSourceType = "grpc"
)

type QueryDataSourceAttributes struct {
	baseDataSourceAttributes
	Type     QueryDataSourceType
	Endpoint string
}
