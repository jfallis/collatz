package collatz

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuccess(t *testing.T) {
	t.Parallel()

	testCases := []int64{-100, -1, 0, 2, 100}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("test success returns true for %d", tc), func(t *testing.T) {
			t.Parallel()

			actual := New("1")
			assert.NoError(t, actual.Calculate(true))
			actual.steps = []*big.Int{big.NewInt(tc)}
			assert.Equal(t, 1, len(actual.Steps()))
			assert.True(t, actual.Success())
		})
	}
}
