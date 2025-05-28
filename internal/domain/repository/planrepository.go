package repository

import "context"

type PlanRepository interface {
	Running(ctx context.Context) bool
	Run(ctx context.Context)
	Stop(ctx context.Context)
}
