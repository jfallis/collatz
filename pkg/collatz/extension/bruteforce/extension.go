package bruteforce

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/jfallis/collatz/pkg/collatz/extension"

	"github.com/jfallis/collatz/pkg/collatz"
	"golang.org/x/sync/errgroup"
)

var (
	zero = big.NewInt(0)
	one  = big.NewInt(1)
)

type Request struct {
	Start      string
	End        string
	BatchSize  string
	PrintAll   bool
	Logging    bool
	EnableStep bool
}

type Results struct {
	Number        string
	StepLen       collatz.KeyValue
	LargestNumber collatz.KeyValue
}

func Run(ctx context.Context, request Request) (*Results, error) {
	req, err := bruteforceHandler(request)
	if err != nil {
		return nil, err
	}

	errGroup, ctx := errgroup.WithContext(ctx)
	sem := make(chan struct{}, req.batchSize.Int64())

	index := new(big.Int).Set(req.start)
	for x := index; x.Cmp(req.end) < 0; x.Add(x, one) {
		isBreakPoint := new(big.Int).Mod(new(big.Int).Add(x, one), req.breakPoint).Cmp(zero) == 0
		number := new(big.Int).Add(x, one)
		errGroup.Go(func() error {
			return routine(ctx, sem, number, isBreakPoint, request)
		})
	}

	if err := extension.WaitErrHandling(errGroup.Wait()); err != nil {
		return nil, err
	}

	return &Results{Number: index.String()}, nil
}

type RequestValues struct {
	start, end, batchSize, breakPoint *big.Int
}

func bruteforceHandler(request Request) (*RequestValues, error) {
	start, ok := new(big.Int).SetString(request.Start, collatz.Base)
	if !ok {
		return nil, fmt.Errorf("failed to set start value: %s", request.Start)
	}

	end, ok := new(big.Int).SetString(request.End, collatz.Base)
	if !ok {
		return nil, fmt.Errorf("failed to set end value: %s", request.End)
	}

	batchSize, ok := new(big.Int).SetString(request.BatchSize, collatz.Base)
	if !ok {
		return nil, fmt.Errorf("failed to set batch size value: %s", request.BatchSize)
	}

	difference := new(big.Int).Sub(end, start)
	if difference.Cmp(zero) <= 0 {
		return nil, fmt.Errorf("difference between start and end must be greater than 0")
	}

	if difference.Cmp(batchSize) < 0 {
		batchSize.Set(end)
	}

	breakPoint := new(big.Int).Div(difference, batchSize)
	if request.PrintAll || difference.Cmp(extension.DefaultBatchSize()) < 0 {
		breakPoint.Set(one)
	}

	return &RequestValues{
		start:      start,
		end:        end,
		batchSize:  batchSize,
		breakPoint: breakPoint,
	}, nil
}

func routine(ctx context.Context, sem chan struct{}, number *big.Int, breakPoint bool, request Request) error {
	sem <- struct{}{}
	defer func() { <-sem }()

	col := collatz.New(number.String())
	if err := col.Calculate(request.EnableStep && breakPoint); err != nil {
		return fmt.Errorf("routine error: %w", err)
	}

	if col.Success() {
		return collatz.SuccessError{String: col.String()}
	}

	if request.Logging && breakPoint {
		slog.InfoContext(ctx, col.String())
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
