package geom2d

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func TestSimpleConvexPolygon_ContainsPoint(t *testing.T) {
	tests := map[string]struct {
		polygon         simpleConvexPolygon[int]
		point           Point[int]
		expectedContain bool
	}{
		"Point RelationshipPointCircleContainedByCircle Polygon": {
			polygon: simpleConvexPolygon[int]{
				Points: []Point[int]{
					NewPoint(0, 0),
					NewPoint(4, 0),
					NewPoint(4, 4),
					NewPoint(0, 4),
				},
			},
			point:           NewPoint(2, 2), // Point inside the square
			expectedContain: true,
		},
		"Point RelationshipPointCircleMiss Polygon": {
			polygon: simpleConvexPolygon[int]{
				Points: []Point[int]{
					NewPoint(0, 0),
					NewPoint(4, 0),
					NewPoint(4, 4),
					NewPoint(0, 4),
				},
			},
			point:           NewPoint(5, 5), // Point outside the square
			expectedContain: false,
		},
		"Point On Edge": {
			polygon: simpleConvexPolygon[int]{
				Points: []Point[int]{
					NewPoint(0, 0),
					NewPoint(4, 0),
					NewPoint(4, 4),
					NewPoint(0, 4),
				},
			},
			point:           NewPoint(4, 2), // Point on the edge of the square
			expectedContain: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			contains := tc.polygon.ContainsPoint(tc.point)
			assert.Equal(t, tc.expectedContain, contains, "Expected ContainsPoint result does not match")
		})
	}
}

func TestContour_contains(t *testing.T) {
	tests := map[string]struct {
		contour  contour[int]
		point    Point[int]
		expected bool
	}{
		"point in contour": {
			contour: contour[int]{
				{point: Point[int]{0, 0}},
				{point: Point[int]{1, 2}},
				{point: Point[int]{3, 4}},
			},
			point:    Point[int]{1, 2},
			expected: true,
		},
		"point not in contour": {
			contour: contour[int]{
				{point: Point[int]{0, 0}},
				{point: Point[int]{1, 2}},
				{point: Point[int]{3, 4}},
			},
			point:    Point[int]{5, 6},
			expected: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tc.contour.contains(tc.point)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestContour_EnsureClockwise(t *testing.T) {
	tests := map[string]struct {
		contour        contour[int]
		expectedPoints []Point[int]
	}{
		"Clockwise": {
			contour: contour[int]{
				{point: NewPoint(0, 0)},
				{point: NewPoint(0, 10)},
				{point: NewPoint(10, 0)},
			},
			expectedPoints: []Point[int]{
				NewPoint(0, 0),
				NewPoint(0, 10),
				NewPoint(10, 0),
			},
		},
		"CounterClockwise": {
			contour: contour[int]{
				{point: NewPoint(0, 0)},
				{point: NewPoint(10, 0)},
				{point: NewPoint(0, 10)},
			},
			expectedPoints: []Point[int]{
				NewPoint(0, 10),
				NewPoint(10, 0),
				NewPoint(0, 0),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.contour.ensureClockwise()
			for i, p := range tc.contour {
				assert.Equal(t, tc.expectedPoints[i], p.point, "contour points mismatch")
			}
		})
	}
}

func TestContour_EnsureCounterClockwise(t *testing.T) {
	tests := map[string]struct {
		contour        contour[int]
		expectedPoints []Point[int]
	}{
		"Clockwise": {
			contour: contour[int]{
				{point: NewPoint(0, 0)},
				{point: NewPoint(0, 10)},
				{point: NewPoint(10, 0)},
			},
			expectedPoints: []Point[int]{
				NewPoint(10, 0),
				NewPoint(0, 10),
				NewPoint(0, 0),
			},
		},
		"CounterClockwise": {
			contour: contour[int]{
				{point: NewPoint(0, 0)},
				{point: NewPoint(10, 0)},
				{point: NewPoint(0, 10)},
			},
			expectedPoints: []Point[int]{
				NewPoint(0, 0),
				NewPoint(10, 0),
				NewPoint(0, 10),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.contour.ensureCounterClockwise()
			for i, p := range tc.contour {
				assert.Equal(t, tc.expectedPoints[i], p.point, "contour points mismatch")
			}
		})
	}
}

func TestContour_Eq(t *testing.T) {
	tests := map[string]struct {
		contour1 contour[int]
		contour2 contour[int]
		expected bool
	}{
		"identical contours": {
			contour1: contour[int]{
				{point: Point[int]{0, 0}},
				{point: Point[int]{10, 0}},
				{point: Point[int]{10, 10}},
				{point: Point[int]{0, 10}},
			},
			contour2: contour[int]{
				{point: Point[int]{0, 0}},
				{point: Point[int]{10, 0}},
				{point: Point[int]{10, 10}},
				{point: Point[int]{0, 10}},
			},
			expected: true,
		},
		"offset contours": {
			contour1: contour[int]{
				{point: Point[int]{0, 0}},
				{point: Point[int]{10, 0}},
				{point: Point[int]{10, 10}},
				{point: Point[int]{0, 10}},
			},
			contour2: contour[int]{
				{point: Point[int]{10, 10}},
				{point: Point[int]{0, 10}},
				{point: Point[int]{0, 0}},
				{point: Point[int]{10, 0}},
			},
			expected: true,
		},
		"rotated contours": {
			contour1: contour[int]{
				{point: Point[int]{0, 0}},
				{point: Point[int]{10, 0}},
				{point: Point[int]{10, 10}},
				{point: Point[int]{0, 10}},
			},
			contour2: contour[int]{
				{point: Point[int]{10, 0}},
				{point: Point[int]{10, 10}},
				{point: Point[int]{0, 10}},
				{point: Point[int]{0, 0}},
			},
			expected: true,
		},
		"reversed contours": {
			contour1: contour[int]{
				{point: Point[int]{0, 0}},
				{point: Point[int]{10, 0}},
				{point: Point[int]{10, 10}},
				{point: Point[int]{0, 10}},
			},
			contour2: contour[int]{
				{point: Point[int]{0, 10}},
				{point: Point[int]{10, 10}},
				{point: Point[int]{10, 0}},
				{point: Point[int]{0, 0}},
			},
			expected: true,
		},
		"offset + reversed contours": {
			contour1: contour[int]{
				{point: Point[int]{0, 0}},
				{point: Point[int]{10, 0}},
				{point: Point[int]{10, 10}},
				{point: Point[int]{0, 10}},
			},
			contour2: contour[int]{
				{point: Point[int]{10, 0}},
				{point: Point[int]{0, 0}},
				{point: Point[int]{0, 10}},
				{point: Point[int]{10, 10}},
			},
			expected: true,
		},
		"empty contours": {
			contour1: contour[int]{},
			contour2: contour[int]{},
			expected: true,
		},
		"one empty contour": {
			contour1: contour[int]{
				{point: Point[int]{0, 0}},
				{point: Point[int]{10, 0}},
				{point: Point[int]{10, 10}},
				{point: Point[int]{0, 10}},
			},
			contour2: contour[int]{},
			expected: false,
		},
		"mismatched lengths": {
			contour1: contour[int]{
				{point: Point[int]{0, 0}},
				{point: Point[int]{10, 0}},
				{point: Point[int]{10, 10}},
				{point: Point[int]{0, 10}},
			},
			contour2: contour[int]{
				{point: Point[int]{0, 0}},
				{point: Point[int]{10, 0}},
				{point: Point[int]{10, 10}},
			},
			expected: false,
		},
		"mismatched points": {
			contour1: contour[int]{
				{point: Point[int]{0, 0}},
				{point: Point[int]{10, 0}},
				{point: Point[int]{10, 10}},
				{point: Point[int]{0, 10}},
			},
			contour2: contour[int]{
				{point: Point[int]{5, 5}},
				{point: Point[int]{15, 5}},
				{point: Point[int]{15, 15}},
				{point: Point[int]{5, 15}},
			},
			expected: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.contour1.eq(tc.contour2))
		})
	}
}

func TestContour_findLowestLeftmost(t *testing.T) {
	tests := map[string]struct {
		contourPoints []Point[int]
		expectedPoint Point[int]
	}{
		"Single point": {
			contourPoints: []Point[int]{
				NewPoint(5, 5),
			},
			expectedPoint: NewPoint(5, 5),
		},
		"Multiple points with distinct y": {
			contourPoints: []Point[int]{
				NewPoint(10, 10),
				NewPoint(5, 5),
				NewPoint(15, 20),
			},
			expectedPoint: NewPoint(5, 5),
		},
		"Multiple points with same y": {
			contourPoints: []Point[int]{
				NewPoint(10, 10),
				NewPoint(5, 10),
				NewPoint(15, 10),
			},
			expectedPoint: NewPoint(5, 10),
		},
		"Negative coordinates": {
			contourPoints: []Point[int]{
				NewPoint(-5, -10),
				NewPoint(-10, -5),
				NewPoint(-15, -20),
			},
			expectedPoint: NewPoint(-15, -20),
		},
		"Duplicate lowest-leftmost points": {
			contourPoints: []Point[int]{
				NewPoint(5, 5),
				NewPoint(5, 5),
				NewPoint(10, 10),
			},
			expectedPoint: NewPoint(5, 5),
		},
		"Points with mixed y and x values": {
			contourPoints: []Point[int]{
				NewPoint(10, 10),
				NewPoint(10, 5),
				NewPoint(5, 5),
			},
			expectedPoint: NewPoint(5, 5),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			contour := contour[int]{}
			for _, pt := range tt.contourPoints {
				contour = append(contour, polyTreePoint[int]{point: pt})
			}

			result := contour.findLowestLeftmost()
			assert.Equal(t, tt.expectedPoint, result, "unexpected lowest-leftmost point")
		})
	}
}

func TestContour_insertIntersectionPoint(t *testing.T) {
	tests := map[string]struct {
		contour       []polyTreePoint[int]
		start, end    int
		intersection  polyTreePoint[int]
		expectedOrder []polyTreePoint[int]
	}{
		"Insert at middle": {
			contour: []polyTreePoint[int]{
				{point: NewPoint(0, 0)},
				{point: NewPoint(10, 10)},
			},
			start:        0,
			end:          1,
			intersection: polyTreePoint[int]{point: NewPoint(5, 5)},
			expectedOrder: []polyTreePoint[int]{
				{point: NewPoint(0, 0)},
				{point: NewPoint(5, 5)},
				{point: NewPoint(10, 10)},
			},
		},
		"Insert closer to start": {
			contour: []polyTreePoint[int]{
				{point: NewPoint(0, 0)},
				{point: NewPoint(10, 10)},
				{point: NewPoint(20, 20)},
			},
			start:        0,
			end:          1,
			intersection: polyTreePoint[int]{point: NewPoint(5, 5)},
			expectedOrder: []polyTreePoint[int]{
				{point: NewPoint(0, 0)},
				{point: NewPoint(5, 5)},
				{point: NewPoint(10, 10)},
				{point: NewPoint(20, 20)},
			},
		},
		"Insert multiple intermediate points": {
			contour: []polyTreePoint[int]{
				{point: NewPoint(0, 0)},
				{point: NewPoint(20, 20)},
				{point: NewPoint(40, 40)},
			},
			start:        0,
			end:          1,
			intersection: polyTreePoint[int]{point: NewPoint(10, 10)},
			expectedOrder: []polyTreePoint[int]{
				{point: NewPoint(0, 0)},
				{point: NewPoint(10, 10)},
				{point: NewPoint(20, 20)},
				{point: NewPoint(40, 40)},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			contour := contour[int](tt.contour)
			contour.insertIntersectionPoint(tt.start, tt.end, tt.intersection)
			assert.Equal(t, tt.expectedOrder, []polyTreePoint[int](contour), "unexpected contour order after insertion")
		})
	}
}

