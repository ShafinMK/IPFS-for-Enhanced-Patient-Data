#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."

echo "[generate] crypto material"
cryptogen generate --config=config/crypto-config.yaml --output=organizations

echo "[generate] genesis block"
mkdir -p channel-artifacts
FABRIC_CFG_PATH=config configtxgen -profile Genesis -outputBlock channel-artifacts/genesis.block -channelID system-channel

echo "[generate] channel tx"
FABRIC_CFG_PATH=config configtxgen -profile HealthcareChannel -outputCreateChannelTx channel-artifacts/healthcare.tx -channelID healthcare