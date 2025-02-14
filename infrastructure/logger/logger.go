package logger

import (
	"context"
	"log/slog"
	"os"
)

func NewLogger(handler *Handler) *slog.Logger {
	return slog.New(handler)
}

type Handler struct {
	handler slog.Handler
}

func NewHandler() *Handler {
	return &Handler{
		handler: slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
	}
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	if requestId, ok := ctx.Value("X-Request-ID").(string); ok {
		r.AddAttrs(slog.String("X-Request-ID", requestId))
	}

	return h.handler.Handle(ctx, r)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.handler.WithAttrs(attrs)
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return h.handler.WithGroup(name)
}
