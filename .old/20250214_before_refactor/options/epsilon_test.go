package options

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWithEpsilon(t *testing.T) {
	tests := map[string]struct {
		defaultOptions  GeometryOptions
		inputEpsilon    float64
		expectedEpsilon float64
	}{
		"Negative epsilon value (should clamp to zero)": {
			defaultOptions:  GeometryOptions{Epsilon: 0.01},
			inputEpsilon:    -1e-9,
			expectedEpsilon: 0,
		},
		"Zero epsilon value": {
			defaultOptions:  GeometryOptions{Epsilon: 0.01},
			inputEpsilon:    0,
			expectedEpsilon: 0,
		},
		"Positive epsilon value": {
			defaultOptions:  GeometryOptions{Epsilon: 0.01},
			inputEpsilon:    1e-9,
			expectedEpsilon: 1e-9,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			opts := ApplyGeometryOptions(tc.defaultOptions, WithEpsilon(tc.inputEpsilon))
			assert.Equal(t, tc.expectedEpsilon, opts.Epsilon)
		})
	}
}
