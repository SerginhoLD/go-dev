package usecase

import (
	"context"
	"exampleapp/internal/domain/errors"
	"exampleapp/internal/domain/messages"
	"exampleapp/internal/domain/messenger"
	"exampleapp/internal/domain/repository"
	"math"
)

type StartParseObjectsUseCase struct {
	planRepository     repository.PlanRepository
	veObjectRepository repository.VeObjectRepository
	bus                messenger.Bus
}

func NewStartParseObjectsUseCase(
	planRepository repository.PlanRepository,
	veObjectRepository repository.VeObjectRepository,
	bus messenger.Bus,
) *StartParseObjectsUseCase {
	return &StartParseObjectsUseCase{
		planRepository:     planRepository,
		veObjectRepository: veObjectRepository,
		bus:                bus,
	}
}

func (u *StartParseObjectsUseCase) Handle(ctx context.Context) error {
	if u.planRepository.Running(ctx) {
		return errors.NewContext(ctx, "StartParseObjects: is running")
	}

	total, size := u.veObjectRepository.Total(ctx)

	if size == 0 {
		return errors.NewContext(ctx, "StartParseObjects: nil size")
	}

	// todo: defer Stop
	pages := uint64(math.Ceil(float64(total) / float64(size)))

	for page := range pages {
		u.bus.Send(ctx, &messages.ParsePageMessage{Page: page + 1})
	}

	u.planRepository.Run(ctx)
	return nil
}
