package logger

import (
	"exampleapp/internal/domain/errors"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

var metrics = make(map[string]prometheus.Collector)

func init() {
	NewHistogram(prometheus.HistogramOpts{
		Name:    "app_http_request_duration_ms",
		Buckets: []float64{1, 5, 10, 25, 50, 100, 250, 500, 1000},
	}, "pattern", "status")

	NewHistogram(prometheus.HistogramOpts{
		Name:    "app_search_request_duration_ms",
		Buckets: []float64{1, 2, 4, 8, 16, 32, 64, 128, 256, 512},
	})
}

func NewCounter(opts prometheus.CounterOpts, labels ...string) {
	if _, has := metrics[opts.Name]; has {
		panic(errors.New(fmt.Sprintf(`metric "%s" already exists`, opts.Name)))
	}

	if len(labels) > 0 {
		metrics[opts.Name] = prometheus.NewCounterVec(opts, labels)
	} else {
		metrics[opts.Name] = prometheus.NewCounter(opts)
	}

	prometheus.MustRegister(metrics[opts.Name])
}

func NewHistogram(opts prometheus.HistogramOpts, labels ...string) {
	if _, has := metrics[opts.Name]; has {
		panic(errors.New(fmt.Sprintf(`metric "%s" already exists`, opts.Name)))
	}

	if len(labels) > 0 {
		metrics[opts.Name] = prometheus.NewHistogramVec(opts, labels)
	} else {
		metrics[opts.Name] = prometheus.NewHistogram(opts)
	}

	prometheus.MustRegister(metrics[opts.Name])
}

func AddMetric(name string, value float64, labelValues ...string) {
	defer func() {
		if rec := recover(); rec != nil {
			// логирование не должно приводить к падению приложения
		}
	}()

	metric, ok := metrics[name]

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
	}
}
