package linesegment

import (
	"github.com/mikenye/geom2d/options"
	"github.com/mikenye/geom2d/point"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestLineSegment_AsFloat32(t *testing.T) {
	tests := map[string]struct {
		start, end any
		expected   LineSegment[float32]
	}{
		"int: start: (1,2), end: (3,4)": {
			start: point.New[int](1, 2),
			end:   point.New[int](3, 4),
			expected: NewFromPoints[float32](
				point.New[float32](1, 2),
				point.New[float32](3, 4),
			),
		},
		"float64: start: (1,2), end: (3,4)": {
			start: point.New[float64](1, 2),
			end:   point.New[float64](3, 4),
			expected: NewFromPoints[float32](
				point.New[float32](1, 2),
				point.New[float32](3, 4),
			),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			switch start := tt.start.(type) {
			case point.Point[int]:
				end := tt.end.(point.Point[int])
				ls := NewFromPoints(start, end)
				assert.Equal(t, tt.expected, ls.AsFloat32())
			case point.Point[float64]:
				end := tt.end.(point.Point[float64])
				ls := NewFromPoints(start, end)
				assert.Equal(t, tt.expected, ls.AsFloat32())
			}
		})
	}
}

func TestLineSegment_AsFloat64(t *testing.T) {
	tests := map[string]struct {
		start, end any
		expected   LineSegment[float64]
	}{
		"int: start: (1,2), end: (3,4)": {
			start: point.New[int](1, 2),
			end:   point.New[int](3, 4),
			expected: NewFromPoints[float64](
				point.New[float64](1, 2),
				point.New[float64](3, 4),
			),
		},
		"float64: start: (1,2), end: (3,4)": {
			start: point.New[float64](1, 2),
			end:   point.New[float64](3, 4),
			expected: NewFromPoints[float64](
				point.New[float64](1, 2),
				point.New[float64](3, 4),
			),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			switch start := tt.start.(type) {
			case point.Point[int]:
				end := tt.end.(point.Point[int])
				ls := NewFromPoints(start, end)
				assert.Equal(t, tt.expected, ls.AsFloat64())
			case point.Point[float64]:
				end := tt.end.(point.Point[float64])
				ls := NewFromPoints(start, end)
				assert.Equal(t, tt.expected, ls.AsFloat64())
			}
		})
	}
}

func TestLineSegment_AsInt(t *testing.T) {
	tests := map[string]struct {
		start, end any
		expected   LineSegment[int]
	}{
		"float64: start: (1.2,2.7), end: (3.2,4.7)": {
			start: point.New[float64](1.2, 2.7),
			end:   point.New[float64](3.2, 4.7),
			expected: NewFromPoints[int](
				point.New[int](1, 2),
				point.New[int](3, 4),
			),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			switch start := tt.start.(type) {
			case point.Point[int]:
				end := tt.end.(point.Point[int])
				ls := NewFromPoints(start, end)
				assert.Equal(t, tt.expected, ls.AsInt())
			case point.Point[float64]:
				end := tt.end.(point.Point[float64])
				ls := NewFromPoints(start, end)
				assert.Equal(t, tt.expected, ls.AsInt())
			}
		})
	}
}

func TestLineSegment_AsIntRounded(t *testing.T) {
	tests := map[string]struct {
		start, end any
		expected   LineSegment[int]
	}{
		"float64: start: (1.2,2.7), end: (3.2,4.7)": {
			start: point.New[float64](1.2, 2.7),
			end:   point.New[float64](3.2, 4.7),
			expected: NewFromPoints[int](
				point.New[int](1, 3),
				point.New[int](3, 5),
			),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			switch start := tt.start.(type) {
			case point.Point[int]:
				end := tt.end.(point.Point[int])
				ls := NewFromPoints(start, end)
				assert.Equal(t, tt.expected, ls.AsIntRounded())
			case point.Point[float64]:
				end := tt.end.(point.Point[float64])
				ls := NewFromPoints(start, end)
				assert.Equal(t, tt.expected, ls.AsIntRounded())
			}
		})
	}
}

