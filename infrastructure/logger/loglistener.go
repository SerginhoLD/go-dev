package logger

import (
	"exampleapp/domain/event"
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

func (l *LogListener) OnEvent1(event *event.TestEvent) error {
	l.logger.Info("ff0 Event1", "v", event.Value)
	return nil
}
