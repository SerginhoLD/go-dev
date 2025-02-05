package controller

import (
	"fmt"
	"log/slog"
	"net/http"
)

type HomeController struct {
	logger slog.Logger
}

func NewHomeController(logger slog.Logger) *HomeController {
	return &HomeController{
		logger: logger,
	}
}

func (c *HomeController) Index(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello 3\n")
	c.logger.Info("hello again 3", "key", "val", "age", 25)
}
