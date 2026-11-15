package dto

import "sundance/backend/services/forms/internal/core/domain"

type SubmissionSectionInstanceDto struct {
	Values []SubmissionValueDto `json:"values" validate:"required,dive"`
}

type SubmissionSectionDto struct {
	SectionID domain.SectionID               `json:"sectionId" validate:"required"`
	Instances []SubmissionSectionInstanceDto `json:"instances" validate:"required,dive"`
}

func FromSubmissionSectionRequests(reqs []SubmissionSectionDto) ([]*domain.SubmissionSection, error) {
	sections := make([]*domain.SubmissionSection, 0, len(reqs))
	for _, section := range reqs {
		s, err := FromSubmissionSectionRequest(section)
		if err != nil {
			return nil, err
		}

		sections = append(sections, s)
	}
	return sections, nil
}

func FromSubmissionSectionRequest(req SubmissionSectionDto) (*domain.SubmissionSection, error) {
	instances := make([]domain.SectionInstance, 0, len(req.Instances))

	for _, i := range req.Instances {
		values := make(domain.SectionInstance, 0, len(i.Values))

		for _, v := range i.Values {
			sv, err := domain.NewSubmissionValue(v.ElementID, v.Value, v.ItemIndex)
			if err != nil {
				return nil, err
			}
			values = append(values, sv)
		}

		instances = append(instances, values)
	}

	return domain.NewSubmissionSection(req.SectionID, instances)
}

func toSubmissionSectionResponse(s *domain.SubmissionSection) SubmissionSectionDto {
	instances := s.Instances()
	instanceDtos := make([]SubmissionSectionInstanceDto, 0, len(instances))

	for _, i := range instances {
		values := make([]SubmissionValueDto, 0, len(i))
		for _, value := range i {
			values = append(
				values,
				SubmissionValueDto{
					ElementID: value.ElementID,
					Value:     value.Value,
					ItemIndex: value.ItemIndex,
				},
			)
		}
		instanceDtos = append(instanceDtos, SubmissionSectionInstanceDto{values})
	}

	return SubmissionSectionDto{
		SectionID: s.SectionID,
		Instances: instanceDtos,
	}
}