func TestContour_isPointInside(t *testing.T) {
	tests := map[string]struct {
		c        contour[int]
		p        Point[int]
		expected bool
	}{
		"point inside box": {
			c: contour[int]{
				{point: Point[int]{0, 0}},
				{point: Point[int]{10, 0}},
				{point: Point[int]{10, 10}},
				{point: Point[int]{0, 10}},
			},
			p:        NewPoint(5, 5),
			expected: true,
		},
		"point outside box": {
			c: contour[int]{
				{point: Point[int]{0, 0}},
				{point: Point[int]{10, 0}},
				{point: Point[int]{10, 10}},
				{point: Point[int]{0, 10}},
			},
			p:        NewPoint(-5, 5),
			expected: false,
		},
		"point inside, collinear with contour edge facing right": {
			c: contour[int]{
				{point: Point[int]{0, 0}},
				{point: Point[int]{40, 0}},
				{point: Point[int]{40, 14}},
				{point: Point[int]{54, 14}},
				{point: Point[int]{54, 54}},
				{point: Point[int]{0, 54}},
			},
			p:        NewPoint(30, 14),
			expected: true,
		},
		"point inside, collinear with contour edge facing left": {
			c: contour[int]{
				{point: Point[int]{0, 0}},
				{point: Point[int]{54, 0}},
				{point: Point[int]{54, 14}},
				{point: Point[int]{40, 14}},
				{point: Point[int]{40, 54}},
				{point: Point[int]{0, 54}},
			},
			p:        NewPoint(30, 14),
			expected: true,
		},
		"point outside, collinear with contour edge facing right": {
			c: contour[int]{
				{point: Point[int]{0, 0}},
				{point: Point[int]{40, 0}},
				{point: Point[int]{40, 14}},
				{point: Point[int]{54, 14}},
				{point: Point[int]{54, 54}},
				{point: Point[int]{0, 54}},
			},
			p:        NewPoint(-30, 14),
			expected: false,
		},
		"point outside, collinear with contour edge facing left": {
			c: contour[int]{
				{point: Point[int]{0, 0}},
				{point: Point[int]{54, 0}},
				{point: Point[int]{54, 14}},
				{point: Point[int]{40, 14}},
				{point: Point[int]{40, 54}},
				{point: Point[int]{0, 54}},
			},
			p:        NewPoint(-30, 14),
			expected: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.c.isPointInside(tc.p))
		})
	}
}

func TestContour_iterEdges_EarlyExit(t *testing.T) {
	c := contour[int]{
		{point: Point[int]{0, 0}},
		{point: Point[int]{10, 0}},
		{point: Point[int]{10, 10}},
		{point: Point[int]{0, 10}},
	}
	count := 0
	c.iterEdges(func(edge LineSegment[int]) bool {
		count++
		return count < 2 // Exit after two edges
	})
	require.Equal(t, 2, count, "iterEdges should exit early when yield returns false")
}

func TestContour_iterEdges_Empty(t *testing.T) {
	c := contour[int]{}
	count := 0
	c.iterEdges(func(edge LineSegment[int]) bool {
		count++
		return true
	})
	require.Equal(t, 0, count, "iterEdges should not yield edges for an empty c")
}

func TestContour_iterEdges_FullPolygon(t *testing.T) {
	c := contour[int]{
		{point: Point[int]{0, 0}},
		{point: Point[int]{10, 0}},
		{point: Point[int]{10, 10}},
		{point: Point[int]{0, 10}},
	}
	var edges []LineSegment[int]
	c.iterEdges(func(edge LineSegment[int]) bool {
		edges = append(edges, edge)
		return true
	})
	require.Equal(t, 4, len(edges), "iterEdges should yield one edge per c segment")
	require.Equal(t, NewLineSegment(Point[int]{0, 0}, Point[int]{10, 0}), edges[0])
	require.Equal(t, NewLineSegment(Point[int]{10, 0}, Point[int]{10, 10}), edges[1])
	require.Equal(t, NewLineSegment(Point[int]{10, 10}, Point[int]{0, 10}), edges[2])
	require.Equal(t, NewLineSegment(Point[int]{0, 10}, Point[int]{0, 0}), edges[3])
}

func TestContour_iterEdges_TwoPoints(t *testing.T) {
	c := contour[int]{
		{point: Point[int]{0, 0}},
		{point: Point[int]{10, 0}},
	}
	var edges []LineSegment[int]
	c.iterEdges(func(edge LineSegment[int]) bool {
		edges = append(edges, edge)
		return true
	})
	require.Equal(t, 2, len(edges), "iterEdges should yield exactly two edges for a closed loop with two points")
	require.Equal(t, NewLineSegment(Point[int]{0, 0}, Point[int]{10, 0}), edges[0])
	require.Equal(t, NewLineSegment(Point[int]{10, 0}, Point[int]{0, 0}), edges[1])
}

func TestContour_iterEdges_Triangle(t *testing.T) {
	c := contour[int]{
		{point: Point[int]{0, 0}},
		{point: Point[int]{10, 0}},
		{point: Point[int]{5, 10}},
	}
	var edges []LineSegment[int]
	c.iterEdges(func(edge LineSegment[int]) bool {
		edges = append(edges, edge)
		return true
	})
	require.Equal(t, 3, len(edges), "iterEdges should yield exactly three edges for a triangle")
	require.Equal(t, NewLineSegment(Point[int]{0, 0}, Point[int]{10, 0}), edges[0])
	require.Equal(t, NewLineSegment(Point[int]{10, 0}, Point[int]{5, 10}), edges[1])
	require.Equal(t, NewLineSegment(Point[int]{5, 10}, Point[int]{0, 0}), edges[2])
}

func TestPolyIntersectionType_String(t *testing.T) {
	tests := map[string]struct {
		input          polyIntersectionType
		expectedOutput string
		shouldPanic    bool
	}{
		"NotSet": {
			input:          intersectionTypeNotSet,
			expectedOutput: "not set",
			shouldPanic:    false,
		},
		"Entry": {
			input:          intersectionTypeEntry,
			expectedOutput: "entry",
			shouldPanic:    false,
		},
		"Exit": {
			input:          intersectionTypeExit,
			expectedOutput: "exit",
			shouldPanic:    false,
		},
		"UnsupportedType": {
			input:       polyIntersectionType(999),
			shouldPanic: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if test.shouldPanic {
				assert.Panics(t, func() {
					_ = test.input.String()
				}, "Expected panic for unsupported polyIntersectionType")
			} else {
				assert.Equal(t, test.expectedOutput, test.input.String())
			}
		})
	}
}

func TestPolygonType_String(t *testing.T) {
	tests := map[string]struct {
		input    PolygonType
		expected string
		panics   bool
	}{
		"PTSolid": {
			input:    PTSolid,
			expected: "PTSolid",
			panics:   false,
		},
		"PTHole": {
			input:    PTHole,
			expected: "PTHole",
			panics:   false,
		},
		"Unsupported PolygonType": {
			input:  PolygonType(5), // Arbitrary unsupported value
			panics: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.panics {
				require.Panics(t, func() {
					_ = tt.input.String()
				}, "Expected panic for unsupported PolygonType")
			} else {
				assert.Equal(t, tt.expected, tt.input.String(), "Unexpected string value for PolygonType")
			}
		})
	}
}

func TestPolyTree_AddChild(t *testing.T) {
	t.Run("Adding a nil child", func(t *testing.T) {
		parent, _ := NewPolyTree([]Point[int]{{0, 0}, {10, 0}, {10, 10}, {0, 10}}, PTSolid)
		err := parent.AddChild(nil)
		require.Error(t, err, "expected error when adding a nil child, but got none")
	})

	t.Run("Adding a Child with Opposite polygonType", func(t *testing.T) {
		parent, err := NewPolyTree([]Point[int]{{0, 0}, {10, 0}, {10, 10}, {0, 10}}, PTSolid)
		require.NoError(t, err, "error creating parent polygon, when none was expected")

		child, err := NewPolyTree([]Point[int]{{2, 2}, {8, 2}, {8, 8}, {2, 8}}, PTHole)
		require.NoError(t, err, "error creating child polygon, when none was expected")

		err = parent.AddChild(child)
		require.NoError(t, err, "error calling AddChild, when none was expected")

		assert.Contains(t, parent.children, child)
		assert.Equal(t, parent, child.parent)
	})

	t.Run("Adding a Child with Mismatched polygonType", func(t *testing.T) {
		parent, err := NewPolyTree([]Point[int]{{0, 0}, {10, 0}, {10, 10}, {0, 10}}, PTSolid)
		require.NoError(t, err, "error creating parent polygon, when none was expected")

		child, err := NewPolyTree([]Point[int]{{2, 2}, {8, 2}, {8, 8}, {2, 8}}, PTSolid)
		require.NoError(t, err, "error creating child polygon, when none was expected")

		err = parent.AddChild(child)
		require.Error(t, err, "no error returned from AddChild, when one was expected")
	})

	t.Run("Adding Multiple Children", func(t *testing.T) {
		parent, err := NewPolyTree([]Point[int]{{0, 0}, {10, 0}, {10, 10}, {0, 10}}, PTSolid)
		require.NoError(t, err, "error creating parent polygon, when none was expected")

		child1, err := NewPolyTree([]Point[int]{{2, 2}, {4, 2}, {4, 4}, {2, 4}}, PTHole)
		require.NoError(t, err, "error creating first child polygon, when none was expected")

		child2, err := NewPolyTree([]Point[int]{{6, 6}, {8, 6}, {8, 8}, {6, 8}}, PTHole)
		require.NoError(t, err, "error creating second child polygon, when none was expected")

		err = parent.AddChild(child1)
		require.NoError(t, err, "error calling AddChild for child1, when none was expected")
		err = parent.AddChild(child2)
		require.NoError(t, err, "error calling AddChild for child2, when none was expected")

		assert.Contains(t, parent.children, child1)
		assert.Contains(t, parent.children, child2)
		assert.Equal(t, parent, child1.parent)
		assert.Equal(t, parent, child2.parent)
	})
}

func TestPolyTree_AddSibling(t *testing.T) {
	t.Run("Adding a nil sibling", func(t *testing.T) {
		poly1, _ := NewPolyTree([]Point[int]{{0, 0}, {10, 0}, {10, 10}, {0, 10}}, PTSolid)
		err := poly1.AddSibling(nil)
		require.Error(t, err, "expected error when adding a nil sibling, but got none")
	})

	t.Run("Adding a Sibling with Matching polygonType", func(t *testing.T) {
		poly1, err := NewPolyTree([]Point[int]{{0, 0}, {10, 0}, {10, 10}, {0, 10}}, PTSolid)
		require.NoError(t, err, "error creating poly1, when none was expected")
		poly2, err := NewPolyTree([]Point[int]{{20, 20}, {30, 20}, {30, 30}, {20, 30}}, PTSolid)
		require.NoError(t, err, "error creating poly2, when none was expected")

		err = poly1.AddSibling(poly2)
		require.NoError(t, err, "error calling AddSibling, when none was expected")
		assert.Contains(t, poly1.siblings, poly2)
		assert.Contains(t, poly2.siblings, poly1)
	})

	t.Run("Adding a Sibling with Mismatched polygonType", func(t *testing.T) {
		poly1, err := NewPolyTree([]Point[int]{{0, 0}, {10, 0}, {10, 10}, {0, 10}}, PTSolid)
		require.NoError(t, err, "error creating poly1, when none was expected")
		poly2, err := NewPolyTree([]Point[int]{{20, 20}, {30, 20}, {30, 30}, {20, 30}}, PTHole)
		require.NoError(t, err, "error creating poly2, when none was expected")

		err = poly1.AddSibling(poly2)
		require.Error(t, err, "no error returned from AddSibling, when one was expected")
	})

	t.Run("Adding Multiple Siblings", func(t *testing.T) {
		poly1, err := NewPolyTree([]Point[int]{{0, 0}, {10, 0}, {10, 10}, {0, 10}}, PTSolid)
		require.NoError(t, err, "error creating poly1, when none was expected")
		poly2, err := NewPolyTree([]Point[int]{{20, 20}, {30, 20}, {30, 30}, {20, 30}}, PTSolid)
		require.NoError(t, err, "error creating poly2, when none was expected")
		poly3, err := NewPolyTree([]Point[int]{{40, 40}, {50, 40}, {50, 50}, {40, 50}}, PTSolid)
		require.NoError(t, err, "error creating poly3, when none was expected")

		err = poly1.AddSibling(poly2)
		require.NoError(t, err, "error returned from poly1.AddSibling(poly2) when none was expected")
		err = poly1.AddSibling(poly3)
		require.NoError(t, err, "error returned from poly1.AddSibling(poly3) when none was expected")

		assert.Contains(t, poly1.siblings, poly2)
		assert.Contains(t, poly1.siblings, poly3)
		assert.Contains(t, poly2.siblings, poly1)
		assert.Contains(t, poly2.siblings, poly3)
		assert.Contains(t, poly3.siblings, poly1)
		assert.Contains(t, poly3.siblings, poly2)
	})
}

