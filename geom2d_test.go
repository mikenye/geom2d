package geom2d

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAbs(t *testing.T) {
	t.Run("positive int", func(t *testing.T) {
		input := 42
		expected := 42
		actual := abs(input)
		assert.Equal(t, expected, actual, "abs(%v) should return %v", input, expected)
	})

	t.Run("negative int", func(t *testing.T) {
		input := -42
		expected := 42
		actual := abs(input)
		assert.Equal(t, expected, actual, "abs(%v) should return %v", input, expected)
	})

	t.Run("zero int", func(t *testing.T) {
		input := 0
		expected := 0
		actual := abs(input)
		assert.Equal(t, expected, actual, "abs(%v) should return %v", input, expected)
	})

	t.Run("positive float", func(t *testing.T) {
		input := 3.14
		expected := 3.14
		actual := abs(input)
		assert.Equal(t, expected, actual, "abs(%v) should return %v", input, expected)
	})

	t.Run("negative float", func(t *testing.T) {
		input := -3.14
		expected := 3.14
		actual := abs(input)
		assert.Equal(t, expected, actual, "abs(%v) should return %v", input, expected)
	})
}

func TestRelationship_String(t *testing.T) {
	tests := map[string]struct {
		input       Relationship
		expected    string
		shouldPanic bool
	}{
		"RelationshipDisjoint": {
			input:       RelationshipDisjoint,
			expected:    "RelationshipDisjoint",
			shouldPanic: false,
		},
		"RelationshipIntersection": {
			input:       RelationshipIntersection,
			expected:    "RelationshipIntersection",
			shouldPanic: false,
		},
		"RelationshipContainedBy": {
			input:       RelationshipContainedBy,
			expected:    "RelationshipContainedBy",
			shouldPanic: false,
		},
		"RelationshipContains": {
			input:       RelationshipContains,
			expected:    "RelationshipContains",
			shouldPanic: false,
		},
		"RelationshipEqual": {
			input:       RelationshipEqual,
			expected:    "RelationshipEqual",
			shouldPanic: false,
		},
		"UnsupportedRelationship": {
			input:       Relationship(100), // An unsupported relationship
			shouldPanic: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.shouldPanic {
				require.Panics(t, func() {
					_ = tt.input.String()
				}, "Expected panic for unsupported relationship")
			} else {
				require.NotPanics(t, func() {
					output := tt.input.String()
					assert.Equal(t, tt.expected, output, "Unexpected string for relationship")
				}, "Did not expect panic for supported relationship")
			}
		})
	}
}

func TestWithEpsilon(t *testing.T) {
	t.Run("Positive epsilon", func(t *testing.T) {
		opts := &geomOptions{}
		option := WithEpsilon(0.001)
		option(opts)

		assert.Equal(t, 0.001, opts.epsilon, "Expected epsilon to be set to 0.001")
	})

	t.Run("Zero epsilon", func(t *testing.T) {
		opts := &geomOptions{}
		option := WithEpsilon(0)
		option(opts)

		assert.Equal(t, 0.0, opts.epsilon, "Expected epsilon to be set to 0.0")
	})

	t.Run("Negative epsilon", func(t *testing.T) {
		opts := &geomOptions{}
		option := WithEpsilon(-0.5)
		option(opts)

		assert.Equal(t, 0.0, opts.epsilon, "Expected epsilon to default to 0.0 for negative input")
	})
}
