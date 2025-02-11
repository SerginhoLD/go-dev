package logger

import (
	"exampleapp/domain/event"
	"github.com/prometheus/client_golang/prometheus"
)

type MetricListener struct {
	metrics map[string]prometheus.Collector
}

func NewMetricListener() *MetricListener {
	l := &MetricListener{make(map[string]prometheus.Collector)}

	l.addCounter(prometheus.CounterOpts{
		Name: "app_test_counter",
		//Help: "An example counter metric",
	})

	return l
}

func (l *MetricListener) addCounter(opts prometheus.CounterOpts) {
	l.metrics[opts.Name] = prometheus.NewCounter(opts)
	prometheus.MustRegister(l.metrics[opts.Name])
}

func (l *MetricListener) addGauge(opts prometheus.GaugeOpts) {
	l.metrics[opts.Name] = prometheus.NewGauge(opts)
	prometheus.MustRegister(l.metrics[opts.Name])
}

func (l *MetricListener) addHistogram(opts prometheus.HistogramOpts) {
	l.metrics[opts.Name] = prometheus.NewHistogram(opts)
	prometheus.MustRegister(l.metrics[opts.Name])
}

func (l *MetricListener) OnTestEvent(event *event.TestEvent) {
	l.metrics["app_test_counter"].(prometheus.Counter).Inc()
}