func TestPolyTree_Area(t *testing.T) {
	tests := map[string]struct {
		contour       []Point[int]
		expectedArea  float64
		expectedError error
	}{
		"Square": {
			contour: []Point[int]{
				NewPoint(0, 0),
				NewPoint(0, 10),
				NewPoint(10, 10),
				NewPoint(10, 0),
			},
			expectedArea:  100.0,
			expectedError: nil,
		},
		"Triangle": {
			contour: []Point[int]{
				NewPoint(0, 0),
				NewPoint(5, 10),
				NewPoint(10, 0),
			},
			expectedArea:  50.0,
			expectedError: nil,
		},
		"Line (Degenerate Polygon)": {
			contour: []Point[int]{
				NewPoint(0, 0),
				NewPoint(10, 0),
			},
			expectedArea:  0.0,
			expectedError: fmt.Errorf("new polytree must have at least 3 points"),
		},
		"Point (Degenerate Polygon)": {
			contour: []Point[int]{
				NewPoint(0, 0),
			},
			expectedArea:  0.0,
			expectedError: fmt.Errorf("new polytree must have at least 3 points"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			pt, err := NewPolyTree(test.contour, PTSolid)
			if test.expectedError != nil {
				assert.Error(t, err, "Expected error")
				assert.Equal(t, test.expectedError, err)
				return
			}
			assert.NoError(t, err, "Unexpected error")

			area := pt.Area()
			assert.Equal(t, test.expectedArea, area, "Unexpected area")
		})
	}
}

func TestPolyTree_AsInt(t *testing.T) {
	root, err := NewPolyTree([]Point[float64]{
		NewPoint(0.1, 0.9),
		NewPoint(100.8, 0.8),
		NewPoint(100.7, 100.1),
		NewPoint(0.2, 100.2),
	}, PTSolid)
	require.NoError(t, err, "error defining root poly")

	hole, err := NewPolyTree([]Point[float64]{
		NewPoint(20.1, 20.8),
		NewPoint(80.2, 20.1),
		NewPoint(80.7, 80.6),
		NewPoint(20.7, 80.7),
	}, PTHole)
	require.NoError(t, err, "error defining hole poly")

	island, err := NewPolyTree([]Point[float64]{
		NewPoint(40.1, 40.8),
		NewPoint(60.2, 40.1),
		NewPoint(60.7, 60.6),
		NewPoint(40.7, 60.7),
	}, PTSolid)
	require.NoError(t, err, "error defining island poly")

	err = root.AddChild(hole)
	require.NoError(t, err, "error adding hole as child of root")

	err = hole.AddChild(island)
	require.NoError(t, err, "error adding island as child of hole")

	expectedRoot, err := NewPolyTree([]Point[int]{
		NewPoint(0, 0),
		NewPoint(100, 0),
		NewPoint(100, 100),
		NewPoint(0, 100),
	}, PTSolid)
	require.NoError(t, err, "error defining expectedRoot poly")

	expectedHole, err := NewPolyTree([]Point[int]{
		NewPoint(20, 20),
		NewPoint(80, 20),
		NewPoint(80, 80),
		NewPoint(20, 80),
	}, PTHole)
	require.NoError(t, err, "error defining expectedHole poly")

	expectedIsland, err := NewPolyTree([]Point[int]{
		NewPoint(40, 40),
		NewPoint(60, 40),
		NewPoint(60, 60),
		NewPoint(40, 60),
	}, PTSolid)
	require.NoError(t, err, "error defining expectedIsland poly")

	err = expectedRoot.AddChild(expectedHole)
	require.NoError(t, err, "error adding expectedHole as child of expectedRoot")

	err = expectedHole.AddChild(expectedIsland)
	require.NoError(t, err, "error adding expectedIsland as child of expectedHole")

	// Convert to int (truncated)
	rootInt := root.AsInt()

	assert.True(t, rootInt.contour.eq(expectedRoot.contour), "root contour does not match")
	assert.True(t, rootInt.children[0].contour.eq(expectedRoot.children[0].contour), "hole contour does not match")
	assert.True(t, rootInt.children[0].children[0].contour.eq(expectedRoot.children[0].children[0].contour), "hole contour does not match")

}

func TestPolyTree_AsIntRounded(t *testing.T) {
	root, err := NewPolyTree([]Point[float64]{
		NewPoint(0.1, 0.9),
		NewPoint(100.8, 0.8),
		NewPoint(100.7, 100.1),
		NewPoint(0.2, 100.2),
	}, PTSolid)
	require.NoError(t, err, "error defining root poly")

	hole, err := NewPolyTree([]Point[float64]{
		NewPoint(20.1, 20.8),
		NewPoint(80.2, 20.1),
		NewPoint(80.7, 80.6),
		NewPoint(20.7, 80.7),
	}, PTHole)
	require.NoError(t, err, "error defining hole poly")

	island, err := NewPolyTree([]Point[float64]{
		NewPoint(40.1, 40.8),
		NewPoint(60.2, 40.1),
		NewPoint(60.7, 60.6),
		NewPoint(40.7, 60.7),
	}, PTSolid)
	require.NoError(t, err, "error defining island poly")

	err = root.AddChild(hole)
	require.NoError(t, err, "error adding hole as child of root")

	err = hole.AddChild(island)
	require.NoError(t, err, "error adding island as child of hole")

	expectedRoot, err := NewPolyTree([]Point[int]{
		NewPoint(0, 1),
		NewPoint(101, 1),
		NewPoint(101, 100),
		NewPoint(0, 100),
	}, PTSolid)
	require.NoError(t, err, "error defining expectedRoot poly")

	expectedHole, err := NewPolyTree([]Point[int]{
		NewPoint(20, 21),
		NewPoint(80, 20),
		NewPoint(81, 81),
		NewPoint(21, 81),
	}, PTHole)
	require.NoError(t, err, "error defining expectedHole poly")

	expectedIsland, err := NewPolyTree([]Point[int]{
		NewPoint(40, 41),
		NewPoint(60, 40),
		NewPoint(61, 61),
		NewPoint(41, 61),
	}, PTSolid)
	require.NoError(t, err, "error defining expectedIsland poly")

	err = expectedRoot.AddChild(expectedHole)
	require.NoError(t, err, "error adding expectedHole as child of expectedRoot")

	err = expectedHole.AddChild(expectedIsland)
	require.NoError(t, err, "error adding expectedIsland as child of expectedHole")

	// Convert to int (truncated)
	rootIntRounded := root.AsIntRounded()

	assert.True(t, rootIntRounded.contour.eq(expectedRoot.contour), "root contour does not match")
	assert.True(t, rootIntRounded.children[0].contour.eq(expectedRoot.children[0].contour), "hole contour does not match")
	assert.True(t, rootIntRounded.children[0].children[0].contour.eq(expectedRoot.children[0].children[0].contour), "hole contour does not match")

}

func TestPolyTree_AsFloat32(t *testing.T) {
	root, err := NewPolyTree([]Point[int]{
		NewPoint(0, 0),
		NewPoint(100, 0),
		NewPoint(100, 100),
		NewPoint(0, 100),
	}, PTSolid)
	require.NoError(t, err, "error defining root poly")

	hole, err := NewPolyTree([]Point[int]{
		NewPoint(20, 20),
		NewPoint(80, 20),
		NewPoint(80, 80),
		NewPoint(20, 80),
	}, PTHole)
	require.NoError(t, err, "error defining hole poly")

	island, err := NewPolyTree([]Point[int]{
		NewPoint(40, 40),
		NewPoint(60, 40),
		NewPoint(60, 60),
		NewPoint(40, 60),
	}, PTSolid)
	require.NoError(t, err, "error defining island poly")

	err = root.AddChild(hole)
	require.NoError(t, err, "error adding hole as child of root")

	err = hole.AddChild(island)
	require.NoError(t, err, "error adding island as child of hole")

	expectedRoot, err := NewPolyTree([]Point[float32]{
		NewPoint[float32](0, 0),
		NewPoint[float32](100, 0),
		NewPoint[float32](100, 100),
		NewPoint[float32](0, 100),
	}, PTSolid)
	require.NoError(t, err, "error defining expectedRoot poly")

	expectedHole, err := NewPolyTree([]Point[float32]{
		NewPoint[float32](20, 20),
		NewPoint[float32](80, 20),
		NewPoint[float32](80, 80),
		NewPoint[float32](20, 80),
	}, PTHole)
	require.NoError(t, err, "error defining expectedHole poly")

	expectedIsland, err := NewPolyTree([]Point[float32]{
		NewPoint[float32](40, 40),
		NewPoint[float32](60, 40),
		NewPoint[float32](60, 60),
		NewPoint[float32](40, 60),
	}, PTSolid)
	require.NoError(t, err, "error defining expectedIsland poly")

	err = expectedRoot.AddChild(expectedHole)
	require.NoError(t, err, "error adding expectedHole as child of expectedRoot")

	err = expectedHole.AddChild(expectedIsland)
	require.NoError(t, err, "error adding expectedIsland as child of expectedHole")

	// Convert to int (truncated)
	rootIntRounded := root.AsFloat32()

	assert.True(t, rootIntRounded.contour.eq(expectedRoot.contour), "root contour does not match")
	assert.True(t, rootIntRounded.children[0].contour.eq(expectedRoot.children[0].contour), "hole contour does not match")
	assert.True(t, rootIntRounded.children[0].children[0].contour.eq(expectedRoot.children[0].children[0].contour), "hole contour does not match")

}

