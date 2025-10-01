// package benchmark - For benchmarking operations
package benchmark

import (
	"fmt"
	"time"
)

// Benchmark - for benchmarking operations
type Benchmark struct {
	StartTime time.Time
	EndTime   time.Time
}

// New allocates and returns a pointer to a Benchmark object, with the start and end times set to the current time.
func New() *Benchmark {
	return &Benchmark{
		StartTime: time.Now(),
		EndTime:   time.Now(),
	}
}

// Start records the current time as the start time for the benchmark.
func (b *Benchmark) Start() {
	b.StartTime = time.Now()
}

// Stop records the current time as the end time for the benchmark.
func (b *Benchmark) Stop() {
	b.EndTime = time.Now()
}

// Duration calculates and returns the duration between the start time and end time of the benchmark.
func (b *Benchmark) Duration() time.Duration {
	return b.EndTime.Sub(b.StartTime)
}

// Report returns a formatted string detailing the duration of the operation,
// including the start and end times.
func (b *Benchmark) Report() string {
	return fmt.Sprintf("The duration of the operation was %s with the operation starting at %s and ending at %s", b.Duration(), b.StartTime, b.EndTime)
}
