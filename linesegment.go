package geom2d

import (
	"fmt"
	"math"
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

// LineSegment represents a line segment in a 2D space, defined by two endpoints, the start [Point] and end [Point].
//
// The generic type parameter T must satisfy the [SignedNumber] constraint, allowing the segment
// to use various numeric types such as int or float64 for its coordinates.
type LineSegment[T SignedNumber] struct {
	start Point[T]
	end   Point[T]
}

// NewLineSegment creates a new line segment from two endpoints, a start [Point] and end [Point].
//
// This constructor function initializes a LineSegment with the specified starting and ending points.
// The generic type parameter `T` must satisfy the [SignedNumber] constraint, allowing various numeric types
// (such as `int` or `float64`) to be used for the segment’s coordinates.
//
// Parameters:
//   - start ([Point][T]): The starting [Point] of the LineSegment.
//   - end ([Point][T]): The ending [Point] of the LineSegment.
//
// Returns:
//   - LineSegment[T] - A new line segment defined by the start and end points.
func NewLineSegment[T SignedNumber](start, end Point[T]) LineSegment[T] {
	return LineSegment[T]{
		start: start,
		end:   end,
	}
}

// AddLineSegment adds the start and end points of another line segment to this one.
//
// This method performs an element-wise addition, where the start and end points
// of the other line segment are added to the corresponding start and end points
// of the current line segment.
//
// Parameters:
//   - CD: LineSegment[T] - The line segment to add to the current one.
//
// Returns:
//   - LineSegment[T] - A new line segment where each endpoint is the sum of the corresponding endpoints.
func (AB LineSegment[T]) AddLineSegment(CD LineSegment[T]) LineSegment[T] {
	return NewLineSegment(
		AB.start.Translate(CD.start),
		AB.end.Translate(CD.end),
	)
}

// Area returns the area of the line segment, which is always 0.
//
// This is because a line segment is a one-dimensional geometric entity
// and does not enclose any space in two dimensions.
//
// Note:
//   - This method exists to satisfy the [Measurable] interface. It is not here
//     to insult your intelligence. Rest assured, we know you understand that
//     line segments don't have area. 😊
//
// Returns:
//   - float64: The area of the line segment, which is 0.
func (AB LineSegment[T]) Area() float64 {
	return 0
}

// AsFloat converts the line segment to a LineSegment[float64] type.
//
// This function converts both endpoints of the LineSegment AB to [Point][float64]
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
// This function converts both endpoints of the line segment AB to [Point][int]
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
// This function converts both endpoints of the line segment AB to [Point][int]
// by rounding each coordinate to the nearest integer. It is useful when you need to
// approximate the segment’s position with integer coordinates while minimizing the
// rounding error.
//
// Returns:
//   - LineSegment[int] - The line segment with both endpoints converted to integer coordinates by rounding.
func (AB LineSegment[T]) AsIntRounded() LineSegment[int] {
	return NewLineSegment(AB.start.AsIntRounded(), AB.end.AsIntRounded())
}

// BoundingBox computes the smallest axis-aligned rectangle that fully contains the LineSegment.
//
// Returns:
//   - [Rectangle][T]: A rectangle defined by the opposite corners of the LineSegment.
//
// Behavior:
//   - The rectangle's top-left corner corresponds to the minimum x and y coordinates
//     of the LineSegment's start and end points.
//   - The rectangle's bottom-right corner corresponds to the maximum x and y coordinates
//     of the LineSegment's start and end points.
//
// Notes:
//   - This method is useful for spatial queries, collision detection, or visual rendering.
func (AB LineSegment[T]) BoundingBox() Rectangle[T] {
	points := []Point[T]{
		AB.start,
		NewPoint(AB.start.x, AB.end.y),
		NewPoint(AB.end.x, AB.start.y),
		AB.end,
	}
	return NewRectangle(points)
}

// Center calculates the midpoint of the line segment, optionally applying an epsilon
// threshold to adjust the precision of the result.
//
// Parameters:
//   - opts: A variadic slice of [Option] functions to customize the behavior of the calculation.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for snapping near-integer or
//     near-zero results to cleaner values, improving robustness in floating-point calculations.
//
// Behavior:
//   - The midpoint is calculated by averaging the x and y coordinates of the start and end
//     points of the line segment.
//   - If [WithEpsilon] is provided, the resulting midpoint coordinates are adjusted such that
//     small deviations due to floating-point precision errors are corrected.
//
// Returns:
//   - [Point][float64]: The midpoint of the line segment as a point with floating-point coordinates,
//     optionally adjusted based on epsilon.
//
// Notes:
//   - Epsilon adjustment is particularly useful when working with floating-point coordinates
//     where minor imprecision could affect the midpoint calculation.
//   - The midpoint is always returned as [Point][float64], ensuring precision regardless of the
//     coordinate type of the original line segment.
func (AB LineSegment[T]) Center(opts ...Option) Point[float64] {
	// Apply geomOptions with defaults
	options := applyOptions(geomOptions{epsilon: 0}, opts...)

	start := AB.start.AsFloat()
	end := AB.end.AsFloat()

	midX := (start.x + end.x) / 2
	midY := (start.y + end.y) / 2

	// Apply epsilon if specified
	if options.epsilon > 0 {
		midX = applyEpsilon(midX, options.epsilon)
		midY = applyEpsilon(midY, options.epsilon)
	}

	return NewPoint[float64](midX, midY)
}

func (AB LineSegment[T]) ContainsPoint(p Point[T]) bool {
	return p.IsOnLineSegment(AB)
}

// DistanceToLineSegment calculates the minimum distance between two line segments, AB and CD.
//
// If the segments intersect or touch at any point, the function returns 0, as the distance is effectively zero.
// Otherwise, it calculates the shortest distance by considering distances between segment endpoints and the
// perpendicular projections of each endpoint onto the opposite segment.
//
// Parameters:
//   - CD: LineSegment[T] - The line segment to compare with `AB`.
//   - opts: A variadic slice of Option functions to customize the behavior of the calculation.
//     WithEpsilon(epsilon float64): Specifies a tolerance for snapping near-zero results to zero or
//     for handling small floating-point deviations in distance calculations.
//
// Behavior:
//   - The function checks whether the segments intersect or touch. If so, the distance is immediately returned as 0.
//   - For non-intersecting segments, the function calculates the shortest distance using the following steps:
//     1. Compute direct distances between the endpoints of `AB` and `CD`.
//     2. Compute distances to the perpendicular projections of each endpoint onto the opposite segment.
//     3. Track the minimum distance among all calculations and return this value.
//   - If `WithEpsilon` is provided, epsilon adjustments are applied to the calculated distances and projections
//     to ensure robustness against floating-point precision errors.
//
// Returns:
//   - float64: The minimum distance between the two segments. If the segments intersect or touch, this value is 0.
//
// Example Usage:
//
//	segmentAB := NewLineSegment(NewPoint(0, 0), NewPoint(2, 2))
//	segmentCD := NewLineSegment(NewPoint(3, 3), NewPoint(5, 5))
//
//	// Default behavior (no epsilon adjustment)
//	distance := segmentAB.DistanceToLineSegment(segmentCD)
//
//	// With epsilon adjustment
//	distanceWithEpsilon := segmentAB.DistanceToLineSegment(segmentCD, WithEpsilon(1e-4))
//
// Notes:
//   - Epsilon adjustment is particularly useful when working with floating-point coordinates, where small
//     imprecisions might cause distances to deviate slightly from expected values.
//   - This function is backward-compatible, maintaining exact calculations if no options are provided.
func (AB LineSegment[T]) DistanceToLineSegment(CD LineSegment[T], opts ...Option) float64 {
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
	updateMinDist(ABf.start.DistanceToPoint(CDf.start, opts...))
	updateMinDist(ABf.start.DistanceToPoint(CDf.end, opts...))
	updateMinDist(ABf.end.DistanceToPoint(CDf.start, opts...))
	updateMinDist(ABf.end.DistanceToPoint(CDf.end, opts...))

	// Calculate distances to projections on the opposite segment.
	updateMinDist(ABf.start.DistanceToPoint(ABf.start.ProjectOntoLineSegment(CDf), opts...))
	updateMinDist(ABf.end.DistanceToPoint(ABf.end.ProjectOntoLineSegment(CDf), opts...))
	updateMinDist(CDf.start.DistanceToPoint(CDf.start.ProjectOntoLineSegment(ABf), opts...))
	updateMinDist(CDf.end.DistanceToPoint(CDf.end.ProjectOntoLineSegment(ABf), opts...))

	return minDist
}

// DistanceToPoint calculates the orthogonal (shortest) distance from a specified Point p
// to the LineSegment AB. This distance is determined by projecting the Point p onto the
// LineSegment and measuring the distance from p to the projected Point.
//
// Parameters:
//   - p: The Point for which the distance to the LineSegment AB is calculated.
//   - opts: A variadic slice of Option functions to customize the behavior of the calculation.
//     WithEpsilon(epsilon float64): Specifies a tolerance for snapping near-zero results to zero
//     or handling small floating-point deviations in distance calculations.
//
// Behavior:
//   - The function computes the projection of p onto the given LineSegment AB. This is the closest
//     point on AB to p, whether it lies within the segment bounds or at an endpoint.
//   - The orthogonal distance is then calculated as the Euclidean distance between p and the
//     projected point.
//   - If `WithEpsilon` is provided, epsilon adjustments are applied to the calculated distance
//     to ensure robustness against floating-point precision errors.
//
// Returns:
//   - float64: The shortest distance between the point p and the LineSegment AB, optionally
//     adjusted based on epsilon.
//
// Example Usage:
//
//	segmentAB := NewLineSegment(NewPoint(0, 0), NewPoint(6, 8))
//	p := NewPoint(3, 4)
//
//	// Default behavior (no epsilon adjustment)
//	distance := segmentAB.DistanceToPoint(p)
//
//	// With epsilon adjustment
//	distanceWithEpsilon := segmentAB.DistanceToPoint(p, WithEpsilon(1e-4))
//
// Notes:
//   - This function leverages the Point.DistanceToLineSegment method to perform the calculation,
//     ensuring precision and consistency across related operations.
//   - Epsilon adjustment is particularly useful for applications involving floating-point data,
//     where small deviations can affect the accuracy of results.
func (AB LineSegment[T]) DistanceToPoint(p Point[T], opts ...Option) float64 {
	return p.DistanceToLineSegment(AB, opts...)
}

// End returns the ending point of the line segment.
//
// This function provides access to the ending point of the line segment `AB`, typically representing
// the endpoint of the segment.
func (AB LineSegment[T]) End() Point[T] {
	return AB.end
}

// Eq checks if two line segments are equal by comparing their start and end points.
// Equality can be evaluated either exactly (default) or approximately using an epsilon threshold.
//
// Parameters:
//   - CD: LineSegment[T] - The line segment to compare with the current line segment.
//   - opts: A variadic slice of Option functions to customize the equality check.
//     WithEpsilon(epsilon float64): Specifies a tolerance for comparing the start and end
//     points of the line segments. If the absolute difference between the coordinates of
//     the points is less than epsilon, they are considered equal.
//
// Behavior:
//   - By default, the function performs an exact equality check, returning true only if
//     both the `start` and `end` points of `AB` and `CD` are identical.
//   - If `WithEpsilon` is provided, the function performs an approximate equality check,
//     considering the points equal if their coordinate differences are within the specified
//     epsilon threshold.
//
// Returns:
//   - bool - Returns `true` if both line segments have identical (or approximately equal) start
//     and end points; otherwise, `false`.
//
// Example Usage:
//
//	segment1 := NewLineSegment(NewPoint(1.0, 1.0), NewPoint(4.0, 5.0))
//	segment2 := NewLineSegment(NewPoint(1.0, 1.0), NewPoint(4.0, 5.0))
//
//	// Default behavior (exact equality)
//	equal := segment1.Eq(segment2) // equal will be true
//
//	// Approximate equality with epsilon
//	segment3 := NewLineSegment(NewPoint(1.00001, 1.00001), NewPoint(4.00001, 5.00001))
//	approximatelyEqual := segment1.Eq(segment3, WithEpsilon(1e-4)) // approximatelyEqual will be true
//
// Notes:
//   - Approximate equality is useful when comparing line segments with floating-point coordinates,
//     where small precision errors might otherwise cause inequality.
//   - This function relies on the `Eq` method of the `Point` type, which supports epsilon adjustments.
func (AB LineSegment[T]) Eq(CD LineSegment[T], opts ...Option) bool {
	return AB.start.Eq(CD.start, opts...) && AB.end.Eq(CD.end, opts...)
}

// IntersectsLineSegment checks whether there is any intersection or overlap between LineSegment AB and LineSegment CD.
//
// This function returns true if segments `AB` and `CD` have an intersecting spatial relationship, such as intersection,
// overlap, containment, or endpoint coincidence. It leverages the `RelationshipToLineSegment` function to
// determine if the relationship value is greater than `LLRMiss`, indicating that the segments are not fully
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
	if AB.RelationshipToLineSegment(CD) > LLRMiss {
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
	intersection := A.Translate(dir1.Scale(NewPoint[float64](0, 0), t))
	return intersection, true
}

// Length calculates the length of the line segment, optionally using an epsilon threshold
// to adjust the precision of the calculation.
//
// Parameters:
//   - opts: A variadic slice of Option functions to customize the behavior of the calculation.
//   - WithEpsilon(epsilon float64): Specifies a tolerance for snapping small floating-point
//     deviations in the calculated length to cleaner values, improving robustness.
//
// Behavior:
//   - The function computes the Euclidean distance between the `start` and `end` points of the
//     line segment using the `DistanceToPoint` method.
//   - If `WithEpsilon` is provided, the resulting length is adjusted such that small deviations
//     due to floating-point precision errors are corrected.
//
// Returns:
//   - float64: The length of the line segment, optionally adjusted based on epsilon.
//
// Example Usage:
//
//	segment := NewLineSegment(NewPoint(0, 0), NewPoint(3, 4))
//
//	// Default behavior (no epsilon adjustment)
//	length := segment.Length() // length will be 5.0
//
// Notes:
//   - This function relies on `DistanceToPoint`, which supports epsilon adjustments for distance
//     calculations. Epsilon is particularly useful for floating-point coordinates where minor
//     imprecisions might affect the result.
func (AB LineSegment[T]) Length(opts ...Option) float64 {
	return AB.start.DistanceToPoint(AB.end, opts...)
}

func (AB LineSegment[T]) Perimeter(opts ...Option) float64 {
	return AB.Length(opts...)
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

// RelationshipToCircle determines the spatial relationship of the line segment
// to a circle. It returns one of several possible relationships, such as whether
// the segment is inside, outside, tangent to, or intersects the circle.
//
// Parameters:
//   - c: The circle to analyze.
//   - opts: A variadic slice of Option functions to customize the behavior of the relationship check.
//     WithEpsilon(epsilon float64): Specifies a tolerance for comparing distances to the circle's radius,
//     improving robustness against floating-point precision errors.
//
// Returns:
//   - CircleLineSegmentRelationship: An enum value indicating the relationship.
//
// Possible Relationships:
//   - CLROutside: The segment lies entirely outside the circle.
//   - CLRInside: The segment lies entirely within the circle.
//   - CLRIntersecting: The segment intersects the circle at two points.
//   - CLRTangent: The segment is tangent to the circle, touching it at exactly one point.
//   - CLROneEndOnCircumferenceOutside: One endpoint is on the circle's boundary, and the other is outside.
//   - CLROneEndOnCircumferenceInside: One endpoint is on the circle's boundary, and the other is inside.
//   - CLRBothEndsOnCircumference: Both endpoints lie on the circle's boundary.
//
// Example Usage:
//
//	c := NewCircle(NewPoint(0, 0), 5)
//	segment := NewLineSegment(NewPoint(0, -6), NewPoint(0, 6))
//
//	// Default behavior (no epsilon adjustment)
//	relationship := segment.RelationshipToCircle(c)
//
//	// With epsilon adjustment
//	relationshipWithEpsilon := segment.RelationshipToCircle(c, WithEpsilon(1e-4))
//
// Notes:
//   - Epsilon adjustment is particularly useful for floating-point coordinates, where small precision
//     errors might otherwise cause incorrect classifications.
func (AB LineSegment[T]) RelationshipToCircle(c Circle[T], opts ...Option) CircleLineSegmentRelationship {
	return c.RelationshipToLineSegment(AB, opts...)
}

// RelationshipToLineSegment determines the spatial relationship between two line segments, AB and CD.
//
// This function evaluates the relationship between two line segments, `AB` and `CD`, by checking for
// endpoint coincidences, intersections, collinear relationships, and containment. It returns a
// `LineSegmentLineSegmentRelationship` constant that describes the exact relationship between the segments, such
// as intersection, partial overlap, or full containment.
//
// Parameters:
//   - CD: LineSegment[T] - The line segment to compare with `AB`.
//   - opts: A variadic slice of Option functions to customize the behavior of the relationship check.
//     WithEpsilon(epsilon float64): Specifies a tolerance for comparing points and collinearity calculations,
//     allowing for robust handling of floating-point precision errors.
//
// Behavior:
//   - The function first checks if the two line segments are exactly equal (or approximately equal if epsilon is provided).
//   - It then evaluates endpoint coincidences, collinearity, intersection, and containment using orientation tests,
//     point-on-segment checks, and direct comparisons.
//   - If WithEpsilon is provided, epsilon adjustments are applied to point comparisons, collinearity checks, and
//     point-on-segment tests to ensure robustness against floating-point imprecision.
//
// Returns:
//   - LineSegmentLineSegmentRelationship: A constant that describes the relationship between segments AB and CD.
//
// Possible return values:
//   - LLRCollinearDisjoint: The segments are collinear but do not overlap or touch at any point.
//   - LLRMiss: The segments are not collinear, disjoint, and do not intersect, overlap, or touch at any point.
//   - LLRIntersects: The segments intersect at a unique point that is not an endpoint.
//   - LLRAeqC, LLRAeqD, LLRBeqC, LLRBeqD: An endpoint of AB coincides with an endpoint of CD. For example,
//     LLRAeqC indicates that point A of AB coincides with point C of CD.
//   - LLRAonCD, LLRBonCD, LLRConAB, LLRDonAB: One endpoint of one segment lies on the other segment without
//     the segments being collinear. For example, LLRAonCD indicates that point A of AB lies on segment CD.
//   - LLRCollinearAonCD, LLRCollinearBonCD: The segments are collinear with partial overlap where one endpoint of one
//     segment lies on the other. For example, LLRCollinearAonCD means point A of AB lies on CD with collinearity.
//   - LLRCollinearABinCD: The entire segment AB is contained within segment CD.
//   - LLRCollinearCDinAB: The entire segment CD is contained within segment AB.
//   - LLRCollinearEqual: The segments AB and CD are exactly equal, sharing both endpoints in the same locations.
//
// Example Usage:
//
//	segmentAB := NewLineSegment(NewPoint(0, 0), NewPoint(2, 2))
//	segmentCD := NewLineSegment(NewPoint(1, 1), NewPoint(3, 3))
//
//	// Default behavior (no epsilon adjustment)
//	relationship := segmentAB.RelationshipToLineSegment(segmentCD)
//
//	// With epsilon adjustment
//	relationshipWithEpsilon := segmentAB.RelationshipToLineSegment(segmentCD, WithEpsilon(1e-4))
//
//	// The variable `relationship` should equal `LLRCollinearAonCD` as:
//	//   - `AB` and `CD` are collinear
//	//   - `A` lies on `CD`
//	//   - `AB` and `CD` don't fully overlap.
//
// Notes:
//   - Epsilon adjustment is particularly useful when working with floating-point coordinates, where small
//     precision errors might otherwise cause incorrect results.
//   - This function relies on the `Eq`, `Orientation`, and `IsOnLineSegment` methods, all of which support
//     epsilon adjustments.
func (AB LineSegment[T]) RelationshipToLineSegment(CD LineSegment[T], opts ...Option) LineSegmentLineSegmentRelationship {

	// Check if segments are exactly equal
	if (AB.start.Eq(CD.start, opts...) && AB.end.Eq(CD.end, opts...)) || (AB.start.Eq(CD.end, opts...) && AB.end.Eq(CD.start, opts...)) {
		return LLRCollinearEqual
	}

	switch {

	// Check if A and C coincide
	case AB.start.Eq(CD.start, opts...):
		return LLRAeqC

	// Check if A and D coincide
	case AB.start.Eq(CD.end, opts...):
		return LLRAeqD

	// Check if End and C coincide
	case AB.end.Eq(CD.start, opts...):
		return LLRBeqC

	// Check if End and D coincide
	case AB.end.Eq(CD.end, opts...):
		return LLRBeqD

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
			return LLRAonCD

		// Check if End lies on CD
		case !AB.start.IsOnLineSegment(CD) && AB.end.IsOnLineSegment(CD):
			return LLRBonCD

		// Check if C lies on AB
		case CD.start.IsOnLineSegment(AB) && !CD.end.IsOnLineSegment(AB):
			return LLRConAB

		// Check if D lies on AB
		case !CD.start.IsOnLineSegment(AB) && CD.end.IsOnLineSegment(AB):
			return LLRDonAB

		// Default case that lines intersect without any "edge cases"
		default:
			return LLRIntersects
		}
	}

	// PointsCollinear cases: All orientations are zero
	if o1 == 0 && o2 == 0 && o3 == 0 && o4 == 0 {
		// Check if segments are collinear and disjoint
		if !AB.start.IsOnLineSegment(CD) && !AB.end.IsOnLineSegment(CD) &&
			!CD.start.IsOnLineSegment(AB) && !CD.end.IsOnLineSegment(AB) {
			return LLRCollinearDisjoint
		}
		// Check if AB is fully contained within CD
		if AB.start.IsOnLineSegment(CD) && AB.end.IsOnLineSegment(CD) {
			return LLRCollinearABinCD
		}
		// Check if CD is fully contained within AB
		if CD.start.IsOnLineSegment(AB) && CD.end.IsOnLineSegment(AB) {
			return LLRCollinearCDinAB
		}
		// Check specific collinear partial overlaps
		if AB.start.IsOnLineSegment(CD) {
			return LLRCollinearAonCD
		}
		if AB.end.IsOnLineSegment(CD) {
			return LLRCollinearBonCD
		}
	}

	// If none of the conditions matched, the segments are disjoint
	return LLRMiss
}

func (AB LineSegment[T]) RelationshipToPoint(p Point[T], opts ...Option) PointLineSegmentRelationship {
	return p.RelationshipToLineSegment(AB)
}

// RelationshipToPolyTree determines the spatial relationship of a line segment to a PolyTree.
//
// The function evaluates whether the line segment:
//   - Intersects any boundary within the PolyTree.
//   - Is entirely within a solid or hole polygon in the PolyTree.
//   - Lies entirely outside the PolyTree.
//
// Parameters:
//   - tree: A pointer to the PolyTree to evaluate.
//   - opts: Optional configurations for geometric calculations, such as epsilon tolerance for floating-point comparisons.
//
// Returns:
//   - PolyTreeLineSegmentRelationship: The relationship between the line segment and the PolyTree.
//
// Behavior:
//   - If the line segment intersects any boundary (considering epsilon), the function immediately returns PTLRIntersectsBoundary.
//   - If the segment's endpoints are entirely contained within the same polygon, the function returns PTLRInsideSolid or PTLRInsideHole, depending on the polygon type.
//   - The function uses an epsilon tolerance when checking relationships between line segments and polygon edges.
//   - If no stronger relationship is found, the function returns PTLRMiss, indicating the segment is entirely outside the PolyTree.
func (AB LineSegment[T]) RelationshipToPolyTree(tree *PolyTree[T], opts ...Option) PolyTreeLineSegmentRelationship {

	// as the points in a polytree contour are doubled, we need to also double the input line segment
	lineSegmentDoubled := AB.Scale(NewPoint[T](0, 0), 2)

	highestRel := PTLRMiss // Default to outside

	// Iterate through each polygon in the tree
	for poly := range tree.iterPolys {
		// Check each edge of the polygon's contour
		for edge := range poly.contour.iterEdges {
			// Determine relationship between poly contour & line segment
			rel := edge.RelationshipToLineSegment(lineSegmentDoubled, opts...)
			// any intersection
			if rel > LLRMiss {
				return PTLRIntersectsBoundary
			}
		}

		// check for containment
		if poly.contour.isPointInside(lineSegmentDoubled.start) && poly.contour.isPointInside(lineSegmentDoubled.end) {
			switch poly.polygonType {
			case PTSolid:
				highestRel = PTLRInsideSolid
			case PTHole:
				highestRel = PTLRInsideHole
			}
		}
	}

	return highestRel
}

func (AB LineSegment[T]) RelationshipToRectangle(r Rectangle[T]) RectangleLineSegmentRelationship {
	return r.RelationshipToLineSegment(AB)
}

// Rotate rotates the LineSegment around a given pivot point by a specified angle in radians.
// Optionally, an epsilon threshold can be applied to adjust the precision of the resulting coordinates.
//
// Parameters:
//   - pivot: Point[T] - The point around which to rotate the line segment.
//   - radians: float64 - The rotation angle in radians.
//   - opts: A variadic slice of Option functions to customize the behavior of the rotation.
//     WithEpsilon(epsilon float64): Specifies a tolerance for snapping near-zero or near-integer
//     values in the resulting coordinates to cleaner values, improving robustness.
//
// Behavior:
//   - The function rotates the start and end points of the line segment around the given pivot
//     point by the specified angle using the Point.Rotate method.
//   - If WithEpsilon is provided, epsilon adjustments are applied to the rotated coordinates to
//     handle floating-point precision errors.
//
// Returns:
//   - LineSegment[float64]: A new line segment representing the rotated position, with floating-point coordinates.
//
// Example Usage:
//
//	segment := geom2d.NewLineSegment(geom2d.NewPoint(0, 0), geom2d.NewPoint(10, 0))
//	pivot := geom2d.NewPoint(0, 0)
//	angle := math.Pi / 2 // Rotate 90 degrees
//
//	// Default behavior (no epsilon adjustment)
//	rotatedSegment := segment.Rotate(pivot, angle)
//	// rotatedSegment = LineSegment[(0, 0) -> (6.123233995736757e-16, 10)]
//
//	// With epsilon adjustment
//	rotatedSegmentWithEpsilon := segment.Rotate(pivot, angle, geom2d.WithEpsilon(1e-10))
//	// rotatedSegmentWithEpsilon =LineSegment[(0, 0) -> (0, 10)]
//
// Notes:
//   - Epsilon adjustment is particularly useful when the rotation involves floating-point
//     calculations that could result in minor inaccuracies.
//   - The returned line segment always has `float64` coordinates, ensuring precision regardless
//     of the coordinate type of the original line segment.
func (AB LineSegment[T]) Rotate(pivot Point[T], radians float64, opts ...Option) LineSegment[float64] {
	newStart := AB.start.Rotate(pivot, radians, opts...)
	newEnd := AB.end.Rotate(pivot, radians, opts...)
	return NewLineSegment(newStart, newEnd)
}

// Scale scales the line segment by a given factor from a specified origin point.
//
// Parameters:
//   - origin: ScaleOrigin - The point from which to scale (ScaleFromStart, ScaleFromEnd, or ScaleFromMidpoint).
//   - factor: float64 - The scaling factor.
//
// Returns:
//   - LineSegment[float64] - A new line segment scaled relative to the specified origin.
//
// todo: Implement ScaleFrom as per Point type
func (AB LineSegment[T]) Scale(ref Point[T], factor T) LineSegment[T] {
	return NewLineSegment(
		AB.start.Scale(ref, factor),
		AB.end.Scale(ref, factor),
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

// Translate moves the line segment by a specified vector.
//
// This method shifts the line segment's position in the 2D plane by translating
// both its start and end points by the given vector `delta`. The relative
// orientation and length of the line segment remain unchanged.
//
// Parameters:
//   - delta: Point[T] - The vector by which to translate the line segment.
//
// Returns:
//   - LineSegment[T]: A new LineSegment translated by the specified vector.
//
// Example Usage:
//
//	lineSegment := NewLineSegment(NewPoint(1, 1), NewPoint(4, 4))
//	translationVector := NewPoint(2, 3)
//	translatedLineSegment := lineSegment.Translate(translationVector)
//	// translatedLineSegment has its start point at (3, 4)
//	// and end point at (6, 7), preserving its length and orientation.
func (AB LineSegment[T]) Translate(delta Point[T]) LineSegment[T] {
	return NewLineSegment(
		AB.start.Translate(delta),
		AB.end.Translate(delta),
	)
}
