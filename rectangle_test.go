package geom2d

import (
	"github.com/stretchr/testify/assert"
	"image"
	"testing"
)

func TestRectangle_Area(t *testing.T) {
	tests := map[string]struct {
		rect     Rectangle[int]
		expected int
	}{
		"standard rectangle": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 20)),
			expected: 200,
		},
		"rectangle with swapped corners": {
			rect:     NewRectangleByPoints(NewPoint(10, 20), NewPoint(0, 0)),
			expected: 200,
		},
		"degenerate rectangle (zero width)": {
			rect:     NewRectangleByPoints(NewPoint(5, 5), NewPoint(5, 15)),
			expected: 0,
		},
		"degenerate rectangle (zero height)": {
			rect:     NewRectangleByPoints(NewPoint(5, 5), NewPoint(15, 5)),
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

func TestRectangle_AsFloat(t *testing.T) {
	rect := NewRectangleByPoints(NewPoint[int](1, 2), NewPoint[int](10, 20))
	expected := NewRectangleByPoints(NewPoint[float64](1.0, 2.0), NewPoint[float64](10.0, 20.0))
	assert.Equal(t, expected, rect.AsFloat())
}

func TestRectangle_AsInt(t *testing.T) {
	rect := NewRectangleByPoints(NewPoint[float64](1.7, 2.9), NewPoint[float64](10.5, 20.3))
	expected := NewRectangleByPoints(NewPoint(1, 2), NewPoint(10, 20))
	assert.Equal(t, expected, rect.AsInt())
}

func TestRectangle_AsIntRounded(t *testing.T) {
	rect := NewRectangleByPoints(NewPoint[float64](1.7, 2.9), NewPoint[float64](10.5, 20.3))
	expected := NewRectangleByPoints(NewPoint(2, 3), NewPoint(11, 20))
	assert.Equal(t, expected, rect.AsIntRounded())
}

func TestRectangle_ContainsPoint(t *testing.T) {
	tests := map[string]struct {
		rect     Rectangle[int]
		point    Point[int]
		expected bool
	}{
		"point inside rectangle": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(5, 10),
			expected: true,
		},
		"point on rectangle boundary (top-left corner)": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(0, 0),
			expected: true,
		},
		"point on rectangle boundary (bottom-right corner)": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(10, 20),
			expected: true,
		},
		"point outside rectangle": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(15, 10),
			expected: false,
		},
		"point on boundary (horizontal edge)": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(5, 0),
			expected: true,
		},
		"point on boundary (vertical edge)": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 20)),
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
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 20)),
			expected: 20,
		},
		"negative height": {
			rect:     NewRectangleByPoints(NewPoint(0, 20), NewPoint(10, 0)),
			expected: 20,
		},
		"zero height": {
			rect:     NewRectangleByPoints(NewPoint(5, 5), NewPoint(15, 5)),
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

func TestRectangle_IsLineSegmentOnEdgeWithEndTouchingVertex(t *testing.T) {
	tests := map[string]struct {
		rect     Rectangle[int]
		segment  LineSegment[int]
		expected bool
	}{
		"segment on top edge touching top-left vertex": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 10)),
			segment:  NewLineSegment(NewPoint(0, 0), NewPoint(5, 0)),
			expected: true,
		},
		"segment on bottom edge touching bottom-right vertex": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 10)),
			segment:  NewLineSegment(NewPoint(10, 10), NewPoint(5, 10)),
			expected: true,
		},
		"segment on left edge not touching a vertex": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 10)),
			segment:  NewLineSegment(NewPoint(0, 3), NewPoint(0, 7)),
			expected: false,
		},
		"segment on right edge touching top-right vertex": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 10)),
			segment:  NewLineSegment(NewPoint(10, 0), NewPoint(10, 5)),
			expected: true,
		},
		"segment entirely inside rectangle": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 10)),
			segment:  NewLineSegment(NewPoint(3, 3), NewPoint(7, 7)),
			expected: false,
		},
		"segment completely outside rectangle": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 10)),
			segment:  NewLineSegment(NewPoint(15, 15), NewPoint(20, 20)),
			expected: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.rect.IsLineSegmentOnEdgeWithEndTouchingVertex(tt.segment)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestRectangle_LineSegmentEntersAndExits(t *testing.T) {
	tests := map[string]struct {
		rect     Rectangle[int]
		segment  LineSegment[int]
		expected bool
	}{
		"segment entering through top and exiting through bottom": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 10)),
			segment:  NewLineSegment(NewPoint(5, -5), NewPoint(5, 15)),
			expected: true,
		},
		"segment entering through left and exiting through right": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 10)),
			segment:  NewLineSegment(NewPoint(-5, 5), NewPoint(15, 5)),
			expected: true,
		},
		"segment entirely outside": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 10)),
			segment:  NewLineSegment(NewPoint(15, 15), NewPoint(20, 20)),
			expected: false,
		},
		"segment touching but not entering and exiting": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 10)),
			segment:  NewLineSegment(NewPoint(0, 5), NewPoint(10, 5)), // Lies on the top edge without entering or exiting
			expected: false,
		},
		"segment entering through one edge but not exiting": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 10)),
			segment:  NewLineSegment(NewPoint(5, -5), NewPoint(5, 5)), // Enters but does not exit
			expected: false,
		},
		"segment intersecting through two edges diagonally": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 10)),
			segment:  NewLineSegment(NewPoint(-5, -5), NewPoint(15, 15)),
			expected: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.rect.LineSegmentEntersAndExits(tt.segment)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestRectangle_LineSegmentIntersectsEdges(t *testing.T) {
	tests := map[string]struct {
		rect     Rectangle[int]
		segment  LineSegment[int]
		expected bool
	}{
		"segment intersects top and bottom edges": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 10)),
			segment:  NewLineSegment(NewPoint(5, -5), NewPoint(5, 15)),
			expected: true,
		},
		"segment intersects left and right edges": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 10)),
			segment:  NewLineSegment(NewPoint(-5, 5), NewPoint(15, 5)),
			expected: true,
		},
		"segment does not intersect any edge": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 10)),
			segment:  NewLineSegment(NewPoint(15, 15), NewPoint(20, 20)),
			expected: false,
		},
		"segment touches top edge at a point": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 10)),
			segment:  NewLineSegment(NewPoint(5, 0), NewPoint(5, 0)), // Degenerate line touching top edge
			expected: false,
		},
		"segment lies on the top edge without intersecting": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 10)),
			segment:  NewLineSegment(NewPoint(0, 0), NewPoint(10, 0)),
			expected: false,
		},
		"segment intersects at one vertex": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 10)),
			segment:  NewLineSegment(NewPoint(-5, -5), NewPoint(0, 0)), // Intersects at top-left vertex
			expected: false,
		},
		"diagonal segment intersects two edges": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 10)),
			segment:  NewLineSegment(NewPoint(-5, -5), NewPoint(15, 15)),
			expected: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.rect.LineSegmentIntersectsEdges(tt.segment)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestRectangle_Operations(t *testing.T) {
	tests := map[string]struct {
		rect     Rectangle[int]
		op       string
		value    any
		expected Rectangle[float64]
	}{
		"Add vector": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 20)),
			op:       "Add",
			value:    NewPoint(5, 5),
			expected: NewRectangleByPoints(NewPoint(5.0, 5.0), NewPoint(15.0, 25.0)),
		},
		"Sub vector": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 20)),
			op:       "Sub",
			value:    NewPoint(5, 5),
			expected: NewRectangleByPoints(NewPoint(-5.0, -5.0), NewPoint(5.0, 15.0)),
		},
		"Scale by factor 2": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 20)),
			op:       "Scale",
			value:    2.0,
			expected: NewRectangleByPoints(NewPoint(0.0, 0.0), NewPoint(20.0, 40.0)),
		},
		"Div by factor 2": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 20)),
			op:       "Div",
			value:    2.0,
			expected: NewRectangleByPoints(NewPoint(0.0, 0.0), NewPoint(5.0, 10.0)),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var result Rectangle[float64]
			switch tt.op {
			case "Add":
				result = tt.rect.Add(tt.value.(Point[int])).AsFloat()
			case "Sub":
				result = tt.rect.Sub(tt.value.(Point[int])).AsFloat()
			case "Scale":
				result = tt.rect.Scale(tt.value.(float64))
			case "Div":
				result = tt.rect.Div(tt.value.(float64))
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRectangle_Perimeter(t *testing.T) {
	tests := map[string]struct {
		rect     Rectangle[int]
		expected int
	}{
		"standard rectangle": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 20)),
			expected: 60,
		},
		"rectangle with swapped corners": {
			rect:     NewRectangleByPoints(NewPoint(10, 20), NewPoint(0, 0)),
			expected: 60,
		},
		"degenerate rectangle (zero width)": {
			rect:     NewRectangleByPoints(NewPoint(5, 5), NewPoint(5, 15)),
			expected: 20,
		},
		"degenerate rectangle (zero height)": {
			rect:     NewRectangleByPoints(NewPoint(5, 5), NewPoint(15, 5)),
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
			rect: NewRectangleByPoints(NewPoint(2, 3), NewPoint(8, 6)),
			expected: []Point[int]{
				NewPoint(2, 3), // top-left
				NewPoint(8, 3), // top-right
				NewPoint(8, 6), // bottom-right
				NewPoint(2, 6), // bottom-left
			},
		},
		"rectangle with negative coordinates": {
			rect: NewRectangleByPoints(NewPoint(-5, -5), NewPoint(0, 0)),
			expected: []Point[int]{
				NewPoint(-5, -5), // top-left
				NewPoint(0, -5),  // top-right
				NewPoint(0, 0),   // bottom-right
				NewPoint(-5, 0),  // bottom-left
			},
		},
		"rectangle with zero-width or height": {
			rect: NewRectangleByPoints(NewPoint(1, 1), NewPoint(1, 4)),
			expected: []Point[int]{
				NewPoint(1, 1), // top-left
				NewPoint(1, 1), // top-right (overlapping)
				NewPoint(1, 4), // bottom-right
				NewPoint(1, 4), // bottom-left (overlapping)
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.rect.Points()
			for i, expectedPoint := range tt.expected {
				assert.Equal(t, expectedPoint, actual[i], "Point mismatch at index %d", i)
			}
		})
	}
}

