#!/bin/sh
# Apply network delays to host interfaces (demo), affects docker bridge as well
set -e
if ! command -v tc >/dev/null; then
  echo "tc not found, installing..."; apk add --no-cache iproute2 jq
fi
for IF in $(ls /sys/class/net | grep -E 'eth|ens|br-'); do
  tc qdisc replace dev $IF root netem delay 80ms 20ms loss 0.2%
done
echo "Applied netem on host/bridge interfaces."
sleep infinity
