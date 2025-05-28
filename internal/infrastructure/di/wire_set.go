package di

import (
	"exampleapp/internal/domain/messenger"
	"exampleapp/internal/domain/usecase"
	"exampleapp/internal/infrastructure/keydb"
	"exampleapp/internal/infrastructure/manticore"
	repositoryimpl "exampleapp/internal/infrastructure/repository"

	"github.com/google/wire"
	"github.com/gorilla/schema"
)

var Version = "" // ldflags

var Set = wire.NewSet(
	//wire.NewSet(
	//	slog.New,
	//	logger.NewHandler,
	//	wire.Bind(new(slog.Handler), new(*logger.Handler)),
	//	wire.InterfaceValue(new(io.Writer), os.Stderr),
	//),

	wire.NewSet(
		keydb.NewStream,
		keydb.NewClient,
		//wire.Value(keydb.MaxLen(2000)),
	),
	wire.Bind(new(messenger.Bus), new(*keydb.Stream)),

	manticore.NewClient,
	repositoryimpl.NewHttpClient,

	schema.NewDecoder,

	repositoryimpl.NewObjectRepositoryImpl,
	repositoryimpl.NewVeObjectRepositoryImpl,
	repositoryimpl.NewPlanRepositoryImpl,

	usecase.NewParseObjectsUseCase,
	usecase.NewStartParseObjectsUseCase,
	usecase.NewPaginateObjectsUseCase,
)
