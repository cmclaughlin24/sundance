package processors

import (
	"cmp"
	"context"
	"errors"
	"log/slog"
	"maps"
	"slices"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

var (
	ErrMissingCollectionIndex = errors.New("missing collection index")
)

type candidate struct {
	ftm   *domain.FieldTagMapping
	value *domain.SubmissionFieldValue
}

type tagAggregate struct {
	tag      domain.Tag
	versions []*domain.TagVersion
}

type submissionNormalizer struct {
	logger               *slog.Logger
	tagsRepository       ports.TagsRepository
	tagVersionRepository ports.TagVersionsRepository
}

func newSubmissionNormalizer(logger *slog.Logger, repository *ports.Repository) *submissionNormalizer {
	return &submissionNormalizer{
		logger:               logger,
		tagsRepository:       repository.Tags,
		tagVersionRepository: repository.TagVersions,
	}
}

func (n *submissionNormalizer) normalize(ctx context.Context, resolved []resolveField) ([]*domain.CanonicalFact, error) {
	candidatesByVersion := make(map[domain.TagVersionID][]candidate)
	for _, rf := range resolved {
		for _, t := range rf.field.GetTags() {
			candidatesByVersion[t.TagVersionID] = append(candidatesByVersion[t.TagVersionID], candidate{
				ftm:   t,
				value: rf.value,
			})
		}
	}

	tags, err := n.getTags(ctx, slices.Collect(maps.Keys(candidatesByVersion)))
	if err != nil {
		return nil, err
	}

	facts := make([]*domain.CanonicalFact, 0)
	for _, ta := range tags {
		version, err := domain.ResolveTagVersion(ta.versions)
		if err != nil {
			return nil, err
		}

		var evalFn func(domain.Tag, domain.TagVersion, []candidate) ([]*domain.CanonicalFact, error)
		if ta.tag.HasCollectionAncestor() {
			evalFn = n.evaluateCollectionCandidates
		} else {
			evalFn = n.evaluateScalarCandidates
		}

		f, err := evalFn(ta.tag, *version, candidatesByVersion[version.ID])
		if err != nil {
			return nil, err
		}

		facts = append(facts, f...)
	}

	return facts, nil
}

func (n *submissionNormalizer) getTags(ctx context.Context, ids []domain.TagVersionID) ([]tagAggregate, error) {
	versions, err := n.tagVersionRepository.FindByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	versionsByTag := make(map[domain.TagID][]*domain.TagVersion)
	for _, v := range versions {
		versionsByTag[v.TagID] = append(versionsByTag[v.TagID], v)
	}

	tags, err := n.tagsRepository.FindByIDs(ctx, slices.Collect(maps.Keys(versionsByTag)))
	if err != nil {
		return nil, err
	}

	aggregates := make([]tagAggregate, 0, len(tags))
	for _, t := range tags {
		aggregates = append(aggregates, tagAggregate{*t, versionsByTag[t.ID]})
	}

	return aggregates, nil
}

func (n *submissionNormalizer) evaluateCollectionCandidates(tag domain.Tag, version domain.TagVersion, candidates []candidate) ([]*domain.CanonicalFact, error) {
	facts := make([]*domain.CanonicalFact, 0)

	byCollectionIdx := make(map[int][]candidate)
	for _, c := range candidates {
		value := c.value

		if value == nil || value.CollectionIndex == nil {
			return nil, ErrMissingCollectionIndex
		}

		byCollectionIdx[*value.CollectionIndex] = append(byCollectionIdx[*value.CollectionIndex], c)
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

	return facts, nil
}

func (n *submissionNormalizer) evaluateScalarCandidates(tag domain.Tag, version domain.TagVersion, candidates []candidate) ([]*domain.CanonicalFact, error) {
	facts := make([]*domain.CanonicalFact, 0)
	winner := rankCandidates(candidates)

	var value any
	if winner.value != nil {
		value = winner.value.Value
	}

	facts = append(facts, domain.NewCanonicalFact(
		winner.ftm.FieldID,
		version.ID,
		tag.KeyPath,
		value,
		nil,
	))

	return facts, nil
}

func rankCandidates(candidates []candidate) candidate {
	slices.SortFunc(candidates, func(fc1, fc2 candidate) int {
		return cmp.Compare(fc2.ftm.Priority, fc1.ftm.Priority)
	})
	return candidates[0]
}
