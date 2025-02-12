package usecase

import (
	"exampleapp/domain/event"
	"exampleapp/domain/eventdispatcher"
	"exampleapp/domain/repository"
)

type AllProductsUseCase struct {
	repository      repository.ProductRepository
	eventDispatcher eventdispatcher.EventDispatcher
}

func NewAllProductsUseCase(
	repository repository.ProductRepository,
	eventDispatcher eventdispatcher.EventDispatcher,
) *AllProductsUseCase {
	return &AllProductsUseCase{repository, eventDispatcher}
}

func (u *AllProductsUseCase) Handle() []*GetProductViewModel {
	products, _ := u.repository.All()
	models := make([]*GetProductViewModel, 0)

	for _, p := range products {
		models = append(models, &GetProductViewModel{
			Id:    p.Id,
			Name:  p.Name,
			Price: p.Price,
		})
	}

	u.eventDispatcher.Dispatch(&event.TestEvent{Value: "h"})

	return models
}
