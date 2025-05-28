package manticore

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Conn struct {
	db *sql.DB
}

func NewConn() *Conn {
	db, err := sql.Open("mysql", os.Getenv("MANTICORE_TCP"))

	if err != nil {
		panic(err)
	}

	return &Conn{db}
}

func (c *Conn) Close() error {
	return c.db.Close()
}

func (c *Conn) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	slog.DebugContext(ctx, fmt.Sprintf("sql: %s", query))
	return c.db.ExecContext(ctx, query, args...)
}
