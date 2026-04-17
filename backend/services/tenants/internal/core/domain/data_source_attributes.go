package domain

type DataSourceAttributes any

type StaticDataSourceAttributes struct {
	Data any
}

type ScheduledDataSourceAttributes struct {
}

type QueryDataSourceType string

const (
	QueryDataSourceTypeREST QueryDataSourceType = "rest"
	QueryDataSourceTypeGRPC QueryDataSourceType = "grpc"
)

type QueryDataSourceAttributes struct {
	Type     QueryDataSourceType
	Endpoint string
}
