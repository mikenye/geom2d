package _old

import (
	"github.com/stretchr/testify/require"
	"image"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkConvexHull(b *testing.B) {
	points := []Point[int]{
		{1, 4}, {4, 13}, {8, 17}, {18, 20}, {33, 18}, {38, 11},
		{34, -2}, {21, -3}, {6, -1}, {7, 6}, {10, 14}, {5, 2},
		{16, 0}, {12, 12}, {23, 16}, {14, 6}, {24, 0}, {21, -2},
		{30, 14}, {27, 9}, {29, 4}, {31, -1}, {34, 7}, {35, 12},
	}

	// Reset the timer to ignore any setup cost
	b.ResetTimer()

	// Run the benchmark loop
	for i := 0; i < b.N; i++ {
		_ = ConvexHull(points)
	}
}

func TestPointOrientation_String(t *testing.T) {
	tests := map[string]struct {
		input          PointOrientation
		expectedOutput string
	}{
		"PointsCollinear": {
			input:          PointsCollinear,
			expectedOutput: "PointsCollinear",
		},
		"PointsClockwise": {
			input:          PointsClockwise,
			expectedOutput: "PointsClockwise",
		},
		"PointsCounterClockwise": {
			input:          PointsCounterClockwise,
			expectedOutput: "PointsCounterClockwise",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expectedOutput, tc.input.String())
		})
	}
}

func TestPoint_RelationshipToLineSegment(t *testing.T) {
	segment := NewLineSegment(
		NewPoint(0, 0),
		NewPoint(10, 0),
	)

	tests := map[string]struct {
		point       Point[int]
		expectedRel Relationship
	}{
		"Point lies on the segment": {
			point:       NewPoint(5, 0),
			expectedRel: RelationshipIntersection,
		},
		"Point lies on the segment start": {
			point:       NewPoint(0, 0),
			expectedRel: RelationshipIntersection,
		},
		"Point lies on the segment end": {
			point:       NewPoint(10, 0),
			expectedRel: RelationshipIntersection,
		},
		"Point lies outside the segment": {
			point:       NewPoint(15, 0),
			expectedRel: RelationshipDisjoint,
		},
		"Point lies above the segment": {
			point:       NewPoint(5, 2),
			expectedRel: RelationshipDisjoint,
		},
		"Point lies below the segment": {
			point:       NewPoint(5, -2),
			expectedRel: RelationshipDisjoint,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			rel := tt.point.RelationshipToLineSegment(segment, WithEpsilon(1e-10))
			assert.Equal(t, tt.expectedRel, rel, "unexpected relationship")
		})
	}
}

func TestPoint_RelationshipToPoint(t *testing.T) {
	tests := map[string]struct {
		pointA      Point[int]
		pointB      Point[int]
		expectedRel Relationship
	}{
		"Points are equal": {
			pointA:      NewPoint(5, 5),
			pointB:      NewPoint(5, 5),
			expectedRel: RelationshipEqual,
		},
		"Points are disjoint": {
			pointA:      NewPoint(5, 5),
			pointB:      NewPoint(10, 10),
			expectedRel: RelationshipDisjoint,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			rel := tt.pointA.RelationshipToPoint(tt.pointB, WithEpsilon(1e-10))
			assert.Equal(t, tt.expectedRel, rel, "unexpected relationship")
		})
	}
}

func TestPoint_RelationshipToPolyTree(t *testing.T) {
	t.Run("issue 7: Point is left of and collinear to a PolyTree edge", func(t *testing.T) {
		holeInPolyTree, err := NewPolyTree[int]([]Point[int]{
			NewPoint(299, 191),
			NewPoint(329, 195),
			NewPoint(325, 210),
			NewPoint(298, 211),
		}, PTHole)
		require.NoError(t, err, "error creating holeInPolyTree")
		polyTree, err := NewPolyTree[int]([]Point[int]{
			NewPoint(333, 218),
			NewPoint(345, 195),
			NewPoint(324, 181),
			NewPoint(341, 164),
			NewPoint(307, 169),
			NewPoint(270, 163),
			NewPoint(254, 180),
			NewPoint(263, 193),
			NewPoint(253, 210),
			NewPoint(290, 181),
			NewPoint(288, 218),
		}, PTSolid, WithChildren(holeInPolyTree))
		require.NoError(t, err, "error creating polyTree")
		point := NewPoint(273, 218)
		rel := point.RelationshipToPolyTree(polyTree)
		assert.Equal(t, RelationshipDisjoint, rel[polyTree])
		assert.Equal(t, RelationshipDisjoint, rel[holeInPolyTree])
	})
}

