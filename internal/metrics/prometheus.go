package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rate_limiter_requests_total",
			Help: "Total number of requests processed by the rate limiter.",
		},
		[]string{"limiter"},
	)
	RequestsDenied = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rate_limiter_requests_denied_total",
			Help: "Total number of requests denied by the rate limiter.",
		},
		[]string{"limiter"},
	)
	AverageLatency = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rate_limiter_average_latency_seconds",
			Help: "Average latency of rate limiter checks.",
		},
		[]string{"limiter"},
	)
)

func Register() {
	prometheus.MustRegister(RequestsTotal, RequestsDenied, AverageLatency)
}

func UpdateMetrics(limiterName string, total, denied int64, avgLatency time.Duration) {
	RequestsTotal.WithLabelValues(limiterName).Add(float64(total))
	RequestsDenied.WithLabelValues(limiterName).Add(float64(denied))
	AverageLatency.WithLabelValues(limiterName).Set(avgLatency.Seconds())
}
