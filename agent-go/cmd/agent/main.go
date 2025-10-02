package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/example/l0-btc-regtest-ai-agent/agent-go/internal/overlay"
	"github.com/example/l0-btc-regtest-ai-agent/agent-go/internal/policy"
	"github.com/example/l0-btc-regtest-ai-agent/agent-go/internal/metrics"
	"github.com/example/l0-btc-regtest-ai-agent/agent-go/internal/bitcoin"
)

func main() {
	_ = godotenv.Load()
	zerolog.TimeFieldFormat = time.RFC3339

	mode := getenv("EXPERIMENT_GROUP", "control")
	metricsAddr := getenv("METRICS_ADDR", ":9090")
	overlayListen := getenv("OVERLAY_LISTEN", ":7070")

	log.Info().Str("mode", mode).Msg("starting L0 BTC agent (regtest)")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// metrics
	reg := metrics.NewRegistry()
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.HandlerFor(reg.Gatherer, promhttp.HandlerOpts{}))
		log.Info().Str("addr", metricsAddr).Msg("metrics server")
		log.Fatal().Err(http.ListenAndServe(metricsAddr, mux)).Msg("metrics server failed")
	}()

	// overlay
	ov := overlay.New(reg, overlay.Config{
		ListenAddr:  overlayListen,
		MaxPeers:    64,
		BatchMax:    128,
		FlushEvery:  35 * time.Millisecond,
	})

	// policy (AI placeholder inside)
	pl := policy.New(mode, reg)

	// bitcoin adapter (RPC/ZMQ placeholders)
	btc := bitcoin.NewAdapter(bitcoin.Config{
		RPCHost:   getenv("BTC_RPC_HOST", "btc1"),
		RPCPort:   getenv("BTC_RPC_PORT", "18443"),
		RPCUser:   getenv("BTC_RPC_USER", "user"),
		RPCPass:   getenv("BTC_RPC_PASS", "pass"),
		ZMQRawTx:  getenv("BTC_ZMQ_RAWTX", "tcp://btc1:28332"),
		ZMQRawBlk: getenv("BTC_ZMQ_RAWBLOCK", "tcp://btc1:28333"),
	}, reg)

	// wire adapter to overlay (simplified: only metrics/feedback demo)
	go btc.Run(ctx, func fb(latMs float64, dups int){ 
		pl.Report(overlay.Feedback{LatencyMs: latMs, Duplicates: dups})
	})

	// run overlay
	if err := ov.Run(ctx, pl); err != nil {
		log.Fatal().Err(err).Msg("overlay stopped")
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" { return v }
	return def
}
