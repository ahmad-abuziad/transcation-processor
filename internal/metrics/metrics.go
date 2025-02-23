package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	RequestDuration       *prometheus.HistogramVec
	TransactionsProcessed prometheus.Counter
	CacheHits             prometheus.Counter
	CacheMisses           prometheus.Counter
}

func NewMetrics() Metrics {
	m := Metrics{
		RequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "api_request_duration_seconds",
				Help:    "Histogram of API request durations.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"path"},
		),

		// TODO number of processed transactions per second
		TransactionsProcessed: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "transactions_processed_total",
				Help: "Total number of transactions processed.",
			},
		),

		CacheHits: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "cache_hits_total",
				Help: "Total number of cache hits.",
			},
		),

		CacheMisses: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "cache_misses_total",
				Help: "Total number of cache misses.",
			},
		),
	}
	prometheus.MustRegister(m.RequestDuration, m.TransactionsProcessed, m.CacheHits, m.CacheMisses)
	return m
}