func TestPolyTree_AsFloat64(t *testing.T) {
	root, err := NewPolyTree([]Point[int]{
		NewPoint(0, 0),
		NewPoint(100, 0),
		NewPoint(100, 100),
		NewPoint(0, 100),
	}, PTSolid)
	require.NoError(t, err, "error defining root poly")

	hole, err := NewPolyTree([]Point[int]{
		NewPoint(20, 20),
		NewPoint(80, 20),
		NewPoint(80, 80),
		NewPoint(20, 80),
	}, PTHole)
	require.NoError(t, err, "error defining hole poly")

	island, err := NewPolyTree([]Point[int]{
		NewPoint(40, 40),
		NewPoint(60, 40),
		NewPoint(60, 60),
		NewPoint(40, 60),
	}, PTSolid)
	require.NoError(t, err, "error defining island poly")

	err = root.AddChild(hole)
	require.NoError(t, err, "error adding hole as child of root")

	err = hole.AddChild(island)
	require.NoError(t, err, "error adding island as child of hole")

	expectedRoot, err := NewPolyTree([]Point[float64]{
		NewPoint[float64](0, 0),
		NewPoint[float64](100, 0),
		NewPoint[float64](100, 100),
		NewPoint[float64](0, 100),
	}, PTSolid)
	require.NoError(t, err, "error defining expectedRoot poly")

	expectedHole, err := NewPolyTree([]Point[float64]{
		NewPoint[float64](20, 20),
		NewPoint[float64](80, 20),
		NewPoint[float64](80, 80),
		NewPoint[float64](20, 80),
	}, PTHole)
	require.NoError(t, err, "error defining expectedHole poly")

	expectedIsland, err := NewPolyTree([]Point[float64]{
		NewPoint[float64](40, 40),
		NewPoint[float64](60, 40),
		NewPoint[float64](60, 60),
		NewPoint[float64](40, 60),
	}, PTSolid)
	require.NoError(t, err, "error defining expectedIsland poly")

	err = expectedRoot.AddChild(expectedHole)
	require.NoError(t, err, "error adding expectedHole as child of expectedRoot")

	err = expectedHole.AddChild(expectedIsland)
	require.NoError(t, err, "error adding expectedIsland as child of expectedHole")

	// Convert to int (truncated)
	rootIntRounded := root.AsFloat64()

	assert.True(t, rootIntRounded.contour.eq(expectedRoot.contour), "root contour does not match")
	assert.True(t, rootIntRounded.children[0].contour.eq(expectedRoot.children[0].contour), "hole contour does not match")
	assert.True(t, rootIntRounded.children[0].children[0].contour.eq(expectedRoot.children[0].children[0].contour), "hole contour does not match")

}

func TestPolyTree_BooleanOperation(t *testing.T) {
	tests := map[string]struct {
		poly1     [][]Point[int]
		poly2     [][]Point[int]
		operation BooleanOperation
		expected  func() (*PolyTree[int], error)
		wantErr   bool
	}{
		"Union of non-intersecting polygons": {
			poly1:     [][]Point[int]{{{0, 0}, {10, 0}, {10, 10}, {0, 10}}},
			poly2:     [][]Point[int]{{{20, 20}, {30, 20}, {30, 30}, {20, 30}}},
			operation: BooleanUnion,
			expected: func() (*PolyTree[int], error) {
				root, err := NewPolyTree([]Point[int]{{0, 0}, {10, 0}, {10, 10}, {0, 10}}, PTSolid)
				if err != nil {
					return nil, fmt.Errorf("error creating root: %w", err)
				}
				sibling, err := NewPolyTree([]Point[int]{{20, 20}, {30, 20}, {30, 30}, {20, 30}}, PTSolid)
				if err != nil {
					return nil, fmt.Errorf("error creating sibling: %w", err)
				}
				if err := root.AddSibling(sibling); err != nil {
					return nil, fmt.Errorf("error adding sibling: %w", err)
				}
				return root, nil
			},
			wantErr: false,
		},
		"Intersection of non-intersecting polygons": {
			poly1:     [][]Point[int]{{{0, 0}, {10, 0}, {10, 10}, {0, 10}}},
			poly2:     [][]Point[int]{{{20, 20}, {30, 20}, {30, 30}, {20, 30}}},
			operation: BooleanIntersection,
			expected: func() (*PolyTree[int], error) {
				return nil, nil // No intersection
			},
			wantErr: false,
		},
		"Subtraction with non-intersecting polygons": {
			poly1:     [][]Point[int]{{{0, 0}, {10, 0}, {10, 10}, {0, 10}}},
			poly2:     [][]Point[int]{{{20, 20}, {30, 20}, {30, 30}, {20, 30}}},
			operation: BooleanSubtraction,
			expected: func() (*PolyTree[int], error) {
				root, err := NewPolyTree([]Point[int]{{0, 0}, {10, 0}, {10, 10}, {0, 10}}, PTSolid)
				if err != nil {
					return nil, fmt.Errorf("error creating root: %w", err)
				}
				return root, nil
			},
			wantErr: false,
		},
		"Union with one polygon inside another": {
			poly1:     [][]Point[int]{{{0, 0}, {20, 0}, {20, 20}, {0, 20}}},
			poly2:     [][]Point[int]{{{5, 5}, {15, 5}, {15, 15}, {5, 15}}},
			operation: BooleanUnion,
			expected: func() (*PolyTree[int], error) {
				root, err := NewPolyTree([]Point[int]{{0, 0}, {20, 0}, {20, 20}, {0, 20}}, PTSolid)
				if err != nil {
					return nil, fmt.Errorf("error creating root: %w", err)
				}
				hole, err := NewPolyTree([]Point[int]{{5, 5}, {15, 5}, {15, 15}, {5, 15}}, PTHole)
				if err != nil {
					return nil, fmt.Errorf("error creating hole: %w", err)
				}
				if err := root.AddChild(hole); err != nil {
					return nil, fmt.Errorf("error adding hole: %w", err)
				}
				return root, nil
			},
			wantErr: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			poly1, err := nestPointsToPolyTrees(tc.poly1)
			require.NoError(t, err, "error returned from nestPointsToPolyTrees(tc.poly1) when none was expected")
			poly2, err := nestPointsToPolyTrees(tc.poly2)
			require.NoError(t, err, "error returned from nestPointsToPolyTrees(tc.poly2) when none was expected")
			result, err := poly1.BooleanOperation(poly2, tc.operation)

			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			expected, expErr := tc.expected()
			require.NoError(t, expErr, "expected function returned an error")
			match, _ := expected.Eq(result)
			assert.True(t, match, "result did not match expected")
		})
	}
}

func TestPolyTree_booleanOperationTraversal_Intersection(t *testing.T) {
	poly1HolePoints := []Point[int]{
		{5, 5},
		{15, 5},
		{15, 15},
		{5, 15},
	}
	poly1Hole, err := NewPolyTree(poly1HolePoints, PTHole)
	require.NoError(t, err, "expected no error when creating poly1Hole")

	poly1Points := []Point[int]{
		{0, 0},
		{20, 0},
		{20, 20},
		{0, 20},
	}
	polyTree1, err := NewPolyTree(poly1Points, PTSolid, WithChildren(poly1Hole))
	require.NoError(t, err, "expected no error when creating polyTree1")

	poly2HolePoints := []Point[int]{
		{12, 12},
		{22, 12},
		{22, 22},
		{12, 22},
	}
	poly2Hole, err := NewPolyTree(poly2HolePoints, PTHole)
	require.NoError(t, err, "expected no error when creating poly2Hole")

	poly2Points := []Point[int]{
		{7, 7},
		{27, 7},
		{27, 27},
		{7, 27},
	}
	polyTree2, err := NewPolyTree(poly2Points, PTSolid, WithChildren(poly2Hole))
	require.NoError(t, err, "expected no error when creating polyTree2")

	// find intersection points between all polys
	polyTree1.findIntersections(polyTree2)

	// mark points for Intersection
	polyTree1.markEntryExitPoints(polyTree2, BooleanIntersection)

	// traverse for union
	expectedPointsIntersection := [][]Point[int]{
		{
			{40 / 2, 14 / 2},
			{40 / 2, 24 / 2},
			{30 / 2, 24 / 2},
			{30 / 2, 14 / 2},
		},
		{
			{24 / 2, 40 / 2},
			{14 / 2, 40 / 2},
			{14 / 2, 30 / 2},
			{24 / 2, 30 / 2},
		},
	}
	resultingPointsIntersection := polyTree1.booleanOperationTraversal(polyTree2, BooleanIntersection)
	assert.Equal(t, expectedPointsIntersection, resultingPointsIntersection)
}

func TestPolyTree_booleanOperationTraversal_Subtraction(t *testing.T) {
	poly1HolePoints := []Point[int]{
		{5, 5},
		{15, 5},
		{15, 15},
		{5, 15},
	}
	poly1Hole, err := NewPolyTree(poly1HolePoints, PTHole)
	require.NoError(t, err, "expected no error when creating poly1Hole")

	poly1Points := []Point[int]{
		{0, 0},
		{20, 0},
		{20, 20},
		{0, 20},
	}
	polyTree1, err := NewPolyTree(poly1Points, PTSolid, WithChildren(poly1Hole))
	require.NoError(t, err, "expected no error when creating polyTree1")

	poly2HolePoints := []Point[int]{
		{12, 12},
		{22, 12},
		{22, 22},
		{12, 22},
	}
	poly2Hole, err := NewPolyTree(poly2HolePoints, PTHole)
	require.NoError(t, err, "expected no error when creating poly2Hole")

	poly2Points := []Point[int]{
		{7, 7},
		{27, 7},
		{27, 27},
		{7, 27},
	}
	polyTree2, err := NewPolyTree(poly2Points, PTSolid, WithChildren(poly2Hole))
	require.NoError(t, err, "expected no error when creating polyTree2")

	// find intersection points between all polys
	polyTree1.findIntersections(polyTree2)

	// mark points for Intersection
	polyTree1.markEntryExitPoints(polyTree2, BooleanSubtraction)

	expectedPointsSubtraction := [][]Point[int]{
		{
			{40 / 2, 24 / 2},
			{40 / 2, 40 / 2},
			{24 / 2, 40 / 2},
			{24 / 2, 30 / 2},
			{30 / 2, 30 / 2},
			{30 / 2, 24 / 2},
		},
		{
			{14 / 2, 40 / 2},
			{0 / 2, 40 / 2},
			{0 / 2, 0 / 2},
			{40 / 2, 0 / 2},
			{40 / 2, 14 / 2},
			{30 / 2, 14 / 2},
			{30 / 2, 10 / 2},
			{10 / 2, 10 / 2},
			{10 / 2, 30 / 2},
			{14 / 2, 30 / 2},
		},
	}
	resultingPointsSubtraction := polyTree1.booleanOperationTraversal(polyTree2, BooleanSubtraction)
	require.Equal(t, expectedPointsSubtraction, resultingPointsSubtraction)

	// find intersection points between all polys
	polyTree2.findIntersections(polyTree1)

	// mark points for Intersection
	polyTree2.markEntryExitPoints(polyTree1, BooleanSubtraction)

	expectedPointsSubtraction = [][]Point[int]{
		{
			{14 / 2, 30 / 2},
			{14 / 2, 14 / 2},
			{30 / 2, 14 / 2},
			{30 / 2, 24 / 2},
			{24 / 2, 24 / 2},
			{24 / 2, 30 / 2},
		},
		{
			{40 / 2, 14 / 2},
			{54 / 2, 14 / 2},
			{54 / 2, 54 / 2},
			{14 / 2, 54 / 2},
			{14 / 2, 40 / 2},
			{24 / 2, 40 / 2},
			{24 / 2, 44 / 2},
			{44 / 2, 44 / 2},
			{44 / 2, 24 / 2},
			{40 / 2, 24 / 2},
		},
	}
	resultingPointsSubtraction = polyTree2.booleanOperationTraversal(polyTree1, BooleanSubtraction)
	assert.Equal(t, expectedPointsSubtraction, resultingPointsSubtraction)
}

