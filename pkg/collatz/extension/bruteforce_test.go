package extension_test

import (
	"bytes"
	"fmt"
	"log/slog"
	"testing"

	"github.com/jfallis/collatz/pkg/collatz/extension"
	"github.com/stretchr/testify/assert"
)

func TestBruteforce(t *testing.T) {
	t.Parallel()

	tests := []struct {
		number        uint64
		maxBatchCount uint64
	}{
		{5, 10},
		{5, 5},
		{10, 2},
		{10, 3},
	}

	for _, test := range tests {
		name := fmt.Sprintf("test input number %d, maxBatchCount %d", test.number, test.maxBatchCount)
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := slog.New(slog.NewJSONHandler(&buf, nil))
			slog.SetDefault(logger)

			actual := extension.Bruteforce(test.number, test.maxBatchCount, true)

			assert.NoError(t, actual)
			for i := uint64(0); i < test.number; i++ {
				if (i+1)%test.maxBatchCount == 0 {
					assert.Containsf(t, buf.String(), fmt.Sprintf("bruteforce number: %d", i+1), "expected log to contain bruteforce number: %d", i)
				}
			}
		})
	}
}

func BenchmarkBruteforce(b *testing.B) {
	num := uint64(1_000)
	batchSize := num / 10

	for i := 0; i < b.N; i++ {
		assert.NoError(b, extension.Bruteforce(num, batchSize, false))
	}
}
