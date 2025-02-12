package usecase

import (
	"exampleapp/domain/event"
	"exampleapp/domain/eventdispatcher"
	"exampleapp/domain/repository"
)

type PaginateProductsQuery struct {
	Page  uint64
	Limit uint64
}

type PaginateProductsUseCase struct {
	repository      repository.ProductRepository
	eventDispatcher eventdispatcher.EventDispatcher
}

func NewPaginateProductsUseCase(
	repository repository.ProductRepository,
	eventDispatcher eventdispatcher.EventDispatcher,
) *PaginateProductsUseCase {
	return &PaginateProductsUseCase{repository, eventDispatcher}
}

func (u *PaginateProductsUseCase) Handle(query PaginateProductsQuery) []*GetProductViewModel {
	products, _ := u.repository.Paginate(query.Page, query.Limit)
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
