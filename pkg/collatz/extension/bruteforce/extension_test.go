package bruteforce_test

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"strconv"
	"strings"
	"testing"

	"github.com/jfallis/collatz/pkg/collatz/extension"

	"github.com/jfallis/collatz/pkg/collatz/extension/bruteforce"
	"github.com/stretchr/testify/assert"
)

func TestCtxDone(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := bruteforce.Run(ctx, bruteforce.Request{
		Start:      "0",
		End:        "1000",
		BatchSize:  "10",
		EnableStep: true,
	})

	assert.Error(t, err)
	assert.Equal(t, fmt.Errorf("routine failed: %w", context.Canceled), err)
}

func TestInvalidArguments(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		start, end, batchSize string
		expectedError         string
	}{
		"invalid start value": {
			start:         "invalid",
			end:           "10",
			batchSize:     "10",
			expectedError: "failed to set start value: invalid",
		},
		"empty start value": {
			start:         "",
			end:           "10",
			batchSize:     "10",
			expectedError: "failed to set start value: ",
		},
		"invalid end value": {
			start:         "0",
			end:           "invalid",
			batchSize:     "10",
			expectedError: "failed to set end value: invalid",
		},
		"empty end value": {
			start:         "0",
			end:           "",
			batchSize:     "10",
			expectedError: "failed to set end value: ",
		},
		"invalid batch size value": {
			start:         "0",
			end:           "10",
			batchSize:     "invalid",
			expectedError: "failed to set batch size value: invalid",
		},
		"empty batch size value": {
			start:         "0",
			end:           "10",
			batchSize:     "",
			expectedError: "failed to set batch size value: ",
		},
		"invalid difference": {
			start:         "10",
			end:           "10",
			batchSize:     "10",
			expectedError: "difference between start and end must be greater than 0",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := bruteforce.Run(context.Background(), bruteforce.Request{
				Start:      tc.start,
				End:        tc.end,
				BatchSize:  tc.batchSize,
				EnableStep: true,
			})
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedError)
		})
	}
}

func TestDefaultBatchSize(t *testing.T) {
	t.Parallel()

	assert.Equal(t, big.NewInt(1000), extension.DefaultBatchSize())
}

func TestBruteforce(t *testing.T) {
	type expected struct {
		number string
		logs   []string
	}
	testCases := map[string]struct {
		expected
		maxBatchCount string
	}{
		"start value 5, maxBatchCount 10": {
			expected{
				number: "5",
				logs: []string{
					"number: 1, steps: 3",
					"number: 2, steps: 1",
					"number: 3, steps: 7",
					"number: 4, steps: 2",
					"number: 5, steps: 5",
				},
			},
			"10",
		},
		"start value 5, maxBatchCount 5": {
			expected{
				number: "5",
				logs: []string{
					"number: 1, steps: 3",
					"number: 2, steps: 1",
					"number: 3, steps: 7",
					"number: 4, steps: 2",
					"number: 5, steps: 5",
				},
			},
			"5",
		},
		"start value 10, maxBatchCount 2": {
			expected{
				number: "10",
				logs: []string{
					"number: 1, steps: 3",
					"number: 2, steps: 1",
					"number: 3, steps: 7",
					"number: 4, steps: 2",
					"number: 5, steps: 5",
					"number: 6, steps: 8",
					"number: 7, steps: 16",
					"number: 8, steps: 3",
					"number: 9, steps: 19",
					"number: 10, steps: 6",
				},
			},
			"2",
		},
		"start value 10, maxBatchCount 3": {
			expected{
				number: "10",
				logs: []string{
					"number: 1, steps: 3",
					"number: 2, steps: 1",
					"number: 3, steps: 7",
					"number: 4, steps: 2",
					"number: 5, steps: 5",
					"number: 6, steps: 8",
					"number: 7, steps: 16",
					"number: 8, steps: 3",
					"number: 9, steps: 19",
					"number: 10, steps: 6",
				},
			},
			"3",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer

			logger := slog.New(slog.NewTextHandler(&buf, nil))
			slog.SetDefault(logger)

			actual, err := bruteforce.Run(context.Background(), bruteforce.Request{
				Start:      "0",
				End:        tc.expected.number,
				BatchSize:  tc.maxBatchCount,
				Logging:    true,
				EnableStep: true,
			})

			assert.NoError(t, err)
			assert.Equal(t, tc.expected.number, actual.Number)
			assert.NotEmpty(t, buf.String())

			assert.Len(t, tc.expected.logs, strings.Count(buf.String(), "\n"))
			for _, expectedLog := range tc.expected.logs {
				assert.Contains(t, buf.String(), expectedLog)
			}
		})
	}
}

func TestLargeStepCount(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))
	slog.SetDefault(logger)

	actual, err := bruteforce.Run(context.Background(), bruteforce.Request{
		Start:      "0",
		End:        "30",
		BatchSize:  "1",
		Logging:    true,
		EnableStep: true,
	})
	assert.NoError(t, err)

	assert.Equal(t, "30", actual.Number)
	assert.NotEmpty(t, buf.String())
	assert.Contains(t, buf.String(), "number: 27, steps: 111")
}

func BenchmarkWithStepsEnabled(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := bruteforce.Run(context.Background(), bruteforce.Request{
			Start: "0",
			End:   "1000", BatchSize: strconv.Itoa(1000 / 10),
			EnableStep: true,
		})
		assert.NoError(b, err)
	}
}

func BenchmarkWithoutStepsEnabled(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := bruteforce.Run(context.Background(), bruteforce.Request{
			Start: "0",
			End:   "1000", BatchSize: strconv.Itoa(1000 / 10),
		})
		assert.NoError(b, err)
	}
}
