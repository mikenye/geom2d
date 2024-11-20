// linesegment.go defines the LineSegment type and implements methods for operations on line segments in 2D space.
//
// This file includes the definition of the LineSegment type, with methods for calculations and relationships
// between line segments, including distance calculations, midpoint computations, and intersections. Functions
// that primarily operate on line segments are centralized in this file.

package geom2d

import (
	"fmt"
	"math"
)

// LineSegmentsRelationship (LSR) defines the possible spatial relationships between two line segments, AB and CD.
//
// Where:
//   - LineSegment AB starts at Point A and ends at Point End
//   - LineSegment CD starts at Point C and ends at Point D
//
// This type is used by functions that analyze the geometric relationship between two line segments
// to provide precise information about their interaction, such as intersection, overlap, or disjointedness.
type LineSegmentsRelationship int8

// Valid values for LineSegmentsRelationship:
//
// LineSegmentsRelationship represents possible spatial relationships between two line segments, AB and CD.
//
// # LSR stands for Line Segment Relationship
//
// A positive value indicates that the segments have some form of relationship (e.g., intersection,
// overlap, or containment). A non-positive value (either 0 or -1) indicates no interaction between
// the segments.
//
// Note: `LSRCollinearDisjoint` is intentionally assigned a value of -1. This allows for a simple check
// to determine whether the segments have any interaction by evaluating if `LineSegmentsRelationship > 0`.
// `LSRCollinearDisjoint` specifically represents the case where the segments are collinear but
// do not overlap or touch at any point, making it distinct from `LSRMiss`, which represents disjoint,
// non-collinear segments.
const (
	LSRCollinearDisjoint LineSegmentsRelationship = iota - 1 // Segments are collinear and do not intersect, overlap, or touch at any point.
	LSRMiss                                                  // The segments are not collinear, disjoint and do not intersect, overlap, or touch at any point.
	LSRIntersects                                            // The segments intersect at a unique point that is not an endpoint.
	LSRAeqC                                                  // Point A of segment AB coincides with Point C of segment CD
	LSRAeqD                                                  // Point A of segment AB coincides with Point D of segment CD
	LSRBeqC                                                  // Point End of segment AB coincides with Point C of segment CD
	LSRBeqD                                                  // Point End of segment AB coincides with Point D of segment CD
	LSRAonCD                                                 // Point A lies on LineSegment CD
	LSRBonCD                                                 // Point End lies on LineSegment CD
	LSRConAB                                                 // Point C lies on LineSegment AB
	LSRDonAB                                                 // Point D lies on LineSegment AB
	LSRCollinearAonCD                                        // Point A lies on LineSegment CD (partial overlap), and line segments are collinear
	LSRCollinearBonCD                                        // Point End lies on LineSegment CD (partial overlap), and line segments are collinear
	LSRCollinearABinCD                                       // Segment AB is fully contained within segment CD
	LSRCollinearCDinAB                                       // Segment CD is fully contained within segment AB
	LSRCollinearEqual                                        // The segments AB and CD are exactly equal, sharing both endpoints in the same locations.
)

// ScaleOrigin specifies the origin point from which a line segment should be scaled.
//
// The origin can be set to either the start point, the end point, or the midpoint
// of the line segment. This allows for flexible scaling behavior depending on
// the desired point of reference.
type ScaleOrigin uint8

const (
	// ScaleFromStart scales the line segment from its start point.
	// The start point remains fixed, while the end point is adjusted based on the scaling factor.
	ScaleFromStart ScaleOrigin = iota

	// ScaleFromEnd scales the line segment from its end point.
	// The end point remains fixed, while the start point is adjusted based on the scaling factor.
	ScaleFromEnd

	// ScaleFromMidpoint scales the line segment from its midpoint.
	// Both the start and end points are adjusted proportionally to maintain the midpoint's position.
	ScaleFromMidpoint
)

// LineSegment represents a line segment in a 2D space, defined by two endpoints, A and End.
//
// The generic type parameter T must satisfy the SignedNumber constraint, allowing the segment
// to use various numeric types such as int or float64 for its coordinates.
//
// Fields:
//   - start: Point[T] - The starting point of the line segment.
//   - end: Point[T] - The ending point of the line segment.
type LineSegment[T SignedNumber] struct {
	start Point[T]
	end   Point[T]
}

