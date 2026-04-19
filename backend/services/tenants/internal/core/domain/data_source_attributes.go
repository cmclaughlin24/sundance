package domain

type DataSourceAttributes interface {
	isDataSourceAttributes()
}

type StaticDataSourceAttributes struct {
	Data any
}

func (a StaticDataSourceAttributes) isDataSourceAttributes() {}

type ScheduledDataSourceAttributes struct {
}

func (a ScheduledDataSourceAttributes) isDataSourceAttributes() {}

type QueryDataSourceType string

const (
	QueryDataSourceTypeREST QueryDataSourceType = "rest"
	QueryDataSourceTypeGRPC QueryDataSourceType = "grpc"
)

type QueryDataSourceAttributes struct {
	Type     QueryDataSourceType
	Endpoint string
}

func (a QueryDataSourceAttributes) isDataSourceAttributes() {}
