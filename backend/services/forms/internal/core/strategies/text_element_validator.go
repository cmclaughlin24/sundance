package strategies

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"

	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"
)

type TextElementValidatorStrategy struct {
	logger *slog.Logger
}

func NewTextElementValidatorStrategy(logger *slog.Logger) ports.ElementValidatorStrategy {
	return &TextElementValidatorStrategy{
		logger: logger,
	}
}

func (s *TextElementValidatorStrategy) Validate(ctx context.Context, element domain.Element, sv domain.SubmissionValue) error {
	attr, err := domain.GetElementAttributes[*domain.TextElementAttributes](element.Attributes)
	if err != nil {
		s.logger.ErrorContext(ctx, "strategy attributes mismatch", "element_id", element.ID, "element_type", element.Type, "error", err)
		return err
	}

	val, err := checkValueRequired[string](attr, sv.Value)
	if err != nil {
		return err
	}

	// NOTE: checkValueRequired will return ("", nil) for both nil and an empty string inputs on a non-required element. Check the
	// original value to distinguish "no value submitted" (skip validation) from "empty string submitted".
	if sv.Value == nil {
		return nil
	}

	if val == "" && attr.GetIsRequired() {
		return newValidationErr(element.Key, ErrElementRequired)
	}

	if attr.MinLength != nil && len(val) < *attr.MinLength {
		return newValidationErr(element.Key, fmt.Errorf("min length"))
	}

	if attr.MaxLength != nil && len(val) > *attr.MaxLength {
		return newValidationErr(element.Key, fmt.Errorf("max length"))
	}

	if attr.Pattern != "" {
		pattern, err := regexp.Compile(attr.Pattern)

		if err != nil {
			s.logger.WarnContext(ctx, "cannot validate pattern; failed to compile regex", "element_id", element.ID, "error", err, "pattern", attr.Pattern)
			return newValidationErr(element.Key, err)
		} else if ok := pattern.Match([]byte(val)); !ok {
			return newValidationErr(element.Key, fmt.Errorf("does not match pattern %s", attr.Pattern))
		}
	}

	return nil
}
