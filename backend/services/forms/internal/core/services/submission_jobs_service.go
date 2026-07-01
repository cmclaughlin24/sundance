package services

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"slices"

	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
	"sundance/backend/services/forms/internal/core/strategies"
)

type ruleByTypeGetter interface {
	GetRuleByType(domain.RuleType) *domain.Rule
}

type factCandidate struct {
	ftm             domain.FieldTagMapping
	value           any
	collectionIndex *int
}

type tagAggregate struct {
	tag      domain.Tag
	versions []*domain.TagVersion
}

type submissionJobsService struct {
	logger                   *slog.Logger
	evaluator                ports.RuleEvaluator
	repository               *ports.Repository
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
		repository:               repository,
		fieldValidatorStrategies: strategies.FieldValidator,
	}
}

func (s *submissionJobsService) Find(ctx context.Context, query ports.FindSubmissionJobsQuery) ([]domain.SubmissionID, error) {
	s.logger.DebugContext(ctx, "listing submission jobs")

	if err := query.Validate(); err != nil {
		s.logger.WarnContext(ctx, "submission job listing failed; invalid query", "error", err)
		return nil, err
	}

	ids, err := s.repository.Submissions.FindJobs(ctx, &ports.FindSubmissionsFilter{
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

	submission, err := s.repository.Submissions.FindByID(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to retrieve submission job", "submission_id", id, "error", err)
		return err
	}

	if submission.Status == domain.SubmissionStatusAccepted || submission.Status == domain.SubmissionStatusRejected {
		s.logger.WarnContext(ctx, "skipping submission job; invalid status", "submission_id", submission.ID, "status", submission.Status)
		return nil
	}

	facts, err := s.sanitize(ctx, submission)

	if err == nil {
		submission.Accept(facts)
	} else if shouldReject(err) {
		submission.Reject(err)
	} else {
		submission.Fail(err)
	}

	if err := s.updateSubmission(ctx, submission); err != nil {
		return err
	}

	s.logger.InfoContext(ctx, "submission job recorded", "submission_id", submission.ID, "status", submission.Status)

	return nil
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
	version, err := s.repository.FormVersions.FindByID(ctx, submission.VersionID)
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
					fc := factCandidate{ftm: *ftm}
					fv, ok := submission.GetFieldValue(field.ID)

					if ok {
						fc.value = fv.Value
						fc.collectionIndex = fv.CollectionIndex
					}

					candidates = append(candidates, fc)
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
		version, err := domain.ResolveTagVersion(ta.versions)
		if err != nil {
			return nil, err
		}

		var evalFn func(domain.Tag, domain.TagVersion, []factCandidate) []*domain.CanonicalFact
		if ta.tag.HasCollectionAncestor() {
			evalFn = s.evaluateCollectionCandidates
		} else {
			evalFn = s.evaluateScalarCandidates
		}

		f := evalFn(ta.tag, *version, candidatesByVersion[version.ID])
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

func (s *submissionJobsService) updateSubmission(ctx context.Context, submission *domain.Submission) error {
	txCtx, err := s.repository.Database.BeginTx(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to begin transaction", "submission_id", submission.ID, "error", err)
		return err
	}

	defer s.repository.Database.RollbackTx(txCtx)

	_, err = s.repository.Submissions.Upsert(txCtx, submission)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to upsert submission", "submission_id", submission.ID, "error", err)
		return err
	}

	if err := s.repository.Database.CommitTx(txCtx); err != nil {
		s.logger.ErrorContext(ctx, "failed to commit transaction", "submission_id", submission.ID, "error", err)
		return err
	}

	return nil
}

func (s *submissionJobsService) getTags(ctx context.Context, ids []domain.TagVersionID) ([]tagAggregate, error) {
	versions, err := s.repository.TagVersions.FindByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	versionsByTag := make(map[domain.TagID][]*domain.TagVersion)
	for _, v := range versions {
		versionsByTag[v.TagID] = append(versionsByTag[v.TagID], v)
	}

	tags, err := s.repository.Tags.FindByIDs(ctx, slices.Collect(maps.Keys(versionsByTag)))
	if err != nil {
		return nil, err
	}

	aggregates := make([]tagAggregate, 0, len(tags))
	for _, t := range tags {
		aggregates = append(aggregates, tagAggregate{*t, versionsByTag[t.ID]})
	}

	return aggregates, nil
}

func (s *submissionJobsService) evaluateCollectionCandidates(tag domain.Tag, version domain.TagVersion, candidates []factCandidate) []*domain.CanonicalFact {
	facts := make([]*domain.CanonicalFact, 0)

	byCollectionIdx := make(map[int][]factCandidate)
	for _, fc := range candidates {
		byCollectionIdx[*fc.collectionIndex] = append(byCollectionIdx[*fc.collectionIndex], fc)
	}

	for idx, group := range byCollectionIdx {
		winner := rankCandidates(group)
		facts = append(facts, domain.NewCanonicalFact(
			winner.ftm.FieldID,
			version.ID,
			tag.KeyPath,
			winner.value,
			&idx,
		))
	}

	return facts
}

func (s *submissionJobsService) evaluateScalarCandidates(tag domain.Tag, version domain.TagVersion, candidates []factCandidate) []*domain.CanonicalFact {
	facts := make([]*domain.CanonicalFact, 0)
	winner := rankCandidates(candidates)
	facts = append(facts, domain.NewCanonicalFact(
		winner.ftm.FieldID,
		version.ID,
		tag.KeyPath,
		winner.value,
		nil,
	))

	return facts
}

func rankCandidates(candidates []factCandidate) factCandidate {
	slices.SortFunc(candidates, func(fc1, fc2 factCandidate) int {
		return cmp.Compare(fc2.ftm.Priority, fc1.ftm.Priority)
	})
	return candidates[0]
}

func shouldReject(err error) bool {
	return errors.Is(err, domain.ErrInvalidVersionStatus) ||
		errors.Is(err, strategies.ErrFieldValidation) ||
		errors.Is(err, strategies.ErrFieldRequired) ||
		errors.Is(err, strategies.ErrFieldTypeValue)
}
