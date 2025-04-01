package point

import (
	"encoding/json"
	"github.com/mikenye/geom2d"
	"github.com/mikenye/geom2d/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"image"
	"math"
	"testing"
)

func TestPoint_AngleBetween(t *testing.T) {
	tests := map[string]struct {
		origin, a, b    Point   // The points for the test
		expected        float64 // The expected angle in radians
		shouldReturnNaN bool    // Whether the result should be NaN
	}{
		"basic angle between points": {
			origin:          New(0, 0),
			a:               New(1, 0),
			b:               New(0, 1),
			expected:        math.Pi / 2, // 90 degrees
			shouldReturnNaN: false,
		},
		"collinear points": {
			origin:          New(0, 0),
			a:               New(1, 1),
			b:               New(-1, -1),
			expected:        math.Pi, // 180 degrees
			shouldReturnNaN: false,
		},
		"identical points": {
			origin:          New(0, 0),
			a:               New(1, 1),
			b:               New(1, 1),
			expected:        0,
			shouldReturnNaN: false,
		},
		"zero vector (a or b equal to origin)": {
			origin:          New(0, 0),
			a:               New(0, 0),
			b:               New(1, 1),
			expected:        math.NaN(),
			shouldReturnNaN: true,
		},
		"small angle": {
			origin:          New(0, 0),
			a:               New(1, 0),
			b:               New(1, 0.01),
			expected:        0.01, // Small angle in radians
			shouldReturnNaN: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// global epsilon (default 1e-12) may be too strict for AngleBetween,
			// since floating-point trigonometric functions (like math.Acos) introduce larger errors.
			epsilon := 1e-6

			t.Logf("calculating angle between %s and %s about %s", tc.a, tc.b, tc.origin)
			result := tc.origin.AngleBetween(tc.a, tc.b)

			if tc.shouldReturnNaN {
				assert.True(t, math.IsNaN(result), "expected NaN but got %v", result)
			} else {
				assert.InDelta(t, tc.expected, result, epsilon, "unexpected angle")
			}
		})
	}
}

func TestPoint_Coordinates(t *testing.T) {
	tests := map[string]struct {
		point Point
		wantX float64
		wantY float64
	}{
		"origin":          {New(0, 0), 0, 0},
		"positive values": {New(3, 4), 3, 4},
		"negative values": {New(-5, -10), -5, -10},
		"mixed values":    {New(-7, 9), -7, 9},
		"large values":    {New(1000000, -999999), 1000000, -999999},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			x, y := tc.point.Coordinates()
			assert.Equal(t, tc.wantX, x, "X coordinate mismatch")
			assert.Equal(t, tc.wantY, y, "Y coordinate mismatch")
		})
	}
}

func TestPoint_CrossProduct(t *testing.T) {
	tests := []struct {
		name     string
		p, q     Point   // Support different Point types with `any`
		expected float64 // Expected result for different types
	}{
		{
			name:     "float64: (2.0,3.0) x (4.0,5.0)",
			p:        New(2.0, 3.0),
			q:        New(4.0, 5.0),
			expected: -2.0,
		},
		{
			name:     "float64: (3.5,2.5) x (4.0,6.0)",
			p:        New(3.5, 2.5),
			q:        New(4.0, 6.0),
			expected: 11.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := tt.expected
			actual := tt.p.CrossProduct(tt.q)
			assert.Equal(t, expected, actual)
		})
	}
}

func TestPoint_DistanceToPoint(t *testing.T) {
	tests := []struct {
		name     string
		p, q     Point
		expected float64
	}{
		{
			name:     "float64: distance between (2.0,10.0) and (10.0,2.0)",
			p:        New(2.0, 10.0),
			q:        New(10.0, 2.0),
			expected: math.Sqrt(((2 - 10) * (2 - 10)) + ((10 - 2) * (10 - 2))),
		},
		{
			name:     "float64: distance between (0.0,0.0) and (3.0,4.0)",
			p:        New(0.0, 0.0),
			q:        New(3.0, 4.0),
			expected: math.Sqrt(((0 - 3) * (0 - 3)) + ((0 - 4) * (0 - 4))),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.p.DistanceToPoint(tt.q)
			assert.InDelta(t, tt.expected, actual, geom2d.GetEpsilon())
		})
	}
}

