package geom2d

import (
	"fmt"
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
				{point: NewPoint(4, 2), pointType: PointTypeAddedIntersection, entryExit: intersectionTypeExit, otherPolyIndex: 2},
				{point: NewPoint(4, 4), pointType: PointTypeNormal},
				{point: NewPoint(2, 4), pointType: PointTypeAddedIntersection, entryExit: intersectionTypeEntry, otherPolyIndex: 0},
				{point: NewPoint(0, 4), pointType: PointTypeNormal},
			},
			expectedPoly2: []polyPoint[int]{
				{point: NewPoint(2, 4), pointType: PointTypeAddedIntersection, entryExit: intersectionTypeExit, otherPolyIndex: 4},
				{point: NewPoint(2, 2), pointType: PointTypeNormal},
				{point: NewPoint(4, 2), pointType: PointTypeAddedIntersection, entryExit: intersectionTypeEntry, otherPolyIndex: 2},
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
				{point: NewPoint(4, 2), pointType: PointTypeAddedIntersection, entryExit: intersectionTypeEntry, otherPolyIndex: 2},
				{point: NewPoint(4, 4), pointType: PointTypeNormal},
				{point: NewPoint(2, 4), pointType: PointTypeAddedIntersection, entryExit: intersectionTypeExit, otherPolyIndex: 0},
				{point: NewPoint(0, 4), pointType: PointTypeNormal},
			},
			expectedPoly2: []polyPoint[int]{
				{point: NewPoint(2, 4), pointType: PointTypeAddedIntersection, entryExit: intersectionTypeEntry, otherPolyIndex: 4},
				{point: NewPoint(2, 2), pointType: PointTypeNormal},
				{point: NewPoint(4, 2), pointType: PointTypeAddedIntersection, entryExit: intersectionTypeExit, otherPolyIndex: 2},
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
				{point: NewPoint(4, 2), pointType: PointTypeAddedIntersection, entryExit: intersectionTypeExit, otherPolyIndex: 2},
				{point: NewPoint(4, 4), pointType: PointTypeNormal},
				{point: NewPoint(2, 4), pointType: PointTypeAddedIntersection, entryExit: intersectionTypeEntry, otherPolyIndex: 0},
				{point: NewPoint(0, 4), pointType: PointTypeNormal},
			},
			expectedPoly2: []polyPoint[int]{
				{point: NewPoint(2, 4), pointType: PointTypeAddedIntersection, entryExit: intersectionTypeEntry, otherPolyIndex: 4},
				{point: NewPoint(2, 2), pointType: PointTypeNormal},
				{point: NewPoint(4, 2), pointType: PointTypeAddedIntersection, entryExit: intersectionTypeExit, otherPolyIndex: 2},
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
			polygonType:   PTSolid,
			expectedError: false,
			validateFunc: func(t *testing.T, polygon *Polygon[int]) {
				assert.Equal(t, PTSolid, polygon.polygonType, "Expected polygon type to be PTSolid")
				assert.Len(t, polygon.Points(), 3, "Expected three points")
				assert.NotNil(t, polygon.hull, "Expected convex hull to be created")
			},
		},
		"Invalid Fewer Than Three Points": {
			points: []Point[int]{
				NewPoint(0, 0),
				NewPoint(1, 1),
			},
			polygonType:   PTSolid,
			expectedError: true,
		},
		"Invalid Zero Area": {
			points: []Point[int]{
				NewPoint(0, 0),
				NewPoint(1, 1),
				NewPoint(2, 2),
			},
			polygonType:   PTSolid,
			expectedError: true,
		},
		"Orientation Correction for Solid Polygon": {
			points: []Point[int]{
				NewPoint(0, 0),
				NewPoint(1, 2),
				NewPoint(2, 0),
			},
			polygonType:   PTSolid,
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
			polygonType:   PTHole,
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
			polygon, err := NewPolygon(tc.points, PTSolid)
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
		"Point PCROutside Polygon": {
			polygon: func() Polygon[int] {
				poly, _ := NewPolygon([]Point[int]{
					NewPoint(0, 0),
					NewPoint(4, 0),
					NewPoint(4, 4),
					NewPoint(0, 4),
				}, PTSolid)
				return *poly
			}(),
			point:       NewPoint(5, 5),
			expectedRel: PPRPointOutside,
		},
		"Point PCRInside Polygon": {
			polygon: func() Polygon[int] {
				poly, _ := NewPolygon([]Point[int]{
					NewPoint(0, 0),
					NewPoint(4, 0),
					NewPoint(4, 4),
					NewPoint(0, 4),
				}, PTSolid)
				return *poly
			}(),
			point:       NewPoint(2, 2),
			expectedRel: PPRPointInside,
		},
		"Point On Vertex": {
			polygon: func() Polygon[int] {
				poly, _ := NewPolygon([]Point[int]{
					NewPoint(0, 0),
					NewPoint(4, 0),
					NewPoint(4, 4),
					NewPoint(0, 4),
				}, PTSolid)
				return *poly
			}(),
			point:       NewPoint(4, 0),
			expectedRel: PPRPointOnVertex,
		},
		"Point On Edge": {
			polygon: func() Polygon[int] {
				poly, _ := NewPolygon([]Point[int]{
					NewPoint(0, 0),
					NewPoint(4, 0),
					NewPoint(4, 4),
					NewPoint(0, 4),
				}, PTSolid)
				return *poly
			}(),
			point:       NewPoint(2, 0),
			expectedRel: PPRPointOnEdge,
		},
		"Point in Concave Indent": {
			polygon: func() Polygon[int] {
				poly, _ := NewPolygon([]Point[int]{
					NewPoint(0, 0),
					NewPoint(4, 0),
					NewPoint(2, 2), // Creates the indent in the polygon
					NewPoint(4, 4),
					NewPoint(0, 4),
				}, PTSolid)
				return *poly
			}(),
			point:       NewPoint(3, 2), // Point in the indent but outside the polygon
			expectedRel: PPRPointOutside,
		},
		"Point in Concave Polygon": {
			polygon: func() Polygon[int] {
				poly, _ := NewPolygon([]Point[int]{
					NewPoint(0, 0),
					NewPoint(4, 0),
					NewPoint(2, 2), // Creates the indent in the polygon
					NewPoint(4, 4),
					NewPoint(0, 4),
				}, PTSolid)
				return *poly
			}(),
			point:       NewPoint(2, 3),
			expectedRel: PPRPointInside,
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
		"Point PCRInside Concave Polygon PCROutside Hole": {
			polygon: func() Polygon[int] {
				mainPoly, _ := NewPolygon([]Point[int]{
					NewPoint(0, 0),
					NewPoint(12, 0),
					NewPoint(12, 12),
					NewPoint(0, 12),
				}, PTSolid)

				hole, _ := NewPolygon([]Point[int]{
					NewPoint(4, 4),
					NewPoint(8, 4),
					NewPoint(8, 8),
					NewPoint(4, 8),
				}, PTHole)

				mainPoly.children = []*Polygon[int]{hole}
				return *mainPoly
			}(),
			point:       NewPoint(2, 2),
			expectedRel: PPRPointInside,
		},
		"Point PCRInside Hole": {
			polygon: func() Polygon[int] {
				mainPoly, _ := NewPolygon([]Point[int]{
					NewPoint(0, 0),
					NewPoint(12, 0),
					NewPoint(12, 12),
					NewPoint(0, 12),
				}, PTSolid)

				hole, _ := NewPolygon([]Point[int]{
					NewPoint(4, 4),
					NewPoint(8, 4),
					NewPoint(8, 8),
					NewPoint(4, 8),
				}, PTHole)

				mainPoly.children = []*Polygon[int]{hole}
				return *mainPoly
			}(),
			point:       NewPoint(6, 6),
			expectedRel: PPRPointInHole,
		},
		"Point PCRInside Island": {
			polygon: func() Polygon[int] {
				mainPoly, _ := NewPolygon([]Point[int]{
					NewPoint(0, 0),
					NewPoint(12, 0),
					NewPoint(12, 12),
					NewPoint(0, 12),
				}, PTSolid)

				hole, _ := NewPolygon([]Point[int]{
					NewPoint(4, 4),
					NewPoint(8, 4),
					NewPoint(8, 8),
					NewPoint(4, 8),
				}, PTHole)

				island, _ := NewPolygon([]Point[int]{
					NewPoint(5, 5),
					NewPoint(7, 5),
					NewPoint(7, 7),
					NewPoint(5, 7),
				}, PTSolid)

				hole.children = []*Polygon[int]{island}
				mainPoly.children = []*Polygon[int]{hole}
				return *mainPoly
			}(),
			point:       NewPoint(6, 6),
			expectedRel: PPRPointInsideIsland,
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
				}, PTSolid)

				hole1, _ := NewPolygon([]Point[int]{
					NewPoint(5, 5),
					NewPoint(15, 5),
					NewPoint(15, 15),
					NewPoint(5, 15),
				}, PTHole)

				island1, _ := NewPolygon([]Point[int]{
					NewPoint(6, 6),
					NewPoint(14, 6),
					NewPoint(14, 14),
					NewPoint(6, 14),
				}, PTSolid)

				hole2, _ := NewPolygon([]Point[int]{
					NewPoint(8, 8),
					NewPoint(12, 8),
					NewPoint(12, 12),
					NewPoint(8, 12),
				}, PTHole)

				island2, _ := NewPolygon([]Point[int]{
					NewPoint(9, 9),
					NewPoint(11, 9),
					NewPoint(11, 11),
					NewPoint(9, 11),
				}, PTSolid)

				// Nest the polygons
				hole1.children = []*Polygon[int]{island1}
				island1.children = []*Polygon[int]{hole2}
				hole2.children = []*Polygon[int]{island2}
				mainPoly.children = []*Polygon[int]{hole1}
				return *mainPoly
			}(),
			point:       NewPoint(10, 10),
			expectedRel: PPRPointInsideIsland,
		},
		"Point Deeply Nested in Hole": {
			polygon: func() Polygon[int] {
				mainPoly, _ := NewPolygon([]Point[int]{
					NewPoint(0, 0),
					NewPoint(20, 0),
					NewPoint(20, 20),
					NewPoint(0, 20),
				}, PTSolid)

				hole1, _ := NewPolygon([]Point[int]{
					NewPoint(5, 5),
					NewPoint(15, 5),
					NewPoint(15, 15),
					NewPoint(5, 15),
				}, PTHole)

				island1, _ := NewPolygon([]Point[int]{
					NewPoint(6, 6),
					NewPoint(14, 6),
					NewPoint(14, 14),
					NewPoint(6, 14),
				}, PTSolid)

				hole2, _ := NewPolygon([]Point[int]{
					NewPoint(8, 8),
					NewPoint(12, 8),
					NewPoint(12, 12),
					NewPoint(8, 12),
				}, PTHole)

				// Nest the polygons
				hole1.children = []*Polygon[int]{island1}
				island1.children = []*Polygon[int]{hole2}
				mainPoly.children = []*Polygon[int]{hole1}
				return *mainPoly
			}(),
			point:       NewPoint(10, 10),
			expectedRel: PPRPointInHole,
		},
		"Point in Outer Polygon, Not in Deepest Children": {
			polygon: func() Polygon[int] {
				mainPoly, _ := NewPolygon([]Point[int]{
					NewPoint(0, 0),
					NewPoint(20, 0),
					NewPoint(20, 20),
					NewPoint(0, 20),
				}, PTSolid)

				hole1, _ := NewPolygon([]Point[int]{
					NewPoint(5, 5),
					NewPoint(15, 5),
					NewPoint(15, 15),
					NewPoint(5, 15),
				}, PTHole)

				island1, _ := NewPolygon([]Point[int]{
					NewPoint(6, 6),
					NewPoint(14, 6),
					NewPoint(14, 14),
					NewPoint(6, 14),
				}, PTSolid)

				hole2, _ := NewPolygon([]Point[int]{
					NewPoint(8, 8),
					NewPoint(12, 8),
					NewPoint(12, 12),
					NewPoint(8, 12),
				}, PTHole)

				// Nest the polygons
				hole1.children = []*Polygon[int]{island1}
				island1.children = []*Polygon[int]{hole2}
				mainPoly.children = []*Polygon[int]{hole1}
				return *mainPoly
			}(),
			point:       NewPoint(3, 3),
			expectedRel: PPRPointInside,
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
		"Point PCRInside Polygon": {
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
		"Point PCROutside Polygon": {
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
					{point: NewPoint(4, 2), pointType: PointTypeAddedIntersection, entryExit: intersectionTypeEntry, visited: false, otherPolyIndex: 2},
					{point: NewPoint(4, 4), pointType: PointTypeNormal, entryExit: intersectionTypeNotSet, visited: false, otherPolyIndex: 0},
					{point: NewPoint(2, 4), pointType: PointTypeAddedIntersection, entryExit: intersectionTypeEntry, visited: false, otherPolyIndex: 4},
					{point: NewPoint(2, 2), pointType: PointTypeNormal, entryExit: intersectionTypeNotSet, visited: false, otherPolyIndex: 0},
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
					{point: NewPoint(2, 2), pointType: PointTypeNormal, entryExit: intersectionTypeNotSet, visited: false, otherPolyIndex: 0},
					{point: NewPoint(4, 2), pointType: PointTypeNormal, entryExit: intersectionTypeNotSet, visited: false, otherPolyIndex: 0},
					{point: NewPoint(4, 4), pointType: PointTypeNormal, entryExit: intersectionTypeNotSet, visited: false, otherPolyIndex: 0},
					{point: NewPoint(2, 4), pointType: PointTypeNormal, entryExit: intersectionTypeNotSet, visited: false, otherPolyIndex: 0},
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
					{point: NewPoint(2, 4), pointType: PointTypeAddedIntersection, entryExit: intersectionTypeEntry, visited: false, otherPolyIndex: 0},
					{point: NewPoint(0, 4), pointType: PointTypeNormal, entryExit: intersectionTypeNotSet, visited: false, otherPolyIndex: 0},
					{point: NewPoint(0, 0), pointType: PointTypeNormal, entryExit: intersectionTypeNotSet, visited: false, otherPolyIndex: 0},
					{point: NewPoint(4, 0), pointType: PointTypeNormal, entryExit: intersectionTypeNotSet, visited: false, otherPolyIndex: 0},
					{point: NewPoint(4, 2), pointType: PointTypeAddedIntersection, entryExit: intersectionTypeEntry, visited: false, otherPolyIndex: 2},
					{point: NewPoint(6, 2), pointType: PointTypeNormal, entryExit: intersectionTypeNotSet, visited: false, otherPolyIndex: 0},
					{point: NewPoint(6, 6), pointType: PointTypeNormal, entryExit: intersectionTypeNotSet, visited: false, otherPolyIndex: 0},
					{point: NewPoint(2, 6), pointType: PointTypeNormal, entryExit: intersectionTypeNotSet, visited: false, otherPolyIndex: 0},
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
					{point: NewPoint(0, 0), pointType: PointTypeNormal, entryExit: intersectionTypeNotSet, visited: false, otherPolyIndex: 0},
					{point: NewPoint(4, 0), pointType: PointTypeNormal, entryExit: intersectionTypeNotSet, visited: false, otherPolyIndex: 0},
					{point: NewPoint(4, 4), pointType: PointTypeNormal, entryExit: intersectionTypeNotSet, visited: false, otherPolyIndex: 0},
					{point: NewPoint(0, 4), pointType: PointTypeNormal, entryExit: intersectionTypeNotSet, visited: false, otherPolyIndex: 0},
				},
				{
					{point: NewPoint(6, 6), pointType: PointTypeNormal, entryExit: intersectionTypeNotSet, visited: false, otherPolyIndex: 0},
					{point: NewPoint(10, 6), pointType: PointTypeNormal, entryExit: intersectionTypeNotSet, visited: false, otherPolyIndex: 0},
					{point: NewPoint(10, 10), pointType: PointTypeNormal, entryExit: intersectionTypeNotSet, visited: false, otherPolyIndex: 0},
					{point: NewPoint(6, 10), pointType: PointTypeNormal, entryExit: intersectionTypeNotSet, visited: false, otherPolyIndex: 0},
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
