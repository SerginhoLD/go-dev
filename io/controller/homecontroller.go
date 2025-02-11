package controller

import (
	"exampleapp/domain/event"
	"exampleapp/domain/eventdispatcher"
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
	fmt.Fprintf(w, "hello 7\n")
	c.eventDispatcher.Dispatch(&event.TestEvent{Value: "stopped"})
}
