package app

import (
	"context"
	"log/slog"
)

type TestJob struct {
	logger *slog.Logger
}

func NewTestJob(logger *slog.Logger) *TestJob {
	return &TestJob{logger}
}

func (job *TestJob) Handler(ctx context.Context) {
	job.logger.InfoContext(ctx, "*TestJob")
}
