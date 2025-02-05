package logger

import "log/slog"

type LogListener struct {
	logger *slog.Logger
}

func NewLogListener(logger *slog.Logger) *LogListener {
	return &LogListener{logger}
}

func (l *LogListener) OnEvent1(data map[string]any) {
	l.logger.Info("ff9 Event1", "arg0", data["arg0"])
}
