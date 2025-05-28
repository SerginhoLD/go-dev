package logger

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"os"

	"github.com/google/uuid"
)

func init() {
	slog.SetDefault(slog.New(NewHandler(os.Stderr)))
}

type Handler struct {
	handler slog.Handler
}

func NewHandler(writer io.Writer) *Handler {
	return &Handler{
		handler: slog.NewJSONHandler(writer, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
	}
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = errors.New("!LOG")
			// логирование не должно приводить к падению приложения
		}
	}()

	record := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
	contextMap := make(map[string]any)

	r.Attrs(func(attr slog.Attr) bool {
		if "err" == attr.Key && slog.KindString == attr.Value.Kind() {
			record.AddAttrs(attr)
		} else if slog.KindGroup == attr.Value.Kind() {
			contextMap[attr.Key] = h.handleGroup(attr.Value.Group())
		} else {
			contextMap[attr.Key] = attr.Value.Any()
		}

		return true
	})

	if len(contextMap) > 0 {
		contextBytes, _ := json.Marshal(contextMap)
		record.AddAttrs(slog.String("context", string(contextBytes)))
	}

	if requestId, ok := ctx.Value("X-Request-ID").(string); ok {
		record.AddAttrs(slog.String("X-Request-ID", requestId))
	}

	id, _ := uuid.NewV7()
	record.AddAttrs(slog.String("id", id.String()))

	return h.handler.Handle(ctx, record)
}

func (h *Handler) handleGroup(attrs []slog.Attr) map[string]any {
	contextMap := make(map[string]any)

	for _, attr := range attrs {
		if slog.KindGroup == attr.Value.Kind() {
			contextMap[attr.Key] = h.handleGroup(attr.Value.Group())
		} else {
			contextMap[attr.Key] = attr.Value.Any()
		}
	}

	return contextMap
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.handler.WithAttrs(attrs)
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return h.handler.WithGroup(name)
}
