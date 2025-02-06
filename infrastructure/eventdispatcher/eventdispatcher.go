package eventdispatcher

import (
	domainevent "example.com/m/domain/event"
	"example.com/m/infrastructure/logger"
	"fmt"
)

type EventDispatcherImpl struct {
	logListener *logger.LogListener
}

func New(logListener *logger.LogListener) *EventDispatcherImpl {
	d := &EventDispatcherImpl{logListener}
	return d
}

func (d *EventDispatcherImpl) Dispatch(event interface{}) {
	switch e := event.(type) {
	case *domainevent.TestEvent:
		d.logListener.OnEvent1(e)
	default:
		panic(fmt.Sprintf("unhandled type %T", event))
	}
}
