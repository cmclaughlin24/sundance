package domain

import (
	"errors"
	"time"
)

type DataSourceID string

type DataSourceType string

const (
	DataSourceTypeStatic    DataSourceType = "static"
	DataSourceTypeScheduled DataSourceType = "scheduled"
	DataSourceTypeQuery     DataSourceType = "query"
)

var (
	ErrInvalidSourceTypeAttributes = errors.New("invalid data source attributes for type")
)

type DataSource struct {
	ID         DataSourceID
	TenantID   TenantID
	Type       DataSourceType
	Attributes DataSourceAttributes
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewDataSource(id DataSourceID, tenantID TenantID, sourceType DataSourceType, attributes DataSourceAttributes) (*DataSource, error) {
	if !isValidAttributeType(sourceType, attributes) {
		return nil, ErrInvalidSourceTypeAttributes
	}

	return &DataSource{
		ID:         id,
		TenantID:   tenantID,
		Type:       sourceType,
		Attributes: attributes,
	}, nil
}

func isValidAttributeType(sourceType DataSourceType, attributes DataSourceAttributes) bool {
	switch attributes.(type) {
	case StaticDataSourceAttributes:
		return sourceType == DataSourceTypeStatic
	case ScheduledDataSourceAttributes:
		return sourceType == DataSourceTypeScheduled
	case QueryDataSourceAttributes:
		return sourceType == DataSourceTypeQuery
	default:
		return false
	}
}
