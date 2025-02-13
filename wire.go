//go:build wireinject
// +build wireinject

package main

import (
	"exampleapp/domain/eventdispatcher"
	"exampleapp/domain/repository"
	"exampleapp/domain/usecase"
	eventdispatcherimpl "exampleapp/infrastructure/eventdispatcher"
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
		logger.NewLogger,
		logger.NewLogListener,
		logger.NewMetricListener,
		postgres.NewConn,
		eventdispatcherimpl.New,
		wire.Bind(new(eventdispatcher.EventDispatcher), new(*eventdispatcherimpl.EventDispatcherImpl)),
		repositoryimpl.NewProductRepositoryImpl,
		wire.Bind(new(repository.ProductRepository), new(*repositoryimpl.ProductRepositoryImpl)),
		usecase.NewPaginateProductsUseCase,
		controller.NewHomeController,
		usecase.NewGetProductUseCase,
		controller.NewGetProductController,
	)

	return &io.App{}
}
