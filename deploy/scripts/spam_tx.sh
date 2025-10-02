#!/usr/bin/env bash
set -euo pipefail
RPCU=user RPCP=pass

ensure_wallet() {
  local s=$1
  if ! docker compose exec "$s" bitcoin-cli -regtest -rpcuser=$RPCU -rpcpassword=$RPCP getwalletinfo >/dev/null 2>&1; then
    docker compose exec "$s" bitcoin-cli -regtest -rpcuser=$RPCU -rpcpassword=$RPCP loadwallet default >/dev/null 2>&1 || true
    if ! docker compose exec "$s" bitcoin-cli -regtest -rpcuser=$RPCU -rpcpassword=$RPCP getwalletinfo >/dev/null 2>&1; then
      docker compose exec "$s" bitcoin-cli -regtest -rpcuser=$RPCU -rpcpassword=$RPCP -named createwallet wallet_name=default descriptors=false >/dev/null
    fi
  fi
}

ensure_wallet btc1
ensure_wallet btc2

ADDR=$(docker compose exec btc1 bitcoin-cli -regtest -rpcuser=$RPCU -rpcpassword=$RPCP getnewaddress)
docker compose exec btc1 bitcoin-cli -regtest -rpcuser=$RPCU -rpcpassword=$RPCP generatetoaddress 1 "$ADDR" >/dev/null || true

for i in $(seq 1 20); do
  TO=$(docker compose exec btc2 bitcoin-cli -regtest -rpcuser=$RPCU -rpcpassword=$RPCP getnewaddress)
  TXID=$(docker compose exec btc1 bitcoin-cli -regtest -rpcuser=$RPCU -rpcpassword=$RPCP sendtoaddress "$TO" 0.001)
  echo "sent $TXID -> $TO"
  sleep 0.2
done

docker compose exec btc1 bitcoin-cli -regtest -rpcuser=$RPCU -rpcpassword=$RPCP generatetoaddress 1 "$ADDR" >/dev/null || true
echo "spam done."
