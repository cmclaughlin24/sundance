package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"slices"

	"sundance/backend/pkg/database"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
	"sundance/backend/services/forms/internal/core/strategies"
)

var (
	ErrNoEligibleTagVersion     = errors.New("")
	ErrMultipleActiveTagVersion = errors.New("")
)

type ruleByTypeGetter interface {
	GetRuleByType(domain.RuleType) *domain.Rule
}

type factCandidate struct {
	ftm   domain.FieldTagMapping
	value any
}

type tagAggregate struct {
	tag      domain.Tag
	versions []*domain.TagVersion
}

type submissionJobsService struct {
	logger                   *slog.Logger
	evaluator                ports.RuleEvaluator
	database                 database.Database
	formVersionsRepository   ports.FormVersionRepository
	submissionsRepository    ports.SubmissionsRepository
	tagsRepository           ports.TagsRepository
	tagVersionsRepository    ports.TagVersionsRepository
	fieldValidatorStrategies ports.FieldValidatorRegistry
}

func NewSubmissionJobsService(
	logger *slog.Logger,
	evaluator ports.RuleEvaluator,
	repository *ports.Repository,
	strategies *ports.Strategies,
) ports.SubmissionJobsAPI {
	return &submissionJobsService{
		logger:                   logger,
		evaluator:                evaluator,
		database:                 repository.Database,
		formVersionsRepository:   repository.FormVersions,
		submissionsRepository:    repository.Submissions,
		tagsRepository:           repository.Tags,
		tagVersionsRepository:    repository.TagVersions,
		fieldValidatorStrategies: strategies.FieldValidator,
	}
}

func (s *submissionJobsService) Find(ctx context.Context, query ports.FindSubmissionJobsQuery) ([]domain.SubmissionID, error) {
	s.logger.DebugContext(ctx, "listing submission jobs")

	if err := query.Validate(); err != nil {
		s.logger.WarnContext(ctx, "submission job listing failed; invalid query", "error", err)
		return nil, err
	}

	ids, err := s.submissionsRepository.FindJobs(ctx, &ports.FindSubmissionsFilter{
		Take:     query.Take,
		Statuses: []domain.SubmissionStatus{domain.SubmissionStatusPending},
	})

	if err != nil {
		s.logger.ErrorContext(ctx, "failed to retrieve submission jobs", "error", err)
		return nil, err
	}

	return ids, nil
}

func (s *submissionJobsService) Process(ctx context.Context, id domain.SubmissionID) error {
	s.logger.DebugContext(ctx, "processing submission job", "submission_id", id)

	submission, err := s.submissionsRepository.FindByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to retrieve submission job", "submission_id", id, "error", err)
		return err
	}

	if submission.Status == domain.SubmissionStatusAccepted || submission.Status == domain.SubmissionStatusRejected {
		s.logger.WarnContext(ctx, "skipping submission job; invalid status", "submission_id", submission.ID, "status", submission.Status)
		return nil
	}

	_, err = s.sanitize(ctx, submission)

	if err := s.recordAttempt(ctx, submission, err); err != nil {
		return err
	}

	s.logger.InfoContext(ctx, "submission job recorded", "submission_id", submission.ID, "status", submission.Status)

	return err
}

func (s *submissionJobsService) sanitize(ctx context.Context, submission *domain.Submission) ([]*domain.CanonicalFact, error) {
	factCandidates, err := s.extractFactCandidates(ctx, submission)
	if err != nil {
		return nil, err
	}

	facts, err := s.normalize(ctx, factCandidates)
	if err != nil {
		return nil, err
	}

	return facts, nil
}

