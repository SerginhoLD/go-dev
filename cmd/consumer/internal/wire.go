//go:build wireinject
// +build wireinject

package internal

import (
	"exampleapp/internal/infrastructure/di"

	"github.com/google/wire"
)

func InitializeApp(group GroupName, stream StreamName) *App {
	wire.Build(
		NewApp,
		di.Set,
	)

	return &App{}
}
