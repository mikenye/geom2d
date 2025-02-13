package simple

import (
	"github.com/mikenye/geom2d/point"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvexHull(t *testing.T) {
	tests := []struct {
		name     string
		points   any // Supports both int and float64 points
		expected any // Expected convex hull points
	}{
		// Integer points test cases
		{
			name: "int: Convex hull with multiple interior points",
			points: []point.Point[int]{
				point.New(1, 4), point.New(4, 13), point.New(8, 17), point.New(18, 20),
				point.New(33, 18), point.New(38, 11), point.New(34, -2), point.New(21, -3),
				point.New(6, -1), point.New(7, 6), point.New(10, 14), point.New(5, 2),
				point.New(16, 0), point.New(12, 12), point.New(23, 16), point.New(14, 6),
				point.New(24, 0), point.New(21, -2), point.New(30, 14), point.New(27, 9),
				point.New(29, 4), point.New(31, -1), point.New(34, 7), point.New(35, 12),
				point.New(26, 2), point.New(25, 1), point.New(27, 3),
			},
			expected: []point.Point[int]{
				point.New(21, -3), point.New(34, -2), point.New(38, 11), point.New(33, 18),
				point.New(18, 20), point.New(8, 17), point.New(4, 13), point.New(1, 4),
				point.New(6, -1),
			},
		},
		{
			name: "int: Simple square convex hull with all points outside",
			points: []point.Point[int]{
				point.New(0, 0), point.New(20, 0), point.New(20, 20), point.New(0, 20),
				point.New(19, 1), point.New(18, 3), point.New(17, 4), point.New(16, 4),
				point.New(15, 3), point.New(14, 2), point.New(13, 2), point.New(12, 2),
			},
			expected: []point.Point[int]{
				point.New(0, 0), point.New(20, 0), point.New(20, 20), point.New(0, 20),
			},
		},

		// Float64 points test cases
		{
			name: "float64: Convex hull with multiple interior points",
			points: []point.Point[float64]{
				point.New(1.0, 4.0), point.New(4.0, 13.0), point.New(8.0, 17.0), point.New(18.0, 20.0),
				point.New(33.0, 18.0), point.New(38.0, 11.0), point.New(34.0, -2.0), point.New(21.0, -3.0),
				point.New(6.0, -1.0), point.New(7.0, 6.0), point.New(10.0, 14.0), point.New(5.0, 2.0),
				point.New(16.0, 0.0), point.New(12.0, 12.0), point.New(23.0, 16.0), point.New(14.0, 6.0),
				point.New(24.0, 0.0), point.New(21.0, -2.0), point.New(30.0, 14.0), point.New(27.0, 9.0),
				point.New(29.0, 4.0), point.New(31.0, -1.0), point.New(34.0, 7.0), point.New(35.0, 12.0),
				point.New(26.0, 2.0), point.New(25.0, 1.0), point.New(27.0, 3.0),
			},
			expected: []point.Point[float64]{
				point.New(21.0, -3.0), point.New(34.0, -2.0), point.New(38.0, 11.0), point.New(33.0, 18.0),
				point.New(18.0, 20.0), point.New(8.0, 17.0), point.New(4.0, 13.0), point.New(1.0, 4.0),
				point.New(6.0, -1.0),
			},
		},
		{
			name: "float64: Simple square convex hull with all points outside",
			points: []point.Point[float64]{
				point.New(0.0, 0.0), point.New(20.0, 0.0), point.New(20.0, 20.0), point.New(0.0, 20.0),
				point.New(19.0, 1.0), point.New(18.0, 3.0), point.New(17.0, 4.0), point.New(16.0, 4.0),
				point.New(15.0, 3.0), point.New(14.0, 2.0), point.New(13.0, 2.0), point.New(12.0, 2.0),
			},
			expected: []point.Point[float64]{
				point.New(0.0, 0.0), point.New(20.0, 0.0), point.New(20.0, 20.0), point.New(0.0, 20.0),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch points := tt.points.(type) {
			case []point.Point[int]:
				expected := tt.expected.([]point.Point[int])
				actual := ConvexHull(points...)
				assert.Equal(t, expected, actual)

			case []point.Point[float64]:
				expected := tt.expected.([]point.Point[float64])
				actual := ConvexHull(points...)
				assert.Equal(t, expected, actual)
			}
		})
	}
}

func TestFindLowestLeftestPoint(t *testing.T) {
	tests := []struct {
		name      string
		points    any // Supports both int and float64 points
		expectedI int // Expected index of the lowest point
		expectedP any // Expected lowest point as either int or float64
	}{
		// Integer points test case
		{
			name: "int: Multiple points, lowest at (0,0)",
			points: []point.Point[int]{
				point.New(10, 10),
				point.New(10, 0),
				point.New(0, 0),
			},
			expectedI: 2,
			expectedP: point.New[int](0, 0),
		},

		// Float64 points test case
		{
			name: "float64: Multiple points, lowest at (0.0,0.0)",
			points: []point.Point[float64]{
				point.New[float64](10.5, 10.5),
				point.New[float64](10.0, 0.0),
				point.New[float64](0.0, 0.0),
			},
			expectedI: 2,
			expectedP: point.New[float64](0.0, 0.0),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			switch points := tc.points.(type) {
			case []point.Point[int]:
				expectedP := tc.expectedP.(point.Point[int])
				actualI, actualP := findLowestLeftestPoint(points...)
				assert.Equal(t, tc.expectedI, actualI)
				assert.Equal(t, expectedP, actualP)

			case []point.Point[float64]:
				expectedP := tc.expectedP.(point.Point[float64])
				actualI, actualP := findLowestLeftestPoint(points...)
				assert.Equal(t, tc.expectedI, actualI)
				assert.Equal(t, expectedP, actualP)
			}
		})
	}
}

