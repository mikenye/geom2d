package geom2d

import (
	"image"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkConvexHull(b *testing.B) {
	points := []Point[int]{
		{1, 4}, {4, 13}, {8, 17}, {18, 20}, {33, 18}, {38, 11},
		{34, -2}, {21, -3}, {6, -1}, {7, 6}, {10, 14}, {5, 2},
		{16, 0}, {12, 12}, {23, 16}, {14, 6}, {24, 0}, {21, -2},
		{30, 14}, {27, 9}, {29, 4}, {31, -1}, {34, 7}, {35, 12},
	}

	// Reset the timer to ignore any setup cost
	b.ResetTimer()

	// Run the benchmark loop
	for i := 0; i < b.N; i++ {
		_ = ConvexHull(points)
	}
}

func TestPoint_Add(t *testing.T) {
	tests := []struct {
		name     string
		p, q     any // Use 'any' to support different Point types
		expected any
	}{
		// Integer points
		{
			name:     "int: (1,2)+(0,0)",
			p:        NewPoint(1, 2),
			q:        NewPoint(0, 0),
			expected: NewPoint(1, 2),
		},
		{
			name:     "int: (1,2)+(3,4)",
			p:        NewPoint(1, 2),
			q:        NewPoint(3, 4),
			expected: NewPoint(4, 6),
		},
		{
			name:     "int: (-1,-2)+(3,4)",
			p:        NewPoint(-1, -2),
			q:        NewPoint(3, 4),
			expected: NewPoint(2, 2),
		},

		// Float64 points
		{
			name:     "float64: (1.0,2.0)+(3.0,4.0)",
			p:        NewPoint(1.0, 2.0),
			q:        NewPoint(3.0, 4.0),
			expected: NewPoint(4.0, 6.0),
		},
		{
			name:     "float64: (-1.5,-2.5)+(3.5,4.5)",
			p:        NewPoint(-1.5, -2.5),
			q:        NewPoint(3.5, 4.5),
			expected: NewPoint(2.0, 2.0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch p := tt.p.(type) {
			case Point[int]:
				q := tt.q.(Point[int])
				expected := tt.expected.(Point[int])
				actual := p.Add(q)
				assert.Equal(t, expected, actual)

			case Point[float64]:
				q := tt.q.(Point[float64])
				expected := tt.expected.(Point[float64])
				actual := p.Add(q)
				assert.Equal(t, expected, actual)
			}
		})
	}
}

func TestPoint_AsFloat(t *testing.T) {
	// Define test cases with various point types
	tests := []struct {
		name     string
		point    Point[int]     // The point to convert
		expected Point[float64] // The expected result after conversion
	}{
		{
			name:     "Integer point conversion",
			point:    NewPoint(3, 4),
			expected: Point[float64]{3.0, 4.0},
		},
		{
			name:     "Negative integer point conversion",
			point:    NewPoint(-7, -5),
			expected: Point[float64]{-7.0, -5.0},
		},
		{
			name:     "Zero point conversion",
			point:    NewPoint(0, 0),
			expected: Point[float64]{0.0, 0.0},
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.point.AsFloat()
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

func TestPoint_CrossProduct(t *testing.T) {
	tests := []struct {
		name     string
		p, q     any // Support different Point types with `any`
		expected any // Expected result for different types
	}{
		// Integer points
		{
			name:     "int: (2,3) x (4,5)",
			p:        NewPoint(2, 3),
			q:        NewPoint(4, 5),
			expected: -2,
		},
		{
			name:     "int: (3,2) x (4,6)",
			p:        NewPoint(3, 2),
			q:        NewPoint(4, 6),
			expected: 10,
		},

		// Float64 points
		{
			name:     "float64: (2.0,3.0) x (4.0,5.0)",
			p:        NewPoint(2.0, 3.0),
			q:        NewPoint(4.0, 5.0),
			expected: -2.0,
		},
		{
			name:     "float64: (3.5,2.5) x (4.0,6.0)",
			p:        NewPoint(3.5, 2.5),
			q:        NewPoint(4.0, 6.0),
			expected: 11.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch p := tt.p.(type) {
			case Point[int]:
				q := tt.q.(Point[int])
				expected := tt.expected.(int)
				actual := p.CrossProduct(q)
				assert.Equal(t, expected, actual)

			case Point[float64]:
				q := tt.q.(Point[float64])
				expected := tt.expected.(float64)
				actual := p.CrossProduct(q)
				assert.Equal(t, expected, actual)
			}
		})
	}
}

func TestPoint_DistanceToLineSegment(t *testing.T) {
	tests := map[string]struct {
		point    any     // Point to be projected (can be int or float64)
		segment  any     // Line segment for projection (can be int or float64)
		expected float64 // Expected distance
	}{
		// Integer points test cases
		"Project onto segment from inside (int)": {
			point:    NewPoint(5, 5),
			segment:  NewLineSegment[int](NewPoint(2, 3), NewPoint(8, 7)),
			expected: 0.0, // Point is on the segment
		},
		"Project onto segment from above (int)": {
			point:    NewPoint(4, 6),
			segment:  NewLineSegment[int](NewPoint(2, 3), NewPoint(8, 7)),
			expected: 1.386,
		},
		"Project onto segment from below (int)": {
			point:    NewPoint(4, 2),
			segment:  NewLineSegment[int](NewPoint(2, 3), NewPoint(8, 7)),
			expected: 1.941,
		},
		"Project off the start of segment (int)": {
			point:    NewPoint(0, 5),
			segment:  NewLineSegment[int](NewPoint(2, 3), NewPoint(8, 7)),
			expected: 2.8284,
		},
		"Project off the end of segment (int)": {
			point:    NewPoint(10, 5),
			segment:  NewLineSegment[int](NewPoint(2, 3), NewPoint(8, 7)),
			expected: 2.8284,
		},

		// Float64 points test cases
		"Project onto segment from inside (float64)": {
			point:    NewPoint(5.5, 5.5),
			segment:  NewLineSegment[float64](NewPoint(2.0, 3.0), NewPoint(8.0, 7.0)),
			expected: 0.1387,
		},
		"Project onto segment from above (float64)": {
			point:    NewPoint(4.0, 6.0),
			segment:  NewLineSegment[float64](NewPoint(2.0, 3.0), NewPoint(8.0, 7.0)),
			expected: 1.386,
		},
		"Project onto segment from below (float64)": {
			point:    NewPoint(4.0, 2.0),
			segment:  NewLineSegment[float64](NewPoint(2.0, 3.0), NewPoint(8.0, 7.0)),
			expected: 1.941,
		},
		"Project off the start of segment (float64)": {
			point:    NewPoint(0.0, 5.0),
			segment:  NewLineSegment[float64](NewPoint(2.0, 3.0), NewPoint(8.0, 7.0)),
			expected: 2.8284,
		},
		"Project off the end of segment (float64)": {
			point:    NewPoint(10.0, 5.0),
			segment:  NewLineSegment[float64](NewPoint(2.0, 3.0), NewPoint(8.0, 7.0)),
			expected: 2.8284,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			switch p := tt.point.(type) {
			case Point[int]:
				seg := tt.segment.(LineSegment[int])
				actual := p.DistanceToLineSegment(seg)
				assert.InDelta(t, tt.expected, actual, 0.001, "Expected distance does not match actual distance for int points")

			case Point[float64]:
				seg := tt.segment.(LineSegment[float64])
				actual := p.DistanceToLineSegment(seg)
				assert.InDelta(t, tt.expected, actual, 0.001, "Expected distance does not match actual distance for float64 points")
			}
		})
	}
}

func TestPoint_DistanceToPoint(t *testing.T) {
	tests := []struct {
		name     string
		p, q     any     // Use `any` to handle different Point types
		expected float64 // Expected distance as a float64
	}{
		// Integer points
		{
			name:     "int: distance between (2,10) and (10,2)",
			p:        NewPoint(2, 10),
			q:        NewPoint(10, 2),
			expected: 11.3137,
		},
		{
			name:     "int: distance between (0,0) and (3,4)",
			p:        NewPoint(0, 0),
			q:        NewPoint(3, 4),
			expected: 5.0,
		},

		// Float64 points
		{
			name:     "float64: distance between (2.0,10.0) and (10.0,2.0)",
			p:        NewPoint(2.0, 10.0),
			q:        NewPoint(10.0, 2.0),
			expected: 11.3137,
		},
		{
			name:     "float64: distance between (0.0,0.0) and (3.0,4.0)",
			p:        NewPoint(0.0, 0.0),
			q:        NewPoint(3.0, 4.0),
			expected: 5.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch p := tt.p.(type) {
			case Point[int]:
				q := tt.q.(Point[int])
				actual := p.DistanceToPoint(q)
				assert.InDelta(t, tt.expected, actual, 0.0001)

			case Point[float64]:
				q := tt.q.(Point[float64])
				actual := p.DistanceToPoint(q)
				assert.InDelta(t, tt.expected, actual, 0.0001)
			}
		})
	}
}

func TestPoint_Div(t *testing.T) {
	tests := []struct {
		name     string
		p        any            // Use `any` to handle different Point types
		k        any            // The divisor, which can be int or float64
		expected Point[float64] // Expected result as a Point[float64] since division produces float results
	}{
		// Integer division cases
		{
			name:     "int: (2,3)/2",
			p:        NewPoint(2, 3),
			k:        2,
			expected: NewPoint(1.0, 1.5),
		},
		{
			name:     "int: (4,6)/2",
			p:        NewPoint(4, 6),
			k:        2,
			expected: NewPoint(2.0, 3.0),
		},

		// Float64 division cases
		{
			name:     "float64: (2.0,3.0)/2.0",
			p:        NewPoint(2.0, 3.0),
			k:        2.0,
			expected: NewPoint(1.0, 1.5),
		},
		{
			name:     "float64: (4.5,6.0)/1.5",
			p:        NewPoint(4.5, 6.0),
			k:        1.5,
			expected: NewPoint(3.0, 4.0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch p := tt.p.(type) {
			case Point[int]:
				k := tt.k.(int)
				actual := p.Div(k)
				assert.Equal(t, tt.expected, actual)

			case Point[float64]:
				k := tt.k.(float64)
				actual := p.Div(k)
				assert.Equal(t, tt.expected, actual)
			}
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
			p:        NewPoint(2, 3),
			q:        NewPoint(4, 5),
			expected: 23,
		},
		{
			name:     "int: (3,4) . (1,2)",
			p:        NewPoint(3, 4),
			q:        NewPoint(1, 2),
			expected: 11,
		},

		// Float64 points
		{
			name:     "float64: (2.0,3.0) . (4.0,5.0)",
			p:        NewPoint(2.0, 3.0),
			q:        NewPoint(4.0, 5.0),
			expected: 23.0,
		},
		{
			name:     "float64: (1.5,2.5) . (3.5,4.5)",
			p:        NewPoint(1.5, 2.5),
			q:        NewPoint(3.5, 4.5),
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
	tests := []struct {
		name     string
		p, q     any  // Supports different Point types with `any`
		expected bool // Expected result of equality comparison
	}{
		// Integer points
		{
			name:     "int: (2,3) == (4,5)",
			p:        NewPoint(2, 3),
			q:        NewPoint(4, 5),
			expected: false,
		},
		{
			name:     "int: (2,3) == (2,3)",
			p:        NewPoint(2, 3),
			q:        NewPoint(2, 3),
			expected: true,
		},

		// Float64 points
		{
			name:     "float64: (2.0,3.0) == (4.0,5.0)",
			p:        NewPoint(2.0, 3.0),
			q:        NewPoint(4.0, 5.0),
			expected: false,
		},
		{
			name:     "float64: (2.0,3.0) == (2.0,3.0)",
			p:        NewPoint(2.0, 3.0),
			q:        NewPoint(2.0, 3.0),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch p := tt.p.(type) {
			case Point[int]:
				q := tt.q.(Point[int])
				actual := p.Eq(q)
				assert.Equal(t, tt.expected, actual)

			case Point[float64]:
				q := tt.q.(Point[float64])
				actual := p.Eq(q)
				assert.Equal(t, tt.expected, actual)
			}
		})
	}
}

func TestPoint_IsOnLineSegment(t *testing.T) {
	tests := map[string]struct {
		point    any
		segment  any
		expected bool
	}{
		"Point on line segment (float64)": {
			point:    NewPoint[float64](1, 1),
			segment:  NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](2, 2)),
			expected: true,
		},
		"Point at endpoint A (float64)": {
			point:    NewPoint[float64](0, 0),
			segment:  NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](2, 2)),
			expected: true,
		},
		"Point at endpoint End (float64)": {
			point:    NewPoint[float64](2, 2),
			segment:  NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](2, 2)),
			expected: true,
		},
		"Point collinear but outside bounding box (float64)": {
			point:    NewPoint[float64](3, 3),
			segment:  NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](2, 2)),
			expected: false,
		},
		"Point not collinear (float64)": {
			point:    NewPoint[float64](1, 2),
			segment:  NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](2, 2)),
			expected: false,
		},
		"Point on line segment (int)": {
			point:    NewPoint[int](1, 1),
			segment:  NewLineSegment(NewPoint[int](0, 0), NewPoint[int](2, 2)),
			expected: true,
		},
		"Point at endpoint A (int)": {
			point:    NewPoint[int](0, 0),
			segment:  NewLineSegment(NewPoint[int](0, 0), NewPoint[int](2, 2)),
			expected: true,
		},
		"Point at endpoint End (int)": {
			point:    NewPoint[int](2, 2),
			segment:  NewLineSegment(NewPoint[int](0, 0), NewPoint[int](2, 2)),
			expected: true,
		},
		"Point collinear but outside bounding box (int)": {
			point:    NewPoint[int](3, 3),
			segment:  NewLineSegment(NewPoint[int](0, 0), NewPoint[int](2, 2)),
			expected: false,
		},
		"Point not collinear (int)": {
			point:    NewPoint[int](1, 2),
			segment:  NewLineSegment(NewPoint[int](0, 0), NewPoint[int](2, 2)),
			expected: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			switch point := tt.point.(type) {
			case Point[int]:
				segment := tt.segment.(LineSegment[int])
				result := point.IsOnLineSegment(segment)
				assert.Equal(t, tt.expected, result, "Test %s failed", name)
			case Point[float64]:
				segment := tt.segment.(LineSegment[float64])
				result := point.IsOnLineSegment(segment)
				assert.Equal(t, tt.expected, result, "Test %s failed", name)
			default:
				t.Errorf("Unsupported point type in test %s", name)
			}
		})
	}
}

