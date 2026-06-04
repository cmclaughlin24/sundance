package strategies

import (
	"context"
	"log/slog"
	"sundance/backend/services/tenants/internal/core/domain"
)

type baseLookupStrategy struct {
	logger *slog.Logger
}

func (s *baseLookupStrategy) toLookups(ctx context.Context, rows []map[string]any, valueField, labelField string, dataSourceID domain.DataSourceID) []*domain.Lookup {
	lookups := make([]*domain.Lookup, 0, len(rows))

	for i, row := range rows {
		value, ok := row[valueField]
		if !ok {
			s.logger.WarnContext(ctx, "skipping lookup row missing value field", "data_source_id", dataSourceID, "row_index", i, "value_field", valueField)
			continue
		}

		label, ok := row[labelField]
		if !ok {
			s.logger.WarnContext(ctx, "skipping lookup row missing label field", "data_source_id", dataSourceID, "row_index", i, "label_field", labelField)
			continue
		}

		switch l := label.(type) {
		case string:
			lookups = append(lookups, domain.NewLookup(value, l))
		default:
			s.logger.WarnContext(ctx, "skipping lookup row label field is not string", "data_source_id", dataSourceID, "row_index", i, "label_field", labelField)
		}
	}
	return lookups
}

func (s baseLookupStrategy) missingRequiredKeys(required []string, params map[string]any) []string {
	if len(required) == 0 {
		return nil
	}

	var missing []string

	for _, key := range required {
		v, ok := params[key]
		if !ok || v == nil {
			missing = append(missing, key)
		}
	}

	return missing
}
