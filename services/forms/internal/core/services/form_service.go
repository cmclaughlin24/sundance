package services

import (
	"context"
	"log"
	"time"

	"github.com/cmclaughlin24/sundance/common"
	"github.com/cmclaughlin24/sundance/forms/internal/core/domain"
	"github.com/cmclaughlin24/sundance/forms/internal/core/ports"
)

type FormsService struct {
	logger     *log.Logger
	repository *ports.Repository
}

func NewFormsService(logger *log.Logger, repository *ports.Repository) *FormsService {
	return &FormsService{
		logger:     logger,
		repository: repository,
	}
}

func (s *FormsService) Find(ctx context.Context) ([]*domain.Form, error) {
	return s.repository.Forms.Find(ctx)
}

func (s *FormsService) FindById(ctx context.Context, query ports.FindByIdQuery) (*domain.Form, error) {
	form, err := s.repository.Forms.FindById(ctx, query.ID)

	if err != nil {
		return nil, err
	}

	if form.TenantID != query.TenantID {
		return nil, common.ErrUnauthorized
	}

	return form, nil
}

func (s *FormsService) Create(ctx context.Context, command *ports.CreateFormCommand) (*domain.Form, error) {
	form, err := domain.NewForm(domain.FormID(""), command.TenantID, command.Name, command.Description)

	if err != nil {
		return nil, err
	}

	form, err = s.repository.Forms.Create(ctx, form)

	if err != nil {
		return nil, err
	}

	return form, nil
}

func (s *FormsService) Update(ctx context.Context, command *ports.UpdateFormCommand) (*domain.Form, error) {
	if err := s.isValidAccess(ctx, command.TenantID, command.ID); err != nil {
		return nil, err
	}

	form, err := s.repository.Forms.FindById(ctx, command.ID)

	if err != nil {
		return nil, err
	}

	if err := form.Update(command.Name, command.Description); err != nil {
		return nil, err
	}

	form, err = s.repository.Forms.Update(ctx, form)

	if err != nil {
		return nil, err
	}

	return form, nil
}

func (s *FormsService) CreateVersion(ctx context.Context, command *ports.CreateVersionCommand) (*domain.Version, error) {
	if err := s.isValidAccess(ctx, command.TenantID, command.FormID); err != nil {
		return nil, err
	}

	txCtx, err := s.repository.Database.BeginTx(ctx)

	if err != nil {
		return nil, err
	}

	defer s.repository.Database.RollbackTx(txCtx)

	versionNum, err := s.repository.Forms.FindNextVersionNumber(txCtx, command.FormID)

	if err != nil {
		return nil, err
	}

	version, err := domain.NewVersion(command.ID, command.FormID, versionNum, domain.VersionStatusDraft)

	if err != nil {
		return nil, err
	}

	version, err = s.repository.Forms.CreateVersion(txCtx, version)

	if err != nil {
		return nil, err
	}

	if err := s.repository.Database.CommitTx(txCtx); err != nil {
		return nil, err
	}

	return version, nil
}

func (s *FormsService) UpdateVersion(ctx context.Context, command *ports.UpdateVersionCommand) (*domain.Version, error) {
	if err := s.isValidAccess(ctx, command.TenantID, command.FormID); err != nil {
		return nil, err
	}

	version, err := s.repository.Forms.FindVersion(ctx, command.FormID, command.ID)

	if err != nil {
		return nil, err
	}

	if err := version.UpdatePages(command.Pages...); err != nil {
		return nil, err
	}

	version, err = s.repository.Forms.UpdateVersion(ctx, version)

	if err != nil {
		return nil, err
	}

	return version, nil
}

func (s *FormsService) PublishVersion(ctx context.Context, command *ports.PublishVersionCommand) (*domain.Version, error) {
	if err := s.isValidAccess(ctx, command.TenantID, command.FormID); err != nil {
		return nil, err
	}

	version, err := s.repository.Forms.FindVersion(ctx, command.FormID, command.VersionID)

	if err != nil {
		return nil, err
	}

	if err := version.Publish(command.UserId, time.Now()); err != nil {
		return nil, err
	}

	version, err = s.repository.Forms.UpdateVersion(ctx, version)

	if err != nil {
		return nil, err
	}

	return version, nil
}

func (s *FormsService) RetireVersion(ctx context.Context, command *ports.RetireVersionCommand) (*domain.Version, error) {
	if err := s.isValidAccess(ctx, command.TenantID, command.FormID); err != nil {
		return nil, err
	}

	version, err := s.repository.Forms.FindVersion(ctx, command.FormID, command.VersionID)

	if err != nil {
		return nil, err
	}

	if err := version.Retire(command.UserId, time.Now()); err != nil {
		return nil, err
	}

	version, err = s.repository.Forms.UpdateVersion(ctx, version)

	if err != nil {
		return nil, err
	}

	return version, nil
}

func (s *FormsService) isValidAccess(ctx context.Context, tenantId string, formId domain.FormID) error {
	form, err := s.repository.Forms.FindById(ctx, formId)

	if err != nil {
		return err
	}

	if form.TenantID != tenantId {
		return common.ErrUnauthorized
	}

	return nil
}
