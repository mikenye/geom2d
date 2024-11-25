package geom2d

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNestPointsToPolyTrees(t *testing.T) {
	tests := map[string]struct {
		contours [][]Point[int]
		expected func() *PolyTree[int]
		wantErr  bool
	}{
		"single polygon": {
			contours: [][]Point[int]{
				{{0, 0}, {10, 0}, {10, 10}, {0, 10}},
			},
			expected: func() *PolyTree[int] {
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
				}
			},
			wantErr: false,
		},
		"nested polygons": {
			contours: [][]Point[int]{
				{{0, 0}, {20, 0}, {20, 20}, {0, 20}}, // Outer solid
				{{5, 5}, {15, 5}, {15, 15}, {5, 15}}, // Inner hole
				{{7, 7}, {13, 7}, {13, 13}, {7, 13}}, // Island inside hole
			},
			expected: func() *PolyTree[int] {
				root := &PolyTree[int]{
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
							point:                         Point[int]{40, 0},
							pointType:                     pointTypeOriginal,
							entryExit:                     intersectionTypeNotSet,
							visited:                       false,
							intersectionPartner:           nil,
							intersectionPartnerPointIndex: -1,
						},
						polyTreePoint[int]{ // 2
							point:                         Point[int]{40, 40},
							pointType:                     pointTypeOriginal,
							entryExit:                     intersectionTypeNotSet,
							visited:                       false,
							intersectionPartner:           nil,
							intersectionPartnerPointIndex: -1,
						},
						polyTreePoint[int]{ // 3
							point:                         Point[int]{0, 40},
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
						Points: []Point[int]{{0, 0}, {20, 0}, {20, 20}, {0, 20}},
					},
					maxX: 41,
				}
				hole := &PolyTree[int]{
					contour: contour[int]{
						polyTreePoint[int]{ // 0
							point:                         Point[int]{10, 10},
							pointType:                     pointTypeOriginal,
							entryExit:                     intersectionTypeNotSet,
							visited:                       false,
							intersectionPartner:           nil,
							intersectionPartnerPointIndex: -1,
						},
						polyTreePoint[int]{ // 1
							point:                         Point[int]{30, 10},
							pointType:                     pointTypeOriginal,
							entryExit:                     intersectionTypeNotSet,
							visited:                       false,
							intersectionPartner:           nil,
							intersectionPartnerPointIndex: -1,
						},
						polyTreePoint[int]{ // 2
							point:                         Point[int]{30, 30},
							pointType:                     pointTypeOriginal,
							entryExit:                     intersectionTypeNotSet,
							visited:                       false,
							intersectionPartner:           nil,
							intersectionPartnerPointIndex: -1,
						},
						polyTreePoint[int]{ // 3
							point:                         Point[int]{10, 30},
							pointType:                     pointTypeOriginal,
							entryExit:                     intersectionTypeNotSet,
							visited:                       false,
							intersectionPartner:           nil,
							intersectionPartnerPointIndex: -1,
						},
					},
					polygonType: PTHole,
					siblings:    nil,
					children:    nil,
					parent:      nil,
					hull: simpleConvexPolygon[int]{
						Points: []Point[int]{{5, 5}, {15, 5}, {15, 15}, {5, 15}},
					},
					maxX: 31,
				}
				island := &PolyTree[int]{
					contour: contour[int]{
						polyTreePoint[int]{ // 0
							point:                         Point[int]{14, 14},
							pointType:                     pointTypeOriginal,
							entryExit:                     intersectionTypeNotSet,
							visited:                       false,
							intersectionPartner:           nil,
							intersectionPartnerPointIndex: -1,
						},
						polyTreePoint[int]{ // 1
							point:                         Point[int]{26, 14},
							pointType:                     pointTypeOriginal,
							entryExit:                     intersectionTypeNotSet,
							visited:                       false,
							intersectionPartner:           nil,
							intersectionPartnerPointIndex: -1,
						},
						polyTreePoint[int]{ // 2
							point:                         Point[int]{26, 26},
							pointType:                     pointTypeOriginal,
							entryExit:                     intersectionTypeNotSet,
							visited:                       false,
							intersectionPartner:           nil,
							intersectionPartnerPointIndex: -1,
						},
						polyTreePoint[int]{
							point:                         Point[int]{14, 26},
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
						Points: []Point[int]{{7, 7}, {13, 7}, {13, 13}, {7, 13}},
					},
					maxX: 27,
				}
				hole.addChild(island)
				root.addChild(hole)
				return root
			},
			wantErr: false,
		},
		"siblings": {
			contours: [][]Point[int]{
				{{0, 0}, {10, 0}, {10, 10}, {0, 10}},     // Polygon 1
				{{20, 20}, {30, 20}, {30, 30}, {20, 30}}, // Polygon 2
			},
			expected: func() *PolyTree[int] {
				poly1 := &PolyTree[int]{
					contour: contour[int]{
						polyTreePoint[int]{ // 0
							point:                         Point[int]{40, 40},
							pointType:                     pointTypeOriginal,
							entryExit:                     intersectionTypeNotSet,
							visited:                       false,
							intersectionPartner:           nil,
							intersectionPartnerPointIndex: -1,
						},
						polyTreePoint[int]{ // 1
							point:                         Point[int]{60, 40},
							pointType:                     pointTypeOriginal,
							entryExit:                     intersectionTypeNotSet,
							visited:                       false,
							intersectionPartner:           nil,
							intersectionPartnerPointIndex: -1,
						},
						polyTreePoint[int]{ // 2
							point:                         Point[int]{60, 60},
							pointType:                     pointTypeOriginal,
							entryExit:                     intersectionTypeNotSet,
							visited:                       false,
							intersectionPartner:           nil,
							intersectionPartnerPointIndex: -1,
						},
						polyTreePoint[int]{ // 3
							point:                         Point[int]{40, 60},
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
						Points: []Point[int]{{20, 20}, {30, 20}, {30, 30}, {20, 30}},
					},
					maxX: 61,
				}
				poly2 := &PolyTree[int]{
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
				}
				poly1.addSibling(poly2)
				return poly1
			},
			wantErr: false,
		},
		"no input polygons": {
			contours: [][]Point[int]{},
			expected: nil,
			wantErr:  true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := nestPointsToPolyTrees(tc.contours)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expected(), got)
		})
	}
}

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