// AddLineSegment adds the start and end points of another line segment to this one.
//
// This method performs an element-wise addition, where the `start` and `end` points
// of the `other` line segment are added to the corresponding `start` and `end` points
// of the current line segment.
//
// Parameters:
//   - CD: LineSegment[T] - The line segment to add to the current one.
//
// Returns:
//   - LineSegment[T] - A new line segment where each endpoint is the sum of the corresponding endpoints.
//
// Example usage:
//
//	segment1 := NewLineSegment(NewPoint(1, 1), NewPoint(4, 5))
//	segment2 := NewLineSegment(NewPoint(2, 3), NewPoint(1, 2))
//	result := segment1.AddLineSegment(segment2)
//	// `result` will be a new segment from (3, 4) to (5, 7)
func (AB LineSegment[T]) AddLineSegment(CD LineSegment[T]) LineSegment[T] {
	return NewLineSegment(
		AB.start.Add(CD.start),
		AB.end.Add(CD.end),
	)
}

// AddVector translates the line segment by adding a given vector to both the start and end points.
//
// This method moves the line segment by the specified vector, effectively shifting
// both endpoints by the vector’s coordinates. It’s useful for translating the entire
// segment in a given direction.
//
// Parameters:
//   - v: Point[T] - The vector to add to both endpoints of the segment.
//
// Returns:
//   - LineSegment[T] - A new line segment translated by the given vector.
//
// Example usage:
//
//	segment := NewLineSegment(NewPoint(1, 1), NewPoint(4, 5))
//	translated := segment.AddVector(NewPoint(2, 3))
//	// `translated` will be a new segment from (3, 4) to (6, 8)
func (AB LineSegment[T]) AddVector(v Point[T]) LineSegment[T] {
	return NewLineSegment(
		AB.start.Add(v),
		AB.end.Add(v),
	)
}

// AsFloat converts the line segment to a LineSegment[float64] type.
//
// This function converts both endpoints of the line segment `AB` to `Point[float64]`
// values, creating a new line segment with floating-point coordinates.
// It is useful for precise calculations where floating-point accuracy is needed.
//
// Returns:
//   - LineSegment[float64] - The line segment with both endpoints converted to float64.
func (AB LineSegment[T]) AsFloat() LineSegment[float64] {
	return NewLineSegment(AB.start.AsFloat(), AB.end.AsFloat())
}

// AsInt converts the line segment to a LineSegment[int] type.
//
// This function converts both endpoints of the line segment `AB` to `Point[int]`
// by truncating any decimal places. It is useful for converting a floating-point
// line segment to integer coordinates without rounding.
//
// Returns:
//   - LineSegment[int] - The line segment with both endpoints converted to integer coordinates by truncation.
func (AB LineSegment[T]) AsInt() LineSegment[int] {
	return NewLineSegment(AB.start.AsInt(), AB.end.AsInt())
}

// AsIntRounded converts the line segment to a LineSegment[int] type with rounded coordinates.
//
// This function converts both endpoints of the line segment `AB` to `Point[int]`
// by rounding each coordinate to the nearest integer. It is useful when you need to
// approximate the segment’s position with integer coordinates while minimizing the
// rounding error.
//
// Returns:
//   - LineSegment[int] - The line segment with both endpoints converted to integer coordinates by rounding.
func (AB LineSegment[T]) AsIntRounded() LineSegment[int] {
	return NewLineSegment(AB.start.AsIntRounded(), AB.end.AsIntRounded())
}

