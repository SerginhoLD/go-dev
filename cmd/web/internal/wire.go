//go:build wireinject
// +build wireinject

package internal

import (
	"exampleapp/internal/infrastructure/di"

	"github.com/google/wire"
)

func InitializeApp() *App {
	wire.Build(
		NewApp,
		di.Set,
		NewCoverageController,
		NewHomeController,
	)

	return &App{}
}
