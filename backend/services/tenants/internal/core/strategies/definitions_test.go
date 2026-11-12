package strategies_test

import (
	"reflect"
	"sundance/backend/services/tenants/internal/core/ports"
	"sundance/backend/services/tenants/internal/core/strategies"
	"testing"
)

func TestDataSourceDefinitionStrategy_Descriptors(t *testing.T) {
	clients := &ports.Clients{}

	httpRequestAttributes := []ports.AttributeDescriptor{
		{Type: "string", Name: "url", Description: "Endpoint URL to request lookup options from.", Required: true},
		{Type: "string", Name: "method", Description: "HTTP method used for the request (e.g. GET, POST).", Required: true},
		{Type: "object", Name: "headers", Description: "HTTP headers sent with the request, such as authorization.", Required: true},
		{Type: "string", Name: "valueField", Description: "Field in each response item to use as the option value.", Required: true},
		{Type: "string", Name: "labelField", Description: "Field in each response item to use as the option label.", Required: true},
	}

	tests := []struct {
		name     string
		strategy ports.DataSourceDefinitionStrategy
		want     []ports.AttributeDescriptor
	}{
		{
			"should describe the static strategy",
			strategies.NewStaticDataSourceDefinitionStrategy(logger),
			[]ports.AttributeDescriptor{
				{Type: "array", Name: "data", Description: "List of value and label options returned by the data source.", Required: true},
			},
		},
		{
			"should describe the scheduled strategy without system-managed attributes",
			strategies.NewScheduledDataSourceDefinitionStrategy(logger),
			append(append([]ports.AttributeDescriptor{}, httpRequestAttributes...),
				ports.AttributeDescriptor{Type: "number", Name: "intervalHours", Description: "Number of hours between data refreshes.", Required: true},
			),
		},
		{
			"should describe the webhook strategy",
			strategies.NewWebhookDataSourceDefinitionStrategy(logger, clients),
			append(append([]ports.AttributeDescriptor{}, httpRequestAttributes...),
				ports.AttributeDescriptor{Type: "array", Name: "requiredKeys", Description: "Parameter keys that must be provided on each lookup request.", Required: true},
			),
		},
		{
			"should describe the data lake strategy",
			strategies.NewDataLakeDataSourceDefinitionStrategy(logger, clients),
			[]ports.AttributeDescriptor{
				{Type: "string", Name: "query", Description: "Parameterized SQL query used to resolve lookup options.", Required: true},
				{Type: "array", Name: "requiredKeys", Description: "Query parameter keys that must be provided on each lookup request.", Required: true},
				{Type: "array", Name: "optionalKeys", Description: "Query parameter keys that may optionally be provided on each lookup request.", Required: true},
				{Type: "string", Name: "catalog", Description: "Data lake catalog to run the query against.", Required: true},
				{Type: "string", Name: "schema", Description: "Data lake schema to run the query against.", Required: true},
				{Type: "string", Name: "valueField", Description: "Column in the query result to use as the option value.", Required: true},
				{Type: "string", Name: "labelField", Description: "Column in the query result to use as the option label.", Required: true},
				{Type: "integer", Name: "timeoutMs", Description: "Maximum query execution time in milliseconds.", Required: true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act.
			gotDescription := tt.strategy.Description()
			gotAttributes := tt.strategy.Attributes()

			// Assert.
			if gotDescription == "" {
				t.Errorf("expected a description but got an empty string")
			}

			if !reflect.DeepEqual(tt.want, gotAttributes) {
				t.Errorf("expected %+v but got %+v", tt.want, gotAttributes)
			}
		})
	}
}
