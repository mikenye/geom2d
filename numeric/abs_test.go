package numeric

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAbs(t *testing.T) {
	tests := map[string]struct {
		input    any // Supports different SignedNumber types
		expected any // Expected absolute value
	}{
		// Test cases for int
		"int: positive number": {
			input:    42,
			expected: 42,
		},
		"int: negative number": {
			input:    -42,
			expected: 42,
		},
		"int: zero": {
			input:    0,
			expected: 0,
		},

		// Test cases for int64
		"int64: positive number": {
			input:    int64(1000000),
			expected: int64(1000000),
		},
		"int64: negative number": {
			input:    int64(-1000000),
			expected: int64(1000000),
		},

		// Test cases for float64
		"float64: positive number": {
			input:    42.42,
			expected: 42.42,
		},
		"float64: negative number": {
			input:    -42.42,
			expected: 42.42,
		},
		"float64: zero": {
			input:    0.0,
			expected: 0.0,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			switch input := tt.input.(type) {
			case int:
				expected := tt.expected.(int)
				assert.Equal(t, expected, Abs(input))
			case int64:
				expected := tt.expected.(int64)
				assert.Equal(t, expected, Abs(input))
			case float64:
				expected := tt.expected.(float64)
				assert.Equal(t, expected, Abs(input))
			}
		})
	}
}
