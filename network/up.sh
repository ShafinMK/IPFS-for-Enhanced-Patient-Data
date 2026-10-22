#!/usr/bin/env bash
set -euo pipefail
echo "[up] waiting for peers..."
sleep 8
docker compose -f docker-compose.yaml ps