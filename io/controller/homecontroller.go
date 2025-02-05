package controller

import (
	"example.com/m/domain/eventdispatcher"
	"fmt"
	"net/http"
)

type HomeController struct {
	eventDispatcher eventdispatcher.EventDispatcher
}

func NewHomeController(
	eventDispatcher eventdispatcher.EventDispatcher,
) *HomeController {
	return &HomeController{
		eventDispatcher: eventDispatcher,
	}
}

func (c *HomeController) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello 6\n")
	c.eventDispatcher.Dispatch("Event1", map[string]any{"arg0": "val0"}) // todo: struct
}
