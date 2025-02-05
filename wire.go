//go:build wireinject
// +build wireinject

package main

import (
	"example.com/m/io/controller"
	"github.com/google/wire"
	"log/slog"
	"os"
)

func InitializeHomeController() *controller.HomeController {
	wire.Build(
		InitializeLogger,
		controller.NewHomeController,
	)

	return &controller.HomeController{}
}

func InitializeHeadersController() *controller.HeadersController {
	wire.Build(
		InitializeLogger,
		controller.NewHeadersController,
	)

	return &controller.HeadersController{}
}

func InitializeLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stderr, nil))
}
