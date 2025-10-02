package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/example/l0-btc-regtest-ai-agent/agent-go/internal/metrics"
	"github.com/example/l0-btc-regtest-ai-agent/agent-go/internal/overlay"
)

func getenv(k, def string) string { if v:=os.Getenv(k); v!="" {return v}; return def }

func main() {
	_ = godotenv.Load()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	mode := getenv("EXPERIMENT_GROUP", "control")
	addr := getenv("METRICS_ADDR", ":9090")
	zmq  := getenv("BTC_ZMQ_RAWTX", "tcp://btc1:28332")

	reg := metrics.New(mode)
	if err := overlay.NewCollector(mode, zmq, reg).Start(); err != nil {
		log.Fatal().Err(err).Msg("collector start failed")
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", reg.MetricsHandler())
	mux.HandleFunc("/firstseen", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type","application/json")
		_ = json.NewEncoder(w).Encode(reg.FirstSeenSnapshot())
	})

	log.Info().Str("mode",mode).Str("addr",addr).Str("zmq",zmq).Msg("agent up")
	log.Fatal().Err(http.ListenAndServe(addr, mux)).Msg("http stopped")
}