func TestPoint_ProjectOntoLineSegment(t *testing.T) {
	tests := map[string]struct {
		point    any            // Point to be projected (can be int or float64)
		segment  any            // Line segment for projection (can be int or float64)
		expected Point[float64] // Expected projected point (float64 type)
	}{
		// Integer points test cases
		"int: Project onto segment from inside": {
			point:    NewPoint(5, 5),
			segment:  NewLineSegment[int](NewPoint(2, 3), NewPoint(8, 7)),
			expected: NewPoint[float64](5, 5), // Should project onto the line segment itself
		},
		"int: Project onto segment from above": {
			point:    NewPoint(4, 6),
			segment:  NewLineSegment[int](NewPoint(2, 3), NewPoint(8, 7)),
			expected: NewPoint[float64](4.769, 4.846),
		},
		"int: Project onto segment from below": {
			point:    NewPoint(4, 2),
			segment:  NewLineSegment[int](NewPoint(2, 3), NewPoint(8, 7)),
			expected: NewPoint[float64](2.9231, 3.6154), // Projected point should be (5,5)
		},
		"int: Project off the start of segment": {
			point:    NewPoint(0, 5),
			segment:  NewLineSegment[int](NewPoint(2, 3), NewPoint(8, 7)),
			expected: NewPoint[float64](2, 3), // Should return a point of the segment
		},
		"int: Project off the end of segment": {
			point:    NewPoint(10, 5),
			segment:  NewLineSegment[int](NewPoint(2, 3), NewPoint(8, 7)),
			expected: NewPoint[float64](8, 7), // Should return end point of the segment
		},

		// Float64 points test cases
		"float64: Project onto segment from inside": {
			point:    NewPoint(5.5, 5.5),
			segment:  NewLineSegment[float64](NewPoint(2.0, 3.0), NewPoint(8.0, 7.0)),
			expected: NewPoint[float64](5.5772, 5.3848), // Should project onto the line segment itself
		},
		"float64: Project off the start of segment": {
			point:    NewPoint(0.0, 5.0),
			segment:  NewLineSegment[float64](NewPoint(2.0, 3.0), NewPoint(8.0, 7.0)),
			expected: NewPoint[float64](2.0, 3.0), // Should return a point of the segment
		},
		"float64: Project off the end of segment": {
			point:    NewPoint(10.0, 5.0),
			segment:  NewLineSegment[float64](NewPoint(2.0, 3.0), NewPoint(8.0, 7.0)),
			expected: NewPoint[float64](8.0, 7.0), // Should return end point of the segment
		},

		// Zero-length segment test cases
		"int: Project onto zero-length segment": {
			point:    NewPoint(3, 4),
			segment:  NewLineSegment[int](NewPoint(2, 2), NewPoint(2, 2)), // Zero-length segment
			expected: NewPoint[float64](2, 2),                             // Should return point A (or End), since the segment is a single point
		},
		"float64: Project onto zero-length segment": {
			point:    NewPoint(5.0, 5.0),
			segment:  NewLineSegment[float64](NewPoint(2.5, 2.5), NewPoint(2.5, 2.5)), // Zero-length segment
			expected: NewPoint[float64](2.5, 2.5),                                     // Should return point A (or End), since the segment is a single point
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			switch p := tt.point.(type) {
			case Point[int]:
				seg := tt.segment.(LineSegment[int])
				actual := p.ProjectOntoLineSegment(seg)
				assert.InDelta(t, tt.expected.x, actual.x, 0.001)
				assert.InDelta(t, tt.expected.y, actual.y, 0.001)

			case Point[float64]:
				seg := tt.segment.(LineSegment[float64])
				actual := p.ProjectOntoLineSegment(seg)
				assert.InDelta(t, tt.expected.x, actual.x, 0.001)
				assert.InDelta(t, tt.expected.y, actual.y, 0.001)
			}
		})
	}
}

