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

type WebhookDataSourceAttributes struct {
	baseDataSourceAttributes
	URL     string
	Method  string
	Headers map[string]string
}
