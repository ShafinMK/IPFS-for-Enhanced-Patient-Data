package benchmark

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// FabricClient is the minimal interface the driver needs; wire it to your
// fabric-gateway client in main().
type FabricClient interface {
	Query(ctx context.Context, args ...string) ([]byte, error)
	Invoke(ctx context.Context, args ...string) ([]byte, error)
}

type Result struct {
	RunID      int           `json:"runID"`
	Op         Op            `json:"op"`
	Volume     int           `json:"volume"`
	Successful int           `json:"successful"`
	Failed     int           `json:"failed"`
	Duration   time.Duration `json:"duration"`
	TPS        float64       `json:"tps"`
	MeanLatency time.Duration `json:"meanLatency"`
}

// RunOnce executes one full run of `volume` tasks with the given concurrency.
func RunOnce(ctx context.Context, client FabricClient, op Op, volume, concurrency, runID int) (*Result, error) {
	tasks := Generate(op, volume)
	start := time.Now()
	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrency)
	var mu sync.Mutex
	success, failed := 0, 0
	var totalLatency time.Duration

	for _, t := range tasks {
		wg.Add(1)
		sem <- struct{}{}
		go func(t Task) {
			defer wg.Done()
			defer func() { <-sem }()
			issued := time.Now()
			var err error
			if t.Op == OpQuery {
				_, err = client.Query(ctx, "ReadRecord", string(t.Payload))
			} else {
				_, err = client.Invoke(ctx, "CreateRecord", string(t.Payload))
			}
			mu.Lock()
			defer mu.Unlock()
			totalLatency += time.Since(issued)
			if err != nil { failed++ } else { success++ }
		}(t)
	}
	wg.Wait()
	dur := time.Since(start)
	return &Result{
		RunID: runID, Op: op, Volume: volume,
		Successful: success, Failed: failed,
		Duration: dur, TPS: float64(success) / dur.Seconds(),
		MeanLatency: totalLatency / time.Duration(max(success, 1)),
	}, nil
}

func max(a, b int) int { if a > b { return a }; return b }

// Warmup issues a short ramp before measurement (Section 4.3 protocol).
func Warmup(ctx context.Context, client FabricClient, op Op, n int) error {
	_, err := RunOnce(ctx, client, op, n, 1, 0)
	return err
}

// RESTEndpoint returns a stub error to make the placeholder explicit.
var _ = fmt.Sprintf