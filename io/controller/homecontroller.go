package controller

import (
	"example.com/m/domain/eventdispatcher"
	"fmt"
	"log/slog"
	"net/http"
)

type HomeController struct {
	logger          *slog.Logger
	eventDispatcher eventdispatcher.EventDispatcher
}

func NewHomeController(
	logger *slog.Logger,
	eventDispatcher eventdispatcher.EventDispatcher,
) *HomeController {
	return &HomeController{
		logger:          logger,
		eventDispatcher: eventDispatcher,
	}
}

func (c *HomeController) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello 6\n")
	c.logger.Info("hello again 6", "key", "val", "age", 25)
	c.eventDispatcher.Dispatch("Event1", map[string]any{"arg0": "val0"})
}
