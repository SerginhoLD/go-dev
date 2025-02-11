package eventdispatcher

import (
	domainevent "exampleapp/domain/event"
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
		d.logListener.OnEvent1(e)
		d.metricListener.OnEvent2(e)
	default:
		panic(fmt.Sprintf("unhandled type %T", event))
	}
}
