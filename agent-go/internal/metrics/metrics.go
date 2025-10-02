package metrics

import (
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Registry struct {
	reg         *prometheus.Registry
	FirstSeenTS *prometheus.GaugeVec
	FirstCount  *prometheus.CounterVec

	mu        sync.RWMutex
	firstSeen map[string]float64
	mode      string
}

func New(mode string) *Registry {
	r := prometheus.NewRegistry()
	firstTS := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "l0_first_seen_ts_seconds",
		Help: "Unix timestamp when txid first seen by this agent",
	}, []string{"txid", "mode"})
	firstCnt := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "l0_first_seen_count",
		Help: "Number of first-seen transactions observed",
	}, []string{"mode"})

	r.MustRegister(firstTS, firstCnt)
	return &Registry{
		reg:         r,
		FirstSeenTS: firstTS,
		FirstCount:  firstCnt,
		firstSeen:   make(map[string]float64),
		mode:        mode,
	}
}

func (m *Registry) ObserveFirstSeen(txid string, t time.Time) {
	labels := prometheus.Labels{"txid": txid, "mode": m.mode}
	m.FirstSeenTS.With(labels).Set(float64(t.Unix()))
	m.FirstCount.With(prometheus.Labels{"mode": m.mode}).Inc()
	m.mu.Lock()
	m.firstSeen[txid] = float64(t.Unix())
	m.mu.Unlock()
}

func (m *Registry) MetricsHandler() http.Handler { return promhttp.HandlerFor(m.reg, promhttp.HandlerOpts{}) }

func (m *Registry) FirstSeenSnapshot() map[string]float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make(map[string]float64, len(m.firstSeen))
	for k, v := range m.firstSeen {
		out[k] = v
	}
	return out
}
