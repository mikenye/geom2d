package _old

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"image"
	"testing"
)

func TestRectangle_RelationshipToCircle(t *testing.T) {
	// Define a rectangle
	rect := NewRectangle([]Point[int]{
		NewPoint[int](0, 0),
		NewPoint[int](100, 0),
		NewPoint[int](100, 100),
		NewPoint[int](0, 100),
	})

	tests := map[string]struct {
		circle      Circle[int]
		expectedRel Relationship
	}{
		"Circle inside rectangle": {
			circle:      NewCircle(NewPoint[int](50, 50), 10),
			expectedRel: RelationshipContains,
		},
		"Circle intersecting rectangle": {
			circle:      NewCircle(NewPoint[int](50, 50), 60),
			expectedRel: RelationshipIntersection,
		},
		"Circle outside rectangle": {
			circle:      NewCircle(NewPoint[int](200, 200), 20),
			expectedRel: RelationshipDisjoint,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			rel := rect.RelationshipToCircle(test.circle)
			assert.Equal(t, test.expectedRel, rel)
		})
	}
}

func TestRectangle_RelationshipToLineSegment(t *testing.T) {
	// Define a rectangle
	rect := NewRectangle([]Point[int]{
		NewPoint[int](0, 0),
		NewPoint[int](100, 0),
		NewPoint[int](100, 100),
		NewPoint[int](0, 100),
	})

	tests := map[string]struct {
		segment     LineSegment[int]
		expectedRel Relationship
	}{
		"Segment inside rectangle": {
			segment:     NewLineSegment(NewPoint[int](10, 10), NewPoint[int](90, 90)),
			expectedRel: RelationshipContains,
		},
		"Segment intersecting rectangle": {
			segment:     NewLineSegment(NewPoint[int](-10, 50), NewPoint[int](110, 50)),
			expectedRel: RelationshipIntersection,
		},
		"Segment outside rectangle": {
			segment:     NewLineSegment(NewPoint[int](200, 200), NewPoint[int](300, 300)),
			expectedRel: RelationshipDisjoint,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			rel := rect.RelationshipToLineSegment(test.segment)
			assert.Equal(t, test.expectedRel, rel)
		})
	}
}

func TestRectangle_RelationshipToPoint(t *testing.T) {
	// Define a rectangle
	rect := NewRectangle([]Point[int]{
		NewPoint[int](0, 0),
		NewPoint[int](100, 0),
		NewPoint[int](100, 100),
		NewPoint[int](0, 100),
	})

	tests := map[string]struct {
		point       Point[int]
		expectedRel Relationship
	}{
		"Point inside rectangle": {
			point:       NewPoint[int](50, 50),
			expectedRel: RelationshipContains,
		},
		"Point on rectangle edge": {
			point:       NewPoint[int](0, 50),
			expectedRel: RelationshipIntersection,
		},
		"Point on rectangle vertex": {
			point:       NewPoint[int](0, 0),
			expectedRel: RelationshipIntersection,
		},
		"Point outside rectangle": {
			point:       NewPoint[int](200, 200),
			expectedRel: RelationshipDisjoint,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			rel := rect.RelationshipToPoint(test.point, WithEpsilon(1e-10))
			assert.Equal(t, test.expectedRel, rel)
		})
	}
}

func TestRectangle_RelationshipToPolyTree(t *testing.T) {
	// Create a PolyTree
	root, err := NewPolyTree([]Point[int]{
		NewPoint(5, 5),
		NewPoint(5, 15),
		NewPoint(15, 15),
		NewPoint(15, 5),
	}, PTSolid)
	require.NoError(t, err, "error creating root polygon")

	hole, err := NewPolyTree([]Point[int]{
		NewPoint(6, 6),
		NewPoint(6, 9),
		NewPoint(9, 9),
		NewPoint(9, 6),
	}, PTHole)
	require.NoError(t, err, "error creating hole polygon")
	require.NoError(t, root.AddChild(hole), "error adding hole to root polygon")

	// Define test cases
	tests := []struct {
		name                     string
		rect                     Rectangle[int]
		pt                       *PolyTree[int]
		expectedRootRelationship Relationship
		expectedHoleRelationship Relationship
	}{
		{
			name: "Disjoint relationship",
			rect: NewRectangle([]Point[int]{
				NewPoint(20, 20),
				NewPoint(30, 20),
				NewPoint(30, 30),
				NewPoint(20, 30),
			}),
			pt:                       root,
			expectedRootRelationship: RelationshipDisjoint,
			expectedHoleRelationship: RelationshipDisjoint,
		},
		{
			name: "Intersection relationship",
			rect: NewRectangle([]Point[int]{
				NewPoint(13, 13),
				NewPoint(20, 13),
				NewPoint(20, 20),
				NewPoint(13, 20),
			}),
			pt:                       root,
			expectedRootRelationship: RelationshipIntersection,
			expectedHoleRelationship: RelationshipDisjoint,
		},
		{
			name: "Containment relationship",
			rect: NewRectangle([]Point[int]{
				NewPoint(4, 4),
				NewPoint(16, 4),
				NewPoint(16, 16),
				NewPoint(4, 16),
			}),
			pt:                       root,
			expectedRootRelationship: RelationshipContains,
			expectedHoleRelationship: RelationshipContains,
		},
		{
			name: "Contains relationship",
			rect: NewRectangle([]Point[int]{
				NewPoint(0, 0),
				NewPoint(20, 0),
				NewPoint(20, 20),
				NewPoint(0, 20),
			}),
			pt:                       root,
			expectedRootRelationship: RelationshipContains,
			expectedHoleRelationship: RelationshipContains,
		},
		{
			name: "ContainedBy relationship",
			rect: NewRectangle([]Point[int]{
				NewPoint(6, 6),
				NewPoint(14, 6),
				NewPoint(14, 14),
				NewPoint(6, 14),
			}),
			pt:                       root,
			expectedRootRelationship: RelationshipContainedBy,
			expectedHoleRelationship: RelationshipIntersection,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rels := tt.rect.RelationshipToPolyTree(tt.pt, WithEpsilon(1e-10))
			assert.Equal(t, tt.expectedRootRelationship, rels[root], "unexpected root relationship")
			assert.Equal(t, tt.expectedHoleRelationship, rels[hole], "unexpected hole relationship")
		})
	}
}

