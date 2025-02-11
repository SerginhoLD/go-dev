package logger

import (
	"exampleapp/domain/event"
	"log/slog"
)

type MetricListener struct {
	logger *slog.Logger
}

func NewMetricListener(logger *slog.Logger) *MetricListener {
	return &MetricListener{logger}
}

func (l *MetricListener) OnEvent2(event *event.TestEvent) {
	l.logger.Info("todo: metrics")
}
