// Package extension Common utilities for the collatz extension.
package extension

import (
	"fmt"
	"runtime"
	"strconv"
)

const (
	// CPUMultiplier defines the multiplier for the number of CPU cores to determine the batch size.
	CPUMultiplier = 100
	// DefaultBatchSize defines the default batch size if not set.
	DefaultBatchSize = 100
)

// CPUBatchSize returns the batch size based on the number of CPU cores and a predefined multiplier.
func CPUBatchSize() string {
	return strconv.Itoa(runtime.NumCPU() * CPUMultiplier)
}

// WaitErrHandling handles errors from goroutines, wrapping them with additional context if necessary.
func WaitErrHandling(err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("routine failed: %w", err)
}