func TestOrderPointsByAngleAboutLowestPoint(t *testing.T) {
	tests := map[string]struct {
		points, expected any
	}{
		"int: non-collinear points": {
			points: []point.Point[int]{ // input points are expected points, but have been randomized.
				point.New(-10, 10), point.New(-8, 10), point.New(10, 9), point.New(-10, 0),
				point.New(-2, 10), point.New(10, 4), point.New(-9, 10), point.New(-10, 6),
				point.New(-10, 3), point.New(-10, 1), point.New(-6, 10), point.New(-10, 7),
				point.New(-4, 10), point.New(-10, 4), point.New(6, 10), point.New(10, 5),
				point.New(10, 3), point.New(-5, 10), point.New(10, 1), point.New(-7, 10),
				point.New(3, 10), point.New(-10, 10), point.New(8, 10), point.New(-10, 8),
				point.New(-3, 10), point.New(7, 10), point.New(10, 2), point.New(-10, 9),
				point.New(0, 10), point.New(9, 10), point.New(10, 10), point.New(2, 10),
				point.New(0, -1), point.New(-10, 5), point.New(4, 10), point.New(-10, 2),
				point.New(10, 0), point.New(5, 10), point.New(10, 6), point.New(10, 8),
				point.New(10, 7), point.New(-1, 10), point.New(1, 10),
			},
			expected: []point.Point[int]{
				point.New(0, -1), // lowest point
				point.New(10, 0), point.New(10, 1), point.New(10, 2), point.New(10, 3),
				point.New(10, 4), point.New(10, 5), point.New(10, 6), point.New(10, 7),
				point.New(10, 8), point.New(10, 9), point.New(10, 10), point.New(9, 10),
				point.New(8, 10), point.New(7, 10), point.New(6, 10), point.New(5, 10),
				point.New(4, 10), point.New(3, 10), point.New(2, 10), point.New(1, 10),
				point.New(0, 10), point.New(-1, 10), point.New(-2, 10), point.New(-3, 10),
				point.New(-4, 10), point.New(-5, 10), point.New(-6, 10), point.New(-7, 10),
				point.New(-8, 10), point.New(-9, 10), point.New(-10, 10), point.New(-10, 10),
				point.New(-10, 9), point.New(-10, 8), point.New(-10, 7), point.New(-10, 6),
				point.New(-10, 5), point.New(-10, 4), point.New(-10, 3), point.New(-10, 2),
				point.New(-10, 1), point.New(-10, 0),
			},
		},
		"int: collinear points": {
			points: []point.Point[int]{ // input points are expected points, but have been randomized.
				point.New(1, 8), point.New(-3, 8), point.New(-7, 7), point.New(3, 3),
				point.New(-6, 8), point.New(-3, 3), point.New(7, 8), point.New(7, 7),
				point.New(5, 8), point.New(4, 8), point.New(0, 0), point.New(3, 8),
				point.New(4, 4), point.New(0, 8), point.New(8, 8), point.New(-4, 4),
				point.New(-5, 5), point.New(-5, 8), point.New(6, 6), point.New(5, 5),
				point.New(6, 8), point.New(1, 1), point.New(-2, 8), point.New(2, 8),
				point.New(-7, 8), point.New(2, 2), point.New(-1, 1), point.New(-2, 2),
				point.New(-1, 8), point.New(-8, 8), point.New(-6, 6), point.New(-4, 8),
			},
			expected: []point.Point[int]{
				point.New(0, 0), // lowest point
				point.New(1, 1), point.New(2, 2), point.New(3, 3), point.New(4, 4),
				point.New(5, 5), point.New(6, 6), point.New(7, 7), point.New(8, 8),
				point.New(7, 8), point.New(6, 8), point.New(5, 8), point.New(4, 8),
				point.New(3, 8), point.New(2, 8), point.New(1, 8), point.New(0, 8),
				point.New(-1, 8), point.New(-2, 8), point.New(-3, 8), point.New(-4, 8),
				point.New(-5, 8), point.New(-6, 8), point.New(-7, 8), point.New(-1, 1),
				point.New(-2, 2), point.New(-3, 3), point.New(-4, 4), point.New(-5, 5),
				point.New(-6, 6), point.New(-7, 7), point.New(-8, 8),
			},
		},
		"float64: non-collinear points": {
			points: []point.Point[float64]{ // input points are expected points, but have been randomized.
				point.New(-10.0, 10.0), point.New(-8.0, 10.0), point.New(10.0, 9.0),
				point.New(-10.0, 0.0), point.New(-2.0, 10.0), point.New(10.0, 4.0),
				point.New(-9.0, 10.0), point.New(-10.0, 6.0), point.New(-10.0, 3.0),
				point.New(-10.0, 1.0), point.New(-6.0, 10.0), point.New(-10.0, 7.0),
				point.New(-4.0, 10.0), point.New(-10.0, 4.0), point.New(6.0, 10.0),
				point.New(10.0, 5.0), point.New(10.0, 3.0), point.New(-5.0, 10.0),
				point.New(10.0, 1.0), point.New(-7.0, 10.0), point.New(3.0, 10.0),
				point.New(-10.0, 10.0), point.New(8.0, 10.0), point.New(-10.0, 8.0),
				point.New(-3.0, 10.0), point.New(7.0, 10.0), point.New(10.0, 2.0),
				point.New(-10.0, 9.0), point.New(0.0, 10.0), point.New(9.0, 10.0),
				point.New(10.0, 10.0), point.New(2.0, 10.0), point.New(0.0, -1.0),
				point.New(-10.0, 5.0), point.New(4.0, 10.0), point.New(-10.0, 2.0),
				point.New(10.0, 0.0), point.New(5.0, 10.0), point.New(10.0, 6.0),
				point.New(10.0, 8.0), point.New(10.0, 7.0), point.New(-1.0, 10.0),
				point.New(1.0, 10.0),
			},
			expected: []point.Point[float64]{
				point.New(0.0, -1.0), // lowest point
				point.New(10.0, 0.0), point.New(10.0, 1.0), point.New(10.0, 2.0),
				point.New(10.0, 3.0), point.New(10.0, 4.0), point.New(10.0, 5.0),
				point.New(10.0, 6.0), point.New(10.0, 7.0), point.New(10.0, 8.0),
				point.New(10.0, 9.0), point.New(10.0, 10.0), point.New(9.0, 10.0),
				point.New(8.0, 10.0), point.New(7.0, 10.0), point.New(6.0, 10.0),
				point.New(5.0, 10.0), point.New(4.0, 10.0), point.New(3.0, 10.0),
				point.New(2.0, 10.0), point.New(1.0, 10.0), point.New(0.0, 10.0),
				point.New(-1.0, 10.0), point.New(-2.0, 10.0), point.New(-3.0, 10.0),
				point.New(-4.0, 10.0), point.New(-5.0, 10.0), point.New(-6.0, 10.0),
				point.New(-7.0, 10.0), point.New(-8.0, 10.0), point.New(-9.0, 10.0),
				point.New(-10.0, 10.0), point.New(-10.0, 10.0), point.New(-10.0, 9.0),
				point.New(-10.0, 8.0), point.New(-10.0, 7.0), point.New(-10.0, 6.0),
				point.New(-10.0, 5.0), point.New(-10.0, 4.0), point.New(-10.0, 3.0),
				point.New(-10.0, 2.0), point.New(-10.0, 1.0), point.New(-10.0, 0.0),
			},
		},
		"float64: collinear points": {
			points: []point.Point[float64]{ // input points are expected points, but have been randomized.
				point.New(1.0, 8.0), point.New(-3.0, 8.0), point.New(-7.0, 7.0),
				point.New(3.0, 3.0), point.New(-6.0, 8.0), point.New(-3.0, 3.0),
				point.New(7.0, 8.0), point.New(7.0, 7.0), point.New(5.0, 8.0),
				point.New(4.0, 8.0), point.New(0.0, 0.0), point.New(3.0, 8.0),
				point.New(4.0, 4.0), point.New(0.0, 8.0), point.New(8.0, 8.0),
				point.New(-4.0, 4.0), point.New(-5.0, 5.0), point.New(-5.0, 8.0),
				point.New(6.0, 6.0), point.New(5.0, 5.0), point.New(6.0, 8.0),
				point.New(1.0, 1.0), point.New(-2.0, 8.0), point.New(2.0, 8.0),
				point.New(-7.0, 8.0), point.New(2.0, 2.0), point.New(-1.0, 1.0),
				point.New(-2.0, 2.0), point.New(-1.0, 8.0), point.New(-8.0, 8.0),
				point.New(-6.0, 6.0), point.New(-4.0, 8.0),
			},
			expected: []point.Point[float64]{
				point.New(0.0, 0.0), // lowest point
				point.New(1.0, 1.0), point.New(2.0, 2.0), point.New(3.0, 3.0),
				point.New(4.0, 4.0), point.New(5.0, 5.0), point.New(6.0, 6.0),
				point.New(7.0, 7.0), point.New(8.0, 8.0), point.New(7.0, 8.0),
				point.New(6.0, 8.0), point.New(5.0, 8.0), point.New(4.0, 8.0),
				point.New(3.0, 8.0), point.New(2.0, 8.0), point.New(1.0, 8.0),
				point.New(0.0, 8.0), point.New(-1.0, 8.0), point.New(-2.0, 8.0),
				point.New(-3.0, 8.0), point.New(-4.0, 8.0), point.New(-5.0, 8.0),
				point.New(-6.0, 8.0), point.New(-7.0, 8.0), point.New(-1.0, 1.0),
				point.New(-2.0, 2.0), point.New(-3.0, 3.0), point.New(-4.0, 4.0),
				point.New(-5.0, 5.0), point.New(-6.0, 6.0), point.New(-7.0, 7.0),
				point.New(-8.0, 8.0),
			},
		},
	}
	for name, tc := range tests {

		t.Run(name, func(t *testing.T) {
			switch points := tc.points.(type) {
			case []point.Point[int]:
				expected := tc.expected.([]point.Point[int])
				orderPointsByAngleAboutLowestPoint(expected[0], points)
				assert.Equal(t, tc.expected, tc.points)

			case []point.Point[float64]:
				expected := tc.expected.([]point.Point[float64])
				orderPointsByAngleAboutLowestPoint(expected[0], points)
				assert.Equal(t, tc.expected, tc.points)
			}
		})
	}
}
