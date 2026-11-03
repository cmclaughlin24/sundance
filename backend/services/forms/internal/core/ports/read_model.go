package ports

import "sundance/backend/services/forms/internal/core/domain"

type AttributeDescriptor struct {
	Type        string
	Name        string
	Description string
	Required    bool
}

type ElementTypeDescriptor struct {
	Type        domain.ElementType
	Description string
	Attributes  []AttributeDescriptor
}
