package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
)

type Conn struct {
	logger *slog.Logger
	db     []*sql.DB
}

func NewConn(logger *slog.Logger) *Conn {
	master, err := sql.Open("postgres", os.Getenv("GOOSE_DBSTRING"))

	if err != nil {
		panic(err)
	}

	slave1, err := sql.Open("postgres", os.Getenv("GOOSE_DBSTRING_SLAVE1"))

	if err != nil {
		master.Close()
		panic(err)
	}

	return &Conn{logger, []*sql.DB{master, slave1}}
}

//func (c *Conn) Close() error {
//	return c.db.Close()
//}

func (c *Conn) Master() *sql.DB {
	return c.db[0]
}

func (c *Conn) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	c.logger.DebugContext(ctx, fmt.Sprintf("sql: %s", query))

	if tx, ok := ctx.Value("*sql.Tx").(*sql.Tx); ok {
		return tx.QueryContext(ctx, query, args...)
	}

	return c.db[1].QueryContext(ctx, query, args...)
}

func (c *Conn) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	c.logger.DebugContext(ctx, fmt.Sprintf("sql: %s", query))

	if tx, ok := ctx.Value("*sql.Tx").(*sql.Tx); ok {
		return tx.QueryRowContext(ctx, query, args...)
	}

	return c.db[1].QueryRowContext(ctx, query, args...)
}

func (c *Conn) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	c.logger.DebugContext(ctx, fmt.Sprintf("sql: %s", query))

	if tx, ok := ctx.Value("*sql.Tx").(*sql.Tx); ok {
		return tx.ExecContext(ctx, query, args...)
	}

	return c.db[0].ExecContext(ctx, query, args...)
}