func TestRectangle_RelationshipToLineSegment(t *testing.T) {
	tests := map[string]struct {
		rect     Rectangle[int]
		segment  LineSegment[int]
		expected RectangleSegmentRelationship
	}{
		"segment completely outside": {
			rect:     NewRectangleByPoints(NewPoint(1, 1), NewPoint(5, 5)),
			segment:  NewLineSegment(NewPoint(6, 6), NewPoint(7, 7)),
			expected: RSROutside,
		},
		"segment outside with one end touching edge": {
			rect:     NewRectangleByPoints(NewPoint(1, 1), NewPoint(5, 5)),
			segment:  NewLineSegment(NewPoint(0, 3), NewPoint(1, 3)),
			expected: RSROutsideEndTouchesEdge,
		},
		"segment outside with one end touching vertex": {
			rect:     NewRectangleByPoints(NewPoint(1, 1), NewPoint(5, 5)),
			segment:  NewLineSegment(NewPoint(0, 0), NewPoint(1, 1)),
			expected: RSROutsideEndTouchesVertex,
		},
		"segment completely inside": {
			rect:     NewRectangleByPoints(NewPoint(1, 1), NewPoint(5, 5)),
			segment:  NewLineSegment(NewPoint(2, 2), NewPoint(3, 3)),
			expected: RSRInside,
		},
		"segment inside with one end touching edge": {
			rect:     NewRectangleByPoints(NewPoint(1, 1), NewPoint(5, 5)),
			segment:  NewLineSegment(NewPoint(2, 2), NewPoint(1, 3)),
			expected: RSRInsideEndTouchesEdge,
		},
		"segment inside with one end touching vertex": {
			rect:     NewRectangleByPoints(NewPoint(1, 1), NewPoint(5, 5)),
			segment:  NewLineSegment(NewPoint(3, 3), NewPoint(1, 1)),
			expected: RSRInsideEndTouchesVertex,
		},
		"segment lying on edge": {
			rect:     NewRectangleByPoints(NewPoint(1, 1), NewPoint(5, 5)),
			segment:  NewLineSegment(NewPoint(2, 1), NewPoint(4, 1)),
			expected: RSROnEdge,
		},
		"segment on edge with end touching vertex": {
			rect:     NewRectangleByPoints(NewPoint(1, 1), NewPoint(5, 5)),
			segment:  NewLineSegment(NewPoint(1, 1), NewPoint(1, 5)),
			expected: RSROnEdgeEndTouchesVertex,
		},
		"segment intersecting through one edge diagonally": {
			rect:     NewRectangleByPoints(NewPoint(1, 1), NewPoint(5, 5)),
			segment:  NewLineSegment(NewPoint(0, 0), NewPoint(3, 3)),
			expected: RSRIntersects,
		},
		"segment entering through edge": {
			rect:     NewRectangleByPoints(NewPoint(1, 1), NewPoint(5, 5)),
			segment:  NewLineSegment(NewPoint(0, 3), NewPoint(3, 3)),
			expected: RSRIntersects,
		},
		"segment intersecting through two edges diagonally": {
			rect:     NewRectangleByPoints(NewPoint(1, 1), NewPoint(5, 5)),
			segment:  NewLineSegment(NewPoint(0, 0), NewPoint(6, 6)),
			expected: RSREntersAndExits,
		},
		"segment entering and exiting through different edges": {
			rect:     NewRectangleByPoints(NewPoint(1, 1), NewPoint(5, 5)),
			segment:  NewLineSegment(NewPoint(0, 3), NewPoint(6, 3)),
			expected: RSREntersAndExits,
		},
		"degenerate segment on vertex": {
			rect:     NewRectangleByPoints(NewPoint(1, 1), NewPoint(5, 5)),
			segment:  NewLineSegment(NewPoint(1, 1), NewPoint(1, 1)),
			expected: RSROutsideEndTouchesVertex,
		},
		"segment through opposite vertices": {
			rect:     NewRectangleByPoints(NewPoint(1, 1), NewPoint(5, 5)),
			segment:  NewLineSegment(NewPoint(1, 1), NewPoint(5, 5)),
			expected: RSRInsideEndTouchesVertex,
		},
		"partially collinear segment on one edge": {
			rect:     NewRectangleByPoints(NewPoint(1, 1), NewPoint(5, 5)),
			segment:  NewLineSegment(NewPoint(1, 2), NewPoint(1, 0)),
			expected: RSROutsideEndTouchesEdge,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := tt.rect.RelationshipToLineSegment(tt.segment)
			assert.Equal(t, tt.expected, result)

			// flip line segment, should yield same result
			segFlipped := NewLineSegment(tt.segment.End(), tt.segment.Start())
			result = tt.rect.RelationshipToLineSegment(segFlipped)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRectangle_RelationshipToPoint(t *testing.T) {
	tests := map[string]struct {
		rect     Rectangle[int]
		point    Point[int]
		expected PointRectangleRelationship
	}{
		"point inside rectangle": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(5, 10),
			expected: PRRInside,
		},
		"point outside rectangle": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(15, 10),
			expected: PRROutside,
		},
		"point on top-left vertex": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(0, 0),
			expected: PRROnVertex,
		},
		"point on bottom-right vertex": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(10, 20),
			expected: PRROnVertex,
		},
		"point on bottom-left vertex": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(0, 20),
			expected: PRROnVertex,
		},
		"point on top-right vertex": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(10, 0),
			expected: PRROnVertex,
		},
		"point on top edge": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(5, 0),
			expected: PRROnEdge,
		},
		"point on left edge": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(0, 10),
			expected: PRROnEdge,
		},
		"point on bottom edge": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(5, 20),
			expected: PRROnEdge,
		},
		"point on right edge": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(10, 10),
			expected: PRROnEdge,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.rect.RelationshipToPoint(tt.point)
			assert.Equal(t, tt.expected, actual)
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
			rect:           NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 20)),
			scaleWidth:     1.5,
			scaleHeight:    0.5,
			expectedWidth:  15.0,
			expectedHeight: 10.0,
		},
		"scale width only": {
			rect:           NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 20)),
			scaleWidth:     2.0,
			scaleHeight:    1.0,
			expectedWidth:  20.0,
			expectedHeight: 20.0,
		},
		"scale height only": {
			rect:           NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 20)),
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
	rect := NewRectangleByPoints(NewPoint(0, 0), NewPoint(100, 200))
	expected := image.Rect(0, 0, 100, 200)
	assert.Equal(t, expected, rect.ToImageRect())
}

