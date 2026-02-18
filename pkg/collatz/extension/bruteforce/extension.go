// Package bruteforce provides functionality for executing Collatz Conjecture calculations
// over a range of numbers using concurrent routines and batch processing.
package bruteforce

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/jfallis/collatz/pkg/collatz"
	"github.com/jfallis/collatz/pkg/collatz/extension"
	"golang.org/x/sync/errgroup"
)

const (
	activeRoutines = 1000
)

type bigIntValues struct {
	zero *big.Int
	one  *big.Int
}

// Results is the output from the Run bruteforce function.
type Results struct {
	Num           string
	StepLen       collatz.KeyValue
	LargestNumber collatz.KeyValue
}

// Run executes the Collatz Conjecture bruteforce calculations.
func Run(ctx context.Context, request Request) (*Results, error) {
	req, err := handler(request)
	if err != nil {
		return nil, err
	}

	errGroup, errGroupCtx := errgroup.WithContext(ctx)
	errGroup.SetLimit(activeRoutines)

	sem := make(chan struct{}, req.batchSize.Int64())

	startingPoint := new(big.Int).Set(req.start)
	for index := startingPoint; index.Cmp(req.end) < 0; index.Add(index, req.calcVals.one) {
		if errGroupCtx.Err() != nil {
			break
		}

		isBreakPoint := new(big.Int).Mod(new(big.Int).Add(index, req.calcVals.one), req.breakPoint).
			Cmp(req.calcVals.zero) == 0
		number := new(big.Int).Add(index, req.calcVals.one)
		errGroup.TryGo(func() error {
			col := collatz.New(number.String())
			return routine(errGroupCtx, sem, col, isBreakPoint, request)
		})
	}

	if errGroupCtx.Err() != nil {
		return nil, fmt.Errorf("errgroup err: %w", errGroupCtx.Err())
	}

	if wgErr := extension.WaitErrHandling(errGroup.Wait()); wgErr != nil {
		return nil, fmt.Errorf("run error: %w", wgErr)
	}

	return &Results{Num: startingPoint.String()}, nil
}

// Request is the input for the Run bruteforce function.
type Request struct {
	Start       string
	End         string
	BatchSize   string
	PrintAll    bool
	Logger      *slog.Logger
	EnableSteps bool
}

type requestValues struct {
	start, end, batchSize, breakPoint *big.Int
	calcVals                          bigIntValues
}

func newRequest(request Request) (*requestValues, error) {
	reqVals := new(requestValues)
	reqVals.calcVals = bigIntValues{
		zero: big.NewInt(0),
		one:  big.NewInt(1),
	}

	start, err := reqVals.value("start", request.Start, false)
	if err != nil {
		return nil, err
	}
	reqVals.start = start

	end, err := reqVals.value("end", request.End, true)
	if err != nil {
		return nil, err
	}
	reqVals.end = end

	batchSize, err := reqVals.value("batch size", request.BatchSize, true)
	if err != nil {
		return nil, err
	}
	reqVals.batchSize = batchSize

	return reqVals, nil
}

func (req *requestValues) value(name, val string, validate bool) (*big.Int, error) {
	bigVal, exists := new(big.Int).SetString(val, collatz.Base)
	if !exists || (validate && bigVal.Cmp(req.calcVals.zero) == 0) {
		return nil, fmt.Errorf("failed to set %s value: %s", name, val)
	}

	return bigVal, nil
}

func handler(request Request) (*requestValues, error) {
	req, err := newRequest(request)
	if err != nil {
		return nil, err
	}

	difference := new(big.Int).Sub(req.end, req.start)
	if difference.Cmp(req.calcVals.zero) <= 0 {
		return nil, fmt.Errorf(
			"difference between start and end must be greater than 0, start: %s, end: %s",
			req.start,
			req.end,
		)
	}

	if difference.Cmp(req.batchSize) < 0 {
		req.batchSize.Set(req.end)
	}

	req.breakPoint = new(big.Int).Div(difference, req.batchSize)
	if request.PrintAll ||
		difference.Cmp(big.NewInt(extension.DefaultBatchSize)) < 0 ||
		req.breakPoint.Cmp(req.calcVals.one) < 0 {
		req.breakPoint.Set(req.calcVals.one)
	}

	return req, nil
}

func routine(ctx context.Context, sem chan struct{}, col *collatz.Collatz, breakPoint bool, request Request) error {
	sem <- struct{}{}
	defer func() { <-sem }()
	if ctx.Err() != nil {
		return fmt.Errorf("routine: num=%s, %w", col.Number(), ctx.Err())
	}

	if err := col.CalculateWithContext(ctx, request.EnableSteps && breakPoint); err != nil {
		return fmt.Errorf("calculate: %w", err)
	}

	if col.Success() {
		return collatz.NewSuccessErr(col.String())
	}

	if request.Logger != nil && breakPoint {
		request.Logger.InfoContext(ctx, col.String())
	}

	return nil
}
