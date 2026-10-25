package spec

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sundance/backend/pkg/common/validate"
)

// Universal Standard for policy engines is the Least Privileged Algorithm:
// 1) Explicit Deny Rules: If no operation matches any pattern, it is immediately rejected, regardless of any allow rules.
// 2) Default-Allow when allow is empty: If allow is empty, all operations are permitted except those explicitly listed in deny. This makes blacklisting easy.
// 3) Default-Deny when allow is populated: If allow contains, only operations explicitly matching an allow rule (and did not match a deny rule) are permitted.

const (
	singlePathSegment string = "/*"
	doublePathSegment string = "/**"
	methodWildCard    string = "*"
)

var (
	ErrEmptyPolicyPattern     = errors.New("policy cannot be empty")
	ErrMalformedPolicyPattern = errors.New("malformed policy rule; expected '[METHOD] [PATH]' or '[PATH]'")
	ErrInvalidPolicyMethod    = errors.New("invalid HTTP method in policy rule")
)

var methodValidator = validate.NewTypeValidator([]string{
	methodWildCard,
	http.MethodGet,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodOptions,
	http.MethodHead,
})

type PolicyConfig struct {
	Allow []string `json:"allow"`
	Deny  []string `json:"deny"`
}

func (p *PolicyConfig) Validate() error {
	for _, r := range p.Deny {
		if _, _, err := parseRule(r); err != nil {
			return fmt.Errorf("invalid deny rule: %w", err)
		}
	}

	for _, r := range p.Allow {
		if _, _, err := parseRule(r); err != nil {
			return fmt.Errorf("invalid allow rule: %w", err)
		}
	}

	return nil
}

func (p *PolicyConfig) IsAllowed(method, path string) bool {
	if isDenied := hasMatch(p.Deny, method, path); isDenied {
		return false
	}

	if len(p.Allow) == 0 {
		return true
	}

	return hasMatch(p.Allow, method, path)
}

func parseRule(rule string) (string, string, error) {
	parts := strings.Fields(rule)

	if len(parts) == 0 {
		return "", "", ErrEmptyPolicyPattern
	}

	if len(parts) > 2 {
		return "", "", fmt.Errorf("%w: %q (has %d parts)", ErrMalformedPolicyPattern, rule, len(parts))
	}

	if len(parts) == 1 {
		return methodWildCard, parts[0], nil
	}

	method := strings.ToUpper(parts[0])

	if isValid := methodValidator(method); !isValid {
		return "", "", fmt.Errorf("%w: %q in rule %q", ErrInvalidPolicyMethod, parts[0], rule)
	}

	return method, parts[1], nil
}

func hasMatch(rules []string, method, path string) bool {
	for _, r := range rules {
		rMethod, rPath, err := parseRule(r)
		if err != nil {
			// NOTE: PolicyConfig.Validate is called in the composition root on start-up to ensure malformed
			//       rules are caught during start-up.
			continue
		}

		if rMethod != methodWildCard && !strings.EqualFold(rMethod, method) {
			continue
		}

		if isPathMatch(rPath, path) {
			return true
		}
	}

	return false
}

func isPathMatch(pattern, path string) bool {
	pattern = normalizePath(pattern)
	path = normalizePath(path)

	if pattern == path {
		return true
	}

	if pattern == doublePathSegment {
		return true
	}

	if strings.HasPrefix(pattern, doublePathSegment+"/") {
		if suffix := strings.TrimPrefix(pattern, doublePathSegment); strings.HasSuffix(path, suffix) {
			return true
		}
	}

	if strings.HasSuffix(pattern, doublePathSegment) {
		if prefix := strings.TrimSuffix(pattern, doublePathSegment); path == prefix || strings.HasPrefix(path, prefix+"/") {
			return true
		}
	}

	if strings.HasSuffix(pattern, singlePathSegment) {
		if prefix := strings.TrimSuffix(pattern, singlePathSegment); strings.HasPrefix(path, prefix+"/") {
			rest := strings.TrimPrefix(path, prefix+"/")

			if !strings.Contains(rest, "/") && len(rest) > 0 {
				return true
			}
		}
	}

	return false
}

func normalizePath(path string) string {
	return "/" + strings.Trim(path, "/")
}
