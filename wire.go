//go:build wireinject
// +build wireinject

package main

import (
	"exampleapp/domain/errors"
	"exampleapp/domain/repository"
	"exampleapp/domain/usecase"
	"exampleapp/domain/validator"
	errorsimpl "exampleapp/infrastructure/errors"
	"exampleapp/infrastructure/logger"
	"exampleapp/infrastructure/postgres"
	repositoryimpl "exampleapp/infrastructure/repository"
	validatorimpl "exampleapp/infrastructure/validator"
	appio "exampleapp/io"
	"exampleapp/io/controller"
	"github.com/google/wire"
	"io"
	"log/slog"
	"os"
)

func InitializeApp() *appio.App {
	wire.Build(
		appio.NewApp,
		wire.NewSet(
			slog.New,
			logger.NewHandler,
			wire.Bind(new(slog.Handler), new(*logger.Handler)),
			wire.InterfaceValue(new(io.Writer), os.Stderr),
		),
		logger.NewMetrics,
		errorsimpl.NewFactory,
		wire.Bind(new(errors.Factory), new(*errorsimpl.FactoryImpl)),
		validatorimpl.New,
		wire.Bind(new(validator.Validator), new(*validatorimpl.ValidatorImpl)),
		postgres.NewConn,
		repositoryimpl.NewProductRepositoryImpl,
		wire.Bind(new(repository.ProductRepository), new(*repositoryimpl.ProductRepositoryImpl)),
		controller.NewCoverageController,
		usecase.NewPaginateProductsUseCase,
		controller.NewHomeController,
		usecase.NewGetProductUseCase,
		controller.NewGetProductController,
		usecase.NewCreateProductUseCase,
		controller.NewCreateProductController,
	)

	return &appio.App{}
}