// DistanceToLineSegment calculates the minimum distance between two line segments, AB and CD.
//
// If the segments intersect or touch at any point, the function returns 0, as the distance is effectively zero.
// Otherwise, it calculates the shortest distance by considering distances between segment endpoints and the
// perpendicular projections of each endpoint onto the opposite segment.
//
// Parameters:
//   - CD: LineSegment[T] - The line segment to compare with `AB`.
//
// Returns:
//   - float64 - The minimum distance between the two segments. If the segments intersect or touch, this value is 0.
//
// The function converts the segments to float64 for precise calculations, checks the intersection status,
// and computes distances in three main parts:
//  1. Direct distances between endpoints of `AB` and `CD`.
//  2. Perpendicular projections of each endpoint onto the opposite segment.
//  3. Tracks the minimum distance among all calculations and returns this value.
//
// Example usage:
//
//	segmentAB := NewLineSegment(NewPoint(0, 0), NewPoint(2, 2))
//	segmentCD := NewLineSegment(NewPoint(3, 3), NewPoint(5, 5))
//	distance := segmentAB.DistanceToLineSegment(segmentCD)
//
// `distance` will hold the minimum distance between `AB` and `CD`.
func (AB LineSegment[T]) DistanceToLineSegment(CD LineSegment[T]) float64 {
	// If line segments intersect, the distance is zero.
	if AB.IntersectsLineSegment(CD) {
		return 0
	}

	// Convert segments to float for precise calculations.
	ABf, CDf := AB.AsFloat(), CD.AsFloat()

	// Track the minimum distance.
	minDist := math.MaxFloat64

	// Helper function to update minimum distance.
	updateMinDist := func(d float64) {
		if d < minDist {
			minDist = d
		}
	}

	// Calculate distances between endpoints.
	updateMinDist(ABf.start.DistanceToPoint(CDf.start))
	updateMinDist(ABf.start.DistanceToPoint(CDf.end))
	updateMinDist(ABf.end.DistanceToPoint(CDf.start))
	updateMinDist(ABf.end.DistanceToPoint(CDf.end))

	// Calculate distances to projections on the opposite segment.
	updateMinDist(ABf.start.DistanceToPoint(ABf.start.ProjectOntoLineSegment(CDf)))
	updateMinDist(ABf.end.DistanceToPoint(ABf.end.ProjectOntoLineSegment(CDf)))
	updateMinDist(CDf.start.DistanceToPoint(CDf.start.ProjectOntoLineSegment(ABf)))
	updateMinDist(CDf.end.DistanceToPoint(CDf.end.ProjectOntoLineSegment(ABf)))

	return minDist
}

// DistanceToPoint calculates the orthogonal distance from the specified Point p
// to the LineSegment AB. This distance is determined by projecting
// the Point p onto the LineSegment and measuring the distance from p to
// the projected Point.
//
// Returns the distance as a float64.
//
// The function first computes the projection of p onto the line segment defined
// by the endpoints of l. It then converts the original point p to a
// Point[float64] to ensure accurate distance calculation, as the
// DistanceToPoint function operates on float64 points.
func (AB LineSegment[T]) DistanceToPoint(p Point[T]) float64 {
	return p.DistanceToLineSegment(AB)
}

// End returns the ending point of the line segment.
//
// This function provides access to the ending point of the line segment `AB`, typically representing
// the endpoint of the segment.
func (AB LineSegment[T]) End() Point[T] {
	return AB.end
}

// Eq checks if two line segments are equal by comparing their start and end points.
//
// This method returns true if both the `start` and `end` points of the current line segment
// (`AB`) are equal to the corresponding `start` and `end` points of the `CD` line segment.
// Equality is determined based on the `Eq` method of the `Point` type.
//
// Parameters:
//   - CD: LineSegment[T] - The line segment to compare with the current line segment.
//
// Returns:
//   - bool - Returns `true` if both line segments have identical start and end points, and `false` otherwise.
//
// Example usage:
//
//	segment1 := NewLineSegment(NewPoint(1, 1), NewPoint(4, 5))
//	segment2 := NewLineSegment(NewPoint(1, 1), NewPoint(4, 5))
//	equal := segment1.Eq(segment2) // equal will be true
func (AB LineSegment[T]) Eq(CD LineSegment[T]) bool {
	return AB.start.Eq(CD.start) && AB.end.Eq(CD.end)
}

