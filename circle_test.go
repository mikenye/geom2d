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
			result := tt.circle.AsFloat()
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
	tests := map[string]struct {
		c1       Circle[float64]
		c2       Circle[float64]
		epsilon  float64
		expected CircleCircleRelationship
	}{
		"Disjoint Circles": {
			c1:       NewCircle[float64](NewPoint[float64](0, 0), 5),
			c2:       NewCircle[float64](NewPoint[float64](15, 0), 5),
			epsilon:  1e-10,
			expected: CCRMiss,
		},
		"Externally Tangent Circles": {
			c1:       NewCircle[float64](NewPoint[float64](0, 0), 5),
			c2:       NewCircle[float64](NewPoint[float64](10, 0), 5),
			epsilon:  1e-10,
			expected: CCRTouchingExternal,
		},
		"Overlapping Circles": {
			c1:       NewCircle[float64](NewPoint[float64](0, 0), 5),
			c2:       NewCircle[float64](NewPoint[float64](6, 0), 5),
			epsilon:  1e-10,
			expected: CCROverlapping,
		},
		"Internally Tangent Circles": {
			c1:       NewCircle[float64](NewPoint[float64](0, 0), 5),
			c2:       NewCircle[float64](NewPoint[float64](2, 0), 3),
			epsilon:  1e-10,
			expected: CCRTouchingInternal,
		},
		"Circle Fully Contained in Another": {
			c1:       NewCircle[float64](NewPoint[float64](0, 0), 10),
			c2:       NewCircle[float64](NewPoint[float64](2, 2), 5),
			epsilon:  1e-10,
			expected: CCRContained,
		},
		"Circle Fully Contains Another": {
			c1:       NewCircle[float64](NewPoint[float64](2, 2), 5),
			c2:       NewCircle[float64](NewPoint[float64](0, 0), 10),
			epsilon:  1e-10,
			expected: CCRContained,
		},
		"Equal Circles": {
			c1:       NewCircle[float64](NewPoint[float64](0, 0), 5),
			c2:       NewCircle[float64](NewPoint[float64](0, 0), 5),
			epsilon:  1e-10,
			expected: CCREqual,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.c1.RelationshipToCircle(tc.c2, WithEpsilon(tc.epsilon))
			assert.Equal(t, tc.expected, result, "Expected %v, got %v", tc.expected, result)
		})
	}
}

