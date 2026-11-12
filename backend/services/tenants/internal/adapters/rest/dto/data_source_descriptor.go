package dto

import (
	"sundance/backend/services/tenants/internal/core/domain"
	"sundance/backend/services/tenants/internal/core/ports"
)

type AttributeDescriptorResponse struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

type DataSourceTypeDescriptorResponse struct {
	Type        domain.DataSourceType         `json:"type"`
	Description string                        `json:"description"`
	Attributes  []AttributeDescriptorResponse `json:"attributes"`
}

func DataSourceTypeDescriptorToResponse(d ports.DataSourceTypeDescriptor) DataSourceTypeDescriptorResponse {
	attrs := make([]AttributeDescriptorResponse, 0, len(d.Attributes))
	for _, a := range d.Attributes {
		attrs = append(attrs, AttributeDescriptorResponse{
			Type:        a.Type,
			Name:        a.Name,
			Description: a.Description,
			Required:    a.Required,
		})
	}

	return DataSourceTypeDescriptorResponse{
		Type:        d.Type,
		Description: d.Description,
		Attributes:  attrs,
	}
}

func DataSourceTypeDescriptorsToResponse(descriptors []ports.DataSourceTypeDescriptor) []DataSourceTypeDescriptorResponse {
	dtos := make([]DataSourceTypeDescriptorResponse, 0, len(descriptors))
	for _, d := range descriptors {
		dtos = append(dtos, DataSourceTypeDescriptorToResponse(d))
	}

	return dtos
}
