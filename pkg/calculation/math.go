package calculation

import (
	"math"
)

const (
	Base                = 10
	Bit                 = 64
	leadingDigitMinimum = 10
)

type Statement interface {
	Number() uint64
	Sequence(val uint64) uint64
	Calculate() error
	Steps() []uint64
	Success() bool
}

func LeadingDigitEstimate(value uint64) uint64 {
	if value < leadingDigitMinimum {
		return value
	}

	digits := math.Log10(float64(value))
	pow := math.Pow10(int(digits))

	for value >= leadingDigitMinimum {
		value /= Base
	}

	value *= uint64(pow)

	return value
}

type IntCallback func(uint64) uint64
