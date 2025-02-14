package postgres

import (
	"context"
	"database/sql"
	"exampleapp/domain/eventdispatcher"
	_ "github.com/lib/pq"
	"os"
)

type Conn struct {
	db              *sql.DB
	eventDispatcher eventdispatcher.EventDispatcher
}

func NewConn(eventDispatcher eventdispatcher.EventDispatcher) *Conn {
	dsn, _ := os.LookupEnv("GOOSE_DBSTRING")
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		panic(err)
	}

	return &Conn{db, eventDispatcher}
}

func (c *Conn) Close() error {
	return c.db.Close()
}

func (c *Conn) DB() *sql.DB {
	return c.db
}

func (c *Conn) Query(query string, args ...any) (*sql.Rows, error) {
	return c.db.QueryContext(context.Background(), query, args...)
}

func (c *Conn) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	c.eventDispatcher.Dispatch(ctx, &QueryEvent{query})
	return c.db.QueryContext(ctx, query, args...)
}

func (c *Conn) QueryRow(query string, args ...any) *sql.Row {
	return c.db.QueryRowContext(context.Background(), query, args...)
}

func (c *Conn) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	c.eventDispatcher.Dispatch(ctx, &QueryEvent{query})
	return c.db.QueryRowContext(ctx, query, args...)
}

type QueryEvent struct {
	Query string
}
