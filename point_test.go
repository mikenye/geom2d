package geom2d

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

func TestPoint_CrossProduct(t *testing.T) {
	tests := []struct {
		name     string
		p, q     any // Support different Point types with `any`
		expected any // Expected result for different types
	}{
		// Integer points
		{
			name:     "int: (2,3) x (4,5)",
			p:        NewPoint(2, 3),
			q:        NewPoint(4, 5),
			expected: -2,
		},
		{
			name:     "int: (3,2) x (4,6)",
			p:        NewPoint(3, 2),
			q:        NewPoint(4, 6),
			expected: 10,
		},

		// Float64 points
		{
			name:     "float64: (2.0,3.0) x (4.0,5.0)",
			p:        NewPoint(2.0, 3.0),
			q:        NewPoint(4.0, 5.0),
			expected: -2.0,
		},
		{
			name:     "float64: (3.5,2.5) x (4.0,6.0)",
			p:        NewPoint(3.5, 2.5),
			q:        NewPoint(4.0, 6.0),
			expected: 11.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch p := tt.p.(type) {
			case Point[int]:
				q := tt.q.(Point[int])
				expected := tt.expected.(int)
				actual := p.CrossProduct(q)
				assert.Equal(t, expected, actual)

			case Point[float64]:
				q := tt.q.(Point[float64])
				expected := tt.expected.(float64)
				actual := p.CrossProduct(q)
				assert.Equal(t, expected, actual)
			}
		})
	}
}

func TestPoint_DistanceToLineSegment(t *testing.T) {
	tests := map[string]struct {
		point    any     // Point to be projected (can be int or float64)
		segment  any     // Line segment for projection (can be int or float64)
		expected float64 // Expected distance
	}{
		// Integer points test cases
		"Project onto segment from inside (int)": {
			point:    NewPoint(5, 5),
			segment:  NewLineSegment[int](NewPoint(2, 3), NewPoint(8, 7)),
			expected: 0.0, // Point is on the segment
		},
		"Project onto segment from above (int)": {
			point:    NewPoint(4, 6),
			segment:  NewLineSegment[int](NewPoint(2, 3), NewPoint(8, 7)),
			expected: 1.386,
		},
		"Project onto segment from below (int)": {
			point:    NewPoint(4, 2),
			segment:  NewLineSegment[int](NewPoint(2, 3), NewPoint(8, 7)),
			expected: 1.941,
		},
		"Project off the start of segment (int)": {
			point:    NewPoint(0, 5),
			segment:  NewLineSegment[int](NewPoint(2, 3), NewPoint(8, 7)),
			expected: 2.8284,
		},
		"Project off the end of segment (int)": {
			point:    NewPoint(10, 5),
			segment:  NewLineSegment[int](NewPoint(2, 3), NewPoint(8, 7)),
			expected: 2.8284,
		},

		// Float64 points test cases
		"Project onto segment from inside (float64)": {
			point:    NewPoint(5.5, 5.5),
			segment:  NewLineSegment[float64](NewPoint(2.0, 3.0), NewPoint(8.0, 7.0)),
			expected: 0.1387,
		},
		"Project onto segment from above (float64)": {
			point:    NewPoint(4.0, 6.0),
			segment:  NewLineSegment[float64](NewPoint(2.0, 3.0), NewPoint(8.0, 7.0)),
			expected: 1.386,
		},
		"Project onto segment from below (float64)": {
			point:    NewPoint(4.0, 2.0),
			segment:  NewLineSegment[float64](NewPoint(2.0, 3.0), NewPoint(8.0, 7.0)),
			expected: 1.941,
		},
		"Project off the start of segment (float64)": {
			point:    NewPoint(0.0, 5.0),
			segment:  NewLineSegment[float64](NewPoint(2.0, 3.0), NewPoint(8.0, 7.0)),
			expected: 2.8284,
		},
		"Project off the end of segment (float64)": {
			point:    NewPoint(10.0, 5.0),
			segment:  NewLineSegment[float64](NewPoint(2.0, 3.0), NewPoint(8.0, 7.0)),
			expected: 2.8284,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			switch p := tt.point.(type) {
			case Point[int]:
				seg := tt.segment.(LineSegment[int])
				actual := p.DistanceToLineSegment(seg)
				assert.InDelta(t, tt.expected, actual, 0.001, "Expected distance does not match actual distance for int points")

			case Point[float64]:
				seg := tt.segment.(LineSegment[float64])
				actual := p.DistanceToLineSegment(seg)
				assert.InDelta(t, tt.expected, actual, 0.001, "Expected distance does not match actual distance for float64 points")
			}
		})
	}
}

