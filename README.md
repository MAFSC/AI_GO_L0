# L0 AI Go Agent — Bitcoin Regtest Only

This repository is a focused framework for the L0 AI accelerator for the Bitcoin testnet (regtest).
Includes:
- Go agent (overlay + policy + Prometheus metrics)
- Bitcoin adapter (RPC + ZMQ stubs)
- Docker training ground: 5 Bitcoin nodes (regtest) + 2 agent instances (A/B), Prometheus + Grafana
- Scripts: regtest setup, block and transaction generation, network latency
- ONNX inference stub (can be replaced with your own model)

## Quick Start
```bash
cd deploy/docker
docker compose up --build -d
# initialization of regtest, wallets, connections and block/tx generation
../scripts/bootstrap_regtest.sh
../scripts/spam_tx.sh
# Grafana: http://localhost:3000 (Armaturenbrett L0 BTC Regtest)
# Prometheus: http://localhost:9090
```