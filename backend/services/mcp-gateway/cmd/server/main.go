package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"sundance/backend/pkg/common"
	"sundance/backend/pkg/common/logger"
	"sundance/backend/services/mcp-gateway/internal/mcp"
	"sundance/backend/services/mcp-gateway/internal/spec"
)

type settings struct {
	LogLevel  string              `json:"logLevel"`
	Transport mcp.ServerTransport `json:"transport"`
	Targets   []spec.TargetConfig `json:"targets"`
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

	tools, err := load(l, settings.Targets)
	if err != nil {
		panic(err)
	}

	server := mcp.NewServer(
		mcp.ServerWithLogger(l),
		mcp.ServerWithName("forms-hub-server"),
		mcp.ServerWithVersion("1.0.0"),
		mcp.ServerWithTransport(settings.Transport),
		mcp.ServerWithPort(settings.Port),
		mcp.ServerWithToolRegistry(tools),
	)

	if err := server.Run(context.Background()); err != nil {
		l.Error("application exited with error", "error", err)
	}

	l.Info("application stopped")
}
