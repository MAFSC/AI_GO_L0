package policy

import (
	"math"
	"math/rand"
	"time"
	"fmt"
	"context"

	"github.com/example/l0-btc-regtest-ai-agent/agent-go/internal/overlay"
	"github.com/example/l0-btc-regtest-ai-agent/agent-go/internal/metrics"
	"github.com/example/l0-btc-regtest-ai-agent/agent-go/internal/policy/onnx"
)

type Policy struct{
	mode string
	reg *metrics.Registry
	ema float64
	model *onnx.Model // may be nil if not loaded
}

func New(mode string, reg *metrics.Registry) *Policy {
	// try to load ONNX model if present (optional)
	m, _ := onnx.Load("ai-model.onnx")
	return &Policy{mode: mode, reg: reg, ema: 0, model: m}
}

func (p *Policy) Suggest(now time.Time, s overlay.Snapshot) overlay.Decision {
	// If ONNX model is loaded, use it
	if p.model != nil {
		b, flush, err := p.model.Predict(context.Background(), []float32{float32(s.RTTms), float32(s.Loss), float32(s.Peers), float32(s.Queue)})
		if err == nil {
			return overlay.Decision{
				BatchMax: int(b),
				FlushEvery: time.Duration(flush) * time.Millisecond,
				PeerWeight: map[string]float64{"peer-1":1.0},
			}
		}
	}
	// Fallback: simple contextual control
	scale := math.Max(0.25, 1.5 - s.RTTms/200.0)
	batch := int(scale * float64(128))
	flush := time.Duration(35 * float64(time.Millisecond) * (2.0 - scale))

	weights := map[string]float64{}
	for i:=0;i<s.Peers;i++ { 
		weights[fmt.Sprintf("peer-%d", i+1)] = 1.0 + 0.1*rand.Float64() 
	}

	return overlay.Decision{
		BatchMax: batch, FlushEvery: flush, PeerWeight: weights,
	}
}

func (p *Policy) Report(f overlay.Feedback) {
	// EMA for stability / could feed to online learner
	p.ema = 0.9*p.ema + 0.1*f.LatencyMs
}
