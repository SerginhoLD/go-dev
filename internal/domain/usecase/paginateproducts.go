package usecase

import (
	"context"
	"exampleapp/internal/domain/repository"
)

type PaginateProductsQuery struct {
	Page  uint64
	Limit uint64
}

type PaginateProductsUseCase struct {
	repository repository.ProductRepository
}

func NewPaginateProductsUseCase(
	repository repository.ProductRepository,
) *PaginateProductsUseCase {
	return &PaginateProductsUseCase{repository}
}

func (u *PaginateProductsUseCase) Handle(ctx context.Context, query PaginateProductsQuery) []*PaginateProductViewModel {
	products, _ := u.repository.Paginate(ctx, query.Page, query.Limit)
	models := make([]*PaginateProductViewModel, 0)

	for _, p := range products {
		models = append(models, &PaginateProductViewModel{
			Id:   p.Id,
			Name: p.Name,
		})
	}

	return models
}

type PaginateProductViewModel struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}
