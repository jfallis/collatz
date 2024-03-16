package extension

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jfallis/collatz/pkg/collatz"
	"golang.org/x/sync/errgroup"
)

const (
	DefaultBruteforceRoutineBatchLimit = 1000
)

func Bruteforce(num, maxBatchCount uint64, logging bool) error {
	if num < maxBatchCount {
		maxBatchCount = num
	}

	errorGroup, ctx := errgroup.WithContext(context.Background())
	sem := make(chan struct{}, maxBatchCount)

	for index := uint64(0); index < num; index++ {
		isBreakPoint := (index+1)%maxBatchCount == 0

		func(number uint64, breakPoint bool) {
			errorGroup.Go(func() error {
				sem <- struct{}{}
				defer func() { <-sem }()

				col := collatz.New(number)
				if err := col.Calculate(); err != nil {
					return fmt.Errorf("bruteforce failed: %w", err)
				}

				if logging && breakPoint {
					slog.Info(fmt.Sprintf("bruteforce number: %d", number))
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
		}(index+1, isBreakPoint)
	}

	return waitErrHandling(errorGroup.Wait())
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
