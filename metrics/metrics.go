// Package metrics contains prometheus metrics registry.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

type (
	// Registry interface.
	Registry interface {
		prometheus.Registerer
		prometheus.Gatherer
	}

	// Duration is a HistogramVec metric.
	Duration interface {
		WithLabelValues(lvs ...string) prometheus.Observer
	}

	// Count is a CounterVec metric.
	Count interface {
		WithLabelValues(lvs ...string) prometheus.Counter
	}
)

// New method create new metrics registry with default collectors.
func New() Registry {

	var r = prometheus.NewRegistry()
	r.MustRegister(collectors.NewBuildInfoCollector())
	r.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	r.MustRegister(collectors.NewGoCollector())

	return r
}

// globalRegistry is auxilary registry for application.
var globalRegistry Registry = prometheus.NewRegistry()

// Global return registered global logger (or default).
func Global() Registry {
	return globalRegistry
}

// SetGlobal setup global logger.
func SetGlobal(registry Registry) {
	globalRegistry = registry
}
