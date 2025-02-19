//go:build wireinject
// +build wireinject

package main

import (
	"exampleapp/cmd/web/internal/app"
	"exampleapp/internal/domain/errors"
	"exampleapp/internal/domain/repository"
	"exampleapp/internal/domain/usecase"
	"exampleapp/internal/domain/validator"
	errorsimpl "exampleapp/internal/infrastructure/errors"
	"exampleapp/internal/infrastructure/logger"
	"exampleapp/internal/infrastructure/postgres"
	repositoryimpl "exampleapp/internal/infrastructure/repository"
	validatorimpl "exampleapp/internal/infrastructure/validator"
	"github.com/google/wire"
	"io"
	"log/slog"
	"os"
)

func InitializeApp() *app.App {
	wire.Build(
		app.New,
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
		app.NewCoverageController,
		usecase.NewPaginateProductsUseCase,
		app.NewHomeController,
		usecase.NewGetProductUseCase,
		app.NewGetProductController,
		usecase.NewCreateProductUseCase,
		app.NewCreateProductController,
	)

	return &app.App{}
}
