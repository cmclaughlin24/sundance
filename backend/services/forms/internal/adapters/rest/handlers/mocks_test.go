package handlers

import (
	"context"

	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

type mockFormsService struct {
	findFn           func(context.Context, *ports.FindFormsQuery) ([]*domain.Form, error)
	findByIDFn       func(context.Context, ports.FindByIDQuery[domain.FormID]) (*domain.Form, error)
	createFn         func(context.Context, *ports.CreateFormCommand) (*domain.Form, error)
	updateFn         func(context.Context, *ports.UpdateFormCommand) (*domain.Form, error)
	deleteFn         func(context.Context, *ports.DeleteCommand[domain.FormID]) error
	findVersionsFn   func(context.Context, *ports.FindFormVersionsQuery) ([]*domain.FormVersion, error)
	findVersionFn    func(context.Context, *ports.FindFormVersionByIDQuery) (*domain.FormVersion, error)
	createVersionFn  func(context.Context, *ports.CreateFormVersionCommand) (*domain.FormVersion, error)
	updateVersionFn  func(context.Context, *ports.UpdateFormVersionCommand) (*domain.FormVersion, error)
	publishVersionFn func(context.Context, *ports.PublishFormVersionCommand) (*domain.FormVersion, error)
	retireVersionFn  func(context.Context, *ports.RetireFormVersionCommand) (*domain.FormVersion, error)
}

func (s *mockFormsService) Find(ctx context.Context, query *ports.FindFormsQuery) ([]*domain.Form, error) {
	return s.findFn(ctx, query)
}

func (s *mockFormsService) FindByID(ctx context.Context, query ports.FindByIDQuery[domain.FormID]) (*domain.Form, error) {
	return s.findByIDFn(ctx, query)
}

func (s *mockFormsService) Create(ctx context.Context, command *ports.CreateFormCommand) (*domain.Form, error) {
	return s.createFn(ctx, command)
}

func (s *mockFormsService) Update(ctx context.Context, command *ports.UpdateFormCommand) (*domain.Form, error) {
	return s.updateFn(ctx, command)
}

func (s *mockFormsService) Delete(ctx context.Context, command *ports.DeleteCommand[domain.FormID]) error {
	return s.deleteFn(ctx, command)
}

func (s *mockFormsService) FindVersions(ctx context.Context, query *ports.FindFormVersionsQuery) ([]*domain.FormVersion, error) {
	return s.findVersionsFn(ctx, query)
}

func (s *mockFormsService) FindVersion(ctx context.Context, query *ports.FindFormVersionByIDQuery) (*domain.FormVersion, error) {
	return s.findVersionFn(ctx, query)
}

func (s *mockFormsService) CreateVersion(ctx context.Context, command *ports.CreateFormVersionCommand) (*domain.FormVersion, error) {
	return s.createVersionFn(ctx, command)
}

func (s *mockFormsService) UpdateVersion(ctx context.Context, command *ports.UpdateFormVersionCommand) (*domain.FormVersion, error) {
	return s.updateVersionFn(ctx, command)
}

func (s *mockFormsService) PublishVersion(ctx context.Context, command *ports.PublishFormVersionCommand) (*domain.FormVersion, error) {
	return s.publishVersionFn(ctx, command)
}

func (s *mockFormsService) RetireVersion(ctx context.Context, command *ports.RetireFormVersionCommand) (*domain.FormVersion, error) {
	return s.retireVersionFn(ctx, command)
}

type mockClaims struct {
	subject string
}

func (c *mockClaims) GetSubject() string {
	return c.subject
}
