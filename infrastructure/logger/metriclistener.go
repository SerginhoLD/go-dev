package logger

import (
	"exampleapp/domain/event"
	"github.com/prometheus/client_golang/prometheus"
)

type MetricListener struct {
	metrics map[string]prometheus.Collector
}

func NewMetricListener() *MetricListener {
	metrics := make(map[string]prometheus.Collector)

	metrics["app_test_counter"] = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "app_test_counter",
		//Help: "An example counter metric",
	})

	prometheus.MustRegister(metrics["app_test_counter"])

	return &MetricListener{metrics}
}

func (l *MetricListener) OnTestEvent(event *event.TestEvent) {
	l.metrics["app_test_counter"].(prometheus.Counter).Inc()
}
