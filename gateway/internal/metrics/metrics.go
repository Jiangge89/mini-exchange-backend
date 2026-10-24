package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

/*
同一个指标名不能注册两次，否则会 panic
*/

var (
	RequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "gateway",
		Name:      "requests_total",
		Help:      "Total JSON-RPC requests by method.",
	}, []string{"method"})
	
	RequestsErrorTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "gateway",
		Name:      "requests_error_total",
		Help:      "Total JSON-RPC error requests by method and error code.",
	}, []string{"method", "error"})

	RequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "gateway",
		Name:      "request_duration_seconds",
		Help:      "JSON-RPC request duration in seconds.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"method"})

	UpstreamErrors = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "gateway",
		Name:      "upstream_errors_total",
		Help:      "Total upstream request errors.",
	})

	CacheHits = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "gateway",
		Name:      "cache_hits_total",
		Help:      "Total cache hits.",
	})

	CacheMisses = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "gateway",
		Name:      "cache_misses_total",
		Help:      "Total cache misses.",
	})

	CacheErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "gateway",
		Name:      "cache_errors_total",
		Help:      "Total cache errors.",
	}, []string{"method"})

	UpstreamHealth = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "gateway",
		Name:      "upstream_health_score",
		Help:      "Health score of each upstream (0.0 to 1.0).",
	}, []string{"upstream"})
)

func Register() {
	// promauto handles registration automatically.
}
