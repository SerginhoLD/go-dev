package internal

import (
	"context"
	"encoding/json"
	"exampleapp/internal/domain/messages"
	"exampleapp/internal/domain/usecase"
	"exampleapp/internal/infrastructure/di"
	"exampleapp/internal/infrastructure/keydb"
	"fmt"
	"log/slog"
	"os"

	"github.com/google/uuid"
)

type GroupName string  // for wire
type StreamName string // for wire

type App struct {
	group    string
	stream   string
	consumer *keydb.Stream
	useCase  *usecase.ParseObjectsUseCase
}

func NewApp(
	group GroupName,
	stream StreamName,
	consumer *keydb.Stream,
	useCase *usecase.ParseObjectsUseCase,
) *App {
	return &App{
		group:    string(group),
		stream:   string(stream),
		consumer: consumer,
		useCase:  useCase,
	}
}

func (app *App) Run() {
	slog.Debug(fmt.Sprintf("consumer: start (env=%s, ver=%s)", os.Getenv("APP_ENV"), di.Version))

	err := app.consumer.Consume(context.Background(), keydb.ConsumeOpts{Group: app.group, Stream: app.stream, DelAfterAck: true}, requestMiddleware(func(ctx context.Context, msg *keydb.Message) error {
		switch msg.Headers()["type"] {
		case "*messages.ParsePageMessage":
			command := deserialize[messages.ParsePageMessage](msg.Payload())
			err := app.useCase.Handle(ctx, command)

			if err != nil {
				return err
			}
		default:
			slog.Warn(fmt.Sprintf("consumer: unsupported message type: %s", msg.Headers()["type"]), "stream", app.stream, "id", msg.Id(), "payload", msg.Payload(), "headers", msg.Headers())
		}

		return msg.Ack(ctx)
	}))

	if err != nil {
		slog.Error(fmt.Sprintf("consumer: %s", err.Error()))
	}
}

func deserialize[T any](payload string) *T {
	command := new(T)
	err := json.Unmarshal([]byte(payload), command)

	if err != nil {
		panic(err)
	}

	return command
}

func requestMiddleware(next func(context.Context, *keydb.Message) error) func(context.Context, *keydb.Message) error {
	return func(ctx context.Context, msg *keydb.Message) (err error) {
		defer func() {
			if rec := recover(); rec != nil {
				err = fmt.Errorf("%v", rec)
			}
		}()

		requestId, _ := uuid.NewV7()
		return next(context.WithValue(ctx, "X-Request-ID", requestId.String()), msg)
	}
}
