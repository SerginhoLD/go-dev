package logger

import (
	"context"
	"exampleapp/domain/event"
	"exampleapp/infrastructure/postgres"
	"exampleapp/io/controller"
	"fmt"
	"log/slog"
)

type LogListener struct {
	logger *slog.Logger
}

func NewLogListener(logger *slog.Logger) *LogListener {
	return &LogListener{logger}
}

func (l *LogListener) OnUnhandledEvent(ctx context.Context, event interface{}) {
	l.logger.ErrorContext(ctx, fmt.Sprintf("Unhandled event \"%T\"", event))
}

func (l *LogListener) OnHttpResponse(ctx context.Context, event *controller.ResponseEvent) {
	l.logger.DebugContext(ctx, fmt.Sprintf("%s", event.Request.Pattern), slog.Int("statusCode", event.StatusCode))
}

func (l *LogListener) OnSqlQuery(ctx context.Context, event *postgres.QueryEvent) {
	l.logger.DebugContext(ctx, event.Query)
}

func (l *LogListener) OnTestEvent(ctx context.Context, event *event.TestEvent) error {
	l.logger.InfoContext(ctx, "ff0 Event1", "v", event.Value)
	return nil
}
