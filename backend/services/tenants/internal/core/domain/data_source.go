package domain

import (
	"errors"
	"time"

	"sundance/backend/pkg/common/validate"
)

type DataSourceID string

type DataSourceType string

const (
	DataSourceTypeStatic    DataSourceType = "static"
	DataSourceTypeScheduled DataSourceType = "scheduled"
	DataSourceTypeWebhook   DataSourceType = "webhook"
)

var (
	ErrInvalidSourceType           = errors.New("invalid data source type")
	ErrInvalidSourceTypeAttributes = errors.New("invalid data source attributes for type")
)

type DataSource struct {
	ID          DataSourceID
	TenantID    TenantID `validate:"required"`
	Name        string   `validate:"required"`
	Description string
	Type        DataSourceType
	Attributes  DataSourceAttributes
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewDataSource(
	tenantID TenantID,
	name,
	description string,
	sourceType DataSourceType,
	attr DataSourceAttributes,
) (*DataSource, error) {
	if !isValidSourceType(sourceType) {
		return nil, ErrInvalidSourceType
	}

	if !isValidAttributeType(sourceType, attr) {
		return nil, ErrInvalidSourceTypeAttributes
	}

	ds := &DataSource{
		ID:          DataSourceID(NewID()),
		TenantID:    tenantID,
		Name:        name,
		Description: description,
		Type:        sourceType,
		Attributes:  attr,
		CreatedAt:   Now(),
	}

	if err := validate.ValidateStruct(ds); err != nil {
		return nil, err
	}

	return ds, nil
}

func HydrateDataSource(
	ID DataSourceID,
	tenantID TenantID,
	name,
	description string,
	sourceType DataSourceType,
	attr DataSourceAttributes,
	createdAt,
	updatedAt time.Time,
) *DataSource {
	return &DataSource{
		ID:          ID,
		TenantID:    tenantID,
		Name:        name,
		Description: description,
		Type:        sourceType,
		Attributes:  attr,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func (ds *DataSource) Update(name, description string, sourceType DataSourceType, attr DataSourceAttributes) error {
	if !isValidSourceType(sourceType) {
		return ErrInvalidSourceType
	}

	if !isValidAttributeType(sourceType, attr) {
		return ErrInvalidSourceTypeAttributes
	}

	cpy := *ds
	cpy.Name = name
	cpy.Description = description
	cpy.Type = sourceType
	cpy.Attributes = attr

	if err := validate.ValidateStruct(cpy); err != nil {
		return err
	}

	*ds = cpy
	ds.UpdatedAt = Now()

	return nil
}

func (ds *DataSource) UpdateAttributes(attr DataSourceAttributes) error {
	if !isValidAttributeType(ds.Type, attr) {
		return ErrInvalidSourceTypeAttributes
	}

	ds.Attributes = attr
	ds.UpdatedAt = Now()

	return nil
}

var isValidSourceType = validate.NewTypeValidator([]DataSourceType{
	DataSourceTypeStatic,
	DataSourceTypeScheduled,
	DataSourceTypeWebhook,
})

func isValidAttributeType(sourceType DataSourceType, attributes DataSourceAttributes) bool {
	switch attributes.(type) {
	case StaticDataSourceAttributes:
		return sourceType == DataSourceTypeStatic
	case ScheduledDataSourceAttributes:
		return sourceType == DataSourceTypeScheduled
	case WebhookDataSourceAttributes:
		return sourceType == DataSourceTypeWebhook
	default:
		return false
	}
}
