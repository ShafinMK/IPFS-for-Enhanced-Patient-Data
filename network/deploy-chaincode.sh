#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."
export FABRIC_CFG_PATH=$PWD/config

CC_NAME=${1:-accesscontrol}
CC_PATH=./chaincode/$CC_NAME
CC_VERSION=1.0
CC_SEQ=1
CHANNEL=healthcare

for ORG in hospital lab; do
  case $ORG in
    hospital) PORT=7051; MSP=HospitalMSP; DOMAIN=hospital.example.com ;;
    lab)      PORT=9051; MSP=LabMSP;      DOMAIN=lab.example.com ;;
  esac
  export CORE_PEER_LOCALMSPID=$MSP
  export CORE_PEER_ADDRESS=peer0.$ORG:$PORT
  export CORE_PEER_TLS_ROOTCERT_FILE=$PWD/organizations/peerOrganizations/$DOMAIN/peers/peer0.$ORG.$DOMAIN/tls/ca.crt
  export CORE_PEER_MSPCONFIGPATH=$PWD/organizations/peerOrganizations/$DOMAIN/users/Admin@$DOMAIN/msp

  peer lifecycle chaincode package ${CC_NAME}.tar.gz \
    --path $CC_PATH --lang golang --label ${CC_NAME}_${CC_VERSION}
  peer lifecycle chaincode install ${CC_NAME}.tar.gz
  peer lifecycle chaincode approveformyorg -o orderer.example.com:7050 \
    --channelID $CHANNEL --name $CC_NAME --version $CC_VERSION --sequence $CC_SEQ \
    --package-id $(peer lifecycle chaincode calculatepackageid ${CC_NAME}.tar.gz) \
    --tls --cafile $PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
done

peer lifecycle chaincode commit -o orderer.example.com:7050 \
  --channelID $CHANNEL --name $CC_NAME --version $CC_VERSION --sequence $CC_SEQ \
  --tls --cafile $PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
  --peerAddresses peer0.hospital:7051 --tlsRootCertFiles $PWD/organizations/peerOrganizations/hospital.example.com/peers/peer0.hospital.example.com/tls/ca.crt \
  --peerAddresses peer0.lab:9051 --tlsRootCertFiles $PWD/organizations/peerOrganizations/lab.example.com/peers/peer0.lab.example.com/tls/ca.crt