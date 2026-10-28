package spec

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sundance/backend/pkg/common/stratreg"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
)

type parserFn func(context.Context, []byte) (*openapi3.T, error)

type TargetLoader struct {
	logger  *slog.Logger
	parsers stratreg.StrategyRegistry[string, parserFn]
}

func NewTargetLoader(logger *slog.Logger) *TargetLoader {
	parsers := stratreg.New[string, parserFn]().
		Set("2.0", parseV2).
		Set("3.0", parseV3x).
		Set("3.1", parseV3x)

	return &TargetLoader{
		logger:  logger,
		parsers: parsers,
	}
}

func (p *TargetLoader) Load(ctx context.Context, target TargetConfig) (*openapi3.T, error) {
	spec, err := p.readSpecBytes(ctx, target)
	if err != nil {
		return nil, err
	}

	parser, err := p.parsers.Get(target.Spec.Version)
	if err != nil {
		return nil, err
	}

	return parser(ctx, spec)
}

func (p *TargetLoader) readSpecBytes(ctx context.Context, target TargetConfig) ([]byte, error) {
	if strings.Trim(target.Spec.Path, " ") != "" {
		return os.ReadFile(target.Spec.Path)
	}

	if strings.Trim(target.Spec.URL, " ") != "" {
		return p.requestFile(ctx, target)
	}

	return nil, fmt.Errorf("target %q has neither spec.Path nor spec.URL configured: %w", target.Name, ErrInvalidTargetConfig)
}

func (p *TargetLoader) requestFile(ctx context.Context, target TargetConfig) ([]byte, error) {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	if target.Timeout > 0 {
		client.Timeout = target.Timeout
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target.Spec.URL, nil)
	if err != nil {
		return nil, ErrInvalidTargetConfig
	}

	for k, v := range target.Headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to target %q failed: %w", target.Name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("target %q spec endpoint returned status %d", target.Name, resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
