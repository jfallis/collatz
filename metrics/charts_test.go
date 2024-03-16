package metrics_test

import (
	"testing"

	"github.com/jfallis/collatz/metrics"
	"github.com/stretchr/testify/assert"
)

func TestBuildMetrics(t *testing.T) {
	t.Parallel()

	assert.Equal(
		t,
		[][]float64{{0, 0}, {1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5}, {6, 6}, {7, 7}, {8, 8}, {9, 9}},
		metrics.BuildMetrics([]float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}),
	)
}

func BenchmarkBuildMetrics(b *testing.B) {
	for i := 0; i < b.N; i++ {
		metrics.BuildMetrics([]float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	}
}
