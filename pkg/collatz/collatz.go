// Package collatz provides functionality to solve the Collatz Conjecture problem.
// The Collatz Conjecture is a mathematical problem that is easy to understand but difficult to solve.
// More information about the Collatz Conjecture can be found at:
// https://en.wikipedia.org/wiki/Collatz_conjecture
// https://www.youtube.com/watch?v=094y1Z2wpJg
package collatz

import (
	"context"
	"errors"
	"fmt"
	"math/big"
)

const (
	seqOne   = 4
	seqTwo   = 2
	seqThree = 1
)

const (
	mul  = 3
	inc  = 1
	mini = 1
	// Base is the numerical Base used for parsing input numbers.
	Base = 10
)

var (
	// ErrInvalidNumber is returned when the input number is less than 1 or cannot be parsed.
	ErrInvalidNumber = errors.New("number must be greater than or equal to 1")
)

// KeyValue is a struct that holds a key-value pair where the key is an integer and the value can be of any type.
type KeyValue struct {
	Key   int
	Value any
}

func (kv *KeyValue) String() string {
	return fmt.Sprintf("key: %d, value: %s", kv.Key, kv.Value)
}

// SuccessError is a custom error type that is returned when the Collatz Conjecture problem is successfully solved.
type SuccessError struct {
	Num string
}

// NewSuccessErr creates a new instance of SuccessError when the conjecture is solved.
func NewSuccessErr(s string) SuccessError {
	return SuccessError{Num: s}
}

// Error method for the SuccessError type.
func (n SuccessError) Error() string {
	return fmt.Sprintf("🎉 did you solve the collatz conjecture: %s", n.Num)
}

type collatzConfig struct {
	mul *big.Int
	inc *big.Int
	min *big.Int
}

type steps []*big.Int

func newCollatzConfig() *collatzConfig {
	return &collatzConfig{mul: big.NewInt(mul), inc: big.NewInt(inc), min: big.NewInt(mini)}
}

// Collatz holds the configuration and state for calculating the Collatz sequence.
type Collatz struct {
	config     *collatzConfig
	num        string
	seq        *big.Int
	err        error
	steps      steps
	hasStarted bool
}

// New creates a new instance of Collatz with the given number as a string.
func New(num string) *Collatz {
	return &Collatz{
		config: newCollatzConfig(),
		num:    num,
		steps:  make(steps, 0),
	}
}

// Err returns the error encountered during the calculation, if any.
func (c *Collatz) Err() error {
	return c.err
}

// Calculate computes the Collatz sequence for the given number.
func (c *Collatz) Calculate(enableSteps bool) error {
	return c.CalculateWithContext(context.Background(), enableSteps)
}

// CalculateWithContext computes the Collatz sequence for the given number with context support.
func (c *Collatz) CalculateWithContext(ctx context.Context, enableSteps bool) error {
	defer c.updateErrorState()

	var ok bool
	c.seq, ok = new(big.Int).SetString(c.num, Base)
	if !ok {
		return fmt.Errorf("invalid num: %s; %w", c.num, ErrInvalidNumber)
	}
	numberComparison := c.seq.Cmp(c.config.min)
	if numberComparison < 0 {
		return fmt.Errorf("invalid num: %s; %w", c.num, ErrInvalidNumber)
	}

	if numberComparison == 0 {
		c.steps = steps{big.NewInt(seqOne), big.NewInt(seqTwo), big.NewInt(seqThree)}
		return nil
	}

	c.hasStarted = true

	for c.seq.Cmp(c.config.min) != 0 {
		select {
		case <-ctx.Done():
			return NewSuccessErr(c.num)
		default:
			c.sequence()
			if enableSteps {
				c.steps = append(c.steps, new(big.Int).Set(c.seq))
			}
		}
	}

	return nil
}

// Number returns the original number input as a string.
func (c *Collatz) Number() string {
	return c.num
}

// Steps returns the sequence of steps taken to reach 1 in the Collatz sequence.
func (c *Collatz) Steps() []string {
	vals := make([]string, len(c.steps))
	for i, step := range c.steps {
		vals[i] = fmt.Sprintf("%v", step)
	}

	return vals
}

// StepsMaxValue returns the index and value of the maximum step in the Collatz sequence.
func (c *Collatz) StepsMaxValue() (key int, value string) {
	key = -1
	value = "-1"
	if len(c.steps) == 0 {
		return
	}

	maxVal := big.NewInt(-1)
	for x, y := range c.steps {
		comp := y.Cmp(maxVal)
		if comp > 0 {
			key = x
			maxVal = y
		}
	}

	return key, maxVal.String()

}

func (c *Collatz) updateErrorState() {
	if r := recover(); r != nil {
		c.err = fmt.Errorf("collatz panic: %+v", r)
	}
	if c.hasStarted && c.seq.Cmp(c.config.min) != 0 {
		c.err = fmt.Errorf("%w: %w", NewSuccessErr(c.num), c.err)
	}
}

func (c *Collatz) sequence() {
	if c.seq.Bit(0) == 0 {
		c.seq.Rsh(c.seq, 1)
		return
	}

	c.seq.Mul(c.seq, c.config.mul)
	c.seq.Add(c.seq, c.config.inc)
}

// Success returns true if the Collatz sequence successfully reached.
func (c *Collatz) Success() bool {
	return c.hasStarted && c.seq.Cmp(c.config.min) != 0
}

// String returns a string output of the numeric input.
func (c *Collatz) String() string {
	_, s := c.StepsMaxValue()
	return fmt.Sprintf(
		"number: %s, steps: %d, max: %s, success: %t",
		c.Number(), len(c.Steps()), s, c.Success(),
	)
}
