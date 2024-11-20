package geom2d

import (
	fmt "fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEdgesFromPolyPoints(t *testing.T) {
	tests := map[string]struct {
		points []polyPoint[int]
		edges  []polyEdge[int]
	}{
		"right triangle": {
			points: []polyPoint[int]{
				{point: NewPoint(0, 0)},
				{point: NewPoint(10, 0)},
				{point: NewPoint(10, 10)},
			},
			edges: []polyEdge[int]{
				{lineSegment: NewLineSegment(NewPoint(0, 0), NewPoint(10, 0))},
				{lineSegment: NewLineSegment(NewPoint(10, 0), NewPoint(10, 10))},
				{lineSegment: NewLineSegment(NewPoint(10, 10), NewPoint(0, 0))},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := edgesFromPolyPoints(tc.points)
			assert.Equal(t, tc.edges, actual)
		})
	}
}

func TestIsInsidePolygon(t *testing.T) {
	tests := map[string]struct {
		polygon  []polyPoint[int]
		point    Point[int]
		expected bool
	}{
		"Point inside convex polygon": {
			polygon: []polyPoint[int]{
				{point: NewPoint(0, 0)},
				{point: NewPoint(4, 0)},
				{point: NewPoint(4, 4)},
				{point: NewPoint(0, 4)},
			},
			point:    NewPoint(2, 2),
			expected: true,
		},
		"Point outside convex polygon": {
			polygon: []polyPoint[int]{
				{point: NewPoint(0, 0)},
				{point: NewPoint(4, 0)},
				{point: NewPoint(4, 4)},
				{point: NewPoint(0, 4)},
			},
			point:    NewPoint(5, 5),
			expected: false,
		},
		"Point outside convex polygon #2": {
			polygon: []polyPoint[int]{
				{point: NewPoint(3, 4)},
				{point: NewPoint(2, 2)},
				{point: NewPoint(4, 2)},
				{point: NewPoint(6, 2)},
				{point: NewPoint(6, 6)},
				{point: NewPoint(2, 6)},
			},
			point:    NewPoint(1, 4),
			expected: false,
		},
		"Point on edge of polygon": {
			polygon: []polyPoint[int]{
				{point: NewPoint(0, 0)},
				{point: NewPoint(4, 0)},
				{point: NewPoint(4, 4)},
				{point: NewPoint(0, 4)},
			},
			point:    NewPoint(2, 0),
			expected: true, // Considered inside for edge cases
		},
		"Point on vertex of convex polygon": {
			polygon: []polyPoint[int]{
				{point: NewPoint(0, 0)},
				{point: NewPoint(4, 0)},
				{point: NewPoint(4, 4)},
				{point: NewPoint(0, 4)},
			},
			point:    NewPoint(0, 0),
			expected: true, // Considered inside for vertex cases
		},
		"Point inside concave polygon": {
			polygon: []polyPoint[int]{
				{point: NewPoint(0, 0)},
				{point: NewPoint(6, 0)},
				{point: NewPoint(3, 3)},
				{point: NewPoint(6, 6)},
				{point: NewPoint(0, 6)},
			},
			point:    NewPoint(3, 2),
			expected: true,
		},
		"Point outside concave polygon": {
			polygon: []polyPoint[int]{
				{point: NewPoint(0, 0)},
				{point: NewPoint(6, 0)},
				{point: NewPoint(3, 3)},
				{point: NewPoint(6, 6)},
				{point: NewPoint(0, 6)},
			},
			point:    NewPoint(5, 3),
			expected: false,
		},
		"Point on concave polygon edge": {
			polygon: []polyPoint[int]{
				{point: NewPoint(0, 0)},
				{point: NewPoint(6, 0)},
				{point: NewPoint(3, 3)},
				{point: NewPoint(6, 6)},
				{point: NewPoint(0, 6)},
			},
			point:    NewPoint(4, 2),
			expected: true, // Considered inside for edge cases
		},
		"inside concave polygon in line with vertex 'stalagmite'": {
			polygon: []polyPoint[int]{
				{point: NewPoint(0, 0)},
				{point: NewPoint(3, 0)},
				{point: NewPoint(5, 7)},
				{point: NewPoint(7, 0)},
				{point: NewPoint(14, 0)},
				{point: NewPoint(14, 10)},
				{point: NewPoint(11, 10)},
				{point: NewPoint(9, 3)},
				{point: NewPoint(7, 10)},
				{point: NewPoint(0, 10)},
			},
			point:    NewPoint(2, 7),
			expected: true, // Considered inside for edge cases
		},
		"inside concave polygon in line with vertex 'stalagtite'": {
			polygon: []polyPoint[int]{
				{point: NewPoint(0, 0)},
				{point: NewPoint(3, 0)},
				{point: NewPoint(5, 7)},
				{point: NewPoint(7, 0)},
				{point: NewPoint(14, 0)},
				{point: NewPoint(14, 10)},
				{point: NewPoint(11, 10)},
				{point: NewPoint(9, 3)},
				{point: NewPoint(7, 10)},
				{point: NewPoint(0, 10)},
			},
			point:    NewPoint(2, 3),
			expected: true, // Considered inside for edge cases
		},
		"outside concave polygon in line with vertex 'stalagmite'": {
			polygon: []polyPoint[int]{
				{point: NewPoint(0, 0)},
				{point: NewPoint(3, 0)},
				{point: NewPoint(5, 7)},
				{point: NewPoint(7, 0)},
				{point: NewPoint(14, 0)},
				{point: NewPoint(14, 10)},
				{point: NewPoint(11, 10)},
				{point: NewPoint(9, 3)},
				{point: NewPoint(7, 10)},
				{point: NewPoint(0, 10)},
			},
			point:    NewPoint(-3, 7),
			expected: false,
		},
		"outside concave polygon in line with vertex 'stalagtite'": {
			polygon: []polyPoint[int]{
				{point: NewPoint(0, 0)},
				{point: NewPoint(3, 0)},
				{point: NewPoint(5, 7)},
				{point: NewPoint(7, 0)},
				{point: NewPoint(14, 0)},
				{point: NewPoint(14, 10)},
				{point: NewPoint(11, 10)},
				{point: NewPoint(9, 3)},
				{point: NewPoint(7, 10)},
				{point: NewPoint(0, 10)},
			},
			point:    NewPoint(-3, 3),
			expected: false,
		},
		"outside concave polygon with point inside 'stalagtite'": {
			polygon: []polyPoint[int]{
				{point: NewPoint(0, 0)},
				{point: NewPoint(3, 0)},
				{point: NewPoint(5, 7)},
				{point: NewPoint(7, 0)},
				{point: NewPoint(14, 0)},
				{point: NewPoint(14, 10)},
				{point: NewPoint(11, 10)},
				{point: NewPoint(9, 3)},
				{point: NewPoint(7, 10)},
				{point: NewPoint(0, 10)},
			},
			point:    NewPoint(9, 10),
			expected: false,
		},
		"outside concave polygon with point inside 'stalagmite'": {
			polygon: []polyPoint[int]{
				{point: NewPoint(0, 0)},
				{point: NewPoint(3, 0)},
				{point: NewPoint(5, 7)},
				{point: NewPoint(7, 0)},
				{point: NewPoint(14, 0)},
				{point: NewPoint(14, 10)},
				{point: NewPoint(11, 10)},
				{point: NewPoint(9, 3)},
				{point: NewPoint(7, 10)},
				{point: NewPoint(0, 10)},
			},
			point:    NewPoint(5, 0),
			expected: false,
		},
		"outside concave polygon with point inside 'right chevron'": {
			polygon: []polyPoint[int]{
				{point: NewPoint(0, 0)},
				{point: NewPoint(3, 0)},
				{point: NewPoint(5, 7)},
				{point: NewPoint(7, 0)},
				{point: NewPoint(14, 0)},
				{point: NewPoint(12, 5)},
				{point: NewPoint(14, 10)},
				{point: NewPoint(11, 10)},
				{point: NewPoint(9, 3)},
				{point: NewPoint(7, 10)},
				{point: NewPoint(0, 10)},
				{point: NewPoint(2, 5)},
			},
			point:    NewPoint(1, 5),
			expected: false,
		},
		"outside concave polygon with point inside 'left chevron'": {
			polygon: []polyPoint[int]{
				{point: NewPoint(0, 0)},
				{point: NewPoint(3, 0)},
				{point: NewPoint(5, 7)},
				{point: NewPoint(7, 0)},
				{point: NewPoint(14, 0)},
				{point: NewPoint(12, 5)},
				{point: NewPoint(14, 10)},
				{point: NewPoint(11, 10)},
				{point: NewPoint(9, 3)},
				{point: NewPoint(7, 10)},
				{point: NewPoint(0, 10)},
				{point: NewPoint(2, 5)},
			},
			point:    NewPoint(13, 5),
			expected: false,
		},
		"inside concave polygon with point beside 'right chevron'": {
			polygon: []polyPoint[int]{
				{point: NewPoint(0, 0)},
				{point: NewPoint(3, 0)},
				{point: NewPoint(5, 7)},
				{point: NewPoint(7, 0)},
				{point: NewPoint(14, 0)},
				{point: NewPoint(12, 5)},
				{point: NewPoint(14, 10)},
				{point: NewPoint(11, 10)},
				{point: NewPoint(9, 3)},
				{point: NewPoint(7, 10)},
				{point: NewPoint(0, 10)},
				{point: NewPoint(2, 5)},
			},
			point:    NewPoint(3, 5),
			expected: true,
		},
		"inside concave polygon with point beside 'left chevron'": {
			polygon: []polyPoint[int]{
				{point: NewPoint(0, 0)},
				{point: NewPoint(3, 0)},
				{point: NewPoint(5, 7)},
				{point: NewPoint(7, 0)},
				{point: NewPoint(14, 0)},
				{point: NewPoint(12, 5)},
				{point: NewPoint(14, 10)},
				{point: NewPoint(11, 10)},
				{point: NewPoint(9, 3)},
				{point: NewPoint(7, 10)},
				{point: NewPoint(0, 10)},
				{point: NewPoint(2, 5)},
			},
			point:    NewPoint(11, 5),
			expected: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := isInsidePolygon(tc.polygon, tc.point)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFindIntersectionsBetweenPolys(t *testing.T) {
	tests := map[string]struct {
		poly1            []polyPoint[int]
		poly2            []polyPoint[int]
		expectedPoly1Out []polyPoint[int]
		expectedPoly2Out []polyPoint[int]
	}{
		"No intersections": {
			poly1: []polyPoint[int]{
				{point: NewPoint(0, 0), pointType: PointTypeNormal},
				{point: NewPoint(4, 0), pointType: PointTypeNormal},
				{point: NewPoint(2, 2), pointType: PointTypeNormal},
			},
			poly2: []polyPoint[int]{
				{point: NewPoint(6, 6), pointType: PointTypeNormal},
				{point: NewPoint(8, 6), pointType: PointTypeNormal},
				{point: NewPoint(7, 8), pointType: PointTypeNormal},
			},
			expectedPoly1Out: []polyPoint[int]{
				{point: NewPoint(0, 0), pointType: PointTypeNormal},
				{point: NewPoint(4, 0), pointType: PointTypeNormal},
				{point: NewPoint(2, 2), pointType: PointTypeNormal},
			},
			expectedPoly2Out: []polyPoint[int]{
				{point: NewPoint(6, 6), pointType: PointTypeNormal},
				{point: NewPoint(8, 6), pointType: PointTypeNormal},
				{point: NewPoint(7, 8), pointType: PointTypeNormal},
			},
		},
		"Multiple intersections": {
			poly1: []polyPoint[int]{
				{point: NewPoint(0, 0), pointType: PointTypeNormal},
				{point: NewPoint(4, 0), pointType: PointTypeNormal},
				{point: NewPoint(4, 4), pointType: PointTypeNormal},
				{point: NewPoint(0, 4), pointType: PointTypeNormal},
			},
			poly2: []polyPoint[int]{
				{point: NewPoint(2, 2), pointType: PointTypeNormal},
				{point: NewPoint(6, 2), pointType: PointTypeNormal},
				{point: NewPoint(6, 6), pointType: PointTypeNormal},
				{point: NewPoint(2, 6), pointType: PointTypeNormal},
			},
			expectedPoly1Out: []polyPoint[int]{
				{point: NewPoint(0, 0), pointType: PointTypeNormal},
				{point: NewPoint(4, 0), pointType: PointTypeNormal},
				{point: NewPoint(4, 2), pointType: PointTypeAddedIntersection},
				{point: NewPoint(4, 4), pointType: PointTypeNormal},
				{point: NewPoint(2, 4), pointType: PointTypeAddedIntersection},
				{point: NewPoint(0, 4), pointType: PointTypeNormal},
			},
			expectedPoly2Out: []polyPoint[int]{
				{point: NewPoint(2, 4), pointType: PointTypeAddedIntersection},
				{point: NewPoint(2, 2), pointType: PointTypeNormal},
				{point: NewPoint(4, 2), pointType: PointTypeAddedIntersection},
				{point: NewPoint(6, 2), pointType: PointTypeNormal},
				{point: NewPoint(6, 6), pointType: PointTypeNormal},
				{point: NewPoint(2, 6), pointType: PointTypeNormal},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			poly1Out, poly2Out := findIntersectionsBetweenPolys(tc.poly1, tc.poly2)
			assert.Equal(t, tc.expectedPoly1Out, poly1Out)
			assert.Equal(t, tc.expectedPoly2Out, poly2Out)
		})
	}
}

func TestMarkEntryExitPoints(t *testing.T) {
	tests := map[string]struct {
		poly1         []polyPoint[int]
		poly2         []polyPoint[int]
		operation     BooleanOperation
		expectedPoly1 []polyPoint[int]
		expectedPoly2 []polyPoint[int]
	}{
		"Simple intersection for union": {
			poly1: []polyPoint[int]{
				{point: NewPoint(0, 0), pointType: PointTypeNormal},
				{point: NewPoint(4, 0), pointType: PointTypeNormal},
				{point: NewPoint(4, 4), pointType: PointTypeNormal},
				{point: NewPoint(0, 4), pointType: PointTypeNormal},
			},
			poly2: []polyPoint[int]{
				{point: NewPoint(2, 2), pointType: PointTypeNormal},
				{point: NewPoint(6, 2), pointType: PointTypeNormal},
				{point: NewPoint(6, 6), pointType: PointTypeNormal},
				{point: NewPoint(2, 6), pointType: PointTypeNormal},
			},
			operation: BooleanUnion,
			expectedPoly1: []polyPoint[int]{
				{point: NewPoint(0, 0), pointType: PointTypeNormal},
				{point: NewPoint(4, 0), pointType: PointTypeNormal},
				{point: NewPoint(4, 2), pointType: PointTypeAddedIntersection, entryExit: Exit, otherPolyIndex: 2},
				{point: NewPoint(4, 4), pointType: PointTypeNormal},
				{point: NewPoint(2, 4), pointType: PointTypeAddedIntersection, entryExit: Entry, otherPolyIndex: 0},
				{point: NewPoint(0, 4), pointType: PointTypeNormal},
			},
			expectedPoly2: []polyPoint[int]{
				{point: NewPoint(2, 4), pointType: PointTypeAddedIntersection, entryExit: Exit, otherPolyIndex: 4},
				{point: NewPoint(2, 2), pointType: PointTypeNormal},
				{point: NewPoint(4, 2), pointType: PointTypeAddedIntersection, entryExit: Entry, otherPolyIndex: 2},
				{point: NewPoint(6, 2), pointType: PointTypeNormal},
				{point: NewPoint(6, 6), pointType: PointTypeNormal},
				{point: NewPoint(2, 6), pointType: PointTypeNormal},
			},
		},
		"Simple intersection for intersection": {
			poly1: []polyPoint[int]{
				{point: NewPoint(0, 0), pointType: PointTypeNormal},
				{point: NewPoint(4, 0), pointType: PointTypeNormal},
				{point: NewPoint(4, 4), pointType: PointTypeNormal},
				{point: NewPoint(0, 4), pointType: PointTypeNormal},
			},
			poly2: []polyPoint[int]{
				{point: NewPoint(2, 2), pointType: PointTypeNormal},
				{point: NewPoint(6, 2), pointType: PointTypeNormal},
				{point: NewPoint(6, 6), pointType: PointTypeNormal},
				{point: NewPoint(2, 6), pointType: PointTypeNormal},
			},
			operation: BooleanIntersection,
			expectedPoly1: []polyPoint[int]{
				{point: NewPoint(0, 0), pointType: PointTypeNormal},
				{point: NewPoint(4, 0), pointType: PointTypeNormal},
				{point: NewPoint(4, 2), pointType: PointTypeAddedIntersection, entryExit: Entry, otherPolyIndex: 2},
				{point: NewPoint(4, 4), pointType: PointTypeNormal},
				{point: NewPoint(2, 4), pointType: PointTypeAddedIntersection, entryExit: Exit, otherPolyIndex: 0},
				{point: NewPoint(0, 4), pointType: PointTypeNormal},
			},
			expectedPoly2: []polyPoint[int]{
				{point: NewPoint(2, 4), pointType: PointTypeAddedIntersection, entryExit: Entry, otherPolyIndex: 4},
				{point: NewPoint(2, 2), pointType: PointTypeNormal},
				{point: NewPoint(4, 2), pointType: PointTypeAddedIntersection, entryExit: Exit, otherPolyIndex: 2},
				{point: NewPoint(6, 2), pointType: PointTypeNormal},
				{point: NewPoint(6, 6), pointType: PointTypeNormal},
				{point: NewPoint(2, 6), pointType: PointTypeNormal},
			},
		},
		"Simple intersection for subtraction": {
			poly1: []polyPoint[int]{
				{point: NewPoint(0, 0), pointType: PointTypeNormal},
				{point: NewPoint(4, 0), pointType: PointTypeNormal},
				{point: NewPoint(4, 4), pointType: PointTypeNormal},
				{point: NewPoint(0, 4), pointType: PointTypeNormal},
			},
			poly2: []polyPoint[int]{
				{point: NewPoint(2, 2), pointType: PointTypeNormal},
				{point: NewPoint(6, 2), pointType: PointTypeNormal},
				{point: NewPoint(6, 6), pointType: PointTypeNormal},
				{point: NewPoint(2, 6), pointType: PointTypeNormal},
			},
			operation: BooleanSubtraction,
			expectedPoly1: []polyPoint[int]{
				{point: NewPoint(0, 0), pointType: PointTypeNormal},
				{point: NewPoint(4, 0), pointType: PointTypeNormal},
				{point: NewPoint(4, 2), pointType: PointTypeAddedIntersection, entryExit: Exit, otherPolyIndex: 2},
				{point: NewPoint(4, 4), pointType: PointTypeNormal},
				{point: NewPoint(2, 4), pointType: PointTypeAddedIntersection, entryExit: Entry, otherPolyIndex: 0},
				{point: NewPoint(0, 4), pointType: PointTypeNormal},
			},
			expectedPoly2: []polyPoint[int]{
				{point: NewPoint(2, 4), pointType: PointTypeAddedIntersection, entryExit: Entry, otherPolyIndex: 4},
				{point: NewPoint(2, 2), pointType: PointTypeNormal},
				{point: NewPoint(4, 2), pointType: PointTypeAddedIntersection, entryExit: Exit, otherPolyIndex: 2},
				{point: NewPoint(6, 2), pointType: PointTypeNormal},
				{point: NewPoint(6, 6), pointType: PointTypeNormal},
				{point: NewPoint(2, 6), pointType: PointTypeNormal},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			poly1, poly2 := findIntersectionsBetweenPolys(tc.poly1, tc.poly2)
			markEntryExitPoints(poly1, poly2, tc.operation)
			assert.Equal(t, tc.expectedPoly1, poly1, "Poly1 mismatch")
			assert.Equal(t, tc.expectedPoly2, poly2, "Poly2 mismatch")
		})
	}
}

func TestTraverse_Intersection(t *testing.T) {
	tests := map[string]struct {
		poly1           []polyPoint[int]
		poly2           []polyPoint[int]
		expectedResults [][]polyPoint[int]
	}{
		"Simple intersection": {
			poly1: []polyPoint[int]{
				{point: NewPoint(0, 0)},
				{point: NewPoint(4, 0)},
				{point: NewPoint(4, 4)},
				{point: NewPoint(0, 4)},
			},
			poly2: []polyPoint[int]{
				{point: NewPoint(2, 2)},
				{point: NewPoint(6, 2)},
				{point: NewPoint(6, 6)},
				{point: NewPoint(2, 6)},
			},
			expectedResults: [][]polyPoint[int]{
				{
					{point: NewPoint(4, 2), pointType: PointTypeAddedIntersection, entryExit: Entry, visited: false, otherPolyIndex: 2},
					{point: NewPoint(4, 4), pointType: PointTypeNormal, entryExit: NotSet, visited: false, otherPolyIndex: 0},
					{point: NewPoint(2, 4), pointType: PointTypeAddedIntersection, entryExit: Entry, visited: false, otherPolyIndex: 4},
					{point: NewPoint(2, 2), pointType: PointTypeNormal, entryExit: NotSet, visited: false, otherPolyIndex: 0},
				},
			},
		},
		"No intersection": {
			poly1: []polyPoint[int]{
				{point: NewPoint(0, 0)},
				{point: NewPoint(4, 0)},
				{point: NewPoint(4, 4)},
				{point: NewPoint(0, 4)},
			},
			poly2: []polyPoint[int]{
				{point: NewPoint(6, 6)},
				{point: NewPoint(10, 6)},
				{point: NewPoint(10, 10)},
				{point: NewPoint(6, 10)},
			},
			expectedResults: [][]polyPoint[int]{}, // No intersection
		},
		"Nested polygons (poly2 inside poly1)": {
			poly1: []polyPoint[int]{
				{point: NewPoint(0, 0)},
				{point: NewPoint(6, 0)},
				{point: NewPoint(6, 6)},
				{point: NewPoint(0, 6)},
			},
			poly2: []polyPoint[int]{
				{point: NewPoint(2, 2)},
				{point: NewPoint(4, 2)},
				{point: NewPoint(4, 4)},
				{point: NewPoint(2, 4)},
			},
			expectedResults: [][]polyPoint[int]{
				{
					{point: NewPoint(2, 2), pointType: PointTypeNormal, entryExit: NotSet, visited: false, otherPolyIndex: 0},
					{point: NewPoint(4, 2), pointType: PointTypeNormal, entryExit: NotSet, visited: false, otherPolyIndex: 0},
					{point: NewPoint(4, 4), pointType: PointTypeNormal, entryExit: NotSet, visited: false, otherPolyIndex: 0},
					{point: NewPoint(2, 4), pointType: PointTypeNormal, entryExit: NotSet, visited: false, otherPolyIndex: 0},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Run("Debug Simple intersection", func(t *testing.T) {
				poly1, poly2 := findIntersectionsBetweenPolys(tc.poly1, tc.poly2)
				fmt.Println("Poly1 after intersections:", poly1)
				fmt.Println("Poly2 after intersections:", poly2)

				markEntryExitPoints(poly1, poly2, BooleanIntersection)
				fmt.Println("Poly1 after marking entry/exit:", poly1)
				fmt.Println("Poly2 after marking entry/exit:", poly2)

				results := traverse(poly1, poly2, BooleanIntersection)
				fmt.Println("Traversal results:", results)

				assert.Equal(t, tc.expectedResults, results, "Traversal result mismatch")
			})
		})
	}
}

func TestTraverse_Union(t *testing.T) {
	tests := map[string]struct {
		poly1           []polyPoint[int]
		poly2           []polyPoint[int]
		expectedResults [][]polyPoint[int]
	}{
		"Multiple intersections": {
			poly1: []polyPoint[int]{
				{point: NewPoint(0, 0)},
				{point: NewPoint(4, 0)},
				{point: NewPoint(4, 4)},
				{point: NewPoint(0, 4)},
			},
			poly2: []polyPoint[int]{
				{point: NewPoint(2, 2)},
				{point: NewPoint(6, 2)},
				{point: NewPoint(6, 6)},
				{point: NewPoint(2, 6)},
			},
			expectedResults: [][]polyPoint[int]{
				{
					{point: NewPoint(2, 4), pointType: PointTypeAddedIntersection, entryExit: Entry, visited: false, otherPolyIndex: 0},
					{point: NewPoint(0, 4), pointType: PointTypeNormal, entryExit: NotSet, visited: false, otherPolyIndex: 0},
					{point: NewPoint(0, 0), pointType: PointTypeNormal, entryExit: NotSet, visited: false, otherPolyIndex: 0},
					{point: NewPoint(4, 0), pointType: PointTypeNormal, entryExit: NotSet, visited: false, otherPolyIndex: 0},
					{point: NewPoint(4, 2), pointType: PointTypeAddedIntersection, entryExit: Entry, visited: false, otherPolyIndex: 2},
					{point: NewPoint(6, 2), pointType: PointTypeNormal, entryExit: NotSet, visited: false, otherPolyIndex: 0},
					{point: NewPoint(6, 6), pointType: PointTypeNormal, entryExit: NotSet, visited: false, otherPolyIndex: 0},
					{point: NewPoint(2, 6), pointType: PointTypeNormal, entryExit: NotSet, visited: false, otherPolyIndex: 0},
				},
			},
		},
		"Non-intersecting polygons": {
			poly1: []polyPoint[int]{
				{point: NewPoint(0, 0)},
				{point: NewPoint(4, 0)},
				{point: NewPoint(4, 4)},
				{point: NewPoint(0, 4)},
			},
			poly2: []polyPoint[int]{
				{point: NewPoint(6, 6)},
				{point: NewPoint(10, 6)},
				{point: NewPoint(10, 10)},
				{point: NewPoint(6, 10)},
			},
			expectedResults: [][]polyPoint[int]{
				{
					{point: NewPoint(0, 0), pointType: PointTypeNormal, entryExit: NotSet, visited: false, otherPolyIndex: 0},
					{point: NewPoint(4, 0), pointType: PointTypeNormal, entryExit: NotSet, visited: false, otherPolyIndex: 0},
					{point: NewPoint(4, 4), pointType: PointTypeNormal, entryExit: NotSet, visited: false, otherPolyIndex: 0},
					{point: NewPoint(0, 4), pointType: PointTypeNormal, entryExit: NotSet, visited: false, otherPolyIndex: 0},
				},
				{
					{point: NewPoint(6, 6), pointType: PointTypeNormal, entryExit: NotSet, visited: false, otherPolyIndex: 0},
					{point: NewPoint(10, 6), pointType: PointTypeNormal, entryExit: NotSet, visited: false, otherPolyIndex: 0},
					{point: NewPoint(10, 10), pointType: PointTypeNormal, entryExit: NotSet, visited: false, otherPolyIndex: 0},
					{point: NewPoint(6, 10), pointType: PointTypeNormal, entryExit: NotSet, visited: false, otherPolyIndex: 0},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			poly1, poly2 := findIntersectionsBetweenPolys(tc.poly1, tc.poly2)
			markEntryExitPoints(poly1, poly2, BooleanUnion)
			results := traverse(poly1, poly2, BooleanUnion)
			assert.Equal(t, tc.expectedResults, results, "Union traversal mismatch")
		})
	}
}

//func TestFindIntersectionsWithOtherPoly(t *testing.T) {
//	tests := map[string]struct {
//		poly1            []polyPoint[int]
//		poly2            []polyPoint[int]
//		expectedPoly1Out []polyPoint[int]
//		expectedPoly2Out []polyPoint[int]
//	}{
//		"LSRMiss #1": {
//			poly1: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(4, 0), pointType: PointTypeNormal},
//				{point: NewPoint(2, 2), pointType: PointTypeNormal},
//			},
//			poly2: []polyPoint[int]{
//				{point: NewPoint(6, 6), pointType: PointTypeNormal},
//				{point: NewPoint(8, 6), pointType: PointTypeNormal},
//				{point: NewPoint(7, 8), pointType: PointTypeNormal},
//			},
//			expectedPoly1Out: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(4, 0), pointType: PointTypeNormal},
//				{point: NewPoint(2, 2), pointType: PointTypeNormal},
//			},
//			expectedPoly2Out: []polyPoint[int]{
//				{point: NewPoint(6, 6), pointType: PointTypeNormal},
//				{point: NewPoint(8, 6), pointType: PointTypeNormal},
//				{point: NewPoint(7, 8), pointType: PointTypeNormal},
//			},
//		},
//		"LSRMiss #2": {
//			poly2: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(4, 0), pointType: PointTypeNormal},
//				{point: NewPoint(2, 2), pointType: PointTypeNormal},
//			},
//			poly1: []polyPoint[int]{
//				{point: NewPoint(6, 6), pointType: PointTypeNormal},
//				{point: NewPoint(8, 6), pointType: PointTypeNormal},
//				{point: NewPoint(7, 8), pointType: PointTypeNormal},
//			},
//			expectedPoly2Out: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(4, 0), pointType: PointTypeNormal},
//				{point: NewPoint(2, 2), pointType: PointTypeNormal},
//			},
//			expectedPoly1Out: []polyPoint[int]{
//				{point: NewPoint(6, 6), pointType: PointTypeNormal},
//				{point: NewPoint(8, 6), pointType: PointTypeNormal},
//				{point: NewPoint(7, 8), pointType: PointTypeNormal},
//			},
//		},
//		"LSRIntersects + LSRAeqC + LSRAeqD": {
//			poly1: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//				{point: NewPoint(0, 5), pointType: PointTypeNormal},
//			},
//			poly2: []polyPoint[int]{
//				{point: NewPoint(3, 3), pointType: PointTypeNormal},
//				{point: NewPoint(8, 3), pointType: PointTypeNormal},
//				{point: NewPoint(8, 8), pointType: PointTypeNormal},
//				{point: NewPoint(3, 8), pointType: PointTypeNormal},
//			},
//			expectedPoly1Out: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 3), pointType: PointTypeAddedIntersection},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//				{point: NewPoint(3, 5), pointType: PointTypeAddedIntersection},
//				{point: NewPoint(0, 5), pointType: PointTypeNormal},
//			},
//			expectedPoly2Out: []polyPoint[int]{
//				{point: NewPoint(3, 5), pointType: PointTypeAddedIntersection},
//				{point: NewPoint(3, 3), pointType: PointTypeNormal},
//				{point: NewPoint(5, 3), pointType: PointTypeAddedIntersection},
//				{point: NewPoint(8, 3), pointType: PointTypeNormal},
//				{point: NewPoint(8, 8), pointType: PointTypeNormal},
//				{point: NewPoint(3, 8), pointType: PointTypeNormal},
//			},
//		},
//		"LSRIntersects + LSRAeqC + LSRAeqD + LSRBeqC + LSRBeqD": {
//			poly2: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//				{point: NewPoint(0, 5), pointType: PointTypeNormal},
//			},
//			poly1: []polyPoint[int]{
//				{point: NewPoint(3, 3), pointType: PointTypeNormal},
//				{point: NewPoint(8, 3), pointType: PointTypeNormal},
//				{point: NewPoint(8, 8), pointType: PointTypeNormal},
//				{point: NewPoint(3, 8), pointType: PointTypeNormal},
//			},
//			expectedPoly2Out: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 3), pointType: PointTypeAddedIntersection},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//				{point: NewPoint(3, 5), pointType: PointTypeAddedIntersection},
//				{point: NewPoint(0, 5), pointType: PointTypeNormal},
//			},
//			expectedPoly1Out: []polyPoint[int]{
//				{point: NewPoint(3, 5), pointType: PointTypeAddedIntersection},
//				{point: NewPoint(3, 3), pointType: PointTypeNormal},
//				{point: NewPoint(5, 3), pointType: PointTypeAddedIntersection},
//				{point: NewPoint(8, 3), pointType: PointTypeNormal},
//				{point: NewPoint(8, 8), pointType: PointTypeNormal},
//				{point: NewPoint(3, 8), pointType: PointTypeNormal},
//			},
//		},
//		"LSRCollinearDisjoint #1": {
//			poly1: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(4, 4), pointType: PointTypeNormal},
//				{point: NewPoint(2, 6), pointType: PointTypeNormal},
//			},
//			poly2: []polyPoint[int]{
//				{point: NewPoint(8, 8), pointType: PointTypeNormal},
//				{point: NewPoint(10, 10), pointType: PointTypeNormal},
//				{point: NewPoint(6, 10), pointType: PointTypeNormal},
//			},
//			expectedPoly1Out: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(4, 4), pointType: PointTypeNormal},
//				{point: NewPoint(2, 6), pointType: PointTypeNormal},
//			},
//			expectedPoly2Out: []polyPoint[int]{
//				{point: NewPoint(8, 8), pointType: PointTypeNormal},
//				{point: NewPoint(10, 10), pointType: PointTypeNormal},
//				{point: NewPoint(6, 10), pointType: PointTypeNormal},
//			},
//		},
//		"LSRCollinearDisjoint #2": {
//			poly2: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(4, 4), pointType: PointTypeNormal},
//				{point: NewPoint(2, 6), pointType: PointTypeNormal},
//			},
//			poly1: []polyPoint[int]{
//				{point: NewPoint(8, 8), pointType: PointTypeNormal},
//				{point: NewPoint(10, 10), pointType: PointTypeNormal},
//				{point: NewPoint(6, 10), pointType: PointTypeNormal},
//			},
//			expectedPoly2Out: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(4, 4), pointType: PointTypeNormal},
//				{point: NewPoint(2, 6), pointType: PointTypeNormal},
//			},
//			expectedPoly1Out: []polyPoint[int]{
//				{point: NewPoint(8, 8), pointType: PointTypeNormal},
//				{point: NewPoint(10, 10), pointType: PointTypeNormal},
//				{point: NewPoint(6, 10), pointType: PointTypeNormal},
//			},
//		},
//		"LSRConAB + LSRAeqC + LSRDonAB + LSRAeqD": {
//			poly1: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//				{point: NewPoint(0, 5), pointType: PointTypeNormal},
//			},
//			poly2: []polyPoint[int]{
//				{point: NewPoint(5, 1), pointType: PointTypeNormal},
//				{point: NewPoint(9, 1), pointType: PointTypeNormal},
//				{point: NewPoint(9, 4), pointType: PointTypeNormal},
//				{point: NewPoint(5, 4), pointType: PointTypeNormal},
//			},
//			expectedPoly1Out: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 1), pointType: PointTypeAddedIntersection},
//				{point: NewPoint(5, 4), pointType: PointTypeAddedIntersection},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//				{point: NewPoint(0, 5), pointType: PointTypeNormal},
//			},
//			expectedPoly2Out: []polyPoint[int]{
//				{point: NewPoint(5, 1), pointType: PointTypeIntersection},
//				{point: NewPoint(9, 1), pointType: PointTypeNormal},
//				{point: NewPoint(9, 4), pointType: PointTypeNormal},
//				{point: NewPoint(5, 4), pointType: PointTypeIntersection},
//			},
//		},
//		"LSRAonCD + LSRAeqC + LSRBonCD + LSRBeqC + LSRBeqD + LSRCollinearEqual": {
//			poly2: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//				{point: NewPoint(0, 5), pointType: PointTypeNormal},
//			},
//			poly1: []polyPoint[int]{
//				{point: NewPoint(5, 1), pointType: PointTypeNormal},
//				{point: NewPoint(9, 1), pointType: PointTypeNormal},
//				{point: NewPoint(9, 4), pointType: PointTypeNormal},
//				{point: NewPoint(5, 4), pointType: PointTypeNormal},
//			},
//			expectedPoly2Out: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 1), pointType: PointTypeAddedIntersection},
//				{point: NewPoint(5, 4), pointType: PointTypeAddedIntersection},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//				{point: NewPoint(0, 5), pointType: PointTypeNormal},
//			},
//			expectedPoly1Out: []polyPoint[int]{
//				{point: NewPoint(5, 1), pointType: PointTypeIntersection},
//				{point: NewPoint(9, 1), pointType: PointTypeNormal},
//				{point: NewPoint(9, 4), pointType: PointTypeNormal},
//				{point: NewPoint(5, 4), pointType: PointTypeIntersection},
//			},
//		},
//		"LSRBonCD + LSRBeqC + LSRAeqD + LSRAeqC + LSRDonAB + LSRCollinearDisjoint": {
//			poly1: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//				{point: NewPoint(0, 5), pointType: PointTypeNormal},
//			},
//			poly2: []polyPoint[int]{
//				{point: NewPoint(2, 5), pointType: PointTypeNormal},
//				{point: NewPoint(8, 5), pointType: PointTypeNormal},
//				{point: NewPoint(8, 9), pointType: PointTypeNormal},
//				{point: NewPoint(2, 9), pointType: PointTypeNormal},
//			},
//			expectedPoly1Out: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeIntersection},
//				{point: NewPoint(2, 5), pointType: PointTypeAddedIntersection},
//				{point: NewPoint(0, 5), pointType: PointTypeNormal},
//			},
//			expectedPoly2Out: []polyPoint[int]{
//				{point: NewPoint(2, 5), pointType: PointTypeIntersection},
//				{point: NewPoint(5, 5), pointType: PointTypeAddedIntersection},
//				{point: NewPoint(8, 5), pointType: PointTypeNormal},
//				{point: NewPoint(8, 9), pointType: PointTypeNormal},
//				{point: NewPoint(2, 9), pointType: PointTypeNormal},
//			},
//		},
//		"LSRDonAB + LSRAeqD + LSRAeqC + LSRBonCD + LSRBeqC": {
//			poly2: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//				{point: NewPoint(0, 5), pointType: PointTypeNormal},
//			},
//			poly1: []polyPoint[int]{
//				{point: NewPoint(2, 5), pointType: PointTypeNormal},
//				{point: NewPoint(8, 5), pointType: PointTypeNormal},
//				{point: NewPoint(8, 9), pointType: PointTypeNormal},
//				{point: NewPoint(2, 9), pointType: PointTypeNormal},
//			},
//			expectedPoly2Out: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeIntersection},
//				{point: NewPoint(2, 5), pointType: PointTypeAddedIntersection},
//				{point: NewPoint(0, 5), pointType: PointTypeNormal},
//			},
//			expectedPoly1Out: []polyPoint[int]{
//				{point: NewPoint(2, 5), pointType: PointTypeIntersection},
//				{point: NewPoint(5, 5), pointType: PointTypeAddedIntersection},
//				{point: NewPoint(8, 5), pointType: PointTypeNormal},
//				{point: NewPoint(8, 9), pointType: PointTypeNormal},
//				{point: NewPoint(2, 9), pointType: PointTypeNormal},
//			},
//		},
//		"LSRBeqC + LSRBeqD + LSRAeqC + LSRCollinearEqual + LSRAeqD": {
//			poly1: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//				{point: NewPoint(0, 5), pointType: PointTypeNormal},
//			},
//			poly2: []polyPoint[int]{
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(10, 0), pointType: PointTypeNormal},
//				{point: NewPoint(10, 5), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//			},
//			expectedPoly1Out: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeIntersection},
//				{point: NewPoint(5, 5), pointType: PointTypeIntersection},
//				{point: NewPoint(0, 5), pointType: PointTypeNormal},
//			},
//			expectedPoly2Out: []polyPoint[int]{
//				{point: NewPoint(5, 0), pointType: PointTypeIntersection},
//				{point: NewPoint(10, 0), pointType: PointTypeNormal},
//				{point: NewPoint(10, 5), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeIntersection},
//			},
//		},
//		"LSRAeqD + LSRAeqC + LSRBeqD + LSRBeqC + LSRCollinearEqual": {
//			poly2: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//				{point: NewPoint(0, 5), pointType: PointTypeNormal},
//			},
//			poly1: []polyPoint[int]{
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(10, 0), pointType: PointTypeNormal},
//				{point: NewPoint(10, 5), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//			},
//			expectedPoly2Out: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeIntersection},
//				{point: NewPoint(5, 5), pointType: PointTypeIntersection},
//				{point: NewPoint(0, 5), pointType: PointTypeNormal},
//			},
//			expectedPoly1Out: []polyPoint[int]{
//				{point: NewPoint(5, 0), pointType: PointTypeIntersection},
//				{point: NewPoint(10, 0), pointType: PointTypeNormal},
//				{point: NewPoint(10, 5), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeIntersection},
//			},
//		},
//		"LSRCollinearCDinAB + LSRCollinearEqual + LSRBeqC + LSRAeqD + LSRAeqC": {
//			poly1: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//				{point: NewPoint(0, 5), pointType: PointTypeNormal},
//			},
//			poly2: []polyPoint[int]{
//				{point: NewPoint(1, 0), pointType: PointTypeNormal},
//				{point: NewPoint(4, 0), pointType: PointTypeNormal},
//				{point: NewPoint(4, -3), pointType: PointTypeNormal},
//				{point: NewPoint(1, -3), pointType: PointTypeNormal},
//			},
//			expectedPoly1Out: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(1, 0), pointType: PointTypeAddedIntersection},
//				{point: NewPoint(4, 0), pointType: PointTypeAddedIntersection},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//				{point: NewPoint(0, 5), pointType: PointTypeNormal},
//			},
//			expectedPoly2Out: []polyPoint[int]{
//				{point: NewPoint(1, 0), pointType: PointTypeIntersection},
//				{point: NewPoint(4, 0), pointType: PointTypeIntersection},
//				{point: NewPoint(4, -3), pointType: PointTypeNormal},
//				{point: NewPoint(1, -3), pointType: PointTypeNormal},
//			},
//		},
//		"LSRCollinearABinCD + LSRCollinearEqual + LSRBeqC + LSRAeqD + LSRAeqC + LSRBeqD + LSRBeqC": {
//			poly2: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//				{point: NewPoint(0, 5), pointType: PointTypeNormal},
//			},
//			poly1: []polyPoint[int]{
//				{point: NewPoint(1, 0), pointType: PointTypeNormal},
//				{point: NewPoint(4, 0), pointType: PointTypeNormal},
//				{point: NewPoint(4, -3), pointType: PointTypeNormal},
//				{point: NewPoint(1, -3), pointType: PointTypeNormal},
//			},
//			expectedPoly2Out: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(1, 0), pointType: PointTypeAddedIntersection},
//				{point: NewPoint(4, 0), pointType: PointTypeAddedIntersection},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//				{point: NewPoint(0, 5), pointType: PointTypeNormal},
//			},
//			expectedPoly1Out: []polyPoint[int]{
//				{point: NewPoint(1, 0), pointType: PointTypeIntersection},
//				{point: NewPoint(4, 0), pointType: PointTypeIntersection},
//				{point: NewPoint(4, -3), pointType: PointTypeNormal},
//				{point: NewPoint(1, -3), pointType: PointTypeNormal},
//			},
//		},
//		"LSRCollinearAonCD + LSRAeqC + LSRConAB + LSRCollinearDisjoint + LSRAeqD + LSRAeqC + LSRBeqD + LSRBeqC": {
//			poly1: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//				{point: NewPoint(0, 5), pointType: PointTypeNormal},
//			},
//			poly2: []polyPoint[int]{
//				{point: NewPoint(-2, 0), pointType: PointTypeNormal},
//				{point: NewPoint(2, 0), pointType: PointTypeNormal},
//				{point: NewPoint(2, -3), pointType: PointTypeNormal},
//				{point: NewPoint(-2, -3), pointType: PointTypeNormal},
//			},
//			expectedPoly1Out: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeIntersection},
//				{point: NewPoint(2, 0), pointType: PointTypeAddedIntersection},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//				{point: NewPoint(0, 5), pointType: PointTypeNormal},
//			},
//			expectedPoly2Out: []polyPoint[int]{
//				{point: NewPoint(-2, 0), pointType: PointTypeNormal},
//				{point: NewPoint(0, 0), pointType: PointTypeAddedIntersection},
//				{point: NewPoint(2, 0), pointType: PointTypeIntersection},
//				{point: NewPoint(2, -3), pointType: PointTypeNormal},
//				{point: NewPoint(-2, -3), pointType: PointTypeNormal},
//			},
//		},
//		"LSRCollinearBonCD + LSRBeqC + LSRDonAB + LSRCollinearEqual + LSRAeqD + LSRAeqC": {
//			poly1: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//				{point: NewPoint(0, 5), pointType: PointTypeNormal},
//			},
//			poly2: []polyPoint[int]{
//				{point: NewPoint(3, 0), pointType: PointTypeNormal},
//				{point: NewPoint(8, 0), pointType: PointTypeNormal},
//				{point: NewPoint(8, -3), pointType: PointTypeNormal},
//				{point: NewPoint(3, -3), pointType: PointTypeNormal},
//			},
//			expectedPoly1Out: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(3, 0), pointType: PointTypeAddedIntersection},
//				{point: NewPoint(5, 0), pointType: PointTypeIntersection},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//				{point: NewPoint(0, 5), pointType: PointTypeNormal},
//			},
//			expectedPoly2Out: []polyPoint[int]{
//				{point: NewPoint(3, 0), pointType: PointTypeIntersection},
//				{point: NewPoint(5, 0), pointType: PointTypeAddedIntersection},
//				{point: NewPoint(8, 0), pointType: PointTypeNormal},
//				{point: NewPoint(8, -3), pointType: PointTypeNormal},
//				{point: NewPoint(3, -3), pointType: PointTypeNormal},
//			},
//		},
//		"LSRCollinearEqual + LSRBeqC + LSRAeqD + LSRAeqC + LSRBeqD": {
//			poly1: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//				{point: NewPoint(0, 5), pointType: PointTypeNormal},
//			},
//			poly2: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, -5), pointType: PointTypeNormal},
//				{point: NewPoint(0, -5), pointType: PointTypeNormal},
//			},
//			expectedPoly1Out: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeIntersection},
//				{point: NewPoint(5, 0), pointType: PointTypeIntersection},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//				{point: NewPoint(0, 5), pointType: PointTypeNormal},
//			},
//			expectedPoly2Out: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeIntersection},
//				{point: NewPoint(5, 0), pointType: PointTypeIntersection},
//				{point: NewPoint(5, -5), pointType: PointTypeNormal},
//				{point: NewPoint(0, -5), pointType: PointTypeNormal},
//			},
//		},
//		"Nested polygon with no intersection": {
//			poly1: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(10, 0), pointType: PointTypeNormal},
//				{point: NewPoint(10, 10), pointType: PointTypeNormal},
//				{point: NewPoint(0, 10), pointType: PointTypeNormal},
//			},
//			poly2: []polyPoint[int]{
//				{point: NewPoint(3, 3), pointType: PointTypeNormal},
//				{point: NewPoint(7, 3), pointType: PointTypeNormal},
//				{point: NewPoint(7, 7), pointType: PointTypeNormal},
//				{point: NewPoint(3, 7), pointType: PointTypeNormal},
//			},
//			expectedPoly1Out: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(10, 0), pointType: PointTypeNormal},
//				{point: NewPoint(10, 10), pointType: PointTypeNormal},
//				{point: NewPoint(0, 10), pointType: PointTypeNormal},
//			},
//			expectedPoly2Out: []polyPoint[int]{
//				{point: NewPoint(3, 3), pointType: PointTypeNormal},
//				{point: NewPoint(7, 3), pointType: PointTypeNormal},
//				{point: NewPoint(7, 7), pointType: PointTypeNormal},
//				{point: NewPoint(3, 7), pointType: PointTypeNormal},
//			},
//		},
//		"Touching polygon corners": {
//			poly1: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//				{point: NewPoint(0, 5), pointType: PointTypeNormal},
//			},
//			poly2: []polyPoint[int]{
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//				{point: NewPoint(10, 5), pointType: PointTypeNormal},
//				{point: NewPoint(10, 10), pointType: PointTypeNormal},
//				{point: NewPoint(5, 10), pointType: PointTypeNormal},
//			},
//			expectedPoly1Out: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeIntersection},
//				{point: NewPoint(0, 5), pointType: PointTypeNormal},
//			},
//			expectedPoly2Out: []polyPoint[int]{
//				{point: NewPoint(5, 5), pointType: PointTypeIntersection},
//				{point: NewPoint(10, 5), pointType: PointTypeNormal},
//				{point: NewPoint(10, 10), pointType: PointTypeNormal},
//				{point: NewPoint(5, 10), pointType: PointTypeNormal},
//			},
//		},
//		"Fully overlapping polygons": {
//			poly1: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//				{point: NewPoint(0, 5), pointType: PointTypeNormal},
//			},
//			poly2: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//				{point: NewPoint(0, 5), pointType: PointTypeNormal},
//			},
//			expectedPoly1Out: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeIntersection},
//				{point: NewPoint(5, 0), pointType: PointTypeIntersection},
//				{point: NewPoint(5, 5), pointType: PointTypeIntersection},
//				{point: NewPoint(0, 5), pointType: PointTypeIntersection},
//			},
//			expectedPoly2Out: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeIntersection},
//				{point: NewPoint(5, 0), pointType: PointTypeIntersection},
//				{point: NewPoint(5, 5), pointType: PointTypeIntersection},
//				{point: NewPoint(0, 5), pointType: PointTypeIntersection},
//			},
//		},
//		"Touching edges without overlap": {
//			poly1: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//				{point: NewPoint(0, 5), pointType: PointTypeNormal},
//			},
//			poly2: []polyPoint[int]{
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(10, 0), pointType: PointTypeNormal},
//				{point: NewPoint(10, 5), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//			},
//			expectedPoly1Out: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeIntersection},
//				{point: NewPoint(5, 5), pointType: PointTypeIntersection},
//				{point: NewPoint(0, 5), pointType: PointTypeNormal},
//			},
//			expectedPoly2Out: []polyPoint[int]{
//				{point: NewPoint(5, 0), pointType: PointTypeIntersection},
//				{point: NewPoint(10, 0), pointType: PointTypeNormal},
//				{point: NewPoint(10, 5), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeIntersection},
//			},
//		},
//		"Concave polygon intersection": {
//			poly1: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(6, 0), pointType: PointTypeNormal},
//				{point: NewPoint(3, 3), pointType: PointTypeNormal},
//				{point: NewPoint(6, 6), pointType: PointTypeNormal},
//				{point: NewPoint(0, 6), pointType: PointTypeNormal},
//			},
//			poly2: []polyPoint[int]{
//				{point: NewPoint(2, 1), pointType: PointTypeNormal},
//				{point: NewPoint(5, 1), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//				{point: NewPoint(2, 5), pointType: PointTypeNormal},
//			},
//			expectedPoly1Out: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(6, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 1), pointType: PointTypeAddedIntersection},
//				{point: NewPoint(3, 3), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeAddedIntersection},
//				{point: NewPoint(6, 6), pointType: PointTypeNormal},
//				{point: NewPoint(0, 6), pointType: PointTypeNormal},
//			},
//			expectedPoly2Out: []polyPoint[int]{
//				{point: NewPoint(2, 1), pointType: PointTypeNormal},
//				{point: NewPoint(5, 1), pointType: PointTypeIntersection},
//				{point: NewPoint(5, 5), pointType: PointTypeIntersection},
//				{point: NewPoint(2, 5), pointType: PointTypeNormal},
//			},
//		},
//	}
//	for name, tc := range tests {
//		t.Run(name, func(t *testing.T) {
//			poly1Out, poly2Out := findIntersectionsWithOtherPoly(tc.poly1, tc.poly2)
//			assert.Equal(t, tc.expectedPoly1Out, poly1Out)
//			assert.Equal(t, tc.expectedPoly2Out, poly2Out)
//		})
//	}
//}

//func TestFindIntersections_PanicsOnInvalidInput(t *testing.T) {
//	tests := map[string]struct {
//		poly1 []polyPoint[int]
//		poly2 []polyPoint[int]
//	}{
//		"poly1 with less than three points": {
//			poly1: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(1, 1), pointType: PointTypeNormal},
//			},
//			poly2: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//			},
//		},
//		"poly2 with less than three points": {
//			poly1: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//			},
//			poly2: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(1, 1), pointType: PointTypeNormal},
//			},
//		},
//		"poly1 with zero area": {
//			poly1: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(10, 0), pointType: PointTypeNormal},
//			},
//			poly2: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//			},
//		},
//		"poly2 with zero area": {
//			poly1: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 5), pointType: PointTypeNormal},
//			},
//			poly2: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(10, 0), pointType: PointTypeNormal},
//			},
//		},
//		"both poly1 and poly2 with zero area": {
//			poly1: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(10, 0), pointType: PointTypeNormal},
//			},
//			poly2: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(10, 0), pointType: PointTypeNormal},
//			},
//		},
//	}
//
//	for name, tc := range tests {
//		t.Run(name, func(t *testing.T) {
//			assert.Panics(t, func() {
//				findIntersectionsWithOtherPoly(tc.poly1, tc.poly2)
//			}, "Expected findIntersectionsWithOtherPoly to panic on invalid input")
//		})
//	}
//}

//func TestInsertIntersectionIfNotDuplicate(t *testing.T) {
//	tests := map[string]struct {
//		initialPoly    []polyPoint[int]
//		index          int
//		pointToInsert  polyPoint[int]
//		expectedOutput []polyPoint[int]
//	}{
//		"Insert at end without duplicate": {
//			initialPoly: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//			},
//			index: 2,
//			pointToInsert: polyPoint[int]{
//				point:     NewPoint(10, 0),
//				pointType: PointTypeIntersection,
//			},
//			expectedOutput: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(10, 0), pointType: PointTypeIntersection},
//			},
//		},
//		"Insert in middle without duplicate": {
//			initialPoly: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(10, 0), pointType: PointTypeNormal},
//			},
//			index: 1,
//			pointToInsert: polyPoint[int]{
//				point:     NewPoint(5, 0),
//				pointType: PointTypeIntersection,
//			},
//			expectedOutput: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeIntersection},
//				{point: NewPoint(10, 0), pointType: PointTypeNormal},
//			},
//		},
//		"Skip insertion due to existing duplicate with Intersection type": {
//			initialPoly: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeIntersection},
//				{point: NewPoint(10, 0), pointType: PointTypeNormal},
//			},
//			index: 1,
//			pointToInsert: polyPoint[int]{
//				point:     NewPoint(5, 0),
//				pointType: PointTypeIntersection,
//			},
//			expectedOutput: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeNormal},
//				{point: NewPoint(5, 0), pointType: PointTypeIntersection},
//				{point: NewPoint(10, 0), pointType: PointTypeNormal},
//			},
//		},
//		"Insert at beginning without duplicate": {
//			initialPoly: []polyPoint[int]{
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(10, 0), pointType: PointTypeNormal},
//			},
//			index: 0,
//			pointToInsert: polyPoint[int]{
//				point:     NewPoint(0, 0),
//				pointType: PointTypeIntersection,
//			},
//			expectedOutput: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeIntersection},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(10, 0), pointType: PointTypeNormal},
//			},
//		},
//		"Insert in middle with different point type (Normal)": {
//			initialPoly: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeIntersection},
//				{point: NewPoint(10, 0), pointType: PointTypeNormal},
//			},
//			index: 1,
//			pointToInsert: polyPoint[int]{
//				point:     NewPoint(5, 0),
//				pointType: PointTypeNormal, // Insert even though it's the same point as others with a different type
//			},
//			expectedOutput: []polyPoint[int]{
//				{point: NewPoint(0, 0), pointType: PointTypeIntersection},
//				{point: NewPoint(5, 0), pointType: PointTypeNormal},
//				{point: NewPoint(10, 0), pointType: PointTypeNormal},
//			},
//		},
//	}
//
//	for name, tc := range tests {
//		t.Run(name, func(t *testing.T) {
//			poly := make([]polyPoint[int], len(tc.initialPoly))
//			copy(poly, tc.initialPoly)
//
//			insertIntersectionIfNotDuplicate(&poly, tc.index, tc.pointToInsert)
//
//			assert.Equal(t, tc.expectedOutput, poly)
//		})
//	}
//}

func TestNewPolygon(t *testing.T) {
	tests := map[string]struct {
		points        []Point[int]
		polygonType   PolygonType
		expectedError bool
		validateFunc  func(t *testing.T, polygon *Polygon[int])
	}{
		"Valid Solid Polygon": {
			points: []Point[int]{
				NewPoint(0, 0),
				NewPoint(2, 0),
				NewPoint(1, 2),
			},
			polygonType:   Solid,
			expectedError: false,
			validateFunc: func(t *testing.T, polygon *Polygon[int]) {
				assert.Equal(t, Solid, polygon.polygonType, "Expected polygon type to be Solid")
				assert.Len(t, polygon.Points(), 3, "Expected three points")
				assert.NotNil(t, polygon.hull, "Expected convex hull to be created")
			},
		},
		"Invalid Fewer Than Three Points": {
			points: []Point[int]{
				NewPoint(0, 0),
				NewPoint(1, 1),
			},
			polygonType:   Solid,
			expectedError: true,
		},
		"Invalid Zero Area": {
			points: []Point[int]{
				NewPoint(0, 0),
				NewPoint(1, 1),
				NewPoint(2, 2),
			},
			polygonType:   Solid,
			expectedError: true,
		},
		"Orientation Correction for Solid Polygon": {
			points: []Point[int]{
				NewPoint(0, 0),
				NewPoint(1, 2),
				NewPoint(2, 0),
			},
			polygonType:   Solid,
			expectedError: false,
			validateFunc: func(t *testing.T, polygon *Polygon[int]) {
				assert.Greater(t, SignedArea2X(polygon.Points()), 0, "Expected counter-clockwise orientation for solid polygon")
			},
		},
		"Orientation Correction for Hole Polygon": {
			points: []Point[int]{
				NewPoint(0, 0),
				NewPoint(1, 2),
				NewPoint(2, 0),
			},
			polygonType:   Hole,
			expectedError: false,
			validateFunc: func(t *testing.T, polygon *Polygon[int]) {
				assert.Less(t, SignedArea2X(polygon.Points()), 0, "Expected clockwise orientation for hole polygon")
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			polygon, err := NewPolygon(tc.points, tc.polygonType)

			if tc.expectedError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			if tc.validateFunc != nil {
				tc.validateFunc(t, polygon)
			}
		})
	}
}

func TestPolygon_Points(t *testing.T) {
	tests := map[string]struct {
		points         []Point[int]
		pointTypes     []polyPointType
		expectedPoints []Point[int]
	}{
		"Normal points only": {
			points: []Point[int]{
				NewPoint(0, 0),
				NewPoint(5, 0),
				NewPoint(5, 5),
				NewPoint(5, 5),
			},
			pointTypes: []polyPointType{
				PointTypeNormal, PointTypeNormal, PointTypeNormal, PointTypeNormal,
			},
			expectedPoints: []Point[int]{
				NewPoint(0, 0),
				NewPoint(5, 0),
				NewPoint(5, 5),
				NewPoint(5, 5),
			},
		},
		"Mixed point types": {
			points: []Point[int]{
				NewPoint(0, 0),
				NewPoint(5, 0),
				NewPoint(5, 2),
				NewPoint(5, 5),
				NewPoint(5, 5),
			},
			pointTypes: []polyPointType{
				PointTypeNormal, PointTypeNormal, PointTypeAddedIntersection, PointTypeIntersection, PointTypeNormal,
			},
			expectedPoints: []Point[int]{
				NewPoint(0, 0),
				NewPoint(5, 0),
				NewPoint(5, 5),
				NewPoint(5, 5),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Create the polygon using NewPolygon, then modify the point types
			polygon, err := NewPolygon(tc.points, Solid)
			assert.NoError(t, err)

			// Set point types according to the test case
			for i := range tc.pointTypes {
				polygon.points[i].pointType = tc.pointTypes[i]
			}

			// Check the result from Points method
			points := polygon.Points()
			assert.Equal(t, tc.expectedPoints, points)
		})
	}
}

func TestPolygon_RelationshipToPoint(t *testing.T) {
	tests := map[string]struct {
		polygon     Polygon[int]
		point       Point[int]
		expectedRel PointPolygonRelationship
	}{
		"Point Outside Polygon": {
			polygon: func() Polygon[int] {
				poly, _ := NewPolygon([]Point[int]{
					NewPoint(0, 0),
					NewPoint(4, 0),
					NewPoint(4, 4),
					NewPoint(0, 4),
				}, Solid)
				return *poly
			}(),
			point:       NewPoint(5, 5),
			expectedRel: PointOutside,
		},
		"Point Inside Polygon": {
			polygon: func() Polygon[int] {
				poly, _ := NewPolygon([]Point[int]{
					NewPoint(0, 0),
					NewPoint(4, 0),
					NewPoint(4, 4),
					NewPoint(0, 4),
				}, Solid)
				return *poly
			}(),
			point:       NewPoint(2, 2),
			expectedRel: PointInside,
		},
		"Point On Vertex": {
			polygon: func() Polygon[int] {
				poly, _ := NewPolygon([]Point[int]{
					NewPoint(0, 0),
					NewPoint(4, 0),
					NewPoint(4, 4),
					NewPoint(0, 4),
				}, Solid)
				return *poly
			}(),
			point:       NewPoint(4, 0),
			expectedRel: PointOnVertex,
		},
		"Point On Edge": {
			polygon: func() Polygon[int] {
				poly, _ := NewPolygon([]Point[int]{
					NewPoint(0, 0),
					NewPoint(4, 0),
					NewPoint(4, 4),
					NewPoint(0, 4),
				}, Solid)
				return *poly
			}(),
			point:       NewPoint(2, 0),
			expectedRel: PointOnEdge,
		},
		"Point in Concave Indent": {
			polygon: func() Polygon[int] {
				poly, _ := NewPolygon([]Point[int]{
					NewPoint(0, 0),
					NewPoint(4, 0),
					NewPoint(2, 2), // Creates the indent in the polygon
					NewPoint(4, 4),
					NewPoint(0, 4),
				}, Solid)
				return *poly
			}(),
			point:       NewPoint(3, 2), // Point in the indent but outside the polygon
			expectedRel: PointOutside,
		},
		"Point in Concave Polygon": {
			polygon: func() Polygon[int] {
				poly, _ := NewPolygon([]Point[int]{
					NewPoint(0, 0),
					NewPoint(4, 0),
					NewPoint(2, 2), // Creates the indent in the polygon
					NewPoint(4, 4),
					NewPoint(0, 4),
				}, Solid)
				return *poly
			}(),
			point:       NewPoint(2, 3),
			expectedRel: PointInside,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Compute the relationship of the point to the polygon
			relationship := tc.polygon.RelationshipToPoint(tc.point)

			// Assert that the relationship matches the expected result
			assert.Equal(t, tc.expectedRel, relationship, "Expected relationship %v, got %v", tc.expectedRel, relationship)
		})
	}
}