// IntersectsLineSegment checks whether there is any intersection or overlap between LineSegment AB and LineSegment CD.
//
// This function returns true if segments `AB` and `CD` have an intersecting spatial relationship, such as intersection,
// overlap, containment, or endpoint coincidence. It leverages the `RelationshipToLineSegment` function to
// determine if the relationship value is greater than `LSRMiss`, indicating that the segments are not fully
// disjoint.
//
// Parameters:
//   - CD: LineSegment[T] - The line segment to compare with `AB`.
//
// Returns:
//   - bool - Returns `true` if `AB` and `CD` intersect, overlap, or share any form of intersecting relationship, and `false` if they are completely disjoint.
//
// Example usage:
//
//	segmentAB := NewLineSegment(NewPoint(0, 0), NewPoint(2, 2))
//	segmentCD := NewLineSegment(NewPoint(1, 1), NewPoint(3, 3))
//	intersects := segmentAB.IntersectsLineSegment(segmentCD)
//
// `intersects` will be `true` as there is an intersecting relationship between `AB` and `CD`.
func (AB LineSegment[T]) IntersectsLineSegment(CD LineSegment[T]) bool {
	if AB.RelationshipToLineSegment(CD) > LSRMiss {
		return true
	}
	return false
}

// IntersectionPoint calculates the intersection point between two line segments, if one exists.
//
// This method checks if the LineSegment AB and LineSegment CD intersect within their boundaries
// and, if so, calculates and returns the intersection point. If the segments do not intersect
// or are parallel, it returns a zero-value Point and false.
//
// It uses the parametric form of the line segments to solve for intersection parameters `t` and `u`.
// If `t` and `u` are both in the range [0, 1], the intersection point lies within the bounds of
// both segments.
//
// Parameters:
// - CD (LineSegment[T]): The second line segment to check for an intersection.
//
// Returns:
//   - Point[T]: The intersection point.
//   - bool: If `true`, the first element is the intersection point. If `false`, there is
//     no intersection within the segments’ bounds or the segments are parallel.
//
// Usage:
//
//	intersection, intersects := seg1.IntersectionPoint(seg2)
//	if intersects {
//	    // Process the intersection point
//	} else {
//	    // Handle non-intersecting segments
//	}
func (AB LineSegment[T]) IntersectionPoint(CD LineSegment[T]) (Point[float64], bool) {
	// Define segment endpoints for AB and CD
	A, B := AB.start.AsFloat(), AB.end.AsFloat()
	C, D := CD.start.AsFloat(), CD.end.AsFloat()

	// Calculate the direction vectors

	dir1 := B.Sub(A)
	dir2 := D.Sub(C)

	// Calculate the determinants
	denominator := dir1.CrossProduct(dir2)

	// Check if the lines are parallel (no intersection)
	if denominator == 0 {
		return Point[float64]{}, false // No intersection
	}

	// Calculate parameters t and u
	AC := C.Sub(A)
	tNumerator := AC.CrossProduct(dir2)
	uNumerator := AC.CrossProduct(dir1)

	t := tNumerator / denominator
	u := uNumerator / denominator

	// Check if intersection occurs within the segment bounds
	if t < 0 || t > 1 || u < 0 || u > 1 {
		return Point[float64]{}, false // Intersection is outside the segments
	}

	// Calculate the intersection point
	intersection := A.Add(dir1.Scale(t))
	return intersection, true
}

// Length calculates the length of the line segment.
//
// This function computes the distance between the starting point and the ending point of the line segment `AB`,
// using the Euclidean distance formula.
//
// Returns:
//   - float64 - The length of the line segment.
func (AB LineSegment[T]) Length() float64 {
	return AB.start.DistanceToPoint(AB.end)
}

// Midpoint calculates the midpoint of the line segment.
//
// This function finds the midpoint by averaging the x and y coordinates of the start and end points.
// The result is returned as a `Point[float64]`, regardless of the segment’s original coordinate type,
// ensuring precision.
//
// Returns:
//   - Point[float64] - The midpoint of the line segment as a point with floating-point coordinates.
func (AB LineSegment[T]) Midpoint() Point[float64] {
	start := AB.start.AsFloat()
	end := AB.end.AsFloat()
	return NewPoint[float64](
		(start.x+end.x)/2,
		(start.y+end.y)/2,
	)
}

