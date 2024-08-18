// Package collatz provides functionality to solve the Collatz Conjecture problem.
// The Collatz Conjecture is a mathematical problem that is easy to understand but difficult to solve.
// More information about the Collatz Conjecture can be found at:
// https://en.wikipedia.org/wiki/Collatz_conjecture
// https://www.youtube.com/watch?v=094y1Z2wpJg
package collatz

import (
	"fmt"
	"math/big"
	"os"
	"sync"
)

const (
	SuccessMsg = "You found an infinite loop 🎉"
	SeqOne     = 4
	SeqTwo     = 2
	SeqThree   = 1
	Base       = 10
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
	String string
}

// Error method for the SuccessError type.
func (e SuccessError) Error() string {
	return fmt.Sprintf("%s - %s", SuccessMsg, e.String)
}

type Collatz struct {
	number string
	steps  Steps
}

func New(num string) *Collatz {
	return &Collatz{
		number: num,
		steps:  make(Steps, 0),
	}
}

func (c *Collatz) Calculate(enableSteps bool) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic: %s\n", c.String())
			fmt.Printf("Recover message: %+v\n", r)
			os.Exit(1)
		}
	}()

	num, ok := new(big.Int).SetString(c.number, Base)
	if !ok {
		return ErrInvalidNumber()
	}
	numberComparison := num.Cmp(minimum)
	if numberComparison < 0 {
		return ErrInvalidNumber()
	}

	if numberComparison == 0 {
		c.steps = defaultResponse
		return nil
	}

	for num.Cmp(minimum) != 0 {
		c.Sequence(num)

		if enableSteps {
			c.steps = append(c.steps, new(big.Int).Set(num))
		}
	}

	return nil
}

func (c *Collatz) Sequence(val *big.Int) {
	if val.Bit(0) == 0 {
		val.Rsh(val, 1)
		return
	}

	val.Mul(val, multiplication)
	val.Add(val, increment)
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
	length := len(c.Steps())
	return length != 0 && c.Steps()[length-1].Cmp(minimum) != 0
}

func (c *Collatz) String() string {
	return fmt.Sprintf("number: %s, steps: %d, max: %s", c.Number(), len(c.Steps()), c.Steps().MaxStepValue().Value)
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
