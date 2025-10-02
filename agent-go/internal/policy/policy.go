package policy

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	m "github.com/example/l0-btc-regtest-ai-agent/agent-go/internal/metrics"
	o "github.com/example/l0-btc-regtest-ai-agent/agent-go/internal/overlay"
)

type Policy struct {
	mode string
	reg  *m.Registry
	ema  float64
}

func New(mode string, reg *m.Registry) *Policy {
	return &Policy{mode: mode, reg: reg}
}

func (p *Policy) Suggest(now time.Time, s o.Snapshot) o.Decision {
	// Простая контекстная политика: уменьшаем batch при росте RTT
	scale := math.Max(0.25, 1.5 - s.RTTms/200.0)
	batch := int(scale*128.0)
	flush := time.Duration(35.0*(2.0 - scale)) * time.Millisecond

	weights := make(map[string]float64, s.Peers)
	for i := 0; i < s.Peers; i++ {
		weights[fmt.Sprintf("peer-%d", i+1)] = 1.0 + 0.1*rand.Float64()
	}

	return o.Decision{BatchMax: batch, FlushEvery: flush, PeerWeight: weights}
}

func (p *Policy) Report(f o.Feedback) {
	p.ema = 0.9*p.ema + 0.1*f.LatencyMs
}
