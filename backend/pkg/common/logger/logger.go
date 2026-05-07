package logger

import (
	"context"
	"log/slog"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
)

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelError LogLevel = "error"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
)

func LogLevelToLevel(level string) slog.Level {
	switch LogLevel(strings.ToLower(level)) {
	case LogLevelDebug:
		return slog.LevelDebug
	case LogLevelError:
		return slog.LevelError
	case LogLevelWarn:
		return slog.LevelWarn
	case LogLevelInfo:
		fallthrough
	default:
		return slog.LevelInfo
	}
}

type RequestContextHandler struct {
	slog.Handler
}

func (h *RequestContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if id := middleware.GetReqID(ctx); id != "" {
		r.AddAttrs(slog.String("request_id", id))
	}
	return h.Handler.Handle(ctx, r)
}

func (h *RequestContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &RequestContextHandler{h.Handler.WithAttrs(attrs)}
}

func (h *RequestContextHandler) WithGroup(name string) slog.Handler {
	return &RequestContextHandler{h.Handler.WithGroup(name)}
}
