//go:build wireinject
// +build wireinject

package internal

import (
	"exampleapp/internal/infrastructure/di"

	"github.com/google/wire"
)

func InitializeScheduler() *Scheduler {
	wire.Build(
		NewApp,
		NewStartParseJob,
		di.Set,
	)

	return &Scheduler{}
}
