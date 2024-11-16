package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Counter describes a metric that accumulates values monotonically.
// An example of a counter is the number of received HTTP requests.
type Counter interface {
	With(labelValues ...string) Counter
	Add(delta float64)
}

// Gauge describes a metric that takes specific values over time.
// An example of a gauge is the current depth of a job queue.
type Gauge interface {
	With(labelValues ...string) Gauge
	Set(value float64)
	Add(delta float64)
}

// Histogram describes a metric that takes repeated observations of the same
// kind of thing, and produces a statistical summary of those observations,
// typically expressed as quantiles or buckets. An example of a histogram is
// HTTP request latencies.
type Histogram interface {
	With(labelValues ...string) Histogram
	Observe(value float64)
}

var (
	// NewCounterVec works like the function of the same name in the prometheus package
	// but it automatically registers the CounterVec with the prometheus.DefaultRegisterer.

	// CounterVec is a Collector that bundles a set of Counters that all share the same Desc,
	// but have different values for their variable labels. This is used if you want to count
	// the same thing partitioned by various dimensions (e.g. number of HTTP requests, partitioned
	// by response code and method).
	requestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "app",
			Name:      "handler_request_total_counter",
			Help:      "Total amount of request by handler",
		},
		[]string{"handler"},
	)

	handlerHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "app",
			Name:      "handler_duration_histogram",
			Help:      "Total duration of handler processing by request",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"handler"},
	)
)

func IncRequestCounter(handler string) {
	requestCounter.WithLabelValues(handler).Inc()
}

func StoreHandlerDuration(handler string, since time.Duration) {
	handlerHistogram.WithLabelValues(handler).Observe(since.Seconds())
}

type Metrics struct {
	RequestCount   Counter
	RequestLatency Histogram
	ErrorCount     Counter
}

func Setup() *Metrics {
	return &Metrics{
		RequestCount: NewCounterFrom(prometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "string_service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),

		RequestLatency: NewHistogramFrom(prometheus.HistogramOpts{
			Namespace: "api",
			Subsystem: "string_service",
			Name:      "request_latency_seconds",
			Help:      "Request duration in seconds.",
		}, []string{"method"}),

		ErrorCount: NewCounterFrom(prometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "string_service",
			Name:      "error_count",
			Help:      "Number of errors occurred.",
		}, []string{"method"}),
	}
}
