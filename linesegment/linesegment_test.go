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
		defined     bool
	}{
		"positive slope": {
			lineSegment: NewFromPoints(point.New(1, 1), point.New(3, 3)),
			expected:    1.0,
			defined:     true,
		},
		"negative slope": {
			lineSegment: NewFromPoints(point.New(3, 3), point.New(1, 1)),
			expected:    1.0,
			defined:     true,
		},
		"zero slope": {
			lineSegment: NewFromPoints(point.New(1, 2), point.New(5, 2)),
			expected:    0.0,
			defined:     true,
		},
		"undefined slope": {
			lineSegment: NewFromPoints(point.New(2, 1), point.New(2, 5)),
			expected:    0.0,
			defined:     false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			slope, defined := tc.lineSegment.Slope()
			if tc.defined {
				assert.True(t, defined, "expected the slope to be defined")
				assert.InDelta(t, tc.expected, slope, 1e-6, "expected slope to match")
			} else {
				assert.False(t, defined, "expected the slope to be undefined")
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
