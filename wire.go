//go:build wireinject
// +build wireinject

package main

import (
	"exampleapp/domain/errors"
	"exampleapp/domain/repository"
	"exampleapp/domain/usecase"
	errorsimpl "exampleapp/infrastructure/errors"
	"exampleapp/infrastructure/logger"
	"exampleapp/infrastructure/postgres"
	repositoryimpl "exampleapp/infrastructure/repository"
	"exampleapp/io"
	"exampleapp/io/controller"
	"github.com/google/wire"
)

func InitializeApp() *io.App {
	wire.Build(
		io.NewApp,
		logger.NewHandler,
		logger.NewLogger,
		logger.NewMetrics,
		errorsimpl.NewFactory,
		wire.Bind(new(errors.Factory), new(*errorsimpl.FactoryImpl)),
		postgres.NewConn,
		repositoryimpl.NewProductRepositoryImpl,
		wire.Bind(new(repository.ProductRepository), new(*repositoryimpl.ProductRepositoryImpl)),
		usecase.NewPaginateProductsUseCase,
		controller.NewHomeController,
		usecase.NewGetProductUseCase,
		controller.NewGetProductController,
	)

	return &io.App{}
}
