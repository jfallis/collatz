package calculation_test

import (
	"testing"

	"github.com/jfallis/collatz/calculation"
	"github.com/stretchr/testify/assert"
)

func TestLeadingDigit(t *testing.T) {
	t.Parallel()

	assert.Equal(t, uint64(0), calculation.LeadingDigit(0))
	assert.Equal(t, uint64(1), calculation.LeadingDigit(1))
	assert.Equal(t, uint64(5), calculation.LeadingDigit(5))
	assert.Equal(t, uint64(1), calculation.LeadingDigit(10))
	assert.Equal(t, uint64(1), calculation.LeadingDigit(100))
	assert.Equal(t, uint64(1), calculation.LeadingDigit(102))
	assert.Equal(t, uint64(9), calculation.LeadingDigit(9098))
}

func TestCreateHistogram(t *testing.T) {
	t.Parallel()

	actual := calculation.CreateHistogram(
		[]uint64{1, 2, 21, 3, 31, 311, 4, 41, 411, 4111},
		calculation.LeadingDigit,
	)
	assert.Equal(t, []uint64{0, 1, 2, 3, 4, 0, 0, 0, 0, 0}, actual)
}

func TestConvertIntToFloat(t *testing.T) {
	t.Parallel()

	actual := calculation.ConvertIntToFloat([]uint64{1, 2, 21, 3, 31, 311, 4, 41, 411, 4111})
	assert.Equal(t, []float64{1, 2, 21, 3, 31, 311, 4, 41, 411, 4111}, actual)
}

func TestConvertFloatToInt(t *testing.T) {
	t.Parallel()

	actual := calculation.ConvertFloatToInt([]float64{1, 2, 21, 3, 31, 311, 4, 41, 411, 4111})
	assert.Equal(t, []uint64{1, 2, 21, 3, 31, 311, 4, 41, 411, 4111}, actual)
}

func BenchmarkLeadingDigitUnder10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		calculation.LeadingDigit(9)
	}
}

func BenchmarkLeadingDigit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		calculation.LeadingDigit(35468743)
	}
}

func BenchmarkConvertIntToFloat(b *testing.B) {
	example := []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 11, 12}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		calculation.ConvertIntToFloat(example)
	}
}

func BenchmarkConvertFloatToInt(b *testing.B) {
	example := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 11, 12}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		calculation.ConvertFloatToInt(example)
	}
}
