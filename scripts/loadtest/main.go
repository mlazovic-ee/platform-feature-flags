package main

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	flagd "github.com/open-feature/go-sdk-contrib/providers/flagd/pkg"
	"github.com/open-feature/go-sdk/openfeature"
)

var (
	client        *openfeature.Client
	requestsTotal uint64
	errorsTotal   uint64
)

func main() {
	provider, err := flagd.NewProvider(
		flagd.WithHost("flags.easygo-drop-kick.com"),
		flagd.WithPort(443),
		flagd.WithTLS(""),
		flagd.WithoutCache(),
	)
	if err != nil {
		panic("Failed to create Flagd provider: " + err.Error())
	}
	openfeature.SetProviderAndWait(provider)
	client = openfeature.NewClient("loadtest")

	// === CONFIG ===
	concurrency := 200    // Number of goroutines
	flagKey := "web-test" // Flag key to evaluate
	ctx := context.Background()
	stopAfter := 30 * time.Second // Duration to run the test

	// === RUN ===
	fmt.Printf("Starting load test with %d goroutines for %v...\n", concurrency, stopAfter)
	start := time.Now()

	// Start workers
	for i := 0; i < concurrency; i++ {
		go func() {
			for {
				_, err := client.BooleanValue(ctx, flagKey, false, openfeature.EvaluationContext{})
				atomic.AddUint64(&requestsTotal, 1)
				if err != nil {
					atomic.AddUint64(&errorsTotal, 1)
				}
			}
		}()
	}

	// Print stats every second
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	lastRequests := uint64(0)
	lastErrors := uint64(0)

	done := time.After(stopAfter)

	for {
		select {
		case <-ticker.C:
			curRequests := atomic.LoadUint64(&requestsTotal)
			curErrors := atomic.LoadUint64(&errorsTotal)
			rate := curRequests - lastRequests
			errRate := curErrors - lastErrors
			fmt.Printf("Requests/sec: %6d | Errors/sec: %6d | Total: %8d\n", rate, errRate, curRequests)
			lastRequests = curRequests
			lastErrors = curErrors

		case <-done:
			totalTime := time.Since(start).Seconds()
			fmt.Println("Load test finished.")
			fmt.Printf("Total requests: %d in %.2f seconds (%.2f req/sec)\n",
				requestsTotal, totalTime, float64(requestsTotal)/totalTime)
			fmt.Printf("Total errors:   %d\n", errorsTotal)
			return
		}
	}
}
