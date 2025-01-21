package circle

import (
	"github.com/mikenye/geom2d/options"
	"github.com/mikenye/geom2d/point"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestCircle_Area(t *testing.T) {
	tests := map[string]struct {
		circle   Circle[float64]
		inDelta  float64
		expected float64
	}{
		"radius 1": {
			circle:   New[float64](0, 0, 1),
			inDelta:  0.0001,
			expected: math.Pi,
		},
		"radius 2": {
			circle:   New[float64](0, 0, 2),
			inDelta:  0.0001,
			expected: 4 * math.Pi,
		},
		"radius 0": {
			circle:   New[float64](0, 0, 0),
			inDelta:  0.0001,
			expected: 0,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.InDelta(t, tc.expected, tc.circle.Area(), tc.inDelta)
		})
	}
}

func TestCircle_AsFloat32(t *testing.T) {
	tests := map[string]struct {
		circle   Circle[int]
		expected Circle[float32]
	}{
		"integer center and radius": {
			circle:   New[int](3, 4, 5),
			expected: New[float32](3, 4, 5),
		},
		"zero center and radius": {
			circle:   NewFromPoint(point.New(0, 0), 0),
			expected: New[float32](0, 0, 0),
		},
		"negative center and radius": {
			circle:   New(-3, -4, 5),
			expected: NewFromPoint[float32](point.New[float32](-3, -4), 5),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.circle.AsFloat32()
			assert.Equal(t, tc.expected.center, result.center)
			assert.Equal(t, tc.expected.radius, result.radius)
		})
	}
}

func TestCircle_AsFloat64(t *testing.T) {
	tests := map[string]struct {
		circle   Circle[int]
		expected Circle[float64]
	}{
		"integer center and radius": {
			circle:   New[int](3, 4, 5),
			expected: New[float64](3, 4, 5),
		},
		"zero center and radius": {
			circle:   New[int](0, 0, 0),
			expected: New[float64](0, 0, 0),
		},
		"negative center and radius": {
			circle:   New[int](-3, -4, 5),
			expected: New[float64](-3, -4, 5),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.circle.AsFloat64()
			assert.Equal(t, tc.expected.center, result.center)
			assert.Equal(t, tc.expected.radius, result.radius)
		})
	}
}

func TestCircle_AsInt(t *testing.T) {
	tests := map[string]struct {
		circle   Circle[float64]
		expected Circle[int]
	}{
		"positive float center and radius": {
			circle:   New[float64](3.9, 4.5, 5.8),
			expected: New[int](3, 4, 5),
		},
		"zero center and radius": {
			circle:   New[float64](0, 0, 0),
			expected: New[int](0, 0, 0),
		},
		"negative float center and radius": {
			circle:   New[float64](-3.7, -4.2, 5.9),
			expected: New[int](-3, -4, 5),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.circle.AsInt()
			assert.Equal(t, tc.expected.center, result.center)
			assert.Equal(t, tc.expected.radius, result.radius)
		})
	}
}

func TestCircle_AsIntRounded(t *testing.T) {
	tests := map[string]struct {
		circle   Circle[float64]
		expected Circle[int]
	}{
		"positive float center and radius with rounding up": {
			circle:   New[float64](3.6, 4.5, 5.7),
			expected: New[int](4, 5, 6),
		},
		"positive float center and radius with rounding down": {
			circle:   New[float64](3.4, 4.4, 5.2),
			expected: New[int](3, 4, 5),
		},
		"zero center and radius": {
			circle:   New[float64](0, 0, 0),
			expected: New[int](0, 0, 0),
		},
		"negative float center and radius with rounding up": {
			circle:   New[float64](-3.6, -4.5, 5.7),
			expected: New[int](-4, -5, 6),
		},
		"negative float center and radius with rounding down": {
			circle:   New[float64](-3.4, -4.4, 5.2),
			expected: New[int](-3, -4, 5),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.circle.AsIntRounded()
			assert.Equal(t, tc.expected.center, result.center)
			assert.Equal(t, tc.expected.radius, result.radius)
		})
	}
}

