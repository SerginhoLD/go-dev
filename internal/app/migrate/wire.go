//go:build wireinject
// +build wireinject

package migrate

import (
	//"exampleapp/internal/infrastructure/di"
	"exampleapp/internal/infrastructure/manticore"

	"github.com/google/wire"
)

func InitializeConn() *manticore.Conn {
	wire.Build(
		manticore.NewConn,
		//di.Set,
	)

	return &manticore.Conn{}
}
