package geom2d

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewBetterPolygon(t *testing.T) {
	tests := map[string]struct {
		points   []Point[int]
		t        PolygonType
		expected PolyTree[int]
	}{
		"solid": {
			points: []Point[int]{
				{x: 0, y: 0},
				{x: 6, y: 0},
				{x: 6, y: 6},
				{x: 0, y: 6},
			},
			t: PTSolid,
			expected: PolyTree[int]{
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
				traversalDirection: polyTraversalDirectionCounterClockwise,
				polygonType:        PTSolid,
				children:           nil,
				parent:             nil,
				hull: simpleConvexPolygon[int]{
					Points: []Point[int]{
						{x: 0, y: 0},
						{x: 6, y: 0},
						{x: 6, y: 6},
						{x: 0, y: 6},
					},
				},
				maxX: 13,
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
			expected: PolyTree[int]{
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
				traversalDirection: polyTraversalDirectionClockwise,
				polygonType:        PTHole,
				children:           nil,
				parent:             nil,
				hull: simpleConvexPolygon[int]{
					Points: []Point[int]{
						{x: 0, y: 0},
						{x: 6, y: 0},
						{x: 6, y: 6},
						{x: 0, y: 6},
					},
				},
				maxX: 13,
			},
		},
	}
	for _, tc := range tests {
		p, err := NewPolyTree(tc.points, tc.t)
		require.NoError(t, err, "expected no error")
		assert.Equal(t, &tc.expected, p)
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

func TestNewBetterPolygon_findIntersectionsBetweenBetterPolys(t *testing.T) {
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

	// mark points for Union
	polyTree1.markEntryExitPoints(polyTree2, BooleanUnion)

	// traverse for union
	_ = polyTree1.traverse(polyTree2, BooleanUnion)

	fmt.Println("not yet finished!")

}
