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
	"exampleapp/io/controller"
	"github.com/google/wire"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type App struct {
	homeController       *controller.HomeController
	getProductController *controller.GetProductController
}

func NewApp(
	homeController *controller.HomeController,
	getProductController *controller.GetProductController,
) *App {
	return &App{
		homeController,
		getProductController,
	}
}

func (app *App) Run() {
	http.HandleFunc("/", app.homeController.ServeHTTP)
	http.HandleFunc("/product/{id}", app.getProductController.ServeHTTP)

	http.Handle("/metrics", promhttp.Handler())

	http.ListenAndServe(":8080", nil)
}

func InitializeApp() *App {
	wire.Build(
		NewApp,
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

	return &App{}
}
