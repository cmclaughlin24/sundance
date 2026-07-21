package documents

import "sundance/backend/services/forms/internal/core/domain"

type canonicalFact struct {
	ElementID       string `bson:"element_id"`
	TagVersionID    string `bson:"tag_version_id"`
	TagKey          string `bson:"tag_key"`
	Value           any    `bson:"value"`
	CollectionIndex *int   `bson:"collection_index,omitempty"`
}

func toCanonicalFactDocument(cf *domain.CanonicalFact) *canonicalFact {
	return &canonicalFact{
		ElementID:       string(cf.ElementID),
		TagVersionID:    string(cf.TagVersionID),
		TagKey:          cf.TagKey,
		Value:           cf.Value,
		CollectionIndex: cf.CollectionIndex,
	}
}

func fromCanonicalFactDocument(doc *canonicalFact) *domain.CanonicalFact {
	return domain.HydrateCanonicalFact(
		domain.ElementID(doc.ElementID),
		domain.TagVersionID(doc.TagVersionID),
		doc.TagKey,
		doc.Value,
		doc.CollectionIndex,
	)
}
