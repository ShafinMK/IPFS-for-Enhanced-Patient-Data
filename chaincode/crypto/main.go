// main.go
package main

import (
	"log"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	cc, err := contractapi.NewChaincode(&CryptoContract{})
	if err != nil { log.Fatalf("create chaincode: %v", err) }
	if err := cc.Start(); err != nil { log.Fatalf("start chaincode: %v", err) }
}