func TestPolygon_RelationshipToPoint_WithHolesAndIslands(t *testing.T) {
	tests := map[string]struct {
		polygon     Polygon[int]
		point       Point[int]
		expectedRel PointPolygonRelationship
	}{
		"Point Inside Concave Polygon Outside Hole": {
			polygon: func() Polygon[int] {
				mainPoly, _ := NewPolygon([]Point[int]{
					NewPoint(0, 0),
					NewPoint(12, 0),
					NewPoint(12, 12),
					NewPoint(0, 12),
				}, Solid)

				hole, _ := NewPolygon([]Point[int]{
					NewPoint(4, 4),
					NewPoint(8, 4),
					NewPoint(8, 8),
					NewPoint(4, 8),
				}, Hole)

				mainPoly.children = []*Polygon[int]{hole}
				return *mainPoly
			}(),
			point:       NewPoint(2, 2),
			expectedRel: PointInside,
		},
		"Point Inside Hole": {
			polygon: func() Polygon[int] {
				mainPoly, _ := NewPolygon([]Point[int]{
					NewPoint(0, 0),
					NewPoint(12, 0),
					NewPoint(12, 12),
					NewPoint(0, 12),
				}, Solid)

				hole, _ := NewPolygon([]Point[int]{
					NewPoint(4, 4),
					NewPoint(8, 4),
					NewPoint(8, 8),
					NewPoint(4, 8),
				}, Hole)

				mainPoly.children = []*Polygon[int]{hole}
				return *mainPoly
			}(),
			point:       NewPoint(6, 6),
			expectedRel: PointInHole,
		},
		"Point Inside Island": {
			polygon: func() Polygon[int] {
				mainPoly, _ := NewPolygon([]Point[int]{
					NewPoint(0, 0),
					NewPoint(12, 0),
					NewPoint(12, 12),
					NewPoint(0, 12),
				}, Solid)

				hole, _ := NewPolygon([]Point[int]{
					NewPoint(4, 4),
					NewPoint(8, 4),
					NewPoint(8, 8),
					NewPoint(4, 8),
				}, Hole)

				island, _ := NewPolygon([]Point[int]{
					NewPoint(5, 5),
					NewPoint(7, 5),
					NewPoint(7, 7),
					NewPoint(5, 7),
				}, Solid)

				hole.children = []*Polygon[int]{island}
				mainPoly.children = []*Polygon[int]{hole}
				return *mainPoly
			}(),
			point:       NewPoint(6, 6),
			expectedRel: PointInsideIsland,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Compute the relationship of the point to the polygon
			relationship := tc.polygon.RelationshipToPoint(tc.point)

			// Assert that the relationship matches the expected result
			assert.Equal(t, tc.expectedRel, relationship, "Expected relationship %v, got %v", tc.expectedRel, relationship)
		})
	}
}

