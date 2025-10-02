#!/usr/bin/env bash
set -euo pipefail
# Bootstrap regtest: create wallets, connect peers, generate funds
BTCCLI="docker compose exec -T btc1 bitcoin-cli -regtest -rpcuser=user -rpcpassword=pass"
cd "$(dirname "$0")/../docker"

# Create wallets on btc1..btc5
for n in 1 2 3 4 5; do
  docker compose exec -T btc$n bitcoin-cli -regtest -rpcuser=user -rpcpassword=pass -named createwallet wallet_name=default descriptors=false >/dev/null || true
done

# Connect nodes
for n in 2 3 4 5; do
  ADDR=$(docker compose exec -T btc$n bitcoin-cli -regtest -rpcuser=user -rpcpassword=pass getnetworkinfo | jq -r '.localaddresses[0].address' || echo "btc$n")
  docker compose exec -T btc1 bitcoin-cli -regtest -rpcuser=user -rpcpassword=pass addnode btc$n onetry || true
done

# Generate initial coins to btc1
ADDR1=$($BTCCLI getnewaddress)
$BTCCLI generatetoaddress 101 "$ADDR1"

echo "Bootstrap done."
