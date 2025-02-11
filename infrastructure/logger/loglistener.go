package logger

import (
	"exampleapp/domain/event"
	"log/slog"
)

type LogListener struct {
	logger *slog.Logger
}

func NewLogListener(logger *slog.Logger) *LogListener {
	return &LogListener{logger}
}

func (l *LogListener) OnEvent1(event *event.TestEvent) {
	l.logger.Info("ff0 Event1", "v", event.Value)
}
