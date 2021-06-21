package collect

import "github.com/prometheus/client_golang/prometheus"

type XenPromCollector interface {
	Name() string
	DefaultEnabled() bool
	PromCollector() prometheus.Collector
}
