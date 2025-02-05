package gookiteventdispatcher

import (
	"github.com/gookit/event"
	"log/slog"
)

type GookitEventDispatcher struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *GookitEventDispatcher {
	d := &GookitEventDispatcher{logger: logger}
	d.init()
	return d
}

func (d *GookitEventDispatcher) Dispatch(name string, data map[string]any) {
	event.MustFire(name, data)
}

func (d *GookitEventDispatcher) listen(name string, callback func(map[string]any), priority ...int) {
	event.On(name, event.ListenerFunc(func(e event.Event) error {
		callback(e.Data())
		return nil
	}), priority...)
}

// todo
func (d *GookitEventDispatcher) init() {
	d.listen("Event1", func(data map[string]any) {
		d.logger.Info("ff8 Event1", "arg0", data["arg0"])
	})
}
