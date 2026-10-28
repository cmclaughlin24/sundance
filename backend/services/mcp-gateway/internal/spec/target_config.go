package spec

import (
	"errors"
	"strings"
	"time"
)

type HeaderOverridePolicy string

const (
	HeaderOverridePolicyAllow     HeaderOverridePolicy = "allow"
	HeaderOverridePolicyForbidden HeaderOverridePolicy = "forbidden"
)

var ErrInvalidTargetConfig = errors.New("invalid target configuration")

type TargetConfig struct {
	Name    string           `json:"name"`
	BaseURL string           `json:"baseUrl"`
	Spec    SpecSourceConfig `json:"spec"`
	Headers HeadersConfig    `json:"headers"`
	Timeout time.Duration    `json:"timeout"`
	Policy  PolicyConfig     `json:"policy"`
}

type SpecSourceConfig struct {
	Version string `json:"version"`
	Path    string `json:"path"`
	URL     string `json:"url"`
}

type Header struct {
	Value          any                  `json:"value"`
	OverridePolicy HeaderOverridePolicy `json:"overridePolicy"`
}

type HeadersConfig map[string]Header

func (hc HeadersConfig) AllowHeader(key string) bool {
	if len(hc) == 0 {
		return true
	}

	for header, config := range hc {
		if strings.ToLower(header) != strings.ToLower(key) {
			continue
		}

		if config.OverridePolicy == HeaderOverridePolicyForbidden {
			return false
		}

		if config.OverridePolicy == HeaderOverridePolicyAllow {
			return true
		}
	}

	return true
}
