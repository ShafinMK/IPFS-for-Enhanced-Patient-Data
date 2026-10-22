#!/usr/bin/env bash
set -euo pipefail
docker compose -f docker-compose.yaml down -v
rm -rf organizations channel-artifacts