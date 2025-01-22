package simple

import (
	"github.com/mikenye/geom2d/point"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPolygonArea2XSigned(t *testing.T) {
	tests := map[string]struct {
		points   []point.Point[int]
		expected int
	}{
		"convex triangle": {
			points: []point.Point[int]{
				point.New(0, 0),
				point.New(4, 0),
				point.New(0, 3),
			},
			expected: 12, // 2 * area of the triangle
		},
		"concave quadrilateral": {
			points: []point.Point[int]{
				point.New(0, 0),
				point.New(4, 0),
				point.New(2, 1), // A point that makes it concave
				point.New(0, 3),
			},
			expected: 10, // 2 * signed area
		},
		"regular quadrilateral": {
			points: []point.Point[int]{
				point.New(0, 0),
				point.New(4, 0),
				point.New(4, 3),
				point.New(0, 3),
			},
			expected: 24, // 2 * area of the rectangle
		},
		"degenerate polygon (line segment)": {
			points: []point.Point[int]{
				point.New(0, 0),
				point.New(4, 0),
			},
			expected: 0, // No area
		},
		"single point polygon": {
			points: []point.Point[int]{
				point.New(0, 0),
			},
			expected: 0, // No area
		},
		"empty point slice": {
			points:   []point.Point[int]{},
			expected: 0, // No area
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Call the function under test
			actual := Area2XSigned(tc.points...)

			// Assert the result
			assert.Equal(t, tc.expected, actual)
		})
	}
}
