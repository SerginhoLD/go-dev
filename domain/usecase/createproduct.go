package usecase

import (
	"context"
	"exampleapp/domain/entity"
	"exampleapp/domain/errors"
	"exampleapp/domain/repository"
)

type CreateProductCommand struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type CreateProductUseCase struct {
	errFactory errors.Factory
	repository repository.ProductRepository
}

func NewCreateProductUseCase(errFactory errors.Factory, repository repository.ProductRepository) *CreateProductUseCase {
	return &CreateProductUseCase{errFactory, repository}
}

func (u *CreateProductUseCase) Handle(ctx context.Context, command CreateProductCommand) (*CreateProductViewModel, error) {
	product := u.repository.FindByName(ctx, command.Name)

	if product != nil {
		return nil, u.errFactory.NewContext(ctx, "CreateProduct: product already exists")
	}

	product = new(entity.Product)
	product.Name = command.Name
	product.Price = command.Price

	u.repository.Create(ctx, product)

	return &CreateProductViewModel{
		Id: product.Id,
	}, nil
}

type CreateProductViewModel struct {
	Id uint64 `json:"id"`
}
