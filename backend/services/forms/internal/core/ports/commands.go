package ports

import (
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
)

type baseFormCommand struct {
	Name        string `validate:"required,max=75"`
	Description string `validate:"required,max=250"`
}

type CreateFormCommand struct {
	baseFormCommand
}

func NewCreateFormCommand(name, description string) *CreateFormCommand {
	return &CreateFormCommand{
		baseFormCommand: baseFormCommand{
			Name:        name,
			Description: description,
		},
	}
}

type UpdateFormCommand struct {
	baseFormCommand
	ID domain.FormID `validate:"required"`
}

func NewUpdateFormCommand(id domain.FormID, name, description string) *UpdateFormCommand {
	return &UpdateFormCommand{
		baseFormCommand: baseFormCommand{
			Name:        name,
			Description: description,
		},
		ID: id,
	}
}

type baseVersionCommand struct {
	FormID domain.FormID `validate:"required"`
}

type CreateVersionCommand struct {
	baseVersionCommand
}

func NewCreateVersionCommand(formId domain.FormID) *CreateVersionCommand {
	return &CreateVersionCommand{
		baseVersionCommand{
			FormID: formId,
		},
	}
}

type UpdateVersionCommand struct {
	baseVersionCommand
	VersionID domain.VersionID `validate:"required"`
	Pages     []*domain.Page
}

func NewUpdateVersionCommand(id domain.VersionID, formID domain.FormID, pages []*domain.Page) *UpdateVersionCommand {
	return &UpdateVersionCommand{
		baseVersionCommand: baseVersionCommand{
			FormID: formID,
		},
		VersionID: id,
		Pages:     pages,
	}
}

type PublishVersionCommand struct {
	baseVersionCommand
	VersionID domain.VersionID `validate:"required"`
	UserID    string           `validate:"required"`
}

func NewPublishVersionCommand(formID domain.FormID, versionID domain.VersionID, userID string) *PublishVersionCommand {
	return &PublishVersionCommand{
		baseVersionCommand: baseVersionCommand{
			FormID: formID,
		},
		VersionID: versionID,
		UserID:    userID,
	}
}

type RetireVersionCommand struct {
	baseVersionCommand
	VersionID domain.VersionID `validate:"required"`
	UserID    string           `validate:"required"`
}

func NewRetireVersionCommand(formID domain.FormID, versionID domain.VersionID, userID string) *RetireVersionCommand {
	return &RetireVersionCommand{
		baseVersionCommand: baseVersionCommand{
			FormID: formID,
		},
		VersionID: versionID,
		UserID:    userID,
	}
}
