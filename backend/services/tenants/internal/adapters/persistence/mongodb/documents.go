package mongodb

import (
	"time"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/strategy"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type tenantDocument struct {
	ID          string    `bson:"_id,omitempty"`
	Name        string    `bson:"name"`
	Description string    `bson:"description"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}

func toTenantDocument(t *domain.Tenant) *tenantDocument {
	return &tenantDocument{
		ID:          string(t.ID),
		Name:        t.Name,
		Description: t.Description,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

func fromTenantDocument(t *tenantDocument) *domain.Tenant {
	return &domain.Tenant{
		ID:          domain.TenantID(t.ID),
		Name:        t.Name,
		Description: t.Description,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

type dataSourceDocument struct {
	ID          string    `bson:"_id,omitempty"`
	TenantID    string    `bson:"tenant_id"`
	Name        string    `bson:"name"`
	Description string    `bson:"description"`
	Type        string    `bson:"type"`
	Attributes  bson.Raw  `bson:"attributes"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}

func toDataSourceDocument(ds *domain.DataSource) (*dataSourceDocument, error) {
	attr, err := bson.Marshal(ds.Attributes)

	if err != nil {
		return nil, err
	}

	return &dataSourceDocument{
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

func fromDataSourceDocument(ds *dataSourceDocument) (*domain.DataSource, error) {
	sourceType := domain.DataSourceType(ds.Type)
	attr, err := unmarshalDataSourceAttributes(sourceType, ds.Attributes)

	if err != nil {
		return nil, err
	}

	return &domain.DataSource{
		ID:          domain.DataSourceID(ds.ID),
		TenantID:    domain.TenantID(ds.TenantID),
		Name:        ds.Name,
		Description: ds.Description,
		Type:        sourceType,
		Attributes:  attr,
		CreatedAt:   ds.CreatedAt,
		UpdatedAt:   ds.UpdatedAt,
	}, nil
}

type attributeParser func(bson.Raw) (domain.DataSourceAttributes, error)

var attributeParserStrategies = strategy.NewStrategies[domain.DataSourceType, attributeParser]().
	Set(domain.DataSourceTypeStatic, func(raw bson.Raw) (domain.DataSourceAttributes, error) {
		return parseDataSourceAttributes[domain.StaticDataSourceAttributes](raw)
	}).
	Set(domain.DataSourceTypeScheduled, func(raw bson.Raw) (domain.DataSourceAttributes, error) {
		return parseDataSourceAttributes[domain.ScheduledDataSourceAttributes](raw)
	}).
	Set(domain.DataSourceTypeQuery, func(raw bson.Raw) (domain.DataSourceAttributes, error) {
		return parseDataSourceAttributes[domain.QueryDataSourceAttributes](raw)
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
