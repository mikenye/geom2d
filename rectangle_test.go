package geom2d

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"image"
	"testing"
)

func TestRectangle_Area(t *testing.T) {
	tests := map[string]struct {
		rect     Rectangle[int]
		expected int
	}{
		"standard rectangle": {
			rect:     newRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			expected: 200,
		},
		"rectangle with swapped corners": {
			rect:     newRectangleByOppositeCorners(NewPoint(10, 20), NewPoint(0, 0)),
			expected: 200,
		},
		"degenerate rectangle (zero width)": {
			rect:     newRectangleByOppositeCorners(NewPoint(5, 5), NewPoint(5, 15)),
			expected: 0,
		},
		"degenerate rectangle (zero height)": {
			rect:     newRectangleByOppositeCorners(NewPoint(5, 5), NewPoint(15, 5)),
			expected: 0,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.rect.Area()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestRectangle_AsFloat32(t *testing.T) {
	rect := newRectangleByOppositeCorners(NewPoint[int](1, 2), NewPoint[int](10, 20))
	expected := newRectangleByOppositeCorners(NewPoint[float32](1.0, 2.0), NewPoint[float32](10.0, 20.0))
	assert.Equal(t, expected, rect.AsFloat32())
}

func TestRectangle_AsFloat64(t *testing.T) {
	rect := newRectangleByOppositeCorners(NewPoint[int](1, 2), NewPoint[int](10, 20))
	expected := newRectangleByOppositeCorners(NewPoint[float64](1.0, 2.0), NewPoint[float64](10.0, 20.0))
	assert.Equal(t, expected, rect.AsFloat64())
}

func TestRectangle_AsInt(t *testing.T) {
	rect := newRectangleByOppositeCorners(NewPoint[float64](1.7, 2.9), NewPoint[float64](10.5, 20.3))
	expected := newRectangleByOppositeCorners(NewPoint(1, 2), NewPoint(10, 20))
	assert.Equal(t, expected, rect.AsInt())
}

func TestRectangle_AsIntRounded(t *testing.T) {
	rect := newRectangleByOppositeCorners(NewPoint[float64](1.7, 2.9), NewPoint[float64](10.5, 20.3))
	expected := newRectangleByOppositeCorners(NewPoint(2, 3), NewPoint(11, 20))
	assert.Equal(t, expected, rect.AsIntRounded())
}

func TestRectangle_ContainsPoint(t *testing.T) {
	tests := map[string]struct {
		rect     Rectangle[int]
		point    Point[int]
		expected bool
	}{
		"point inside rectangle": {
			rect:     newRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(5, 10),
			expected: true,
		},
		"point on rectangle boundary (top-left corner)": {
			rect:     newRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(0, 0),
			expected: true,
		},
		"point on rectangle boundary (bottom-right corner)": {
			rect:     newRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(10, 20),
			expected: true,
		},
		"point outside rectangle": {
			rect:     newRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(15, 10),
			expected: false,
		},
		"point on boundary (horizontal edge)": {
			rect:     newRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(5, 0),
			expected: true,
		},
		"point on boundary (vertical edge)": {
			rect:     newRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(10, 5),
			expected: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.rect.ContainsPoint(tt.point)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestRectangle_Height(t *testing.T) {
	tests := map[string]struct {
		rect     Rectangle[int]
		expected int
	}{
		"positive height": {
			rect:     newRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			expected: 20,
		},
		"negative height": {
			rect:     newRectangleByOppositeCorners(NewPoint(0, 20), NewPoint(10, 0)),
			expected: 20,
		},
		"zero height": {
			rect:     newRectangleByOppositeCorners(NewPoint(5, 5), NewPoint(15, 5)),
			expected: 0,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.rect.Height()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestRectangle_Perimeter(t *testing.T) {
	tests := map[string]struct {
		rect     Rectangle[int]
		expected int
	}{
		"standard rectangle": {
			rect:     newRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			expected: 60,
		},
		"rectangle with swapped corners": {
			rect:     newRectangleByOppositeCorners(NewPoint(10, 20), NewPoint(0, 0)),
			expected: 60,
		},
		"degenerate rectangle (zero width)": {
			rect:     newRectangleByOppositeCorners(NewPoint(5, 5), NewPoint(5, 15)),
			expected: 20,
		},
		"degenerate rectangle (zero height)": {
			rect:     newRectangleByOppositeCorners(NewPoint(5, 5), NewPoint(15, 5)),
			expected: 20,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.rect.Perimeter()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestRectangle_Points(t *testing.T) {
	tests := map[string]struct {
		rect     Rectangle[int]
		expected []Point[int]
	}{
		"rectangle with positive coordinates": {
			rect: newRectangleByOppositeCorners(NewPoint(2, 3), NewPoint(8, 6)),
			expected: []Point[int]{
				NewPoint(2, 3), // top-left
				NewPoint(8, 3), // top-right
				NewPoint(8, 6), // bottom-right
				NewPoint(2, 6), // bottom-left
			},
		},
		"rectangle with negative coordinates": {
			rect: newRectangleByOppositeCorners(NewPoint(-5, -5), NewPoint(0, 0)),
			expected: []Point[int]{
				NewPoint(-5, -5), // top-left
				NewPoint(0, -5),  // top-right
				NewPoint(0, 0),   // bottom-right
				NewPoint(-5, 0),  // bottom-left
			},
		},
		"rectangle with zero-width or height": {
			rect: newRectangleByOppositeCorners(NewPoint(1, 1), NewPoint(1, 4)),
			expected: []Point[int]{
				NewPoint(1, 1), // top-left
				NewPoint(1, 1), // top-right (overlapping)
				NewPoint(1, 4), // bottom-right
				NewPoint(1, 4), // bottom-left (overlapping)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tc.rect.Points()
			assert.Len(t, actual, 4)
			for _, expectedPoint := range tc.expected {
				assert.Contains(t, actual, expectedPoint)
			}
		})
	}
}

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

func TestRectangle_Scale(t *testing.T) {
	tests := map[string]struct {
		rect     Rectangle[float64]
		ref      Point[float64]
		k        float64
		expected Rectangle[float64]
	}{
		"scale by 2 from origin": {
			rect: NewRectangle([]Point[float64]{
				NewPoint[float64](0, 10),
				NewPoint[float64](10, 10),
				NewPoint[float64](0, 0),
				NewPoint[float64](10, 0),
			}),
			ref: NewPoint[float64](0, 0),
			k:   2,
			expected: NewRectangle([]Point[float64]{
				NewPoint[float64](0, 20),
				NewPoint[float64](20, 20),
				NewPoint[float64](0, 0),
				NewPoint[float64](20, 0),
			}),
		},
		"scale by 0.5 from center": {
			rect: NewRectangle([]Point[float64]{
				NewPoint[float64](-10, 10),
				NewPoint[float64](10, 10),
				NewPoint[float64](-10, -10),
				NewPoint[float64](10, -10),
			}),
			ref: NewPoint[float64](0, 0),
			k:   0.5,
			expected: NewRectangle([]Point[float64]{
				NewPoint[float64](-5, -5),
				NewPoint[float64](5, -5),
				NewPoint[float64](5, 5),
				NewPoint[float64](-5, 5),
			}),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			scaled := test.rect.Scale(test.ref, test.k)
			if !scaled.Eq(test.expected) {
				t.Errorf("expected %v, got %v", test.expected, scaled)
			}
		})
	}
}

func TestRectangle_ScaleWidthHeight(t *testing.T) {
	tests := map[string]struct {
		rect           Rectangle[int]
		scaleWidth     float64
		scaleHeight    float64
		expectedWidth  float64
		expectedHeight float64
	}{
		"scale both dimensions": {
			rect:           newRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			scaleWidth:     1.5,
			scaleHeight:    0.5,
			expectedWidth:  15.0,
			expectedHeight: 10.0,
		},
		"scale width only": {
			rect:           newRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			scaleWidth:     2.0,
			scaleHeight:    1.0,
			expectedWidth:  20.0,
			expectedHeight: 20.0,
		},
		"scale height only": {
			rect:           newRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			scaleWidth:     1.0,
			scaleHeight:    2.0,
			expectedWidth:  10.0,
			expectedHeight: 40.0,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			widthScaled := tt.rect.ScaleWidth(tt.scaleWidth)
			heightScaled := tt.rect.ScaleHeight(tt.scaleHeight)

			assert.InDelta(t, tt.expectedWidth, widthScaled.Width(), 0.001)
			assert.InDelta(t, tt.expectedHeight, heightScaled.Height(), 0.001)
		})
	}
}

func TestRectangle_ToImageRect(t *testing.T) {
	rect := newRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(100, 200))
	expected := image.Rect(0, 0, 100, 200)
	assert.Equal(t, expected, rect.ToImageRect())
}

func TestRectangle_Translate(t *testing.T) {
	tests := map[string]struct {
		inputRect    Rectangle[int]
		translateBy  Point[int]
		expectedRect Rectangle[int]
	}{
		"translate up and right": {
			inputRect: NewRectangle([]Point[int]{
				NewPoint(0, 10),
				NewPoint(10, 10),
				NewPoint(0, 0),
				NewPoint(10, 0),
			}),
			translateBy: NewPoint(5, 5),
			expectedRect: NewRectangle([]Point[int]{
				NewPoint(5, 15),
				NewPoint(15, 15),
				NewPoint(5, 5),
				NewPoint(15, 5),
			}),
		},
		"translate down and left": {
			inputRect: NewRectangle([]Point[int]{
				NewPoint(0, 10),
				NewPoint(10, 10),
				NewPoint(0, 0),
				NewPoint(10, 0),
			}),
			translateBy: NewPoint(-5, -5),
			expectedRect: NewRectangle([]Point[int]{
				NewPoint(-5, 5),
				NewPoint(5, 5),
				NewPoint(-5, -5),
				NewPoint(5, -5),
			}),
		},
		"translate by zero": {
			inputRect: NewRectangle([]Point[int]{
				NewPoint(0, 10),
				NewPoint(10, 10),
				NewPoint(0, 0),
				NewPoint(10, 0),
			}),
			translateBy: NewPoint(0, 0),
			expectedRect: NewRectangle([]Point[int]{
				NewPoint(0, 10),
				NewPoint(10, 10),
				NewPoint(0, 0),
				NewPoint(10, 0),
			}),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := test.inputRect.Translate(test.translateBy)
			if !actual.Eq(test.expectedRect) {
				t.Errorf("expected %v, got %v", test.expectedRect, actual)
			}
		})
	}
}

func TestRectangle_Width(t *testing.T) {
	tests := map[string]struct {
		rect     Rectangle[int]
		expected int
	}{
		"positive width": {
			rect:     newRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			expected: 10,
		},
		"negative width": {
			rect:     newRectangleByOppositeCorners(NewPoint(10, 0), NewPoint(0, 20)),
			expected: 10,
		},
		"zero width": {
			rect:     newRectangleByOppositeCorners(NewPoint(5, 5), NewPoint(5, 15)),
			expected: 0,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.rect.Width()
			assert.Equal(t, tt.expected, actual)
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

func TestNewRectangleFromImageRect(t *testing.T) {
	tests := map[string]struct {
		imgRect  image.Rectangle
		expected []Point[int]
	}{
		"standard image rectangle": {
			imgRect: image.Rect(0, 0, 100, 200),
			expected: []Point[int]{
				NewPoint(0, 0),
				NewPoint(0, 200),
				NewPoint(100, 200),
				NewPoint(100, 0),
			},
		},
		"negative coordinates": {
			imgRect: image.Rect(-50, -50, 50, 50),
			expected: []Point[int]{
				NewPoint(-50, -50),
				NewPoint(-50, 50),
				NewPoint(50, 50),
				NewPoint(50, -50),
			},
		},
		"degenerate rectangle": {
			imgRect: image.Rect(10, 10, 10, 10),
			expected: []Point[int]{
				NewPoint(10, 10),
				NewPoint(10, 10),
				NewPoint(10, 10),
				NewPoint(10, 10),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := NewRectangleFromImageRect(tc.imgRect)
			for _, p := range actual.Points() {
				assert.Contains(t, tc.expected, p)
			}
		})
	}
}