func TestLineSegment_Bresenham(t *testing.T) {
	tests := map[string]struct {
		lineSegment LineSegment[int]
		expected    []point.Point[int]
	}{
		"horizontal line": {
			lineSegment: NewFromPoints(point.New(0, 0), point.New(5, 0)),
			expected: []point.Point[int]{
				point.New(0, 0),
				point.New(1, 0),
				point.New(2, 0),
				point.New(3, 0),
				point.New(4, 0),
				point.New(5, 0),
			},
		},
		"vertical line": {
			lineSegment: NewFromPoints(point.New(0, 0), point.New(0, 5)),
			expected: []point.Point[int]{
				point.New(0, 0),
				point.New(0, 1),
				point.New(0, 2),
				point.New(0, 3),
				point.New(0, 4),
				point.New(0, 5),
			},
		},
		"diagonal line": {
			lineSegment: NewFromPoints(point.New(0, 0), point.New(5, 5)),
			expected: []point.Point[int]{
				point.New(0, 0),
				point.New(1, 1),
				point.New(2, 2),
				point.New(3, 3),
				point.New(4, 4),
				point.New(5, 5),
			},
		},
		"reverse diagonal": {
			lineSegment: NewFromPoints(point.New(5, 5), point.New(0, 0)),
			expected: []point.Point[int]{
				point.New(5, 5),
				point.New(4, 4),
				point.New(3, 3),
				point.New(2, 2),
				point.New(1, 1),
				point.New(0, 0),
			},
		},
		"steep slope": {
			lineSegment: NewFromPoints(point.New(0, 0), point.New(2, 5)),
			expected: []point.Point[int]{
				point.New(0, 0),
				point.New(0, 1),
				point.New(1, 2),
				point.New(1, 3),
				point.New(2, 4),
				point.New(2, 5),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var actual []point.Point[int]
			test.lineSegment.Bresenham(func(p point.Point[int]) bool {
				actual = append(actual, p)
				return true
			})
			assert.ElementsMatch(t, test.expected, actual, "Bresenham points mismatch")
		})
	}
}

func TestLineSegment_Center(t *testing.T) {
	tests := map[string]struct {
		lineSegment LineSegment[int]
		epsilon     float64
		expected    point.Point[float64]
	}{
		"No epsilon, simple case": {
			lineSegment: New(0, 0, 4, 4),
			epsilon:     0,
			expected:    point.New[float64](2, 2),
		},
		"No epsilon, negative coordinates": {
			lineSegment: New(-4, -4, 4, 4),
			epsilon:     0,
			expected:    point.New[float64](0, 0),
		},
		"With epsilon, rounding applied": {
			lineSegment: New(0, 0, 3, 3),
			epsilon:     0.1,
			expected:    point.New[float64](1.5, 1.5), // No rounding as it's precise
		},
		"With epsilon, midpoint near integer": {
			lineSegment: New(0, 0, 4, 5),
			epsilon:     0.5,
			expected:    point.New[float64](2, 2.5), // Epsilon not applied due to midpoint already exact
		},
		"With epsilon, midpoint adjusted to integer": {
			lineSegment: New(0, 0, 5, 5),
			epsilon:     0.5,
			expected:    point.New[float64](2.5, 2.5), // Exact match without adjustment
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Calculate center
			center := tc.lineSegment.Center(options.WithEpsilon(tc.epsilon))

			// Assert the result
			assert.InDelta(t, tc.expected.X(), center.X(), 1e-9, "Unexpected x-coordinate for center")
			assert.InDelta(t, tc.expected.Y(), center.Y(), 1e-9, "Unexpected y-coordinate for center")
		})
	}
}

