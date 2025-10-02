#!/usr/bin/env bash
set -euo pipefail

PROM=${PROM:-http://localhost:9090}
Q_RULE='l0_first_seen_delta_ms'
Q_FALLBACK='(l0_first_seen_ts_seconds{mode="experiment"} - on (txid) l0_first_se                                                                                                             en_ts_seconds{mode="control"}) * 1000'

fetch_vals () {
  local q="$1"
  curl -fsS "$PROM/api/v1/query" --data-urlencode "query=$q" \
  | jq -r '.data.result[].value[1]' 2>/dev/null || true
}

vals="$(fetch_vals "$Q_RULE")"
[ -z "$vals" ] && vals="$(fetch_vals "$Q_FALLBACK")"

if [ -z "$vals" ]; then
  echo "no data from Prometheus ($PROM)"
  exit 0
../scripts/report_first_seen.sht_seen.sh} ms  mean=${mean} ms  min=${min} ms  ma
count=22  p50=0.000 ms  p95=0.000 ms  mean=0.000 ms  min=0 ms  max=0 ms
root@vmi2125901:~/1/AI_GO_L0/deploy/docker# cd ~/1/AI_GO_L0/deploy/docker

# override для агентов — укажем разные ZMQ и RPC
cat > docker-compose.override.yml <<'EOF'
services:
  agent_control:
    command: ["/agent","-mode","control",
              "-zmq","tcp://btc2:28332",
              "-rpc","http://user:pass@btc2:18443"]
  agent_experiment:
    command: ["/agent","-mode","experiment",
              "-zmq","tcp://btc3:28332",
              "-rpc","http://user:pass@btc3:18443"]
