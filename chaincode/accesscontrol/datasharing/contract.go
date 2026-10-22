// Package main implements the consent-gated Data Sharing smart contract (Section 3.3, Alg. 4).
package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Grant struct {
	RecordID  string `json:"recordID"`
	Grantee   string `json:"grantee"`
	AccessType string `json:"accessType"` // read | read-write
	GrantedAt string `json:"grantedAt"`
	RevokedAt string `json:"revokedAt,omitempty"`
	TokenHash string `json:"tokenHash"`
}

type DataSharingContract struct{ contractapi.Contract }

func grantKey(recordID, grantee string) string {
	return "grant~" + recordID + "~" + grantee
}

func (c *DataSharingContract) GrantAccess(ctx contractapi.TransactionContextInterface, g *Grant) error {
	g.GrantedAt = time.Now().UTC().Format(time.RFC3339)
	bytes, _ := json.Marshal(g)
	return ctx.GetStub().PutState(grantKey(g.RecordID, g.Grantee), bytes)
}

func (c *DataSharingContract) RevokeAccess(ctx contractapi.TransactionContextInterface, recordID, grantee string) error {
	bytes, err := ctx.GetStub().GetState(grantKey(recordID, grantee))
	if err != nil { return err }
	if bytes == nil { return fmt.Errorf("no grant for %s/%s", recordID, grantee) }
	var g Grant
	_ = json.Unmarshal(bytes, &g)
	g.RevokedAt = time.Now().UTC().Format(time.RFC3339)
	out, _ := json.Marshal(g)
	return ctx.GetStub().PutState(grantKey(recordID, grantee), out)
}

func (c *DataSharingContract) CheckAccess(ctx contractapi.TransactionContextInterface, recordID, grantee, accessType string) (bool, error) {
	bytes, err := ctx.GetStub().GetState(grantKey(recordID, grantee))
	if err != nil || bytes == nil { return false, err }
	var g Grant
	_ = json.Unmarshal(bytes, &g)
	if g.RevokedAt != "" { return false, nil }
	if g.AccessType == "read-write" { return true, nil }
	return g.AccessType == accessType, nil
}