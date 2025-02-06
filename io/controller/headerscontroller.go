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
}
