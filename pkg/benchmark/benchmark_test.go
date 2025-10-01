// Package benchmark - For benchmarking operations
package benchmark

import (
	"testing"
	"time"
)

// TestBenchmark - a simple benchmark test.
func TestBenchmark(t *testing.T) {
	// Test the end to end cycle of benchmarking
	benchmark := New()
	benchmark.Start()
	time.Sleep(time.Second)
	benchmark.Stop()
	t.Log(benchmark.Report())

	// Ensure that the duration is greater than one second.
	if benchmark.Duration() < time.Second {
		t.Errorf("Expected duration to be greater than 1 second, but got %s", benchmark.Duration())
	}
}