func TestRelativeAngle(t *testing.T) {
	tests := map[string]struct {
		A, B, O  any
		expected float64
	}{
		"int: ~22 deg": {
			A:        NewPoint[int](10, 0),
			B:        NewPoint[int](10, 4),
			O:        NewPoint[int](0, 0),
			expected: 0.38050637717509, // ~22 deg
		},
		"int: ~66 deg": {
			A:        NewPoint[int](10, 0),
			B:        NewPoint[int](4, 9),
			O:        NewPoint[int](0, 0),
			expected: 1.1525719972927, // ~66 deg
		},
		"int: ~90 deg": {
			A:        NewPoint[int](10, 0),
			B:        NewPoint[int](0, 10),
			O:        NewPoint[int](0, 0),
			expected: 1.5708, // ~90 deg
		},
		"int: ~114 deg": {
			A:        NewPoint[int](10, 0),
			B:        NewPoint[int](-4, 9),
			O:        NewPoint[int](0, 0),
			expected: 1.98902065628929, // ~114 deg
		},
		"int: ~156 deg": {
			A:        NewPoint[int](10, 0),
			B:        NewPoint[int](-9, 4),
			O:        NewPoint[int](0, 0),
			expected: 2.72336832408371, // ~156 deg
		},
		"float64: ~22 deg": {
			A:        NewPoint[float64](10, 0),
			B:        NewPoint[float64](10, 4),
			O:        NewPoint[float64](0, 0),
			expected: 0.38050637717509, // ~22 deg
		},
		"float64: ~66 deg": {
			A:        NewPoint[float64](10, 0),
			B:        NewPoint[float64](4, 9),
			O:        NewPoint[float64](0, 0),
			expected: 1.1525719972927, // ~66 deg
		},
		"float64: ~90 deg": {
			A:        NewPoint[float64](10, 0),
			B:        NewPoint[float64](0, 10),
			O:        NewPoint[float64](0, 0),
			expected: 1.5708, // ~90 deg
		},
		"float64: ~114 deg": {
			A:        NewPoint[float64](10, 0),
			B:        NewPoint[float64](-4, 9),
			O:        NewPoint[float64](0, 0),
			expected: 1.98902065628929, // ~114 deg
		},
		"float64: ~156 deg": {
			A:        NewPoint[float64](10, 0),
			B:        NewPoint[float64](-9, 4),
			O:        NewPoint[float64](0, 0),
			expected: 2.72336832408371, // ~156 deg
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			switch A := tt.A.(type) {
			case Point[int]:
				B := tt.B.(Point[int])
				O := tt.O.(Point[int])
				actual := RelativeAngle(A, B, O)
				assert.InDelta(t, tt.expected, actual, 0.00001)

			case Point[float64]:
				B := tt.B.(Point[float64])
				O := tt.O.(Point[float64])
				actual := RelativeAngle(A, B, O)
				assert.InDelta(t, tt.expected, actual, 0.00001)
			}
		})
	}
}

func TestRelativeAngle_NoOriginGiven(t *testing.T) {
	tests := map[string]struct {
		A, B     any
		expected float64
	}{
		"int: right triangle, origin at 0,0": {
			A:        NewPoint[int](10, 10),
			B:        NewPoint[int](10, 0),
			expected: 0.785398, // 45 deg
		},
		"float64: right triangle, origin at 0,0": {
			A:        NewPoint[float64](10, 10),
			B:        NewPoint[float64](10, 0),
			expected: 0.785398, // 45 deg
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			switch A := tt.A.(type) {
			case Point[int]:
				B := tt.B.(Point[int])
				actual := RelativeAngle(A, B)
				assert.InDelta(t, tt.expected, actual, 0.00001)

			case Point[float64]:
				B := tt.B.(Point[float64])
				actual := RelativeAngle(A, B)
				assert.InDelta(t, tt.expected, actual, 0.00001)
			}
		})
	}
}

func TestSignedArea2X(t *testing.T) {
	tests := map[string]struct {
		points   []Point[int]
		expected int
	}{
		"counterclockwise triangle": {
			points:   []Point[int]{NewPoint(0, 0), NewPoint(4, 0), NewPoint(2, 3)},
			expected: 12,
		},
		"clockwise triangle": {
			points:   []Point[int]{NewPoint(0, 0), NewPoint(2, 3), NewPoint(4, 0)},
			expected: -12,
		},
		"collinear points": {
			points:   []Point[int]{NewPoint(0, 0), NewPoint(2, 2), NewPoint(4, 4)},
			expected: 0,
		},
		"two points, not a polygon, no area": {
			points:   []Point[int]{NewPoint(0, 0), NewPoint(2, 2)},
			expected: 0,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := SignedArea2X(tt.points)
			assert.Equal(t, tt.expected, result)
		})
	}
}
