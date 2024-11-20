// Package geom2d provides foundational definitions, constants, and utility types for 2D geometry operations.
//
// This file contains core type definitions, global variables, constants, and utility functions
// that are shared across the package. It serves as the basis for other types like Point and LineSegment,
// offering foundational functionality and supporting mathematical constants and helper functions.

package geom2d

// CircleSegmentRelationship describes the spatial relationship between a line segment
// and a circle.
type CircleSegmentRelationship int

const (
	// CSROutside indicates that the line segment lies completely outside the circle,
	// with no intersection or tangency.
	CSROutside CircleSegmentRelationship = iota

	// CSRInside indicates that the line segment lies completely within the circle,
	// with both endpoints inside the circle's boundary.
	CSRInside

	// CSRIntersecting indicates that the line segment intersects the circle at exactly
	// two distinct points.
	CSRIntersecting

	// CSRTangent indicates that the line segment is tangent to the circle, touching it
	// at exactly one point with the tangent forming a 90-degree angle to the circle's radius.
	CSRTangent

	// CSROneEndOnCircumferenceOutside indicates that one endpoint of the line segment
	// lies on the circle's boundary, while the other endpoint lies outside the circle.
	CSROneEndOnCircumferenceOutside

	// CSROneEndOnCircumferenceInside indicates that one endpoint of the line segment
	// lies on the circle's boundary, while the other endpoint lies inside the circle.
	CSROneEndOnCircumferenceInside

	// CSRBothEndsOnCircumference indicates that both endpoints of the line segment lie
	// exactly on the circle's boundary. The line segment does not extend inside or outside
	// the circle.
	CSRBothEndsOnCircumference
)

// PointCircleRelationship represents the spatial relationship of a point to a circle.
type PointCircleRelationship int

const (
	Outside         PointCircleRelationship = iota // Point is outside the circle
	OnCircumference                                // Point lies exactly on the circle's circumference
	Inside                                         // Point is inside the circle
)

// PointRectangleRelationship describes the relationship of a point to a rectangle.
type PointRectangleRelationship int

const (
	PRROutside  PointRectangleRelationship = iota // The point lies outside the rectangle.
	PRRInside                                     // The point lies strictly inside the rectangle.
	PRROnVertex                                   // The point lies on a vertex of the rectangle.
	PRROnEdge                                     // The point lies on an edge of the rectangle.
)

// RectangleSegmentRelationship describes the relationship between a rectangle and a line segment.
type RectangleSegmentRelationship int

const (
	RSROutside                 RectangleSegmentRelationship = iota // The segment lies entirely outside the rectangle.
	RSROutsideEndTouchesEdge                                       // The segment lies outside the rectangle and one end touches an edge.
	RSROutsideEndTouchesVertex                                     // The segment lies outside the rectangle and one end touches a vertex.
	RSRInside                                                      // The segment lies entirely within the rectangle.
	RSRInsideEndTouchesEdge                                        // The segment lies within the rectangle and one end touches an edge without crossing the boundary.
	RSRInsideEndTouchesVertex                                      // The segment lies within the rectangle and one end touches a vertex without crossing the boundary.
	RSROnEdge                                                      // The segment lies entirely on one edge of the rectangle.
	RSROnEdgeEndTouchesVertex                                      // The segment lies entirely on one edge of the rectangle and one or both ends touch a vertex.
	RSRIntersects                                                  // The segment crosses one or more edges of the rectangle.
	RSREntersAndExits                                              // The segment enters through one edge and exits through another.
)

// ReflectionAxis specifies the axis or line across which a point or line segment should be reflected.
//
// This type defines the possible axes for reflection, including the standard x-axis and y-axis,
// as well as an arbitrary line defined by a custom line segment.
type ReflectionAxis int

const (
	// XAxis reflects a point or line segment across the x-axis, flipping the y-coordinate.
	XAxis ReflectionAxis = iota

	// YAxis reflects a point or line segment across the y-axis, flipping the x-coordinate.
	YAxis

	// CustomLine reflects a point or line segment across an arbitrary line defined by a LineSegment.
	// This line segment can be specified as an additional argument to the Reflect method.
	CustomLine
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

// TwoLinesRelationship defines the possible spatial relationships between two line segments, AB and CD.
//
// Where:
//   - LineSegment AB starts at Point A and ends at Point End
//   - LineSegment CD starts at Point C and ends at Point D
//
// This type is used by functions that analyze the geometric relationship between two line segments
// to provide precise information about their interaction, such as intersection, overlap, or disjointedness.
type TwoLinesRelationship int8

// Valid values for TwoLinesRelationship:
//
// TwoLinesRelationship represents possible spatial relationships between two line segments, AB and CD.
//
// # LSR stands for Line Segment Relationship
//
// A positive value indicates that the segments have some form of relationship (e.g., intersection,
// overlap, or containment). A non-positive value (either 0 or -1) indicates no interaction between
// the segments.
//
// Note: `LSRCollinearDisjoint` is intentionally assigned a value of -1. This allows for a simple check
// to determine whether the segments have any interaction by evaluating if `TwoLinesRelationship > 0`.
// `LSRCollinearDisjoint` specifically represents the case where the segments are collinear but
// do not overlap or touch at any point, making it distinct from `LSRMiss`, which represents disjoint,
// non-collinear segments.
const (
	LSRCollinearDisjoint TwoLinesRelationship = iota - 1 // Segments are collinear and do not intersect, overlap, or touch at any point.
	LSRMiss                                              // The segments are not collinear, disjoint and do not intersect, overlap, or touch at any point.
	LSRIntersects                                        // The segments intersect at a unique point that is not an endpoint.
	LSRAeqC                                              // Point A of segment AB coincides with Point C of segment CD
	LSRAeqD                                              // Point A of segment AB coincides with Point D of segment CD
	LSRBeqC                                              // Point End of segment AB coincides with Point C of segment CD
	LSRBeqD                                              // Point End of segment AB coincides with Point D of segment CD
	LSRAonCD                                             // Point A lies on LineSegment CD
	LSRBonCD                                             // Point End lies on LineSegment CD
	LSRConAB                                             // Point C lies on LineSegment AB
	LSRDonAB                                             // Point D lies on LineSegment AB
	LSRCollinearAonCD                                    // Point A lies on LineSegment CD (partial overlap), and line segments are collinear
	LSRCollinearBonCD                                    // Point End lies on LineSegment CD (partial overlap), and line segments are collinear
	LSRCollinearABinCD                                   // Segment AB is fully contained within segment CD
	LSRCollinearCDinAB                                   // Segment CD is fully contained within segment AB
	LSRCollinearEqual                                    // The segments AB and CD are exactly equal, sharing both endpoints in the same locations.
)

// adjoiningEdges checks if two values are adjacent in the slice.
func adjoiningEdges[T comparable](s []T, a, b T) bool {
	var i, i2 int
	for i = 0; i < len(s); i++ {
		i2 = i + 1
		if s[i%len(s)] == a && s[i2%len(s)] == b || s[i%len(s)] == b && s[i2%len(s)] == a {
			return true
		}
	}
	return false
}

// countOccurrences returns the number of occurrences of an element in a slice.
func countOccurrences[T comparable](s []T, element T) int {
	count := 0
	for _, v := range s {
		if v == element {
			count++
		}
	}
	return count
}
