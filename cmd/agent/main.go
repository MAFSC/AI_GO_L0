package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"ai-agent/pkg/aiassist"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// --- провайдер фич (пока из ENV; в проде — сбор из ваших метрик) ---
func sampleFeatsFromEnv() aiassist.Feats {
	return aiassist.Feats{
		AvgRttMs:    envF("AVG_RTT_MS", 50),
		MinRttMs:    envF("MIN_RTT_MS", 20),
		MaxRttMs:    envF("MAX_RTT_MS", 120),
		AvgLossPct:  envF("AVG_LOSS_PCT", 0.5),
		DupRatio:    envF("DUP_RATIO", 0.1),
		MempoolSize: envF("MEMPOOL_SIZE", 1500),
		SockQueue:   envF("SOCK_QUEUE", 3),
		OverlayLoad: envF("OVERLAY_LOAD", 0.4),
	}
}
func envF(k string, def float64) float64 {
	if v := os.Getenv(k); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return def
}

var (
	inferCalls = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "l0_infer_calls_total", Help: "inference calls"},
		[]string{"status"},
	)
	inferLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "l0_infer_latency_ms",
			Help:    "inference latency, ms",
			Buckets: []float64{5, 10, 20, 50, 100, 200, 400, 800, 1600},
		},
	)
	inferDelta = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "l0_infer_predicted_delta_ms",
		Help: "predicted delta ms",
	})
	agentBatch = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "l0_agent_batch_size",
		Help: "current batch size",
	})
	agentFlush = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "l0_agent_flush_ms",
		Help: "current flush interval ms",
	})
)

func mustRegister() {
	prometheus.MustRegister(inferCalls, inferLatency, inferDelta, agentBatch, agentFlush)
}

func main() {
	cfg := aiassist.FromEnv()

	pollMs := envI("POLL_MS", 500)
	bind := envS("BIND", ":9109")

	agent := aiassist.NewAgent(cfg)
	mustRegister()

	// HTTP: /metrics, /current
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/current", func(w http.ResponseWriter, r *http.Request) {
		sug, _ := agent.SuggestParams(sampleFeatsFromEnv())
		_ = json.NewEncoder(w).Encode(sug)
	})
	go func() {
		log.Printf("HTTP listening on %s", bind)
		if err := http.ListenAndServe(bind, mux); err != nil {
			log.Fatalf("http: %v", err)
		}
	}()

	// главный цикл
	ticker := time.NewTicker(time.Duration(pollMs) * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		start := time.Now()
		ctx := sampleFeatsFromEnv()

		sug, err := agent.SuggestParams(ctx)
		el := time.Since(start)

		if err != nil && !sug.FromAI {
			inferCalls.WithLabelValues("error").Inc()
		} else {
			inferCalls.WithLabelValues("ok").Inc()
		}
		inferLatency.Observe(float64(el.Milliseconds()))
		inferDelta.Set(sug.DeltaMs)
		agentBatch.Set(float64(sug.Batch))
		agentFlush.Set(float64(sug.FlushMs))

		// TODO: тут вызывай свой applyBatchAndFlush(sug.Batch, sug.FlushMs)
		fmt.Printf("policy → batch=%d flush_ms=%d (ai=%v, delta=%.1f)\n",
			sug.Batch, sug.FlushMs, sug.FromAI, sug.DeltaMs)
	}
}

func envI(k string, def int) int {
	if v := os.Getenv(k); v != "" {
		if n, err := strconv.Atoi(v); err == nil { return n }
	}
	return def
}
func envS(k, def string) string {
	if v := os.Getenv(k); v != "" { return v }
	return def
}
