package logger

import (
	"exampleapp/domain/event"
	"exampleapp/io/controller"
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
	if len(labels) > 0 {
		l.metrics[opts.Name] = prometheus.NewCounterVec(opts, labels)
	} else {
		l.metrics[opts.Name] = prometheus.NewCounter(opts)
	}

	prometheus.MustRegister(l.metrics[opts.Name])
}

func (l *MetricListener) addHistogram(opts prometheus.HistogramOpts, labels ...string) {
	if len(labels) > 0 {
		l.metrics[opts.Name] = prometheus.NewHistogramVec(opts, labels)
	} else {
		l.metrics[opts.Name] = prometheus.NewHistogram(opts)
	}

	prometheus.MustRegister(l.metrics[opts.Name])
}

func (l *MetricListener) Add(name string, value float64, labelValues ...string) {
	//if value < 0 {
	//	return
	//}

	metric, ok := l.metrics[name]

	if !ok {
		return
	}

	switch m := metric.(type) {
	case *prometheus.CounterVec:
		m.WithLabelValues(labelValues...).Add(value)
	case prometheus.Counter:
		m.Add(value)
	case *prometheus.HistogramVec:
		m.WithLabelValues(labelValues...).Observe(value)
	case prometheus.Histogram:
		m.Observe(value)
		//default:
		//	panic(fmt.Sprintf("Unknown metric type (name: %s, type: %T)", name, metric))
	}

}

func (l *MetricListener) OnHttpResponse(event *controller.ResponseEvent) {
	l.Add("app_http_request_duration_ms", float64(event.Duration.Milliseconds()), event.Request.Pattern, strconv.Itoa(event.StatusCode))
}

func (l *MetricListener) OnTestEvent(event *event.TestEvent) {
	l.Add("app_test_counter", 1)
}