func TestRectangle_Width(t *testing.T) {
	tests := map[string]struct {
		rect     Rectangle[int]
		expected int
	}{
		"positive width": {
			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 20)),
			expected: 10,
		},
		"negative width": {
			rect:     NewRectangleByPoints(NewPoint(10, 0), NewPoint(0, 20)),
			expected: 10,
		},
		"zero width": {
			rect:     NewRectangleByPoints(NewPoint(5, 5), NewPoint(5, 15)),
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

func TestNewRectangleByDimensions(t *testing.T) {
	tests := map[string]struct {
		topLeft  Point[int]
		width    int
		height   int
		expected Rectangle[int]
	}{
		"positive width and height": {
			topLeft: NewPoint(0, 0),
			width:   10,
			height:  10,
			expected: Rectangle[int]{
				topLeft:     NewPoint(0, 0),
				bottomRight: NewPoint(10, 10),
			},
		},
		"negative width and height": {
			topLeft: NewPoint(5, 5),
			width:   -10,
			height:  -10,
			expected: Rectangle[int]{
				topLeft:     NewPoint(5, 5),
				bottomRight: NewPoint(-5, -5),
			},
		},
		"zero width and height": {
			topLeft: NewPoint(3, 3),
			width:   0,
			height:  0,
			expected: Rectangle[int]{
				topLeft:     NewPoint(3, 3),
				bottomRight: NewPoint(3, 3),
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := NewRectangleByDimensions(tt.topLeft, tt.width, tt.height)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestNewRectangleByPoints(t *testing.T) {
	tests := map[string]struct {
		topLeft     Point[int]
		bottomRight Point[int]
		expected    Rectangle[int]
	}{
		"positive coordinates": {
			topLeft:     NewPoint(0, 0),
			bottomRight: NewPoint(10, 10),
			expected: Rectangle[int]{
				topLeft:     NewPoint(0, 0),
				bottomRight: NewPoint(10, 10),
			},
		},
		"negative coordinates": {
			topLeft:     NewPoint(-10, -10),
			bottomRight: NewPoint(0, 0),
			expected: Rectangle[int]{
				topLeft:     NewPoint(-10, -10),
				bottomRight: NewPoint(0, 0),
			},
		},
		"mixed coordinates": {
			topLeft:     NewPoint(-5, 5),
			bottomRight: NewPoint(5, -5),
			expected: Rectangle[int]{
				topLeft:     NewPoint(-5, 5),
				bottomRight: NewPoint(5, -5),
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := NewRectangleByPoints(tt.topLeft, tt.bottomRight)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestNewRectangleFromImageRect(t *testing.T) {
	tests := map[string]struct {
		imgRect  image.Rectangle
		expected Rectangle[int]
	}{
		"standard image rectangle": {
			imgRect: image.Rect(0, 0, 100, 200),
			expected: Rectangle[int]{
				topLeft:     NewPoint(0, 0),
				bottomRight: NewPoint(100, 200),
			},
		},
		"negative coordinates": {
			imgRect: image.Rect(-50, -50, 50, 50),
			expected: Rectangle[int]{
				topLeft:     NewPoint(-50, -50),
				bottomRight: NewPoint(50, 50),
			},
		},
		"degenerate rectangle": {
			imgRect: image.Rect(10, 10, 10, 10),
			expected: Rectangle[int]{
				topLeft:     NewPoint(10, 10),
				bottomRight: NewPoint(10, 10),
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := NewRectangleFromImageRect(tt.imgRect)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
