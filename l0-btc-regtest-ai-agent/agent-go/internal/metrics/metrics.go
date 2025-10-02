package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Registry struct {
	Gatherer *prometheus.Registry
}

func NewRegistry() *Registry {
	reg := prometheus.NewRegistry()
	return &Registry{Gatherer: reg}
}

type Collector struct {
	lat prometheus.Summary
}

func NewCollector(r *Registry) *Collector {
	lat := prometheus.NewSummary(prometheus.SummaryOpts{
		Name: "l0_latency_ms",
		Help: "Latency observed by the agent",
		Objectives: map[float64]float64{0.5:0.01, 0.95:0.005, 0.99:0.001},
	})
	r.Gatherer.MustRegister(lat)
	return &Collector{lat: lat}
}

func (c *Collector) ObserveLatency(ms float64) { c.lat.Observe(ms) }
