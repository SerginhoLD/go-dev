package eventdispatcher

// todo: domain | infrastructure

import (
	"github.com/gookit/event"
	"log/slog"
)

type EventDispatcher struct {
	logger *slog.Logger
}

func NewEventDispatcher(logger *slog.Logger) *EventDispatcher {
	d := &EventDispatcher{logger: logger}
	d.init()
	return d
}

func (d *EventDispatcher) Dispatch(name string, data map[string]any) {
	event.MustFire(name, data)
}

// todo: domain | infrastructure
func (d *EventDispatcher) listen(name string, callback func(map[string]any), priority int) {
	event.On(name, event.ListenerFunc(func(e event.Event) error {
		callback(e.Data())
		return nil
	}), priority)
}

func (d *EventDispatcher) init() {
	d.listen("Event1", func(data map[string]any) {
		d.logger.Info("ff5 Event1", "arg0", data["arg0"])
	}, 0)
}
