package server

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	buckets = []float64{300, 1200, 5000}
)
var (
	httpCall = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        "guts_theater_http_request_total",
			Help:        "How many HTTP requests processed, partitioned by status code, method and HTTP path (with patterns).",
			ConstLabels: prometheus.Labels{"service": "guts"},
		},
		[]string{"code", "method", "path"},
	)
	httpDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "guts_theater_http_duration_seconds",
		Help:        "Duration of HTTP requests.",
		ConstLabels: prometheus.Labels{"service": "guts"},
		Buckets:     buckets,
	}, []string{"code", "method", "path"})
)

func registerPromVec() {
	prometheus.MustRegister(httpCall)
	prometheus.MustRegister(httpDuration)
}
