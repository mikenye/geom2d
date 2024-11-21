// Package geom2d provides foundational definitions, constants, and utility types for 2D geometry operations.
//
// This file contains core type definitions, global variables, constants, and utility functions
// that are shared across the package. It serves as the basis for other types like Point and LineSegment,
// offering foundational functionality and supporting mathematical constants and helper functions.

package geom2d

// CircleLineSegmentRelationship (CLR) describes the spatial relationship between a line segment
// and a circle.
type CircleLineSegmentRelationship int

const (
	// CLROutside indicates that the line segment lies completely outside the circle,
	// with no intersection or tangency.
	CLROutside CircleLineSegmentRelationship = iota

	// CLRInside indicates that the line segment lies completely within the circle,
	// with both endpoints inside the circle's boundary.
	CLRInside

	// CLRIntersecting indicates that the line segment intersects the circle at exactly
	// two distinct points.
	CLRIntersecting

	// CLRTangent indicates that the line segment is tangent to the circle, touching it
	// at exactly one point with the tangent forming a 90-degree angle to the circle's radius.
	CLRTangent

	// CLROneEndOnCircumferenceOutside indicates that one endpoint of the line segment
	// lies on the circle's boundary, while the other endpoint lies outside the circle.
	CLROneEndOnCircumferenceOutside

	// CLROneEndOnCircumferenceInside indicates that one endpoint of the line segment
	// lies on the circle's boundary, while the other endpoint lies inside the circle.
	CLROneEndOnCircumferenceInside

	// CLRBothEndsOnCircumference indicates that both endpoints of the line segment lie
	// exactly on the circle's boundary. The line segment does not extend inside or outside
	// the circle.
	CLRBothEndsOnCircumference
)

// PointCircleRelationship (PCR) represents the spatial relationship of a point to a circle.
type PointCircleRelationship int

const (
	PCROutside         PointCircleRelationship = iota // Point is outside the circle
	PCROnCircumference                                // Point lies exactly on the circle's circumference
	PCRInside                                         // Point is inside the circle
)

// PointRectangleRelationship (PRR) describes the relationship of a point to a rectangle.
type PointRectangleRelationship int

const (
	PRROutside  PointRectangleRelationship = iota // The point lies outside the rectangle.
	PRRInside                                     // The point lies strictly inside the rectangle.
	PRROnVertex                                   // The point lies on a vertex of the rectangle.
	PRROnEdge                                     // The point lies on an edge of the rectangle.
)

// RectangleLineSegmentRelationship (RLR) describes the relationship between a rectangle and a line segment.
type RectangleLineSegmentRelationship int

const (
	RLROutside                 RectangleLineSegmentRelationship = iota // The segment lies entirely outside the rectangle.
	RLROutsideEndTouchesEdge                                           // The segment lies outside the rectangle and one end touches an edge.
	RLROutsideEndTouchesVertex                                         // The segment lies outside the rectangle and one end touches a vertex.
	RLRInside                                                          // The segment lies entirely within the rectangle.
	RLRInsideEndTouchesEdge                                            // The segment lies within the rectangle and one end touches an edge without crossing the boundary.
	RLRInsideEndTouchesVertex                                          // The segment lies within the rectangle and one end touches a vertex without crossing the boundary.
	RLROnEdge                                                          // The segment lies entirely on one edge of the rectangle.
	RLROnEdgeEndTouchesVertex                                          // The segment lies entirely on one edge of the rectangle and one or both ends touch a vertex.
	RLRIntersects                                                      // The segment crosses one or more edges of the rectangle.
	RLREntersAndExits                                                  // The segment enters through one edge and exits through another.
)

// ReflectionAxis specifies the axis or line across which a point or line segment should be reflected.
//
// This type defines the possible axes for reflection, including the standard x-axis and y-axis,
// as well as an arbitrary line defined by a custom line segment.
type ReflectionAxis int

const (
	// ReflectAcrossXAxis reflects a point or line segment across the x-axis, flipping the y-coordinate.
	ReflectAcrossXAxis ReflectionAxis = iota

	// ReflectAcrossYAxis reflects a point or line segment across the y-axis, flipping the x-coordinate.
	ReflectAcrossYAxis

	// ReflectAcrossCustomLine reflects a point or line segment across an arbitrary line defined by a LineSegment.
	// This line segment can be specified as an additional argument to the Reflect method.
	ReflectAcrossCustomLine
)

// SignedNumber is a generic interface representing signed numeric types supported by this package.
// This interface allows functions and structs to operate generically on various numeric types,
// including integer and floating-point types, while restricting to signed values only.
//
// Supported types:
//   - int
//   - int32
//   - int64
//   - float32
//   - float64
//
// By using SignedNumber, functions can handle multiple numeric types without needing to be rewritten
// for each specific type, enabling flexible and type-safe operations across different numeric data.
type SignedNumber interface {
	int | int32 | int64 | float32 | float64
}

// adjacentInSlice checks if two values are adjacent in the slice.
func adjacentInSlice[T comparable](s []T, a, b T) bool {
	var i, i2 int
	for i = 0; i < len(s); i++ {
		i2 = i + 1 // todo: (i+1)%len(s) and remove all the %len(s) below?
		if s[i%len(s)] == a && s[i2%len(s)] == b || s[i%len(s)] == b && s[i2%len(s)] == a {
			return true
		}
	}
	return false
}

// countOccurrencesInSlice returns the number of occurrences of an element in a slice.
func countOccurrencesInSlice[T comparable](s []T, element T) int {
	count := 0
	for _, v := range s {
		if v == element {
			count++
		}
	}
	return count
}

// inOrder returns true if b lies between a and c
func inOrder[T SignedNumber](a, b, c T) bool {
	return (a-b)*(b-c) > 0
}
