package dto

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cmclaughlin24/sundance/backend/pkg/common/stratreg"
	"github.com/cmclaughlin24/sundance/backend/services/tenants/internal/core/domain"
)

var (
	ErrDataSourceAttrParse = errors.New("failed to deserialize data source attributes")
)

type attributeParser func([]byte) (domain.DataSourceAttributes, error)

var attributeParserStrategies = stratreg.New[domain.DataSourceType, attributeParser]().
	Set(domain.DataSourceTypeStatic, func(data []byte) (domain.DataSourceAttributes, error) {
		return parseAttributes[domain.StaticDataSourceAttributes](data)
	}).
	Set(domain.DataSourceTypeScheduled, func(data []byte) (domain.DataSourceAttributes, error) {
		return parseAttributes[domain.ScheduledDataSourceAttributes](data)
	}).
	Set(domain.DataSourceTypeWebhook, func(data []byte) (domain.DataSourceAttributes, error) {
		return parseAttributes[domain.WebhookDataSourceAttributes](data)
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

type webhookDataSourceAttributesResponse struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers,omitempty"`
}

func dataSourceAttributesToResponse(attr domain.DataSourceAttributes) any {
	switch t := attr.(type) {
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
	case domain.WebhookDataSourceAttributes:
		return webhookDataSourceAttributesResponse{
			URL:     t.URL,
			Method:  t.Method,
			Headers: t.Headers,
		}
	default:
		return attr
	}
}
