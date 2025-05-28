package internal

import (
	"context"
	"exampleapp/internal/infrastructure/di"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	_ "time/tzdata"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type Scheduler struct {
	scheduler     gocron.Scheduler
	startParseJob *StartParseJob
}

func NewApp(
	startParseJob *StartParseJob,
) *Scheduler {
	s, _ := gocron.NewScheduler()

	return &Scheduler{
		s,
		startParseJob,
	}
}

func (app *Scheduler) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	slog.DebugContext(ctx, fmt.Sprintf("scheduler: start (env=%s, ver=%s)", os.Getenv("APP_ENV"), di.Version))

	app.cron(ctx, "TZ=Europe/Moscow * * * * *", "*StartParseJob", app.startParseJob.Handle)

	app.scheduler.Start()

	<-ctx.Done()
	slog.DebugContext(ctx, "scheduler: stop")

	if err := app.scheduler.Shutdown(); err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("scheduler: %s", err.Error()))
	}
}

func (app *Scheduler) cron(ctx context.Context, crontab string, taskName string, next func(context.Context)) {
	// https://github.com/go-co-op/gocron/blob/e1b7d52/example_test.go#L617
	_, err := app.scheduler.NewJob(gocron.CronJob(crontab, false), gocron.NewTask(app.requestMiddleware(taskName, next)), gocron.WithContext(ctx))

	if err != nil {
		panic(err)
	}
}

func (app *Scheduler) requestMiddleware(taskName string, next func(context.Context)) func(context.Context) {
	return func(ctx context.Context) {
		requestId, _ := uuid.NewV7()
		ctx = context.WithValue(ctx, "X-Request-ID", requestId.String())
		slog.DebugContext(ctx, fmt.Sprintf(`scheduler: run "%s"`, taskName))

		defer func() {
			slog.DebugContext(ctx, fmt.Sprintf(`scheduler: done "%s"`, taskName))
		}()

		next(ctx)
	}
}
