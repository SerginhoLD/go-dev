package app

import (
	"context"
	"database/sql"
	"errors"
	"exampleapp/internal/infrastructure/postgres"
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"log/slog"
	"os"
	_ "time/tzdata"
)

type Scheduler struct {
	scheduler gocron.Scheduler
	logger    *slog.Logger
	conn      *postgres.Conn
	testJob   *TestJob
}

func New(
	logger *slog.Logger,
	conn *postgres.Conn,
	testJob *TestJob,
) *Scheduler {
	s, _ := gocron.NewScheduler()

	return &Scheduler{
		s,
		logger,
		conn,
		testJob,
	}
}

func (app *Scheduler) Run() {
	app.logger.Debug(fmt.Sprintf("scheduler: start (env=%s)", os.Getenv("APP_ENV")))

	app.cron("*testJob", "TZ=Europe/Moscow * * * * *", app.transactionMiddleware(app.testJob.Handler))

	app.scheduler.Start()

	// block until you are ready to shut down
	select {
	//case <-time.After(time.Minute):
	}

	//if err = s.Shutdown(); err != nil {
	//	app.logger.Error(fmt.Sprintf("scheduler: %s", err.Error()))
	//}
}

func (app *Scheduler) cron(taskName string, crontab string, next func(context.Context)) {
	_, err := app.scheduler.NewJob(gocron.CronJob(crontab, false), gocron.NewTask(func() {
		requestId, _ := uuid.NewV7()
		ctx := context.WithValue(context.Background(), "X-Request-ID", requestId.String())
		ctx = context.WithValue(ctx, "SchedulerTaskName", taskName)
		app.logger.DebugContext(ctx, fmt.Sprintf(`scheduler: run "%s"`, taskName))

		defer func() {
			app.logger.DebugContext(ctx, fmt.Sprintf(`scheduler: done "%s"`, taskName))
		}()

		next(ctx)
	}))

	if err != nil {
		panic(err)
	}
}

func (app *Scheduler) transactionMiddleware(next func(context.Context)) func(context.Context) {
	return func(ctx context.Context) {
		tx, err := app.conn.DB().BeginTx(ctx, nil)

		if err != nil {
			app.logger.ErrorContext(ctx, fmt.Sprintf("scheduler: %s", err.Error()))
			panic(err)
		}

		app.logger.DebugContext(ctx, fmt.Sprintf(`sql: begin "%s"`, ctx.Value("SchedulerTaskName")))
		ctx = context.WithValue(ctx, "*sql.Tx", tx)

		defer func() {
			errRollback := tx.Rollback()

			if errRollback == nil {
				app.logger.DebugContext(ctx, fmt.Sprintf(`sql: rollback "%s"`, ctx.Value("SchedulerTaskName")))
			} else if !errors.Is(errRollback, sql.ErrTxDone) {
				panic(errRollback)
			}
		}()

		defer func() {
			if r := recover(); r != nil {
				panic(r) // вызываются все предыдущее defer
			}

			errCommit := tx.Commit()

			if errCommit != nil {
				panic(errCommit)
			}

			app.logger.DebugContext(ctx, fmt.Sprintf(`sql: commit "%s"`, ctx.Value("SchedulerTaskName")))
		}()

		next(ctx)
	}
}