func TestPoint_DistanceToPoint(t *testing.T) {
	tests := []struct {
		name     string
		p, q     any     // Use `any` to handle different Point types
		expected float64 // Expected distance as a float64
	}{
		// Integer points
		{
			name:     "int: distance between (2,10) and (10,2)",
			p:        NewPoint(2, 10),
			q:        NewPoint(10, 2),
			expected: 11.3137,
		},
		{
			name:     "int: distance between (0,0) and (3,4)",
			p:        NewPoint(0, 0),
			q:        NewPoint(3, 4),
			expected: 5.0,
		},

		// Float64 points
		{
			name:     "float64: distance between (2.0,10.0) and (10.0,2.0)",
			p:        NewPoint(2.0, 10.0),
			q:        NewPoint(10.0, 2.0),
			expected: 11.3137,
		},
		{
			name:     "float64: distance between (0.0,0.0) and (3.0,4.0)",
			p:        NewPoint(0.0, 0.0),
			q:        NewPoint(3.0, 4.0),
			expected: 5.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch p := tt.p.(type) {
			case Point[int]:
				q := tt.q.(Point[int])
				actual := p.DistanceToPoint(q)
				assert.InDelta(t, tt.expected, actual, 0.0001)

			case Point[float64]:
				q := tt.q.(Point[float64])
				actual := p.DistanceToPoint(q)
				assert.InDelta(t, tt.expected, actual, 0.0001)
			}
		})
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

func TestPoint_ProjectOntoLineSegment(t *testing.T) {
	tests := map[string]struct {
		point    any            // Point to be projected (can be int or float64)
		segment  any            // Line segment for projection (can be int or float64)
		expected Point[float64] // Expected projected point (float64 type)
	}{
		// Integer points test cases
		"int: Project onto segment from inside": {
			point:    NewPoint(5, 5),
			segment:  NewLineSegment[int](NewPoint(2, 3), NewPoint(8, 7)),
			expected: NewPoint[float64](5, 5), // Should project onto the line segment itself
		},
		"int: Project onto segment from above": {
			point:    NewPoint(4, 6),
			segment:  NewLineSegment[int](NewPoint(2, 3), NewPoint(8, 7)),
			expected: NewPoint[float64](4.769, 4.846),
		},
		"int: Project onto segment from below": {
			point:    NewPoint(4, 2),
			segment:  NewLineSegment[int](NewPoint(2, 3), NewPoint(8, 7)),
			expected: NewPoint[float64](2.9231, 3.6154), // Projected point should be (5,5)
		},
		"int: Project off the start of segment": {
			point:    NewPoint(0, 5),
			segment:  NewLineSegment[int](NewPoint(2, 3), NewPoint(8, 7)),
			expected: NewPoint[float64](2, 3), // Should return a point of the segment
		},
		"int: Project off the end of segment": {
			point:    NewPoint(10, 5),
			segment:  NewLineSegment[int](NewPoint(2, 3), NewPoint(8, 7)),
			expected: NewPoint[float64](8, 7), // Should return end point of the segment
		},

		// Float64 points test cases
		"float64: Project onto segment from inside": {
			point:    NewPoint(5.5, 5.5),
			segment:  NewLineSegment[float64](NewPoint(2.0, 3.0), NewPoint(8.0, 7.0)),
			expected: NewPoint[float64](5.5772, 5.3848), // Should project onto the line segment itself
		},
		"float64: Project off the start of segment": {
			point:    NewPoint(0.0, 5.0),
			segment:  NewLineSegment[float64](NewPoint(2.0, 3.0), NewPoint(8.0, 7.0)),
			expected: NewPoint[float64](2.0, 3.0), // Should return a point of the segment
		},
		"float64: Project off the end of segment": {
			point:    NewPoint(10.0, 5.0),
			segment:  NewLineSegment[float64](NewPoint(2.0, 3.0), NewPoint(8.0, 7.0)),
			expected: NewPoint[float64](8.0, 7.0), // Should return end point of the segment
		},

		// Zero-length segment test cases
		"int: Project onto zero-length segment": {
			point:    NewPoint(3, 4),
			segment:  NewLineSegment[int](NewPoint(2, 2), NewPoint(2, 2)), // Zero-length segment
			expected: NewPoint[float64](2, 2),                             // Should return point A (or End), since the segment is a single point
		},
		"float64: Project onto zero-length segment": {
			point:    NewPoint(5.0, 5.0),
			segment:  NewLineSegment[float64](NewPoint(2.5, 2.5), NewPoint(2.5, 2.5)), // Zero-length segment
			expected: NewPoint[float64](2.5, 2.5),                                     // Should return point A (or End), since the segment is a single point
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			switch p := tt.point.(type) {
			case Point[int]:
				seg := tt.segment.(LineSegment[int])
				actual := p.ProjectOntoLineSegment(seg)
				assert.InDelta(t, tt.expected.x, actual.x, 0.001)
				assert.InDelta(t, tt.expected.y, actual.y, 0.001)

			case Point[float64]:
				seg := tt.segment.(LineSegment[float64])
				actual := p.ProjectOntoLineSegment(seg)
				assert.InDelta(t, tt.expected.x, actual.x, 0.001)
				assert.InDelta(t, tt.expected.y, actual.y, 0.001)
			}
		})
	}
}

