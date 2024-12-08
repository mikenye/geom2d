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
			rect:     NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			expected: 200,
		},
		"rectangle with swapped corners": {
			rect:     NewRectangleByOppositeCorners(NewPoint(10, 20), NewPoint(0, 0)),
			expected: 200,
		},
		"degenerate rectangle (zero width)": {
			rect:     NewRectangleByOppositeCorners(NewPoint(5, 5), NewPoint(5, 15)),
			expected: 0,
		},
		"degenerate rectangle (zero height)": {
			rect:     NewRectangleByOppositeCorners(NewPoint(5, 5), NewPoint(15, 5)),
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
	rect := NewRectangleByOppositeCorners(NewPoint[int](1, 2), NewPoint[int](10, 20))
	expected := NewRectangleByOppositeCorners(NewPoint[float64](1.0, 2.0), NewPoint[float64](10.0, 20.0))
	assert.Equal(t, expected, rect.AsFloat())
}

func TestRectangle_AsInt(t *testing.T) {
	rect := NewRectangleByOppositeCorners(NewPoint[float64](1.7, 2.9), NewPoint[float64](10.5, 20.3))
	expected := NewRectangleByOppositeCorners(NewPoint(1, 2), NewPoint(10, 20))
	assert.Equal(t, expected, rect.AsInt())
}

func TestRectangle_AsIntRounded(t *testing.T) {
	rect := NewRectangleByOppositeCorners(NewPoint[float64](1.7, 2.9), NewPoint[float64](10.5, 20.3))
	expected := NewRectangleByOppositeCorners(NewPoint(2, 3), NewPoint(11, 20))
	assert.Equal(t, expected, rect.AsIntRounded())
}

func TestRectangle_ContainsPoint(t *testing.T) {
	tests := map[string]struct {
		rect     Rectangle[int]
		point    Point[int]
		expected bool
	}{
		"point inside rectangle": {
			rect:     NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(5, 10),
			expected: true,
		},
		"point on rectangle boundary (top-left corner)": {
			rect:     NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(0, 0),
			expected: true,
		},
		"point on rectangle boundary (bottom-right corner)": {
			rect:     NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(10, 20),
			expected: true,
		},
		"point outside rectangle": {
			rect:     NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(15, 10),
			expected: false,
		},
		"point on boundary (horizontal edge)": {
			rect:     NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(5, 0),
			expected: true,
		},
		"point on boundary (vertical edge)": {
			rect:     NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
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
			rect:     NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			expected: 20,
		},
		"negative height": {
			rect:     NewRectangleByOppositeCorners(NewPoint(0, 20), NewPoint(10, 0)),
			expected: 20,
		},
		"zero height": {
			rect:     NewRectangleByOppositeCorners(NewPoint(5, 5), NewPoint(15, 5)),
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
			rect:     NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 10)),
			segment:  NewLineSegment(NewPoint(0, 0), NewPoint(5, 0)),
			expected: true,
		},
		"segment on bottom edge touching bottom-right vertex": {
			rect:     NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 10)),
			segment:  NewLineSegment(NewPoint(10, 10), NewPoint(5, 10)),
			expected: true,
		},
		"segment on left edge not touching a vertex": {
			rect:     NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 10)),
			segment:  NewLineSegment(NewPoint(0, 3), NewPoint(0, 7)),
			expected: false,
		},
		"segment on right edge touching top-right vertex": {
			rect:     NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 10)),
			segment:  NewLineSegment(NewPoint(10, 0), NewPoint(10, 5)),
			expected: true,
		},
		"segment entirely inside rectangle": {
			rect:     NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 10)),
			segment:  NewLineSegment(NewPoint(3, 3), NewPoint(7, 7)),
			expected: false,
		},
		"segment completely outside rectangle": {
			rect:     NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 10)),
			segment:  NewLineSegment(NewPoint(15, 15), NewPoint(20, 20)),
			expected: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.rect.isLineSegmentOnEdgeWithEndTouchingVertex(tt.segment)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