func TestPoint_Reflect(t *testing.T) {
	tests := map[string]struct {
		point    Point[float64]       // The point to reflect
		axis     ReflectionAxis       // The axis or line type for reflection
		line     LineSegment[float64] // Custom line segment for ReflectAcrossCustomLine reflection
		expected Point[float64]       // Expected reflected point
	}{
		"reflect across x-axis": {
			point:    NewPoint[float64](3, 4),
			axis:     ReflectAcrossXAxis,
			expected: NewPoint[float64](3, -4),
		},
		"reflect across y-axis": {
			point:    NewPoint[float64](3, 4),
			axis:     ReflectAcrossYAxis,
			expected: NewPoint[float64](-3, 4),
		},
		"reflect across y = x line (ReflectAcrossCustomLine)": {
			point:    NewPoint[float64](3, 4),
			axis:     ReflectAcrossCustomLine,
			line:     NewLineSegment[float64](NewPoint[float64](0, 0), NewPoint[float64](1, 1)),
			expected: NewPoint[float64](4, 3),
		},
		"reflect across y = -x line (ReflectAcrossCustomLine)": {
			point:    NewPoint[float64](3, 4),
			axis:     ReflectAcrossCustomLine,
			line:     NewLineSegment[float64](NewPoint[float64](0, 0), NewPoint[float64](-1, 1)),
			expected: NewPoint[float64](-4, -3),
		},
		"reflect across degenerate line segment": {
			point:    NewPoint[float64](3, 4),
			axis:     ReflectAcrossCustomLine,
			line:     NewLineSegment[float64](NewPoint[float64](1, 1), NewPoint[float64](1, 1)), // Degenerate line
			expected: NewPoint[float64](3, 4),                                                   // Expect the point to remain unchanged
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var result Point[float64]
			if tt.axis == ReflectAcrossCustomLine {
				result = tt.point.Reflect(ReflectAcrossCustomLine, tt.line)
			} else {
				result = tt.point.Reflect(tt.axis)
			}

			assert.InDelta(t, tt.expected.x, result.x, 0.001)
			assert.InDelta(t, tt.expected.y, result.y, 0.001)
		})
	}

	t.Run("ReflectAcrossCustomLine with no given line segment", func(t *testing.T) {
		point := NewPoint[float64](3, 4)
		result := point.Reflect(ReflectAcrossCustomLine)
		assert.Equal(t, point, result, 0.001)
	})

	t.Run("invalid ReflectionAxis", func(t *testing.T) {
		point := NewPoint[float64](3, 4)
		result := point.Reflect(-1)
		assert.Equal(t, point, result, 0.001)
	})
}

