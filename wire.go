//go:build wireinject
// +build wireinject

package main

import (
	"example.com/m/domain/eventdispatcher"
	eventdispatcherimpl "example.com/m/infrastructure/eventdispatcher"
	"example.com/m/infrastructure/logger"
	"example.com/m/io/controller"
	"github.com/google/wire"
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

	http.ListenAndServe(":8080", nil)
}

func InitializeApp() *App {
	wire.Build(
		NewApp,
		logger.NewLogger,
		logger.NewLogListener,
		eventdispatcherimpl.New,
		wire.Bind(new(eventdispatcher.EventDispatcher), new(*eventdispatcherimpl.EventDispatcherImpl)),
		controller.NewHomeController,
		controller.NewHeadersController,
	)

	return &App{}
}