// Points returns the two endpoints of the line segment as a slice of Points.
// The order of the points is [start, end].
//
// Returns:
//   - []Point[T]: A slice containing the start and end points of the line segment.
func (AB LineSegment[T]) Points() []Point[T] {
	return []Point[T]{AB.start, AB.end}
}

// Reflect reflects the line segment across the specified axis or custom line.
//
// Parameters:
//   - axis: ReflectionAxis - The axis or line to reflect across (ReflectAcrossXAxis, ReflectAcrossYAxis, or ReflectAcrossCustomLine).
//   - line: Optional LineSegment[float64] - The line segment for ReflectAcrossCustomLine reflection.
//
// Returns:
//   - LineSegment[float64] - A new line segment where both endpoints are reflected accordingly.
//
// Example usage:
//
//	segment := NewLineSegment[float64](NewPoint[float64](2, 3), NewPoint[float64](4, 5))
//	reflected := segment.Reflect(ReflectAcrossXAxis) // Reflects both points across the x-axis.
func (AB LineSegment[float64]) Reflect(axis ReflectionAxis, line ...LineSegment[float64]) LineSegment[float64] {
	var startReflected, endReflected Point[float64]
	switch axis {
	case ReflectAcrossXAxis:
		// Reflect across the x-axis
		startReflected = AB.start.Reflect(ReflectAcrossXAxis)
		endReflected = AB.end.Reflect(ReflectAcrossXAxis)
	case ReflectAcrossYAxis:
		// Reflect across the y-axis
		startReflected = AB.start.Reflect(ReflectAcrossYAxis)
		endReflected = AB.end.Reflect(ReflectAcrossYAxis)
	case ReflectAcrossCustomLine:
		// Reflect across a custom line if provided
		if len(line) > 0 {
			startReflected = AB.start.Reflect(ReflectAcrossCustomLine, line[0])
			endReflected = AB.end.Reflect(ReflectAcrossCustomLine, line[0])
		} else {
			// No custom line provided; return the original line segment unchanged
			return AB
		}
	default:
		// Invalid axis, return the line segment unchanged
		return AB
	}

	// Return a new line segment with reflected points
	return NewLineSegment(startReflected, endReflected)
}

