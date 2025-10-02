#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/../docker"
# Send a bunch of tx from btc1 to btc2..btc5 and mine some blocks
for i in $(seq 1 20); do
  for n in 2 3 4 5; do
    ADDR=$(docker compose exec -T btc$n bitcoin-cli -regtest -rpcuser=user -rpcpassword=pass getnewaddress)
    docker compose exec -T btc1 bitcoin-cli -regtest -rpcuser=user -rpcpassword=pass sendtoaddress "$ADDR" 0.001 >/dev/null
  done
  docker compose exec -T btc1 bitcoin-cli -regtest -rpcuser=user -rpcpassword=pass generatetoaddress 1 $(docker compose exec -T btc1 bitcoin-cli -regtest -rpcuser=user -rpcpassword=pass getnewaddress) >/dev/null
done
echo "Spam complete."
