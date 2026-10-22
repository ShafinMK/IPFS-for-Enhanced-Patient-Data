package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Envelope struct {
	RecordID   string            `json:"recordID"`
	Ciphertext string            `json:"ciphertext"`
	IV         string            `json:"iv"`
	Tag        string            `json:"tag"`
	Wrapped    map[string]string `json:"wrapped"` // recipientID -> base64(RSA-OAEP(DK))
	Epoch      int               `json:"epoch"`
	Signature  string            `json:"signature"`
}

type CryptoContract struct{ contractapi.Contract }

func envKey(id string, epoch int) string { return fmt.Sprintf("env~%s~%d", id, epoch) }

func (c *CryptoContract) StoreEnvelope(ctx contractapi.TransactionContextInterface, e *Envelope) error {
	if e.Epoch == 0 { e.Epoch = 1 }
	bytes, _ := json.Marshal(e)
	return ctx.GetStub().PutState(envKey(e.RecordID, e.Epoch), bytes)
}

// RevokeRecipient removes a recipient's wrapped key and bumps the key epoch.
// The caller must supply the re-encrypted envelope for the new epoch.
func (c *CryptoContract) RevokeRecipient(ctx contractapi.TransactionContextInterface,
	recordID, recipient string, reencrypted *Envelope) error {

	// remove recipient from old epoch's wrapped map
	old, err := ctx.GetStub().GetState(envKey(recordID, reencrypted.Epoch-1))
	if err == nil && old != nil {
		var e Envelope
		_ = json.Unmarshal(old, &e)
		delete(e.Wrapped, recipient)
		b, _ := json.Marshal(&e)
		_ = ctx.GetStub().PutState(envKey(recordID, reencrypted.Epoch-1), b)
	}

	// store new epoch without the revoked recipient
	delete(reencrypted.Wrapped, recipient)
	bytes, _ := json.Marshal(reencrypted)
	return ctx.GetStub().PutState(envKey(recordID, reencrypted.Epoch), bytes)
}