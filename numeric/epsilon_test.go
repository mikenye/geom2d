package numeric

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoundToEpsilon(t *testing.T) {
	tests := map[string]struct {
		value    float64
		epsilon  float64
		expected float64
	}{
		"close to whole number":   {value: -0.9999999999, epsilon: 1e-9, expected: -1.0},
		"far from whole number":   {value: 1.0001, epsilon: 1e-9, expected: 1.0001},
		"exactly at whole number": {value: 2.0, epsilon: 1e-9, expected: 2.0},
		"just within epsilon":     {value: 1.9999, epsilon: 1e-3, expected: 2.0},
		"just outside epsilon":    {value: 1.9999, epsilon: 1e-5, expected: 1.9999},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, SnapToEpsilon(tc.value, tc.epsilon))
		})
	}
}