func TestPolygon_RelationshipToPoint_WithDeepNesting(t *testing.T) {
	tests := map[string]struct {
		polygon     Polygon[int]
		point       Point[int]
		expectedRel PointPolygonRelationship
	}{
		"Point Deeply Nested in Island": {
			polygon: func() Polygon[int] {
				mainPoly, _ := NewPolygon([]Point[int]{
					NewPoint(0, 0),
					NewPoint(20, 0),
					NewPoint(20, 20),
					NewPoint(0, 20),
				}, Solid)

				hole1, _ := NewPolygon([]Point[int]{
					NewPoint(5, 5),
					NewPoint(15, 5),
					NewPoint(15, 15),
					NewPoint(5, 15),
				}, Hole)

				island1, _ := NewPolygon([]Point[int]{
					NewPoint(6, 6),
					NewPoint(14, 6),
					NewPoint(14, 14),
					NewPoint(6, 14),
				}, Solid)

				hole2, _ := NewPolygon([]Point[int]{
					NewPoint(8, 8),
					NewPoint(12, 8),
					NewPoint(12, 12),
					NewPoint(8, 12),
				}, Hole)

				island2, _ := NewPolygon([]Point[int]{
					NewPoint(9, 9),
					NewPoint(11, 9),
					NewPoint(11, 11),
					NewPoint(9, 11),
				}, Solid)

				// Nest the polygons
				hole1.children = []*Polygon[int]{island1}
				island1.children = []*Polygon[int]{hole2}
				hole2.children = []*Polygon[int]{island2}
				mainPoly.children = []*Polygon[int]{hole1}
				return *mainPoly
			}(),
			point:       NewPoint(10, 10),
			expectedRel: PointInsideIsland,
		},
		"Point Deeply Nested in Hole": {
			polygon: func() Polygon[int] {
				mainPoly, _ := NewPolygon([]Point[int]{
					NewPoint(0, 0),
					NewPoint(20, 0),
					NewPoint(20, 20),
					NewPoint(0, 20),
				}, Solid)

				hole1, _ := NewPolygon([]Point[int]{
					NewPoint(5, 5),
					NewPoint(15, 5),
					NewPoint(15, 15),
					NewPoint(5, 15),
				}, Hole)

				island1, _ := NewPolygon([]Point[int]{
					NewPoint(6, 6),
					NewPoint(14, 6),
					NewPoint(14, 14),
					NewPoint(6, 14),
				}, Solid)

				hole2, _ := NewPolygon([]Point[int]{
					NewPoint(8, 8),
					NewPoint(12, 8),
					NewPoint(12, 12),
					NewPoint(8, 12),
				}, Hole)

				// Nest the polygons
				hole1.children = []*Polygon[int]{island1}
				island1.children = []*Polygon[int]{hole2}
				mainPoly.children = []*Polygon[int]{hole1}
				return *mainPoly
			}(),
			point:       NewPoint(10, 10),
			expectedRel: PointInHole,
		},
		"Point in Outer Polygon, Not in Deepest Children": {
			polygon: func() Polygon[int] {
				mainPoly, _ := NewPolygon([]Point[int]{
					NewPoint(0, 0),
					NewPoint(20, 0),
					NewPoint(20, 20),
					NewPoint(0, 20),
				}, Solid)

				hole1, _ := NewPolygon([]Point[int]{
					NewPoint(5, 5),
					NewPoint(15, 5),
					NewPoint(15, 15),
					NewPoint(5, 15),
				}, Hole)

				island1, _ := NewPolygon([]Point[int]{
					NewPoint(6, 6),
					NewPoint(14, 6),
					NewPoint(14, 14),
					NewPoint(6, 14),
				}, Solid)

				hole2, _ := NewPolygon([]Point[int]{
					NewPoint(8, 8),
					NewPoint(12, 8),
					NewPoint(12, 12),
					NewPoint(8, 12),
				}, Hole)

				// Nest the polygons
				hole1.children = []*Polygon[int]{island1}
				island1.children = []*Polygon[int]{hole2}
				mainPoly.children = []*Polygon[int]{hole1}
				return *mainPoly
			}(),
			point:       NewPoint(3, 3),
			expectedRel: PointInside,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Compute the relationship of the point to the polygon
			relationship := tc.polygon.RelationshipToPoint(tc.point)

			// Assert that the relationship matches the expected result
			assert.Equal(t, tc.expectedRel, relationship, "Expected relationship %v, got %v", tc.expectedRel, relationship)
		})
	}
}

func TestSimpleConvexPolygon_ContainsPoint(t *testing.T) {
	tests := map[string]struct {
		polygon         simpleConvexPolygon[int]
		point           Point[int]
		expectedContain bool
	}{
		"Point Inside Polygon": {
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
		"Point Outside Polygon": {
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
