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

type DataSourceRequest struct {
	URL        string
	Method     string
	Headers    map[string]string
	ValueField string
	LabelField string
}

type baseDataSourceAttributes struct{}

func (b baseDataSourceAttributes) isDataSourceAttributes() {}

type StaticDataSourceAttributes struct {
	baseDataSourceAttributes
	Data []*Lookup
}

type ScheduledDataSourceAttributes struct {
	baseDataSourceAttributes
	DataSourceRequest
	Data           []*Lookup
	IntervalHours  float64
	ExpirationDate time.Time
	Attempts       int
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
	DataSourceRequest
	RequiredKeys []string
}

type DataLakeDataSourceAttributes struct {
	baseDataSourceAttributes
	Query        string
	RequiredKeys []string
	OptionalKeys []string
	Catalog      string
	Schema       string
	ValueField   string
	LabelField   string
	TimeoutMs    int
}

func GetDataSourceAttributes[T DataSourceAttributes](attr DataSourceAttributes) (T, error) {
	switch t := attr.(type) {
	case T:
		return t, nil
	default:
		return *new(T), ErrDataSourceAttributeMismatch
	}
}