func TestPoint_Reflect(t *testing.T) {
	tests := map[string]struct {
		point    Point[float64]       // The point to reflect
		axis     ReflectionAxis       // The axis or line type for reflection
		line     LineSegment[float64] // Custom line segment for ReflectAcrossCustomLine reflection
		expected Point[float64]       // Expected reflected point
	}{
		"reflect across x-axis": {
			point:    NewPoint[float64](3, 4),
			axis:     ReflectAcrossXAxis,
			expected: NewPoint[float64](3, -4),
		},
		"reflect across y-axis": {
			point:    NewPoint[float64](3, 4),
			axis:     ReflectAcrossYAxis,
			expected: NewPoint[float64](-3, 4),
		},
		"reflect across y = x line (ReflectAcrossCustomLine)": {
			point:    NewPoint[float64](3, 4),
			axis:     ReflectAcrossCustomLine,
			line:     NewLineSegment[float64](NewPoint[float64](0, 0), NewPoint[float64](1, 1)),
			expected: NewPoint[float64](4, 3),
		},
		"reflect across y = -x line (ReflectAcrossCustomLine)": {
			point:    NewPoint[float64](3, 4),
			axis:     ReflectAcrossCustomLine,
			line:     NewLineSegment[float64](NewPoint[float64](0, 0), NewPoint[float64](-1, 1)),
			expected: NewPoint[float64](-4, -3),
		},
		"reflect across degenerate line segment": {
			point:    NewPoint[float64](3, 4),
			axis:     ReflectAcrossCustomLine,
			line:     NewLineSegment[float64](NewPoint[float64](1, 1), NewPoint[float64](1, 1)), // Degenerate line
			expected: NewPoint[float64](3, 4),                                                   // Expect the point to remain unchanged
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var result Point[float64]
			if tt.axis == ReflectAcrossCustomLine {
				result = tt.point.Reflect(ReflectAcrossCustomLine, tt.line)
			} else {
				result = tt.point.Reflect(tt.axis)
			}

			assert.InDelta(t, tt.expected.x, result.x, 0.001)
			assert.InDelta(t, tt.expected.y, result.y, 0.001)
		})
	}

	t.Run("ReflectAcrossCustomLine with no given line segment", func(t *testing.T) {
		point := NewPoint[float64](3, 4)
		result := point.Reflect(ReflectAcrossCustomLine)
		assert.Equal(t, point, result, 0.001)
	})

	t.Run("invalid ReflectionAxis", func(t *testing.T) {
		point := NewPoint[float64](3, 4)
		result := point.Reflect(-1)
		assert.Equal(t, point, result, 0.001)
	})
}

