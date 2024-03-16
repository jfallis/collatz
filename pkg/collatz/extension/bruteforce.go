package extension

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/jfallis/collatz/pkg/collatz"
	"golang.org/x/sync/errgroup"
)

const (
	DefaultBatchLimit = 1_000
)

var (
	DefaultBruteforceRoutineBatchLimit = big.NewInt(DefaultBatchLimit)

	zero = big.NewInt(0)
	one  = big.NewInt(1)
)

func Bruteforce(start, end, maxBatchCount *big.Int, logging bool) (*big.Int, error) {
	if end.Cmp(maxBatchCount) < 0 {
		maxBatchCount.Set(end)
	}

	errorGroup, ctx := errgroup.WithContext(context.Background())
	sem := make(chan struct{}, maxBatchCount.Int64())

	iStart := new(big.Int).Set(start)
	lastIndex := new(big.Int)
	for index := iStart; index.Cmp(end) < 0; index.Add(index, one) {
		isBreakPoint := new(big.Int).Mod(new(big.Int).Add(index, one), maxBatchCount).Cmp(zero) == 0

		func(number *big.Int, breakPoint bool) {
			errorGroup.Go(func() error {
				sem <- struct{}{}
				defer func() { <-sem }()

				col := collatz.New(new(big.Int).Set(number))
				if err := col.Calculate(); err != nil {
					return fmt.Errorf("bruteforce failed: %w", err)
				}

				if logging && breakPoint {
					slog.Info(fmt.Sprintf("bruteforce number: %s, steps: %d", col.Number().String(), len(col.Steps())))
				}

				if success := col.Success(); success {
					return collatz.SuccessError{Number: col.Number(), Steps: col.Steps()}
				}

				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					return nil
				}
			})
		}(new(big.Int).Add(index, one), isBreakPoint)

		lastIndex.Set(index)
	}

	return lastIndex.Add(lastIndex, big.NewInt(1)), waitErrHandling(errorGroup.Wait())
}

func waitErrHandling(err error) error {
	if err == nil {
		return nil
	}

	var success collatz.SuccessError
	if errors.As(err, &success) {
		return fmt.Errorf("successfully found the number %w", err)
	}

	return fmt.Errorf("bruteforce failed: %w", err)
}