// RelationshipToLineSegment determines the spatial relationship between two line segments, AB and CD.
//
// This function evaluates the relationship between two line segments, `AB` and `CD`, by checking for
// endpoint coincidences, intersections, collinear relationships, and containment. It returns a
// `LineSegmentsRelationship` constant that describes the exact relationship between the segments, such
// as intersection, partial overlap, or full containment.
//
// Parameters:
//   - CD: LineSegment[T] - The line segment to compare with `AB`.
//
// Returns:
//   - LineSegmentsRelationship - A constant that describes the relationship between segments `AB` and `CD`.
//
// Possible return values:
//   - LSRCollinearDisjoint: The segments are collinear but do not overlap or touch at any point.
//   - LSRMiss: The segments are not collinear, disjoint, and do not intersect, overlap, or touch at any point.
//   - LSRIntersects: The segments intersect at a unique point that is not an endpoint.
//   - LSRAeqC, LSRAeqD, LSRBeqC, LSRBeqD: An endpoint of `AB` coincides with an endpoint of `CD`. For example,
//     `LSRAeqC` indicates that point A of `AB` coincides with point C of `CD`.
//   - LSRAonCD, LSRBonCD, LSRConAB, LSRDonAB: One endpoint of one segment lies on the other segment without
//     the segments being collinear. For example, `LSRAonCD` indicates that point A of `AB` lies on segment `CD`.
//   - LSRCollinearAonCD, LSRCollinearBonCD: The segments are collinear
//     with partial overlap where one endpoint of one segment lies on the other. For example, `LSRCollinearAonCD`
//     means point A of `AB` lies on `CD` with collinearity.
//   - LSRCollinearABinCD: The entire segment `AB` is contained within segment `CD`.
//   - LSRCollinearCDinAB: The entire segment `CD` is contained within segment `AB`.
//   - LSRCollinearEqual: The segments `AB` and `CD` are exactly equal, sharing both endpoints in the same locations.
//
// Example usage:
//
//	segmentAB := NewLineSegment(NewPoint(0, 0), NewPoint(2, 2))
//	segmentCD := NewLineSegment(NewPoint(1, 1), NewPoint(3, 3))
//	relationship := segmentAB.RelationshipToLineSegment(segmentCD)
//
// The variable `relationship` should equal `LSRCollinearAonCD` as:
//   - `AB` and `CD` are Collinear
//   - `A` lies on `CD`
//   - `AB` and `CD` don't fully overlap.
func (AB LineSegment[T]) RelationshipToLineSegment(CD LineSegment[T]) LineSegmentsRelationship {

	// Check if segments are exactly equal
	if (AB.start.Eq(CD.start) && AB.end.Eq(CD.end)) || (AB.start.Eq(CD.end) && AB.end.Eq(CD.start)) {
		return LSRCollinearEqual
	}

	switch {

	// Check if A and C coincide
	case AB.start.Eq(CD.start):
		return LSRAeqC

	// Check if A and D coincide
	case AB.start.Eq(CD.end):
		return LSRAeqD

	// Check if End and C coincide
	case AB.end.Eq(CD.start):
		return LSRBeqC

	// Check if End and D coincide
	case AB.end.Eq(CD.end):
		return LSRBeqD

	}

	// Determine orientations for intersection and collinearity checks
	o1 := Orientation(AB.start, AB.end, CD.start)
	o2 := Orientation(AB.start, AB.end, CD.end)
	o3 := Orientation(CD.start, CD.end, AB.start)
	o4 := Orientation(CD.start, CD.end, AB.end)

	// Non-collinear intersection cases
	if o1 != o2 && o3 != o4 {

		switch {

		// Check if A lies on CD
		case AB.start.IsOnLineSegment(CD) && !AB.end.IsOnLineSegment(CD):
			return LSRAonCD

		// Check if End lies on CD
		case !AB.start.IsOnLineSegment(CD) && AB.end.IsOnLineSegment(CD):
			return LSRBonCD

		// Check if C lies on AB
		case CD.start.IsOnLineSegment(AB) && !CD.end.IsOnLineSegment(AB):
			return LSRConAB

		// Check if D lies on AB
		case !CD.start.IsOnLineSegment(AB) && CD.end.IsOnLineSegment(AB):
			return LSRDonAB

		// Default case that lines intersect without any "edge cases"
		default:
			return LSRIntersects
		}
	}

	// PointsCollinear cases: All orientations are zero
	if o1 == 0 && o2 == 0 && o3 == 0 && o4 == 0 {
		// Check if segments are collinear and disjoint
		if !AB.start.IsOnLineSegment(CD) && !AB.end.IsOnLineSegment(CD) &&
			!CD.start.IsOnLineSegment(AB) && !CD.end.IsOnLineSegment(AB) {
			return LSRCollinearDisjoint
		}
		// Check if AB is fully contained within CD
		if AB.start.IsOnLineSegment(CD) && AB.end.IsOnLineSegment(CD) {
			return LSRCollinearABinCD
		}
		// Check if CD is fully contained within AB
		if CD.start.IsOnLineSegment(AB) && CD.end.IsOnLineSegment(AB) {
			return LSRCollinearCDinAB
		}
		// Check specific collinear partial overlaps
		if AB.start.IsOnLineSegment(CD) {
			return LSRCollinearAonCD
		}
		if AB.end.IsOnLineSegment(CD) {
			return LSRCollinearBonCD
		}
	}

	// If none of the conditions matched, the segments are disjoint
	return LSRMiss
}

// Rotate rotates the LineSegment around a given pivot point by a specified angle in radians.
//
// Parameters:
//   - radians: float64 - The rotation angle in radians.
//   - pivot: Point[T] - The point around which to rotate the line segment.
//
// Returns:
//   - LineSegment[float64] - A new line segment representing the rotated position.
func (AB LineSegment[T]) Rotate(pivot Point[T], radians float64) LineSegment[float64] {
	newStart := AB.start.Rotate(pivot, radians)
	newEnd := AB.end.Rotate(pivot, radians)
	return NewLineSegment(newStart, newEnd)
}