func (s *submissionJobsService) extractFactCandidates(ctx context.Context, submission *domain.Submission) ([]factCandidate, error) {
	version, err := s.formVersionsRepository.FindByID(ctx, submission.VersionID)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to retrieve version for submission job", "submission_id", submission.ID, "form_id", submission.FormID, "version_id", submission.VersionID, "error", err)
		return nil, err
	}

	if version.Status == domain.FormVersionStatusDraft {
		s.logger.WarnContext(ctx, "skipping submission job; invalid status", "submission_id", submission.ID, "form_id", submission.FormID, "version_id", submission.VersionID, "version_status", version.Status)
		return nil, domain.ErrInvalidVersionStatus
	}

	evalCtx := make(ports.RuleEvaluationContext, len(submission.Values))
	for _, field := range version.FlatFields() {
		if val, ok := submission.GetFieldValue(field.ID); ok {
			evalCtx[field.Key] = val.Value
		}
	}

	candidates := make([]factCandidate, 0)

pageLoop:
	for _, page := range version.GetPages() {
		visible, err := s.shouldValidate(ctx, page, evalCtx)

		if err != nil {
			return nil, err
		}

		if !visible {
			continue pageLoop
		}

	sectionLoop:
		for _, section := range page.GetSections() {
			visible, err := s.shouldValidate(ctx, section, evalCtx)

			if err != nil {
				return nil, err
			}

			if !visible {
				continue sectionLoop
			}

		fieldLoop:
			for _, field := range section.GetFields() {
				visible, err := s.shouldValidate(ctx, field, evalCtx)

				if err != nil {
					return nil, err
				}

				if !visible {
					continue fieldLoop
				}

				required, err := s.isRequired(ctx, field, evalCtx)
				if err != nil {
					return nil, err
				}

				if required != nil {
					field.Attributes.SetIsRequired(*required)
				}

				if err := s.validateField(ctx, field, submission); err != nil {
					return nil, err
				}

				for _, ftm := range field.GetTags() {
					// NOTE: The second return fv (ok) from submission.GetFieldValue(field.ID) is ignored here
					// because it was checked during the field validation.
					fv, _ := submission.GetFieldValue(field.ID)
					candidates = append(candidates, factCandidate{*ftm, fv.Value})
				}
			}
		}
	}

	return candidates, nil
}

func (s *submissionJobsService) normalize(ctx context.Context, factCandidates []factCandidate) ([]*domain.CanonicalFact, error) {
	candidatesByVersion := make(map[domain.TagVersionID][]factCandidate)
	for _, fc := range factCandidates {
		candidatesByVersion[fc.ftm.TagVersionID] = append(candidatesByVersion[fc.ftm.TagVersionID], fc)
	}

	tags, err := s.getTags(ctx, slices.Collect(maps.Keys(candidatesByVersion)))
	if err != nil {
		return nil, err
	}

	facts := make([]*domain.CanonicalFact, 0)
	for _, ta := range tags {
		version, err := s.selectTagVersion(ctx, ta.versions)
		if err != nil {
			return nil, err
		}

		f, err := s.evaluateCandidates(ctx, ta.tag, *version, candidatesByVersion[version.ID])
		if err != nil {
			return nil, err
		}

		facts = append(facts, f...)
	}

	return facts, nil
}

func (s *submissionJobsService) validateField(ctx context.Context, field *domain.Field, submission *domain.Submission) error {
	fieldValidator, err := s.fieldValidatorStrategies.Get(field.Type)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to process submission; missing field validation strategy", "submission_id", submission.ID, "form_id", submission.FormID, "version_id", submission.VersionID, "field_id", field.ID, "field_type", field.Type)
		return err
	}

	value, ok := submission.GetFieldValue(field.ID)
	if !ok {
		if field.Attributes.GetIsRequired() {
			s.logger.WarnContext(ctx, "submission validation failed; required field missing", "submission_id", submission.ID, "form_id", submission.FormID, "version_id", submission.VersionID, "field_id", field.ID, "field_key", field.Key)
			return fmt.Errorf("%w; id=%s key=%s", strategies.ErrFieldRequired, field.ID, field.Key)
		}

		return nil
	}

	if err = fieldValidator.Validate(ctx, *field, *value); err != nil {
		return err
	}

	return nil
}

