//go:build wireinject
// +build wireinject

package main

import (
	"exampleapp/domain/eventdispatcher"
	eventdispatcherimpl "exampleapp/infrastructure/eventdispatcher"
	"exampleapp/infrastructure/logger"
	"exampleapp/infrastructure/postgres"
	"exampleapp/io/controller"
	"github.com/google/wire"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type App struct {
	homeController    *controller.HomeController
	headersController *controller.HeadersController
}

func NewApp(
	homeController *controller.HomeController,
	headersController *controller.HeadersController,
) *App {
	return &App{
		homeController,
		headersController,
	}
}

func (app *App) Run() {
	http.HandleFunc("/hello", app.homeController.ServeHTTP)
	http.HandleFunc("/headers", app.headersController.ServeHTTP)

	http.Handle("/metrics", promhttp.Handler())

	http.ListenAndServe(":8080", nil)
}

func InitializeApp() *App {
	wire.Build(
		NewApp,
		logger.NewLogger,
		logger.NewLogListener,
		logger.NewMetricListener,
		postgres.NewDB,
		eventdispatcherimpl.New,
		wire.Bind(new(eventdispatcher.EventDispatcher), new(*eventdispatcherimpl.EventDispatcherImpl)),
		controller.NewHomeController,
		controller.NewHeadersController,
	)

	return &App{}
}
