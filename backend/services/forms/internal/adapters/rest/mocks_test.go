package rest

import (
	"context"

	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/ports"
)

type mockFormsService struct {
	findFn           func(context.Context, *ports.FindFormsQuery) ([]*domain.Form, error)
	findByIDFn       func(context.Context, *ports.FindFormsByIDQuery) (*domain.Form, error)
	createFn         func(context.Context, *ports.CreateFormCommand) (*domain.Form, error)
	updateFn         func(context.Context, *ports.UpdateFormCommand) (*domain.Form, error)
	deleteFn         func(context.Context, *ports.RemoveFormCommand) error
	findVersionsFn   func(context.Context, *ports.FindVersionsQuery) ([]*domain.Version, error)
	findVersionFn    func(context.Context, *ports.FindVersionByIDQuery) (*domain.Version, error)
	createVersionFn  func(context.Context, *ports.CreateVersionCommand) (*domain.Version, error)
	updateVersionFn  func(context.Context, *ports.UpdateVersionCommand) (*domain.Version, error)
	publishVersionFn func(context.Context, *ports.PublishVersionCommand) (*domain.Version, error)
	retireVersionFn  func(context.Context, *ports.RetireVersionCommand) (*domain.Version, error)
}

func (s *mockFormsService) Find(ctx context.Context, query *ports.FindFormsQuery) ([]*domain.Form, error) {
	return s.findFn(ctx, query)
}

func (s *mockFormsService) FindByID(ctx context.Context, query *ports.FindFormsByIDQuery) (*domain.Form, error) {
	return s.findByIDFn(ctx, query)
}

func (s *mockFormsService) Create(ctx context.Context, command *ports.CreateFormCommand) (*domain.Form, error) {
	return s.createFn(ctx, command)
}

func (s *mockFormsService) Update(ctx context.Context, command *ports.UpdateFormCommand) (*domain.Form, error) {
	return s.updateFn(ctx, command)
}

func (s *mockFormsService) Delete(ctx context.Context, command *ports.RemoveFormCommand) error {
	return s.deleteFn(ctx, command)
}

func (s *mockFormsService) FindVersions(ctx context.Context, query *ports.FindVersionsQuery) ([]*domain.Version, error) {
	return s.findVersionsFn(ctx, query)
}

func (s *mockFormsService) FindVersion(ctx context.Context, query *ports.FindVersionByIDQuery) (*domain.Version, error) {
	return s.findVersionFn(ctx, query)
}

func (s *mockFormsService) CreateVersion(ctx context.Context, command *ports.CreateVersionCommand) (*domain.Version, error) {
	return s.createVersionFn(ctx, command)
}

func (s *mockFormsService) UpdateVersion(ctx context.Context, command *ports.UpdateVersionCommand) (*domain.Version, error) {
	return s.updateVersionFn(ctx, command)
}

func (s *mockFormsService) PublishVersion(ctx context.Context, command *ports.PublishVersionCommand) (*domain.Version, error) {
	return s.publishVersionFn(ctx, command)
}

func (s *mockFormsService) RetireVersion(ctx context.Context, command *ports.RetireVersionCommand) (*domain.Version, error) {
	return s.retireVersionFn(ctx, command)
}

type mockClaims struct {
	subject string
}

func (c *mockClaims) GetSubject() string {
	return c.subject
}
