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
	return nil, nil
}

func (s *FormsService) FindById(context.Context, domain.FormID) (*domain.Form, error) {
	return nil, nil
}

func (s *FormsService) Create(context.Context, *ports.CreateFormCommand) (*domain.Form, error) {
	return nil, nil
}

func (s *FormsService) Update(context.Context, *ports.UpdateFormCommand) (*domain.Form, error) {
	return nil, nil
}

func (s *FormsService) CreateVersion(ctx context.Context, command *ports.CreateVersionCommand) (*domain.Version, error) {
	if err := s.isValidAccess(ctx, command.TenantID, command.FormID); err != nil {
		return nil, err
	}

	versionNum, err := s.nextVersionNumber(ctx, command.FormID)

	if err != nil {
		return nil, err
	}

	version, err := domain.NewVersion(command.ID, command.FormID, versionNum, domain.VersionStatusDraft)

	if err != nil {
		return nil, err
	}

	version, err = s.repository.Forms.CreateVersion(ctx, version)

	if err != nil {
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

func (s *FormsService) nextVersionNumber(ctx context.Context, formId domain.FormID) (int, error) {
	versions, err := s.repository.Forms.FindVersions(ctx, formId)

	if err != nil {
		if err == common.ErrNotFound {
			return 1, nil
		}

		return 0, err
	}

	maxVersion := 0
	for _, version := range versions {
		if version.Version > maxVersion {
			maxVersion = version.Version
		}
	}

	return maxVersion + 1, nil
}
