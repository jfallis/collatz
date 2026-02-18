package bruteforce_test

import (
	"bytes"
	"context"
	"log/slog"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jfallis/collatz/pkg/collatz/extension/bruteforce"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCtxDone(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	cancel()

	_, err := bruteforce.Run(ctx, bruteforce.Request{
		Start:     "1",
		End:       "1000",
		BatchSize: "10",
	})

	require.Error(t, err)
	deadline, ok := ctx.Deadline()
	assert.True(t, ok)
	assert.WithinDuration(t, time.Now(), deadline, time.Second)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestInvalidArguments(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		start, end, batchSize string
		expectedError         string
	}{
		"invalid start value": {
			start: "invalid", end: "10", batchSize: "10",
			expectedError: "failed to set start value: invalid",
		},
		"empty start value": {
			start: "", end: "10", batchSize: "10",
			expectedError: "failed to set start value: ",
		},
		"zero end value": {
			start: "1", end: "0", batchSize: "10",
			expectedError: "failed to set end value: 0",
		},
		"invalid end value": {
			start: "1", end: "invalid", batchSize: "10",
			expectedError: "failed to set end value: invalid",
		},
		"empty end value": {
			start: "1", end: "", batchSize: "10",
			expectedError: "failed to set end value: ",
		},
		"zero batch size value": {
			start: "1", end: "10", batchSize: "0",
			expectedError: "failed to set batch size value: 0",
		},
		"invalid batch size value": {
			start: "1", end: "10", batchSize: "invalid",
			expectedError: "failed to set batch size value: invalid",
		},
		"empty batch size value": {
			start: "1", end: "10", batchSize: "",
			expectedError: "failed to set batch size value: ",
		},
		"invalid difference": {
			start: "10", end: "10", batchSize: "10",
			expectedError: "difference between start and end must be greater than 0",
		},
	}

	for name, testcase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := bruteforce.Run(context.Background(), bruteforce.Request{
				Start:       testcase.start,
				End:         testcase.end,
				BatchSize:   testcase.batchSize,
				EnableSteps: true,
			})
			require.Error(t, err)
			assert.Contains(t, err.Error(), testcase.expectedError)
		})
	}
}

func TestBruteforceVariousBatchSizes(t *testing.T) {
	t.Parallel()

	type expected struct {
		number string
		logs   []string
	}
	testCases := map[string]struct {
		expected      expected
		maxBatchCount string
	}{
		"start value 5, maxBatchCount 10": {maxBatchCount: "10", expected: expected{number: "5", logs: []string{
			"number: 1, steps: 3", "number: 2, steps: 1", "number: 3, steps: 7", "number: 4, steps: 2", "number: 5, steps: 5",
		}}},
		"start value 5, maxBatchCount 5": {maxBatchCount: "5", expected: expected{number: "5", logs: []string{
			"number: 1, steps: 3", "number: 2, steps: 1", "number: 3, steps: 7", "number: 4, steps: 2", "number: 5, steps: 5",
		}}},
		"start value 10, maxBatchCount 2": {maxBatchCount: "2", expected: expected{number: "10", logs: []string{
			"number: 1, steps: 3", "number: 2, steps: 1", "number: 3, steps: 7", "number: 4, steps: 2", "number: 5, steps: 5",
			"number: 6, steps: 8", "number: 7, steps: 16", "number: 8, steps: 3", "number: 9, steps: 19", "number: 10, steps: 6",
		}}},
		"start value 10, maxBatchCount 3": {maxBatchCount: "3", expected: expected{number: "10", logs: []string{
			"number: 1, steps: 3", "number: 2, steps: 1", "number: 3, steps: 7", "number: 4, steps: 2", "number: 5, steps: 5",
			"number: 6, steps: 8", "number: 7, steps: 16", "number: 8, steps: 3", "number: 9, steps: 19", "number: 10, steps: 6",
		}}},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&buf, nil))

			actual, err := bruteforce.Run(context.Background(), bruteforce.Request{
				Start: "0", End: testCase.expected.number, BatchSize: testCase.maxBatchCount, Logger: logger, EnableSteps: true,
			})

			require.NoError(t, err)
			assert.Equal(t, testCase.expected.number, actual.Num)
			assert.NotEmpty(t, buf.String())

			assert.Len(t, testCase.expected.logs, strings.Count(buf.String(), "\n"))
			for _, expectedLog := range testCase.expected.logs {
				assert.Contains(t, buf.String(), expectedLog)
			}
		})
	}
}

func TestLargeStepCount(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))
	actual, err := bruteforce.Run(context.Background(), bruteforce.Request{
		Start:       "0",
		End:         "30",
		BatchSize:   "1",
		Logger:      logger,
		EnableSteps: true,
	})
	require.NoError(t, err)

	assert.Equal(t, "30", actual.Num)
	assert.NotEmpty(t, buf.String())
	assert.Contains(t, buf.String(), "number: 27, steps: 111")
}

func BenchmarkWithStepsEnabled(b *testing.B) {
	for i := 1; i < b.N; i++ {
		_, err := bruteforce.Run(context.Background(), bruteforce.Request{
			Start:       strconv.Itoa(i),
			End:         strconv.Itoa(i * 1000),
			BatchSize:   strconv.Itoa(1000 / 10),
			EnableSteps: true,
		})
		assert.NoError(b, err)
	}
}

func BenchmarkWithoutStepsEnabled(b *testing.B) {
	for i := 1; i < b.N; i++ {
		_, err := bruteforce.Run(context.Background(), bruteforce.Request{
			Start:       strconv.Itoa(i),
			End:         strconv.Itoa(i * 1000),
			BatchSize:   strconv.Itoa(1000 / 10),
			EnableSteps: false,
		})

		assert.NoError(b, err)
	}
}

func FuzzBruteforce(f *testing.F) {
	for i := 0; i < 100; i++ {
		f.Add(rand.Intn(100_000)+1, rand.Intn(100_000)+1, rand.Intn(100_000)+1) //nolint:gosec
	}
	f.Fuzz(func(t *testing.T, start, end, batchSize int) {
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 4096)
				n := runtime.Stack(buf, false)
				t.Logf("start: %d, end: %d, batchSize: %d\n", start, end, batchSize)
				t.Logf("Stack trace:\n%s\n", buf[:n])
				t.Errorf("panic: %v\n", r)
			}
		}()

		_, _ = bruteforce.Run(context.Background(), bruteforce.Request{
			Start:       strconv.Itoa(start),
			End:         strconv.Itoa(end),
			BatchSize:   strconv.Itoa(batchSize),
			EnableSteps: false,
		})
	})
}
