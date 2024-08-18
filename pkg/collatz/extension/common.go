package extension

import (
	"errors"
	"fmt"
	"math/big"
	"runtime"
	"strconv"
	"sync"

	"github.com/jfallis/collatz/pkg/collatz"
)

const (
	CPUMultiplier       = 100
	DefaultBatchSizeInt = 1000
)

var (
	defaultBatchSize     *big.Int
	defaultBatchSizeOnce sync.Once
)

func DefaultBatchSize() *big.Int {
	defaultBatchSizeOnce.Do(func() {
		defaultBatchSize = big.NewInt(DefaultBatchSizeInt)
	})

	return defaultBatchSize
}

func CPUBatchSize() string {
	return strconv.Itoa(runtime.NumCPU() * CPUMultiplier)
}

func WaitErrHandling(err error) error {
	if err == nil {
		return nil
	}

	var success collatz.SuccessError
	if errors.As(err, &success) {
		return fmt.Errorf("successfully found the number %w", err)
	}

	return fmt.Errorf("routine failed: %w", err)
}
