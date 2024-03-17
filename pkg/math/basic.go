package math

import (
	"math"
)

const (
	Base                = 10
	leadingDigitMinimum = 10
)

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
