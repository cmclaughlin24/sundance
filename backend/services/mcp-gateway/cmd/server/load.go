package main

import (
	"context"
	"log/slog"
	"sundance/backend/services/mcp-gateway/internal/spec"
	"sundance/backend/services/mcp-gateway/internal/tools"
)

func load(logger *slog.Logger, targets []spec.TargetConfig) (tools.Registry, error) {
	loader := spec.NewTargetLoader(logger)
	registry := make(tools.Registry)

	for _, target := range targets {
		doc, err := loader.Load(context.Background(), target)
		if err != nil {
			return nil, err
		}

		compiled, err := tools.ExtractTools(target, doc)
		if err != nil {
			return nil, err
		}

		registry.Register(compiled...)
	}

	return registry, nil
}
