package math_test

import (
	"testing"

	"github.com/jfallis/collatz/pkg/math"

	"github.com/stretchr/testify/assert"
)

func TestLeadingDigitEstimate(t *testing.T) {
	t.Parallel()

	assert.Equal(t, uint64(5), math.LeadingDigitEstimate(5))
	assert.Equal(t, uint64(10), math.LeadingDigitEstimate(10))
	assert.Equal(t, uint64(10), math.LeadingDigitEstimate(12))
	assert.Equal(t, uint64(300), math.LeadingDigitEstimate(320))
	assert.Equal(t, uint64(100), math.LeadingDigitEstimate(184))
	assert.Equal(t, uint64(100), math.LeadingDigitEstimate(148))
	assert.Equal(t, uint64(9000), math.LeadingDigitEstimate(9098))
}

func BenchmarkLeadingDigit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		math.LeadingDigitEstimate(35468743)
	}
}