func TestPolyTree_booleanOperationTraversal_Union(t *testing.T) {
	// These polygons were chosen to test union with overlapping regions, holes, and different orientations.

	// Step 1: Create the first polygon tree (polyTree1) with a hole
	polyTree1HolePoints := []Point[int]{{5, 5}, {15, 5}, {15, 15}, {5, 15}}
	polyTree1Hole, err := NewPolyTree(polyTree1HolePoints, PTHole)
	require.NoError(t, err, "unexpected error when creating polyTree1Hole")

	polyTree1Points := []Point[int]{{0, 0}, {20, 0}, {20, 20}, {0, 20}}
	polyTree1, err := NewPolyTree(polyTree1Points, PTSolid, WithChildren(polyTree1Hole))
	require.NoError(t, err, "unexpected error when creating polyTree1")

	// Step 2: Create the second polygon tree (polyTree2) with a hole
	polyTree2HolePoints := []Point[int]{{12, 12}, {22, 12}, {22, 22}, {12, 22}}
	polyTree2Hole, err := NewPolyTree(polyTree2HolePoints, PTHole)
	require.NoError(t, err, "unexpected error when creating polyTree2Hole")

	polyTree2Points := []Point[int]{{7, 7}, {27, 7}, {27, 27}, {7, 27}}
	polyTree2, err := NewPolyTree(polyTree2Points, PTSolid, WithChildren(polyTree2Hole))
	require.NoError(t, err, "unexpected error when creating polyTree2")

	// Step 3: Find intersection points between polyTree1 and polyTree2
	polyTree1.findIntersections(polyTree2)

	// Step 4: Mark entry and exit points for the union operation
	polyTree1.markEntryExitPoints(polyTree2, BooleanUnion)

	// Step 5: Define the expected traversal output for the union operation
	expectedPointsUnion := [][]Point[int]{
		{ // Hole: chevron pointing up-right. Points are ordered based on traversal direction, starting with region entry point.
			{20, 12}, {20, 20}, {12, 20}, {12, 22}, {22, 22}, {22, 12},
		},
		{ // Outer contour. Points are ordered based on traversal direction, starting with region entry point.
			{7, 20}, {0, 20}, {0, 0}, {20, 0}, {20, 7}, {27, 7}, {27, 27}, {7, 27},
		},
		{ // Hole: square in center. Points are ordered based on traversal direction, starting with region entry point.
			{12, 15}, {15, 15}, {15, 12}, {12, 12},
		},
		{ // Hole: chevron pointing down-left. Points are ordered based on traversal direction, starting with region entry point.
			{15, 7}, {15, 5}, {5, 5}, {5, 15}, {7, 15}, {7, 7},
		},
	}

	// Step 6: Perform the traversal for the union operation
	resultingPointsUnion := polyTree1.booleanOperationTraversal(polyTree2, BooleanUnion)

	// Step 7: Assert the resulting points match the expected output
	assert.Equal(t, expectedPointsUnion, resultingPointsUnion, "unexpected output of booleanOperationTraversal for union")
}

