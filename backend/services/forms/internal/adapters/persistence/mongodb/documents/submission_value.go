package documents

import "sundance/backend/services/forms/internal/core/domain"

type submissionValueDocument struct {
	ElementID string `bson:"element_id"`
	Value     any    `bson:"value"`
	ItemIndex *int   `bson:"item_index,omitempty"`
}

func toSubmissionValueDocument(sv *domain.SubmissionValue) submissionValueDocument {
	return submissionValueDocument{
		ElementID: string(sv.ElementID),
		Value:     sv.Value,
		ItemIndex: sv.ItemIndex,
	}
}

func fromSubmissionValueDocument(doc *submissionValueDocument) *domain.SubmissionValue {
	return domain.HydrateSubmissionValue(domain.ElementID(doc.ElementID), doc.Value, doc.ItemIndex)
}
