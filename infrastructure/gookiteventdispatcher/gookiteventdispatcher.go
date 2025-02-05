package gookiteventdispatcher

import (
	"example.com/m/infrastructure/logger"
	"github.com/gookit/event"
)

type GookitEventDispatcher struct {
	logListener *logger.LogListener
}

func New(logListener *logger.LogListener) *GookitEventDispatcher {
	d := &GookitEventDispatcher{logListener}
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

func (d *GookitEventDispatcher) init() {
	d.listen("Event1", func(data map[string]any) {
		d.logListener.OnEvent1(data)
	})
}
