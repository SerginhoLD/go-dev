package logger

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	metrics map[string]prometheus.Collector
}

func NewMetrics() *Metrics {
	l := &Metrics{make(map[string]prometheus.Collector)}

	l.addHistogram(prometheus.HistogramOpts{
		Name:    "app_http_request_duration_ms",
		Buckets: []float64{1, 5, 10, 25, 50, 100, 250, 500, 1000},
	}, "pattern", "status")

	return l
}

func (l *Metrics) addCounter(opts prometheus.CounterOpts, labels ...string) {
	if len(labels) > 0 {
		l.metrics[opts.Name] = prometheus.NewCounterVec(opts, labels)
	} else {
		l.metrics[opts.Name] = prometheus.NewCounter(opts)
	}

	prometheus.MustRegister(l.metrics[opts.Name])
}

func (l *Metrics) addHistogram(opts prometheus.HistogramOpts, labels ...string) {
	if len(labels) > 0 {
		l.metrics[opts.Name] = prometheus.NewHistogramVec(opts, labels)
	} else {
		l.metrics[opts.Name] = prometheus.NewHistogram(opts)
	}

	prometheus.MustRegister(l.metrics[opts.Name])
}

func (l *Metrics) Add(name string, value float64, labelValues ...string) {
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
