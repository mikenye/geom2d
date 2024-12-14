package geom2d

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func TestCircle_Area(t *testing.T) {
	tests := map[string]struct {
		circle   Circle[float64]
		expected float64
	}{
		"radius 1": {circle: Circle[float64]{center: NewPoint[float64](0, 0), radius: 1}, expected: math.Pi},
		"radius 2": {circle: Circle[float64]{center: NewPoint[float64](0, 0), radius: 2}, expected: 4 * math.Pi},
		"radius 0": {circle: Circle[float64]{center: NewPoint[float64](0, 0), radius: 0}, expected: 0},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.InDelta(t, tt.expected, tt.circle.Area(), 0.001)
		})
	}
}

func TestCircle_AsFloat(t *testing.T) {
	tests := map[string]struct {
		circle   Circle[int]
		expected Circle[float64]
	}{
		"integer center and radius": {
			circle:   Circle[int]{center: NewPoint(3, 4), radius: 5},
			expected: Circle[float64]{center: NewPoint[float64](3.0, 4.0), radius: 5.0},
		},
		"zero center and radius": {
			circle:   Circle[int]{center: NewPoint(0, 0), radius: 0},
			expected: Circle[float64]{center: NewPoint[float64](0.0, 0.0), radius: 0.0},
		},
		"negative center and radius": {
			circle:   Circle[int]{center: NewPoint(-3, -4), radius: 5},
			expected: Circle[float64]{center: NewPoint[float64](-3.0, -4.0), radius: 5.0},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := tt.circle.AsFloat64()
			assert.Equal(t, tt.expected.center, result.center)
			assert.Equal(t, tt.expected.radius, result.radius)
		})
	}
}

func TestCircle_AsInt(t *testing.T) {
	tests := map[string]struct {
		circle   Circle[float64]
		expected Circle[int]
	}{
		"positive float center and radius": {
			circle:   Circle[float64]{center: NewPoint[float64](3.9, 4.5), radius: 5.8},
			expected: Circle[int]{center: NewPoint(3, 4), radius: 5},
		},
		"zero center and radius": {
			circle:   Circle[float64]{center: NewPoint[float64](0.0, 0.0), radius: 0.0},
			expected: Circle[int]{center: NewPoint(0, 0), radius: 0},
		},
		"negative float center and radius": {
			circle:   Circle[float64]{center: NewPoint[float64](-3.7, -4.2), radius: 5.9},
			expected: Circle[int]{center: NewPoint(-3, -4), radius: 5},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := tt.circle.AsInt()
			assert.Equal(t, tt.expected.center, result.center)
			assert.Equal(t, tt.expected.radius, result.radius)
		})
	}
}

func TestCircle_AsIntRounded(t *testing.T) {
	tests := map[string]struct {
		circle   Circle[float64]
		expected Circle[int]
	}{
		"positive float center and radius with rounding up": {
			circle:   Circle[float64]{center: NewPoint[float64](3.6, 4.5), radius: 5.7},
			expected: Circle[int]{center: NewPoint(4, 5), radius: 6},
		},
		"positive float center and radius with rounding down": {
			circle:   Circle[float64]{center: NewPoint[float64](3.4, 4.4), radius: 5.2},
			expected: Circle[int]{center: NewPoint(3, 4), radius: 5},
		},
		"zero center and radius": {
			circle:   Circle[float64]{center: NewPoint[float64](0.0, 0.0), radius: 0.0},
			expected: Circle[int]{center: NewPoint(0, 0), radius: 0},
		},
		"negative float center and radius with rounding up": {
			circle:   Circle[float64]{center: NewPoint[float64](-3.6, -4.5), radius: 5.7},
			expected: Circle[int]{center: NewPoint(-4, -5), radius: 6},
		},
		"negative float center and radius with rounding down": {
			circle:   Circle[float64]{center: NewPoint[float64](-3.4, -4.4), radius: 5.2},
			expected: Circle[int]{center: NewPoint(-3, -4), radius: 5},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := tt.circle.AsIntRounded()
			assert.Equal(t, tt.expected.center, result.center)
			assert.Equal(t, tt.expected.radius, result.radius)
		})
	}
}

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

