package controller

import (
	"fmt"
	"log/slog"
	"net/http"
)

type HeadersController struct {
	logger *slog.Logger
}

func NewHeadersController(logger *slog.Logger) *HeadersController {
	return &HeadersController{logger}
}

func (c *HeadersController) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}
