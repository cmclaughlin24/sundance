package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"sundance/backend/pkg/common"
	"sundance/backend/pkg/common/logger"
	"sundance/backend/services/mcp-gateway/internal/mcp"
)

type settings struct {
	LogLevel  string              `json:"logLevel"`
	Transport mcp.ServerTransport `json:"transport"`
	Targets   []mcp.TargetConfig  `json:"targets"`
	Port      *int                `json:"port"`
}

func main() {
	settingsPath := flag.String("settings", "settings.json", "Path to settings JSON/YAML file")
	flag.Parse()

	var settings settings
	if err := common.ReadSettings(*settingsPath, &settings); err != nil {
		slog.Warn("failed to read settings from file; defaulting to environment variables", "error", err)
		return
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logger.LogLevelToLevel(settings.LogLevel),
	})
	l := slog.New(handler)

	loadTargets(l, settings.Targets)

	server := mcp.NewServer(
		mcp.ServerWithLogger(l),
		mcp.ServerWithName("forms-hub-server"),
		mcp.ServerWithVersion("1.0.0"),
		mcp.ServerWithTransport(settings.Transport),
		mcp.ServerWithTargets(settings.Targets...),
		mcp.ServerWithPort(settings.Port),
	)

	if err := server.Run(context.Background()); err != nil {
		l.Error("application exited with error", "error", err)
	}

	l.Info("application stopped")
}

func loadTargets(logger *slog.Logger, targets []mcp.TargetConfig) error {
	loader := mcp.NewTargetLoader(logger)

	for _, t := range targets {
		doc, err := loader.Load(context.Background(), t)
		if err != nil {
			fmt.Errorf("failed to load target %q: %w", t.Name, err)
		}

		fmt.Println(doc)
	}

	return nil
}
