package overlay

import (
	"context"
	"time"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/example/l0-btc-regtest-ai-agent/agent-go/internal/metrics"
)

type Config struct {
	ListenAddr string
	MaxPeers   int
	BatchMax   int
	FlushEvery time.Duration
}

type Policy interface {
	Suggest(now time.Time, snapshot Snapshot) Decision
	Report(feedback Feedback)
}

type Snapshot struct {
	Peers int
	Queue int
	RTTms float64
	Loss  float64
}

type Decision struct {
	BatchMax   int
	FlushEvery time.Duration
	PeerWeight map[string]float64
}

type Feedback struct {
	LatencyMs float64
	Duplicates int
}

type Overlay struct {
	cfg Config
	metrics *metrics.Collector
}

func New(reg *metrics.Registry, cfg Config) *Overlay {
	return &Overlay{cfg: cfg, metrics: metrics.NewCollector(reg)}
}

func (o *Overlay) Run(ctx context.Context, p Policy) error {
	t := time.NewTicker(250 * time.Millisecond)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case now := <-t.C:
			// TODO: collect real snapshot from sockets/queues
			snap := Snapshot{Peers: 5, Queue: 200, RTTms: 80, Loss: 0.2}
			dec := p.Suggest(now, snap)
			_ = dec // apply to send queues

			// feedback demo
			p.Report(Feedback{LatencyMs: 120, Duplicates: 3})
			o.metrics.ObserveLatency(120)
			log.Debug().Str("apply", fmt.Sprintf("%+v", dec)).Msg("tick")
		}
	}
}
