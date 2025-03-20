package linesegment

import (
	"encoding/json"
	"github.com/mikenye/geom2d"
	"github.com/mikenye/geom2d/point"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func TestLineSegment_Bresenham(t *testing.T) {
	tests := map[string]struct {
		lineSegment LineSegment
		expected    []point.Point
	}{
		"horizontal line": {
			lineSegment: NewFromPoints(point.New(0, 0), point.New(5, 0)),
			expected: []point.Point{
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
			expected: []point.Point{
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
			expected: []point.Point{
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
			expected: []point.Point{
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
			expected: []point.Point{
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
			var actual []point.Point
			test.lineSegment.Bresenham(func(p point.Point) bool {
				actual = append(actual, p)
				return true
			})
			assert.ElementsMatch(t, test.expected, actual, "Bresenham points mismatch")
		})
	}
}

func TestLineSegment_Center(t *testing.T) {
	tests := map[string]struct {
		lineSegment LineSegment
		epsilon     float64
		expected    point.Point
	}{
		"No epsilon, simple case": {
			lineSegment: New(0, 0, 4, 4),
			epsilon:     0,
			expected:    point.New(2, 2),
		},
		"No epsilon, negative coordinates": {
			lineSegment: New(-4, -4, 4, 4),
			epsilon:     0,
			expected:    point.New(0, 0),
		},
		"With epsilon, rounding applied": {
			lineSegment: New(0, 0, 3, 3),
			epsilon:     0.1,
			expected:    point.New(1.5, 1.5), // No rounding as it's precise
		},
		"With epsilon, midpoint near integer": {
			lineSegment: New(0, 0, 4, 5),
			epsilon:     0.5,
			expected:    point.New(2, 2.5), // Epsilon not applied due to midpoint already exact
		},
		"With epsilon, midpoint adjusted to integer": {
			lineSegment: New(0, 0, 5, 5),
			epsilon:     0.5,
			expected:    point.New(2.5, 2.5), // Exact match without adjustment
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Calculate center
			center := tc.lineSegment.Center()

			// Assert the result
			assert.InDelta(t, tc.expected.X(), center.X(), geom2d.GetEpsilon(), "Unexpected x-coordinate for center")
			assert.InDelta(t, tc.expected.Y(), center.Y(), geom2d.GetEpsilon(), "Unexpected y-coordinate for center")
		})
	}
}

func TestLineSegment_ContainsPoint(t *testing.T) {
	tests := map[string]struct {
		segment  LineSegment
		point    point.Point
		expected bool
	}{
		"point lies on the segment": {
			segment:  New(0.0, 0.0, 10.0, 10.0),
			point:    point.New(5.0, 5.0),
			expected: true,
		},
		"point is off the segment but collinear": {
			segment:  New(0.0, 0.0, 10.0, 10.0),
			point:    point.New(15.0, 15.0),
			expected: false,
		},
		"point lies at the start of the segment": {
			segment:  New(0.0, 0.0, 10.0, 10.0),
			point:    point.New(0.0, 0.0),
			expected: true,
		},
		"point lies at the end of the segment": {
			segment:  New(0.0, 0.0, 10.0, 10.0),
			point:    point.New(10.0, 10.0),
			expected: true,
		},
		"floating point precision issue - outside of default epsilon": {
			segment:  New(0.0, 0.0, 10.0, 10.0),
			point:    point.New(5.00000000000001, 5.0),
			expected: true,
		},
		"floating point precision issue - within default epsilon": {
			segment:  New(0.0, 0.0, 10.0, 10.0),
			point:    point.New(5.0000000001, 5.0),
			expected: false,
		},
		"point is not on the segment": {
			segment:  New(0.0, 0.0, 10.0, 10.0),
			point:    point.New(5.0, 6.0),
			expected: false,
		},
		"point lies on degenerate segment": {
			segment:  New(5.0, 5.0, 5.0, 5.0),
			point:    point.New(5.0, 5.0),
			expected: true,
		},
		"point does not lie on degenerate segment": {
			segment:  New(5.0, 5.0, 5.0, 5.0),
			point:    point.New(6.0, 5.0),
			expected: false,
		},
		"FuzzFindIntersections_2segments/7500abbb0aa68f4e/A": {
			segment:  New(-44.285714285714285, 151, 20, 10),
			point:    point.New(-15.467349551856598, 87.79172001707214),
			expected: true,
		},
		"FuzzFindIntersections_2segments/7500abbb0aa68f4e/B": {
			segment:  New(-10, 88, -640, 64),
			point:    point.New(-15.467349551856598, 87.79172001707214),
			expected: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := tt.segment.ContainsPoint(tt.point)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLineSegment_DistanceToLineSegment(t *testing.T) {
	tests := map[string]struct {
		segA, segB LineSegment
		expected   float64
		expectZero bool
	}{
		"intersecting segments": {
			segA:       New(0, 0, 4, 4),
			segB:       New(0, 4, 4, 0),
			expected:   0,
			expectZero: true,
		},
		"parallel non-intersecting segments": {
			segA:       New(0, 0, 4, 0),
			segB:       New(0, 2, 4, 2),
			expected:   2,
			expectZero: false,
		},
		"segments touching at one endpoint": {
			segA:       New(0, 0, 4, 0),
			segB:       New(4, 0, 4, 4),
			expected:   0,
			expectZero: true,
		},
		"skew non-intersecting segments": {
			segA:       New(0, 0, 1, 1),
			segB:       New(2, 0, 3, 1),
			expected:   1.4142135623731,
			expectZero: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.segA.DistanceToLineSegment(tt.segB)
			if tt.expectZero {
				assert.InDelta(t, 0, actual, geom2d.GetEpsilon(), "Expected distance to be zero")
			} else {
				assert.InDelta(t, tt.expected, actual, geom2d.GetEpsilon(), "Expected distances to match")
			}
		})
	}
}

func TestLineSegment_DistanceToPoint(t *testing.T) {
	tests := map[string]struct {
		point    point.Point // Point to be projected (can be int or float64)
		segment  LineSegment // Line segment for projection (can be int or float64)
		expected float64     // Expected distance
	}{
		"Project onto segment from inside": {
			point:    point.New(5.5, 5.5),
			segment:  New(2.0, 3.0, 8.0, 7.0),
			expected: 0.1386750490563,
		},
		"Project onto segment from above": {
			point:    point.New(4.0, 6.0),
			segment:  New(2.0, 3.0, 8.0, 7.0),
			expected: 1.3867504905631,
		},
		"Project onto segment from below": {
			point:    point.New(4.0, 2.0),
			segment:  New(2.0, 3.0, 8.0, 7.0),
			expected: 1.9414506867883,
		},
		"Project off the start of segment": {
			point:    point.New(0.0, 5.0),
			segment:  New(2.0, 3.0, 8.0, 7.0),
			expected: 2.8284271247462,
		},
		"Project off the end of segment": {
			point:    point.New(10.0, 5.0),
			segment:  New(2.0, 3.0, 8.0, 7.0),
			expected: 2.8284271247462,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.segment.DistanceToPoint(tt.point)
			assert.InDelta(t, tt.expected, actual, geom2d.GetEpsilon(), "Expected distance does not match actual distance for float64 points")
		})
	}
}

func TestLineSegment_Eq(t *testing.T) {
	tests := map[string]struct {
		segment1 LineSegment // First line segment (can be int or float64)
		segment2 LineSegment // Second line segment to compare (can be int or float64)
		expected bool        // Expected result of equality check
	}{
		"float64: Equal segments": {
			segment1: NewFromPoints(point.New(1.0, 1.0), point.New(4.0, 5.0)),
			segment2: NewFromPoints(point.New(1.0, 1.0), point.New(4.0, 5.0)),
			expected: true,
		},
		"float64: Unequal segments": {
			segment1: NewFromPoints(point.New(1.5, 1.5), point.New(3.5, 4.5)),
			segment2: NewFromPoints(point.New(1.5, 1.5), point.New(5.5, 6.5)),
			expected: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.segment1.Eq(tc.segment2)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestLineSegment_Length(t *testing.T) {
	tests := map[string]struct {
		lineSegment LineSegment
		expected    float64
	}{
		"horizontal segment": {
			lineSegment: New(0, 0, 5, 0),
			expected:    5.0,
		},
		"vertical segment": {
			lineSegment: New(0, 0, 0, 7),
			expected:    7.0,
		},
		"diagonal segment": {
			lineSegment: New(0, 0, 3, 4),
			expected:    5.0,
		},
		"near-zero length segment": {
			lineSegment: New(1e-10, 1e-10, 2e-10, 2e-10),
			expected:    1.4142e-10,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.lineSegment.Length()
			assert.InDelta(t, tt.expected, actual, geom2d.GetEpsilon(), "Expected length to match")
		})
	}
}

func TestLineSegment_MarshalUnmarshalJSON(t *testing.T) {
	tests := map[string]struct {
		segment  LineSegment // Input segment
		expected LineSegment // Expected output after Marshal -> Unmarshal
	}{
		"LineSegment[float64]": {
			segment:  New(3.5, 7.2, -4.1, 2.8),
			expected: New(3.5, 7.2, -4.1, 2.8),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Marshal
			data, err := json.Marshal(tc.segment)
			require.NoErrorf(t, err, "Failed to marshal %s: %v", tc.segment, err)

			// Determine the correct type for unmarshalling
			var result LineSegment
			err = json.Unmarshal(data, &result)
			require.NoErrorf(t, err, "Failed to unmarshal `%s`: %v", string(data), err)
			assert.Equalf(t, tc.expected, result, "Expected %s, got %s", tc.expected, result)
		})
	}
}

func TestLineSegment_Points(t *testing.T) {
	tests := map[string]struct {
		segment  LineSegment
		expected []point.Point
	}{
		"horizontal segment": {
			segment: NewFromPoints(point.New(1, 1), point.New(5, 1)),
			expected: []point.Point{
				point.New(1, 1), // "toppest-leftest" first
				point.New(5, 1),
			},
		},
		"vertical segment": {
			segment: NewFromPoints(point.New(3, 2), point.New(3, 6)),
			expected: []point.Point{
				point.New(3, 6), // "toppest-leftest" first
				point.New(3, 2),
			},
		},
		"diagonal segment": {
			segment: NewFromPoints(point.New(0, 0), point.New(3, 4)),
			expected: []point.Point{
				point.New(3, 4), // "toppest-leftest" first
				point.New(0, 0),
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
		point    point.Point // Point to be projected (can be int or float64)
		segment  LineSegment // Line segment for projection (can be int or float64)
		expected point.Point // Expected projected point (float64 type)
	}{
		"Project onto segment from inside": {
			point:    point.New(5.5, 5.5),
			segment:  New(2.0, 3.0, 8.0, 7.0),
			expected: point.New(5.5769230769231, 5.3846153846154), // Should project onto the line segment itself
		},
		"Project off the start of segment": {
			point:    point.New(0.0, 5.0),
			segment:  New(2.0, 3.0, 8.0, 7.0),
			expected: point.New(2.0, 3.0), // Should return a point of the segment
		},
		"Project off the end of segment": {
			point:    point.New(10.0, 5.0),
			segment:  New(2.0, 3.0, 8.0, 7.0),
			expected: point.New(8.0, 7.0), // Should return end point of the segment
		},
		"Project onto zero-length segment": {
			point:    point.New(5.0, 5.0),
			segment:  New(2.5, 2.5, 2.5, 2.5), // Zero-length segment
			expected: point.New(2.5, 2.5),     // Should return point A (or End), since the segment is a single point
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.segment.ProjectPoint(tt.point)
			assert.InDelta(t, tt.expected.X(), actual.X(), geom2d.GetEpsilon())
			assert.InDelta(t, tt.expected.Y(), actual.Y(), geom2d.GetEpsilon())

		})
	}
}

func TestLineSegment_ReflectLineSegment(t *testing.T) {
	tests := map[string]struct {
		lineSegment LineSegment
		other       LineSegment
		expected    LineSegment
	}{
		"horizontal reflection": {
			lineSegment: New(0, 0, 10, 0),
			other:       New(2, 2, 8, 2),
			expected:    New(2, -2, 8, -2),
		},
		"vertical reflection": {
			lineSegment: New(0, 0, 0, 10),
			other:       New(2, 2, 2, 8),
			expected:    New(-2, 2, -2, 8),
		},
		"diagonal reflection": {
			lineSegment: New(0, 0, 10, 10),
			other:       New(2, 6, 6, 2),
			expected:    New(6, 2, 2, 6),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.lineSegment.ReflectLineSegment(tt.other)
			assert.InDelta(t, tt.expected.Upper().X(), actual.Upper().X(), geom2d.GetEpsilon(), "Upper X mismatch")
			assert.InDelta(t, tt.expected.Upper().Y(), actual.Upper().Y(), geom2d.GetEpsilon(), "Upper Y mismatch")
			assert.InDelta(t, tt.expected.Lower().X(), actual.Lower().X(), geom2d.GetEpsilon(), "Lower X mismatch")
			assert.InDelta(t, tt.expected.Lower().Y(), actual.Lower().Y(), geom2d.GetEpsilon(), "Lower Y mismatch")
		})
	}
}

func TestLineSegment_ReflectPoint(t *testing.T) {
	tests := map[string]struct {
		point    point.Point // The point to reflect
		axis     LineSegment // Axis for reflection
		expected point.Point // Expected reflected point
	}{
		"reflect across x-axis": {
			point:    point.New(3, 4),
			axis:     New(0, 0, 1, 0),
			expected: point.New(3, -4),
		},
		"reflect across y-axis": {
			point:    point.New(3, 4),
			axis:     New(0, 0, 0, 1),
			expected: point.New(-3, 4),
		},
		"reflect across y = x line (ReflectAcrossCustomLine)": {
			point:    point.New(3, 4),
			axis:     New(0, 0, 1, 1),
			expected: point.New(4, 3),
		},
		"reflect across y = -x line (ReflectAcrossCustomLine)": {
			point:    point.New(3, 4),
			axis:     New(0, 0, -1, 1),
			expected: point.New(-4, -3),
		},
		"reflect across degenerate line segment": {
			point:    point.New(3, 4),
			axis:     New(1, 1, 1, 1), // Degenerate line
			expected: point.New(3, 4), // Expect the point to remain unchanged
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.axis.ReflectPoint(tc.point)
			assert.InDelta(t, tc.expected.X(), result.X(), geom2d.GetEpsilon())
			assert.InDelta(t, tc.expected.Y(), result.Y(), geom2d.GetEpsilon())
		})
	}
}

func TestLineSegment_Rotate(t *testing.T) {
	tests := map[string]struct {
		seg      LineSegment
		pivot    point.Point
		radians  float64
		expected LineSegment
	}{
		"Rotate 90 degrees around origin": {
			seg: NewFromPoints(
				point.New(1.0, 0.0),
				point.New(0.0, 1.0),
			),
			pivot:   point.New(0.0, 0.0),
			radians: math.Pi / 2,
			expected: NewFromPoints(
				point.New(0.0, 1.0),
				point.New(-1.0, 0.0),
			),
		},
		"Rotate 90 degrees around custom pivot": {
			seg: NewFromPoints(
				point.New(1.0, 0.0),
				point.New(0.0, 1.0),
			),
			pivot:   point.New(1.0, 1.0),
			radians: math.Pi / 2,
			expected: NewFromPoints(
				point.New(2.0, 1.0),
				point.New(1.0, 0),
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rotatedLine := tc.seg.Rotate(tc.pivot, tc.radians)
			assert.InDelta(t, tc.expected.Upper().X(), rotatedLine.Upper().X(), geom2d.GetEpsilon())
			assert.InDelta(t, tc.expected.Upper().Y(), rotatedLine.Upper().Y(), geom2d.GetEpsilon())
			assert.InDelta(t, tc.expected.Lower().X(), rotatedLine.Lower().X(), geom2d.GetEpsilon())
			assert.InDelta(t, tc.expected.Lower().Y(), rotatedLine.Lower().Y(), geom2d.GetEpsilon())
		})
	}
}

func TestLineSegment_Scale(t *testing.T) {
	tests := map[string]struct {
		segment  LineSegment
		origin   point.Point
		factor   float64
		expected LineSegment
	}{
		// Integer test cases
		"int: Scale from start point by 2": {
			segment:  NewFromPoints(point.New(1, 1), point.New(4, 5)),
			origin:   point.New(1, 1),
			factor:   2,
			expected: NewFromPoints(point.New(1, 1), point.New(7, 9)),
		},
		"int: Scale from end point by 2": {
			segment:  NewFromPoints(point.New(1, 1), point.New(4, 5)),
			origin:   point.New(4, 5),
			factor:   2,
			expected: NewFromPoints(point.New(-2, -3), point.New(4, 5)),
		},
		"int: Scale from midpoint by 2": {
			segment:  NewFromPoints(point.New(0, 0), point.New(10, 10)),
			origin:   point.New(5, 5),
			factor:   2,
			expected: NewFromPoints(point.New(-5, -5), point.New(15, 15)),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.segment.Scale(tc.origin, tc.factor)
			assert.InDelta(t, tc.expected.Upper().X(), result.Upper().X(), geom2d.GetEpsilon())
			assert.InDelta(t, tc.expected.Upper().Y(), result.Upper().Y(), geom2d.GetEpsilon())
			assert.InDelta(t, tc.expected.Lower().X(), result.Lower().X(), geom2d.GetEpsilon())
			assert.InDelta(t, tc.expected.Lower().Y(), result.Lower().Y(), geom2d.GetEpsilon())
			t.Log(result.String())
		})
	}
}

func TestLineSegment_Slope(t *testing.T) {
	tests := map[string]struct {
		lineSegment LineSegment
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

func TestLineSegment_Upper_Lower(t *testing.T) {
	tests := map[string]struct {
		upper, lower point.Point
	}{
		"float64: start: (1,2), end: (3,4)": {
			upper: point.New(3, 4),
			lower: point.New(1, 2),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {

			ls := NewFromPoints(tt.upper, tt.lower)
			assert.Equal(t, tt.upper, ls.Upper())
			assert.Equal(t, tt.lower, ls.Lower())

		})
	}
}

func TestLineSegment_String(t *testing.T) {
	tests := map[string]struct {
		segment  LineSegment // Line segment to test
		expected string      // Expected string output
	}{
		"String representation": {
			segment:  NewFromPoints(point.New(1.5, 1.5), point.New(4.5, 5.5)),
			expected: "(4.5,5.5)(1.5,1.5)",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := tt.segment.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLineSegment_Translate(t *testing.T) {
	tests := map[string]struct {
		lineSegment LineSegment
		delta       point.Point
		expected    LineSegment
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
		lineSegment LineSegment
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
		lineSegment LineSegment
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
