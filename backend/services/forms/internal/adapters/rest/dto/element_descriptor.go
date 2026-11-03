package dto

import (
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

type AttributeDescriptorResponse struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

type ElementTypeDescriptorResponse struct {
	Type        domain.ElementType            `json:"type"`
	Description string                        `json:"description"`
	Attributes  []AttributeDescriptorResponse `json:"attributes"`
}

func ElementTypeDescriptorToResponse(d ports.ElementTypeDescriptor) ElementTypeDescriptorResponse {
	attrs := make([]AttributeDescriptorResponse, 0, len(d.Attributes))
	for _, a := range d.Attributes {
		attrs = append(attrs, AttributeDescriptorResponse{
			Type:        a.Type,
			Name:        a.Name,
			Description: a.Description,
			Required:    a.Required,
		})
	}

	return ElementTypeDescriptorResponse{
		Type:        d.Type,
		Description: d.Description,
		Attributes:  attrs,
	}
}

func ElementTypeDescriptorsToResponse(descriptors []ports.ElementTypeDescriptor) []ElementTypeDescriptorResponse {
	dtos := make([]ElementTypeDescriptorResponse, 0, len(descriptors))
	for _, d := range descriptors {
		dtos = append(dtos, ElementTypeDescriptorToResponse(d))
	}

	return dtos
}

