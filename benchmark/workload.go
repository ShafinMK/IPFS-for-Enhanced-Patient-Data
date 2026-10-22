// Package benchmark implements the load driver referenced in Section 4.3.
package benchmark

import (
	"math/rand"
	"time"
)

type Op string

const (
	OpQuery  Op = "query"
	OpInvoke Op = "invoke"
)

type Task struct {
	Op      Op
	Payload []byte
	Issued  time.Time
}

// Generate produces n tasks for a given operation type.
func Generate(op Op, n int) []Task {
	out := make([]Task, n)
	for i := range out {
		payload := make([]byte, 256)
		rand.Read(payload)
		out[i] = Task{Op: op, Payload: payload, Issued: time.Now()}
	}
	return out
}