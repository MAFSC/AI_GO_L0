#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/../docker"

create_or_load() {
  local svc="$1"
  # если кошелёк уже открыт — просто сообщим
  if docker compose exec "$svc" bitcoin-cli -regtest -rpcuser=user -rpcpassword=pass getwalletinfo >/dev/null 2>&1; then
    echo "[$svc] wallet already loaded"
    return
  fi
  # пробуем создать; если уже есть на диске — загрузим
  docker compose exec "$svc" bitcoin-cli -regtest -rpcuser=user -rpcpassword=pass -named createwallet wallet_name=default descriptors=false >/dev/null 2>&1 \
  || docker compose exec "$svc" bitcoin-cli -regtest -rpcuser=user -rpcpassword=pass loadwallet default >/dev/null 2>&1 || true
  echo "[$svc] wallet ready"
}

for n in 1 2 3 4 5; do create_or_load "btc$n"; done
for n in 2 3 4 5; do docker compose exec btc1 bitcoin-cli -regtest -rpcuser=user -rpcpassword=pass addnode btc$n onetry >/dev/null 2>&1 || true; done

ADDR1=$(docker compose exec btc1 bitcoin-cli -regtest -rpcuser=user -rpcpassword=pass getnewaddress)
HEIGHT=$(docker compose exec btc1 bitcoin-cli -regtest -rpcuser=user -rpcpassword=pass getblockcount)
if [ "$HEIGHT" -lt 101 ]; then
  docker compose exec btc1 bitcoin-cli -regtest -rpcuser=user -rpcpassword=pass generatetoaddress $((101-HEIGHT)) "$ADDR1" >/dev/null
fi

NEW_H=$(docker compose exec btc1 bitcoin-cli -regtest -rpcuser=user -rpcpassword=pass getblockcount)
echo "Bootstrap done. Height: $NEW_H"
