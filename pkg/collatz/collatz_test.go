package collatz_test

import (
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"testing"

	"github.com/jfallis/collatz/pkg/collatz"

	"github.com/stretchr/testify/assert"
)

type testValues struct {
	number         *big.Int
	steps          []string
	totalStepCount int
}

var (
	values     []testValues
	valuesOnce sync.Once
)

func collatzValues() []testValues {
	valuesOnce.Do(func() {
		values = []testValues{
			0: {big.NewInt(0), []string{}, 0},
			1: {big.NewInt(1), []string{"4", "2", "1"}, 3},
			2: {big.NewInt(2), []string{"1"}, 1},
			7: {big.NewInt(7), []string{
				"22", "11", "34", "17", "52", "26", "13", "40", "20", "10", "5", "16", "8", "4", "2", "1",
			}, 16},
			27: {big.NewInt(27), []string{
				"82", "41", "124", "62", "31", "94", "47", "142", "71", "214", "107", "322", "161", "484", "242",
				"121", "364", "182", "91", "274", "137", "412", "206", "103", "310", "155", "466", "233", "700", "350",
				"175", "526", "263", "790", "395", "1186", "593", "1780", "890", "445", "1336", "668", "334", "167",
				"502", "251", "754", "377", "1132", "566", "283", "850", "425", "1276", "638", "319", "958", "479",
				"1438", "719", "2158", "1079", "3238", "1619", "4858", "2429", "7288", "3644", "1822", "911", "2734",
				"1367", "4102", "2051", "6154", "3077", "9232", "4616", "2308", "1154", "577", "1732", "866", "433",
				"1300", "650", "325", "976", "488", "244", "122", "61", "184", "92", "46", "23", "70", "35", "106", "53",
				"160", "80", "40", "20", "10", "5", "16", "8", "4", "2", "1",
			}, 111},
		}
	})

	return values
}

func TestSuccessError(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input    collatz.SuccessError
		expected string
	}{
		"Test with number 5 and steps [1, 2]": {
			input:    collatz.SuccessError{Number: big.NewInt(5), Steps: []string{"1", "2"}},
			expected: "You found an infinite loop 🎉 number: 5, steps: [1 2]",
		},
		"Test with number 10 and steps [1, 2, 3]": {
			input:    collatz.SuccessError{Number: big.NewInt(10), Steps: []string{"1", "2", "3"}},
			expected: "You found an infinite loop 🎉 number: 10, steps: [1 2 3]",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual := test.input.Error()
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestCollatzCalculateErrorHandling(t *testing.T) {
	t.Parallel()

	actual := collatz.New(collatzValues()[0].number)

	assert.ErrorIs(t, actual.Calculate(), collatz.ErrInvalidNumber())
	assert.Equal(t, collatzValues()[0].totalStepCount, len(actual.Steps()))
}

func TestCollatzCalculate(t *testing.T) {
	t.Parallel()

	tests := map[string]testValues{
		"test value 1":  {collatzValues()[1].number, collatzValues()[1].steps, collatzValues()[1].totalStepCount},
		"test value 2":  {collatzValues()[2].number, collatzValues()[2].steps, collatzValues()[2].totalStepCount},
		"test value 7":  {collatzValues()[7].number, collatzValues()[7].steps, collatzValues()[7].totalStepCount},
		"test value 27": {collatzValues()[27].number, collatzValues()[27].steps, collatzValues()[27].totalStepCount},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual := collatz.New(test.number)

			assert.NoError(t, actual.Calculate())

			assert.Equal(t, test.steps, actual.Steps())
			assert.Equal(t, test.totalStepCount, len(actual.Steps()))
			assert.False(t, actual.Success())
		})
	}
}

func TestCollatzLargestStepCount(t *testing.T) {
	t.Parallel()

	var largestStepCount *collatz.Collatz

	for i := big.NewInt(4); i.Cmp(big.NewInt(30)) <= 0; i.Add(i, big.NewInt(1)) {
		c := collatz.New(i)
		assert.NoError(t, c.Calculate())
		if largestStepCount == nil || len(largestStepCount.Steps()) < len(c.Steps()) {
			largestStepCount = c
		}
	}

	assert.Equal(t, big.NewInt(27), largestStepCount.Number())
	assert.Len(t, largestStepCount.Steps(), 111)
}

func FuzzCalculate(f *testing.F) {
	f.Skip()
	for i := 0; i < 1000; i++ {
		f.Add(rand.Int63n(100_000) + 1) //nolint:gosec
	}
	expectedValue := "1"
	f.Fuzz(func(t *testing.T, num int64) {
		fmt.Println(num)
		c := collatz.New(big.NewInt(num))
		if err := c.Calculate(); err != nil {
			t.Error(err)
		}

		steps := c.Steps()
		if len(steps) == 0 || steps[len(steps)-1] == expectedValue {
			t.Errorf("Expected last step to be 1, got %v", steps[len(steps)-1])
		}
	})
}
