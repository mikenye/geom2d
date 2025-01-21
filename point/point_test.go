package point

import (
	"github.com/mikenye/geom2d/options"
	"github.com/stretchr/testify/assert"
	"image"
	"math"
	"testing"
)

func TestPoint_AsFloat32(t *testing.T) {
	// Define test cases with various point types
	tests := []struct {
		name     string
		point    Point[int]     // The point to convert
		expected Point[float32] // The expected result after conversion
	}{
		{
			name:     "Integer point conversion",
			point:    New(3, 4),
			expected: Point[float32]{3.0, 4.0},
		},
		{
			name:     "Negative integer point conversion",
			point:    New(-7, -5),
			expected: Point[float32]{-7.0, -5.0},
		},
		{
			name:     "Zero point conversion",
			point:    New(0, 0),
			expected: Point[float32]{0.0, 0.0},
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.point.AsFloat32()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPoint_AsFloat64(t *testing.T) {
	// Define test cases with various point types
	tests := []struct {
		name     string
		point    Point[int]     // The point to convert
		expected Point[float64] // The expected result after conversion
	}{
		{
			name:     "Integer point conversion",
			point:    New(3, 4),
			expected: Point[float64]{3.0, 4.0},
		},
		{
			name:     "Negative integer point conversion",
			point:    New(-7, -5),
			expected: Point[float64]{-7.0, -5.0},
		},
		{
			name:     "Zero point conversion",
			point:    New(0, 0),
			expected: Point[float64]{0.0, 0.0},
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.point.AsFloat64()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPoint_AsInt(t *testing.T) {
	// Define test cases with various point types
	tests := []struct {
		name     string
		point    Point[float64] // The point to convert
		expected Point[int]     // The expected result after conversion
	}{
		{
			name:     "Positive float point conversion",
			point:    Point[float64]{3.7, 4.9},
			expected: Point[int]{3, 4},
		},
		{
			name:     "Negative float point conversion",
			point:    Point[float64]{-7.3, -5.6},
			expected: Point[int]{-7, -5},
		},
		{
			name:     "Mixed sign float point conversion",
			point:    Point[float64]{-2.9, 3.1},
			expected: Point[int]{-2, 3},
		},
		{
			name:     "Zero point conversion",
			point:    Point[float64]{0.0, 0.0},
			expected: Point[int]{0, 0},
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.point.AsInt()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPoint_AsIntRounded(t *testing.T) {
	// Define test cases with various point types
	tests := []struct {
		name     string
		point    Point[float64] // The point to convert
		expected Point[int]     // The expected result after rounding conversion
	}{
		{
			name:     "Positive float point rounding up",
			point:    Point[float64]{3.7, 4.5},
			expected: Point[int]{4, 5},
		},
		{
			name:     "Positive float point rounding down",
			point:    Point[float64]{3.2, 4.4},
			expected: Point[int]{3, 4},
		},
		{
			name:     "Negative float point rounding up",
			point:    Point[float64]{-7.6, -5.5},
			expected: Point[int]{-8, -6},
		},
		{
			name:     "Negative float point rounding down",
			point:    Point[float64]{-2.2, -3.4},
			expected: Point[int]{-2, -3},
		},
		{
			name:     "Exact half values",
			point:    Point[float64]{2.5, -2.5},
			expected: Point[int]{3, -3},
		},
		{
			name:     "Zero point conversion",
			point:    Point[float64]{0.0, 0.0},
			expected: Point[int]{0, 0},
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.point.AsIntRounded()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPoint_DotProduct(t *testing.T) {
	tests := []struct {
		name     string
		p, q     any // Use `any` to handle different Point types
		expected any // Expected result can be int or float64
	}{
		// Integer points
		{
			name:     "int: (2,3) . (4,5)",
			p:        New(2, 3),
			q:        New(4, 5),
			expected: 23,
		},
		{
			name:     "int: (3,4) . (1,2)",
			p:        New(3, 4),
			q:        New(1, 2),
			expected: 11,
		},

		// Float64 points
		{
			name:     "float64: (2.0,3.0) . (4.0,5.0)",
			p:        New(2.0, 3.0),
			q:        New(4.0, 5.0),
			expected: 23.0,
		},
		{
			name:     "float64: (1.5,2.5) . (3.5,4.5)",
			p:        New(1.5, 2.5),
			q:        New(3.5, 4.5),
			expected: 16.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch p := tt.p.(type) {
			case Point[int]:
				q := tt.q.(Point[int])
				expected := tt.expected.(int)
				actual := p.DotProduct(q)
				assert.Equal(t, expected, actual)

			case Point[float64]:
				q := tt.q.(Point[float64])
				expected := tt.expected.(float64)
				actual := p.DotProduct(q)
				assert.Equal(t, expected, actual)
			}
		})
	}
}

func TestPoint_Eq(t *testing.T) {
	tests := map[string]struct {
		p, q     Point[float64] // Supports different Point types with `any`
		opts     []options.GeometryOptionsFunc
		expected bool // Expected result of equality comparison
	}{
		"(2.0,3.0) == (4.0,5.0)": {
			p:        New(2.0, 3.0),
			q:        New(4.0, 5.0),
			expected: false,
		},
		"(2.0,3.0) == (2.0,3.0)": {
			p:        New(2.0, 3.0),
			q:        New(2.0, 3.0),
			expected: true,
		},
		"(2.0,3.0) ~= (1.999999999,2.999999999)": {
			p:        New(2.0, 3.0),
			q:        New(1.999999999, 2.999999999),
			opts:     []options.GeometryOptionsFunc{options.WithEpsilon(1e-8)},
			expected: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tc.p.Eq(tc.q, tc.opts...)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestPoint_Rotate(t *testing.T) {
	tests := map[string]struct {
		point    Point[float64] // The point to rotate
		origin   Point[float64] // The origin for RotateFrom (or (0,0) for Rotate)
		angle    float64        // The angle in radians
		opts     []options.GeometryOptionsFunc
		expected Point[float64] // The expected result
		inDelta  float64
	}{
		"rotate 90 degrees around origin": {
			point:    New[float64](1, 0),
			origin:   New[float64](0, 0),
			angle:    math.Pi / 2,
			expected: New[float64](0, 1),
			opts:     []options.GeometryOptionsFunc{options.WithEpsilon(1e-9)},
			inDelta:  0.0001,
		},
		"rotate 180 degrees around origin": {
			point:    New[float64](1, 1),
			origin:   New[float64](0, 0),
			angle:    math.Pi,
			expected: New[float64](-1, -1),
			opts:     []options.GeometryOptionsFunc{options.WithEpsilon(1e-9)},
			inDelta:  0.0001,
		},
		"rotate 90 degrees around (1,1)": {
			point:    New[float64](2, 1),
			origin:   New[float64](1, 1),
			angle:    math.Pi / 2,
			expected: New[float64](1, 2),
			opts:     []options.GeometryOptionsFunc{options.WithEpsilon(1e-9)},
			inDelta:  0.0001,
		},
		"rotate 45 degrees around (2,2)": {
			point:    New[float64](3, 2),
			origin:   New[float64](2, 2),
			angle:    math.Pi / 4,
			expected: New[float64](2.7071, 2.7071),
			opts:     []options.GeometryOptionsFunc{options.WithEpsilon(1e-4)},
			inDelta:  0.0001,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.point.Rotate(tc.origin, tc.angle, tc.opts...)
			assert.InDelta(t, tc.expected.x, result.x, tc.inDelta)
			assert.InDelta(t, tc.expected.y, result.y, tc.inDelta)
		})
	}
}

func TestPoint_Negate(t *testing.T) {
	p := New(1, 2)
	assert.Equal(t, New(-1, -2), p.Negate())
}

func TestPoint_Scale(t *testing.T) {
	tests := map[string]struct {
		point    any            // Point to be scaled (can be int or float64)
		refPoint any            // Reference point for scaling (can be int or float64)
		scale    float64        // Scaling factor
		expected Point[float64] // Expected resulting point (float64 type)
	}{
		// Integer test cases
		"int: Scale point from reference by factor 2": {
			point:    New(3, 4),
			refPoint: New(1, 1),
			scale:    2.0,
			expected: New[float64](5, 7),
		},
		"int: Scale point from reference by factor 0.5": {
			point:    New(3, 4),
			refPoint: New(1, 1),
			scale:    0.5,
			expected: New[float64](2, 2.5),
		},

		// Float64 test cases
		"float64: Scale point from reference by factor 1.5": {
			point:    New(2.0, 3.0),
			refPoint: New(1.0, 1.0),
			scale:    1.5,
			expected: New[float64](2.5, 4.0),
		},
		"float64: Scale point from reference by factor 0.25": {
			point:    New(4.0, 8.0),
			refPoint: New(2.0, 2.0),
			scale:    0.25,
			expected: New[float64](2.5, 3.5),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			switch point := tt.point.(type) {
			case Point[int]:
				ref := tt.refPoint.(Point[int])
				result := point.AsFloat64().Scale(ref.AsFloat64(), tt.scale)
				assert.InDelta(t, tt.expected.x, result.x, 0.001)
				assert.InDelta(t, tt.expected.y, result.y, 0.001)

			case Point[float64]:
				ref := tt.refPoint.(Point[float64])
				result := point.Scale(ref, tt.scale)
				assert.InDelta(t, tt.expected.x, result.x, 0.001)
				assert.InDelta(t, tt.expected.y, result.y, 0.001)
			}
		})
	}
}

func TestPoint_String(t *testing.T) {
	tests := map[string]struct {
		p        any    // Supports different Point types with `any`
		expected string // Expected string representation of the point
	}{
		// Integer points
		"int: (1,2)": {
			p:        New(1, 2),
			expected: "(1,2)",
		},
		"int: (0,-3)": {
			p:        New(0, -3),
			expected: "(0,-3)",
		},

		// Float64 points
		"float64: (1.2,3.4)": {
			p:        New(1.2, 3.4),
			expected: "(1.2,3.4)",
		},
		"float64: (-1.5,-2.5)": {
			p:        New(-1.5, -2.5),
			expected: "(-1.5,-2.5)",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			switch p := tc.p.(type) {
			case Point[int]:
				actual := p.String()
				assert.Equal(t, tc.expected, actual)

			case Point[float64]:
				actual := p.String()
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestPoint_Translate(t *testing.T) {
	tests := []struct {
		name     string
		p, q     any // Use 'any' to support different Point types
		expected any
	}{
		// Integer points
		{
			name:     "int: (1,2)+(0,0)",
			p:        New(1, 2),
			q:        New(0, 0),
			expected: New(1, 2),
		},
		{
			name:     "int: (1,2)+(3,4)",
			p:        New(1, 2),
			q:        New(3, 4),
			expected: New(4, 6),
		},
		{
			name:     "int: (-1,-2)+(3,4)",
			p:        New(-1, -2),
			q:        New(3, 4),
			expected: New(2, 2),
		},

		// Float64 points
		{
			name:     "float64: (1.0,2.0)+(3.0,4.0)",
			p:        New(1.0, 2.0),
			q:        New(3.0, 4.0),
			expected: New(4.0, 6.0),
		},
		{
			name:     "float64: (-1.5,-2.5)+(3.5,4.5)",
			p:        New(-1.5, -2.5),
			q:        New(3.5, 4.5),
			expected: New(2.0, 2.0),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			switch p := tc.p.(type) {
			case Point[int]:
				q := tc.q.(Point[int])
				expected := tc.expected.(Point[int])
				actual := p.Translate(q)
				assert.Equal(t, expected, actual)

			case Point[float64]:
				q := tc.q.(Point[float64])
				expected := tc.expected.(Point[float64])
				actual := p.Translate(q)
				assert.Equal(t, expected, actual)
			}
		})
	}
}

func TestPoint_X(t *testing.T) {
	tests := []struct {
		name     string
		point    any // Supports different Point types with `any`
		expected any // Expected x-coordinate value as either int or float64
	}{
		// Integer points
		{
			name:     "int: Positive coordinates",
			point:    New(3, 4),
			expected: 3,
		},
		{
			name:     "int: Negative coordinates",
			point:    New(-7, -5),
			expected: -7,
		},
		{
			name:     "int: Zero x-coordinate",
			point:    New(0, 4),
			expected: 0,
		},

		// Float64 points
		{
			name:     "float64: Positive coordinates",
			point:    New(3.5, 4.5),
			expected: 3.5,
		},
		{
			name:     "float64: Negative coordinates",
			point:    New(-7.1, -5.2),
			expected: -7.1,
		},
		{
			name:     "float64: Zero x-coordinate",
			point:    New(0.0, 4.5),
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch point := tt.point.(type) {
			case Point[int]:
				expected := tt.expected.(int)
				assert.Equal(t, expected, point.X())

			case Point[float64]:
				expected := tt.expected.(float64)
				assert.Equal(t, expected, point.X())
			}
		})
	}
}

func TestPoint_Y(t *testing.T) {
	tests := []struct {
		name     string
		point    any // Supports different Point types with `any`
		expected any // Expected y-coordinate value as either int or float64
	}{
		// Integer points
		{
			name:     "int: Positive coordinates",
			point:    New(3, 4),
			expected: 4,
		},
		{
			name:     "int: Negative coordinates",
			point:    New(-7, -5),
			expected: -5,
		},
		{
			name:     "int: Zero y-coordinate",
			point:    New(3, 0),
			expected: 0,
		},

		// Float64 points
		{
			name:     "float64: Positive coordinates",
			point:    New(3.5, 4.5),
			expected: 4.5,
		},
		{
			name:     "float64: Negative coordinates",
			point:    New(-7.1, -5.2),
			expected: -5.2,
		},
		{
			name:     "float64: Zero y-coordinate",
			point:    New(3.0, 0.0),
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch point := tt.point.(type) {
			case Point[int]:
				expected := tt.expected.(int)
				assert.Equal(t, expected, point.Y())

			case Point[float64]:
				expected := tt.expected.(float64)
				assert.Equal(t, expected, point.Y())
			}
		})
	}
}

func TestNewPointFromImagePoint(t *testing.T) {
	// Define test cases with various image.Point values
	tests := []struct {
		name     string
		imgPoint image.Point // The image.Point to convert
		expected Point[int]  // The expected Point[int] result after conversion
	}{
		{
			name:     "Positive coordinates",
			imgPoint: image.Point{X: 10, Y: 20},
			expected: Point[int]{10, 20},
		},
		{
			name:     "Negative coordinates",
			imgPoint: image.Point{X: -15, Y: -25},
			expected: Point[int]{-15, -25},
		},
		{
			name:     "Zero coordinates",
			imgPoint: image.Point{X: 0, Y: 0},
			expected: Point[int]{0, 0},
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewFromImagePoint(tt.imgPoint)
			assert.Equal(t, tt.expected, result)
		})
	}
}
