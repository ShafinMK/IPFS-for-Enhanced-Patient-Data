// Package main implements the Access Control smart contract (Section 3.3).
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Record struct {
	RecordID   string `json:"recordID"`
	PatientID  string `json:"patientID"`
	Issuer     string `json:"issuer"`
	Ciphertext string `json:"ciphertext"` // base64 AES-256-GCM envelope
	IV         string `json:"iv"`
	Tag        string `json:"tag"`
	Metadata   string `json:"metadata"`
	Epoch      int    `json:"epoch"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
	IsDeleted  bool   `json:"isDeleted"`
}

type AccessControlContract struct{ contractapi.Contract }

func (c *AccessControlContract) CreateRecord(ctx contractapi.TransactionContextInterface, r *Record) error {
	exists, err := c.RecordExists(ctx, r.RecordID)
	if err != nil { return err }
	if exists { return fmt.Errorf("record %s already exists", r.RecordID) }
	r.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	r.UpdatedAt = r.CreatedAt
	r.Epoch = 1
	bytes, _ := json.Marshal(r)
	return ctx.GetStub().PutState(r.RecordID, bytes)
}

func (c *AccessControlContract) ReadRecord(ctx contractapi.TransactionContextInterface, id string) (*Record, error) {
	bytes, err := ctx.GetStub().GetState(id)
	if err != nil { return nil, err }
	if bytes == nil { return nil, fmt.Errorf("record %s not found", id) }
	var r Record
	if err := json.Unmarshal(bytes, &r); err != nil { return nil, err }
	if r.IsDeleted { return nil, fmt.Errorf("record %s revoked", id) }
	return &r, nil
}

func (c *AccessControlContract) UpdateRecord(ctx contractapi.TransactionContextInterface, r *Record) error {
	exists, err := c.RecordExists(ctx, r.RecordID)
	if err != nil { return err }
	if !exists { return fmt.Errorf("record %s not found", r.RecordID) }
	r.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	bytes, _ := json.Marshal(r)
	return ctx.GetStub().PutState(r.RecordID, bytes)
}

// RevokeRecord implements the soft-delete described in Section 3.3.
func (c *AccessControlContract) RevokeRecord(ctx contractapi.TransactionContextInterface, id string) error {
	r, err := c.ReadRecord(ctx, id)
	if err != nil { return err }
	r.IsDeleted = true
	r.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	bytes, _ := json.Marshal(r)
	return ctx.GetStub().PutState(id, bytes)
}

func (c *AccessControlContract) RecordExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	bytes, err := ctx.GetStub().GetState(id)
	if err != nil { return false, err }
	return bytes != nil, nil
}

// History returns the audit trail required by Section 3.3 (Historian registry).
func (c *AccessControlContract) History(ctx contractapi.TransactionContextInterface, id string) ([]map[string]interface{}, error) {
	it, err := ctx.GetStub().GetHistoryForKey(id)
	if err != nil { return nil, err }
	defer it.Close()
	var out []map[string]interface{}
	for it.HasNext() {
		m, err := it.Next()
		if err != nil { return nil, err }
		h := sha256.Sum256(m.Value)
		out = append(out, map[string]interface{}{
			"txId": m.TxId, "timestamp": m.Timestamp.String(),
			"isDelete": m.IsDelete, "valueHash": hex.EncodeToString(h[:]),
		})
	}
	return out, nil
}