func TestRectangle_RelationshipToRectangle(t *testing.T) {
	// Define test rectangles
	r1 := NewRectangle([]Point[int]{
		NewPoint(0, 0),
		NewPoint(10, 0),
		NewPoint(10, 10),
		NewPoint(0, 10),
	})
	r2 := NewRectangle([]Point[int]{
		NewPoint(5, 5),
		NewPoint(15, 5),
		NewPoint(15, 15),
		NewPoint(5, 15),
	})
	r3 := NewRectangle([]Point[int]{
		NewPoint(2, 2),
		NewPoint(8, 2),
		NewPoint(8, 8),
		NewPoint(2, 8),
	})
	r4 := NewRectangle([]Point[int]{
		NewPoint(10, 10),
		NewPoint(20, 10),
		NewPoint(20, 20),
		NewPoint(10, 20),
	}) // Touching at a vertex
	r5 := NewRectangle([]Point[int]{
		NewPoint(20, 20),
		NewPoint(30, 20),
		NewPoint(30, 30),
		NewPoint(20, 30),
	}) // Disjoint

	tests := []struct {
		name        string
		r1, r2      Rectangle[int]
		expectedRel Relationship
	}{
		{
			name: "RelationshipEqual",
			r1:   r1,
			r2: NewRectangle([]Point[int]{
				NewPoint(0, 0),
				NewPoint(10, 0),
				NewPoint(10, 10),
				NewPoint(0, 10),
			}),
			expectedRel: RelationshipEqual,
		},
		{
			name:        "RelationshipIntersection",
			r1:          r1,
			r2:          r2,
			expectedRel: RelationshipIntersection,
		},
		{
			name:        "RelationshipContains",
			r1:          r1,
			r2:          r3,
			expectedRel: RelationshipContains,
		},
		{
			name:        "RelationshipContainedBy",
			r1:          r3,
			r2:          r1,
			expectedRel: RelationshipContainedBy,
		},
		{
			name:        "RelationshipDisjoint",
			r1:          r1,
			r2:          r5,
			expectedRel: RelationshipDisjoint,
		},
		{
			name:        "Touching Vertex",
			r1:          r1,
			r2:          r4,
			expectedRel: RelationshipIntersection,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rel := tt.r1.RelationshipToRectangle(tt.r2, WithEpsilon(1e-10))
			assert.Equal(t, tt.expectedRel, rel, "unexpected relationship")
		})
	}
}

func TestNewRectangleByPoints(t *testing.T) {
	tests := map[string]struct {
		corner              Point[int]
		oppositeCorner      Point[int]
		expectedTopLeft     Point[int]
		expectedBottomRight Point[int]
	}{
		"Standard case with top-left first": {
			corner:              NewPoint(0, 10),
			oppositeCorner:      NewPoint(10, 0),
			expectedTopLeft:     NewPoint(0, 10),
			expectedBottomRight: NewPoint(10, 0),
		},
		"Standard case with bottom-right first": {
			corner:              NewPoint(10, 0),
			oppositeCorner:      NewPoint(0, 10),
			expectedTopLeft:     NewPoint(0, 10),
			expectedBottomRight: NewPoint(10, 0),
		},
		"Negative coordinates": {
			corner:              NewPoint(-10, -10),
			oppositeCorner:      NewPoint(-5, -5),
			expectedTopLeft:     NewPoint(-10, -5),
			expectedBottomRight: NewPoint(-5, -10),
		},
		"Mixed positive and negative coordinates": {
			corner:              NewPoint(-10, 10),
			oppositeCorner:      NewPoint(5, -5),
			expectedTopLeft:     NewPoint(-10, 10),
			expectedBottomRight: NewPoint(5, -5),
		},
		"Single point rectangle": {
			corner:              NewPoint(5, 5),
			oppositeCorner:      NewPoint(5, 5),
			expectedTopLeft:     NewPoint(5, 5),
			expectedBottomRight: NewPoint(5, 5),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			rect := newRectangleByOppositeCorners(tt.corner, tt.oppositeCorner)

			assert.Equal(t, tt.expectedTopLeft, rect.topLeft, "Top-left corner mismatch")
			assert.Equal(t, tt.expectedBottomRight, rect.bottomRight, "Bottom-right corner mismatch")
		})
	}
}
