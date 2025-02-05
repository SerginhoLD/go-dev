//go:build wireinject
// +build wireinject

package main

import (
	"example.com/m/domain/eventdispatcher"
	"example.com/m/infrastructure/gookiteventdispatcher"
	"example.com/m/infrastructure/logger"
	"example.com/m/io/controller"
	"github.com/google/wire"
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

func InitializeApp() *App {
	wire.Build(
		NewApp,
		logger.NewLogger,
		logger.NewLogListener,
		gookiteventdispatcher.New,
		wire.Bind(new(eventdispatcher.EventDispatcher), new(*gookiteventdispatcher.GookitEventDispatcher)),
		controller.NewHomeController,
		controller.NewHeadersController,
	)

	return &App{}
}
