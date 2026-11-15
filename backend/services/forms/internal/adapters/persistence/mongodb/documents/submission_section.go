package documents

import "sundance/backend/services/forms/internal/core/domain"

type submissionInstanceDocument struct {
	Values []submissionValueDocument `bson:"values"`
}

type submissionSectionDocument struct {
	SectionID string                       `bson:"section_id"`
	Instances []submissionInstanceDocument `bson:"instances"`
}

func toSubmissionSectionDocument(s *domain.SubmissionSection) *submissionSectionDocument {
	instances := s.Instances()
	instanceDocs := make([]submissionInstanceDocument, 0, len(instances))

	for _, i := range instances {
		vals := make([]submissionValueDocument, 0, len(i))
		for _, v := range i {
			vals = append(vals, toSubmissionValueDocument(v))
		}
		instanceDocs = append(instanceDocs, submissionInstanceDocument{vals})
	}

	return &submissionSectionDocument{
		SectionID: string(s.SectionID),
		Instances: instanceDocs,
	}
}

func fromSubmissionSectionDocument(doc *submissionSectionDocument) *domain.SubmissionSection {
	instances := make([]domain.SectionInstance, 0, len(doc.Instances))
	for _, i := range doc.Instances {
		vals := make(domain.SectionInstance, 0, len(i.Values))
		for _, val := range i.Values {
			vals = append(vals, fromSubmissionValueDocument(&val))
		}
		instances = append(instances, vals)
	}

	return domain.HydrateSubmissionSection(domain.SectionID(doc.SectionID), instances)
}