func TestPoint_RelationshipToCircle(t *testing.T) {
	testCases := map[string]struct {
		point       Point[float64]
		circle      Circle[float64]
		expectedRel Relationship
	}{
		"Point inside circle": {
			point:       NewPoint[float64](2, 2),
			circle:      NewCircle(NewPoint[float64](0, 0), 5),
			expectedRel: RelationshipContainedBy,
		},
		"Point on circle boundary": {
			point:       NewPoint[float64](3, 4),
			circle:      NewCircle(NewPoint[float64](0, 0), 5),
			expectedRel: RelationshipIntersection,
		},
		"Point outside circle": {
			point:       NewPoint[float64](6, 8),
			circle:      NewCircle(NewPoint[float64](0, 0), 5),
			expectedRel: RelationshipDisjoint,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := tc.point.RelationshipToCircle(tc.circle, WithEpsilon(1e-10))
			assert.Equal(t, tc.expectedRel, result, "unexpected relationship")
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

func TestPoint_RelationshipToRectangle(t *testing.T) {
	rect := newRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 10))

	tests := map[string]struct {
		point       Point[int]
		expectedRel Relationship
	}{
		"Point inside rectangle": {
			point:       NewPoint(5, 5),
			expectedRel: RelationshipContainedBy,
		},
		"Point on rectangle edge": {
			point:       NewPoint(10, 5),
			expectedRel: RelationshipIntersection,
		},
		"Point outside rectangle": {
			point:       NewPoint(15, 5),
			expectedRel: RelationshipDisjoint,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			rel := tt.point.RelationshipToRectangle(rect, WithEpsilon(1e-10))
			assert.Equal(t, tt.expectedRel, rel, "unexpected relationship")
		})
	}
}

func TestConvexHull(t *testing.T) {
	tests := []struct {
		name     string
		points   any // Supports both int and float64 points
		expected any // Expected convex hull points
	}{
		// Integer points test cases
		{
			name: "int: Convex hull with multiple interior points",
			points: []Point[int]{
				{1, 4}, {4, 13}, {8, 17}, {18, 20}, {33, 18}, {38, 11}, {34, -2}, {21, -3}, {6, -1},
				{7, 6}, {10, 14}, {5, 2}, {16, 0}, {12, 12}, {23, 16}, {14, 6}, {24, 0}, {21, -2},
				{30, 14}, {27, 9}, {29, 4}, {31, -1}, {34, 7}, {35, 12}, {26, 2}, {25, 1}, {27, 3},
			},
			expected: []Point[int]{
				{21, -3}, {34, -2}, {38, 11}, {33, 18}, {18, 20}, {8, 17}, {4, 13}, {1, 4}, {6, -1},
			},
		},
		{
			name: "int: Simple square convex hull with all points outside",
			points: []Point[int]{
				{0, 0}, {20, 0}, {20, 20}, {0, 20},
				{19, 1}, {18, 3}, {17, 4}, {16, 4}, {15, 3}, {14, 2}, {13, 2}, {12, 2},
			},
			expected: []Point[int]{
				{0, 0}, {20, 0}, {20, 20}, {0, 20},
			},
		},

		// Float64 points test cases
		{
			name: "float64: Convex hull with multiple interior points",
			points: []Point[float64]{
				{1.0, 4.0}, {4.0, 13.0}, {8.0, 17.0}, {18.0, 20.0}, {33.0, 18.0}, {38.0, 11.0},
				{34.0, -2.0}, {21.0, -3.0}, {6.0, -1.0}, {7.0, 6.0}, {10.0, 14.0}, {5.0, 2.0},
				{16.0, 0.0}, {12.0, 12.0}, {23.0, 16.0}, {14.0, 6.0}, {24.0, 0.0}, {21.0, -2.0},
				{30.0, 14.0}, {27.0, 9.0}, {29.0, 4.0}, {31.0, -1.0}, {34.0, 7.0}, {35.0, 12.0},
				{26.0, 2.0}, {25.0, 1.0}, {27.0, 3.0},
			},
			expected: []Point[float64]{
				{21.0, -3.0}, {34.0, -2.0}, {38.0, 11.0}, {33.0, 18.0}, {18.0, 20.0}, {8.0, 17.0},
				{4.0, 13.0}, {1.0, 4.0}, {6.0, -1.0},
			},
		},
		{
			name: "float64: Simple square convex hull with all points outside",
			points: []Point[float64]{
				{0.0, 0.0}, {20.0, 0.0}, {20.0, 20.0}, {0.0, 20.0},
				{19.0, 1.0}, {18.0, 3.0}, {17.0, 4.0}, {16.0, 4.0}, {15.0, 3.0}, {14.0, 2.0},
				{13.0, 2.0}, {12.0, 2.0},
			},
			expected: []Point[float64]{
				{0.0, 0.0}, {20.0, 0.0}, {20.0, 20.0}, {0.0, 20.0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch points := tt.points.(type) {
			case []Point[int]:
				expected := tt.expected.([]Point[int])
				actual := ConvexHull(points)
				assert.Equal(t, expected, actual)

			case []Point[float64]:
				expected := tt.expected.([]Point[float64])
				actual := ConvexHull(points)
				assert.Equal(t, expected, actual)
			}
		})
	}
}

func TestEnsureClockwise(t *testing.T) {
	tests := map[string]struct {
		points   []Point[int]
		expected []Point[int]
	}{
		"already clockwise": {
			points:   []Point[int]{NewPoint(0, 0), NewPoint(2, 3), NewPoint(4, 0)},
			expected: []Point[int]{NewPoint(0, 0), NewPoint(2, 3), NewPoint(4, 0)},
		},
		"counterclockwise input": {
			points:   []Point[int]{NewPoint(0, 0), NewPoint(4, 0), NewPoint(2, 3)},
			expected: []Point[int]{NewPoint(2, 3), NewPoint(4, 0), NewPoint(0, 0)},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			EnsureClockwise(tt.points)
			assert.Equal(t, tt.expected, tt.points)
		})
	}
}