func (s *submissionJobsService) isRequired(ctx context.Context, rg ruleByTypeGetter, evalCtx ports.RuleEvaluationContext) (*bool, error) {
	rule := rg.GetRuleByType(domain.RuleTypeRequired)
	if rule == nil {
		return nil, nil
	}

	result, err := s.evaluator.Evaluate(ctx, rule, evalCtx)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *submissionJobsService) shouldValidate(ctx context.Context, rg ruleByTypeGetter, evalCtx ports.RuleEvaluationContext) (bool, error) {
	rule := rg.GetRuleByType(domain.RuleTypeVisible)

	if rule == nil {
		return true, nil
	}

	return s.evaluator.Evaluate(ctx, rule, evalCtx)
}

func (s *submissionJobsService) recordAttempt(ctx context.Context, submission *domain.Submission, err error) error {
	if err == nil {
		submission.Accept()
	} else if shouldReject(err) {
		submission.Reject(err)
	} else {
		submission.Fail(err)
	}

	txCtx, err := s.database.BeginTx(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to begin transaction", "submission_id", submission.ID, "error", err)
		return err
	}

	defer s.database.RollbackTx(txCtx)

	_, err = s.submissionsRepository.Upsert(txCtx, submission)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to upsert submission", "submission_id", submission.ID, "error", err)
		return err
	}

	if err := s.database.CommitTx(txCtx); err != nil {
		s.logger.ErrorContext(ctx, "failed to commit transaction", "submission_id", submission.ID, "error", err)
		return err
	}

	return nil
}

func (s *submissionJobsService) getTags(ctx context.Context, ids []domain.TagVersionID) ([]tagAggregate, error) {
	versions, err := s.tagVersionsRepository.FindByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	versionsByTag := make(map[domain.TagID][]*domain.TagVersion)
	for _, v := range versions {
		versionsByTag[v.TagID] = append(versionsByTag[v.TagID], v)
	}

	tags, err := s.tagsRepository.FindByIDs(ctx, slices.Collect(maps.Keys(versionsByTag)))
	if err != nil {
		return nil, err
	}

	aggregates := make([]tagAggregate, 0, len(tags))
	for _, t := range tags {
		aggregates = append(aggregates, tagAggregate{*t, versionsByTag[t.ID]})
	}

	return aggregates, nil
}

func (s *submissionJobsService) selectTagVersion(ctx context.Context, versions []*domain.TagVersion) (*domain.TagVersion, error) {
	var active *domain.TagVersion
	var deprecated *domain.TagVersion

	for _, v := range versions {
		switch v.Status {
		case domain.TagStatusDraft, domain.TagStatusRetired:
			s.logger.ErrorContext(ctx, "invalid tag version status in form definition", "tag_version_id", v.ID, "status", v.Status)
		case domain.TagStatusDeprecated:
			if deprecated == nil {
				deprecated = v
			} else if deprecated.Version < v.Version {
				deprecated = v
			}
		case domain.TagStatusActive:
			if active != nil {
				return nil, ErrMultipleActiveTagVersion
			}

			active = v
		}
	}

	if active != nil {
		return active, nil
	}

	return deprecated, nil
}

func (s *submissionJobsService) evaluateCandidates(
	ctx context.Context,
	tag domain.Tag,
	version domain.TagVersion,
	candidates []factCandidate,
) ([]*domain.CanonicalFact, error) {
	facts := make([]*domain.CanonicalFact, 0)

	return facts, nil
}

func shouldReject(err error) bool {
	return errors.Is(err, domain.ErrInvalidVersionStatus) ||
		errors.Is(err, strategies.ErrFieldValidation) ||
		errors.Is(err, strategies.ErrFieldRequired) ||
		errors.Is(err, strategies.ErrFieldTypeValue)
}
