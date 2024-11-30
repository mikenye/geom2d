package geom2d

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPolyTree_OrderConsistency(t *testing.T) {
	// Create root and children
	root, err := NewPolyTree([]Point[int]{{10, 10}, {20, 10}, {20, 20}, {10, 20}}, PTSolid)
	require.NoError(t, err, "unexpected error returned when creating root")
	child1, err := NewPolyTree([]Point[int]{{30, 30}, {40, 30}, {40, 40}, {30, 40}}, PTHole)
	require.NoError(t, err, "unexpected error returned when creating child1")
	child2, err := NewPolyTree([]Point[int]{{5, 5}, {15, 5}, {15, 15}, {5, 15}}, PTHole)
	require.NoError(t, err, "unexpected error returned when creating child2")

	// Add children to root
	err = root.addChild(child1)
	require.NoError(t, err, "unexpected error returned when adding child1 as a child of root")
	err = root.addChild(child2)
	require.NoError(t, err, "unexpected error returned when adding child2 as a child of root")

	// Verify children order
	expectedChildOrder := []*PolyTree[int]{child2, child1}
	assert.Equal(t, expectedChildOrder, root.children, "Children should be ordered by lowest, leftmost point")

	// Create siblings
	sibling1, err := NewPolyTree([]Point[int]{{50, 50}, {60, 50}, {60, 60}, {50, 60}}, PTSolid)
	require.NoError(t, err, "unexpected error returned when creating sibling1")
	sibling2, err := NewPolyTree([]Point[int]{{25, 25}, {35, 25}, {35, 35}, {25, 35}}, PTSolid)
	require.NoError(t, err, "unexpected error returned when creating sibling2")

	// Add siblings to root
	err = root.addSibling(sibling1)
	require.NoError(t, err, "unexpected error returned when adding sibling1 as a sibling of root")
	err = root.addSibling(sibling2)
	require.NoError(t, err, "unexpected error returned when adding sibling2 as a sibling of root")

	// Verify sibling order
	expectedSiblingOrder := []*PolyTree[int]{sibling2, sibling1}
	assert.Equal(t, expectedSiblingOrder, root.siblings, "Siblings should be ordered by lowest, leftmost point")
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
				if err := root.addSibling(sibling); err != nil {
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
				if err := root.addChild(hole); err != nil {
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

func TestPolyTree_AddSibling(t *testing.T) {
	t.Run("Adding a nil sibling", func(t *testing.T) {
		poly1, _ := NewPolyTree([]Point[int]{{0, 0}, {10, 0}, {10, 10}, {0, 10}}, PTSolid)
		err := poly1.addSibling(nil)
		require.Error(t, err, "expected error when adding a nil sibling, but got none")
	})

	t.Run("Adding a Sibling with Matching polygonType", func(t *testing.T) {
		poly1, err := NewPolyTree([]Point[int]{{0, 0}, {10, 0}, {10, 10}, {0, 10}}, PTSolid)
		require.NoError(t, err, "error creating poly1, when none was expected")
		poly2, err := NewPolyTree([]Point[int]{{20, 20}, {30, 20}, {30, 30}, {20, 30}}, PTSolid)
		require.NoError(t, err, "error creating poly2, when none was expected")

		err = poly1.addSibling(poly2)
		require.NoError(t, err, "error calling addSibling, when none was expected")
		assert.Contains(t, poly1.siblings, poly2)
		assert.Contains(t, poly2.siblings, poly1)
	})

	t.Run("Adding a Sibling with Mismatched polygonType", func(t *testing.T) {
		poly1, err := NewPolyTree([]Point[int]{{0, 0}, {10, 0}, {10, 10}, {0, 10}}, PTSolid)
		require.NoError(t, err, "error creating poly1, when none was expected")
		poly2, err := NewPolyTree([]Point[int]{{20, 20}, {30, 20}, {30, 30}, {20, 30}}, PTHole)
		require.NoError(t, err, "error creating poly2, when none was expected")

		err = poly1.addSibling(poly2)
		require.Error(t, err, "no error returned from addSibling, when one was expected")
	})

	t.Run("Adding Multiple Siblings", func(t *testing.T) {
		poly1, err := NewPolyTree([]Point[int]{{0, 0}, {10, 0}, {10, 10}, {0, 10}}, PTSolid)
		require.NoError(t, err, "error creating poly1, when none was expected")
		poly2, err := NewPolyTree([]Point[int]{{20, 20}, {30, 20}, {30, 30}, {20, 30}}, PTSolid)
		require.NoError(t, err, "error creating poly2, when none was expected")
		poly3, err := NewPolyTree([]Point[int]{{40, 40}, {50, 40}, {50, 50}, {40, 50}}, PTSolid)
		require.NoError(t, err, "error creating poly3, when none was expected")

		err = poly1.addSibling(poly2)
		require.NoError(t, err, "error returned from poly1.addSibling(poly2) when none was expected")
		err = poly1.addSibling(poly3)
		require.NoError(t, err, "error returned from poly1.addSibling(poly3) when none was expected")

		assert.Contains(t, poly1.siblings, poly2)
		assert.Contains(t, poly1.siblings, poly3)
		assert.Contains(t, poly2.siblings, poly1)
		assert.Contains(t, poly2.siblings, poly3)
		assert.Contains(t, poly3.siblings, poly1)
		assert.Contains(t, poly3.siblings, poly2)
	})
}

func TestPolyTree_AddChild(t *testing.T) {
	t.Run("Adding a nil child", func(t *testing.T) {
		parent, _ := NewPolyTree([]Point[int]{{0, 0}, {10, 0}, {10, 10}, {0, 10}}, PTSolid)
		err := parent.addChild(nil)
		require.Error(t, err, "expected error when adding a nil child, but got none")
	})

	t.Run("Adding a Child with Opposite polygonType", func(t *testing.T) {
		parent, err := NewPolyTree([]Point[int]{{0, 0}, {10, 0}, {10, 10}, {0, 10}}, PTSolid)
		require.NoError(t, err, "error creating parent polygon, when none was expected")

		child, err := NewPolyTree([]Point[int]{{2, 2}, {8, 2}, {8, 8}, {2, 8}}, PTHole)
		require.NoError(t, err, "error creating child polygon, when none was expected")

		err = parent.addChild(child)
		require.NoError(t, err, "error calling addChild, when none was expected")

		assert.Contains(t, parent.children, child)
		assert.Equal(t, parent, child.parent)
	})

	t.Run("Adding a Child with Mismatched polygonType", func(t *testing.T) {
		parent, err := NewPolyTree([]Point[int]{{0, 0}, {10, 0}, {10, 10}, {0, 10}}, PTSolid)
		require.NoError(t, err, "error creating parent polygon, when none was expected")

		child, err := NewPolyTree([]Point[int]{{2, 2}, {8, 2}, {8, 8}, {2, 8}}, PTSolid)
		require.NoError(t, err, "error creating child polygon, when none was expected")

		err = parent.addChild(child)
		require.Error(t, err, "no error returned from addChild, when one was expected")
	})

	t.Run("Adding Multiple Children", func(t *testing.T) {
		parent, err := NewPolyTree([]Point[int]{{0, 0}, {10, 0}, {10, 10}, {0, 10}}, PTSolid)
		require.NoError(t, err, "error creating parent polygon, when none was expected")

		child1, err := NewPolyTree([]Point[int]{{2, 2}, {4, 2}, {4, 4}, {2, 4}}, PTHole)
		require.NoError(t, err, "error creating first child polygon, when none was expected")

		child2, err := NewPolyTree([]Point[int]{{6, 6}, {8, 6}, {8, 8}, {6, 8}}, PTHole)
		require.NoError(t, err, "error creating second child polygon, when none was expected")

		err = parent.addChild(child1)
		require.NoError(t, err, "error calling addChild for child1, when none was expected")
		err = parent.addChild(child2)
		require.NoError(t, err, "error calling addChild for child2, when none was expected")

		assert.Contains(t, parent.children, child1)
		assert.Contains(t, parent.children, child2)
		assert.Equal(t, parent, child1.parent)
		assert.Equal(t, parent, child2.parent)
	})
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
				if err := hole.addChild(island); err != nil {
					return nil, fmt.Errorf("failed to add island: %w", err)
				}
				if err := root.addChild(hole); err != nil {
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

func TestContour_iterEdges_Empty(t *testing.T) {
	contour := contour[int]{}
	count := 0
	contour.iterEdges(func(edge LineSegment[int]) bool {
		count++
		return true
	})
	require.Equal(t, 0, count, "iterEdges should not yield edges for an empty contour")
}

func TestContour_iterEdges_TwoPoints(t *testing.T) {
	contour := contour[int]{
		{point: Point[int]{0, 0}},
		{point: Point[int]{10, 0}},
	}
	var edges []LineSegment[int]
	contour.iterEdges(func(edge LineSegment[int]) bool {
		edges = append(edges, edge)
		return true
	})
	require.Equal(t, 2, len(edges), "iterEdges should yield exactly two edges for a closed loop with two points")
	require.Equal(t, NewLineSegment(Point[int]{0, 0}, Point[int]{10, 0}), edges[0])
	require.Equal(t, NewLineSegment(Point[int]{10, 0}, Point[int]{0, 0}), edges[1])
}

func TestContour_iterEdges_Triangle(t *testing.T) {
	contour := contour[int]{
		{point: Point[int]{0, 0}},
		{point: Point[int]{10, 0}},
		{point: Point[int]{5, 10}},
	}
	var edges []LineSegment[int]
	contour.iterEdges(func(edge LineSegment[int]) bool {
		edges = append(edges, edge)
		return true
	})
	require.Equal(t, 3, len(edges), "iterEdges should yield exactly three edges for a triangle")
	require.Equal(t, NewLineSegment(Point[int]{0, 0}, Point[int]{10, 0}), edges[0])
	require.Equal(t, NewLineSegment(Point[int]{10, 0}, Point[int]{5, 10}), edges[1])
	require.Equal(t, NewLineSegment(Point[int]{5, 10}, Point[int]{0, 0}), edges[2])
}

func TestContour_iterEdges_FullPolygon(t *testing.T) {
	contour := contour[int]{
		{point: Point[int]{0, 0}},
		{point: Point[int]{10, 0}},
		{point: Point[int]{10, 10}},
		{point: Point[int]{0, 10}},
	}
	var edges []LineSegment[int]
	contour.iterEdges(func(edge LineSegment[int]) bool {
		edges = append(edges, edge)
		return true
	})
	require.Equal(t, 4, len(edges), "iterEdges should yield one edge per contour segment")
	require.Equal(t, NewLineSegment(Point[int]{0, 0}, Point[int]{10, 0}), edges[0])
	require.Equal(t, NewLineSegment(Point[int]{10, 0}, Point[int]{10, 10}), edges[1])
	require.Equal(t, NewLineSegment(Point[int]{10, 10}, Point[int]{0, 10}), edges[2])
	require.Equal(t, NewLineSegment(Point[int]{0, 10}, Point[int]{0, 0}), edges[3])
}

func TestContour_iterEdges_EarlyExit(t *testing.T) {
	contour := contour[int]{
		{point: Point[int]{0, 0}},
		{point: Point[int]{10, 0}},
		{point: Point[int]{10, 10}},
		{point: Point[int]{0, 10}},
	}
	count := 0
	contour.iterEdges(func(edge LineSegment[int]) bool {
		count++
		return count < 2 // Exit after two edges
	})
	require.Equal(t, 2, count, "iterEdges should exit early when yield returns false")
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
