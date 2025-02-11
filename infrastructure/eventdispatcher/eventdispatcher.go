package eventdispatcher

import (
	domainevent "exampleapp/domain/event"
	"exampleapp/domain/eventdispatcher"
	"exampleapp/infrastructure/logger"
	"fmt"
)

type EventDispatcherImpl struct {
	logListener    *logger.LogListener
	metricListener *logger.MetricListener
}

func New(logListener *logger.LogListener, metricListener *logger.MetricListener) *EventDispatcherImpl {
	d := &EventDispatcherImpl{logListener, metricListener}
	return d
}

func (d *EventDispatcherImpl) Dispatch(event interface{}) {
	switch e := event.(type) {
	case *domainevent.TestEvent:
		d.dispatchCallbacks(e, []func(interface{}){
			func(interface{}) { d.logListener.OnEvent1(e) },
			func(interface{}) { d.metricListener.OnEvent2(e) },
		})
	default:
		panic(fmt.Sprintf("unhandled type %T", event))
	}
}

func (d *EventDispatcherImpl) dispatchCallbacks(event interface{}, callbacks []func(interface{})) {
	for _, callback := range callbacks {
		callback(event)

		if stoppableEvent, ok := event.(eventdispatcher.StoppableEvent); ok {
			if stoppableEvent.IsPropagationStopped() {
				return
			}
		}
	}
}
