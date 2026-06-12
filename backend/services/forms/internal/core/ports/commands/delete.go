package commands

import "sundance/backend/pkg/common/validate"

type DeleteCommand[T comparable] struct {
	TenantID string `validate:"required"`
	ID       T      `validate:"required"`
}

func NewDeleteCommand[T comparable](tenantID string, id T) DeleteCommand[T] {
	return DeleteCommand[T]{
		TenantID: tenantID,
		ID:       id,
	}
}

func (c DeleteCommand[T]) Validate() error {
	return validate.ValidateStruct(c)
}
