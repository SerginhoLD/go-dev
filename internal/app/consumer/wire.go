//go:build wireinject
// +build wireinject

package consumer

import (
	"exampleapp/internal/infrastructure/di"

	"github.com/google/wire"
)

func Initialize(group GroupName, stream StreamName) *consumer {
	wire.Build(
		newConsumer,
		di.Set,
	)

	return &consumer{}
}
