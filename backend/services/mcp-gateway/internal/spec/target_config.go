package spec

import (
	"errors"
	"time"
)

var ErrInvalidTargetConfig = errors.New("invalid target configuration")

type TargetConfig struct {
	Name    string            `json:"name"`
	BaseURL string            `json:"baseUrl"`
	Spec    SpecSourceConfig  `json:"spec"`
	Headers map[string]string `json:"headers"`
	Timeout time.Duration     `json:"timeout"`
	Policy  PolicyConfig      `json:"policy"`
}

type SpecSourceConfig struct {
	Version string `json:"version"`
	Path    string `json:"path"`
	URL     string `json:"url"`
}
