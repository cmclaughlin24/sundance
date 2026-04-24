package dto

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/strategy"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
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

func parseAttributes[T domain.DataSourceAttributes](data []byte) (domain.DataSourceAttributes, error) {
	var attributes T

	if err := json.Unmarshal(data, &attributes); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDataSourceAttrParse, err)
	}

	return attributes, nil
}

func RequestToDataSourceAttributes(dataSourceType domain.DataSourceType, raw any) (domain.DataSourceAttributes, error) {
	if dataSourceType == "" {
		return nil, errors.New("data source type is required")
	}

	attrBytes, err := json.Marshal(raw)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDataSourceAttrParse, err)
	}

	strategy, err := attributeParserStrategies.Get(dataSourceType)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDataSourceAttrParse, err)
	}

	return strategy(attrBytes)
}

type staticDataSourceAttributesResponse struct {
	Data []*LookupResponse `json:"data"`
}

type scheduledDataSourceAttributesResponse struct {
	Data []*LookupResponse `json:"data"`
}

type queryDataSourceAttributesResponse struct {
	Type     domain.QueryDataSourceType `json:"type"`
	Resource string                     `json:"resource"`
}

func dataSourceAttributesToResponse(attr domain.DataSourceAttributes) any {
	switch t := attr.(type) {
	case domain.QueryDataSourceAttributes:
		return queryDataSourceAttributesResponse{
			Type:     t.Type,
			Resource: t.Resource,
		}
	case domain.ScheduledDataSourceAttributes:
		data := LookupsToResponse(t.Data)

		return scheduledDataSourceAttributesResponse{
			Data: data,
		}
	case domain.StaticDataSourceAttributes:
		data := LookupsToResponse(t.Data)

		return staticDataSourceAttributesResponse{
			Data: data,
		}
	default:
		return attr
	}
}