func TestEnsureCounterClockwise(t *testing.T) {
	tests := map[string]struct {
		points   []Point[int]
		expected []Point[int]
	}{
		"already counterclockwise": {
			points:   []Point[int]{NewPoint(0, 0), NewPoint(4, 0), NewPoint(2, 3)},
			expected: []Point[int]{NewPoint(0, 0), NewPoint(4, 0), NewPoint(2, 3)},
		},
		"clockwise input": {
			points:   []Point[int]{NewPoint(0, 0), NewPoint(2, 3), NewPoint(4, 0)},
			expected: []Point[int]{NewPoint(4, 0), NewPoint(2, 3), NewPoint(0, 0)},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			EnsureCounterClockwise(tt.points)
			assert.Equal(t, tt.expected, tt.points)
		})
	}
}

func TestFindLowestPoint(t *testing.T) {
	tests := []struct {
		name      string
		points    any // Supports both int and float64 points
		expectedI int // Expected index of the lowest point
		expectedP any // Expected lowest point as either int or float64
	}{
		// Integer points test case
		{
			name: "int: Multiple points, lowest at (0,0)",
			points: []Point[int]{
				{10, 10}, {10, 0}, {0, 0},
			},
			expectedI: 2,
			expectedP: Point[int]{0, 0},
		},

		// Float64 points test case
		{
			name: "float64: Multiple points, lowest at (0.0,0.0)",
			points: []Point[float64]{
				{10.5, 10.5}, {10.0, 0.0}, {0.0, 0.0},
			},
			expectedI: 2,
			expectedP: Point[float64]{0.0, 0.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch points := tt.points.(type) {
			case []Point[int]:
				expectedP := tt.expectedP.(Point[int])
				actualI, actualP := findLowestLeftestPoint(points...)
				assert.Equal(t, tt.expectedI, actualI)
				assert.Equal(t, expectedP, actualP)

			case []Point[float64]:
				expectedP := tt.expectedP.(Point[float64])
				actualI, actualP := findLowestLeftestPoint(points...)
				assert.Equal(t, tt.expectedI, actualI)
				assert.Equal(t, expectedP, actualP)
			}
		})
	}
}