// Scale scales the line segment by a given factor from a specified origin point.
//
// Parameters:
//   - origin: ScaleOrigin - The point from which to scale (start, end, or midpoint).
//   - factor: float64 - The scaling factor.
//
// Returns:
//   - LineSegment[T] - A new line segment scaled relative to the specified origin.
func (AB LineSegment[T]) Scale(origin ScaleOrigin, factor float64) LineSegment[float64] {
	var refPoint Point[float64]
	switch origin {
	case ScaleFromStart:
		refPoint = AB.start.AsFloat()
	case ScaleFromEnd:
		refPoint = AB.end.AsFloat()
	case ScaleFromMidpoint:
		refPoint = AB.Midpoint()
	}

	return NewLineSegment(
		AB.start.AsFloat().ScaleFrom(refPoint, factor),
		AB.end.AsFloat().ScaleFrom(refPoint, factor),
	)
}

// Start returns the starting point of the line segment.
//
// This function provides access to the starting point of the line segment `AB`, typically representing
// the beginning of the segment.
func (AB LineSegment[T]) Start() Point[T] {
	return AB.start
}

// String returns a formatted string representation of the line segment for debugging and logging purposes.
//
// The string representation includes the coordinates of the `start` and `end` points in the format:
// "LineSegment[(x1, y1) -> (x2, y2)]", where (x1, y1) are the coordinates of the `start` point,
// and (x2, y2) are the coordinates of the `end` point.
//
// Returns:
//   - string - A string representing the line segment's `start` and `end` coordinates.
//
// Example usage:
//
//	segment := NewLineSegment(NewPoint(1, 1), NewPoint(4, 5))
//	fmt.Println(segment.String()) // Output: "LineSegment[(1, 1) -> (4, 5)]"
func (AB LineSegment[T]) String() string {
	return fmt.Sprintf("LineSegment[(%v, %v) -> (%v, %v)]", AB.start.x, AB.start.y, AB.end.x, AB.end.y)
}

// SubLineSegment subtracts the start and end points of another line segment from this one.
//
// This function performs an element-wise subtraction, where the `start` and `end` points
// of the `other` line segment are subtracted from the corresponding `start` and `end` points
// of the current line segment.
//
// Parameters:
//   - CD: LineSegment[T] - The line segment to subtract from the current one.
//
// Returns:
//   - LineSegment[T] - A new line segment where each endpoint is the result of the element-wise subtraction.
func (AB LineSegment[T]) SubLineSegment(CD LineSegment[T]) LineSegment[T] {
	return NewLineSegment(
		AB.start.Sub(CD.start),
		AB.end.Sub(CD.end),
	)
}

// SubVector translates the line segment by subtracting a given vector from both the start and end points.
//
// This function moves the line segment by the inverse of the vector, effectively shifting
// both endpoints in the opposite direction of the vector.
//
// Parameters:
//   - v: Point[T] - The vector to subtract from both endpoints of the segment.
//
// Returns:
//   - LineSegment[T] - A new line segment translated by the inverse of the given vector.
func (AB LineSegment[T]) SubVector(v Point[T]) LineSegment[T] {
	return NewLineSegment(
		AB.start.Sub(v),
		AB.end.Sub(v),
	)
}

// NewLineSegment creates a new line segment from two endpoints, `start` and `end`.
//
// This constructor function initializes a `LineSegment` with the specified starting and ending points.
// The generic type parameter `T` must satisfy the `SignedNumber` constraint, allowing various numeric types
// (such as `int` or `float64`) to be used for the segment’s coordinates.
//
// Parameters:
//   - start: Point[T] - The starting point of the line segment.
//   - end: Point[T] - The ending point of the line segment.
//
// Returns:
//   - LineSegment[T] - A new line segment defined by the `start` and `end` points.
//
// Example usage:
//
//	segment := NewLineSegment(NewPoint(0, 0), NewPoint(3, 4))
//
// Creates a line segment with generic type parameter int, from (0,0) to (3,4).
func NewLineSegment[T SignedNumber](start, end Point[T]) LineSegment[T] {
	return LineSegment[T]{
		start: start,
		end:   end,
	}
}