func TestLineSegment_ContainsPoint(t *testing.T) {
	tests := map[string]struct {
		segment  LineSegment[float64]
		point    point.Point[float64]
		opts     []options.GeometryOptionsFunc
		expected bool
	}{
		"point lies on the segment": {
			segment:  New(0.0, 0.0, 10.0, 10.0),
			point:    point.New(5.0, 5.0),
			opts:     []options.GeometryOptionsFunc{options.WithEpsilon(1e-7)},
			expected: true,
		},
		"point is off the segment but collinear": {
			segment:  New(0.0, 0.0, 10.0, 10.0),
			point:    point.New(15.0, 15.0),
			opts:     []options.GeometryOptionsFunc{options.WithEpsilon(1e-7)},
			expected: false,
		},
		"point lies at the start of the segment": {
			segment:  New(0.0, 0.0, 10.0, 10.0),
			point:    point.New(0.0, 0.0),
			opts:     []options.GeometryOptionsFunc{options.WithEpsilon(1e-7)},
			expected: true,
		},
		"point lies at the end of the segment": {
			segment:  New(0.0, 0.0, 10.0, 10.0),
			point:    point.New(10.0, 10.0),
			opts:     []options.GeometryOptionsFunc{options.WithEpsilon(1e-7)},
			expected: true,
		},
		"point is slightly off due to precision issues": {
			segment:  New(0.0, 0.0, 10.0, 10.0),
			point:    point.New(5.0000001, 5.0),
			opts:     []options.GeometryOptionsFunc{options.WithEpsilon(1e-7)},
			expected: true,
		},
		"point is not on the segment": {
			segment:  New(0.0, 0.0, 10.0, 10.0),
			point:    point.New(5.0, 6.0),
			opts:     []options.GeometryOptionsFunc{options.WithEpsilon(1e-7)},
			expected: false,
		},
		"point lies on degenerate segment": {
			segment:  New(5.0, 5.0, 5.0, 5.0),
			point:    point.New(5.0, 5.0),
			opts:     []options.GeometryOptionsFunc{options.WithEpsilon(1e-7)},
			expected: true,
		},
		"point does not lie on degenerate segment": {
			segment:  New(5.0, 5.0, 5.0, 5.0),
			point:    point.New(6.0, 5.0),
			opts:     []options.GeometryOptionsFunc{options.WithEpsilon(1e-7)},
			expected: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := tt.segment.ContainsPoint(tt.point, tt.opts...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLineSegment_DistanceToLineSegment(t *testing.T) {
	tests := map[string]struct {
		segA, segB LineSegment[float64]
		expected   float64
		epsilon    float64
		expectZero bool
	}{
		"intersecting segments": {
			segA:       New[float64](0, 0, 4, 4),
			segB:       New[float64](0, 4, 4, 0),
			expected:   0,
			epsilon:    1e-9,
			expectZero: true,
		},
		"parallel non-intersecting segments": {
			segA:       New[float64](0, 0, 4, 0),
			segB:       New[float64](0, 2, 4, 2),
			expected:   2,
			epsilon:    1e-9,
			expectZero: false,
		},
		"segments touching at one endpoint": {
			segA:       New[float64](0, 0, 4, 0),
			segB:       New[float64](4, 0, 4, 4),
			expected:   0,
			epsilon:    1e-9,
			expectZero: true,
		},
		"skew non-intersecting segments": {
			segA:       New[float64](0, 0, 1, 1),
			segB:       New[float64](2, 0, 3, 1),
			expected:   1.4142,
			epsilon:    1e-4,
			expectZero: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			opts := options.WithEpsilon(tt.epsilon)
			actual := tt.segA.DistanceToLineSegment(tt.segB, opts)
			if tt.expectZero {
				assert.InDelta(t, 0, actual, tt.epsilon, "Expected distance to be zero")
			} else {
				assert.InDelta(t, tt.expected, actual, tt.epsilon, "Expected distances to match")
			}
		})
	}
}

func TestLineSegment_DistanceToPoint(t *testing.T) {
	tests := map[string]struct {
		point    any     // Point to be projected (can be int or float64)
		segment  any     // Line segment for projection (can be int or float64)
		expected float64 // Expected distance
	}{
		// Integer points test cases
		"Project onto segment from inside (int)": {
			point:    point.New(5, 5),
			segment:  New[int](2, 3, 8, 7),
			expected: 0.0, // Point is on the segment
		},
		"Project onto segment from above (int)": {
			point:    point.New(4, 6),
			segment:  New[int](2, 3, 8, 7),
			expected: 1.386,
		},
		"Project onto segment from below (int)": {
			point:    point.New(4, 2),
			segment:  New[int](2, 3, 8, 7),
			expected: 1.941,
		},
		"Project off the start of segment (int)": {
			point:    point.New(0, 5),
			segment:  New[int](2, 3, 8, 7),
			expected: 2.8284,
		},
		"Project off the end of segment (int)": {
			point:    point.New(10, 5),
			segment:  New[int](2, 3, 8, 7),
			expected: 2.8284,
		},

		// Float64 points test cases
		"Project onto segment from inside (float64)": {
			point:    point.New(5.5, 5.5),
			segment:  New[float64](2.0, 3.0, 8.0, 7.0),
			expected: 0.1387,
		},
		"Project onto segment from above (float64)": {
			point:    point.New(4.0, 6.0),
			segment:  New[float64](2.0, 3.0, 8.0, 7.0),
			expected: 1.386,
		},
		"Project onto segment from below (float64)": {
			point:    point.New(4.0, 2.0),
			segment:  New[float64](2.0, 3.0, 8.0, 7.0),
			expected: 1.941,
		},
		"Project off the start of segment (float64)": {
			point:    point.New(0.0, 5.0),
			segment:  New[float64](2.0, 3.0, 8.0, 7.0),
			expected: 2.8284,
		},
		"Project off the end of segment (float64)": {
			point:    point.New(10.0, 5.0),
			segment:  New[float64](2.0, 3.0, 8.0, 7.0),
			expected: 2.8284,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			switch p := tt.point.(type) {
			case point.Point[int]:
				seg := tt.segment.(LineSegment[int])
				actual := seg.DistanceToPoint(p)
				assert.InDelta(t, tt.expected, actual, 0.001, "Expected distance does not match actual distance for int points")

			case point.Point[float64]:
				seg := tt.segment.(LineSegment[float64])
				actual := seg.DistanceToPoint(p)
				assert.InDelta(t, tt.expected, actual, 0.001, "Expected distance does not match actual distance for float64 points")
			}
		})
	}
}

func TestLineSegment_Eq(t *testing.T) {
	tests := map[string]struct {
		segment1 any  // First line segment (can be int or float64)
		segment2 any  // Second line segment to compare (can be int or float64)
		expected bool // Expected result of equality check
	}{
		// Integer segment test cases
		"int: Equal segments": {
			segment1: NewFromPoints[int](point.New(1, 1), point.New(4, 5)),
			segment2: NewFromPoints[int](point.New(1, 1), point.New(4, 5)),
			expected: true,
		},
		"int: Unequal segments": {
			segment1: NewFromPoints[int](point.New(1, 1), point.New(4, 5)),
			segment2: NewFromPoints[int](point.New(2, 2), point.New(3, 3)),
			expected: false,
		},

		// Float64 segment test cases
		"float64: Equal segments": {
			segment1: NewFromPoints[float64](point.New(1.0, 1.0), point.New(4.0, 5.0)),
			segment2: NewFromPoints[float64](point.New(1.0, 1.0), point.New(4.0, 5.0)),
			expected: true,
		},
		"float64: Unequal segments": {
			segment1: NewFromPoints[float64](point.New(1.5, 1.5), point.New(3.5, 4.5)),
			segment2: NewFromPoints[float64](point.New(1.5, 1.5), point.New(5.5, 6.5)),
			expected: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			switch segment1 := tc.segment1.(type) {
			case LineSegment[int]:
				segment2 := tc.segment2.(LineSegment[int])
				result := segment1.Eq(segment2)
				assert.Equal(t, tc.expected, result)

			case LineSegment[float64]:
				segment2 := tc.segment2.(LineSegment[float64])
				result := segment1.Eq(segment2)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestLineSegment_Length(t *testing.T) {
	tests := map[string]struct {
		lineSegment LineSegment[float64]
		expected    float64
		epsilon     float64
	}{
		"horizontal segment": {
			lineSegment: New[float64](0, 0, 5, 0),
			expected:    5.0,
			epsilon:     1e-9,
		},
		"vertical segment": {
			lineSegment: New[float64](0, 0, 0, 7),
			expected:    7.0,
			epsilon:     1e-9,
		},
		"diagonal segment": {
			lineSegment: New[float64](0, 0, 3, 4),
			expected:    5.0,
			epsilon:     1e-9,
		},
		"near-zero length segment": {
			lineSegment: New[float64](1e-10, 1e-10, 2e-10, 2e-10),
			expected:    1.4142e-10,
			epsilon:     1e-12,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.lineSegment.Length(options.WithEpsilon(tt.epsilon))
			assert.InDelta(t, tt.expected, actual, tt.epsilon, "Expected length to match")
		})
	}
}

func TestLineSegment_Points(t *testing.T) {
	tests := map[string]struct {
		segment  LineSegment[int]
		expected []point.Point[int]
	}{
		"horizontal segment": {
			segment: NewFromPoints(point.New(1, 1), point.New(5, 1)),
			expected: []point.Point[int]{
				point.New(1, 1),
				point.New(5, 1),
			},
		},
		"vertical segment": {
			segment: NewFromPoints(point.New(3, 2), point.New(3, 6)),
			expected: []point.Point[int]{
				point.New(3, 2),
				point.New(3, 6),
			},
		},
		"diagonal segment": {
			segment: NewFromPoints(point.New(0, 0), point.New(3, 4)),
			expected: []point.Point[int]{
				point.New(0, 0),
				point.New(3, 4),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actualStart, actualEnd := tc.segment.Points()
			assert.Equal(t, tc.expected[0], actualStart)
			assert.Equal(t, tc.expected[1], actualEnd)
		})
	}
}

func TestLineSegment_ProjectPoint(t *testing.T) {
	tests := map[string]struct {
		point    any                  // Point to be projected (can be int or float64)
		segment  any                  // Line segment for projection (can be int or float64)
		expected point.Point[float64] // Expected projected point (float64 type)
	}{
		// Integer points test cases
		"int: Project onto segment from inside": {
			point:    point.New(5, 5),
			segment:  New[int](2, 3, 8, 7),
			expected: point.New[float64](5, 5), // Should project onto the line segment itself
		},
		"int: Project onto segment from above": {
			point:    point.New(4, 6),
			segment:  New[int](2, 3, 8, 7),
			expected: point.New[float64](4.769, 4.846),
		},
		"int: Project onto segment from below": {
			point:    point.New(4, 2),
			segment:  New[int](2, 3, 8, 7),
			expected: point.New[float64](2.9231, 3.6154),
		},
		"int: Project off the start of segment": {
			point:    point.New(0, 5),
			segment:  New[int](2, 3, 8, 7),
			expected: point.New[float64](2, 3), // Should return a point of the segment
		},
		"int: Project off the end of segment": {
			point:    point.New(10, 5),
			segment:  New[int](2, 3, 8, 7),
			expected: point.New[float64](8, 7), // Should return end point of the segment
		},

		// Float64 points test cases
		"float64: Project onto segment from inside": {
			point:    point.New(5.5, 5.5),
			segment:  New[float64](2.0, 3.0, 8.0, 7.0),
			expected: point.New[float64](5.5772, 5.3848), // Should project onto the line segment itself
		},
		"float64: Project off the start of segment": {
			point:    point.New(0.0, 5.0),
			segment:  New[float64](2.0, 3.0, 8.0, 7.0),
			expected: point.New[float64](2.0, 3.0), // Should return a point of the segment
		},
		"float64: Project off the end of segment": {
			point:    point.New(10.0, 5.0),
			segment:  New[float64](2.0, 3.0, 8.0, 7.0),
			expected: point.New[float64](8.0, 7.0), // Should return end point of the segment
		},

		// Zero-length segment test cases
		"int: Project onto zero-length segment": {
			point:    point.New(3, 4),
			segment:  New[int](2, 2, 2, 2),     // Zero-length segment
			expected: point.New[float64](2, 2), // Should return point A (or End), since the segment is a single point
		},
		"float64: Project onto zero-length segment": {
			point:    point.New(5.0, 5.0),
			segment:  New[float64](2.5, 2.5, 2.5, 2.5), // Zero-length segment
			expected: point.New[float64](2.5, 2.5),     // Should return point A (or End), since the segment is a single point
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			switch p := tt.point.(type) {
			case point.Point[int]:
				seg := tt.segment.(LineSegment[int])
				actual := seg.ProjectPoint(p)
				assert.InDelta(t, tt.expected.X(), actual.X(), 0.001)
				assert.InDelta(t, tt.expected.Y(), actual.Y(), 0.001)

			case point.Point[float64]:
				seg := tt.segment.(LineSegment[float64])
				actual := seg.ProjectPoint(p)
				assert.InDelta(t, tt.expected.X(), actual.X(), 0.001)
				assert.InDelta(t, tt.expected.Y(), actual.Y(), 0.001)
			}
		})
	}
}

func TestLineSegment_ReflectLineSegment(t *testing.T) {
	tests := map[string]struct {
		lineSegment LineSegment[float64]
		other       LineSegment[float64]
		expected    LineSegment[float64]
		epsilon     float64
	}{
		"horizontal reflection": {
			lineSegment: New[float64](0, 0, 10, 0),
			other:       New[float64](2, 2, 8, 2),
			expected:    New[float64](2, -2, 8, -2),
			epsilon:     1e-9,
		},
		"vertical reflection": {
			lineSegment: New[float64](0, 0, 0, 10),
			other:       New[float64](2, 2, 2, 8),
			expected:    New[float64](-2, 2, -2, 8),
			epsilon:     1e-9,
		},
		"diagonal reflection": {
			lineSegment: New[float64](0, 0, 10, 10),
			other:       New[float64](2, 6, 6, 2),
			expected:    New[float64](6, 2, 2, 6),
			epsilon:     1e-9,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.lineSegment.ReflectLineSegment(tt.other)
			assert.InDelta(t, tt.expected.Start().X(), actual.Start().X(), tt.epsilon, "Start X mismatch")
			assert.InDelta(t, tt.expected.Start().Y(), actual.Start().Y(), tt.epsilon, "Start Y mismatch")
			assert.InDelta(t, tt.expected.End().X(), actual.End().X(), tt.epsilon, "End X mismatch")
			assert.InDelta(t, tt.expected.End().Y(), actual.End().Y(), tt.epsilon, "End Y mismatch")
		})
	}
}

func TestLineSegment_ReflectPoint(t *testing.T) {
	tests := map[string]struct {
		point    point.Point[float64] // The point to reflect
		axis     LineSegment[float64] // Axis for reflection
		expected point.Point[float64] // Expected reflected point
	}{
		"reflect across x-axis": {
			point:    point.New[float64](3, 4),
			axis:     New[float64](0, 0, 1, 0),
			expected: point.New[float64](3, -4),
		},
		"reflect across y-axis": {
			point:    point.New[float64](3, 4),
			axis:     New[float64](0, 0, 0, 1),
			expected: point.New[float64](-3, 4),
		},
		"reflect across y = x line (ReflectAcrossCustomLine)": {
			point:    point.New[float64](3, 4),
			axis:     New[float64](0, 0, 1, 1),
			expected: point.New[float64](4, 3),
		},
		"reflect across y = -x line (ReflectAcrossCustomLine)": {
			point:    point.New[float64](3, 4),
			axis:     New[float64](0, 0, -1, 1),
			expected: point.New[float64](-4, -3),
		},
		"reflect across degenerate line segment": {
			point:    point.New[float64](3, 4),
			axis:     New[float64](1, 1, 1, 1), // Degenerate line
			expected: point.New[float64](3, 4), // Expect the point to remain unchanged
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.axis.ReflectPoint(tc.point)
			assert.InDelta(t, tc.expected.X(), result.X(), 0.001)
			assert.InDelta(t, tc.expected.Y(), result.Y(), 0.001)
		})
	}
}

func TestLineSegment_Rotate(t *testing.T) {
	tests := map[string]struct {
		seg      LineSegment[float64]
		pivot    point.Point[float64]
		radians  float64
		opts     []options.GeometryOptionsFunc
		inDelta  float64
		expected LineSegment[float64]
	}{
		"Rotate 90 degrees around origin": {
			seg: NewFromPoints[float64](
				point.New[float64](1.0, 0.0),
				point.New[float64](0.0, 1.0),
			),
			pivot:   point.New[float64](0.0, 0.0),
			radians: math.Pi / 2,
			opts:    []options.GeometryOptionsFunc{options.WithEpsilon(1e-9)},
			inDelta: 0.0001,
			expected: NewFromPoints[float64](
				point.New[float64](0.0, 1.0),
				point.New[float64](-1.0, 0.0),
			),
		},
		"Rotate 90 degrees around custom pivot": {
			seg: NewFromPoints[float64](
				point.New[float64](1.0, 0.0),
				point.New[float64](0.0, 1.0),
			),
			pivot:   point.New[float64](1.0, 1.0),
			radians: math.Pi / 2,
			opts:    []options.GeometryOptionsFunc{options.WithEpsilon(1e-9)},
			inDelta: 0.0001,
			expected: NewFromPoints[float64](
				point.New[float64](2.0, 1.0),
				point.New[float64](1.0, 0),
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rotatedLine := tc.seg.Rotate(tc.pivot, tc.radians, tc.opts...)
			assert.InDelta(t, tc.expected.start.X(), rotatedLine.start.X(), tc.inDelta)
			assert.InDelta(t, tc.expected.start.Y(), rotatedLine.start.Y(), tc.inDelta)
			assert.InDelta(t, tc.expected.end.X(), rotatedLine.end.X(), tc.inDelta)
			assert.InDelta(t, tc.expected.end.Y(), rotatedLine.end.Y(), tc.inDelta)
		})
	}
}

func TestLineSegment_Scale(t *testing.T) {
	tests := map[string]struct {
		segment  LineSegment[int]
		origin   point.Point[int]
		factor   int
		inDelta  float64
		expected LineSegment[int]
	}{
		// Integer test cases
		"int: Scale from start point by 2": {
			segment:  NewFromPoints[int](point.New(1, 1), point.New(4, 5)),
			origin:   point.New(1, 1),
			factor:   2,
			inDelta:  0.0001,
			expected: NewFromPoints[int](point.New[int](1, 1), point.New[int](7, 9)),
		},
		"int: Scale from end point by 2": {
			segment:  NewFromPoints[int](point.New(1, 1), point.New(4, 5)),
			origin:   point.New(4, 5),
			factor:   2,
			inDelta:  0.0001,
			expected: NewFromPoints[int](point.New[int](-2, -3), point.New[int](4, 5)),
		},
		"int: Scale from midpoint by 2": {
			segment:  NewFromPoints[int](point.New(0, 0), point.New(10, 10)),
			origin:   point.New(5, 5),
			factor:   2,
			inDelta:  0.0001,
			expected: NewFromPoints[int](point.New[int](-5, -5), point.New[int](15, 15)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.segment.Scale(tc.origin, tc.factor)
			assert.InDelta(t, tc.expected.start.X(), result.start.X(), tc.inDelta)
			assert.InDelta(t, tc.expected.start.Y(), result.start.Y(), tc.inDelta)
			assert.InDelta(t, tc.expected.end.X(), result.end.X(), tc.inDelta)
			assert.InDelta(t, tc.expected.end.Y(), result.end.Y(), tc.inDelta)
			t.Log(result.String())
		})
	}
}

func TestLineSegment_Slope(t *testing.T) {
	tests := map[string]struct {
		lineSegment LineSegment[int]
		expected    float64
		isNaN       bool
	}{
		"positive slope bottom left to top right": {
			lineSegment: NewFromPoints(point.New(1, 1), point.New(3, 3)),
			expected:    1.0,
			isNaN:       false,
		},
		"positive slope top right to bottom left": {
			lineSegment: NewFromPoints(point.New(3, 3), point.New(1, 1)),
			expected:    1.0,
			isNaN:       false,
		},
		"negative slope bottom right to top left": {
			lineSegment: New(10, 0, 0, 10),
			expected:    -1.0,
			isNaN:       false,
		},
		"negative slope top left to bottom right": {
			lineSegment: New(0, 10, 10, 0),
			expected:    -1.0,
			isNaN:       false,
		},
		"zero slope (horizontal)": {
			lineSegment: NewFromPoints(point.New(1, 2), point.New(5, 2)),
			expected:    0.0,
			isNaN:       false,
		},
		"undefined slope (vertical)": {
			lineSegment: NewFromPoints(point.New(2, 1), point.New(2, 5)),
			isNaN:       true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			slope := tc.lineSegment.Slope()

			if tc.isNaN {
				assert.True(t, math.IsNaN(slope), "expected the slope to be NaN")
			} else {
				assert.False(t, math.IsNaN(slope), "expected the slope to be a valid number")
				assert.InDelta(t, tc.expected, slope, 1e-6, "expected slope to match")
			}
		})
	}
}

func TestLineSegment_Start_End(t *testing.T) {
	tests := map[string]struct {
		start, end any
	}{
		"int: start: (1,2), end: (3,4)": {
			start: point.New[int](1, 2),
			end:   point.New[int](3, 4),
		},
		"float64: start: (1,2), end: (3,4)": {
			start: point.New[float64](1, 2),
			end:   point.New[float64](3, 4),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			switch start := tt.start.(type) {
			case point.Point[int]:
				end := tt.end.(point.Point[int])
				ls := NewFromPoints(start, end)
				assert.Equal(t, start, ls.Start())
				assert.Equal(t, end, ls.End())
			case point.Point[float64]:
				end := tt.end.(point.Point[float64])
				ls := NewFromPoints(start, end)
				assert.Equal(t, start, ls.Start())
				assert.Equal(t, end, ls.End())
			}
		})
	}
}