func TestIsWellFormedPolygon(t *testing.T) {
	tests := []struct {
		name      string
		points    []Point[int]
		expected  bool
		errSubstr string // Substring expected in the error message
	}{
		{
			name: "Valid triangle",
			points: []Point[int]{
				NewPoint(0, 0),
				NewPoint(10, 0),
				NewPoint(5, 5),
			},
			expected: true,
		},
		{
			name: "Too few points",
			points: []Point[int]{
				NewPoint(0, 0),
				NewPoint(10, 0),
			},
			expected:  false,
			errSubstr: "at least 3 points",
		},
		{
			name: "Zero area (collinear points)",
			points: []Point[int]{
				NewPoint(0, 0),
				NewPoint(5, 0),
				NewPoint(10, 0),
			},
			expected:  false,
			errSubstr: "zero area",
		},
		{
			name: "Self-intersecting polygon",
			points: []Point[int]{
				NewPoint(0, 0),
				NewPoint(10, 10),
				NewPoint(10, 0),
				NewPoint(0, 2),
			},
			expected:  false,
			errSubstr: "self-intersecting",
		},
		{
			name: "Valid large polygon",
			points: []Point[int]{
				NewPoint(0, 0),
				NewPoint(10, 0),
				NewPoint(10, 10),
				NewPoint(0, 10),
			},
			expected: true,
		},
		{
			name: "Polygon with duplicate points",
			points: []Point[int]{
				NewPoint(0, 0),
				NewPoint(5, 0),
				NewPoint(5, 5),
				NewPoint(0, 5),
				NewPoint(0, 0), // Duplicate
			},
			expected:  false,
			errSubstr: "self-intersecting",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := IsWellFormedPolygon(tc.points)

			assert.Equal(t, tc.expected, result)

			if tc.errSubstr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errSubstr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTriangleAreaX2Signed(t *testing.T) {
	tests := []struct {
		name       string
		p0, p1, p2 any // Supports both int and float64 points
		expected   any // Expected result as either int or float64
	}{
		// Integer points test cases
		{
			name:     "int: (0,0), (10,10), (10,0)",
			p0:       NewPoint(0, 0),
			p1:       NewPoint(10, 10),
			p2:       NewPoint(10, 0),
			expected: -100,
		},
		{
			name:     "int: (0,0), (10,0), (10,10)",
			p0:       NewPoint(0, 0),
			p1:       NewPoint(10, 0),
			p2:       NewPoint(10, 10),
			expected: 100,
		},

		// Float64 points test cases
		{
			name:     "float64: (0.0,0.0), (10.0,10.0), (10.0,0.0)",
			p0:       NewPoint(0.0, 0.0),
			p1:       NewPoint(10.0, 10.0),
			p2:       NewPoint(10.0, 0.0),
			expected: -100.0,
		},
		{
			name:     "float64: (0.0,0.0), (10.0,0.0), (10.0,10.0)",
			p0:       NewPoint(0.0, 0.0),
			p1:       NewPoint(10.0, 0.0),
			p2:       NewPoint(10.0, 10.0),
			expected: 100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch p0 := tt.p0.(type) {
			case Point[int]:
				p1 := tt.p1.(Point[int])
				p2 := tt.p2.(Point[int])
				expected := tt.expected.(int)
				actual := triangleAreaX2Signed(p0, p1, p2)
				assert.Equal(t, expected, actual)

			case Point[float64]:
				p1 := tt.p1.(Point[float64])
				p2 := tt.p2.(Point[float64])
				expected := tt.expected.(float64)
				actual := triangleAreaX2Signed(p0, p1, p2)
				assert.Equal(t, expected, actual)
			}
		})
	}
}

func TestOrderPointsByAngleAboutLowestPoint(t *testing.T) {
	tests := map[string]struct {
		points, expected any
	}{
		"int: non-collinear points": {
			points: []Point[int]{ // input points are expected points, but have been randomized.
				{-10, 10}, {-8, 10}, {10, 9}, {-10, 0}, {-2, 10}, {10, 4},
				{-9, 10}, {-10, 6}, {-10, 3}, {-10, 1}, {-6, 10}, {-10, 7},
				{-4, 10}, {-10, 4}, {6, 10}, {10, 5}, {10, 3}, {-5, 10},
				{10, 1}, {-7, 10}, {3, 10}, {-10, 10}, {8, 10}, {-10, 8},
				{-3, 10}, {7, 10}, {10, 2}, {-10, 9}, {0, 10}, {9, 10},
				{10, 10}, {2, 10}, {0, -1}, {-10, 5}, {4, 10}, {-10, 2},
				{10, 0}, {5, 10}, {10, 6}, {10, 8}, {10, 7}, {-1, 10},
				{1, 10},
			},
			expected: []Point[int]{
				{0, -1}, // lowest point
				{10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {10, 5}, {10, 6},
				{10, 7}, {10, 8}, {10, 9}, {10, 10}, {9, 10}, {8, 10}, {7, 10},
				{6, 10}, {5, 10}, {4, 10}, {3, 10}, {2, 10}, {1, 10}, {0, 10},
				{-1, 10}, {-2, 10}, {-3, 10}, {-4, 10}, {-5, 10}, {-6, 10}, {-7, 10},
				{-8, 10}, {-9, 10}, {-10, 10}, {-10, 10}, {-10, 9}, {-10, 8}, {-10, 7},
				{-10, 6}, {-10, 5}, {-10, 4}, {-10, 3}, {-10, 2}, {-10, 1}, {-10, 0},
			},
		},
		"int: collinear points": {
			points: []Point[int]{ // input points are expected points, but have been randomized.
				{1, 8}, {-3, 8}, {-7, 7}, {3, 3}, {-6, 8}, {-3, 3}, {7, 8},
				{7, 7}, {5, 8}, {4, 8}, {0, 0}, {3, 8}, {4, 4}, {0, 8},
				{8, 8}, {-4, 4}, {-5, 5}, {-5, 8}, {6, 6}, {5, 5}, {6, 8},
				{1, 1}, {-2, 8}, {2, 8}, {-7, 8}, {2, 2}, {-1, 1}, {-2, 2},
				{-1, 8}, {-8, 8}, {-6, 6}, {-4, 8}},

			expected: []Point[int]{
				{0, 0}, // lowest point
				{1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5}, {6, 6}, {7, 7}, {8, 8},
				{7, 8}, {6, 8}, {5, 8}, {4, 8}, {3, 8}, {2, 8}, {1, 8},
				{0, 8}, {-1, 8}, {-2, 8}, {-3, 8}, {-4, 8}, {-5, 8},
				{-6, 8}, {-7, 8}, {-1, 1}, {-2, 2}, {-3, 3}, {-4, 4},
				{-5, 5}, {-6, 6}, {-7, 7}, {-8, 8},
			},
		},
		"float64: non-collinear points": {
			points: []Point[float64]{ // input points are expected points, but have been randomized.
				{-10, 10}, {-8, 10}, {10, 9}, {-10, 0}, {-2, 10}, {10, 4},
				{-9, 10}, {-10, 6}, {-10, 3}, {-10, 1}, {-6, 10}, {-10, 7},
				{-4, 10}, {-10, 4}, {6, 10}, {10, 5}, {10, 3}, {-5, 10},
				{10, 1}, {-7, 10}, {3, 10}, {-10, 10}, {8, 10}, {-10, 8},
				{-3, 10}, {7, 10}, {10, 2}, {-10, 9}, {0, 10}, {9, 10},
				{10, 10}, {2, 10}, {0, -1}, {-10, 5}, {4, 10}, {-10, 2},
				{10, 0}, {5, 10}, {10, 6}, {10, 8}, {10, 7}, {-1, 10},
				{1, 10},
			},
			expected: []Point[float64]{
				{0, -1}, // lowest point
				{10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {10, 5}, {10, 6},
				{10, 7}, {10, 8}, {10, 9}, {10, 10}, {9, 10}, {8, 10}, {7, 10},
				{6, 10}, {5, 10}, {4, 10}, {3, 10}, {2, 10}, {1, 10}, {0, 10},
				{-1, 10}, {-2, 10}, {-3, 10}, {-4, 10}, {-5, 10}, {-6, 10}, {-7, 10},
				{-8, 10}, {-9, 10}, {-10, 10}, {-10, 10}, {-10, 9}, {-10, 8}, {-10, 7},
				{-10, 6}, {-10, 5}, {-10, 4}, {-10, 3}, {-10, 2}, {-10, 1}, {-10, 0},
			},
		},
		"float64: collinear points": {
			points: []Point[float64]{ // input points are expected points, but have been randomized.
				{1, 8}, {-3, 8}, {-7, 7}, {3, 3}, {-6, 8}, {-3, 3}, {7, 8},
				{7, 7}, {5, 8}, {4, 8}, {0, 0}, {3, 8}, {4, 4}, {0, 8},
				{8, 8}, {-4, 4}, {-5, 5}, {-5, 8}, {6, 6}, {5, 5}, {6, 8},
				{1, 1}, {-2, 8}, {2, 8}, {-7, 8}, {2, 2}, {-1, 1}, {-2, 2},
				{-1, 8}, {-8, 8}, {-6, 6}, {-4, 8}},

			expected: []Point[float64]{
				{0, 0}, // lowest point
				{1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5}, {6, 6}, {7, 7}, {8, 8},
				{7, 8}, {6, 8}, {5, 8}, {4, 8}, {3, 8}, {2, 8}, {1, 8},
				{0, 8}, {-1, 8}, {-2, 8}, {-3, 8}, {-4, 8}, {-5, 8},
				{-6, 8}, {-7, 8}, {-1, 1}, {-2, 2}, {-3, 3}, {-4, 4},
				{-5, 5}, {-6, 6}, {-7, 7}, {-8, 8},
			},
		},
	}
	for name, tt := range tests {

		t.Run(name, func(t *testing.T) {
			switch points := tt.points.(type) {
			case []Point[int]:
				expected := tt.expected.([]Point[int])
				orderPointsByAngleAboutLowestPoint(expected[0], points)
				assert.Equal(t, tt.expected, tt.points)

			case []Point[float64]:
				expected := tt.expected.([]Point[float64])
				orderPointsByAngleAboutLowestPoint(expected[0], points)
				assert.Equal(t, tt.expected, tt.points)
			}
		})
	}
}

func TestOrientation(t *testing.T) {
	tests := map[string]struct {
		p0, p1, p2 any
		expected   PointOrientation
	}{
		"int: (0,0), (10,10), (10,0)": {
			p0:       NewPoint[int](0, 0),
			p1:       NewPoint[int](10, 10),
			p2:       NewPoint[int](10, 0),
			expected: PointsClockwise,
		},
		"int: (0,0), (10,0), (10,10)": {
			p0:       NewPoint[int](0, 0),
			p1:       NewPoint[int](10, 0),
			p2:       NewPoint[int](10, 10),
			expected: PointsCounterClockwise,
		},
		"int: (0,0), (10,0), (20,0)": {
			p0:       NewPoint[int](0, 0),
			p1:       NewPoint[int](10, 0),
			p2:       NewPoint[int](20, 0),
			expected: PointsCollinear,
		},
		"float64: (0,0), (10,10), (10,0)": {
			p0:       NewPoint[float64](0, 0),
			p1:       NewPoint[float64](10, 10),
			p2:       NewPoint[float64](10, 0),
			expected: PointsClockwise,
		},
		"float64: (0,0), (10,0), (10,10)": {
			p0:       NewPoint[float64](0, 0),
			p1:       NewPoint[float64](10, 0),
			p2:       NewPoint[float64](10, 10),
			expected: PointsCounterClockwise,
		},
		"float64: (0,0), (10,0), (20,0)": {
			p0:       NewPoint[float64](0, 0),
			p1:       NewPoint[float64](10, 0),
			p2:       NewPoint[float64](20, 0),
			expected: PointsCollinear,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			switch p0 := tt.p0.(type) {
			case Point[int]:
				p1 := tt.p1.(Point[int])
				p2 := tt.p2.(Point[int])
				actual := Orientation(p0, p1, p2)
				assert.Equal(t, tt.expected, actual)

			case Point[float64]:
				p1 := tt.p1.(Point[float64])
				p2 := tt.p2.(Point[float64])
				actual := Orientation(p0, p1, p2)
				assert.Equal(t, tt.expected, actual)
			}
		})
	}
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
