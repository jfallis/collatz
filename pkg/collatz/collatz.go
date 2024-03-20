// Package collatz provides functionality to solve the Collatz Conjecture problem.
// The Collatz Conjecture is a mathematical problem that is easy to understand but difficult to solve.
// More information about the Collatz Conjecture can be found at:
// https://en.wikipedia.org/wiki/Collatz_conjecture
// https://www.youtube.com/watch?v=094y1Z2wpJg
package collatz

import (
	"fmt"
	"math/big"
)

const (
	SuccessMsg = "You found an infinite loop 🎉"
	StepsLimit = 1_000_000
)

var (
	multiplication = big.NewInt(3)
	increment      = big.NewInt(1)
	minimum        = big.NewInt(1)
)

// ErrInvalidNumber is an error that is returned when the input number is less than the minimum value.
var ErrInvalidNumber = fmt.Errorf("number must be greater than or equal to %d", minimum)

// SuccessError is a custom error type that is returned when the Collatz Conjecture problem is successfully solved.
type SuccessError struct {
	Number *big.Int
	Steps  []*big.Int
}

// Error method for the SuccessError type.
func (e SuccessError) Error() string {
	return fmt.Sprintf("%s number: %d, steps: %+v", SuccessMsg, e.Number, e.Steps)
}

type Collatz struct {
	number *big.Int
	steps  []*big.Int
	cache  map[string][]*big.Int
}

func New(num *big.Int) *Collatz {
	return &Collatz{
		number: new(big.Int).Set(num),
		steps:  make([]*big.Int, 0),
		cache:  make(map[string][]*big.Int),
	}
}

func (c *Collatz) Calculate() error {
	if c.number.Cmp(minimum) < 0 {
		return ErrInvalidNumber
	}

	counter := 0
	num := new(big.Int).Set(c.number)

	for num.Cmp(minimum) != 0 || (counter == 0 && counter <= StepsLimit) {
		c.Sequence(num)
		c.steps = append(c.steps, new(big.Int).Set(num))
		counter++
	}

	return nil
}

func (c *Collatz) Sequence(val *big.Int) {
	if new(big.Int).Set(val).Bit(0) == 0 {
		val.Rsh(val, 1)
		return
	}

	val.Mul(val, multiplication)
	val.Add(val, increment)
}

func (c *Collatz) Number() *big.Int {
	return c.number
}

func (c *Collatz) Steps() []*big.Int {
	return c.steps
}

func (c *Collatz) Success() bool {
	length := len(c.Steps())

	return length != 0 && c.Steps()[length-1].Cmp(minimum) != 0
}
