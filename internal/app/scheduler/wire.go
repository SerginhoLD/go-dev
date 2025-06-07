//go:build wireinject
// +build wireinject

package scheduler

import (
	"exampleapp/internal/app/scheduler/internal"
	"exampleapp/internal/infrastructure/di"

	"github.com/google/wire"
)

func Initialize() *scheduler {
	wire.Build(
		new,
		internal.NewStartParseJob,
		di.Set,
	)

	return &scheduler{}
}
