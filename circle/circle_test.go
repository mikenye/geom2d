package circle

import (
	"encoding/json"
	"github.com/mikenye/geom2d"
	"github.com/mikenye/geom2d/point"
	"github.com/mikenye/geom2d/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func TestCircle_Area(t *testing.T) {
	tests := map[string]struct {
		circle   Circle
		expected float64
	}{
		"radius 1": {
			circle:   New(0, 0, 1),
			expected: math.Pi,
		},
		"radius 2": {
			circle:   New(0, 0, 2),
			expected: 4 * math.Pi,
		},
		"radius 0": {
			circle:   New(0, 0, 0),
			expected: 0,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.InDelta(t, tc.expected, tc.circle.Area(), geom2d.GetEpsilon())
		})
	}
}

func TestCircle_Bresenham(t *testing.T) {
	tests := map[string]struct {
		center   point.Point
		radius   float64
		expected []point.Point
	}{
		"circle at origin, radius 2": {
			center: point.New(0, 0),
			radius: 2,
			expected: []point.Point{
				point.New(-1, -2),
				point.New(-1, 2),
				point.New(-2, -1),
				point.New(-2, 0),
				point.New(-2, 1),
				point.New(0, -2),
				point.New(0, 2),
				point.New(1, -2),
				point.New(1, 2),
				point.New(2, -1),
				point.New(2, 0),
				point.New(2, 1),
			},
		},
		"circle offset, radius 3": {
			center: point.New(5, 5),
			radius: 3,
			expected: []point.Point{
				point.New(2, 4),
				point.New(2, 5),
				point.New(2, 6),
				point.New(3, 3),
				point.New(3, 7),
				point.New(4, 2),
				point.New(4, 8),
				point.New(5, 2),
				point.New(5, 8),
				point.New(6, 2),
				point.New(6, 8),
				point.New(7, 3),
				point.New(7, 7),
				point.New(8, 4),
				point.New(8, 5),
				point.New(8, 6),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			c := NewFromPoint(tc.center, tc.radius)
			c.Bresenham(func(p point.Point) bool {
				assert.Contains(t, tc.expected, p, "Points should match expected circle perimeter")
				return true
			})
		})
	}
}

