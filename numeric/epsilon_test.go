package numeric

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFloatEquals(t *testing.T) {
	tests := map[string]struct {
		a, b, epsilon float64
		expected      bool
	}{
		"exactly equal":                  {a: 1.0, b: 1.0, epsilon: 0.0001, expected: true},
		"within epsilon":                 {a: 1.00001, b: 1.00002, epsilon: 0.0001, expected: true},
		"just outside epsilon":           {a: 1.00001, b: 1.0002, epsilon: 0.0001, expected: false},
		"very different values":          {a: 1.0, b: 2.0, epsilon: 0.0001, expected: false},
		"negative values match":          {a: -1.00001, b: -1.00002, epsilon: 0.0001, expected: true},
		"float imprecision":              {a: 2.759493670886076, b: 2.75949367088608, epsilon: 0, expected: false},
		"float imprecision with epsilon": {a: 2.759493670886076, b: 2.75949367088608, epsilon: 1e-14, expected: true},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expected, FloatEquals(tt.a, tt.b, tt.epsilon))
		})
	}
}

func TestFloatGreaterThan(t *testing.T) {
	tests := map[string]struct {
		a, b, epsilon float64
		expected      bool
	}{
		"greater":            {5.0, 4.9, 0.0001, true},
		"equal within eps":   {5.0, 5.00000001, 0.0001, false},
		"equal exactly":      {5.0, 5.0, 0.0001, false},
		"lesser":             {4.9, 5.0, 0.0001, false},
		"greater at epsilon": {5.0002, 5.0, 0.0001, true},
		"smaller at epsilon": {5.00005, 5.0, 0.0001, false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := FloatGreaterThan(tc.a, tc.b, tc.epsilon)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestSnapToEpsilon(t *testing.T) {
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