func TestCircle_RelationshipToLineSegment(t *testing.T) {
	tests := map[string]struct {
		segment  LineSegment[float64]
		circle   Circle[float64]
		epsilon  float64
		expected CircleLineSegmentRelationship
	}{
		"segment completely inside circle": {
			segment:  NewLineSegment(NewPoint[float64](1, 1), NewPoint[float64](2, 2)),
			circle:   Circle[float64]{center: NewPoint[float64](0, 0), radius: 5},
			epsilon:  0,
			expected: CLRInside,
		},
		"segment completely outside circle": {
			segment:  NewLineSegment(NewPoint[float64](10, 10), NewPoint[float64](15, 15)),
			circle:   Circle[float64]{center: NewPoint[float64](0, 0), radius: 5},
			epsilon:  0,
			expected: CLROutside,
		},
		"segment intersects circle at two points": {
			segment:  NewLineSegment(NewPoint[float64](-6, 0), NewPoint[float64](6, 0)),
			circle:   Circle[float64]{center: NewPoint[float64](0, 0), radius: 5},
			epsilon:  0,
			expected: CLRIntersecting,
		},
		"segment is tangent to circle": {
			segment:  NewLineSegment(NewPoint[float64](5, -5), NewPoint[float64](5, 5)),
			circle:   Circle[float64]{center: NewPoint[float64](0, 0), radius: 5},
			epsilon:  1e-10,
			expected: CLRTangent,
		},
		"segment partially inside circle": {
			segment:  NewLineSegment(NewPoint[float64](1, 1), NewPoint[float64](10, 10)),
			circle:   Circle[float64]{center: NewPoint[float64](0, 0), radius: 5},
			epsilon:  0,
			expected: CLRIntersecting,
		},
		"segment with one endpoint on circumference and other outside": {
			segment:  NewLineSegment(NewPoint[float64](5, 0), NewPoint[float64](10, 10)),
			circle:   Circle[float64]{center: NewPoint[float64](0, 0), radius: 5},
			epsilon:  1e-10,
			expected: CLROneEndOnCircumferenceOutside,
		},
		"segment with one endpoint on circumference and other inside": {
			segment:  NewLineSegment(NewPoint[float64](5, 0), NewPoint[float64](2, 2)),
			circle:   Circle[float64]{center: NewPoint[float64](0, 0), radius: 5},
			epsilon:  1e-10,
			expected: CLROneEndOnCircumferenceInside,
		},
		"segment with both endpoints on circumference": {
			segment:  NewLineSegment(NewPoint[float64](5, 0), NewPoint[float64](-5, 0)),
			circle:   Circle[float64]{center: NewPoint[float64](0, 0), radius: 5},
			epsilon:  1e-10,
			expected: CLRBothEndsOnCircumference,
		},
		"degenerate segment inside circle": {
			segment:  NewLineSegment(NewPoint[float64](1, 1), NewPoint[float64](1, 1)),
			circle:   Circle[float64]{center: NewPoint[float64](0, 0), radius: 5},
			epsilon:  0,
			expected: CLRInside,
		},
		"degenerate segment on circle boundary": {
			segment:  NewLineSegment(NewPoint[float64](5, 0), NewPoint[float64](5, 0)),
			circle:   Circle[float64]{center: NewPoint[float64](0, 0), radius: 5},
			epsilon:  1e-10,
			expected: CLRBothEndsOnCircumference,
		},
		"degenerate segment outside circle": {
			segment:  NewLineSegment(NewPoint[float64](10, 10), NewPoint[float64](10, 10)),
			circle:   Circle[float64]{center: NewPoint[float64](0, 0), radius: 5},
			epsilon:  0,
			expected: CLROutside,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.circle.RelationshipToLineSegment(tc.segment, WithEpsilon(tc.epsilon))
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestCircle_RelationshipToPoint(t *testing.T) {
	tests := map[string]struct {
		circle   Circle[float64]
		point    Point[float64]
		epsilon  float64
		expected PointCircleRelationship
	}{
		"point inside circle": {
			circle:   NewCircle(NewPoint[float64](0.0, 0.0), 5.0),
			point:    NewPoint[float64](-3.0, -2.0),
			epsilon:  0,
			expected: PCRInside,
		},
		"point on circle boundary": {
			circle:   NewCircle(NewPoint[float64](0.0, 0.0), 5.0),
			point:    NewPoint[float64](3.0, 4.0),
			epsilon:  1e-10,
			expected: PCROnCircumference,
		},
		"point outside circle": {
			circle:   NewCircle(NewPoint[float64](0.0, 0.0), 5.0),
			point:    NewPoint[float64](6.0, 8.0),
			epsilon:  0,
			expected: PCROutside,
		},
		"point at center of circle": {
			circle:   NewCircle(NewPoint[float64](0.0, 0.0), 5.0),
			point:    NewPoint[float64](0.0, 0.0),
			epsilon:  0,
			expected: PCRInside,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.circle.RelationshipToPoint(tc.point, WithEpsilon(tc.epsilon))
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestCircle_RelationshipToRectangle(t *testing.T) {
	tests := map[string]struct {
		circle      Circle[int]
		rectangle   Rectangle[int]
		epsilon     float64
		expectedRel CircleRectangleRelationship
	}{
		"Disjoint": {
			circle:      NewCircle(NewPoint(50, 50), 10),
			rectangle:   NewRectangle([]Point[int]{NewPoint(0, 0), NewPoint(10, 0), NewPoint(10, 10), NewPoint(0, 10)}),
			epsilon:     1e-10,
			expectedRel: CRRMiss,
		},
		"Circle Inside Rectangle": {
			circle:      NewCircle(NewPoint(5, 5), 4),
			rectangle:   NewRectangle([]Point[int]{NewPoint(0, 0), NewPoint(10, 0), NewPoint(10, 10), NewPoint(0, 10)}),
			epsilon:     1e-10,
			expectedRel: CRRCircleInRect,
		},
		"Rectangle Inside Circle": {
			circle:      NewCircle(NewPoint(5, 5), 10),
			rectangle:   NewRectangle([]Point[int]{NewPoint(3, 3), NewPoint(7, 3), NewPoint(7, 7), NewPoint(3, 7)}),
			epsilon:     1e-10,
			expectedRel: CRRRectInCircle,
		},
		"Circle exactly in Rect (Circle Touching Edges)": {
			circle:      NewCircle(NewPoint(5, 5), 5),
			rectangle:   NewRectangle([]Point[int]{NewPoint(0, 0), NewPoint(10, 0), NewPoint(10, 10), NewPoint(0, 10)}),
			epsilon:     1e-10,
			expectedRel: CRRCircleInRect,
		},
		"Intersecting (Circle Center On Edge)": {
			circle:      NewCircle(NewPoint(5, 0), 5),
			rectangle:   NewRectangle([]Point[int]{NewPoint(0, 0), NewPoint(10, 0), NewPoint(10, 10), NewPoint(0, 10)}),
			epsilon:     1e-10,
			expectedRel: CRRIntersection,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actualRel := test.circle.RelationshipToRectangle(test.rectangle, WithEpsilon(test.epsilon))
			require.Equal(t, test.expectedRel, actualRel, "Relationship mismatch for test: %s", name)
		})
	}
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
