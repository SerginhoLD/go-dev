package usecase

import (
	"exampleapp/domain/repository"
)

type GetProductQuery struct {
	Id uint64
}

type GetProductUseCase struct {
	repository repository.ProductRepository
}

func NewGetProductUseCase(repository repository.ProductRepository) *GetProductUseCase {
	return &GetProductUseCase{repository}
}

func (u *GetProductUseCase) Handle(query GetProductQuery) *GetProductViewModel {
	p := u.repository.Find(query.Id)

	if p == nil {
		return nil
	}

	return &GetProductViewModel{
		Id:    p.Id,
		Name:  p.Name,
		Price: p.Price,
	}
}

type GetProductViewModel struct {
	Id    uint64  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
