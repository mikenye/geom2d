package rectangle

import (
	"encoding/json"
	"github.com/mikenye/geom2d/linesegment"
	"github.com/mikenye/geom2d/point"
	"github.com/mikenye/geom2d/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"image"
	"testing"
)

func TestNewFromImageRect(t *testing.T) {
	tests := map[string]struct {
		imageRect image.Rectangle
		expected  Rectangle
	}{
		"simple rectangle": {
			imageRect: image.Rect(0, 0, 10, 20),
			expected: NewFromPoints(
				point.New(0, 0),
				point.New(10, 20),
				point.New(0, 20),
				point.New(10, 0),
			),
		},
		"negative coordinates": {
			imageRect: image.Rect(-5, -10, 5, 10),
			expected: NewFromPoints(
				point.New(-5, -10),
				point.New(5, 10),
				point.New(-5, 10),
				point.New(5, -10),
			),
		},
		"zero size rectangle": {
			imageRect: image.Rect(0, 0, 0, 0),
			expected: NewFromPoints(
				point.New(0, 0),
				point.New(0, 0),
				point.New(0, 0),
				point.New(0, 0),
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, NewFromImageRect(tc.imageRect))
		})
	}
}

func TestRectangle_Area(t *testing.T) {
	tests := map[string]struct {
		rect     Rectangle
		expected float64
	}{
		"standard rectangle": {
			rect:     New(0, 0, 10, 20),
			expected: 200,
		},
		"rectangle with swapped corners": {
			rect:     New(10, 20, 0, 0),
			expected: 200,
		},
		"degenerate rectangle (zero width)": {
			rect:     New(5, 5, 5, 15),
			expected: 0,
		},
		"degenerate rectangle (zero height)": {
			rect:     New(5, 5, 15, 5),
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

func TestRectangle_ContainsPoint(t *testing.T) {
	tests := map[string]struct {
		rect     Rectangle
		point    point.Point
		expected bool
	}{
		"point inside rectangle": {
			rect:     New(0, 0, 10, 20),
			point:    point.New(5, 10),
			expected: true,
		},
		"point on rectangle boundary (top-left corner)": {
			rect:     New(0, 0, 10, 20),
			point:    point.New(0, 0),
			expected: true,
		},
		"point on rectangle boundary (bottom-right corner)": {
			rect:     New(0, 0, 10, 20),
			point:    point.New(10, 20),
			expected: true,
		},
		"point outside rectangle": {
			rect:     New(0, 0, 10, 20),
			point:    point.New(15, 10),
			expected: false,
		},
		"point on boundary (horizontal edge)": {
			rect:     New(0, 0, 10, 20),
			point:    point.New(5, 0),
			expected: true,
		},
		"point on boundary (vertical edge)": {
			rect:     New(0, 0, 10, 20),
			point:    point.New(10, 5),
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

func TestRectangle_Contour(t *testing.T) {
	// Define a rectangle with specific corners
	bottomLeft := point.New(0, 0)
	bottomRight := point.New(4, 0)
	topRight := point.New(4, 3)
	topLeft := point.New(0, 3)

	rect := NewFromPoints(bottomLeft, bottomRight, topRight, topLeft)

	// Call the Contour method
	bl, br, tr, tl := rect.Contour()

	// Verify the output matches the expected corners
	assert.Equal(t, bottomLeft, bl, "bottom-left corner mismatch")
	assert.Equal(t, bottomRight, br, "bottom-right corner mismatch")
	assert.Equal(t, topRight, tr, "top-right corner mismatch")
	assert.Equal(t, topLeft, tl, "top-left corner mismatch")
}

func TestRectangle_Edges(t *testing.T) {
	// Define a rectangle with specific corners
	bottomLeft := point.New(0, 0)
	bottomRight := point.New(4, 0)
	topRight := point.New(4, 3)
	topLeft := point.New(0, 3)

	rect := NewFromPoints(bottomLeft, bottomRight, topRight, topLeft)

	// Call the Edges method
	bottom, right, top, left := rect.Edges()

	// Define the expected edges
	expectedBottom := linesegment.NewFromPoints(bottomLeft, bottomRight)
	expectedRight := linesegment.NewFromPoints(bottomRight, topRight)
	expectedTop := linesegment.NewFromPoints(topRight, topLeft)
	expectedLeft := linesegment.NewFromPoints(topLeft, bottomLeft)

	// Verify the edges match the expected line segments
	assert.Equal(t, expectedBottom, bottom, "bottom edge mismatch")
	assert.Equal(t, expectedRight, right, "right edge mismatch")
	assert.Equal(t, expectedTop, top, "top edge mismatch")
	assert.Equal(t, expectedLeft, left, "left edge mismatch")
}

func TestRectangle_Eq(t *testing.T) {
	tests := map[string]struct {
		rect1       Rectangle
		rect2       Rectangle
		expectEqual bool
	}{
		"equal rectangles": {
			rect1: NewFromPoints(
				point.New(0, 0),
				point.New(4, 0),
				point.New(4, 3),
				point.New(0, 3),
			),
			rect2: NewFromPoints(
				point.New(0, 0),
				point.New(4, 0),
				point.New(4, 3),
				point.New(0, 3),
			),
			expectEqual: true,
		},
		"different rectangles": {
			rect1: NewFromPoints(
				point.New(0, 0),
				point.New(4, 0),
				point.New(4, 3),
				point.New(0, 3),
			),
			rect2: NewFromPoints(
				point.New(1, 1),
				point.New(5, 1),
				point.New(5, 4),
				point.New(1, 4),
			),
			expectEqual: false,
		},
		"rectangles equal within default epsilon": {
			rect1: NewFromPoints(
				point.New(0.0000000000001, 0.0000000000001),
				point.New(4.0000000000001, 0.0000000000001),
				point.New(4.0000000000001, 3.0000000000001),
				point.New(0.0000000000001, 3.0000000000001),
			),
			rect2: NewFromPoints(
				point.New(0, 0),
				point.New(4, 0),
				point.New(4, 3),
				point.New(0, 3),
			),
			expectEqual: true,
		},
		"rectangles not equal just outside default epsilon": {
			rect1: NewFromPoints(
				point.New(0.000000000001, 0.000000000001),
				point.New(4.000000000001, 0.000000000001),
				point.New(4.000000000001, 3.000000000001),
				point.New(0.000000000001, 3.000000000001),
			),
			rect2: NewFromPoints(
				point.New(0, 0),
				point.New(4, 0),
				point.New(4, 3),
				point.New(0, 3),
			),
			expectEqual: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expectEqual, tc.rect1.Eq(tc.rect2))
		})
	}
}

func TestRectangle_Height(t *testing.T) {
	tests := map[string]struct {
		rect     Rectangle
		expected float64
	}{
		"positive height": {
			rect:     New(0, 0, 10, 20),
			expected: 20,
		},
		"negative height": {
			rect:     New(0, 20, 10, 0),
			expected: 20,
		},
		"zero height": {
			rect:     New(5, 5, 15, 5),
			expected: 0,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.rect.Height())
		})
	}
}

func TestRectangle_MarshalUnmarshalJSON(t *testing.T) {
	tests := map[string]struct {
		rectangle Rectangle // Input rectangle
		expected  Rectangle // Expected output after Marshal -> Unmarshal
	}{
		"Rectangle": {
			rectangle: NewFromPoints(
				point.New(3.5, 7.2),
				point.New(8.4, 7.2),
				point.New(3.5, 2.1),
				point.New(8.4, 2.1),
			),
			expected: NewFromPoints(
				point.New(3.5, 7.2),
				point.New(8.4, 7.2),
				point.New(3.5, 2.1),
				point.New(8.4, 2.1),
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Marshal
			data, err := json.Marshal(tc.rectangle)
			require.NoErrorf(t, err, "Failed to marshal %s: %v", tc.rectangle, err)

			var result Rectangle
			err = json.Unmarshal(data, &result)
			require.NoErrorf(t, err, "Failed to unmarshal %s: %v", string(data), err)
			assert.Equalf(t, tc.expected, result, "Expected %v, got %v", tc.expected, result)
		})
	}
}

func TestRectangle_Perimeter(t *testing.T) {
	tests := map[string]struct {
		rect     Rectangle
		expected float64
	}{
		"standard rectangle": {
			rect:     New(0, 0, 10, 20),
			expected: 60,
		},
		"rectangle with swapped corners": {
			rect:     New(10, 20, 0, 0),
			expected: 60,
		},
		"degenerate rectangle (zero width)": {
			rect:     New(5, 5, 5, 15),
			expected: 20,
		},
		"degenerate rectangle (zero height)": {
			rect:     New(5, 5, 15, 5),
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

func TestRectangle_RelationshipToPoint(t *testing.T) {
	rect := New(0, 0, 10, 10)

	tests := map[string]struct {
		point       point.Point
		expectedRel types.Relationship
	}{
		"Point inside rectangle": {
			point:       point.New(5, 5),
			expectedRel: types.RelationshipContainedBy,
		},
		"Point on rectangle edge": {
			point:       point.New(10, 5),
			expectedRel: types.RelationshipIntersection,
		},
		"Point outside rectangle": {
			point:       point.New(15, 5),
			expectedRel: types.RelationshipDisjoint,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expectedRel, rect.RelationshipToPoint(tt.point), "unexpected relationship")
		})
	}
}

func TestRectangle_Scale(t *testing.T) {
	tests := map[string]struct {
		rect     Rectangle
		ref      point.Point
		k        float64
		expected Rectangle
	}{
		"scale by 2 from origin": {
			rect: NewFromPoints(
				point.New(0, 10),
				point.New(10, 10),
				point.New(0, 0),
				point.New(10, 0),
			),
			ref: point.New(0, 0),
			k:   2,
			expected: NewFromPoints(
				point.New(0, 20),
				point.New(20, 20),
				point.New(0, 0),
				point.New(20, 0),
			),
		},
		"scale by 0.5 from center": {
			rect: NewFromPoints(
				point.New(-10, 10),
				point.New(10, 10),
				point.New(-10, -10),
				point.New(10, -10),
			),
			ref: point.New(0, 0),
			k:   0.5,
			expected: NewFromPoints(
				point.New(-5, -5),
				point.New(5, -5),
				point.New(5, 5),
				point.New(-5, 5),
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.rect.Scale(tc.ref, tc.k))
		})
	}
}

func TestRectangle_ScaleWidthHeight(t *testing.T) {
	tests := map[string]struct {
		rect           Rectangle
		scaleWidth     float64
		scaleHeight    float64
		expectedWidth  float64
		expectedHeight float64
	}{
		"scale both dimensions": {
			rect:           New(0, 0, 10, 20),
			scaleWidth:     1.5,
			scaleHeight:    0.5,
			expectedWidth:  15.0,
			expectedHeight: 10.0,
		},
		"scale width only": {
			rect:           New(0, 0, 10, 20),
			scaleWidth:     2.0,
			scaleHeight:    1.0,
			expectedWidth:  20.0,
			expectedHeight: 20.0,
		},
		"scale height only": {
			rect:           New(0, 0, 10, 20),
			scaleWidth:     1.0,
			scaleHeight:    2.0,
			expectedWidth:  10.0,
			expectedHeight: 40.0,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			widthScaled := tc.rect.ScaleWidth(tc.scaleWidth)
			heightScaled := tc.rect.ScaleHeight(tc.scaleHeight)

			assert.InDelta(t, tc.expectedWidth, widthScaled.Width(), 0.001)
			assert.InDelta(t, tc.expectedHeight, heightScaled.Height(), 0.001)
		})
	}
}

func TestRectangle_String(t *testing.T) {
	tests := map[string]struct {
		rect     Rectangle
		expected string
	}{
		"simple rectangle": {
			rect: NewFromPoints(
				point.New(0, 0),
				point.New(4, 0),
				point.New(4, 3),
				point.New(0, 3),
			),
			expected: "[(0,0),(4,3)]",
		},
		"negative coordinates": {
			rect: NewFromPoints(
				point.New(-3, -2),
				point.New(2, -2),
				point.New(2, 1),
				point.New(-3, 1),
			),
			expected: "[(-3,-2),(2,1)]",
		},
		"decimal values": {
			rect: NewFromPoints(
				point.New(1.123, 2.234),
				point.New(4.567, 2.234),
				point.New(4.567, 5.678),
				point.New(1.123, 5.678),
			),
			expected: "[(1.123,2.234),(4.567,5.678)]",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.rect.String())
		})
	}
}

