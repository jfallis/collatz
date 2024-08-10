package extension_test

import (
	"bytes"
	"errors"
	"log/slog"
	"math/big"
	"testing"

	"github.com/jfallis/collatz/pkg/collatz"

	"github.com/jfallis/collatz/pkg/collatz/extension"
	"github.com/stretchr/testify/assert"
)

func TestWaitErrHandling(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		inputErr    error
		expectedErr string
	}{
		"Test with nil error": {
			inputErr:    nil,
			expectedErr: "",
		},
		"Test with non-success error": {
			inputErr:    errors.New("some error"),
			expectedErr: "bruteforce failed: some error",
		},
		"Test with success error": {
			inputErr:    collatz.SuccessError{Number: big.NewInt(5), Steps: []string{"1", "2"}},
			expectedErr: "successfully found the number You found an infinite loop 🎉 number: 5, steps: [1 2]",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := extension.WaitErrHandling(test.inputErr)
			if test.expectedErr == "" {
				assert.NoError(t, err)
				return
			}

			assert.Error(t, err)
			assert.Equal(t, test.expectedErr, err.Error())
		})
	}
}

func TestBruteforce(t *testing.T) {
	tests := map[string]struct {
		end           *big.Int
		maxBatchCount *big.Int
		expectedLogs  []string
	}{
		"test input end 5, maxBatchCount 10": {
			big.NewInt(5),
			big.NewInt(10),
			[]string{"bruteforce number: 5, steps: 5"},
		},
		"test input end 5, maxBatchCount 5": {
			big.NewInt(5),
			big.NewInt(5),
			[]string{"bruteforce number: 5, steps: 5"},
		},
		"test input end 10, maxBatchCount 2": {
			big.NewInt(10),
			big.NewInt(2),
			[]string{
				"bruteforce number: 2, steps: 1",
				"bruteforce number: 4, steps: 2",
				"bruteforce number: 6, steps: 8",
				"bruteforce number: 8, steps: 3",
				"bruteforce number: 10, steps: 6",
			},
		},
		"test input end 10, maxBatchCount 3": {
			big.NewInt(10),
			big.NewInt(3),
			[]string{
				"bruteforce number: 3, steps: 7",
				"bruteforce number: 6, steps: 8",
				"bruteforce number: 9, steps: 19",
			},
		},
	}

	start := big.NewInt(0)
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := slog.New(slog.NewJSONHandler(&buf, nil))
			slog.SetDefault(logger)

			actual, err := extension.Bruteforce(start, test.end, test.maxBatchCount, true)

			assert.NoError(t, err)
			assert.Equal(t, test.end, actual)
			assert.NotEmpty(t, buf.String())

			for _, expectedLog := range test.expectedLogs {
				assert.Contains(t, buf.String(), expectedLog)
			}
		})
	}
}

func TestBruteforceLargestStepCount(t *testing.T) {
	start := big.NewInt(0)
	end := big.NewInt(30)
	maxBatchCount := big.NewInt(1)

	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	slog.SetDefault(logger)

	actual, err := extension.Bruteforce(start, end, maxBatchCount, true)
	assert.NoError(t, err)

	assert.Equal(t, end, actual)
	assert.NotEmpty(t, buf.String())
	assert.Contains(t, buf.String(), "bruteforce number: 27, steps: 111")
}

func BenchmarkBruteforce(b *testing.B) {
	start := big.NewInt(0)
	end := big.NewInt(1_000)
	batchSize := new(big.Int).Div(end, big.NewInt(10))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := extension.Bruteforce(start, end, batchSize, false)
		assert.NoError(b, err)
	}
}