func TestPoint_DotProduct(t *testing.T) {
	tests := []struct {
		name     string
		p, q     Point
		expected float64
	}{
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
			actual := tt.p.DotProduct(tt.q)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestPoint_Eq(t *testing.T) {
	tests := map[string]struct {
		p, q     Point
		expected bool
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
		"(0.3, 0.3) ~= (0.2+0.1, 0.2+0.1)": {
			p:        New(0.2+0.1, 0.2+0.1),
			q:        New(0.3, 0.3),
			expected: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tc.p.Eq(tc.q)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestPoint_Rotate(t *testing.T) {
	tests := map[string]struct {
		point    Point   // The point to rotate
		origin   Point   // The origin for RotateFrom (or (0,0) for Rotate)
		angle    float64 // The angle in radians
		expected Point   // The expected result
	}{
		"rotate 90 degrees around origin": {
			point:    New(1, 0),
			origin:   New(0, 0),
			angle:    math.Pi / 2,
			expected: New(0, 1),
		},
		"rotate 180 degrees around origin": {
			point:    New(1, 1),
			origin:   New(0, 0),
			angle:    math.Pi,
			expected: New(-1, -1),
		},
		"rotate 90 degrees around (1,1)": {
			point:    New(2, 1),
			origin:   New(1, 1),
			angle:    math.Pi / 2,
			expected: New(1, 2),
		},
		"rotate 45 degrees around (2,2)": {
			point:    New(3, 2),
			origin:   New(2, 2),
			angle:    math.Pi / 4,
			expected: New(2.7071067811865475, 2.7071067811865475),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.point.Rotate(tc.origin, tc.angle)
			assert.InDelta(t, tc.expected.x, result.x, geom2d.GetEpsilon())
			assert.InDelta(t, tc.expected.y, result.y, geom2d.GetEpsilon())
		})
	}
}

func TestPoint_MarshalUnmarshalJSON(t *testing.T) {
	tests := map[string]struct {
		point    Point
		expected Point
	}{
		"Point": {
			point:    New(3.5, 7.2),
			expected: New(3.5, 7.2),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Marshal
			data, err := json.Marshal(tc.point)
			require.NoErrorf(t, err, "Failed to marshal %s: %v", tc.point, err)

			var result Point
			err = json.Unmarshal(data, &result)
			require.NoErrorf(t, err, "Failed to unmarshal `%s`: %v", string(data), err)
			assert.Equalf(t, tc.expected, result, "Expected %v, got %v", tc.expected, result)

		})
	}
}

func TestPoint_Negate(t *testing.T) {
	p := New(1, 2)
	assert.Equal(t, New(-1, -2), p.Negate())
}

func TestPoint_RelationshipToPoint(t *testing.T) {
	tests := map[string]struct {
		pointA      Point
		pointB      Point
		expectedRel types.Relationship
	}{
		"Points are equal": {
			pointA:      New(5, 5),
			pointB:      New(5, 5),
			expectedRel: types.RelationshipEqual,
		},
		"Points are disjoint": {
			pointA:      New(5, 5),
			pointB:      New(10, 10),
			expectedRel: types.RelationshipDisjoint,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rel := tc.pointA.RelationshipToPoint(tc.pointB)
			assert.Equal(t, tc.expectedRel, rel, "unexpected relationship")
		})
	}
}

func TestPoint_Scale(t *testing.T) {
	tests := map[string]struct {
		point    Point   // Point to be scaled
		refPoint Point   // Reference point for scaling
		scale    float64 // Scaling factor
		expected Point   // Expected resulting point
	}{
		"float64: Scale point from reference by factor 1.5": {
			point:    New(2.0, 3.0),
			refPoint: New(1.0, 1.0),
			scale:    1.5,
			expected: New(2.5, 4.0),
		},
		"float64: Scale point from reference by factor 0.25": {
			point:    New(4.0, 8.0),
			refPoint: New(2.0, 2.0),
			scale:    0.25,
			expected: New(2.5, 3.5),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ref := tt.refPoint
			result := tt.point.Scale(ref, tt.scale)
			assert.InDelta(t, tt.expected.x, result.x, geom2d.GetEpsilon())
			assert.InDelta(t, tt.expected.y, result.y, geom2d.GetEpsilon())
		})
	}
}

func TestPoint_String(t *testing.T) {
	tests := map[string]struct {
		p        Point  // Supports different Point types with `any`
		expected string // Expected string representation of the point
	}{
		"float64: (1.2,3.4)": {
			p:        New(1.2, 3.4),
			expected: "(1.200000,3.400000)",
		},
		"float64: (-1.5,-2.5)": {
			p:        New(-1.5, -2.5),
			expected: "(-1.500000,-2.500000)",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tc.p.String()
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestPoint_Translate(t *testing.T) {
	tests := []struct {
		name     string
		p, q     Point
		expected Point
	}{
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
			actual := tc.p.Translate(tc.q)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestPoint_X(t *testing.T) {
	tests := []struct {
		name     string
		point    Point
		expected float64
	}{
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
			assert.Equal(t, tt.expected, tt.point.X())
		})
	}
}

func TestPoint_Y(t *testing.T) {
	tests := []struct {
		name     string
		point    Point
		expected float64
	}{
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
			assert.Equal(t, tt.expected, tt.point.Y())
		})
	}
}

func TestNewPointFromImagePoint(t *testing.T) {
	// Define test cases with various image.Point values
	tests := []struct {
		name     string
		imgPoint image.Point
		expected Point
	}{
		{
			name:     "Positive coordinates",
			imgPoint: image.Point{X: 10, Y: 20},
			expected: Point{10, 20},
		},
		{
			name:     "Negative coordinates",
			imgPoint: image.Point{X: -15, Y: -25},
			expected: Point{-15, -25},
		},
		{
			name:     "Zero coordinates",
			imgPoint: image.Point{X: 0, Y: 0},
			expected: Point{0, 0},
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