func TestPolyTree_BoundingBox(t *testing.T) {
	tests := map[string]struct {
		polyTree *PolyTree[int]
		expected Rectangle[int]
	}{
		"polygon with hole": {
			polyTree: func() *PolyTree[int] {
				root, err := NewPolyTree([]Point[int]{
					NewPoint(-20, -20),
					NewPoint(-20, 40),
					NewPoint(40, 40),
					NewPoint(40, -20),
				}, PTSolid)
				require.NoError(t, err, "error creating root poly")
				hole, err := NewPolyTree([]Point[int]{
					NewPoint(10, 10),
					NewPoint(10, 30),
					NewPoint(30, 30),
					NewPoint(30, 10),
				}, PTHole)
				require.NoError(t, err, "error creating hole poly")
				err = root.AddChild(hole)
				require.NoError(t, err, "error adding hole as child of root")
				return root
			}(),
			expected: NewRectangle([]Point[int]{
				NewPoint(-20, -20),
				NewPoint(-20, 40),
				NewPoint(40, 40),
				NewPoint(40, -20),
			}),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.polyTree.BoundingBox()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestPolyTree_Children(t *testing.T) {
	// Create a root PolyTree
	rootContour := []Point[int]{
		NewPoint(0, 0),
		NewPoint(0, 100),
		NewPoint(100, 100),
		NewPoint(100, 0),
	}
	root, err := NewPolyTree(rootContour, PTSolid)
	require.NoError(t, err, "Unexpected error creating root PolyTree")

	// Create a child PolyTree
	childContour := []Point[int]{
		NewPoint(25, 25),
		NewPoint(25, 75),
		NewPoint(75, 75),
		NewPoint(75, 25),
	}
	child, err := NewPolyTree(childContour, PTHole)
	require.NoError(t, err, "Unexpected error creating child PolyTree")

	// Add the child to the root
	err = root.AddChild(child)
	require.NoError(t, err, "Unexpected error adding child to root")

	// Test the Children method
	children := root.Children()
	assert.Len(t, children, 1, "Expected 1 child")
	assert.Equal(t, child, children[0], "Unexpected child returned")
	assert.Empty(t, child.Children(), "Expected no children for child PolyTree")
}

func TestPolyTree_Edges(t *testing.T) {
	// Create a PolyTree
	contour := []Point[int]{
		NewPoint(0, 0),
		NewPoint(0, 100),
		NewPoint(100, 100),
		NewPoint(100, 0),
	}
	poly, err := NewPolyTree(contour, PTSolid)
	require.NoError(t, err, "Unexpected error creating PolyTree")

	// Retrieve edges
	edges := poly.Edges()
	expectedEdges := []LineSegment[int]{
		NewLineSegment(NewPoint(0, 0), NewPoint(100, 0)),
		NewLineSegment(NewPoint(100, 0), NewPoint(100, 100)),
		NewLineSegment(NewPoint(100, 100), NewPoint(0, 100)),
		NewLineSegment(NewPoint(0, 100), NewPoint(0, 0)),
	}

	// Validate edges
	assert.Len(t, edges, len(expectedEdges), "Expected %d edges", len(expectedEdges))
	for i, edge := range edges {
		assert.Equal(t, expectedEdges[i], edge, "Edge %d mismatch", i)
	}
}

func TestPolyTree_Eq_NilHandling(t *testing.T) {
	var poly1, poly2 *PolyTree[int]

	match, mismatches := poly1.Eq(poly2)
	assert.True(t, match, "Two nil PolyTrees should be considered equal")
	assert.Equal(t, PTMNoMismatch, mismatches, "Expected no mismatches for two nil PolyTrees")

	poly1, _ = NewPolyTree([]Point[int]{{0, 0}, {10, 0}, {10, 10}, {0, 10}}, PTSolid)
	match, mismatches = poly1.Eq(nil)
	assert.False(t, match, "A non-nil PolyTree should not equal a nil PolyTree")
	assert.Equal(t, PTMNilPolygonMismatch, mismatches, "Expected PTMNilPolygonMismatch for a nil comparison")
}

func TestPolyTree_Hull(t *testing.T) {
	// Create a PolyTree with a square contour
	contour := []Point[int]{
		NewPoint(0, 0),
		NewPoint(0, 100),
		NewPoint(100, 100),
		NewPoint(100, 0),
	}
	poly, err := NewPolyTree(contour, PTSolid)
	require.NoError(t, err, "Unexpected error creating PolyTree")

	// Retrieve the convex hull
	hull := poly.Hull()
	expectedHull := []Point[int]{
		NewPoint(0, 0),
		NewPoint(100, 0),
		NewPoint(100, 100),
		NewPoint(0, 100),
	}

	// Validate the hull
	assert.ElementsMatch(t, expectedHull, hull, "Hull points do not match expected")
}

func TestPolyTree_Intersects_EdgeOverlap(t *testing.T) {
	p1 := &PolyTree[int]{contour: contour[int]{
		{point: Point[int]{0, 0}},
		{point: Point[int]{10, 0}},
		{point: Point[int]{10, 10}},
		{point: Point[int]{0, 10}},
	}}
	p2 := &PolyTree[int]{contour: contour[int]{
		{point: Point[int]{5, 0}},
		{point: Point[int]{15, 0}},
		{point: Point[int]{15, 5}},
		{point: Point[int]{5, 5}},
	}}
	require.True(t, p1.Intersects(p2))
}

func TestPolyTree_Intersects_NoIntersection(t *testing.T) {
	p1 := &PolyTree[int]{contour: contour[int]{
		{point: Point[int]{0, 0}},
		{point: Point[int]{10, 0}},
		{point: Point[int]{10, 10}},
		{point: Point[int]{0, 10}},
	}}
	p2 := &PolyTree[int]{contour: contour[int]{
		{point: Point[int]{20, 20}},
		{point: Point[int]{30, 20}},
		{point: Point[int]{30, 30}},
		{point: Point[int]{20, 30}},
	}}
	require.False(t, p1.Intersects(p2))
}

func TestPolyTree_Intersects_OverlappingPolygons(t *testing.T) {
	p1 := &PolyTree[int]{contour: contour[int]{
		{point: Point[int]{0, 0}},
		{point: Point[int]{10, 0}},
		{point: Point[int]{10, 10}},
		{point: Point[int]{0, 10}},
	}}
	p2 := &PolyTree[int]{contour: contour[int]{
		{point: Point[int]{5, 5}},
		{point: Point[int]{15, 5}},
		{point: Point[int]{15, 15}},
		{point: Point[int]{5, 15}},
	}}
	require.True(t, p1.Intersects(p2))
}

func TestPolyTree_Intersects_PointInside(t *testing.T) {
	p1 := &PolyTree[int]{contour: contour[int]{
		{point: Point[int]{0, 0}},
		{point: Point[int]{10, 0}},
		{point: Point[int]{10, 10}},
		{point: Point[int]{0, 10}},
	}}
	p2 := &PolyTree[int]{contour: contour[int]{
		{point: Point[int]{5, 5}},
		{point: Point[int]{6, 5}},
		{point: Point[int]{6, 6}},
		{point: Point[int]{5, 6}},
	}}
	require.True(t, p1.Intersects(p2))
}

func TestPolyTree_IsRoot(t *testing.T) {
	// Create a root PolyTree
	rootContour := []Point[int]{
		NewPoint(0, 0),
		NewPoint(0, 100),
		NewPoint(100, 100),
		NewPoint(100, 0),
	}
	root, err := NewPolyTree(rootContour, PTSolid)
	require.NoError(t, err, "Unexpected error creating root PolyTree")

	// Create a child PolyTree
	childContour := []Point[int]{
		NewPoint(20, 20),
		NewPoint(20, 80),
		NewPoint(80, 80),
		NewPoint(80, 20),
	}
	child, err := NewPolyTree(childContour, PTHole)
	require.NoError(t, err, "Unexpected error creating child PolyTree")
	err = root.AddChild(child)
	require.NoError(t, err, "Unexpected error adding child to root")

	// Validate root
	assert.True(t, root.IsRoot(), "Expected root to be identified as root")

	// Validate child
	assert.False(t, child.IsRoot(), "Expected child not to be identified as root")
}

func TestPolyTree_OrderConsistency(t *testing.T) {
	// Create root and children
	root, err := NewPolyTree([]Point[int]{{10, 10}, {20, 10}, {20, 20}, {10, 20}}, PTSolid)
	require.NoError(t, err, "unexpected error returned when creating root")

	// Valid child polygons
	child1, err := NewPolyTree([]Point[int]{{15, 15}, {18, 15}, {18, 18}, {15, 18}}, PTHole)
	require.NoError(t, err, "unexpected error returned when creating child1")
	child2, err := NewPolyTree([]Point[int]{{11, 11}, {14, 11}, {14, 14}, {11, 14}}, PTHole)
	require.NoError(t, err, "unexpected error returned when creating child2")

	// Add children to root
	err = root.AddChild(child1)
	require.NoError(t, err, "unexpected error returned when adding child1 as a child of root")
	err = root.AddChild(child2)
	require.NoError(t, err, "unexpected error returned when adding child2 as a child of root")

	// Verify children order
	expectedChildOrder := []*PolyTree[int]{child2, child1} // Ordered by lowest, leftmost point
	assert.Equal(t, expectedChildOrder, root.children, "Children should be ordered by lowest, leftmost point")

	// Valid sibling polygons
	sibling1, err := NewPolyTree([]Point[int]{{30, 30}, {40, 30}, {40, 40}, {30, 40}}, PTSolid)
	require.NoError(t, err, "unexpected error returned when creating sibling1")
	sibling2, err := NewPolyTree([]Point[int]{{5, 5}, {9, 5}, {9, 9}, {5, 9}}, PTSolid)
	require.NoError(t, err, "unexpected error returned when creating sibling2")

	// Add siblings to root
	err = root.AddSibling(sibling1)
	require.NoError(t, err, "unexpected error returned when adding sibling1 as a sibling of root")
	err = root.AddSibling(sibling2)
	require.NoError(t, err, "unexpected error returned when adding sibling2 as a sibling of root")

	// Verify sibling order
	expectedSiblingOrder := []*PolyTree[int]{sibling2, sibling1} // Ordered by lowest, leftmost point
	assert.Equal(t, expectedSiblingOrder, root.siblings, "Siblings should be ordered by lowest, leftmost point")
}

func TestPolyTree_Parent(t *testing.T) {
	// Create a root PolyTree
	rootContour := []Point[int]{
		NewPoint(0, 0),
		NewPoint(0, 100),
		NewPoint(100, 100),
		NewPoint(100, 0),
	}
	root, err := NewPolyTree(rootContour, PTSolid)
	require.NoError(t, err, "Unexpected error creating root PolyTree")

	// Create a child PolyTree
	childContour := []Point[int]{
		NewPoint(20, 20),
		NewPoint(20, 80),
		NewPoint(80, 80),
		NewPoint(80, 20),
	}
	child, err := NewPolyTree(childContour, PTHole)
	require.NoError(t, err, "Unexpected error creating child PolyTree")
	err = root.AddChild(child)
	require.NoError(t, err, "Unexpected error adding child to root")

	// Validate root parent
	assert.Nil(t, root.Parent(), "Expected root to have no parent (nil)")

	// Validate child parent
	assert.Equal(t, root, child.Parent(), "Expected child to have root as its parent")
}

func TestPolyTree_Perimeter(t *testing.T) {
	// Create a PolyTree representing a square
	squareContour := []Point[int]{
		NewPoint(0, 0),
		NewPoint(0, 10),
		NewPoint(10, 10),
		NewPoint(10, 0),
	}
	square, err := NewPolyTree(squareContour, PTSolid)
	require.NoError(t, err, "Unexpected error creating PolyTree")

	// Calculate perimeter
	expectedPerimeter := 40.0 // 10 + 10 + 10 + 10
	actualPerimeter := square.Perimeter()

	assert.InDelta(t, expectedPerimeter, actualPerimeter, 1e-10, "Unexpected perimeter value")
}

func TestPolyTree_Points(t *testing.T) {
	// Define a square contour
	cont := []Point[int]{
		NewPoint(0, 0),
		NewPoint(10, 0),
		NewPoint(10, 10),
		NewPoint(0, 10),
	}
	polyTree, err := NewPolyTree(cont, PTSolid)
	require.NoError(t, err, "Unexpected error creating PolyTree")

	// Retrieve the points
	points := polyTree.Points()

	// Verify the points match the contour
	assert.Equal(t, cont, points, "Unexpected points returned from PolyTree")
}

func TestPolyTree_PolygonType(t *testing.T) {
	// Define a solid polygon
	solidContour := []Point[int]{
		NewPoint(0, 0),
		NewPoint(0, 10),
		NewPoint(10, 10),
		NewPoint(10, 0),
	}
	solidPolyTree, err := NewPolyTree(solidContour, PTSolid)
	require.NoError(t, err, "Unexpected error creating solid PolyTree")

	// Check the PolygonType of the solid PolyTree
	assert.Equal(t, PTSolid, solidPolyTree.PolygonType(), "Expected PolygonType to be PTSolid")

	// Define a hole polygon
	holeContour := []Point[int]{
		NewPoint(1, 1),
		NewPoint(1, 9),
		NewPoint(9, 9),
		NewPoint(9, 1),
	}
	holePolyTree, err := NewPolyTree(holeContour, PTHole)
	require.NoError(t, err, "Unexpected error creating hole PolyTree")

	// Check the PolygonType of the hole PolyTree
	assert.Equal(t, PTHole, holePolyTree.PolygonType(), "Expected PolygonType to be PTHole")
}

func TestPolyTree_RelationshipToCircle(t *testing.T) {
	// Create a PolyTree
	root, err := NewPolyTree([]Point[int]{
		NewPoint(0, 0),
		NewPoint(0, 100),
		NewPoint(100, 100),
		NewPoint(100, 0),
	}, PTSolid)
	require.NoError(t, err, "error creating root polygon")

	hole, err := NewPolyTree([]Point[int]{
		NewPoint(20, 20),
		NewPoint(20, 80),
		NewPoint(80, 80),
		NewPoint(80, 20),
	}, PTHole)
	require.NoError(t, err, "error creating hole polygon")
	require.NoError(t, root.AddChild(hole), "error adding hole to root polygon")

	tests := []struct {
		name                     string
		circle                   Circle[int]
		pt                       *PolyTree[int]
		expectedRootRelationship Relationship
		expectedHoleRelationship Relationship
	}{
		{
			name:                     "Circle disjoint from PolyTree",
			circle:                   NewCircle(NewPoint(150, 150), 10),
			pt:                       root,
			expectedRootRelationship: RelationshipDisjoint,
			expectedHoleRelationship: RelationshipDisjoint,
		},
		{
			name:                     "Circle intersects root polygon",
			circle:                   NewCircle(NewPoint(50, 0), 5),
			pt:                       root,
			expectedRootRelationship: RelationshipIntersection,
			expectedHoleRelationship: RelationshipDisjoint,
		},
		{
			name:                     "Circle contained within root polygon",
			circle:                   NewCircle(NewPoint(10, 10), 5),
			pt:                       root,
			expectedRootRelationship: RelationshipContains,
			expectedHoleRelationship: RelationshipDisjoint,
		},
		{
			name:                     "Circle contains root polygon",
			circle:                   NewCircle(NewPoint(50, 50), 100),
			pt:                       root,
			expectedRootRelationship: RelationshipContainedBy,
			expectedHoleRelationship: RelationshipContainedBy,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rels := tt.pt.RelationshipToCircle(tt.circle, WithEpsilon(1e-10))
			assert.Equal(t, tt.expectedRootRelationship, rels[root], "unexpected root relationship")
			assert.Equal(t, tt.expectedHoleRelationship, rels[hole], "unexpected hole relationship")

		})
	}
}

func TestPolyTree_RelationshipToLineSegment(t *testing.T) {
	// Create a PolyTree
	root, err := NewPolyTree([]Point[int]{
		NewPoint(0, 0),
		NewPoint(0, 100),
		NewPoint(100, 100),
		NewPoint(100, 0),
	}, PTSolid)
	require.NoError(t, err, "error creating root polygon")

	hole, err := NewPolyTree([]Point[int]{
		NewPoint(20, 20),
		NewPoint(20, 80),
		NewPoint(80, 80),
		NewPoint(80, 20),
	}, PTHole)
	require.NoError(t, err, "error creating hole polygon")
	require.NoError(t, root.AddChild(hole), "error adding hole to root polygon")

	tests := []struct {
		name                     string
		lineSegment              LineSegment[int]
		pt                       *PolyTree[int]
		expectedRootRelationship Relationship
		expectedHoleRelationship Relationship
	}{
		{
			name:                     "LineSegment disjoint from PolyTree",
			lineSegment:              NewLineSegment(NewPoint(150, 150), NewPoint(200, 200)),
			pt:                       root,
			expectedRootRelationship: RelationshipDisjoint,
			expectedHoleRelationship: RelationshipDisjoint,
		},
		{
			name:                     "LineSegment intersects root polygon",
			lineSegment:              NewLineSegment(NewPoint(-10, 50), NewPoint(10, 50)),
			pt:                       root,
			expectedRootRelationship: RelationshipIntersection,
			expectedHoleRelationship: RelationshipDisjoint,
		},
		{
			name:                     "LineSegment contained within root polygon",
			lineSegment:              NewLineSegment(NewPoint(10, 10), NewPoint(90, 90)),
			pt:                       root,
			expectedRootRelationship: RelationshipContains,
			expectedHoleRelationship: RelationshipIntersection,
		},
		{
			name:                     "LineSegment on edge of root polygon",
			lineSegment:              NewLineSegment(NewPoint(10, 0), NewPoint(90, 0)),
			pt:                       root,
			expectedRootRelationship: RelationshipIntersection,
			expectedHoleRelationship: RelationshipDisjoint,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rels := tt.pt.RelationshipToLineSegment(tt.lineSegment, WithEpsilon(1e-10))
			assert.Equal(t, tt.expectedRootRelationship, rels[root], "unexpected root relationship")
			assert.Equal(t, tt.expectedHoleRelationship, rels[hole], "unexpected hole relationship")
		})
	}
}

func TestPolyTree_RelationshipToPoint(t *testing.T) {
	// Create a PolyTree
	root, err := NewPolyTree([]Point[int]{
		NewPoint(0, 0),
		NewPoint(0, 100),
		NewPoint(100, 100),
		NewPoint(100, 0),
	}, PTSolid)
	require.NoError(t, err, "error creating root polygon")

	hole, err := NewPolyTree([]Point[int]{
		NewPoint(20, 20),
		NewPoint(20, 80),
		NewPoint(80, 80),
		NewPoint(80, 20),
	}, PTHole)
	require.NoError(t, err, "error creating hole polygon")
	require.NoError(t, root.AddChild(hole), "error adding hole to root polygon")

	tests := []struct {
		name                     string
		point                    Point[int]
		pt                       *PolyTree[int]
		expectedRootRelationship Relationship
		expectedHoleRelationship Relationship
	}{
		{
			name:                     "Point outside entire PolyTree",
			point:                    NewPoint(150, 150),
			pt:                       root,
			expectedRootRelationship: RelationshipDisjoint,
			expectedHoleRelationship: RelationshipDisjoint,
		},
		{
			name:                     "Point inside root but outside hole",
			point:                    NewPoint(10, 10),
			pt:                       root,
			expectedRootRelationship: RelationshipContains,
			expectedHoleRelationship: RelationshipDisjoint,
		},
		{
			name:                     "Point inside hole",
			point:                    NewPoint(50, 50),
			pt:                       root,
			expectedRootRelationship: RelationshipContains,
			expectedHoleRelationship: RelationshipContains,
		},
		{
			name:                     "Point on edge of root",
			point:                    NewPoint(0, 50),
			pt:                       root,
			expectedRootRelationship: RelationshipIntersection,
			expectedHoleRelationship: RelationshipDisjoint,
		},
		{
			name:                     "Point on vertex of hole",
			point:                    NewPoint(20, 20),
			pt:                       root,
			expectedRootRelationship: RelationshipContains,
			expectedHoleRelationship: RelationshipIntersection,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rels := tt.pt.RelationshipToPoint(tt.point, WithEpsilon(1e-10))
			assert.Equal(t, tt.expectedRootRelationship, rels[root], "unexpected root relationship")
			assert.Equal(t, tt.expectedHoleRelationship, rels[hole], "unexpected hole relationship")
		})
	}
}

