// Package calculation provides simple math methods.
// The Collatz Conjecture is the simplest math problem no one can solve;
// it is easy enough for almost anyone to understand but notoriously difficult to solve.
//
// https://en.wikipedia.org/wiki/Collatz_conjecture
// https://www.youtube.com/watch?v=094y1Z2wpJg
package calculation

import (
	"context"
	"fmt"
	"reflect"

	"github.com/jfallis/collatz/domain"
	"golang.org/x/sync/errgroup"
)

const (
	maxIterations = 10000
	threads       = 1200
)

const (
	CollatzSuccessError = domain.CollatzError("You found an infinite loop 🎉")
)

const (
	divide         = 2
	multiplication = 3
	increment      = 1
	minimum        = 1
)

type Collatz struct {
	Number     uint64
	Steps      int
	HailStones []float64
}

func Create(num uint64) Statement {
	c := &Collatz{
		Number:     num,
		Steps:      0,
		HailStones: make([]float64, 0),
	}
	c.Execute()

	return c
}

func (c *Collatz) Input() float64 {
	return float64(c.Number)
}

// Execute the Collatz conjecture for a given number until the 4-2-1 loop.
func (c *Collatz) Execute() {
	if c.Number <= minimum {
		return
	}

	c.Steps = 0
	c.HailStones = make([]float64, maxIterations)
	c.HailStones[0] = float64(c.Number)

	x := 0

	for c.HailStones[x] > minimum && x < maxIterations {
		c.HailStones[x+1] = c.Calculate(c.HailStones[x])
		x++
	}

	c.HailStones = c.HailStones[1 : x+1]

	c.Steps = len(c.HailStones)
}

func (c *Collatz) Calculate(val float64) float64 {
	if uint64(val)%divide == 0 {
		return val / divide
	}

	return val*multiplication + increment
}

func (c *Collatz) Len() int {
	return c.Steps
}

func (c *Collatz) Results() []float64 {
	return c.HailStones
}

func Success(hailStones []float64) bool {
	pattern := []float64{4, 2, 1}
	x := len(pattern)
	y := len(hailStones)

	if y < x || (y >= x && reflect.DeepEqual(hailStones[y-x:], pattern)) {
		return false
	}

	return true
}

func Bruteforce(num uint64) (uint64, error) {
	g, _ := errgroup.WithContext(context.Background())

	maxThreads := uint64(threads)
	if maxThreads > num {
		maxThreads = 1
	}

	for p := uint64(0); p < (num / maxThreads); p++ {
		var i uint64

		start := (p + 1) * maxThreads
		end := (p + 2) * maxThreads

		for i = start; i < end; i++ {
			i := i

			g.Go(func() error {
				c := Create(i)
				if success := Success(c.Results()); success {
					return CollatzSuccessError
				}

				return nil
			})
		}

		err := g.Wait()
		if err != nil {
			return i, fmt.Errorf("successfully found the number %w", err)
		}
	}

	return 0, nil
}
