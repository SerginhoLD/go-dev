package postgres

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
)

type Conn struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewConn(logger *slog.Logger) *Conn {
	db, err := sql.Open("postgres", os.Getenv("GOOSE_DBSTRING"))

	if err != nil {
		panic(err)
	}

	return &Conn{db, logger}
}

//func (c *Conn) Close() error {
//	return c.db.Close()
//}

func (c *Conn) DB() *sql.DB {
	return c.db
}

func (c *Conn) Query(query string, args ...any) (*sql.Rows, error) {
	return c.db.QueryContext(context.Background(), query, args...)
}

func (c *Conn) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	c.logger.DebugContext(ctx, fmt.Sprintf("sql: %s", query))

	if tx, ok := ctx.Value("*sql.Tx").(*sql.Tx); ok {
		return tx.QueryContext(ctx, query, args...)
	}

	return c.db.QueryContext(ctx, query, args...)
}

func (c *Conn) QueryRow(query string, args ...any) *sql.Row {
	return c.db.QueryRowContext(context.Background(), query, args...)
}

func (c *Conn) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	c.logger.DebugContext(ctx, fmt.Sprintf("sql: %s", query))

	if tx, ok := ctx.Value("*sql.Tx").(*sql.Tx); ok {
		return tx.QueryRowContext(ctx, query, args...)
	}

	return c.db.QueryRowContext(ctx, query, args...)
}

func (c *Conn) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	c.logger.DebugContext(ctx, fmt.Sprintf("sql: %s", query))

	if tx, ok := ctx.Value("*sql.Tx").(*sql.Tx); ok {
		return tx.ExecContext(ctx, query, args...)
	}

	return c.db.ExecContext(ctx, query, args...)
}