func TestPolyTree_RelationshipToPolyTree(t *testing.T) {
	// Create the first PolyTree
	pt1, err := NewPolyTree([]Point[int]{
		NewPoint(0, 0),
		NewPoint(0, 10),
		NewPoint(10, 10),
		NewPoint(10, 0),
	}, PTSolid)
	require.NoError(t, err, "error creating PolyTree 1")

	// Create the second PolyTree
	pt2, err := NewPolyTree([]Point[int]{
		NewPoint(15, 15),
		NewPoint(15, 25),
		NewPoint(25, 25),
		NewPoint(25, 15),
	}, PTSolid)
	require.NoError(t, err, "error creating PolyTree 2")

	// Create a third PolyTree for testing equality
	pt3, err := NewPolyTree([]Point[int]{
		NewPoint(0, 0),
		NewPoint(0, 10),
		NewPoint(10, 10),
		NewPoint(10, 0),
	}, PTSolid)
	require.NoError(t, err, "error creating PolyTree 3")

	// Perform the relationship checks
	relationships := pt1.RelationshipToPolyTree(pt2)
	require.Len(t, relationships, pt1.Len(), "expected relationships for each polygon in PolyTree 1")

	// Test disjoint relationship
	for _, rel := range relationships[pt1] {
		assert.Equal(t, RelationshipDisjoint, rel, "expected disjoint relationship")
	}

	// Test equality relationship
	relationshipsEqual := pt1.RelationshipToPolyTree(pt3)
	for _, rel := range relationshipsEqual[pt1] {
		assert.Equal(t, RelationshipEqual, rel, "expected equal relationship")
	}

	// Test intersection relationship
	pt4, err := NewPolyTree([]Point[int]{
		NewPoint(5, 0),
		NewPoint(5, 10),
		NewPoint(15, 10),
		NewPoint(15, 0),
	}, PTSolid)
	require.NoError(t, err, "error creating PolyTree 4")

	relationshipsIntersect := pt1.RelationshipToPolyTree(pt4)
	for _, rel := range relationshipsIntersect[pt1] {
		assert.Equal(t, RelationshipIntersection, rel, "expected intersection relationship")
	}

	// Test containment relationships
	pt5, err := NewPolyTree([]Point[int]{
		NewPoint(1, 1),
		NewPoint(1, 9),
		NewPoint(9, 9),
		NewPoint(9, 1),
	}, PTSolid)
	require.NoError(t, err, "error creating PolyTree 5")

	relationshipsContain := pt1.RelationshipToPolyTree(pt5)
	for _, rel := range relationshipsContain[pt1] {
		assert.Equal(t, RelationshipContains, rel, "expected contains relationship")
	}

	relationshipsContainedBy := pt5.RelationshipToPolyTree(pt1)
	for _, rel := range relationshipsContainedBy[pt5] {
		assert.Equal(t, RelationshipContainedBy, rel, "expected contained by relationship")
	}
}

