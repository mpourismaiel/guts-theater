package prometheus

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	buckets = []float64{300, 1200, 5000}
)
var (
	HttpCall = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        "guts_theater_http_request_total",
			Help:        "How many HTTP requests processed, partitioned by status code, method and HTTP path (with patterns).",
			ConstLabels: prometheus.Labels{"service": "guts"},
		},
		[]string{"code", "method", "path"},
	)
	HttpDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "guts_theater_http_duration_seconds",
		Help:        "Duration of HTTP requests.",
		ConstLabels: prometheus.Labels{"service": "guts"},
		Buckets:     buckets,
	}, []string{"code", "method", "path"})
	DbCall = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:        "guts_theater_db_call_total",
		Help:        "The total number of calls to database",
		ConstLabels: prometheus.Labels{"service": "guts"},
	}, []string{"model", "action"})

	once sync.Once
)

func registerPromVec() {
	prometheus.MustRegister(HttpCall)
	prometheus.MustRegister(HttpDuration)
}

func Setup() {
	once.Do(func() {
		registerPromVec()
	})
}
