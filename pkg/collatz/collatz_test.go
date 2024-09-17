package collatz_test

import (
	"context"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/jfallis/collatz/pkg/collatz"

	"github.com/stretchr/testify/assert"
)

type expected struct {
	string string
	steps  string
}

type testValues struct {
	number string
	expected
}

var (
	values     []testValues
	valuesOnce sync.Once
)

func collatzValues() []testValues {
	valuesOnce.Do(func() {
		values = []testValues{
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
	})

	return values
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

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual := tc.input.Error()
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestCollatzCalculateErrorHandling(t *testing.T) {
	t.Parallel()

	actual := collatz.New("0")

	assert.ErrorIs(t, actual.Calculate(true), collatz.ErrInvalidNumber())
	assert.Equal(t, "number: 0, steps: 0, max: -1, success: false", actual.String())
}

func TestCollatzCalculateSuccess(t *testing.T) {
	t.Parallel()

	testCases := collatzValues()
	for _, tc := range testCases {
		t.Run(tc.number, func(t *testing.T) {
			t.Parallel()

			actual := collatz.New(tc.number)

			assert.NoError(t, actual.Calculate(true))

			assert.Equal(t, tc.expected.steps, actual.Steps().String())
			assert.Equal(t, tc.expected.string, actual.String())
			assert.False(t, actual.Success())
		})
	}
}

func TestCalculateWithTimeoutSuccess(t *testing.T) {
	t.Parallel()

	c := collatz.New("27")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	errGroup, ctx := errgroup.WithContext(ctx)

	errGroup.Go(func() error {
		return c.CalculateWithContext(ctx, true)
	})
	if err := errGroup.Wait(); err != nil {
		assert.ErrorAs(t, err, &context.DeadlineExceeded)
		assert.ErrorAs(t, c.Err(), &collatz.SuccessError{Num: "27"})
		assert.True(t, c.Success())
		return
	}

	t.Error("expected error")
}

func TestCollatzCalculateLargeSeed(t *testing.T) {
	t.Parallel()

	number := "9" + strings.Repeat("9", 1000)
	t.Logf("Testing large seed: %s", number)
	actual := collatz.New(number)
	assert.NoError(t, actual.Calculate(true))
	assert.Equal(t, 29855, len(actual.Steps()))
	assert.False(t, actual.Success())
}

func TestMax(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		input    []*big.Int
		expected string
	}{
		"empty slice": {
			input:    []*big.Int{},
			expected: "-1",
		},
		"single element": {
			input:    []*big.Int{big.NewInt(5)},
			expected: "5",
		},
		"multiple elements": {
			input:    []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3), big.NewInt(10), big.NewInt(5)},
			expected: "10",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, actual := collatz.Max(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func FuzzCalculate(f *testing.F) {
	for i := 0; i < 100; i++ {
		f.Add(rand.Intn(100_000) + 1) //nolint:gosec
	}
	expectedValue := "1"
	f.Fuzz(func(t *testing.T, num int) {
		c := collatz.New(strconv.Itoa(num))
		if err := c.Calculate(true); err != nil {
			t.Error(err)
		}

		steps := c.Steps()
		if len(steps) == 0 || steps[len(steps)-1].String() != expectedValue {
			t.Errorf("Expected last step to be 1, got %v", steps[len(steps)-1])
		}
	})
}
