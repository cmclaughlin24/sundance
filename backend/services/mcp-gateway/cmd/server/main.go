package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"sundance/backend/pkg/common"
	"sundance/backend/pkg/common/logger"
	"sundance/backend/services/mcp-gateway/internal/mcp"
	"sundance/backend/services/mcp-gateway/internal/spec"
	"sundance/backend/services/mcp-gateway/internal/tools"
)

type settings struct {
	Name      string              `json:"name"`
	Version   string              `json:"version"`
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

	r, err := load(l, settings.Targets)
	if err != nil {
		panic(err)
	}

	httpDispatcher := tools.NewHTTPDispatcher(&http.Client{}, l)

	server := mcp.NewServer(
		mcp.ServerWithLogger(l),
		mcp.ServerWithName(settings.Name),
		mcp.ServerWithVersion(settings.Version),
		mcp.ServerWithTransport(settings.Transport),
		mcp.ServerWithPort(settings.Port),
		mcp.ServerWithToolRegistry(r),
		mcp.ServerWithDispatcher(httpDispatcher),
	)

	if err := server.Run(context.Background()); err != nil {
		l.Error("application exited with error", "error", err)
	}

	l.Info("application stopped")
}
