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

func (s *TextFieldValidatorStrategy) Validate(ctx context.Context, field domain.Field, fv domain.SubmissionFieldValue) error {
	attr, err := domain.GetFieldAttributes[domain.TextFieldAttributes](field.Attributes)
	if err != nil {
		s.logger.ErrorContext(ctx, "strategy attributes mismatch", "field_id", field.ID, "field_type", field.Type, "error", err)
		return err
	}

	val, err := checkValueRequired[string](attr, fv.Value)
	if err != nil {
		return err
	}

	// NOTE: checkValueRequired will return ("", nil) for both nil and an empty string inputs on a non-required field. Check the
	// original value to distinguish "no value submitted" (skip validation) from "empty string submitted".
	if fv.Value == nil {
		return nil
	}

	if attr.MinLength != nil && len(val) < *attr.MinLength {
		return fmt.Errorf("")
	}

	if attr.MaxLength != nil && len(val) > *attr.MaxLength {
		return fmt.Errorf("")
	}

	if attr.Pattern != "" {
		pattern, err := regexp.Compile(attr.Pattern)

		if err != nil {
			s.logger.WarnContext(ctx, "cannot validate pattern; failed to compile regex", "field_id", field.ID, "error", err, "pattern", attr.Pattern)
		} else if ok := pattern.Match([]byte(val)); !ok {
			return fmt.Errorf("")
		}
	}

	return nil
}
