package collatz_test

import (
	"fmt"
	"testing"

	"github.com/jfallis/collatz/pkg/collatz"

	"github.com/stretchr/testify/assert"
)

type testValues struct {
	number         uint64
	steps          []uint64
	totalStepCount int
}

var (
	testVal0  = testValues{0, []uint64{}, 0}
	testVal1  = testValues{1, []uint64{4, 2, 1}, 3}
	testVal2  = testValues{2, []uint64{1}, 1}
	testVal7  = testValues{7, []uint64{22, 11, 34, 17, 52, 26, 13, 40, 20, 10, 5, 16, 8, 4, 2, 1}, 16}
	testVal27 = testValues{27, []uint64{
		82, 41, 124, 62, 31, 94, 47, 142, 71, 214, 107, 322, 161, 484, 242, 121, 364, 182, 91, 274, 137,
		412, 206, 103, 310, 155, 466, 233, 700, 350, 175, 526, 263, 790, 395, 1186, 593, 1780, 890, 445,
		1336, 668, 334, 167, 502, 251, 754, 377, 1132, 566, 283, 850, 425, 1276, 638, 319, 958, 479, 1438,
		719, 2158, 1079, 3238, 1619, 4858, 2429, 7288, 3644, 1822, 911, 2734, 1367, 4102, 2051, 6154, 3077,
		9232, 4616, 2308, 1154, 577, 1732, 866, 433, 1300, 650, 325, 976, 488, 244, 122, 61, 184, 92, 46,
		23, 70, 35, 106, 53, 160, 80, 40, 20, 10, 5, 16, 8, 4, 2, 1,
	}, 111}
)

func TestCollatzCalculateErrorHandling(t *testing.T) {
	t.Parallel()

	actual := collatz.New(testVal0.number)

	assert.Error(t, actual.Calculate())
	assert.Equal(t, 0, len(actual.Steps()))
}

func TestCollatzCalculate(t *testing.T) {
	t.Parallel()

	tests := []testValues{testVal1, testVal2, testVal7, testVal27}

	for _, test := range tests {
		t.Run(fmt.Sprintf("test input number %d", test.number), func(t *testing.T) {
			t.Parallel()

			actual := collatz.New(test.number)

			assert.NoError(t, actual.Calculate())

			assert.Equal(t, test.steps, actual.Steps())
			assert.Equal(t, test.totalStepCount, len(actual.Steps()))
			assert.False(t, actual.Success())
		})
	}
}

func TestCollatzHighestIterations(t *testing.T) {
	t.Parallel()

	type highestIteration struct {
		index      uint64
		iterations int
	}

	highest := highestIteration{0, 0}

	for i := uint64(4); i <= 250000; i++ {
		c := collatz.New(i)
		assert.NoError(t, c.Calculate())
		if len(c.Steps()) > highest.iterations {
			highest = highestIteration{i, len(c.Steps())}
		}
	}

	assert.Equal(t, uint64(230631), highest.index)
	assert.Equal(t, 442, highest.iterations)
}
