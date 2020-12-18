package prometheus

import "github.com/prometheus/client_golang/prometheus"

var AllMetrics = Metric{}

func RegisterPrometheusMetrics() {
	AllMetrics = Metric{
		HttpRequestsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Count of all HTTP requests",
		}, []string{"code", "method"}),
	}
	prometheus.MustRegister(AllMetrics.HttpRequestsTotal)
}