func TestCircle_Center(t *testing.T) {
	tests := map[string]struct {
		circle   Circle[float64]
		expected Point[float64]
	}{
		"positive center coordinates": {
			circle:   Circle[float64]{center: NewPoint[float64](3.5, 4.5), radius: 5.5},
			expected: NewPoint[float64](3.5, 4.5),
		},
		"zero center coordinates": {
			circle:   Circle[float64]{center: NewPoint[float64](0.0, 0.0), radius: 5.5},
			expected: NewPoint[float64](0.0, 0.0),
		},
		"negative center coordinates": {
			circle:   Circle[float64]{center: NewPoint[float64](-3.5, -4.5), radius: 5.5},
			expected: NewPoint[float64](-3.5, -4.5),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := tt.circle.Center()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCircle_Circumference(t *testing.T) {
	tests := map[string]struct {
		circle   Circle[float64]
		expected float64
	}{
		"radius 1": {circle: Circle[float64]{center: NewPoint[float64](0, 0), radius: 1}, expected: 2 * math.Pi},
		"radius 2": {circle: Circle[float64]{center: NewPoint[float64](0, 0), radius: 2}, expected: 4 * math.Pi},
		"radius 0": {circle: Circle[float64]{center: NewPoint[float64](0, 0), radius: 0}, expected: 0},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.InDelta(t, tt.expected, tt.circle.Circumference(), 0.001)
		})
	}
}

func TestCircle_Eq(t *testing.T) {
	tests := map[string]struct {
		circle1  Circle[float64]
		circle2  Circle[float64]
		expected bool
	}{
		"equal circles with same center and radius": {
			circle1:  Circle[float64]{center: NewPoint[float64](3.0, 4.0), radius: 5.0},
			circle2:  Circle[float64]{center: NewPoint[float64](3.0, 4.0), radius: 5.0},
			expected: true,
		},
		"different center but same radius": {
			circle1:  Circle[float64]{center: NewPoint[float64](3.0, 4.0), radius: 5.0},
			circle2:  Circle[float64]{center: NewPoint[float64](2.0, 4.0), radius: 5.0},
			expected: false,
		},
		"same center but different radius": {
			circle1:  Circle[float64]{center: NewPoint[float64](3.0, 4.0), radius: 5.0},
			circle2:  Circle[float64]{center: NewPoint[float64](3.0, 4.0), radius: 6.0},
			expected: false,
		},
		"different center and different radius": {
			circle1:  Circle[float64]{center: NewPoint[float64](3.0, 4.0), radius: 5.0},
			circle2:  Circle[float64]{center: NewPoint[float64](2.0, 3.0), radius: 6.0},
			expected: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := tt.circle1.Eq(tt.circle2)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCircle_Radius(t *testing.T) {
	tests := map[string]struct {
		circle   Circle[float64]
		expected float64
	}{
		"positive radius": {
			circle:   Circle[float64]{center: NewPoint[float64](3, 4), radius: 5},
			expected: 5,
		},
		"zero radius": {
			circle:   Circle[float64]{center: NewPoint[float64](3, 4), radius: 0},
			expected: 0,
		},
		"small radius": {
			circle:   Circle[float64]{center: NewPoint[float64](3, 4), radius: 0.001},
			expected: 0.001,
		},
		"negative radius (edge case)": {
			circle:   Circle[float64]{center: NewPoint[float64](3, 4), radius: -5},
			expected: -5,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := tt.circle.Radius()
			assert.Equal(t, tt.expected, result)
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

func TestCircle_Rotate(t *testing.T) {
	tests := map[string]struct {
		circle   Circle[float64]
		pivot    Point[float64]
		radians  float64
		expected Circle[float64]
	}{
		"rotate 90 degrees around origin": {
			circle:  NewCircle(NewPoint[float64](3.0, 0.0), 5.0),
			pivot:   NewPoint[float64](0.0, 0.0),
			radians: math.Pi / 2,
			expected: NewCircle(
				NewPoint[float64](0.0, 3.0),
				5.0,
			),
		},
		"rotate 180 degrees around origin": {
			circle:  NewCircle(NewPoint[float64](3.0, 0.0), 5.0),
			pivot:   NewPoint[float64](0.0, 0.0),
			radians: math.Pi,
			expected: NewCircle(
				NewPoint[float64](-3.0, 0.0),
				5.0,
			),
		},
		"rotate 90 degrees around custom pivot": {
			circle:  NewCircle(NewPoint[float64](3.0, 0.0), 5.0),
			pivot:   NewPoint[float64](1.0, 1.0),
			radians: math.Pi / 2,
			expected: NewCircle(
				NewPoint[float64](2.0, 3.0),
				5.0,
			),
		},
		"rotate 0 degrees around custom pivot": {
			circle:  NewCircle(NewPoint[float64](3.0, 0.0), 5.0),
			pivot:   NewPoint[float64](1.0, 1.0),
			radians: 0,
			expected: NewCircle(
				NewPoint[float64](3.0, 0),
				5.0,
			),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := tt.circle.Rotate(tt.pivot, tt.radians)
			assert.InDelta(t, tt.expected.center.x, result.center.x, 0.0001)
			assert.InDelta(t, tt.expected.center.y, result.center.y, 0.0001)
			assert.Equal(t, tt.expected.radius, result.radius)
		})
	}
}

func TestCircle_Scale(t *testing.T) {
	tests := map[string]struct {
		circle   Circle[float64]
		factor   float64
		expected Circle[float64]
	}{
		"scale up by factor of 2": {
			circle:   Circle[float64]{center: NewPoint[float64](3, 4), radius: 5},
			factor:   2,
			expected: Circle[float64]{center: NewPoint[float64](3, 4), radius: 10},
		},
		"scale down by factor of 0.5": {
			circle:   Circle[float64]{center: NewPoint[float64](3, 4), radius: 5},
			factor:   0.5,
			expected: Circle[float64]{center: NewPoint[float64](3, 4), radius: 2.5},
		},
		"no change with factor of 1": {
			circle:   Circle[float64]{center: NewPoint[float64](3, 4), radius: 5},
			factor:   1,
			expected: Circle[float64]{center: NewPoint[float64](3, 4), radius: 5},
		},
		"scale to zero radius with factor of 0": {
			circle:   Circle[float64]{center: NewPoint[float64](3, 4), radius: 5},
			factor:   0,
			expected: Circle[float64]{center: NewPoint[float64](3, 4), radius: 0},
		},
		"scale with negative factor": {
			circle:   Circle[float64]{center: NewPoint[float64](3, 4), radius: 5},
			factor:   -2,
			expected: Circle[float64]{center: NewPoint[float64](3, 4), radius: -10},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := tt.circle.Scale(tt.factor)
			assert.Equal(t, tt.expected.center, result.center)
			assert.Equal(t, tt.expected.radius, result.radius)
		})
	}
}

func TestCircle_String(t *testing.T) {
	tests := map[string]struct {
		circle   Circle[float64]
		expected string
	}{
		"positive center and radius": {
			circle:   Circle[float64]{center: NewPoint[float64](3.5, 4.5), radius: 5.5},
			expected: "Circle[center=(3.5, 4.5), radius=5.5]",
		},
		"zero center and radius": {
			circle:   Circle[float64]{center: NewPoint[float64](0.0, 0.0), radius: 0.0},
			expected: "Circle[center=(0, 0), radius=0]",
		},
		"negative center and radius": {
			circle:   Circle[float64]{center: NewPoint[float64](-3.5, -4.5), radius: -5.5},
			expected: "Circle[center=(-3.5, -4.5), radius=-5.5]",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := tt.circle.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCircle_Translate(t *testing.T) {
	tests := map[string]struct {
		circle   Circle[float64]
		vector   Point[float64]
		expected Circle[float64]
	}{
		"translate circle by positive vector": {
			circle:   Circle[float64]{center: NewPoint[float64](3, 4), radius: 5},
			vector:   NewPoint[float64](2, 3),
			expected: Circle[float64]{center: NewPoint[float64](5, 7), radius: 5},
		},
		"translate circle by negative vector": {
			circle:   Circle[float64]{center: NewPoint[float64](3, 4), radius: 5},
			vector:   NewPoint[float64](-1, -2),
			expected: Circle[float64]{center: NewPoint[float64](2, 2), radius: 5},
		},
		"translate circle by zero vector": {
			circle:   Circle[float64]{center: NewPoint[float64](3, 4), radius: 5},
			vector:   NewPoint[float64](0, 0),
			expected: Circle[float64]{center: NewPoint[float64](3, 4), radius: 5},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.circle.Translate(tc.vector)
			assert.Equal(t, tc.expected.center, result.center)
			assert.Equal(t, tc.expected.radius, result.radius)
		})
	}
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
