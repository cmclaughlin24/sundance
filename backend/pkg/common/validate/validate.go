package validate

import (
	"errors"
	"regexp"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var v *validator.Validate

func init() {
	v = validator.New(validator.WithRequiredStructEnabled())
	v.RegisterValidation("nowhitespace", noWhitespace)
	v.RegisterValidation("uuidv7", isUUIDV7)
}

func IsValidationErr(err error) bool {
	return errors.As(err, &validator.ValidationErrors{})
}

func ValidateStruct(s any) error {
	return v.Struct(s)
}

func RegisterValidation(tag string, fn validator.Func) {
	v.RegisterValidation(tag, fn)
}

func NewTypeValidator[T comparable](types []T) func(T) bool {
	cpy := slices.Clone(types)

	return func(t T) bool {
		return slices.Contains(cpy, t)
	}
}

func NewRegexValidator(pattern string) validator.Func {
	re := regexp.MustCompile(pattern)

	return func(fl validator.FieldLevel) bool {
		return re.MatchString(fl.Field().String())
	}
}

func noWhitespace(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	return strings.TrimSpace(val) != ""
}

func isUUIDV7(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	parsed, err := uuid.Parse(val)
	return err == nil && parsed.Version() == 7
}
