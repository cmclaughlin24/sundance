package domain

import "time"

type DataSourceID string

type DataSourceType string

const (
	DataSourceTypeStatic    DataSourceType = "static"
	DataSourceTypeScheduled DataSourceType = "scheduled"
	DataSourceTypeQuery     DataSourceType = "query"
)

type DataSource struct {
	ID         DataSourceID
	TenantID   TenantID
	Type       DataSourceType
	Attributes DataSourceAttributes
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewDataSource(id DataSourceID, tenantID TenantID, sourceType DataSourceType, attributes DataSourceAttributes) *DataSource {
	return &DataSource{
		ID:         id,
		TenantID:   tenantID,
		Type:       sourceType,
		Attributes: attributes,
	}
}
