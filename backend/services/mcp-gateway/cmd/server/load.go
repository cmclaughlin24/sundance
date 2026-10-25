package main

import (
	"context"
	"fmt"
	"log/slog"
	"sundance/backend/services/mcp-gateway/internal/spec"
	"sundance/backend/services/mcp-gateway/internal/tools"
)

func load(logger *slog.Logger, targets []spec.TargetConfig) (tools.Registry, error) {
	loader := spec.NewTargetLoader(logger)
	registry := make(tools.Registry)

	for _, target := range targets {
		if err := target.Policy.Validate(); err != nil {
			return nil, fmt.Errorf("target %q policy error: %w", target.Name, err)
		}

		doc, err := loader.Load(context.Background(), target)
		if err != nil {
			return nil, err
		}

		compiled := tools.ExtractTools(target, doc)
		registry.Register(compiled...)
	}

	return registry, nil
}