func TestPolyTree_RelationshipToRectangle(t *testing.T) {
	// Create a PolyTree with a root polygon and a hole
	root, err := NewPolyTree([]Point[int]{
		NewPoint(0, 0),
		NewPoint(0, 100),
		NewPoint(100, 100),
		NewPoint(100, 0),
	}, PTSolid)
	require.NoError(t, err, "could not create root polygon")
	hole, err := NewPolyTree([]Point[int]{
		NewPoint(20, 20),
		NewPoint(20, 80),
		NewPoint(80, 80),
		NewPoint(80, 20),
	}, PTHole)
	require.NoError(t, err, "could not create hole polygon")
	err = root.AddChild(hole)
	require.NoError(t, err, "could not add hole as child to root polygon")

	// Test cases
	tests := []struct {
		name     string
		rect     Rectangle[int]
		expected map[*PolyTree[int]]Relationship
	}{
		{
			name: "Rectangle fully contained within root polygon",
			rect: NewRectangle([]Point[int]{
				NewPoint(10, 10),
				NewPoint(90, 10),
				NewPoint(90, 90),
				NewPoint(10, 90),
			}),
			expected: map[*PolyTree[int]]Relationship{
				root: RelationshipContains,
				hole: RelationshipContainedBy,
			},
		},
		{
			name: "Rectangle intersecting the root polygon",
			rect: NewRectangle([]Point[int]{
				NewPoint(-10, 10),
				NewPoint(50, 10),
				NewPoint(50, 50),
				NewPoint(-10, 50),
			}),
			expected: map[*PolyTree[int]]Relationship{
				root: RelationshipIntersection,
				hole: RelationshipIntersection,
			},
		},
		{
			name: "Rectangle fully outside the PolyTree",
			rect: NewRectangle([]Point[int]{
				NewPoint(200, 200),
				NewPoint(300, 200),
				NewPoint(300, 300),
				NewPoint(200, 300),
			}),
			expected: map[*PolyTree[int]]Relationship{
				root: RelationshipDisjoint,
				hole: RelationshipDisjoint,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := root.RelationshipToRectangle(tt.rect)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestPolyTree_Root(t *testing.T) {
	// Create a root polygon
	rootContour := []Point[int]{
		NewPoint(0, 0),
		NewPoint(0, 10),
		NewPoint(10, 10),
		NewPoint(10, 0),
	}
	root, err := NewPolyTree(rootContour, PTSolid)
	require.NoError(t, err)

	// Create a child polygon
	childContour := []Point[int]{
		NewPoint(2, 2),
		NewPoint(2, 8),
		NewPoint(8, 8),
		NewPoint(8, 2),
	}
	child, err := NewPolyTree(childContour, PTHole)
	require.NoError(t, err)
	require.NoError(t, root.AddChild(child))

	// Create a grandchild polygon
	grandchildContour := []Point[int]{
		NewPoint(3, 3),
		NewPoint(3, 7),
		NewPoint(7, 7),
		NewPoint(7, 3),
	}
	grandchild, err := NewPolyTree(grandchildContour, PTSolid)
	require.NoError(t, err)
	require.NoError(t, child.AddChild(grandchild))

	// Test Root method
	assert.Equal(t, root, root.Root(), "Root of the root polygon should be itself")
	assert.Equal(t, root, child.Root(), "Root of the child polygon should be the root polygon")
	assert.Equal(t, root, grandchild.Root(), "Root of the grandchild polygon should be the root polygon")
}

func TestPolyTree_Rotate(t *testing.T) {
	// Create root/parent polygon - large square
	root, err := NewPolyTree([]Point[int]{
		NewPoint(0, 0),
		NewPoint(100, 0),
		NewPoint(100, 100),
		NewPoint(0, 100),
	}, PTSolid)
	require.NoError(t, err)

	// Define pivot point and angle (90° = Pi/2 radians)
	pivot := NewPoint(0, 0)
	angle := math.Pi / 2

	// Perform rotation
	rotated := root.Rotate(pivot, angle)

	// Expected rotated contour (coordinates after 90° rotation around (50, 50))
	expectedContour := []Point[float64]{
		NewPoint[float64](0, 0),
		NewPoint[float64](0, 100),
		NewPoint[float64](-100, 100),
		NewPoint[float64](-100, 0),
	}

	// Verify rotation result
	for i, point := range rotated.Points() {
		assert.InDelta(t, expectedContour[i].X(), point.X(), 1e-10, "Contour should be rotated correctly")
		assert.InDelta(t, expectedContour[i].Y(), point.Y(), 1e-10, "Contour should be rotated correctly")
	}
}

func TestPolyTree_Scale(t *testing.T) {
	// Create root/parent polygon - large square
	root, err := NewPolyTree([]Point[int]{
		NewPoint(0, 0),
		NewPoint(100, 0),
		NewPoint(100, 100),
		NewPoint(0, 100),
	}, PTSolid)
	require.NoError(t, err)

	// Create hole polygon - slightly smaller square
	hole, err := NewPolyTree([]Point[int]{
		NewPoint(20, 20),
		NewPoint(80, 20),
		NewPoint(80, 80),
		NewPoint(20, 80),
	}, PTHole)
	require.NoError(t, err)

	// Create island polygon - even smaller square
	island, err := NewPolyTree([]Point[int]{
		NewPoint(40, 40),
		NewPoint(60, 40),
		NewPoint(60, 60),
		NewPoint(40, 60),
	}, PTSolid)
	require.NoError(t, err)

	// Set up relationships
	require.NoError(t, hole.AddChild(island))
	require.NoError(t, root.AddChild(hole))

	// Scale the PolyTree
	scaled := root.Scale(NewPoint(0, 0), 2)

	// Check root contour
	expectedRoot := []Point[int]{
		NewPoint(0, 0),
		NewPoint(200, 0),
		NewPoint(200, 200),
		NewPoint(0, 200),
	}
	assert.Equal(t, expectedRoot, scaled.Points(), "Root contour should be scaled correctly")

	// Check hole contour
	expectedHole := []Point[int]{
		NewPoint(40, 160),
		NewPoint(160, 160),
		NewPoint(160, 40),
		NewPoint(40, 40),
	}
	assert.Equal(t, expectedHole, scaled.Children()[0].Points(), "Hole contour should be scaled correctly")

	// Check island contour
	expectedIsland := []Point[int]{
		NewPoint(80, 80),
		NewPoint(120, 80),
		NewPoint(120, 120),
		NewPoint(80, 120),
	}
	assert.Equal(t, expectedIsland, scaled.Children()[0].Children()[0].Points(), "Island contour should be scaled correctly")
}

func TestPolyTree_Siblings(t *testing.T) {
	// Create a root polygon
	root, err := NewPolyTree([]Point[int]{
		NewPoint(0, 0),
		NewPoint(0, 100),
		NewPoint(100, 100),
		NewPoint(100, 0),
	}, PTSolid)
	require.NoError(t, err, "failed to create root polygon")

	// Create sibling polygons
	sibling1, err := NewPolyTree([]Point[int]{
		NewPoint(150, 150),
		NewPoint(150, 250),
		NewPoint(250, 250),
		NewPoint(250, 150),
	}, PTSolid)
	require.NoError(t, err, "failed to create sibling polygon 1")

	sibling2, err := NewPolyTree([]Point[int]{
		NewPoint(300, 300),
		NewPoint(300, 400),
		NewPoint(400, 400),
		NewPoint(400, 300),
	}, PTSolid)
	require.NoError(t, err, "failed to create sibling polygon 2")

	// Add siblings
	require.NoError(t, root.AddSibling(sibling1), "failed to add sibling 1")
	require.NoError(t, root.AddSibling(sibling2), "failed to add sibling 2")

	// Test Siblings method
	siblings := root.Siblings()
	assert.ElementsMatch(t, siblings, []*PolyTree[int]{sibling1, sibling2}, "siblings do not match expected values")

	siblingsOfSibling1 := sibling1.Siblings()
	assert.ElementsMatch(t, siblingsOfSibling1, []*PolyTree[int]{root, sibling2}, "sibling1's siblings do not match expected values")

	siblingsOfSibling2 := sibling2.Siblings()
	assert.ElementsMatch(t, siblingsOfSibling2, []*PolyTree[int]{root, sibling1}, "sibling2's siblings do not match expected values")
}

func TestPolyTree_Translate(t *testing.T) {
	// Create root/parent polygon - large square
	root, err := NewPolyTree([]Point[int]{
		NewPoint(0, 0),
		NewPoint(100, 0),
		NewPoint(100, 100),
		NewPoint(0, 100),
	}, PTSolid)
	require.NoError(t, err)

	// Create hole polygon - smaller square
	hole, err := NewPolyTree([]Point[int]{
		NewPoint(20, 20),
		NewPoint(80, 20),
		NewPoint(80, 80),
		NewPoint(20, 80),
	}, PTHole)
	require.NoError(t, err)

	// Add hole to root
	require.NoError(t, root.AddChild(hole))

	// Translate the PolyTree
	delta := NewPoint(10, 10)
	translated := root.Translate(delta)

	// Check root contour
	expectedRoot := []Point[int]{
		NewPoint(10, 10),
		NewPoint(110, 10),
		NewPoint(110, 110),
		NewPoint(10, 110),
	}
	assert.Equal(t, expectedRoot, translated.Points(), "Root contour should be translated correctly")

	// Check hole contour
	expectedHole := []Point[int]{
		NewPoint(30, 90),
		NewPoint(90, 90),
		NewPoint(90, 30),
		NewPoint(30, 30),
	}
	assert.Equal(t, expectedHole, translated.Children()[0].Points(), "Hole contour should be translated correctly")
}

func TestNestPointsToPolyTrees(t *testing.T) {
	tests := map[string]struct {
		contours [][]Point[int]
		expected func() (*PolyTree[int], error)
		wantErr  bool
	}{
		"single polygon": {
			contours: [][]Point[int]{
				{{0, 0}, {10, 0}, {10, 10}, {0, 10}},
			},
			expected: func() (*PolyTree[int], error) {
				return &PolyTree[int]{
					contour: contour[int]{
						polyTreePoint[int]{ // 0
							point:                         Point[int]{0, 0},
							pointType:                     pointTypeOriginal,
							entryExit:                     intersectionTypeNotSet,
							visited:                       false,
							intersectionPartner:           nil,
							intersectionPartnerPointIndex: -1,
						},
						polyTreePoint[int]{ // 1
							point:                         Point[int]{20, 0},
							pointType:                     pointTypeOriginal,
							entryExit:                     intersectionTypeNotSet,
							visited:                       false,
							intersectionPartner:           nil,
							intersectionPartnerPointIndex: -1,
						},
						polyTreePoint[int]{ // 2
							point:                         Point[int]{20, 20},
							pointType:                     pointTypeOriginal,
							entryExit:                     intersectionTypeNotSet,
							visited:                       false,
							intersectionPartner:           nil,
							intersectionPartnerPointIndex: -1,
						},
						polyTreePoint[int]{ // 3
							point:                         Point[int]{0, 20},
							pointType:                     pointTypeOriginal,
							entryExit:                     intersectionTypeNotSet,
							visited:                       false,
							intersectionPartner:           nil,
							intersectionPartnerPointIndex: -1,
						},
					},
					polygonType: PTSolid,
					siblings:    nil,
					children:    nil,
					parent:      nil,
					hull: simpleConvexPolygon[int]{
						Points: []Point[int]{{0, 0}, {10, 0}, {10, 10}, {0, 10}},
					},
					maxX: 21,
				}, nil
			},
			wantErr: false,
		},
		"nested polygons": {
			contours: [][]Point[int]{
				{{0, 0}, {20, 0}, {20, 20}, {0, 20}}, // Outer solid
				{{5, 5}, {15, 5}, {15, 15}, {5, 15}}, // Inner hole
				{{7, 7}, {13, 7}, {13, 13}, {7, 13}}, // Island inside hole
			},
			expected: func() (*PolyTree[int], error) {
				root := &PolyTree[int]{
					contour: contour[int]{
						{Point[int]{0, 0}, pointTypeOriginal, intersectionTypeNotSet, false, nil, -1},
						{Point[int]{40, 0}, pointTypeOriginal, intersectionTypeNotSet, false, nil, -1},
						{Point[int]{40, 40}, pointTypeOriginal, intersectionTypeNotSet, false, nil, -1},
						{Point[int]{0, 40}, pointTypeOriginal, intersectionTypeNotSet, false, nil, -1},
					},
					polygonType: PTSolid,
					hull:        simpleConvexPolygon[int]{Points: []Point[int]{{0, 0}, {20, 0}, {20, 20}, {0, 20}}},
					maxX:        41,
				}
				hole := &PolyTree[int]{
					contour: contour[int]{
						{Point[int]{10, 10}, pointTypeOriginal, intersectionTypeNotSet, false, nil, -1},
						{Point[int]{30, 10}, pointTypeOriginal, intersectionTypeNotSet, false, nil, -1},
						{Point[int]{30, 30}, pointTypeOriginal, intersectionTypeNotSet, false, nil, -1},
						{Point[int]{10, 30}, pointTypeOriginal, intersectionTypeNotSet, false, nil, -1},
					},
					polygonType: PTHole,
					hull:        simpleConvexPolygon[int]{Points: []Point[int]{{5, 5}, {15, 5}, {15, 15}, {5, 15}}},
					maxX:        31,
				}
				island := &PolyTree[int]{
					contour: contour[int]{
						{Point[int]{14, 14}, pointTypeOriginal, intersectionTypeNotSet, false, nil, -1},
						{Point[int]{26, 14}, pointTypeOriginal, intersectionTypeNotSet, false, nil, -1},
						{Point[int]{26, 26}, pointTypeOriginal, intersectionTypeNotSet, false, nil, -1},
						{Point[int]{14, 26}, pointTypeOriginal, intersectionTypeNotSet, false, nil, -1},
					},
					polygonType: PTSolid,
					hull:        simpleConvexPolygon[int]{Points: []Point[int]{{7, 7}, {13, 7}, {13, 13}, {7, 13}}},
					maxX:        27,
				}
				if err := hole.AddChild(island); err != nil {
					return nil, fmt.Errorf("failed to add island: %w", err)
				}
				if err := root.AddChild(hole); err != nil {
					return nil, fmt.Errorf("failed to add hole: %w", err)
				}
				return root, nil
			},
			wantErr: false,
		},
		"no input polygons": {
			contours: [][]Point[int]{},
			expected: func() (*PolyTree[int], error) { return NewPolyTree([]Point[int]{}, PTHole) },
			wantErr:  true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := nestPointsToPolyTrees(tc.contours)
			if tc.wantErr {
				require.Error(t, err, "expected tc.expected() to not raise an error")
				return
			}
			assert.NoError(t, err)

			expected, err := tc.expected()
			require.NoError(t, err, "expected tc.expected() to not raise an error")
			assert.Equal(t, expected, got)
		})
	}
}

func TestNewPolyTree(t *testing.T) {
	tests := map[string]struct {
		points   []Point[int]
		t        PolygonType
		expected func() *PolyTree[int]
	}{
		"solid": {
			points: []Point[int]{
				{x: 0, y: 0},
				{x: 6, y: 0},
				{x: 6, y: 6},
				{x: 0, y: 6},
			},
			t: PTSolid,
			expected: func() *PolyTree[int] {
				return &PolyTree[int]{
					contour: []polyTreePoint[int]{
						{
							point:                         Point[int]{0, 0},
							pointType:                     pointTypeOriginal,
							entryExit:                     intersectionTypeNotSet,
							visited:                       false,
							intersectionPartner:           nil,
							intersectionPartnerPointIndex: -1,
						},
						{
							point:                         Point[int]{12, 0},
							pointType:                     pointTypeOriginal,
							entryExit:                     intersectionTypeNotSet,
							visited:                       false,
							intersectionPartner:           nil,
							intersectionPartnerPointIndex: -1,
						},
						{
							point:                         Point[int]{12, 12},
							pointType:                     pointTypeOriginal,
							entryExit:                     intersectionTypeNotSet,
							visited:                       false,
							intersectionPartner:           nil,
							intersectionPartnerPointIndex: -1,
						},
						{
							point:                         Point[int]{0, 12},
							pointType:                     pointTypeOriginal,
							entryExit:                     intersectionTypeNotSet,
							visited:                       false,
							intersectionPartner:           nil,
							intersectionPartnerPointIndex: -1,
						},
					},
					polygonType: PTSolid,
					children:    nil,
					parent:      nil,
					hull: simpleConvexPolygon[int]{
						Points: []Point[int]{
							{x: 0, y: 0},
							{x: 6, y: 0},
							{x: 6, y: 6},
							{x: 0, y: 6},
						},
					},
					maxX: 13,
				}
			},
		},
		"hole": {
			points: []Point[int]{
				{x: 0, y: 0},
				{x: 6, y: 0},
				{x: 6, y: 6},
				{x: 0, y: 6},
			},
			t: PTHole,
			expected: func() *PolyTree[int] {
				return &PolyTree[int]{
					contour: []polyTreePoint[int]{
						{
							point:                         Point[int]{0, 12},
							pointType:                     pointTypeOriginal,
							entryExit:                     intersectionTypeNotSet,
							visited:                       false,
							intersectionPartner:           nil,
							intersectionPartnerPointIndex: -1,
						},
						{
							point:                         Point[int]{12, 12},
							pointType:                     pointTypeOriginal,
							entryExit:                     intersectionTypeNotSet,
							visited:                       false,
							intersectionPartner:           nil,
							intersectionPartnerPointIndex: -1,
						},
						{
							point:                         Point[int]{12, 0},
							pointType:                     pointTypeOriginal,
							entryExit:                     intersectionTypeNotSet,
							visited:                       false,
							intersectionPartner:           nil,
							intersectionPartnerPointIndex: -1,
						},
						{
							point:                         Point[int]{0, 0},
							pointType:                     pointTypeOriginal,
							entryExit:                     intersectionTypeNotSet,
							visited:                       false,
							intersectionPartner:           nil,
							intersectionPartnerPointIndex: -1,
						},
					},
					polygonType: PTHole,
					children:    nil,
					parent:      nil,
					hull: simpleConvexPolygon[int]{
						Points: []Point[int]{
							{x: 0, y: 0},
							{x: 6, y: 0},
							{x: 6, y: 6},
							{x: 0, y: 6},
						},
					},
					maxX: 13,
				}
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Create the PolyTree using the NewPolyTree function
			result, err := NewPolyTree(tc.points, tc.t)
			require.NoError(t, err, "unexpected error from NewPolyTree")

			// Retrieve the expected result from the test case
			expected := tc.expected()

			// Ensure the points in the contour are ordered correctly for both expected and result
			expected.orderSiblingsAndChildren()
			result.orderSiblingsAndChildren()

			// Use Eq to compare the result and expected trees for a more flexible comparison
			equal, mismatches := result.Eq(expected)
			assert.True(t, equal, "unexpected mismatch: %v", mismatches)
		})
	}
}

func TestNewPolyTree_Errors(t *testing.T) {
	tests := map[string]struct {
		NewPolyFunc    func() (*PolyTree[int], error)
		expectedErrMsg string
	}{
		"Less than three points": {
			NewPolyFunc: func() (*PolyTree[int], error) {
				return NewPolyTree(
					[]Point[int]{NewPoint(0, 0), NewPoint(1, 1)},
					PTSolid,
				)
			},
			expectedErrMsg: "new polytree must have at least 3 points",
		},
		"Zero area polygon": {
			NewPolyFunc: func() (*PolyTree[int], error) {
				return NewPolyTree(
					[]Point[int]{
						NewPoint(0, 0),
						NewPoint(1, 1),
						NewPoint(2, 2),
					},
					PTSolid,
				)
			},
			expectedErrMsg: "new polytree must have non-zero area",
		},
		"Invalid child polygon type for hole": {
			NewPolyFunc: func() (*PolyTree[int], error) {
				hole, err := NewPolyTree(
					[]Point[int]{
						NewPoint(20, 20),
						NewPoint(80, 20),
						NewPoint(80, 80),
						NewPoint(20, 80),
					}, PTSolid)
				if err != nil {
					return nil, err
				}
				return NewPolyTree(
					[]Point[int]{
						NewPoint(0, 0),
						NewPoint(0, 100),
						NewPoint(100, 100),
						NewPoint(100, 0),
					},
					PTSolid,
					WithChildren(hole),
				)
			},
			expectedErrMsg: "cannot add child: mismatched polygon types (parent: PTSolid, child: PTSolid)",
		},
		"Invalid child polygon type for island": {
			NewPolyFunc: func() (*PolyTree[int], error) {
				island, err := NewPolyTree(
					[]Point[int]{
						NewPoint(20, 20),
						NewPoint(80, 20),
						NewPoint(80, 80),
						NewPoint(20, 80),
					}, PTHole)
				if err != nil {
					return nil, err
				}
				return NewPolyTree(
					[]Point[int]{
						NewPoint(0, 0),
						NewPoint(0, 100),
						NewPoint(100, 100),
						NewPoint(100, 0),
					},
					PTHole,
					WithChildren(island),
				)
			},
			expectedErrMsg: "cannot add child: mismatched polygon types (parent: PTHole, child: PTHole)",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := tt.NewPolyFunc()
			require.Error(t, err, "expected error but got nil")
			assert.Contains(t, err.Error(), tt.expectedErrMsg, "unexpected error message")
		})
	}
}
