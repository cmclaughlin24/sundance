package ports

import "sundance/backend/services/tenants/internal/core/domain"

type AttributeDescriptor struct {
	Type        string
	Name        string
	Description string
	Required    bool
}

type DataSourceTypeDescriptor struct {
	Type        domain.DataSourceType
	Description string
	Attributes  []AttributeDescriptor
}
