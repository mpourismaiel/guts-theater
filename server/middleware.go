package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/mpourismaiel/guts-theater/prometheus"
)

func patternHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		prometheus.HttpCall.WithLabelValues(http.StatusText(ww.Status()), r.Method, r.URL.RawPath).Inc()
		prometheus.HttpDuration.WithLabelValues(http.StatusText(ww.Status()), r.Method, r.URL.RawPath).Observe(float64(time.Since(start).Nanoseconds()) / 1000000)
	}
	return http.HandlerFunc(fn)
}
