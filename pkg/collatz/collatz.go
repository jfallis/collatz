// Package collatz provides functionality to solve the Collatz Conjecture problem.
// The Collatz Conjecture is a mathematical problem that is easy to understand but difficult to solve.
// More information about the Collatz Conjecture can be found at:
// https://en.wikipedia.org/wiki/Collatz_conjecture
// https://www.youtube.com/watch?v=094y1Z2wpJg
package collatz

import (
	"context"
	"fmt"
	"math/big"
	"sync"
)

const (
	SeqOne   = 4
	SeqTwo   = 2
	SeqThree = 1
	Base     = 10
)

var (
	multiplication  = big.NewInt(3)
	increment       = big.NewInt(1)
	minimum         = big.NewInt(1)
	defaultResponse = Steps{
		big.NewInt(SeqOne),
		big.NewInt(SeqTwo),
		big.NewInt(SeqThree),
	}
)

var (
	errInvalidNumber     error
	errInvalidNumberOnce sync.Once
)

type KeyValue struct {
	Key   int
	Value any
}

func (kv *KeyValue) String() string {
	return fmt.Sprintf("key: %d, value: %s", kv.Key, kv.Value)
}

// ErrInvalidNumber is an error that is returned when the input number is less than the minimum value.
func ErrInvalidNumber() error {
	errInvalidNumberOnce.Do(func() {
		errInvalidNumber = fmt.Errorf("number must be greater than or equal to %d", minimum)
	})

	return errInvalidNumber
}

// SuccessError is a custom error type that is returned when the Collatz Conjecture problem is successfully solved.
type SuccessError struct {
	Num string
}

func NewSuccessErr(s string) SuccessError {
	return SuccessError{Num: s}
}

// Error method for the SuccessError type.
func (n SuccessError) Error() string {
	return fmt.Sprintf("🎉 did you solve the collatz conjecture: %s", n.Num)
}

type Collatz struct {
	number           string
	seq              *big.Int
	err              error
	steps            Steps
	calculateStarted bool
}

func New(num string) *Collatz {
	return &Collatz{
		number: num,
		steps:  make(Steps, 0),
	}
}

func (c *Collatz) Err() error {
	err := c.err
	return err
}

func (c *Collatz) Calculate(enableSteps bool) error {
	return c.CalculateWithContext(context.Background(), enableSteps)
}

func (c *Collatz) CalculateWithContext(ctx context.Context, enableSteps bool) error {
	defer func() {
		if c.calculateStarted && c.seq.Cmp(minimum) != 0 {
			c.err = NewSuccessErr(c.number)
		}
		if r := recover(); r != nil {
			c.err = fmt.Errorf("%s: %w", r, c.err)
		}
	}()

	var ok bool
	c.seq, ok = new(big.Int).SetString(c.number, Base)
	if !ok {
		return ErrInvalidNumber()
	}
	numberComparison := c.seq.Cmp(minimum)
	if numberComparison < 0 {
		return ErrInvalidNumber()
	}

	if numberComparison == 0 {
		c.steps = defaultResponse
		return nil
	}

	c.calculateStarted = true

	for c.seq.Cmp(minimum) != 0 {
		select {
		case <-ctx.Done():
			return NewSuccessErr(c.number)
		default:
			c.sequence()
			if enableSteps {
				c.steps = append(c.steps, new(big.Int).Set(c.seq))
			}
		}
	}

	return nil
}

func (c *Collatz) sequence() {
	if c.seq.Bit(0) == 0 {
		c.seq.Rsh(c.seq, 1)
		return
	}

	c.seq.Mul(c.seq, multiplication)
	c.seq.Add(c.seq, increment)
}

func (c *Collatz) Number() string {
	return c.number
}

func (c *Collatz) Steps() Steps {
	return c.steps
}

type Steps []*big.Int

func (s Steps) String() string {
	var steps []string
	for _, step := range s {
		steps = append(steps, step.String())
	}

	return fmt.Sprintf("%v", steps)
}

func (s Steps) MaxStepValue() KeyValue {
	i, val := Max(s)
	return KeyValue{
		Key:   i,
		Value: val,
	}
}

func (c *Collatz) Success() bool {
	return c.calculateStarted && c.seq.Cmp(minimum) != 0
}

func (c *Collatz) String() string {
	return fmt.Sprintf("number: %s, steps: %d, max: %s, success: %t", c.Number(), len(c.Steps()), c.Steps().MaxStepValue().Value, c.Success())
}

func Max(slice []*big.Int) (key int, value string) {
	if len(slice) == 0 {
		return key, "-1"
	}

	key = -1
	maxVal := big.NewInt(-1)

	for x, y := range slice {
		comp := y.Cmp(maxVal)
		if comp > 0 {
			key = x
			maxVal = y
		}
	}

	return key, maxVal.String()
}
