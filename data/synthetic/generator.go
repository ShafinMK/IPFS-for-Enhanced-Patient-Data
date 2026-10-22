// Generates synthetic (non-PHI) healthcare records for the evaluation in Section 4.1.
package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type Record struct {
	RecordID    string `json:"recordID"`
	PatientID   string `json:"patientID"`
	Issuer      string `json:"issuer"`
	RecordType  string `json:"recordType"` // prescription | lab | imaging
	Payload     string `json:"payload"`    // synthetic text, NOT real PHI
	PatientPub  string `json:"patientPub"`
}

func randHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func main() {
	n := flag.Int("records", 1000, "number of synthetic records")
	out := flag.String("out", "records.json", "output file")
	flag.Parse()

	types := []string{"prescription", "lab", "imaging"}
	issuers := []string{"hospital", "lab", "pharmacy"}
	var recs []Record
	for i := 0; i < *n; i++ {
		recs = append(recs, Record{
			RecordID:   fmt.Sprintf("rec-%08d", i),
			PatientID:  "patient-" + randHex(4),
			Issuer:     issuers[i%len(issuers)],
			RecordType: types[i%len(types)],
			Payload:    "synthetic-record-" + randHex(8),
			PatientPub: "pub-" + randHex(16),
		})
	}
	f, _ := os.Create(*out)
	defer f.Close()
	_ = json.NewEncoder(f).Encode(recs)
	fmt.Printf("wrote %d synthetic records to %s\n", len(recs), *out)
}