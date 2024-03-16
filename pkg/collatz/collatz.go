// Package collatz provides functionality to solve the Collatz Conjecture problem.
// The Collatz Conjecture is a mathematical problem that is easy to understand but difficult to solve.
// More information about the Collatz Conjecture can be found at:
// https://en.wikipedia.org/wiki/Collatz_conjecture
// https://www.youtube.com/watch?v=094y1Z2wpJg
package collatz

import (
	"fmt"
	"sync"
)

const (
	divide         = 2
	multiplication = 3
	increment      = 1
	minimum        = 1
	SuccessMsg     = "You found an infinite loop 🎉"
)

// ErrInvalidNumber is an error that is returned when the input number is less than the minimum value.
var ErrInvalidNumber = fmt.Errorf("number must be greater than or equal to %d", minimum)

// SuccessError is a custom error type that is returned when the Collatz Conjecture problem is successfully solved.
type SuccessError struct {
	Number uint64
	Steps  []uint64
}

// Error method for the SuccessError type.
func (e SuccessError) Error() string {
	return fmt.Sprintf("%s number: %d, steps: %+v", SuccessMsg, e.Number, e.Steps)
}

type Collatz struct {
	number uint64
	steps  []uint64
	cache  map[uint64][]uint64
	mu     sync.RWMutex
}

func New(num uint64) *Collatz {
	return &Collatz{
		number: num,
		steps:  make([]uint64, 0),
		cache:  make(map[uint64][]uint64),
	}
}

func (c *Collatz) Calculate() error {
	if c.number < minimum {
		return ErrInvalidNumber
	}

	counter := 0
	num := c.number

	for num != minimum || counter == 0 {
		c.mu.RLock()
		cachedSteps, ok := c.cache[num]
		c.mu.RUnlock()

		if ok {
			c.steps = append(c.steps, cachedSteps...)
			break
		}

		num = c.Sequence(num)
		c.steps = append(c.steps, num)
		counter++
	}

	c.mu.Lock()
	c.cache[c.number] = c.steps
	c.mu.Unlock()

	return nil
}

func (c *Collatz) Sequence(val uint64) uint64 {
	if val%divide == 0 {
		return val / divide
	}

	return val*multiplication + increment
}

func (c *Collatz) Number() uint64 {
	return c.number
}

func (c *Collatz) Steps() []uint64 {
	return c.steps
}

func (c *Collatz) Success() bool {
	length := len(c.Steps())

	return length != 0 && c.Steps()[length-1] != 1
}
