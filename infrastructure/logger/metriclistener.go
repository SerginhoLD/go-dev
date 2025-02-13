package logger

import (
	"exampleapp/domain/event"
	"exampleapp/io"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
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

	l.addHistogram(prometheus.HistogramOpts{
		Name:    "app_http_request_duration_ms",
		Buckets: []float64{1, 5, 10, 25, 50, 100, 250, 500, 1000},
	}, "pattern", "status")

	return l
}

func (l *MetricListener) addCounter(opts prometheus.CounterOpts, labels ...string) {
	if (len(labels)) > 0 {
		l.metrics[opts.Name] = prometheus.NewCounterVec(opts, labels)
	} else {
		l.metrics[opts.Name] = prometheus.NewCounter(opts)
	}

	prometheus.MustRegister(l.metrics[opts.Name])
}

func (l *MetricListener) addHistogram(opts prometheus.HistogramOpts, labels ...string) {
	if (len(labels)) > 0 {
		l.metrics[opts.Name] = prometheus.NewHistogramVec(opts, labels)
	} else {
		l.metrics[opts.Name] = prometheus.NewHistogram(opts)
	}

	prometheus.MustRegister(l.metrics[opts.Name])
}

func (l *MetricListener) OnHttpResponse(event *io.ResponseEvent) {
	l.metrics["app_http_request_duration_ms"].(*prometheus.HistogramVec).WithLabelValues(event.Request.Pattern, strconv.Itoa(event.StatusCode)).Observe(float64(event.Duration.Milliseconds()))
}

func (l *MetricListener) OnTestEvent(event *event.TestEvent) {
	l.metrics["app_test_counter"].(prometheus.Counter).Inc()
}
