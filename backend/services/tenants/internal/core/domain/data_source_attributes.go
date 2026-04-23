package domain

type DataSourceAttributes interface {
	isDataSourceAttributes()
}

type baseDataSourceAttributes struct{}

func (b baseDataSourceAttributes) isDataSourceAttributes() {}

type StaticDataSourceAttributes struct {
	baseDataSourceAttributes
	Data []*Lookup
}

type ScheduledDataSourceAttributes struct {
	baseDataSourceAttributes
	Data []*Lookup
}

type QueryDataSourceType string

const (
	QueryDataSourceTypeREST QueryDataSourceType = "rest"
	QueryDataSourceTypeGRPC QueryDataSourceType = "grpc"
)

type QueryDataSourceAttributes struct {
	baseDataSourceAttributes
	Type     QueryDataSourceType
	Resource string
}
