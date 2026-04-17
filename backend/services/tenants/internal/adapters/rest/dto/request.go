package dto

import (
	"encoding/json"
	"errors"

	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/strategy"
)

var (
	ErrDataSourceAttrParse = errors.New("failed to deserialize data source attributes")
)

type attributeParser func([]byte) (domain.DataSourceAttributes, error)

var attributeParserStrategies = strategy.NewStrategies[domain.DataSourceType, attributeParser]().
	Set(domain.DataSourceTypeStatic, func(data []byte) (domain.DataSourceAttributes, error) {
		return parseAttributes[domain.StaticDataSourceAttributes](data)
	}).
	Set(domain.DataSourceTypeScheduled, func(data []byte) (domain.DataSourceAttributes, error) {
		return parseAttributes[domain.ScheduledDataSourceAttributes](data)
	}).
	Set(domain.DataSourceTypeQuery, func(data []byte) (domain.DataSourceAttributes, error) {
		return parseAttributes[domain.QueryDataSourceAttributes](data)
	})

type TenantRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type DataSourceRequest struct {
	Type       domain.DataSourceType `json:"type" validate:"required"`
	Attributes any                   `json:"attributes" validate:"required"`
}

func RequestToDataSourceAttributes(dataSourceType domain.DataSourceType, raw any) (domain.DataSourceAttributes, error) {
	if dataSourceType == "" {
		return nil, errors.New("field type is required")
	}

	attrBytes, err := json.Marshal(raw)

	if err != nil {
		return nil, err
	}

	strategy, err := attributeParserStrategies.Get(dataSourceType)

	if err != nil {
		return nil, err
	}

	return strategy(attrBytes)
}

func parseAttributes[T any](data []byte) (domain.DataSourceAttributes, error) {
	var attributes T

	if err := json.Unmarshal(data, &attributes); err != nil {
		return nil, err
	}

	return attributes, nil
}