//func TestRectangle_LineSegmentEntersAndExits(t *testing.T) {
//	tests := map[string]struct {
//		rect     Rectangle[int]
//		segment  LineSegment[int]
//		expected bool
//	}{
//		"segment entering through top and exiting through bottom": {
//			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 10)),
//			segment:  NewLineSegment(NewPoint(5, -5), NewPoint(5, 15)),
//			expected: true,
//		},
//		"segment entering through left and exiting through right": {
//			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 10)),
//			segment:  NewLineSegment(NewPoint(-5, 5), NewPoint(15, 5)),
//			expected: true,
//		},
//		"segment entirely outside": {
//			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 10)),
//			segment:  NewLineSegment(NewPoint(15, 15), NewPoint(20, 20)),
//			expected: false,
//		},
//		"segment touching but not entering and exiting": {
//			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 10)),
//			segment:  NewLineSegment(NewPoint(0, 5), NewPoint(10, 5)), // Lies on the top edge without entering or exiting
//			expected: false,
//		},
//		"segment entering through one edge but not exiting": {
//			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 10)),
//			segment:  NewLineSegment(NewPoint(5, -5), NewPoint(5, 5)), // Enters but does not exit
//			expected: false,
//		},
//		"segment intersecting through two edges diagonally": {
//			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 10)),
//			segment:  NewLineSegment(NewPoint(-5, -5), NewPoint(15, 15)),
//			expected: true,
//		},
//	}
//
//	for name, tt := range tests {
//		t.Run(name, func(t *testing.T) {
//			actual := tt.rect.LineSegmentEntersAndExits(tt.segment)
//			assert.Equal(t, tt.expected, actual)
//		})
//	}
//}

//func TestRectangle_LineSegmentIntersectsEdges(t *testing.T) {
//	tests := map[string]struct {
//		rect     Rectangle[int]
//		segment  LineSegment[int]
//		expected bool
//	}{
//		"segment intersects top and bottom edges": {
//			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 10)),
//			segment:  NewLineSegment(NewPoint(5, -5), NewPoint(5, 15)),
//			expected: true,
//		},
//		"segment intersects left and right edges": {
//			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 10)),
//			segment:  NewLineSegment(NewPoint(-5, 5), NewPoint(15, 5)),
//			expected: true,
//		},
//		"segment does not intersect any edge": {
//			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 10)),
//			segment:  NewLineSegment(NewPoint(15, 15), NewPoint(20, 20)),
//			expected: false,
//		},
//		"segment touches top edge at a point": {
//			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 10)),
//			segment:  NewLineSegment(NewPoint(5, 0), NewPoint(5, 0)), // Degenerate line touching top edge
//			expected: false,
//		},
//		"segment lies on the top edge without intersecting": {
//			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 10)),
//			segment:  NewLineSegment(NewPoint(0, 0), NewPoint(10, 0)),
//			expected: false,
//		},
//		"segment intersects at one vertex": {
//			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 10)),
//			segment:  NewLineSegment(NewPoint(-5, -5), NewPoint(0, 0)), // Intersects at top-left vertex
//			expected: false,
//		},
//		"diagonal segment intersects two edges": {
//			rect:     NewRectangleByPoints(NewPoint(0, 0), NewPoint(10, 10)),
//			segment:  NewLineSegment(NewPoint(-5, -5), NewPoint(15, 15)),
//			expected: true,
//		},
//	}
//
//	for name, tt := range tests {
//		t.Run(name, func(t *testing.T) {
//			actual := tt.rect.LineSegmentIntersectsEdges(tt.segment)
//			assert.Equal(t, tt.expected, actual)
//		})
//	}
//}

// todo: split into separate tests, add is now translate
//func TestRectangle_Operations(t *testing.T) {
//	tests := map[string]struct {
//		rect     Rectangle[int]
//		op       string
//		value    any
//		expected Rectangle[float64]
//	}{
//		"Add vector": {
//			rect:     NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
//			op:       "Add",
//			value:    NewPoint(5, 5),
//			expected: NewRectangleByOppositeCorners(NewPoint(5.0, 5.0), NewPoint(15.0, 25.0)),
//		},
//		"Sub vector": {
//			rect:     NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
//			op:       "Sub",
//			value:    NewPoint(5, 5),
//			expected: NewRectangleByOppositeCorners(NewPoint(-5.0, -5.0), NewPoint(5.0, 15.0)),
//		},
//		//"Scale by factor 2": {
//		//	rect:     NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
//		//	op:       "Scale",
//		//	value:    2.0,
//		//	expected: NewRectangleByOppositeCorners(NewPoint(0.0, 0.0), NewPoint(20.0, 40.0)),
//		//},
//		//"Div by factor 2": {
//		//	rect:     NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
//		//	op:       "Div",
//		//	value:    2.0,
//		//	expected: NewRectangleByOppositeCorners(NewPoint(0.0, 0.0), NewPoint(5.0, 10.0)),
//		//},
//	}
//
//	for name, tt := range tests {
//		t.Run(name, func(t *testing.T) {
//			var result Rectangle[float64]
//			switch tt.op {
//			case "Add":
//				result = tt.rect.Add(tt.value.(Point[int])).AsFloat()
//			case "Sub":
//				result = tt.rect.Sub(tt.value.(Point[int])).AsFloat()
//				//case "Scale":
//				//	result = tt.rect.Scale(tt.value.(float64))
//				//case "Div":
//				//	result = tt.rect.Div(tt.value.(float64))
//			}
//			assert.Equal(t, tt.expected, result)
//		})
//	}
//}