func TestCircle_Bresenham(t *testing.T) {
	tests := map[string]struct {
		center   point.Point[int]
		radius   int
		expected []point.Point[int]
	}{
		"circle at origin, radius 2": {
			center: point.New(0, 0),
			radius: 2,
			expected: []point.Point[int]{
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
			expected: []point.Point[int]{
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
			c.Bresenham(func(p point.Point[int]) bool {
				assert.Contains(t, tc.expected, p, "Points should match expected circle perimeter")
				return true
			})
		})
	}
}

func TestCircle_Center(t *testing.T) {
	tests := map[string]struct {
		circle   Circle[float64]
		expected point.Point[float64]
	}{
		"positive center coordinates": {
			circle:   New[float64](3.5, 4.5, 5.5),
			expected: point.New[float64](3.5, 4.5),
		},
		"zero center coordinates": {
			circle:   New[float64](0.0, 0.0, 5.5),
			expected: point.New[float64](0.0, 0.0),
		},
		"negative center coordinates": {
			circle:   New[float64](-3.5, -4.5, 5.5),
			expected: point.New[float64](-3.5, -4.5),
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
		circle   Circle[float64]
		expected float64
	}{
		"radius 1": {
			circle:   New[float64](0, 0, 1),
			expected: 2 * math.Pi,
		},
		"radius 2": {
			circle:   New[float64](0, 0, 2),
			expected: 4 * math.Pi,
		},
		"radius 0": {
			circle:   New[float64](0, 0, 0),
			expected: 0,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.InDelta(t, tc.expected, tc.circle.Circumference(), 0.0001)
		})
	}
}

func TestCircle_Eq(t *testing.T) {
	tests := map[string]struct {
		circle1  Circle[float64]
		circle2  Circle[float64]
		opts     []options.GeometryOptionsFunc
		expected bool
	}{
		"equal circles with same center and radius": {
			circle1:  New[float64](3, 4, 5),
			circle2:  New[float64](3, 4, 5),
			expected: true,
		},
		"different center but same radius": {
			circle1:  New[float64](3, 4, 5),
			circle2:  New[float64](2, 4, 5),
			expected: false,
		},
		"same center but different radius": {
			circle1:  New[float64](3, 4, 5),
			circle2:  New[float64](3, 4, 6),
			expected: false,
		},
		"different center and different radius": {
			circle1:  New[float64](3, 4, 5),
			circle2:  New[float64](2, 3, 6),
			expected: false,
		},
		"epsilon-equal circles with same center and radius": {
			circle1: New[float64](3, 4, 5),
			circle2: New[float64](2.999999999, 3.999999999, 4.999999999),
			opts: []options.GeometryOptionsFunc{
				options.WithEpsilon(1e-8),
			},
			expected: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.circle1.Eq(tc.circle2, tc.opts...)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestCircle_Radius(t *testing.T) {
	tests := map[string]struct {
		circle   Circle[float64]
		expected float64
	}{
		"positive radius": {
			circle:   New[float64](3, 4, 5),
			expected: 5,
		},
		"zero radius": {
			circle:   New[float64](3, 4, 0),
			expected: 0,
		},
		"small radius": {
			circle:   New[float64](3, 4, 0.001),
			expected: 0.001,
		},
		"negative radius (edge case)": {
			circle:   New[float64](3, 4, -5),
			expected: -5,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.circle.Radius()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestCircle_Rotate(t *testing.T) {
	tests := map[string]struct {
		circle   Circle[float64]
		pivot    point.Point[float64]
		radians  float64
		expected Circle[float64]
	}{
		"rotate 90 degrees around origin": {
			circle:  NewFromPoint(point.New[float64](3.0, 0.0), 5.0),
			pivot:   point.New[float64](0.0, 0.0),
			radians: math.Pi / 2,
			expected: NewFromPoint(
				point.New[float64](0.0, 3.0),
				5.0,
			),
		},
		"rotate 180 degrees around origin": {
			circle:  NewFromPoint(point.New[float64](3.0, 0.0), 5.0),
			pivot:   point.New[float64](0.0, 0.0),
			radians: math.Pi,
			expected: NewFromPoint(
				point.New[float64](-3.0, 0.0),
				5.0,
			),
		},
		"rotate 90 degrees around custom pivot": {
			circle:  NewFromPoint(point.New[float64](3.0, 0.0), 5.0),
			pivot:   point.New[float64](1.0, 1.0),
			radians: math.Pi / 2,
			expected: NewFromPoint(
				point.New[float64](2.0, 3.0),
				5.0,
			),
		},
		"rotate 0 degrees around custom pivot": {
			circle:  NewFromPoint(point.New[float64](3.0, 0.0), 5.0),
			pivot:   point.New[float64](1.0, 1.0),
			radians: 0,
			expected: NewFromPoint(
				point.New[float64](3.0, 0),
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
		circle   Circle[float64]
		factor   float64
		expected Circle[float64]
	}{
		"scale up by factor of 2": {
			circle:   New[float64](3, 4, 5),
			factor:   2,
			expected: New[float64](3, 4, 10),
		},
		"scale down by factor of 0.5": {
			circle:   New[float64](3, 4, 5),
			factor:   0.5,
			expected: New[float64](3, 4, 2.5),
		},
		"no change with factor of 1": {
			circle:   New[float64](3, 4, 5),
			factor:   1,
			expected: New[float64](3, 4, 5),
		},
		"scale to zero radius with factor of 0": {
			circle:   New[float64](3, 4, 5),
			factor:   0,
			expected: New[float64](3, 4, 0),
		},
		"scale with negative factor": {
			circle:   New[float64](3, 4, 5),
			factor:   -2,
			expected: New[float64](3, 4, -10),
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
		circle   Circle[float64]
		expected string
	}{
		"positive center and radius": {
			circle:   New[float64](3.5, 4.5, 5.5),
			expected: "(3.5,4.5,5.5)",
		},
		"zero center and radius": {
			circle:   New[float64](0, 0, 0),
			expected: "(0,0,0)",
		},
		"negative center and radius": {
			circle:   New[float64](-3.5, -4.5, -5.5),
			expected: "(-3.5,-4.5,-5.5)",
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
		circle   Circle[float64]
		vector   point.Point[float64]
		expected Circle[float64]
	}{
		"translate circle by positive vector": {
			circle:   New[float64](3, 4, 5),
			vector:   point.New[float64](2, 3),
			expected: New[float64](5, 7, 5),
		},
		"translate circle by negative vector": {
			circle:   New[float64](3, 4, 5),
			vector:   point.New[float64](-1, -2),
			expected: New[float64](2, 2, 5),
		},
		"translate circle by zero vector": {
			circle:   New[float64](3, 4, 5),
			vector:   point.New[float64](0, 0),
			expected: New[float64](3, 4, 5),
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
		xc, yc, x, y int
		expected     []point.Point[int]
	}{
		"center at origin, simple point": {
			xc: 0, yc: 0, x: 2, y: 1,
			expected: []point.Point[int]{
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
			expected: []point.Point[int]{
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
