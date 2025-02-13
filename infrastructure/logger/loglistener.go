package logger

import (
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

func (l *LogListener) OnUnhandledEvent(event interface{}) {
	l.logger.Error(fmt.Sprintf("Unhandled event \"%T\"", event))
}

func (l *LogListener) OnHttpResponse(event *controller.ResponseEvent) {
	l.logger.Debug(fmt.Sprintf("%s", event.Request.Pattern), slog.Int("statusCode", event.StatusCode))
}

func (l *LogListener) OnSqlQuery(event *postgres.QueryEvent) {
	l.logger.Debug(event.Query)
}

func (l *LogListener) OnTestEvent(event *event.TestEvent) error {
	l.logger.Info("ff0 Event1", "v", event.Value)
	return nil
}