func TestRectangle_ToImageRect(t *testing.T) {
	rect := New(0, 0, 100, 200)
	expected := image.Rect(0, 0, 100, 200)
	assert.Equal(t, expected, rect.ToImageRect())
}

func TestRectangle_Translate(t *testing.T) {
	tests := map[string]struct {
		inputRect    Rectangle
		translateBy  point.Point
		expectedRect Rectangle
	}{
		"translate up and right": {
			inputRect: NewFromPoints(
				point.New(0, 10),
				point.New(10, 10),
				point.New(0, 0),
				point.New(10, 0),
			),
			translateBy: point.New(5, 5),
			expectedRect: NewFromPoints(
				point.New(5, 15),
				point.New(15, 15),
				point.New(5, 5),
				point.New(15, 5),
			),
		},
		"translate down and left": {
			inputRect: NewFromPoints(
				point.New(0, 10),
				point.New(10, 10),
				point.New(0, 0),
				point.New(10, 0),
			),
			translateBy: point.New(-5, -5),
			expectedRect: NewFromPoints(
				point.New(-5, 5),
				point.New(5, 5),
				point.New(-5, -5),
				point.New(5, -5),
			),
		},
		"translate by zero": {
			inputRect: NewFromPoints(
				point.New(0, 10),
				point.New(10, 10),
				point.New(0, 0),
				point.New(10, 0),
			),
			translateBy: point.New(0, 0),
			expectedRect: NewFromPoints(
				point.New(0, 10),
				point.New(10, 10),
				point.New(0, 0),
				point.New(10, 0),
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expectedRect, tc.inputRect.Translate(tc.translateBy))
		})
	}
}

func TestRectangle_Width(t *testing.T) {
	tests := map[string]struct {
		rect     Rectangle
		expected float64
	}{
		"positive width": {
			rect:     New(0, 0, 10, 20),
			expected: 10,
		},
		"negative width": {
			rect:     New(0, 0, -10, 20),
			expected: 10,
		},
		"zero width": {
			rect:     New(5, 5, 5, 15),
			expected: 0,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.rect.Width())
		})
	}
}
