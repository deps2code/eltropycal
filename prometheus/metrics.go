package prometheus

import "github.com/prometheus/client_golang/prometheus"

//Can add more custom metrics in here
type Metric struct {
	HttpRequestsTotal *prometheus.CounterVec
}