func TestCircle_Center(t *testing.T) {
	tests := map[string]struct {
		circle   Circle
		expected point.Point
	}{
		"positive center coordinates": {
			circle:   New(3.5, 4.5, 5.5),
			expected: point.New(3.5, 4.5),
		},
		"zero center coordinates": {
			circle:   New(0.0, 0.0, 5.5),
			expected: point.New(0.0, 0.0),
		},
		"negative center coordinates": {
			circle:   New(-3.5, -4.5, 5.5),
			expected: point.New(-3.5, -4.5),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.circle.Center()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestCircle_Circumference(t *testing.T) {
	tests := map[string]struct {
		circle   Circle
		expected float64
	}{
		"radius 1": {
			circle:   New(0, 0, 1),
			expected: 2 * math.Pi,
		},
		"radius 2": {
			circle:   New(0, 0, 2),
			expected: 4 * math.Pi,
		},
		"radius 0": {
			circle:   New(0, 0, 0),
			expected: 0,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.InDelta(t, tc.expected, tc.circle.Circumference(), geom2d.GetEpsilon())
		})
	}
}

func TestCircle_Eq(t *testing.T) {
	tests := map[string]struct {
		circle1  Circle
		circle2  Circle
		expected bool
	}{
		"equal circles with same center and radius": {
			circle1:  New(3, 4, 5),
			circle2:  New(3, 4, 5),
			expected: true,
		},
		"different center but same radius": {
			circle1:  New(3, 4, 5),
			circle2:  New(2, 4, 5),
			expected: false,
		},
		"same center but different radius": {
			circle1:  New(3, 4, 5),
			circle2:  New(3, 4, 6),
			expected: false,
		},
		"different center and different radius": {
			circle1:  New(3, 4, 5),
			circle2:  New(2, 3, 6),
			expected: false,
		},
		"epsilon-equal circles with same center and radius": {
			circle1:  New(3, 4, 5),
			circle2:  New(2.9999999999999, 3.9999999999999, 4.9999999999999),
			expected: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.circle1.Eq(tc.circle2)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestCircle_MarshalUnmarshalJSON(t *testing.T) {
	tests := map[string]struct {
		circle   Circle // Input circle
		expected Circle // Expected output after Marshal -> Unmarshal
	}{
		"Circle[float64]": {
			circle:   NewFromPoint(point.New(3.5, 7.2), 2.8),
			expected: NewFromPoint(point.New(3.5, 7.2), 2.8),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Marshal
			data, err := json.Marshal(tc.circle)
			require.NoErrorf(t, err, "Failed to marshal %s: %v", tc.circle, err)

			// Determine the correct type for unmarshalling
			var result Circle
			err = json.Unmarshal(data, &result)
			require.NoErrorf(t, err, "Failed to unmarshal %s: %v", string(data), err)
			assert.Equalf(t, tc.expected, result, "Expected %v, got %v", tc.expected, result)
		})
	}
}

func TestCircle_Radius(t *testing.T) {
	tests := map[string]struct {
		circle   Circle
		expected float64
	}{
		"positive radius": {
			circle:   New(3, 4, 5),
			expected: 5,
		},
		"zero radius": {
			circle:   New(3, 4, 0),
			expected: 0,
		},
		"small radius": {
			circle:   New(3, 4, 0.001),
			expected: 0.001,
		},
		"negative radius (edge case)": {
			circle:   New(3, 4, -5),
			expected: 5,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.circle.Radius()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestCircle_RelationshipToPoint(t *testing.T) {
	testCases := map[string]struct {
		point       point.Point
		circle      Circle
		expectedRel types.Relationship
	}{
		"Point inside circle": {
			point:       point.New(2, 2),
			circle:      New(0, 0, 5),
			expectedRel: types.RelationshipContainedBy,
		},
		"Point on circle boundary": {
			point:       point.New(3, 4),
			circle:      New(0, 0, 5),
			expectedRel: types.RelationshipIntersection,
		},
		"Point outside circle": {
			point:       point.New(6, 8),
			circle:      New(0, 0, 5),
			expectedRel: types.RelationshipDisjoint,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expectedRel, tc.circle.RelationshipToPoint(tc.point), "unexpected relationship")
		})
	}
}

func TestCircle_Rotate(t *testing.T) {
	tests := map[string]struct {
		circle   Circle
		pivot    point.Point
		radians  float64
		expected Circle
	}{
		"rotate 90 degrees around origin": {
			circle:  NewFromPoint(point.New(3.0, 0.0), 5.0),
			pivot:   point.New(0.0, 0.0),
			radians: math.Pi / 2,
			expected: NewFromPoint(
				point.New(0.0, 3.0),
				5.0,
			),
		},
		"rotate 180 degrees around origin": {
			circle:  NewFromPoint(point.New(3.0, 0.0), 5.0),
			pivot:   point.New(0.0, 0.0),
			radians: math.Pi,
			expected: NewFromPoint(
				point.New(-3.0, 0.0),
				5.0,
			),
		},
		"rotate 90 degrees around custom pivot": {
			circle:  NewFromPoint(point.New(3.0, 0.0), 5.0),
			pivot:   point.New(1.0, 1.0),
			radians: math.Pi / 2,
			expected: NewFromPoint(
				point.New(2.0, 3.0),
				5.0,
			),
		},
		"rotate 0 degrees around custom pivot": {
			circle:  NewFromPoint(point.New(3.0, 0.0), 5.0),
			pivot:   point.New(1.0, 1.0),
			radians: 0,
			expected: NewFromPoint(
				point.New(3.0, 0),
				5.0,
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.circle.Rotate(tc.pivot, tc.radians)
			assert.InDelta(t, tc.expected.center.X(), result.center.X(), 0.0001)
			assert.InDelta(t, tc.expected.center.Y(), result.center.Y(), 0.0001)
			assert.Equal(t, tc.expected.radius, result.radius)
		})
	}
}

func TestCircle_Scale(t *testing.T) {
	tests := map[string]struct {
		circle   Circle
		factor   float64
		expected Circle
	}{
		"scale up by factor of 2": {
			circle:   New(3, 4, 5),
			factor:   2,
			expected: New(3, 4, 10),
		},
		"scale down by factor of 0.5": {
			circle:   New(3, 4, 5),
			factor:   0.5,
			expected: New(3, 4, 2.5),
		},
		"no change with factor of 1": {
			circle:   New(3, 4, 5),
			factor:   1,
			expected: New(3, 4, 5),
		},
		"scale to zero radius with factor of 0": {
			circle:   New(3, 4, 5),
			factor:   0,
			expected: New(3, 4, 0),
		},
		"scale with negative factor": {
			circle:   New(3, 4, 5),
			factor:   -2,
			expected: New(3, 4, 10),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.circle.Scale(tc.factor)
			assert.Equal(t, tc.expected.center, result.center)
			assert.Equal(t, tc.expected.radius, result.radius)
		})
	}
}

func TestCircle_String(t *testing.T) {
	tests := map[string]struct {
		circle   Circle
		expected string
	}{
		"positive center and radius": {
			circle:   New(3.5, 4.5, 5.5),
			expected: "(3.500000,4.500000; r=5.500000)",
		},
		"zero center and radius": {
			circle:   New(0, 0, 0),
			expected: "(0.000000,0.000000; r=0.000000)",
		},
		"negative center and radius": {
			circle:   New(-3.5, -4.5, -5.5),
			expected: "(-3.500000,-4.500000; r=5.500000)",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.circle.String()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestCircle_Translate(t *testing.T) {
	tests := map[string]struct {
		circle   Circle
		vector   point.Point
		expected Circle
	}{
		"translate circle by positive vector": {
			circle:   New(3, 4, 5),
			vector:   point.New(2, 3),
			expected: New(5, 7, 5),
		},
		"translate circle by negative vector": {
			circle:   New(3, 4, 5),
			vector:   point.New(-1, -2),
			expected: New(2, 2, 5),
		},
		"translate circle by zero vector": {
			circle:   New(3, 4, 5),
			vector:   point.New(0, 0),
			expected: New(3, 4, 5),
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

func TestReflectAcrossCircleOctants(t *testing.T) {
	tests := map[string]struct {
		xc, yc, x, y float64
		expected     []point.Point
	}{
		"center at origin, simple point": {
			xc: 0, yc: 0, x: 2, y: 1,
			expected: []point.Point{
				point.New(2, 1),   // Octant 1
				point.New(-2, 1),  // Octant 2
				point.New(2, -1),  // Octant 8
				point.New(-2, -1), // Octant 7
				point.New(1, 2),   // Octant 3
				point.New(-1, 2),  // Octant 4
				point.New(1, -2),  // Octant 6
				point.New(-1, -2), // Octant 5
			},
		},
		"center offset, simple point": {
			xc: 3, yc: 4, x: 2, y: 1,
			expected: []point.Point{
				point.New(5, 5), // Octant 1
				point.New(1, 5), // Octant 2
				point.New(5, 3), // Octant 8
				point.New(1, 3), // Octant 7
				point.New(4, 6), // Octant 3
				point.New(2, 6), // Octant 4
				point.New(4, 2), // Octant 6
				point.New(2, 2), // Octant 5
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := reflectAcrossCircleOctants(tc.xc, tc.yc, tc.x, tc.y)
			assert.Equal(t, tc.expected, actual, "Points should match expected octant reflections")
		})
	}
}
