#!/usr/bin/env bash
# 30-run protocol described in Section 4.3.
set -euo pipefail
cd "$(dirname "$0")/.."
mkdir -p benchmark/results
OUT=benchmark/results/raw.csv
echo "run,op,volume,successful,failed,duration_s,tps,mean_latency_ms" > "$OUT"

RUNS=${RUNS:-30}
VOLUMES_INVOKE="15 50 100 200 400 1000 5000 10000 15000"
VOLUMES_QUERY="15 50 100 200 400 1000 5000 10000 15000"

for r in $(seq 1 "$RUNS"); do
  echo "[run $r/$RUNS] restarting network"
  ./network/down.sh >/dev/null 2>&1 || true
  ./network/generate.sh >/dev/null
  docker compose -f docker-compose.yaml up -d >/dev/null
  sleep 8
  ./network/create-channel.sh >/dev/null
  ./network/deploy-chaincode.sh accesscontrol >/dev/null
  ./network/deploy-chaincode.sh datasharing >/dev/null
  ./network/deploy-chaincode.sh crypto >/dev/null

  # warm-up 120 s ramp
  go run ./benchmark -mode warmup -op invoke -n 200

  for v in $VOLUMES_INVOKE; do
    go run ./benchmark -mode run -op invoke -n "$v" -run "$r" -out "$OUT"
  done
  for v in $VOLUMES_QUERY; do
    go run ./benchmark -mode run -op query -n "$v" -run "$r" -out "$OUT"
  done
done

python3 benchmark/stats.py --input "$OUT" --output benchmark/results/table6.csv