func TestRectangle_Perimeter(t *testing.T) {
	tests := map[string]struct {
		rect     Rectangle[int]
		expected int
	}{
		"standard rectangle": {
			rect:     NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			expected: 60,
		},
		"rectangle with swapped corners": {
			rect:     NewRectangleByOppositeCorners(NewPoint(10, 20), NewPoint(0, 0)),
			expected: 60,
		},
		"degenerate rectangle (zero width)": {
			rect:     NewRectangleByOppositeCorners(NewPoint(5, 5), NewPoint(5, 15)),
			expected: 20,
		},
		"degenerate rectangle (zero height)": {
			rect:     NewRectangleByOppositeCorners(NewPoint(5, 5), NewPoint(15, 5)),
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
			rect: NewRectangleByOppositeCorners(NewPoint(2, 3), NewPoint(8, 6)),
			expected: []Point[int]{
				NewPoint(2, 3), // top-left
				NewPoint(8, 3), // top-right
				NewPoint(8, 6), // bottom-right
				NewPoint(2, 6), // bottom-left
			},
		},
		"rectangle with negative coordinates": {
			rect: NewRectangleByOppositeCorners(NewPoint(-5, -5), NewPoint(0, 0)),
			expected: []Point[int]{
				NewPoint(-5, -5), // top-left
				NewPoint(0, -5),  // top-right
				NewPoint(0, 0),   // bottom-right
				NewPoint(-5, 0),  // bottom-left
			},
		},
		"rectangle with zero-width or height": {
			rect: NewRectangleByOppositeCorners(NewPoint(1, 1), NewPoint(1, 4)),
			expected: []Point[int]{
				NewPoint(1, 1), // top-left
				NewPoint(1, 1), // top-right (overlapping)
				NewPoint(1, 4), // bottom-right
				NewPoint(1, 4), // bottom-left (overlapping)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tc.rect.Points()
			assert.Len(t, actual, 4)
			for _, expectedPoint := range tc.expected {
				assert.Contains(t, actual, expectedPoint)
			}
		})
	}
}