func TestPoint_Rotate(t *testing.T) {
	tests := map[string]struct {
		point    Point[float64] // The point to rotate
		origin   Point[float64] // The origin for RotateFrom (or (0,0) for Rotate)
		angle    float64        // The angle in radians
		expected Point[float64] // The expected result
	}{
		"rotate 90 degrees around origin": {
			point:    NewPoint[float64](1, 0),
			origin:   NewPoint[float64](0, 0),
			angle:    math.Pi / 2,
			expected: NewPoint[float64](0, 1),
		},
		"rotate 180 degrees around origin": {
			point:    NewPoint[float64](1, 1),
			origin:   NewPoint[float64](0, 0),
			angle:    math.Pi,
			expected: NewPoint[float64](-1, -1),
		},
		"rotate 90 degrees around (1,1)": {
			point:    NewPoint[float64](2, 1),
			origin:   NewPoint[float64](1, 1),
			angle:    math.Pi / 2,
			expected: NewPoint[float64](1, 2),
		},
		"rotate 45 degrees around (2,2)": {
			point:    NewPoint[float64](3, 2),
			origin:   NewPoint[float64](2, 2),
			angle:    math.Pi / 4,
			expected: NewPoint[float64](2.7071, 2.7071),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := tt.point.Rotate(tt.origin, tt.angle)
			assert.InDelta(t, tt.expected.x, result.x, 0.001)
			assert.InDelta(t, tt.expected.y, result.y, 0.001)
		})
	}
}

func TestPoint_Scale(t *testing.T) {
	tests := []struct {
		name     string
		p        any // Supports different Point types with `any`
		k        any // Scale factor, either int or float64
		expected any // Expected scaled result as either int or float64
	}{
		// Integer points
		{
			name:     "int: (2,3) * 2",
			p:        NewPoint(2, 3),
			k:        2,
			expected: NewPoint(4, 6),
		},
		{
			name:     "int: (3,4) * -1",
			p:        NewPoint(3, 4),
			k:        -1,
			expected: NewPoint(-3, -4),
		},

		// Float64 points
		{
			name:     "float64: (2.0,3.0) * 2.0",
			p:        NewPoint(2.0, 3.0),
			k:        2.0,
			expected: NewPoint(4.0, 6.0),
		},
		{
			name:     "float64: (1.5,2.5) * 0.5",
			p:        NewPoint(1.5, 2.5),
			k:        0.5,
			expected: NewPoint(0.75, 1.25),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch p := tt.p.(type) {
			case Point[int]:
				k := tt.k.(int)
				expected := tt.expected.(Point[int])
				actual := p.Scale(k)
				assert.Equal(t, expected, actual)

			case Point[float64]:
				k := tt.k.(float64)
				expected := tt.expected.(Point[float64])
				actual := p.Scale(k)
				assert.Equal(t, expected, actual)
			}
		})
	}
}

func TestPoint_ScaleFrom(t *testing.T) {
	tests := map[string]struct {
		point    any            // Point to be scaled (can be int or float64)
		refPoint any            // Reference point for scaling (can be int or float64)
		scale    float64        // Scaling factor
		expected Point[float64] // Expected resulting point (float64 type)
	}{
		// Integer test cases
		"int: Scale point from reference by factor 2": {
			point:    NewPoint(3, 4),
			refPoint: NewPoint(1, 1),
			scale:    2.0,
			expected: NewPoint[float64](5, 7),
		},
		"int: Scale point from reference by factor 0.5": {
			point:    NewPoint(3, 4),
			refPoint: NewPoint(1, 1),
			scale:    0.5,
			expected: NewPoint[float64](2, 2.5),
		},

		// Float64 test cases
		"float64: Scale point from reference by factor 1.5": {
			point:    NewPoint(2.0, 3.0),
			refPoint: NewPoint(1.0, 1.0),
			scale:    1.5,
			expected: NewPoint[float64](2.5, 4.0),
		},
		"float64: Scale point from reference by factor 0.25": {
			point:    NewPoint(4.0, 8.0),
			refPoint: NewPoint(2.0, 2.0),
			scale:    0.25,
			expected: NewPoint[float64](2.5, 3.5),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			switch point := tt.point.(type) {
			case Point[int]:
				ref := tt.refPoint.(Point[int])
				result := point.AsFloat().ScaleFrom(ref.AsFloat(), tt.scale)
				assert.InDelta(t, tt.expected.x, result.x, 0.001)
				assert.InDelta(t, tt.expected.y, result.y, 0.001)

			case Point[float64]:
				ref := tt.refPoint.(Point[float64])
				result := point.ScaleFrom(ref, tt.scale)
				assert.InDelta(t, tt.expected.x, result.x, 0.001)
				assert.InDelta(t, tt.expected.y, result.y, 0.001)
			}
		})
	}
}

