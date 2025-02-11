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
	return &EventDispatcherImpl{logListener, metricListener}
}

func (d *EventDispatcherImpl) Dispatch(event interface{}) error {
	var callbacks []func(interface{}) error

	switch e := event.(type) {
	case *domainevent.TestEvent:
		callbacks = append(callbacks, func(interface{}) error { return d.logListener.OnEvent1(e) })
		callbacks = append(callbacks, func(interface{}) error { d.metricListener.OnEvent2(e); return nil })
	default:
		d.logListener.OnUnhandledEvent(e)
		panic(fmt.Sprintf("Unhandled event \"%T\"", event))
	}

	for _, callback := range callbacks {
		err := callback(event)

		if err != nil {
			return err
		}

		if stoppableEvent, ok := event.(eventdispatcher.StoppableEvent); ok {
			if stoppableEvent.IsPropagationStopped() {
				return nil
			}
		}
	}

	return nil
}
