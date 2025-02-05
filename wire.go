//go:build wireinject
// +build wireinject

package main

import (
	"example.com/m/domain/eventdispatcher"
	"example.com/m/io/controller"
	"github.com/google/wire"
	"log/slog"
	"os"
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
		InitializeLogger,
		eventdispatcher.NewEventDispatcher,
		controller.NewHomeController,
		controller.NewHeadersController,
	)

	return &App{}
}

func InitializeLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stderr, nil))
}
