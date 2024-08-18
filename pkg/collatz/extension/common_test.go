package extension_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/jfallis/collatz/pkg/collatz/extension"

	"github.com/jfallis/collatz/pkg/collatz"

	"github.com/stretchr/testify/assert"
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
			inputErr:    collatz.SuccessError{String: "success error message"},
			expectedErr: "successfully found the number You found an infinite loop 🎉 - success error message",
		},
		"wrapped success error": {
			inputErr: fmt.Errorf("[wrapped error message] %w",
				collatz.SuccessError{String: "success error message"},
			),
			expectedErr: "successfully found the number [wrapped error message] You found an infinite loop 🎉 - success error message",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := extension.WaitErrHandling(tc.inputErr)
			if tc.expectedErr == "" {
				assert.NoError(t, err)
				return
			}

			assert.Error(t, err)
			assert.Equal(t, tc.expectedErr, err.Error())
		})
	}
}

func TestCPUBatchSize(t *testing.T) {
	t.Parallel()

	cpu, err := strconv.Atoi(extension.CPUBatchSize())
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, cpu, 100)
}