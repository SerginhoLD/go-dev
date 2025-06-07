package internal

import (
	"context"
	"exampleapp/internal/domain/usecase"
)

type StartParseJob struct {
	useCase *usecase.StartParseObjectsUseCase
}

func NewStartParseJob(useCase *usecase.StartParseObjectsUseCase) *StartParseJob {
	return &StartParseJob{useCase}
}

func (job *StartParseJob) Handle(ctx context.Context) {
	job.useCase.Handle(ctx)
}
