package controller

import (
	"database/sql"
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

	rows, err := c.db.Query("select version()")

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	rows.Next()

	var version string
	rows.Scan(&version)

	fmt.Fprintf(w, "db test: %v\n", version)
}
