package collatz_test

import (
	"math/big"
	"testing"

	"github.com/jfallis/collatz/pkg/collatz"

	"github.com/stretchr/testify/assert"
)

type testValues struct {
	number         *big.Int
	steps          []*big.Int
	totalStepCount int
}

var (
	testVal0 = testValues{big.NewInt(0), []*big.Int{},
		0}
	testVal1 = testValues{big.NewInt(1), []*big.Int{big.NewInt(4), big.NewInt(2), big.NewInt(1)},
		3}
	testVal2 = testValues{big.NewInt(2), []*big.Int{big.NewInt(1)}, 1}
	testVal7 = testValues{big.NewInt(7), []*big.Int{
		big.NewInt(22), big.NewInt(11), big.NewInt(34), big.NewInt(17), big.NewInt(52), big.NewInt(26),
		big.NewInt(13), big.NewInt(40), big.NewInt(20), big.NewInt(10), big.NewInt(5), big.NewInt(16),
		big.NewInt(8), big.NewInt(4), big.NewInt(2), big.NewInt(1)},
		16}
	testVal27 = testValues{big.NewInt(27), []*big.Int{
		big.NewInt(82), big.NewInt(41), big.NewInt(124), big.NewInt(62), big.NewInt(31), big.NewInt(94),
		big.NewInt(47), big.NewInt(142), big.NewInt(71), big.NewInt(214), big.NewInt(107), big.NewInt(322),
		big.NewInt(161), big.NewInt(484), big.NewInt(242), big.NewInt(121), big.NewInt(364), big.NewInt(182),
		big.NewInt(91), big.NewInt(274), big.NewInt(137), big.NewInt(412), big.NewInt(206), big.NewInt(103),
		big.NewInt(310), big.NewInt(155), big.NewInt(466), big.NewInt(233), big.NewInt(700), big.NewInt(350),
		big.NewInt(175), big.NewInt(526), big.NewInt(263), big.NewInt(790), big.NewInt(395), big.NewInt(1186),
		big.NewInt(593), big.NewInt(1780), big.NewInt(890), big.NewInt(445), big.NewInt(1336), big.NewInt(668),
		big.NewInt(334), big.NewInt(167), big.NewInt(502), big.NewInt(251), big.NewInt(754), big.NewInt(377),
		big.NewInt(1132), big.NewInt(566), big.NewInt(283), big.NewInt(850), big.NewInt(425), big.NewInt(1276),
		big.NewInt(638), big.NewInt(319), big.NewInt(958), big.NewInt(479), big.NewInt(1438), big.NewInt(719),
		big.NewInt(2158), big.NewInt(1079), big.NewInt(3238), big.NewInt(1619), big.NewInt(4858), big.NewInt(2429),
		big.NewInt(7288), big.NewInt(3644), big.NewInt(1822), big.NewInt(911), big.NewInt(2734), big.NewInt(1367),
		big.NewInt(4102), big.NewInt(2051), big.NewInt(6154), big.NewInt(3077), big.NewInt(9232), big.NewInt(4616),
		big.NewInt(2308), big.NewInt(1154), big.NewInt(577), big.NewInt(1732), big.NewInt(866), big.NewInt(433),
		big.NewInt(1300), big.NewInt(650), big.NewInt(325), big.NewInt(976), big.NewInt(488), big.NewInt(244),
		big.NewInt(122), big.NewInt(61), big.NewInt(184), big.NewInt(92), big.NewInt(46), big.NewInt(23),
		big.NewInt(70), big.NewInt(35), big.NewInt(106), big.NewInt(53), big.NewInt(160), big.NewInt(80),
		big.NewInt(40), big.NewInt(20), big.NewInt(10), big.NewInt(5), big.NewInt(16), big.NewInt(8),
		big.NewInt(4), big.NewInt(2), big.NewInt(1),
	}, 111}
)

func TestSuccessError(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input    collatz.SuccessError
		expected string
	}{
		"Test with number 5 and steps [1, 2]": {
			input:    collatz.SuccessError{Number: big.NewInt(5), Steps: []*big.Int{big.NewInt(1), big.NewInt(2)}},
			expected: "You found an infinite loop 🎉 number: 5, steps: [+1 +2]",
		},
		"Test with number 10 and steps [1, 2, 3]": {
			input:    collatz.SuccessError{Number: big.NewInt(10), Steps: []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}},
			expected: "You found an infinite loop 🎉 number: 10, steps: [+1 +2 +3]",
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

	actual := collatz.New(testVal0.number)

	assert.Error(t, actual.Calculate())
	assert.Equal(t, testVal0.totalStepCount, len(actual.Steps()))
}

func TestCollatzCalculate(t *testing.T) {
	t.Parallel()

	tests := map[string]testValues{
		"test value 1":  {testVal1.number, testVal1.steps, testVal1.totalStepCount},
		"test value 2":  {testVal2.number, testVal2.steps, testVal2.totalStepCount},
		"test value 7":  {testVal7.number, testVal7.steps, testVal7.totalStepCount},
		"test value 27": {testVal27.number, testVal27.steps, testVal27.totalStepCount},
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
	assert.Equal(t, 111, len(largestStepCount.Steps()))
}
