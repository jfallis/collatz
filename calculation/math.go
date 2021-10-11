// calculation package provides simple math methods.
package calculation

import (
	"strconv"
)

const (
	Base = 10
	Bit  = 64
)

const (
	leadingDigitMinimum = 10
)

// Statement represents.
type Statement interface {
	Input() float64
	Execute()
	Calculate(val float64) float64
	Results() []float64
	Len() int
}

// LeadingDigit returns the leading digit of a given number.
func LeadingDigit(value uint64) uint64 {
	if value < leadingDigitMinimum {
		return value
	}

	if leadingDigit, err := strconv.ParseUint(strconv.FormatUint(value, Base)[0:1], Base, Bit); err == nil {
		return leadingDigit
	}

	return 0
}

type IntCallback func(uint64) uint64

// CreateHistogram generates a histogram from a given number of values.
func CreateHistogram(values []uint64, callback IntCallback) []uint64 {
	newValues := make([]uint64, Base)
	for _, num := range values {
		newValues[callback(num)]++
	}

	return newValues
}

// ConvertIntToFloat transforms int slices to float64.
func ConvertIntToFloat(values []uint64) []float64 {
	newValues := make([]float64, len(values))
	for x, v := range values {
		newValues[x] = float64(v)
	}

	return newValues
}

// ConvertFloatToInt transforms float64 slices to int.
func ConvertFloatToInt(values []float64) []uint64 {
	newValues := make([]uint64, len(values))
	for x, v := range values {
		newValues[x] = uint64(v)
	}

	return newValues
}
