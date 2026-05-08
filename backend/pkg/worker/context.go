package worker

import (
	"context"
	"log/slog"
)

type contextKey string

const workerIDKey contextKey = "tenantID"

func SetWorkerContext(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, workerIDKey, tenantID)
}

func WorkerFromContext(ctx context.Context) string {
	workerID, ok := ctx.Value(workerIDKey).(string)

	if !ok || workerID == "" {
		return ""
	}

	return workerID
}

type WorkerContextHandler struct {
	slog.Handler
}

func (h *WorkerContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if id := WorkerFromContext(ctx); id != "" {
		r.AddAttrs(slog.String("worker_id", id))
	}
	return h.Handler.Handle(ctx, r)
}

func (h *WorkerContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &WorkerContextHandler{h.Handler.WithAttrs(attrs)}
}

func (h *WorkerContextHandler) WithGroup(name string) slog.Handler {
	return &WorkerContextHandler{h.Handler.WithGroup(name)}
}
