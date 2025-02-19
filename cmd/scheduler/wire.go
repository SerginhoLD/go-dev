//go:build wireinject
// +build wireinject

package main

import (
	"exampleapp/cmd/scheduler/internal/app"
	"exampleapp/internal/infrastructure/logger"
	"exampleapp/internal/infrastructure/postgres"
	"github.com/google/wire"
	"io"
	"log/slog"
	"os"
)

func InitializeScheduler() *app.Scheduler {
	wire.Build(
		app.New,
		wire.NewSet(
			slog.New,
			logger.NewHandler,
			wire.Bind(new(slog.Handler), new(*logger.Handler)),
			wire.InterfaceValue(new(io.Writer), os.Stderr),
		),
		postgres.NewConn,
		app.NewTestJob,
	)

	return &app.Scheduler{}
}
