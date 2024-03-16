package calculation_test

import (
	"github.com/jfallis/collatz/pkg/calculation"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLeadingDigitEstimate(t *testing.T) {
	t.Parallel()

	assert.Equal(t, uint64(5), calculation.LeadingDigitEstimate(5))
	assert.Equal(t, uint64(10), calculation.LeadingDigitEstimate(10))
	assert.Equal(t, uint64(10), calculation.LeadingDigitEstimate(12))
	assert.Equal(t, uint64(300), calculation.LeadingDigitEstimate(320))
	assert.Equal(t, uint64(100), calculation.LeadingDigitEstimate(184))
	assert.Equal(t, uint64(100), calculation.LeadingDigitEstimate(148))
	assert.Equal(t, uint64(9000), calculation.LeadingDigitEstimate(9098))
}

func BenchmarkLeadingDigit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		calculation.LeadingDigitEstimate(35468743)
	}
}
