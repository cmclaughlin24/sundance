package documents

import (
	"time"

	"sundance/backend/pkg/common/stratreg"
	"sundance/backend/services/tenants/internal/core/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type DataSourceDocument struct {
	ID          string    `bson:"_id"`
	TenantID    string    `bson:"tenant_id"`
	Name        string    `bson:"name"`
	Description string    `bson:"description"`
	Type        string    `bson:"type"`
	Attributes  bson.Raw  `bson:"attributes"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}

func ToDataSourceDocument(ds *domain.DataSource) (*DataSourceDocument, error) {
	attr, err := bson.Marshal(ds.Attributes)

	if err != nil {
		return nil, err
	}

	return &DataSourceDocument{
		ID:          string(ds.ID),
		TenantID:    string(ds.TenantID),
		Name:        ds.Name,
		Description: ds.Description,
		Type:        string(ds.Type),
		Attributes:  attr,
		CreatedAt:   ds.CreatedAt,
		UpdatedAt:   ds.UpdatedAt,
	}, nil
}

func FromDataSourceDocument(ds *DataSourceDocument) (*domain.DataSource, error) {
	sourceType := domain.DataSourceType(ds.Type)
	attr, err := unmarshalDataSourceAttributes(sourceType, ds.Attributes)

	if err != nil {
		return nil, err
	}

	return domain.HydrateDataSource(
		domain.DataSourceID(ds.ID),
		domain.TenantID(ds.TenantID),
		ds.Name,
		ds.Description,
		sourceType,
		attr,
		ds.CreatedAt,
		ds.UpdatedAt,
	), nil
}

type attributeParser func(bson.Raw) (domain.DataSourceAttributes, error)

var attributeParserStrategies = stratreg.New[domain.DataSourceType, attributeParser]().
	Set(domain.DataSourceTypeStatic, func(raw bson.Raw) (domain.DataSourceAttributes, error) {
		return parseDataSourceAttributes[domain.StaticDataSourceAttributes](raw)
	}).
	Set(domain.DataSourceTypeScheduled, func(raw bson.Raw) (domain.DataSourceAttributes, error) {
		return parseDataSourceAttributes[domain.ScheduledDataSourceAttributes](raw)
	}).
	Set(domain.DataSourceTypeWebhook, func(raw bson.Raw) (domain.DataSourceAttributes, error) {
		return parseDataSourceAttributes[domain.WebhookDataSourceAttributes](raw)
	}).
	Set(domain.DataSourceTypeDataLake, func(raw bson.Raw) (domain.DataSourceAttributes, error) {
		return parseDataSourceAttributes[domain.DataLakeDataSourceAttributes](raw)
	})

func unmarshalDataSourceAttributes(sourceType domain.DataSourceType, raw bson.Raw) (domain.DataSourceAttributes, error) {
	strategy, err := attributeParserStrategies.Get(sourceType)

	if err != nil {
		return nil, err
	}

	return strategy(raw)
}

func parseDataSourceAttributes[T domain.DataSourceAttributes](raw bson.Raw) (domain.DataSourceAttributes, error) {
	var attr T

	if err := bson.Unmarshal(raw, &attr); err != nil {
		return nil, err
	}

	return attr, nil
}
