package domain

import (
	"errors"
	"time"
)

var (
	ErrDataSourceAttributeMismatch = errors.New("data source type and attributes mismatch")
)

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
	Data           []*Lookup
	URL            string
	Method         string
	Headers        map[string]string
	IntervalHours  float64
	ExpirationDate time.Time
}

func (attr *ScheduledDataSourceAttributes) RefreshData(data []*Lookup) {
	attr.Data = data
	attr.ExpirationDate = Now().Add(time.Duration(attr.IntervalHours * float64(time.Hour)))
}

type WebhookDataSourceAttributes struct {
	baseDataSourceAttributes
	URL     string
	Method  string
	Headers map[string]string
}

func GetDataSourceAttributes[T DataSourceAttributes](attr DataSourceAttributes) (T, error) {
	switch t := attr.(type) {
	case T:
		return t, nil
	default:
		return *new(T), ErrDataSourceAttributeMismatch
	}
}
