package collatz_test

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jfallis/collatz/pkg/collatz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	"golang.org/x/sync/errgroup"
)

type expected struct {
	string string
	steps  string
}

type testValues struct {
	number   string
	expected expected
}

func collatzValues() []testValues {
	return []testValues{
		{"1", expected{
			string: "number: 1, steps: 3, max: 4, success: false",
			steps:  "[4 2 1]",
		}},
		{"2", expected{
			string: "number: 2, steps: 1, max: 1, success: false",
			steps:  "[1]",
		}},
		{"7", expected{
			string: "number: 7, steps: 16, max: 52, success: false",
			steps:  "[22 11 34 17 52 26 13 40 20 10 5 16 8 4 2 1]",
		}},
		{"27", expected{
			string: "number: 27, steps: 111, max: 9232, success: false",
			steps: "[82 41 124 62 31 94 47 142 71 214 107 322 161 484 242 " +
				"121 364 182 91 274 137 412 206 103 310 155 466 233 700 350 " +
				"175 526 263 790 395 1186 593 1780 890 445 1336 668 334 167 " +
				"502 251 754 377 1132 566 283 850 425 1276 638 319 958 479 " +
				"1438 719 2158 1079 3238 1619 4858 2429 7288 3644 1822 911 2734 " +
				"1367 4102 2051 6154 3077 9232 4616 2308 1154 577 1732 866 433 " +
				"1300 650 325 976 488 244 122 61 184 92 46 23 70 35 106 53 " +
				"160 80 40 20 10 5 16 8 4 2 1]",
		}},
	}
}

func TestKeyValueString(t *testing.T) {
	t.Parallel()

	testCase := &collatz.KeyValue{Key: 9, Value: "100"}
	actual := testCase.String()
	assert.Equal(t, "key: 9, value: 100", actual)
}

func TestSuccessError(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		input    collatz.SuccessError
		expected string
	}{
		"example success error 1": {
			input:    collatz.NewSuccessErr("example success error 1"),
			expected: "🎉 did you solve the collatz conjecture: example success error 1",
		},
		"example success error 2": {
			input:    collatz.NewSuccessErr("example success error 2"),
			expected: "🎉 did you solve the collatz conjecture: example success error 2",
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual := testCase.input.Error()
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestCollatzCalculateErrorHandling(t *testing.T) {
	t.Parallel()

	actual := collatz.New("0")

	require.ErrorIs(t, actual.Calculate(true), collatz.ErrInvalidNumber)
	assert.Equal(t, "number: 0, steps: 0, max: -1, success: false", actual.String())
}

func TestCollatzCalculateSuccess(t *testing.T) {
	t.Parallel()

	testCases := collatzValues()
	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.number, func(t *testing.T) {
			t.Parallel()

			actual := collatz.New(testCase.number)

			require.NoError(t, actual.Calculate(true))

			assert.Equal(t, testCase.expected.steps, fmt.Sprintf("%+v", actual.Steps()))
			assert.Equal(t, testCase.expected.string, actual.String())
			assert.False(t, actual.Success())
		})
	}
}

func TestCalculateWithTimeoutSuccess(t *testing.T) {
	t.Parallel()

	cal := collatz.New("27")

	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
	defer cancel()
	errGroup, ctx := errgroup.WithContext(ctx)
	errGroup.Go(func() error {
		return cal.CalculateWithContext(ctx, true)
	})
	require.Error(t, errGroup.Wait())
	deadline, ok := ctx.Deadline()
	assert.True(t, ok)
	assert.WithinDuration(t, time.Now(), deadline, time.Second)
	require.ErrorAs(t, cal.Err(), &collatz.SuccessError{Num: "27"})
	assert.True(t, cal.Success())
}

func TestCollatzCalculateLargeSeed(t *testing.T) {
	t.Parallel()

	number := "9" + strings.Repeat("9", 1000)
	t.Logf("Testing large seed: %s", number)
	actual := collatz.New(number)
	require.NoError(t, actual.Calculate(true))
	assert.Len(t, actual.Steps(), 29855)
	assert.False(t, actual.Success())
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func BenchmarkCollatz(b *testing.B) {
	testCases := collatzValues()[3]
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = collatz.New(testCases.number)
	}
}

func FuzzCalculate(f *testing.F) {
	for i := 0; i < 100; i++ {
		f.Add(rand.Intn(100_000) + 1) //nolint:gosec
	}
	f.Fuzz(func(t *testing.T, num int) {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("num: %d\n", num)
				t.Errorf("panic: %v\n", r)
			}
		}()

		c := collatz.New(strconv.Itoa(num))
		if err := c.Calculate(true); err != nil {
			return
		}

		_ = c.Steps()
	})
}
