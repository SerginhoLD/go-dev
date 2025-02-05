package controller

import (
	"fmt"
	"log/slog"
	"net/http"
)

type HomeController struct {
	logger *slog.Logger
}

func NewHomeController(logger *slog.Logger) *HomeController {
	return &HomeController{
		logger: logger,
	}
}

func (c *HomeController) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello 6\n")
	c.logger.Info("hello again 6", "key", "val", "age", 25)
}
