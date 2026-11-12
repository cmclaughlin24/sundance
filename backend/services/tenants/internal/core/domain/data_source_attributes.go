package domain

import (
	"errors"
	"time"
)

var (
	ErrDataSourceAttributeMismatch = errors.New("data source type and attributes mismatch")
	ErrMissingRequiredKeys         = errors.New("missing required keys in parameters")
)

type DataSourceAttributes interface {
	isDataSourceAttributes()
}

type DataSourceHTTPRequest struct {
	URL        string            `desc:"Endpoint URL to request lookup options from."`
	Method     string            `desc:"HTTP method used for the request (e.g. GET, POST)."`
	Headers    map[string]string `desc:"HTTP headers sent with the request, such as authorization."`
	ValueField string            `desc:"Field in each response item to use as the option value."`
	LabelField string            `desc:"Field in each response item to use as the option label."`
}

type baseDataSourceAttributes struct{}

func (b baseDataSourceAttributes) isDataSourceAttributes() {}

type StaticDataSourceAttributes struct {
	baseDataSourceAttributes
	Data []*Lookup `desc:"List of value and label options returned by the data source."`
}

type ScheduledDataSourceAttributes struct {
	baseDataSourceAttributes
	DataSourceHTTPRequest
	Data           []*Lookup `desc:"-"`
	IntervalHours  float64   `desc:"Number of hours between data refreshes."`
	ExpirationDate time.Time `desc:"-"`
	Attempts       int       `desc:"-"`
}

func (attr *ScheduledDataSourceAttributes) RecordAttempt() {
	attr.Attempts += 1
}

func (attr *ScheduledDataSourceAttributes) RefreshData(data []*Lookup) {
	attr.Data = data
	attr.ExpirationDate = Now().Add(time.Duration(attr.IntervalHours * float64(time.Hour)))
	attr.Attempts = 0
}

type WebhookDataSourceAttributes struct {
	baseDataSourceAttributes
	DataSourceHTTPRequest
	RequiredKeys []string `desc:"Parameter keys that must be provided on each lookup request."`
}

type DataLakeDataSourceAttributes struct {
	baseDataSourceAttributes
	Query        string   `desc:"Parameterized SQL query used to resolve lookup options."`
	RequiredKeys []string `desc:"Query parameter keys that must be provided on each lookup request."`
	OptionalKeys []string `desc:"Query parameter keys that may optionally be provided on each lookup request."`
	Catalog      string   `desc:"Data lake catalog to run the query against."`
	Schema       string   `desc:"Data lake schema to run the query against."`
	ValueField   string   `desc:"Column in the query result to use as the option value."`
	LabelField   string   `desc:"Column in the query result to use as the option label."`
	TimeoutMs    int      `desc:"Maximum query execution time in milliseconds."`
}

func GetDataSourceAttributes[T DataSourceAttributes](attr DataSourceAttributes) (T, error) {
	switch t := attr.(type) {
	case T:
		return t, nil
	default:
		return *new(T), ErrDataSourceAttributeMismatch
	}
}
