package strategies

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"

	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
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
	attr, err := domain.GetFieldAttributes[*domain.TextFieldAttributes](field.Attributes)
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

	if val == "" && attr.GetIsRequired() {
		return newValidationErr(field.Key, ErrFieldRequired)
	}

	if attr.MinLength != nil && len(val) < *attr.MinLength {
		return newValidationErr(field.Key, fmt.Errorf("min length"))
	}

	if attr.MaxLength != nil && len(val) > *attr.MaxLength {
		return newValidationErr(field.Key, fmt.Errorf("max length"))
	}

	if attr.Pattern != "" {
		pattern, err := regexp.Compile(attr.Pattern)

		if err != nil {
			s.logger.WarnContext(ctx, "cannot validate pattern; failed to compile regex", "field_id", field.ID, "error", err, "pattern", attr.Pattern)
			return newValidationErr(field.Key, err)
		} else if ok := pattern.Match([]byte(val)); !ok {
			return newValidationErr(field.Key, fmt.Errorf("does not match pattern %s", attr.Pattern))
		}
	}

	return nil
}