func TestNewBetterPolygon_booleanOperationTraversal_Union(t *testing.T) {
	polyTree1HolePoints := []Point[int]{
		{5, 5},
		{15, 5},
		{15, 15},
		{5, 15},
	}
	polyTree1Hole, err := NewPolyTree(polyTree1HolePoints, PTHole)
	require.NoError(t, err, "expected no error when creating polyTree1Hole")

	expectedPolyTree1Hole := &PolyTree[int]{
		contour: contour[int]{
			polyTreePoint[int]{ // 0
				point:                         Point[int]{10, 30},
				pointType:                     pointTypeOriginal,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 1
				point:                         Point[int]{30, 30},
				pointType:                     pointTypeOriginal,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 2
				point:                         Point[int]{30, 10},
				pointType:                     pointTypeOriginal,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 3
				point:                         Point[int]{10, 10},
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
				{5, 5},
				{15, 5},
				{15, 15},
				{5, 15},
			},
		},
		maxX: 31,
	}
	require.Equal(t, expectedPolyTree1Hole, polyTree1Hole, "unexpected output of NewPolyTree for polyTree1Hole")

	polyTree1Points := []Point[int]{
		{0, 0},
		{20, 0},
		{20, 20},
		{0, 20},
	}
	polyTree1, err := NewPolyTree(polyTree1Points, PTSolid, WithChildren(polyTree1Hole))
	require.NoError(t, err, "expected no error when creating polyTree1")

	expectedPolyTree1 := &PolyTree[int]{
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
				point:                         Point[int]{40, 0},
				pointType:                     pointTypeOriginal,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 2
				point:                         Point[int]{40, 40},
				pointType:                     pointTypeOriginal,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 3
				point:                         Point[int]{0, 40},
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
				{0, 0}, {20, 0}, {20, 20}, {0, 20},
			},
		},
		maxX: 41,
	}
	expectedPolyTree1Hole.parent = expectedPolyTree1
	expectedPolyTree1.children = append(expectedPolyTree1.children, expectedPolyTree1Hole)
	require.Equal(t, expectedPolyTree1, polyTree1, "unexpected output of NewPolyTree for polyTree1")

	polyTree2HolePoints := []Point[int]{
		{12, 12},
		{22, 12},
		{22, 22},
		{12, 22},
	}
	polyTree2Hole, err := NewPolyTree(polyTree2HolePoints, PTHole)
	require.NoError(t, err, "expected no error when creating polyTree2Hole")

	expectedPolyTree2Hole := &PolyTree[int]{
		contour: contour[int]{
			polyTreePoint[int]{ // 0
				point:                         Point[int]{24, 44},
				pointType:                     pointTypeOriginal,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 1
				point:                         Point[int]{44, 44},
				pointType:                     pointTypeOriginal,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 2
				point:                         Point[int]{44, 24},
				pointType:                     pointTypeOriginal,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 3
				point:                         Point[int]{24, 24},
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
				{12, 12}, {22, 12}, {22, 22}, {12, 22},
			},
		},
		maxX: 45,
	}
	require.Equal(t, expectedPolyTree2Hole, polyTree2Hole, "unexpected output of NewPolyTree for polyTree2Hole")

	poly2Points := []Point[int]{
		{7, 7},
		{27, 7},
		{27, 27},
		{7, 27},
	}
	polyTree2, err := NewPolyTree(poly2Points, PTSolid, WithChildren(polyTree2Hole))
	require.NoError(t, err, "expected no error when creating polyTree2")

	expectedPolyTree2 := &PolyTree[int]{
		contour: contour[int]{
			polyTreePoint[int]{ // 0
				point:                         Point[int]{14, 14},
				pointType:                     pointTypeOriginal,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 1
				point:                         Point[int]{54, 14},
				pointType:                     pointTypeOriginal,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 2
				point:                         Point[int]{54, 54},
				pointType:                     pointTypeOriginal,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 3
				point:                         Point[int]{14, 54},
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
				{7, 7}, {27, 7}, {27, 27}, {7, 27},
			},
		},
		maxX: 55,
	}
	expectedPolyTree2Hole.parent = expectedPolyTree2
	expectedPolyTree2.children = append(expectedPolyTree2.children, expectedPolyTree2Hole)
	require.Equal(t, expectedPolyTree2, polyTree2, "unexpected output of NewPolyTree for polyTree2")

	// find intersection points between all polys
	polyTree1.findIntersections(polyTree2)

	expectedPolyTree1Hole = &PolyTree[int]{
		contour: contour[int]{
			polyTreePoint[int]{ // 0
				point:                         Point[int]{10, 30},
				pointType:                     pointTypeOriginal,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 1
				point:                         Point[int]{14, 30},
				pointType:                     pointTypeAddedIntersection,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 2
				point:                         Point[int]{24, 30},
				pointType:                     pointTypeAddedIntersection,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 3
				point:                         Point[int]{30, 30},
				pointType:                     pointTypeOriginal,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 4
				point:                         Point[int]{30, 24},
				pointType:                     pointTypeAddedIntersection,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 5
				point:                         Point[int]{30, 14},
				pointType:                     pointTypeAddedIntersection,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 6
				point:                         Point[int]{30, 10},
				pointType:                     pointTypeOriginal,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 7
				point:                         Point[int]{10, 10},
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
				{5, 5},
				{15, 5},
				{15, 15},
				{5, 15},
			},
		},
		maxX: 31,
	}
	expectedPolyTree1 = &PolyTree[int]{
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
				point:                         Point[int]{40, 0},
				pointType:                     pointTypeOriginal,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 2
				point:                         Point[int]{40, 14},
				pointType:                     pointTypeAddedIntersection,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 3
				point:                         Point[int]{40, 24},
				pointType:                     pointTypeAddedIntersection,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 4
				point:                         Point[int]{40, 40},
				pointType:                     pointTypeOriginal,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 5
				point:                         Point[int]{24, 40},
				pointType:                     pointTypeAddedIntersection,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 6
				point:                         Point[int]{14, 40},
				pointType:                     pointTypeAddedIntersection,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 7
				point:                         Point[int]{0, 40},
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
				{0, 0},
				{20, 0},
				{20, 20},
				{0, 20},
			},
		},
		maxX: 41,
	}
	expectedPolyTree1Hole.parent = expectedPolyTree1
	expectedPolyTree1.children = append(expectedPolyTree1.children, expectedPolyTree1Hole)
	require.Equal(t, expectedPolyTree1, polyTree1, "unexpected output of findIntersections for polyTree1")

	expectedPolyTree2Hole = &PolyTree[int]{
		contour: contour[int]{
			polyTreePoint[int]{ // 0
				point:                         Point[int]{24, 30},
				pointType:                     pointTypeAddedIntersection,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 1
				point:                         Point[int]{24, 40},
				pointType:                     pointTypeAddedIntersection,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 2
				point:                         Point[int]{24, 44},
				pointType:                     pointTypeOriginal,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 3
				point:                         Point[int]{44, 44},
				pointType:                     pointTypeOriginal,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 4
				point:                         Point[int]{44, 24},
				pointType:                     pointTypeOriginal,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 5
				point:                         Point[int]{40, 24},
				pointType:                     pointTypeAddedIntersection,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 6
				point:                         Point[int]{30, 24},
				pointType:                     pointTypeAddedIntersection,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 7
				point:                         Point[int]{24, 24},
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
				{12, 12},
				{22, 12},
				{22, 22},
				{12, 22},
			},
		},
		maxX: 45,
	}
	expectedPolyTree2 = &PolyTree[int]{
		contour: contour[int]{
			polyTreePoint[int]{ // 0
				point:                         Point[int]{14, 40},
				pointType:                     pointTypeAddedIntersection,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 1
				point:                         Point[int]{14, 30},
				pointType:                     pointTypeAddedIntersection,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 2
				point:                         Point[int]{14, 14},
				pointType:                     pointTypeOriginal,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 3
				point:                         Point[int]{30, 14},
				pointType:                     pointTypeAddedIntersection,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 4
				point:                         Point[int]{40, 14},
				pointType:                     pointTypeAddedIntersection,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 5
				point:                         Point[int]{54, 14},
				pointType:                     pointTypeOriginal,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 6
				point:                         Point[int]{54, 54},
				pointType:                     pointTypeOriginal,
				entryExit:                     intersectionTypeNotSet,
				visited:                       false,
				intersectionPartner:           nil,
				intersectionPartnerPointIndex: -1,
			},
			polyTreePoint[int]{ // 7
				point:                         Point[int]{14, 54},
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
				{7, 7},
				{27, 7},
				{27, 27},
				{7, 27},
			},
		},
		maxX: 55,
	}
	expectedPolyTree2Hole.parent = expectedPolyTree2
	expectedPolyTree2.children = append(expectedPolyTree2.children, expectedPolyTree2Hole)
	require.Equal(t, expectedPolyTree2, polyTree2, "unexpected output of findIntersections for polyTree2")

	// mark points for Union
	polyTree1.markEntryExitPoints(polyTree2, BooleanUnion)

	expectedPolyTree1.contour[2].entryExit = intersectionTypeExit
	expectedPolyTree1.contour[2].intersectionPartner = expectedPolyTree2
	expectedPolyTree1.contour[2].intersectionPartnerPointIndex = 4

	expectedPolyTree1.contour[3].entryExit = intersectionTypeEntry
	expectedPolyTree1.contour[3].intersectionPartner = expectedPolyTree2Hole
	expectedPolyTree1.contour[3].intersectionPartnerPointIndex = 5

	expectedPolyTree1.contour[5].entryExit = intersectionTypeExit
	expectedPolyTree1.contour[5].intersectionPartner = expectedPolyTree2Hole
	expectedPolyTree1.contour[5].intersectionPartnerPointIndex = 1

	expectedPolyTree1.contour[6].entryExit = intersectionTypeEntry
	expectedPolyTree1.contour[6].intersectionPartner = expectedPolyTree2
	expectedPolyTree1.contour[6].intersectionPartnerPointIndex = 0

	expectedPolyTree1Hole.contour[1].entryExit = intersectionTypeExit
	expectedPolyTree1Hole.contour[1].intersectionPartner = expectedPolyTree2
	expectedPolyTree1Hole.contour[1].intersectionPartnerPointIndex = 1

	expectedPolyTree1Hole.contour[2].entryExit = intersectionTypeEntry
	expectedPolyTree1Hole.contour[2].intersectionPartner = expectedPolyTree2Hole
	expectedPolyTree1Hole.contour[2].intersectionPartnerPointIndex = 0

	expectedPolyTree1Hole.contour[4].entryExit = intersectionTypeExit
	expectedPolyTree1Hole.contour[4].intersectionPartner = expectedPolyTree2Hole
	expectedPolyTree1Hole.contour[4].intersectionPartnerPointIndex = 6

	expectedPolyTree1Hole.contour[5].entryExit = intersectionTypeEntry
	expectedPolyTree1Hole.contour[5].intersectionPartner = expectedPolyTree2
	expectedPolyTree1Hole.contour[5].intersectionPartnerPointIndex = 3

	expectedPolyTree2.contour[0].entryExit = intersectionTypeExit
	expectedPolyTree2.contour[0].intersectionPartner = expectedPolyTree1
	expectedPolyTree2.contour[0].intersectionPartnerPointIndex = 6

	expectedPolyTree2.contour[1].entryExit = intersectionTypeEntry
	expectedPolyTree2.contour[1].intersectionPartner = expectedPolyTree1Hole
	expectedPolyTree2.contour[1].intersectionPartnerPointIndex = 1

	expectedPolyTree2.contour[3].entryExit = intersectionTypeExit
	expectedPolyTree2.contour[3].intersectionPartner = expectedPolyTree1Hole
	expectedPolyTree2.contour[3].intersectionPartnerPointIndex = 5

	expectedPolyTree2.contour[4].entryExit = intersectionTypeEntry
	expectedPolyTree2.contour[4].intersectionPartner = expectedPolyTree1
	expectedPolyTree2.contour[4].intersectionPartnerPointIndex = 2

	expectedPolyTree2Hole.contour[0].entryExit = intersectionTypeExit
	expectedPolyTree2Hole.contour[0].intersectionPartner = expectedPolyTree1Hole
	expectedPolyTree2Hole.contour[0].intersectionPartnerPointIndex = 2

	expectedPolyTree2Hole.contour[1].entryExit = intersectionTypeEntry
	expectedPolyTree2Hole.contour[1].intersectionPartner = expectedPolyTree1
	expectedPolyTree2Hole.contour[1].intersectionPartnerPointIndex = 5

	expectedPolyTree2Hole.contour[5].entryExit = intersectionTypeExit
	expectedPolyTree2Hole.contour[5].intersectionPartner = expectedPolyTree1
	expectedPolyTree2Hole.contour[5].intersectionPartnerPointIndex = 3

	expectedPolyTree2Hole.contour[6].entryExit = intersectionTypeEntry
	expectedPolyTree2Hole.contour[6].intersectionPartner = expectedPolyTree1Hole
	expectedPolyTree2Hole.contour[6].intersectionPartnerPointIndex = 4

	require.Equal(t, expectedPolyTree1, polyTree1, "unexpected output of markEntryExitPoints for polyTree1")

	require.Equal(t, expectedPolyTree2, polyTree2, "unexpected output of markEntryExitPoints for polyTree2")

	// traverse for union
	expectedPointsUnion := [][]Point[int]{
		{
			{40 / 2, 24 / 2},
			{40 / 2, 40 / 2},
			{24 / 2, 40 / 2},
			{24 / 2, 44 / 2},
			{44 / 2, 44 / 2},
			{44 / 2, 24 / 2},
		},
		{
			{14 / 2, 40 / 2},
			{0 / 2, 40 / 2},
			{0 / 2, 0 / 2},
			{40 / 2, 0 / 2},
			{40 / 2, 14 / 2},
			{54 / 2, 14 / 2},
			{54 / 2, 54 / 2},
			{14 / 2, 54 / 2},
		},
		{
			{24 / 2, 30 / 2},
			{30 / 2, 30 / 2},
			{30 / 2, 24 / 2},
			{24 / 2, 24 / 2},
		},
		{
			{30 / 2, 14 / 2},
			{30 / 2, 10 / 2},
			{10 / 2, 10 / 2},
			{10 / 2, 30 / 2},
			{14 / 2, 30 / 2},
			{14 / 2, 14 / 2},
		},
	}
	resultingPointsUnion := polyTree1.booleanOperationTraversal(polyTree2, BooleanUnion)
	assert.Equal(t, expectedPointsUnion, resultingPointsUnion)
}

func TestNewBetterPolygon_booleanOperationTraversal_Intersection(t *testing.T) {
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

func TestNewBetterPolygon_booleanOperationTraversal_Subtraction(t *testing.T) {
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
