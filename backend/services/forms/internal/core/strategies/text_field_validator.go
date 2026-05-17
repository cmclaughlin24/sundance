package strategies

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"

	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/domain"
	"github.com/cmclaughlin24/sundance/backend/services/forms/internal/core/ports"
)

type TextFieldValidatorStrategy struct {
	logger *slog.Logger
}

func NewTextFieldValidatorStrategy(logger *slog.Logger) ports.FieldValidatorStrategy {
	return &TextFieldValidatorStrategy{
		logger: logger,
	}
}

func (s *TextFieldValidatorStrategy) Validate(ctx context.Context, field domain.Field, value domain.SubmissionFieldValue) error {
	attr, err := domain.GetFieldAttributes[domain.TextFieldAttributes](field.Attributes)
	if err != nil {
		s.logger.ErrorContext(ctx, "strategy attributes mismatch", "field_id", field.ID, "field_type", field.Type, "error", err)
		return err
	}

	// TODO: Implement validation.

	if attr.MinLength != nil {
		return fmt.Errorf("")
	}

	if attr.MaxLength != nil {
		return fmt.Errorf("")
	}

	if len(attr.Pattern) > 0 {
		pattern, err := regexp.Compile(attr.Pattern)

		if err != nil {
			s.logger.WarnContext(ctx, "cannot validate pattern; failed to compile regex", "field_id", field.ID, "error", err, "pattern", attr.Pattern)
		} else if ok := pattern.Match([]byte("")); !ok {
			return fmt.Errorf("")
		}
	}

	return nil
}
