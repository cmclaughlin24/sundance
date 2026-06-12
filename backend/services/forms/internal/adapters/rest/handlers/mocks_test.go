package handlers

import (
	"context"

	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
	"sundance/backend/services/forms/internal/core/ports/commands"
)

type mockFormsService struct {
	findFn           func(context.Context, ports.FindFormsQuery) ([]*domain.Form, error)
	findByIDFn       func(context.Context, ports.FindByIDQuery[domain.FormID]) (*domain.Form, error)
	createFn         func(context.Context, commands.CreateFormCommand) (*domain.Form, error)
	updateFn         func(context.Context, commands.UpdateFormCommand) (*domain.Form, error)
	deleteFn         func(context.Context, commands.DeleteCommand[domain.FormID]) error
	findVersionsFn   func(context.Context, ports.FindFormVersionsQuery) ([]*domain.FormVersion, error)
	findVersionFn    func(context.Context, ports.FindFormVersionByIDQuery) (*domain.FormVersion, error)
	createVersionFn  func(context.Context, *commands.CreateFormVersionCommand) (*domain.FormVersion, error)
	updateVersionFn  func(context.Context, *commands.UpdateFormVersionCommand) (*domain.FormVersion, error)
	publishVersionFn func(context.Context, commands.PublishFormVersionCommand) (*domain.FormVersion, error)
	retireVersionFn  func(context.Context, commands.RetireFormVersionCommand) (*domain.FormVersion, error)
}

func (s *mockFormsService) Find(ctx context.Context, query ports.FindFormsQuery) ([]*domain.Form, error) {
	return s.findFn(ctx, query)
}

func (s *mockFormsService) FindByID(ctx context.Context, query ports.FindByIDQuery[domain.FormID]) (*domain.Form, error) {
	return s.findByIDFn(ctx, query)
}

func (s *mockFormsService) Create(ctx context.Context, command commands.CreateFormCommand) (*domain.Form, error) {
	return s.createFn(ctx, command)
}

func (s *mockFormsService) Update(ctx context.Context, command commands.UpdateFormCommand) (*domain.Form, error) {
	return s.updateFn(ctx, command)
}

func (s *mockFormsService) Delete(ctx context.Context, command commands.DeleteCommand[domain.FormID]) error {
	return s.deleteFn(ctx, command)
}

func (s *mockFormsService) FindVersions(ctx context.Context, query ports.FindFormVersionsQuery) ([]*domain.FormVersion, error) {
	return s.findVersionsFn(ctx, query)
}

func (s *mockFormsService) FindVersion(ctx context.Context, query ports.FindFormVersionByIDQuery) (*domain.FormVersion, error) {
	return s.findVersionFn(ctx, query)
}

func (s *mockFormsService) CreateVersion(ctx context.Context, command *commands.CreateFormVersionCommand) (*domain.FormVersion, error) {
	return s.createVersionFn(ctx, command)
}

func (s *mockFormsService) UpdateVersion(ctx context.Context, command *commands.UpdateFormVersionCommand) (*domain.FormVersion, error) {
	return s.updateVersionFn(ctx, command)
}

func (s *mockFormsService) PublishVersion(ctx context.Context, command commands.PublishFormVersionCommand) (*domain.FormVersion, error) {
	return s.publishVersionFn(ctx, command)
}

func (s *mockFormsService) RetireVersion(ctx context.Context, command commands.RetireFormVersionCommand) (*domain.FormVersion, error) {
	return s.retireVersionFn(ctx, command)
}

type mockClaims struct {
	subject string
}

func (c *mockClaims) GetSubject() (string, error) {
	return c.subject, nil
}