func TestRectangle_RelationshipToLineSegment(t *testing.T) {
	tests := map[string]struct {
		rect     Rectangle[int]
		segment  LineSegment[int]
		expected RelationshipLineSegmentRectangle
	}{
		"segment completely outside": {
			rect:     NewRectangleByOppositeCorners(NewPoint(1, 1), NewPoint(5, 5)),
			segment:  NewLineSegment(NewPoint(6, 6), NewPoint(7, 7)),
			expected: RelationshipLineSegmentRectangleMiss,
		},
		"segment outside with one end touching edge": {
			rect:     NewRectangleByOppositeCorners(NewPoint(1, 1), NewPoint(5, 5)),
			segment:  NewLineSegment(NewPoint(0, 3), NewPoint(1, 3)),
			expected: RelationshipLineSegmentRectangleEndTouchesEdgeExternally,
		},
		"segment outside with one end touching vertex": {
			rect:     NewRectangleByOppositeCorners(NewPoint(1, 1), NewPoint(5, 5)),
			segment:  NewLineSegment(NewPoint(0, 0), NewPoint(1, 1)),
			expected: RelationshipLineSegmentRectangleEndTouchesVertexExternally,
		},
		"segment completely inside": {
			rect:     NewRectangleByOppositeCorners(NewPoint(1, 1), NewPoint(5, 5)),
			segment:  NewLineSegment(NewPoint(2, 2), NewPoint(3, 3)),
			expected: RelationshipLineSegmentRectangleContainedByRectangle,
		},
		"segment inside with one end touching edge": {
			rect:     NewRectangleByOppositeCorners(NewPoint(1, 1), NewPoint(5, 5)),
			segment:  NewLineSegment(NewPoint(2, 2), NewPoint(1, 3)),
			expected: RelationshipLineSegmentRectangleEndTouchesEdgeInternally,
		},
		"segment inside with one end touching vertex": {
			rect:     NewRectangleByOppositeCorners(NewPoint(1, 1), NewPoint(5, 5)),
			segment:  NewLineSegment(NewPoint(3, 3), NewPoint(1, 1)),
			expected: RelationshipLineSegmentRectangleEndTouchesVertexInternally,
		},
		"segment lying on edge": {
			rect:     NewRectangleByOppositeCorners(NewPoint(1, 1), NewPoint(5, 5)),
			segment:  NewLineSegment(NewPoint(2, 1), NewPoint(4, 1)),
			expected: RelationshipLineSegmentRectangleEdgeCollinear,
		},
		"segment on edge with end touching vertex": {
			rect:     NewRectangleByOppositeCorners(NewPoint(1, 1), NewPoint(5, 5)),
			segment:  NewLineSegment(NewPoint(1, 1), NewPoint(1, 5)),
			expected: RelationshipLineSegmentRectangleEdgeCollinearTouchingVertex,
		},
		"segment intersecting through one edge diagonally": {
			rect:     NewRectangleByOppositeCorners(NewPoint(1, 1), NewPoint(5, 5)),
			segment:  NewLineSegment(NewPoint(0, 0), NewPoint(3, 3)),
			expected: RelationshipLineSegmentRectangleIntersects,
		},
		"segment entering through edge": {
			rect:     NewRectangleByOppositeCorners(NewPoint(1, 1), NewPoint(5, 5)),
			segment:  NewLineSegment(NewPoint(0, 3), NewPoint(3, 3)),
			expected: RelationshipLineSegmentRectangleIntersects,
		},
		"segment intersecting through two edges diagonally": {
			rect:     NewRectangleByOppositeCorners(NewPoint(1, 1), NewPoint(5, 5)),
			segment:  NewLineSegment(NewPoint(0, 0), NewPoint(6, 6)),
			expected: RelationshipLineSegmentRectangleEntersAndExits,
		},
		"segment entering and exiting through different edges": {
			rect:     NewRectangleByOppositeCorners(NewPoint(1, 1), NewPoint(5, 5)),
			segment:  NewLineSegment(NewPoint(0, 3), NewPoint(6, 3)),
			expected: RelationshipLineSegmentRectangleEntersAndExits,
		},
		"degenerate segment on vertex": {
			rect:     NewRectangleByOppositeCorners(NewPoint(1, 1), NewPoint(5, 5)),
			segment:  NewLineSegment(NewPoint(1, 1), NewPoint(1, 1)),
			expected: RelationshipLineSegmentRectangleEndTouchesVertexExternally,
		},
		"segment through opposite vertices": {
			rect:     NewRectangleByOppositeCorners(NewPoint(1, 1), NewPoint(5, 5)),
			segment:  NewLineSegment(NewPoint(1, 1), NewPoint(5, 5)),
			expected: RelationshipLineSegmentRectangleEndTouchesVertexInternally,
		},
		"partially collinear segment on one edge": {
			rect:     NewRectangleByOppositeCorners(NewPoint(1, 1), NewPoint(5, 5)),
			segment:  NewLineSegment(NewPoint(1, 2), NewPoint(1, 0)),
			expected: RelationshipLineSegmentRectangleEndTouchesEdgeExternally,
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
		expected RelationshipPointRectangle
	}{
		"point inside rectangle": {
			rect:     NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(5, 10),
			expected: RelationshipPointRectangleContainedByRectangle,
		},
		"point outside rectangle": {
			rect:     NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(15, 10),
			expected: RelationshipPointRectangleMiss,
		},
		"point on top-left vertex": {
			rect:     NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(0, 0),
			expected: RelationshipPointRectanglePointOnVertex,
		},
		"point on bottom-right vertex": {
			rect:     NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(10, 20),
			expected: RelationshipPointRectanglePointOnVertex,
		},
		"point on bottom-left vertex": {
			rect:     NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(0, 20),
			expected: RelationshipPointRectanglePointOnVertex,
		},
		"point on top-right vertex": {
			rect:     NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(10, 0),
			expected: RelationshipPointRectanglePointOnVertex,
		},
		"point on top edge": {
			rect:     NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(5, 0),
			expected: RelationshipPointRectanglePointOnEdge,
		},
		"point on left edge": {
			rect:     NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(0, 10),
			expected: RelationshipPointRectanglePointOnEdge,
		},
		"point on bottom edge": {
			rect:     NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(5, 20),
			expected: RelationshipPointRectanglePointOnEdge,
		},
		"point on right edge": {
			rect:     NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			point:    NewPoint(10, 10),
			expected: RelationshipPointRectanglePointOnEdge,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.rect.RelationshipToPoint(tt.point)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestRectangleRelationshipToRectangle(t *testing.T) {
	tests := map[string]struct {
		rect1       Rectangle[int]
		rect2       Rectangle[int]
		expectedRel RelationshipRectangleRectangle
	}{
		"Disjoint Rectangles": {
			rect1: NewRectangle([]Point[int]{
				NewPoint(0, 0),
				NewPoint(10, 0),
				NewPoint(10, 10),
				NewPoint(0, 10),
			}),
			rect2: NewRectangle([]Point[int]{
				NewPoint(20, 20),
				NewPoint(30, 20),
				NewPoint(30, 30),
				NewPoint(20, 30),
			}),
			expectedRel: RelationshipRectangleRectangleMiss,
		},
		"Touching Edge": {
			rect1: NewRectangle([]Point[int]{
				NewPoint(0, 0),
				NewPoint(10, 0),
				NewPoint(10, 10),
				NewPoint(0, 10),
			}),
			rect2: NewRectangle([]Point[int]{
				NewPoint(10, 0),
				NewPoint(20, 0),
				NewPoint(20, 10),
				NewPoint(10, 10),
			}),
			expectedRel: RelationshipRectangleRectangleSharedEdge,
		},
		"Touching Vertex": {
			rect1: NewRectangle([]Point[int]{
				NewPoint(0, 0),
				NewPoint(10, 0),
				NewPoint(10, 10),
				NewPoint(0, 10),
			}),
			rect2: NewRectangle([]Point[int]{
				NewPoint(10, 10),
				NewPoint(20, 10),
				NewPoint(20, 20),
				NewPoint(10, 20),
			}),
			expectedRel: RelationshipRectangleRectangleSharedVertex,
		},
		"Intersecting Rectangles": {
			rect1: NewRectangle([]Point[int]{
				NewPoint(0, 0),
				NewPoint(10, 0),
				NewPoint(10, 10),
				NewPoint(0, 10),
			}),
			rect2: NewRectangle([]Point[int]{
				NewPoint(5, 5),
				NewPoint(15, 5),
				NewPoint(15, 15),
				NewPoint(5, 15),
			}),
			expectedRel: RelationshipRectangleRectangleIntersection,
		},
		"Contained Rectangles": {
			rect1: NewRectangle([]Point[int]{
				NewPoint(0, 0),
				NewPoint(20, 0),
				NewPoint(20, 20),
				NewPoint(0, 20),
			}),
			rect2: NewRectangle([]Point[int]{
				NewPoint(5, 5),
				NewPoint(15, 5),
				NewPoint(15, 15),
				NewPoint(5, 15),
			}),
			expectedRel: RelationshipRectangleRectangleContained,
		},
		"Touching Contained Rectangles": {
			rect1: NewRectangle([]Point[int]{
				NewPoint(0, 0),
				NewPoint(20, 0),
				NewPoint(20, 20),
				NewPoint(0, 20),
			}),
			rect2: NewRectangle([]Point[int]{
				NewPoint(0, 0),
				NewPoint(20, 0),
				NewPoint(20, 20),
				NewPoint(0, 20),
			}),
			expectedRel: RelationshipRectangleRectangleEqual, // Adjusted from RelationshipRectangleRectangleContainedTouching to RelationshipRectangleRectangleEqual
		},
		"Equal Rectangles": {
			rect1: NewRectangle([]Point[int]{
				NewPoint(0, 0),
				NewPoint(10, 0),
				NewPoint(10, 10),
				NewPoint(0, 10),
			}),
			rect2: NewRectangle([]Point[int]{
				NewPoint(0, 0),
				NewPoint(10, 0),
				NewPoint(10, 10),
				NewPoint(0, 10),
			}),
			expectedRel: RelationshipRectangleRectangleEqual,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actualRel := test.rect1.RelationshipToRectangle(test.rect2)
			assert.Equal(t, test.expectedRel, actualRel, "Relationship mismatch for test case: %s", name)
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
			rect:           NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			scaleWidth:     1.5,
			scaleHeight:    0.5,
			expectedWidth:  15.0,
			expectedHeight: 10.0,
		},
		"scale width only": {
			rect:           NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			scaleWidth:     2.0,
			scaleHeight:    1.0,
			expectedWidth:  20.0,
			expectedHeight: 20.0,
		},
		"scale height only": {
			rect:           NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
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
	rect := NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(100, 200))
	expected := image.Rect(0, 0, 100, 200)
	assert.Equal(t, expected, rect.ToImageRect())
}

func TestRectangle_Width(t *testing.T) {
	tests := map[string]struct {
		rect     Rectangle[int]
		expected int
	}{
		"positive width": {
			rect:     NewRectangleByOppositeCorners(NewPoint(0, 0), NewPoint(10, 20)),
			expected: 10,
		},
		"negative width": {
			rect:     NewRectangleByOppositeCorners(NewPoint(10, 0), NewPoint(0, 20)),
			expected: 10,
		},
		"zero width": {
			rect:     NewRectangleByOppositeCorners(NewPoint(5, 5), NewPoint(5, 15)),
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

func TestNewRectangleByPoints(t *testing.T) {
	tests := map[string]struct {
		corner              Point[int]
		oppositeCorner      Point[int]
		expectedTopLeft     Point[int]
		expectedBottomRight Point[int]
	}{
		"Standard case with top-left first": {
			corner:              NewPoint(0, 10),
			oppositeCorner:      NewPoint(10, 0),
			expectedTopLeft:     NewPoint(0, 10),
			expectedBottomRight: NewPoint(10, 0),
		},
		"Standard case with bottom-right first": {
			corner:              NewPoint(10, 0),
			oppositeCorner:      NewPoint(0, 10),
			expectedTopLeft:     NewPoint(0, 10),
			expectedBottomRight: NewPoint(10, 0),
		},
		"Negative coordinates": {
			corner:              NewPoint(-10, -10),
			oppositeCorner:      NewPoint(-5, -5),
			expectedTopLeft:     NewPoint(-10, -5),
			expectedBottomRight: NewPoint(-5, -10),
		},
		"Mixed positive and negative coordinates": {
			corner:              NewPoint(-10, 10),
			oppositeCorner:      NewPoint(5, -5),
			expectedTopLeft:     NewPoint(-10, 10),
			expectedBottomRight: NewPoint(5, -5),
		},
		"Single point rectangle": {
			corner:              NewPoint(5, 5),
			oppositeCorner:      NewPoint(5, 5),
			expectedTopLeft:     NewPoint(5, 5),
			expectedBottomRight: NewPoint(5, 5),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			rect := NewRectangleByOppositeCorners(tt.corner, tt.oppositeCorner)

			assert.Equal(t, tt.expectedTopLeft, rect.topLeft, "Top-left corner mismatch")
			assert.Equal(t, tt.expectedBottomRight, rect.bottomRight, "Bottom-right corner mismatch")
		})
	}
}

func TestNewRectangleFromImageRect(t *testing.T) {
	tests := map[string]struct {
		imgRect  image.Rectangle
		expected []Point[int]
	}{
		"standard image rectangle": {
			imgRect: image.Rect(0, 0, 100, 200),
			expected: []Point[int]{
				NewPoint(0, 0),
				NewPoint(0, 200),
				NewPoint(100, 200),
				NewPoint(100, 0),
			},
		},
		"negative coordinates": {
			imgRect: image.Rect(-50, -50, 50, 50),
			expected: []Point[int]{
				NewPoint(-50, -50),
				NewPoint(-50, 50),
				NewPoint(50, 50),
				NewPoint(50, -50),
			},
		},
		"degenerate rectangle": {
			imgRect: image.Rect(10, 10, 10, 10),
			expected: []Point[int]{
				NewPoint(10, 10),
				NewPoint(10, 10),
				NewPoint(10, 10),
				NewPoint(10, 10),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := NewRectangleFromImageRect(tc.imgRect)
			for _, p := range actual.Points() {
				assert.Contains(t, tc.expected, p)
			}
		})
	}
}
