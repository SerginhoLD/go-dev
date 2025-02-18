package usecase

import (
	"context"
	"exampleapp/domain/entity"
	"exampleapp/domain/errors"
	"exampleapp/domain/repository"
	"exampleapp/domain/validator"
)

type CreateProductCommand struct {
	Name  string  `json:"name" validate:"required,min=3,max=255"`
	Price float64 `json:"price" validate:"required,gte=1"`
}

type CreateProductUseCase struct {
	errFactory errors.Factory
	validator  validator.Validator
	repository repository.ProductRepository
}

func NewCreateProductUseCase(errFactory errors.Factory, validator validator.Validator, repository repository.ProductRepository) *CreateProductUseCase {
	return &CreateProductUseCase{errFactory, validator, repository}
}

func (u *CreateProductUseCase) Handle(ctx context.Context, command CreateProductCommand) (*CreateProductViewModel, error) {
	err := u.validator.Validate(command)

	if err != nil {
		return nil, u.errFactory.WrapContext(ctx, "CreateProduct: %w", err)
	}

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
