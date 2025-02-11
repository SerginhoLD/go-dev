package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

type HeadersController struct {
	db *sql.DB
}

func NewHeadersController(db *sql.DB) *HeadersController {
	return &HeadersController{db}
}

func (c *HeadersController) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}

	var version string
	err := c.db.QueryRow("select name from products where id = $1 and id = $1", 1).Scan(&version)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		version = "not found"
	case err != nil:
		panic(err)
	}

	fmt.Fprintf(w, "db test: %v\n", version)
}
