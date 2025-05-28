package repository

import (
	"context"
	"exampleapp/internal/domain/repository"
	"time"

	"github.com/redis/go-redis/v9"
)

type PlanRepositoryImpl struct {
	client *redis.Client
}

func NewPlanRepositoryImpl(
	client *redis.Client,
) repository.PlanRepository {
	return &PlanRepositoryImpl{client}
}

func (r *PlanRepositoryImpl) Running(ctx context.Context) bool {
	v, err := r.client.Get(ctx, "running").Result()

	if err == redis.Nil {
		return false
	} else if err != nil {
		panic(err)
	}

	return v == "1"
}

func (r *PlanRepositoryImpl) Run(ctx context.Context) {
	r.client.Set(ctx, "running", "1", time.Hour*4)
}

func (r *PlanRepositoryImpl) Stop(ctx context.Context) {
	r.client.Del(ctx, "running")
}
