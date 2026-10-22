#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
export FABRIC_CFG_PATH=$PWD/config

ORDERER_CA=$PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

echo "[channel] creating healthcare channel"
peer channel create \
  -o orderer.example.com:7050 -c healthcare \
  -f channel-artifacts/healthcare.tx \
  --outputBlock channel-artifacts/healthcare.block \
  --tls --cafile "$ORDERER_CA"

for ORG in hospital lab; do
  case $ORG in
    hospital) PORT=7051; MSP=HospitalMSP; DOMAIN=hospital.example.com ;;
    lab)      PORT=9051; MSP=LabMSP;      DOMAIN=lab.example.com ;;
  esac
  export CORE_PEER_LOCALMSPID=$MSP
  export CORE_PEER_ADDRESS=peer0.$ORG:$PORT
  export CORE_PEER_TLS_ROOTCERT_FILE=$PWD/organizations/peerOrganizations/$DOMAIN/peers/peer0.$ORG.$DOMAIN/tls/ca.crt
  export CORE_PEER_MSPCONFIGPATH=$PWD/organizations/peerOrganizations/$DOMAIN/users/Admin@$DOMAIN/msp
  peer channel join -b channel-artifacts/healthcare.block
done