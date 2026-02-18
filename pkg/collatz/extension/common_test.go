package extension_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/jfallis/collatz/pkg/collatz"
	"github.com/jfallis/collatz/pkg/collatz/extension"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWaitErrHandling(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		inputErr    error
		expectedErr string
	}{
		"nil error": {
			inputErr:    nil,
			expectedErr: "",
		},
		"non-success error": {
			inputErr:    errors.New("some error"),
			expectedErr: "routine failed: some error",
		},
		"success error": {
			inputErr:    collatz.NewSuccessErr("success error message"),
			expectedErr: "routine failed: 🎉 did you solve the collatz conjecture: success error message",
		},
		"wrapped success error": {
			inputErr: fmt.Errorf("[wrapped error message] %w",
				collatz.NewSuccessErr("success error message"),
			),
			expectedErr: "routine failed: [wrapped error message] 🎉 did you solve the collatz conjecture: success error message",
		},
	}

	for name, testcase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := extension.WaitErrHandling(testcase.inputErr)
			if testcase.expectedErr == "" {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)
			assert.Equal(t, testcase.expectedErr, err.Error())
		})
	}
}

func TestCPUBatchSize(t *testing.T) {
	t.Parallel()

	cpu, err := strconv.Atoi(extension.CPUBatchSize())
	require.NoError(t, err)
	assert.GreaterOrEqual(t, cpu, 100)
}
