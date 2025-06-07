//go:build wireinject
// +build wireinject

package server

import (
	"exampleapp/internal/app/server/internal"
	"exampleapp/internal/infrastructure/di"

	"github.com/google/wire"
)

func Initialize() *server {
	wire.Build(
		new,
		di.Set,
		internal.NewCoverageController,
		internal.NewHomeController,
	)

	return &server{}
}
