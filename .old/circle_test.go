package _old

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func TestCircle_BoundingBox(t *testing.T) {
	tests := map[string]struct {
		circle   Circle[int]
		expected Rectangle[int]
	}{
		"Unit Circle": {
			circle: NewCircle(NewPoint(0, 0), 1),
			expected: NewRectangle([]Point[int]{
				NewPoint(-1, -1),
				NewPoint(1, -1),
				NewPoint(1, 1),
				NewPoint(-1, 1),
			}),
		},
		"Circle at (10, 10) with radius 5": {
			circle: NewCircle(NewPoint(10, 10), 5),
			expected: NewRectangle([]Point[int]{
				NewPoint(5, 5),
				NewPoint(15, 5),
				NewPoint(15, 15),
				NewPoint(5, 15),
			}),
		},
		"Circle at (-10, -10) with radius 3": {
			circle: NewCircle(NewPoint(-10, -10), 3),
			expected: NewRectangle([]Point[int]{
				NewPoint(-13, -13),
				NewPoint(-7, -13),
				NewPoint(-7, -7),
				NewPoint(-13, -7),
			}),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			boundingBox := tt.circle.BoundingBox()
			assert.Equal(t, tt.expected, boundingBox)
		})
	}
}

func TestCircle_RelationshipToCircle(t *testing.T) {
	circleA := NewCircle(NewPoint[int](0, 0), 10)
	circleB := NewCircle(NewPoint[int](0, 0), 10)
	circleC := NewCircle(NewPoint[int](0, 0), 5)
	circleD := NewCircle(NewPoint[int](20, 0), 10)
	circleE := NewCircle(NewPoint[int](15, 0), 5)
	circleF := NewCircle(NewPoint[int](25, 0), 5)

	tests := map[string]struct {
		circle1  Circle[int]
		circle2  Circle[int]
		expected Relationship
	}{
		"Equal Circles": {
			circle1:  circleA,
			circle2:  circleB,
			expected: RelationshipEqual,
		},
		"ContainedBy Circle": {
			circle1:  circleC,
			circle2:  circleA,
			expected: RelationshipContainedBy,
		},
		"Contains Circle": {
			circle1:  circleA,
			circle2:  circleC,
			expected: RelationshipContains,
		},
		"Intersecting Circles": {
			circle1:  circleA,
			circle2:  circleD,
			expected: RelationshipIntersection,
		},
		"Tangential Circles (Intersection)": {
			circle1:  circleA,
			circle2:  circleE,
			expected: RelationshipIntersection,
		},
		"Disjoint Circles": {
			circle1:  circleA,
			circle2:  circleF,
			expected: RelationshipDisjoint,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := tt.circle1.RelationshipToCircle(tt.circle2, WithEpsilon(1e-10))
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCircle_RelationshipToLineSegment(t *testing.T) {
	circle := NewCircle(NewPoint[int](0, 0), 10)

	tests := map[string]struct {
		line     LineSegment[int]
		expected Relationship
	}{
		"Line Segment Outside Circle": {
			line:     NewLineSegment(NewPoint(15, 0), NewPoint(20, 0)),
			expected: RelationshipDisjoint,
		},
		"Line Segment Intersects Circle": {
			line:     NewLineSegment(NewPoint(0, -15), NewPoint(0, 15)),
			expected: RelationshipIntersection,
		},
		"Line Segment Fully Contained in Circle": {
			line:     NewLineSegment(NewPoint(5, 5), NewPoint(-5, -5)),
			expected: RelationshipContains,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := circle.RelationshipToLineSegment(tt.line, WithEpsilon(1e-10))
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCircle_RelationshipToPoint(t *testing.T) {
	circle := NewCircle(NewPoint[int](0, 0), 10)

	tests := map[string]struct {
		point    Point[int]
		expected Relationship
	}{
		"Point Outside Circle": {
			point:    NewPoint(15, 0),
			expected: RelationshipDisjoint,
		},
		"Point On Circle Boundary": {
			point:    NewPoint(10, 0),
			expected: RelationshipIntersection,
		},
		"Point Inside Circle": {
			point:    NewPoint(5, 0),
			expected: RelationshipContains,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := circle.RelationshipToPoint(tt.point, WithEpsilon(1e-10))
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCircle_RelationshipToPolyTree(t *testing.T) {
	// Create a PolyTree
	root, err := NewPolyTree([]Point[int]{
		NewPoint(0, 0), NewPoint(0, 100), NewPoint(100, 100), NewPoint(100, 0),
	}, PTSolid)
	require.NoError(t, err)
	hole, err := NewPolyTree([]Point[int]{
		NewPoint(20, 20), NewPoint(20, 80), NewPoint(80, 80), NewPoint(80, 20),
	}, PTHole)
	require.NoError(t, err)
	err = root.AddChild(hole)
	require.NoError(t, err)

	t.Run("Circle contains root polygon", func(t *testing.T) {
		circle := NewCircle(NewPoint[int](50, 50), 600)
		rels := circle.RelationshipToPolyTree(root)
		assert.Equal(t, RelationshipContains, rels[root])
		assert.Equal(t, RelationshipContains, rels[hole])
	})

	t.Run("Circle intersects root polygon", func(t *testing.T) {
		circle := NewCircle(NewPoint[int](0, 0), 5)
		rels := circle.RelationshipToPolyTree(root)
		assert.Equal(t, RelationshipIntersection, rels[root])
		assert.Equal(t, RelationshipDisjoint, rels[hole])
	})

	t.Run("Circle disjoint from root polygon", func(t *testing.T) {
		circle := NewCircle(NewPoint[int](200, 200), 10)
		rels := circle.RelationshipToPolyTree(root)
		assert.Equal(t, RelationshipDisjoint, rels[root])
	})

	t.Run("Polygon contains circle", func(t *testing.T) {
		circle := NewCircle(NewPoint[int](50, 50), 10)
		rels := circle.RelationshipToPolyTree(root)
		assert.Equal(t, RelationshipContainedBy, rels[root])
	})

	t.Run("Circle intersects hole", func(t *testing.T) {
		circle := NewCircle(NewPoint[int](50, 50), 70)
		rels := circle.RelationshipToPolyTree(root)
		assert.Equal(t, RelationshipContains, rels[hole])
		assert.Equal(t, RelationshipIntersection, rels[root])
	})
}

func TestCircle_RelationshipToRectangle(t *testing.T) {
	// Define a rectangle
	rect := NewRectangle([]Point[int]{
		NewPoint(0, 0),
		NewPoint(100, 0),
		NewPoint(100, 100),
		NewPoint(0, 100),
	})

	t.Run("Disjoint", func(t *testing.T) {
		circle := NewCircle(NewPoint(-50, -50), 10)
		assert.Equal(t, RelationshipDisjoint, circle.RelationshipToRectangle(rect), "Expected RelationshipDisjoint")
	})

	t.Run("Intersection", func(t *testing.T) {
		circle := NewCircle(NewPoint(50, 120), 30)
		assert.Equal(t, RelationshipIntersection, circle.RelationshipToRectangle(rect), "Expected RelationshipIntersection")
	})

	t.Run("Circle Contains Rectangle", func(t *testing.T) {
		circle := NewCircle(NewPoint(50, 50), 200)
		assert.Equal(t, RelationshipContains, circle.RelationshipToRectangle(rect), "Expected RelationshipContains")
	})

	t.Run("Rectangle Contains Circle", func(t *testing.T) {
		circle := NewCircle(NewPoint(50, 50), 20)
		assert.Equal(t, RelationshipContainedBy, circle.RelationshipToRectangle(rect), "Expected RelationshipContainedBy")
	})
}

func TestNewCircle(t *testing.T) {
	tests := map[string]struct {
		center   Point[float64]
		radius   float64
		expected Circle[float64]
	}{
		"positive center and radius": {
			center:   NewPoint[float64](3.0, 4.0),
			radius:   5.0,
			expected: Circle[float64]{center: NewPoint[float64](3.0, 4.0), radius: 5.0},
		},
		"zero center and radius": {
			center:   NewPoint[float64](0.0, 0.0),
			radius:   0.0,
			expected: Circle[float64]{center: NewPoint[float64](0.0, 0.0), radius: 0.0},
		},
		"negative center and positive radius": {
			center:   NewPoint[float64](-3.0, -4.0),
			radius:   5.0,
			expected: Circle[float64]{center: NewPoint[float64](-3.0, -4.0), radius: 5.0},
		},
		"positive center and negative radius (edge case)": {
			center:   NewPoint[float64](3.0, 4.0),
			radius:   -5.0,
			expected: Circle[float64]{center: NewPoint[float64](3.0, 4.0), radius: -5.0},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := NewCircle(tt.center, tt.radius)
			assert.Equal(t, tt.expected.center, result.center)
			assert.Equal(t, tt.expected.radius, result.radius)
		})
	}
}
