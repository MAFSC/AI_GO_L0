package main

import (
	"time"
	"github.com/prometheus/client_golang/prometheus"
)

var agentAlive = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "l0_agent_metrics_ts",
	Help: "last metrics export timestamp (unix seconds)",
})

func init() {
	prometheus.MustRegister(agentAlive)
	go func() {
		t := time.NewTicker(2 * time.Second)
		for range t.C {
			agentAlive.Set(float64(time.Now().Unix()))
		}
	}()
}