func TestPoint_String(t *testing.T) {
	tests := []struct {
		name     string
		p        any    // Supports different Point types with `any`
		expected string // Expected string representation of the point
	}{
		// Integer points
		{
			name:     "int: (1,2)",
			p:        NewPoint(1, 2),
			expected: "Point[(1, 2)]",
		},
		{
			name:     "int: (0,-3)",
			p:        NewPoint(0, -3),
			expected: "Point[(0, -3)]",
		},

		// Float64 points
		{
			name:     "float64: (1.2,3.4)",
			p:        NewPoint(1.2, 3.4),
			expected: "Point[(1.2, 3.4)]",
		},
		{
			name:     "float64: (-1.5,-2.5)",
			p:        NewPoint(-1.5, -2.5),
			expected: "Point[(-1.5, -2.5)]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch p := tt.p.(type) {
			case Point[int]:
				actual := p.String()
				assert.Equal(t, tt.expected, actual)

			case Point[float64]:
				actual := p.String()
				assert.Equal(t, tt.expected, actual)
			}
		})
	}
}

func TestPoint_Sub(t *testing.T) {
	tests := []struct {
		name     string
		p, q     any // Supports different Point types with `any`
		expected any // Expected result after subtraction
	}{
		// Integer points
		{
			name:     "int: (1,2) - (3,4)",
			p:        NewPoint(1, 2),
			q:        NewPoint(3, 4),
			expected: NewPoint(-2, -2),
		},
		{
			name:     "int: (5,5) - (2,3)",
			p:        NewPoint(5, 5),
			q:        NewPoint(2, 3),
			expected: NewPoint(3, 2),
		},
		{
			name:     "int: (3,4) - (0,0)",
			p:        NewPoint(3, 4),
			q:        NewPoint(0, 0),
			expected: NewPoint(3, 4),
		},

		// Float64 points
		{
			name:     "float64: (1.0,2.0) - (3.0,4.0)",
			p:        NewPoint(1.0, 2.0),
			q:        NewPoint(3.0, 4.0),
			expected: NewPoint(-2.0, -2.0),
		},
		{
			name:     "float64: (5.5,5.5) - (2.0,3.0)",
			p:        NewPoint(5.5, 5.5),
			q:        NewPoint(2.0, 3.0),
			expected: NewPoint(3.5, 2.5),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch p := tt.p.(type) {
			case Point[int]:
				q := tt.q.(Point[int])
				expected := tt.expected.(Point[int])
				actual := p.Sub(q)
				assert.Equal(t, expected, actual)

			case Point[float64]:
				q := tt.q.(Point[float64])
				expected := tt.expected.(Point[float64])
				actual := p.Sub(q)
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
			point:    NewPoint(3, 4),
			expected: 3,
		},
		{
			name:     "int: Negative coordinates",
			point:    NewPoint(-7, -5),
			expected: -7,
		},
		{
			name:     "int: Zero x-coordinate",
			point:    NewPoint(0, 4),
			expected: 0,
		},

		// Float64 points
		{
			name:     "float64: Positive coordinates",
			point:    NewPoint(3.5, 4.5),
			expected: 3.5,
		},
		{
			name:     "float64: Negative coordinates",
			point:    NewPoint(-7.1, -5.2),
			expected: -7.1,
		},
		{
			name:     "float64: Zero x-coordinate",
			point:    NewPoint(0.0, 4.5),
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
			point:    NewPoint(3, 4),
			expected: 4,
		},
		{
			name:     "int: Negative coordinates",
			point:    NewPoint(-7, -5),
			expected: -5,
		},
		{
			name:     "int: Zero y-coordinate",
			point:    NewPoint(3, 0),
			expected: 0,
		},

		// Float64 points
		{
			name:     "float64: Positive coordinates",
			point:    NewPoint(3.5, 4.5),
			expected: 4.5,
		},
		{
			name:     "float64: Negative coordinates",
			point:    NewPoint(-7.1, -5.2),
			expected: -5.2,
		},
		{
			name:     "float64: Zero y-coordinate",
			point:    NewPoint(3.0, 0.0),
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

func TestConvexHull(t *testing.T) {
	tests := []struct {
		name     string
		points   any // Supports both int and float64 points
		expected any // Expected convex hull points
	}{
		// Integer points test cases
		{
			name: "int: Convex hull with multiple interior points",
			points: []Point[int]{
				{1, 4}, {4, 13}, {8, 17}, {18, 20}, {33, 18}, {38, 11}, {34, -2}, {21, -3}, {6, -1},
				{7, 6}, {10, 14}, {5, 2}, {16, 0}, {12, 12}, {23, 16}, {14, 6}, {24, 0}, {21, -2},
				{30, 14}, {27, 9}, {29, 4}, {31, -1}, {34, 7}, {35, 12}, {26, 2}, {25, 1}, {27, 3},
			},
			expected: []Point[int]{
				{21, -3}, {34, -2}, {38, 11}, {33, 18}, {18, 20}, {8, 17}, {4, 13}, {1, 4}, {6, -1},
			},
		},
		{
			name: "int: Simple square convex hull with all points outside",
			points: []Point[int]{
				{0, 0}, {20, 0}, {20, 20}, {0, 20},
				{19, 1}, {18, 3}, {17, 4}, {16, 4}, {15, 3}, {14, 2}, {13, 2}, {12, 2},
			},
			expected: []Point[int]{
				{0, 0}, {20, 0}, {20, 20}, {0, 20},
			},
		},

		// Float64 points test cases
		{
			name: "float64: Convex hull with multiple interior points",
			points: []Point[float64]{
				{1.0, 4.0}, {4.0, 13.0}, {8.0, 17.0}, {18.0, 20.0}, {33.0, 18.0}, {38.0, 11.0},
				{34.0, -2.0}, {21.0, -3.0}, {6.0, -1.0}, {7.0, 6.0}, {10.0, 14.0}, {5.0, 2.0},
				{16.0, 0.0}, {12.0, 12.0}, {23.0, 16.0}, {14.0, 6.0}, {24.0, 0.0}, {21.0, -2.0},
				{30.0, 14.0}, {27.0, 9.0}, {29.0, 4.0}, {31.0, -1.0}, {34.0, 7.0}, {35.0, 12.0},
				{26.0, 2.0}, {25.0, 1.0}, {27.0, 3.0},
			},
			expected: []Point[float64]{
				{21.0, -3.0}, {34.0, -2.0}, {38.0, 11.0}, {33.0, 18.0}, {18.0, 20.0}, {8.0, 17.0},
				{4.0, 13.0}, {1.0, 4.0}, {6.0, -1.0},
			},
		},
		{
			name: "float64: Simple square convex hull with all points outside",
			points: []Point[float64]{
				{0.0, 0.0}, {20.0, 0.0}, {20.0, 20.0}, {0.0, 20.0},
				{19.0, 1.0}, {18.0, 3.0}, {17.0, 4.0}, {16.0, 4.0}, {15.0, 3.0}, {14.0, 2.0},
				{13.0, 2.0}, {12.0, 2.0},
			},
			expected: []Point[float64]{
				{0.0, 0.0}, {20.0, 0.0}, {20.0, 20.0}, {0.0, 20.0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch points := tt.points.(type) {
			case []Point[int]:
				expected := tt.expected.([]Point[int])
				actual := ConvexHull(points)
				assert.Equal(t, expected, actual)

			case []Point[float64]:
				expected := tt.expected.([]Point[float64])
				actual := ConvexHull(points)
				assert.Equal(t, expected, actual)
			}
		})
	}
}

func TestEnsureClockwise(t *testing.T) {
	tests := map[string]struct {
		points   []Point[int]
		expected []Point[int]
	}{
		"already clockwise": {
			points:   []Point[int]{NewPoint(0, 0), NewPoint(2, 3), NewPoint(4, 0)},
			expected: []Point[int]{NewPoint(0, 0), NewPoint(2, 3), NewPoint(4, 0)},
		},
		"counterclockwise input": {
			points:   []Point[int]{NewPoint(0, 0), NewPoint(4, 0), NewPoint(2, 3)},
			expected: []Point[int]{NewPoint(2, 3), NewPoint(4, 0), NewPoint(0, 0)},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			EnsureClockwise(tt.points)
			assert.Equal(t, tt.expected, tt.points)
		})
	}
}

func TestEnsureCounterClockwise(t *testing.T) {
	tests := map[string]struct {
		points   []Point[int]
		expected []Point[int]
	}{
		"already counterclockwise": {
			points:   []Point[int]{NewPoint(0, 0), NewPoint(4, 0), NewPoint(2, 3)},
			expected: []Point[int]{NewPoint(0, 0), NewPoint(4, 0), NewPoint(2, 3)},
		},
		"clockwise input": {
			points:   []Point[int]{NewPoint(0, 0), NewPoint(2, 3), NewPoint(4, 0)},
			expected: []Point[int]{NewPoint(4, 0), NewPoint(2, 3), NewPoint(0, 0)},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			EnsureCounterClockwise(tt.points)
			assert.Equal(t, tt.expected, tt.points)
		})
	}
}

func TestFindLowestPoint(t *testing.T) {
	tests := []struct {
		name      string
		points    any // Supports both int and float64 points
		expectedI int // Expected index of the lowest point
		expectedP any // Expected lowest point as either int or float64
	}{
		// Integer points test case
		{
			name: "int: Multiple points, lowest at (0,0)",
			points: []Point[int]{
				{10, 10}, {10, 0}, {0, 0},
			},
			expectedI: 2,
			expectedP: Point[int]{0, 0},
		},

		// Float64 points test case
		{
			name: "float64: Multiple points, lowest at (0.0,0.0)",
			points: []Point[float64]{
				{10.5, 10.5}, {10.0, 0.0}, {0.0, 0.0},
			},
			expectedI: 2,
			expectedP: Point[float64]{0.0, 0.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch points := tt.points.(type) {
			case []Point[int]:
				expectedP := tt.expectedP.(Point[int])
				actualI, actualP := findLowestLeftestPoint(points...)
				assert.Equal(t, tt.expectedI, actualI)
				assert.Equal(t, expectedP, actualP)

			case []Point[float64]:
				expectedP := tt.expectedP.(Point[float64])
				actualI, actualP := findLowestLeftestPoint(points...)
				assert.Equal(t, tt.expectedI, actualI)
				assert.Equal(t, expectedP, actualP)
			}
		})
	}
}

func TestTriangleAreaX2Signed(t *testing.T) {
	tests := []struct {
		name       string
		p0, p1, p2 any // Supports both int and float64 points
		expected   any // Expected result as either int or float64
	}{
		// Integer points test cases
		{
			name:     "int: (0,0), (10,10), (10,0)",
			p0:       NewPoint(0, 0),
			p1:       NewPoint(10, 10),
			p2:       NewPoint(10, 0),
			expected: -100,
		},
		{
			name:     "int: (0,0), (10,0), (10,10)",
			p0:       NewPoint(0, 0),
			p1:       NewPoint(10, 0),
			p2:       NewPoint(10, 10),
			expected: 100,
		},

		// Float64 points test cases
		{
			name:     "float64: (0.0,0.0), (10.0,10.0), (10.0,0.0)",
			p0:       NewPoint(0.0, 0.0),
			p1:       NewPoint(10.0, 10.0),
			p2:       NewPoint(10.0, 0.0),
			expected: -100.0,
		},
		{
			name:     "float64: (0.0,0.0), (10.0,0.0), (10.0,10.0)",
			p0:       NewPoint(0.0, 0.0),
			p1:       NewPoint(10.0, 0.0),
			p2:       NewPoint(10.0, 10.0),
			expected: 100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch p0 := tt.p0.(type) {
			case Point[int]:
				p1 := tt.p1.(Point[int])
				p2 := tt.p2.(Point[int])
				expected := tt.expected.(int)
				actual := triangleAreaX2Signed(p0, p1, p2)
				assert.Equal(t, expected, actual)

			case Point[float64]:
				p1 := tt.p1.(Point[float64])
				p2 := tt.p2.(Point[float64])
				expected := tt.expected.(float64)
				actual := triangleAreaX2Signed(p0, p1, p2)
				assert.Equal(t, expected, actual)
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
			result := NewPointFromImagePoint(tt.imgPoint)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestOrderPointsByAngleAboutLowestPoint(t *testing.T) {
	tests := map[string]struct {
		points, expected any
	}{
		"int: non-collinear points": {
			points: []Point[int]{ // input points are expected points, but have been randomized.
				{-10, 10}, {-8, 10}, {10, 9}, {-10, 0}, {-2, 10}, {10, 4},
				{-9, 10}, {-10, 6}, {-10, 3}, {-10, 1}, {-6, 10}, {-10, 7},
				{-4, 10}, {-10, 4}, {6, 10}, {10, 5}, {10, 3}, {-5, 10},
				{10, 1}, {-7, 10}, {3, 10}, {-10, 10}, {8, 10}, {-10, 8},
				{-3, 10}, {7, 10}, {10, 2}, {-10, 9}, {0, 10}, {9, 10},
				{10, 10}, {2, 10}, {0, -1}, {-10, 5}, {4, 10}, {-10, 2},
				{10, 0}, {5, 10}, {10, 6}, {10, 8}, {10, 7}, {-1, 10},
				{1, 10},
			},
			expected: []Point[int]{
				{0, -1}, // lowest point
				{10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {10, 5}, {10, 6},
				{10, 7}, {10, 8}, {10, 9}, {10, 10}, {9, 10}, {8, 10}, {7, 10},
				{6, 10}, {5, 10}, {4, 10}, {3, 10}, {2, 10}, {1, 10}, {0, 10},
				{-1, 10}, {-2, 10}, {-3, 10}, {-4, 10}, {-5, 10}, {-6, 10}, {-7, 10},
				{-8, 10}, {-9, 10}, {-10, 10}, {-10, 10}, {-10, 9}, {-10, 8}, {-10, 7},
				{-10, 6}, {-10, 5}, {-10, 4}, {-10, 3}, {-10, 2}, {-10, 1}, {-10, 0},
			},
		},
		"int: collinear points": {
			points: []Point[int]{ // input points are expected points, but have been randomized.
				{1, 8}, {-3, 8}, {-7, 7}, {3, 3}, {-6, 8}, {-3, 3}, {7, 8},
				{7, 7}, {5, 8}, {4, 8}, {0, 0}, {3, 8}, {4, 4}, {0, 8},
				{8, 8}, {-4, 4}, {-5, 5}, {-5, 8}, {6, 6}, {5, 5}, {6, 8},
				{1, 1}, {-2, 8}, {2, 8}, {-7, 8}, {2, 2}, {-1, 1}, {-2, 2},
				{-1, 8}, {-8, 8}, {-6, 6}, {-4, 8}},

			expected: []Point[int]{
				{0, 0}, // lowest point
				{1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5}, {6, 6}, {7, 7}, {8, 8},
				{7, 8}, {6, 8}, {5, 8}, {4, 8}, {3, 8}, {2, 8}, {1, 8},
				{0, 8}, {-1, 8}, {-2, 8}, {-3, 8}, {-4, 8}, {-5, 8},
				{-6, 8}, {-7, 8}, {-1, 1}, {-2, 2}, {-3, 3}, {-4, 4},
				{-5, 5}, {-6, 6}, {-7, 7}, {-8, 8},
			},
		},
		"float64: non-collinear points": {
			points: []Point[float64]{ // input points are expected points, but have been randomized.
				{-10, 10}, {-8, 10}, {10, 9}, {-10, 0}, {-2, 10}, {10, 4},
				{-9, 10}, {-10, 6}, {-10, 3}, {-10, 1}, {-6, 10}, {-10, 7},
				{-4, 10}, {-10, 4}, {6, 10}, {10, 5}, {10, 3}, {-5, 10},
				{10, 1}, {-7, 10}, {3, 10}, {-10, 10}, {8, 10}, {-10, 8},
				{-3, 10}, {7, 10}, {10, 2}, {-10, 9}, {0, 10}, {9, 10},
				{10, 10}, {2, 10}, {0, -1}, {-10, 5}, {4, 10}, {-10, 2},
				{10, 0}, {5, 10}, {10, 6}, {10, 8}, {10, 7}, {-1, 10},
				{1, 10},
			},
			expected: []Point[float64]{
				{0, -1}, // lowest point
				{10, 0}, {10, 1}, {10, 2}, {10, 3}, {10, 4}, {10, 5}, {10, 6},
				{10, 7}, {10, 8}, {10, 9}, {10, 10}, {9, 10}, {8, 10}, {7, 10},
				{6, 10}, {5, 10}, {4, 10}, {3, 10}, {2, 10}, {1, 10}, {0, 10},
				{-1, 10}, {-2, 10}, {-3, 10}, {-4, 10}, {-5, 10}, {-6, 10}, {-7, 10},
				{-8, 10}, {-9, 10}, {-10, 10}, {-10, 10}, {-10, 9}, {-10, 8}, {-10, 7},
				{-10, 6}, {-10, 5}, {-10, 4}, {-10, 3}, {-10, 2}, {-10, 1}, {-10, 0},
			},
		},
		"float64: collinear points": {
			points: []Point[float64]{ // input points are expected points, but have been randomized.
				{1, 8}, {-3, 8}, {-7, 7}, {3, 3}, {-6, 8}, {-3, 3}, {7, 8},
				{7, 7}, {5, 8}, {4, 8}, {0, 0}, {3, 8}, {4, 4}, {0, 8},
				{8, 8}, {-4, 4}, {-5, 5}, {-5, 8}, {6, 6}, {5, 5}, {6, 8},
				{1, 1}, {-2, 8}, {2, 8}, {-7, 8}, {2, 2}, {-1, 1}, {-2, 2},
				{-1, 8}, {-8, 8}, {-6, 6}, {-4, 8}},

			expected: []Point[float64]{
				{0, 0}, // lowest point
				{1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5}, {6, 6}, {7, 7}, {8, 8},
				{7, 8}, {6, 8}, {5, 8}, {4, 8}, {3, 8}, {2, 8}, {1, 8},
				{0, 8}, {-1, 8}, {-2, 8}, {-3, 8}, {-4, 8}, {-5, 8},
				{-6, 8}, {-7, 8}, {-1, 1}, {-2, 2}, {-3, 3}, {-4, 4},
				{-5, 5}, {-6, 6}, {-7, 7}, {-8, 8},
			},
		},
	}
	for name, tt := range tests {

		t.Run(name, func(t *testing.T) {
			switch points := tt.points.(type) {
			case []Point[int]:
				expected := tt.expected.([]Point[int])
				orderPointsByAngleAboutLowestPoint(expected[0], points)
				assert.Equal(t, tt.expected, tt.points)

			case []Point[float64]:
				expected := tt.expected.([]Point[float64])
				orderPointsByAngleAboutLowestPoint(expected[0], points)
				assert.Equal(t, tt.expected, tt.points)
			}
		})
	}
}

func TestOrientation(t *testing.T) {
	tests := map[string]struct {
		p0, p1, p2 any
		expected   PointOrientation
	}{
		"int: (0,0), (10,10), (10,0)": {
			p0:       NewPoint[int](0, 0),
			p1:       NewPoint[int](10, 10),
			p2:       NewPoint[int](10, 0),
			expected: PointsClockwise,
		},
		"int: (0,0), (10,0), (10,10)": {
			p0:       NewPoint[int](0, 0),
			p1:       NewPoint[int](10, 0),
			p2:       NewPoint[int](10, 10),
			expected: PointsCounterClockwise,
		},
		"int: (0,0), (10,0), (20,0)": {
			p0:       NewPoint[int](0, 0),
			p1:       NewPoint[int](10, 0),
			p2:       NewPoint[int](20, 0),
			expected: PointsCollinear,
		},
		"float64: (0,0), (10,10), (10,0)": {
			p0:       NewPoint[float64](0, 0),
			p1:       NewPoint[float64](10, 10),
			p2:       NewPoint[float64](10, 0),
			expected: PointsClockwise,
		},
		"float64: (0,0), (10,0), (10,10)": {
			p0:       NewPoint[float64](0, 0),
			p1:       NewPoint[float64](10, 0),
			p2:       NewPoint[float64](10, 10),
			expected: PointsCounterClockwise,
		},
		"float64: (0,0), (10,0), (20,0)": {
			p0:       NewPoint[float64](0, 0),
			p1:       NewPoint[float64](10, 0),
			p2:       NewPoint[float64](20, 0),
			expected: PointsCollinear,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			switch p0 := tt.p0.(type) {
			case Point[int]:
				p1 := tt.p1.(Point[int])
				p2 := tt.p2.(Point[int])
				actual := Orientation(p0, p1, p2)
				assert.Equal(t, tt.expected, actual)

			case Point[float64]:
				p1 := tt.p1.(Point[float64])
				p2 := tt.p2.(Point[float64])
				actual := Orientation(p0, p1, p2)
				assert.Equal(t, tt.expected, actual)
			}
		})
	}
}

func TestRelativeAngle(t *testing.T) {
	tests := map[string]struct {
		A, B, O  any
		expected float64
	}{
		"int: ~22 deg": {
			A:        NewPoint[int](10, 0),
			B:        NewPoint[int](10, 4),
			O:        NewPoint[int](0, 0),
			expected: 0.38050637717509, // ~22 deg
		},
		"int: ~66 deg": {
			A:        NewPoint[int](10, 0),
			B:        NewPoint[int](4, 9),
			O:        NewPoint[int](0, 0),
			expected: 1.1525719972927, // ~66 deg
		},
		"int: ~90 deg": {
			A:        NewPoint[int](10, 0),
			B:        NewPoint[int](0, 10),
			O:        NewPoint[int](0, 0),
			expected: 1.5708, // ~90 deg
		},
		"int: ~114 deg": {
			A:        NewPoint[int](10, 0),
			B:        NewPoint[int](-4, 9),
			O:        NewPoint[int](0, 0),
			expected: 1.98902065628929, // ~114 deg
		},
		"int: ~156 deg": {
			A:        NewPoint[int](10, 0),
			B:        NewPoint[int](-9, 4),
			O:        NewPoint[int](0, 0),
			expected: 2.72336832408371, // ~156 deg
		},
		"float64: ~22 deg": {
			A:        NewPoint[float64](10, 0),
			B:        NewPoint[float64](10, 4),
			O:        NewPoint[float64](0, 0),
			expected: 0.38050637717509, // ~22 deg
		},
		"float64: ~66 deg": {
			A:        NewPoint[float64](10, 0),
			B:        NewPoint[float64](4, 9),
			O:        NewPoint[float64](0, 0),
			expected: 1.1525719972927, // ~66 deg
		},
		"float64: ~90 deg": {
			A:        NewPoint[float64](10, 0),
			B:        NewPoint[float64](0, 10),
			O:        NewPoint[float64](0, 0),
			expected: 1.5708, // ~90 deg
		},
		"float64: ~114 deg": {
			A:        NewPoint[float64](10, 0),
			B:        NewPoint[float64](-4, 9),
			O:        NewPoint[float64](0, 0),
			expected: 1.98902065628929, // ~114 deg
		},
		"float64: ~156 deg": {
			A:        NewPoint[float64](10, 0),
			B:        NewPoint[float64](-9, 4),
			O:        NewPoint[float64](0, 0),
			expected: 2.72336832408371, // ~156 deg
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			switch A := tt.A.(type) {
			case Point[int]:
				B := tt.B.(Point[int])
				O := tt.O.(Point[int])
				actual := RelativeAngle(A, B, O)
				assert.InDelta(t, tt.expected, actual, 0.00001)

			case Point[float64]:
				B := tt.B.(Point[float64])
				O := tt.O.(Point[float64])
				actual := RelativeAngle(A, B, O)
				assert.InDelta(t, tt.expected, actual, 0.00001)
			}
		})
	}
}

func TestRelativeAngle_NoOriginGiven(t *testing.T) {
	tests := map[string]struct {
		A, B     any
		expected float64
	}{
		"int: right triangle, origin at 0,0": {
			A:        NewPoint[int](10, 10),
			B:        NewPoint[int](10, 0),
			expected: 0.785398, // 45 deg
		},
		"float64: right triangle, origin at 0,0": {
			A:        NewPoint[float64](10, 10),
			B:        NewPoint[float64](10, 0),
			expected: 0.785398, // 45 deg
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			switch A := tt.A.(type) {
			case Point[int]:
				B := tt.B.(Point[int])
				actual := RelativeAngle(A, B)
				assert.InDelta(t, tt.expected, actual, 0.00001)

			case Point[float64]:
				B := tt.B.(Point[float64])
				actual := RelativeAngle(A, B)
				assert.InDelta(t, tt.expected, actual, 0.00001)
			}
		})
	}
}

func TestSignedArea2X(t *testing.T) {
	tests := map[string]struct {
		points   []Point[int]
		expected int
	}{
		"counterclockwise triangle": {
			points:   []Point[int]{NewPoint(0, 0), NewPoint(4, 0), NewPoint(2, 3)},
			expected: 12,
		},
		"clockwise triangle": {
			points:   []Point[int]{NewPoint(0, 0), NewPoint(2, 3), NewPoint(4, 0)},
			expected: -12,
		},
		"collinear points": {
			points:   []Point[int]{NewPoint(0, 0), NewPoint(2, 2), NewPoint(4, 4)},
			expected: 0,
		},
		"two points, not a polygon, no area": {
			points:   []Point[int]{NewPoint(0, 0), NewPoint(2, 2)},
			expected: 0,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := SignedArea2X(tt.points)
			assert.Equal(t, tt.expected, result)
		})
	}
}