func TestLineSegment_String(t *testing.T) {
	tests := map[string]struct {
		segment  any    // Line segment to test (can be int or float64)
		expected string // Expected string output
	}{
		// Integer segment test cases
		"int: String representation": {
			segment:  NewFromPoints[int](point.New(1, 1), point.New(4, 5)),
			expected: "(1,1)(4,5)",
		},

		// Float64 segment test cases
		"float64: String representation": {
			segment:  NewFromPoints[float64](point.New(1.5, 1.5), point.New(4.5, 5.5)),
			expected: "(1.5,1.5)(4.5,5.5)",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			switch segment := tt.segment.(type) {
			case LineSegment[int]:
				result := segment.String()
				assert.Equal(t, tt.expected, result)

			case LineSegment[float64]:
				result := segment.String()
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestLineSegment_Translate(t *testing.T) {
	tests := map[string]struct {
		lineSegment LineSegment[int]
		delta       point.Point[int]
		expected    LineSegment[int]
	}{
		"translate by positive vector": {
			lineSegment: New(1, 1, 3, 3),
			delta:       point.New(2, 2),
			expected:    New(3, 3, 5, 5),
		},
		"translate by negative vector": {
			lineSegment: New(3, 3, 5, 5),
			delta:       point.New(-2, -2),
			expected:    New(1, 1, 3, 3),
		},
		"translate by zero vector": {
			lineSegment: New(1, 1, 3, 3),
			delta:       point.New(0, 0),
			expected:    New(1, 1, 3, 3),
		},
		"translate vertical line segment": {
			lineSegment: New(2, 2, 2, 5),
			delta:       point.New(3, -1),
			expected:    New(5, 1, 5, 4),
		},
		"translate horizontal line segment": {
			lineSegment: New(1, 4, 5, 4),
			delta:       point.New(0, -4),
			expected:    New(1, 0, 5, 0),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tc.lineSegment.Translate(tc.delta)
			assert.Equal(t, tc.expected, actual, "expected translated line segment to match")
		})
	}
}

func TestLineSegment_XAtY(t *testing.T) {
	tests := map[string]struct {
		lineSegment LineSegment[int]
		y           float64
		expected    float64
	}{
		"positive slope": {
			lineSegment: New(1, 1, 4, 4),
			y:           2,
			expected:    2.0,
		},
		"negative slope": {
			lineSegment: New(4, 4, 1, 1),
			y:           3,
			expected:    3.0,
		},
		"horizontal line": {
			lineSegment: New(1, 2, 5, 2),
			y:           2,
			expected:    math.NaN(),
		},
		"vertical line": {
			lineSegment: New(2, 1, 2, 5),
			y:           3,
			expected:    2.0,
		},
		"y outside bounds": {
			lineSegment: New(1, 1, 4, 4),
			y:           5,
			expected:    math.NaN(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tc.lineSegment.XAtY(tc.y)
			if math.IsNaN(tc.expected) {
				assert.True(t, math.IsNaN(actual), "expected NaN")
			} else {
				assert.InDelta(t, tc.expected, actual, 1e-6, "expected x-coordinate to match")
			}
		})
	}
}

func TestLineSegment_YAtX(t *testing.T) {
	tests := map[string]struct {
		lineSegment LineSegment[int]
		x           float64
		expected    float64
	}{
		"positive slope": {
			lineSegment: New(1, 1, 4, 4),
			x:           2,
			expected:    2.0,
		},
		"negative slope": {
			lineSegment: New(4, 4, 1, 1),
			x:           3,
			expected:    3.0,
		},
		"horizontal line": {
			lineSegment: New(1, 2, 5, 2),
			x:           3,
			expected:    2.0,
		},
		"vertical line": {
			lineSegment: New(2, 1, 2, 5),
			x:           2,
			expected:    math.NaN(),
		},
		"x outside bounds": {
			lineSegment: New(1, 1, 4, 4),
			x:           5,
			expected:    math.NaN(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tc.lineSegment.YAtX(tc.x)
			if math.IsNaN(tc.expected) {
				assert.True(t, math.IsNaN(actual), "expected NaN")
			} else {
				assert.InDelta(t, tc.expected, actual, 1e-6, "expected y-coordinate to match")
			}
		})
	}
}
