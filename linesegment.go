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
//   - other (LineSegment[T]): The line segment to add to the current one.
//
// Returns:
//   - LineSegment[T] - A new line segment where each endpoint is the sum of the corresponding endpoints.
func (l LineSegment[T]) AddLineSegment(other LineSegment[T]) LineSegment[T] {
	return NewLineSegment(
		l.start.Translate(other.start),
		l.end.Translate(other.end),
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
func (l LineSegment[T]) Area() float64 {
	return 0
}

// AsFloat converts the line segment to a LineSegment[float64] type.
//
// This function converts both endpoints of the LineSegment l to [Point][float64]
// values, creating a new line segment with floating-point coordinates.
// It is useful for precise calculations where floating-point accuracy is needed.
//
// Returns:
//   - LineSegment[float64] - The line segment with both endpoints converted to float64.
func (l LineSegment[T]) AsFloat() LineSegment[float64] {
	return NewLineSegment(l.start.AsFloat(), l.end.AsFloat())
}

// AsInt converts the line segment to a LineSegment[int] type.
//
// This function converts both endpoints of the line segment l to [Point][int]
// by truncating any decimal places. It is useful for converting a floating-point
// line segment to integer coordinates without rounding.
//
// Returns:
//   - LineSegment[int] - The line segment with both endpoints converted to integer coordinates by truncation.
func (l LineSegment[T]) AsInt() LineSegment[int] {
	return NewLineSegment(l.start.AsInt(), l.end.AsInt())
}

// AsIntRounded converts the line segment to a LineSegment[int] type with rounded coordinates.
//
// This function converts both endpoints of the line segment l to [Point][int]
// by rounding each coordinate to the nearest integer. It is useful when you need to
// approximate the segment’s position with integer coordinates while minimizing the
// rounding error.
//
// Returns:
//   - LineSegment[int] - The line segment with both endpoints converted to integer coordinates by rounding.
func (l LineSegment[T]) AsIntRounded() LineSegment[int] {
	return NewLineSegment(l.start.AsIntRounded(), l.end.AsIntRounded())
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
func (l LineSegment[T]) BoundingBox() Rectangle[T] {
	points := []Point[T]{
		l.start,
		NewPoint(l.start.x, l.end.y),
		NewPoint(l.end.x, l.start.y),
		l.end,
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
func (l LineSegment[T]) Center(opts ...Option) Point[float64] {
	// Apply geomOptions with defaults
	options := applyOptions(geomOptions{epsilon: 0}, opts...)

	start := l.start.AsFloat()
	end := l.end.AsFloat()

	midX := (start.x + end.x) / 2
	midY := (start.y + end.y) / 2

	// Apply epsilon if specified
	if options.epsilon > 0 {
		midX = applyEpsilon(midX, options.epsilon)
		midY = applyEpsilon(midY, options.epsilon)
	}

	return NewPoint[float64](midX, midY)
}

// ContainsPoint determines whether the given [Point] lies on the LineSegment.
//
// This method first checks if the [Point] is collinear with the endpoints of the
// [LineSegment] using an [Orientation]. If the [Point] is not collinear, it
// cannot be on the segment. If the [Point] is collinear, the function then verifies
// if the [Point] lies within the bounding box defined by the segment's endpoints.
//
// Parameters:
//   - p ([Point][T]): The [Point] to test against the LineSegment
//
// Returns:
//   - bool: true if the [Point] lies on the LineSegment, false otherwise.
func (l LineSegment[T]) ContainsPoint(p Point[T]) bool {
	// Check collinearity first; if not collinear, point is not on the line segment
	if Orientation(p, l.start, l.end) != PointsCollinear {
		return false
	}

	// Check if the point lies within the bounding box defined by A and End
	return p.x >= min(l.start.x, l.end.x) && p.x <= max(l.start.x, l.end.x) &&
		p.y >= min(l.start.y, l.end.y) && p.y <= max(l.start.y, l.end.y)
}

// DistanceToLineSegment calculates the minimum distance between two line segments, l and other.
//
// If the segments intersect or touch at any point, the function returns 0, as the distance is effectively zero.
// Otherwise, it calculates the shortest distance by considering distances between segment endpoints and the
// perpendicular projections of each endpoint onto the opposite segment.
//
// Parameters:
//   - other (LineSegment[T]): The line segment to compare with l.
//   - opts: A variadic slice of [Option] functions to customize the behavior of the calculation.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for snapping near-zero results to zero or
//     for handling small floating-point deviations in distance calculations.
//
// Behavior:
//   - The function checks whether the segments intersect or touch. If so, the distance is immediately returned as 0.
//   - For non-intersecting segments, the function calculates the shortest distance using the following steps:
//     1. Compute direct distances between the endpoints of l and other.
//     2. Compute distances to the perpendicular projections of each endpoint onto the opposite segment.
//     3. Track the minimum distance among all calculations and return this value.
//   - If [WithEpsilon] is provided, epsilon adjustments are applied to the calculated distances and projections
//     to ensure robustness against floating-point precision errors.
//
// Returns:
//   - float64: The minimum distance between the two segments. If the segments intersect or touch, this value is 0.
func (l LineSegment[T]) DistanceToLineSegment(other LineSegment[T], opts ...Option) float64 {
	// If line segments intersect, the distance is zero.
	if l.IntersectsLineSegment(other) {
		return 0
	}

	// Convert segments to float for precise calculations.
	ABf, CDf := l.AsFloat(), other.AsFloat()

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

// DistanceToPoint calculates the orthogonal (shortest) distance from a specified [Point] p
// to the LineSegment l. This distance is determined by projecting the [Point] p onto the
// LineSegment and measuring the distance from p to the projected [Point].
//
// Parameters:
//   - p ([Point][T]): The [Point] for which the distance to the LineSegment l is calculated.
//   - opts: A variadic slice of [Option] functions to customize the behavior of the calculation.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for snapping near-zero results to zero
//     or handling small floating-point deviations in distance calculations.
//
// Behavior:
//   - The function computes the projection of p onto the given LineSegment l. This is the closest
//     point on l to p, whether it lies within the segment bounds or at an endpoint.
//   - The orthogonal distance is then calculated as the Euclidean distance between p and the
//     projected point.
//   - If [WithEpsilon] is provided, epsilon adjustments are applied to the calculated distance
//     to ensure robustness against floating-point precision errors.
//
// Returns:
//   - float64: The shortest distance between the [Point] p and the LineSegment l, optionally
//     adjusted based on epsilon.
//
// Notes:
//   - This function leverages the [Point.DistanceToLineSegment] method to perform the calculation,
//     ensuring precision and consistency across related operations.
//   - Epsilon adjustment is particularly useful for applications involving floating-point data,
//     where small deviations can affect the accuracy of results.
func (l LineSegment[T]) DistanceToPoint(p Point[T], opts ...Option) float64 {
	return p.DistanceToLineSegment(l, opts...)
}

// End returns the ending [Point] of the line segment.
//
// This function provides access to the ending [Point] of the line segment l, typically representing
// the endpoint of the segment.
func (l LineSegment[T]) End() Point[T] {
	return l.end
}

// Eq checks if two line segments are equal by comparing their start and end points.
// Equality can be evaluated either exactly (default) or approximately using an epsilon threshold.
//
// Parameters:
//   - other (LineSegment[T]): - The line segment to compare with the current line segment.
//   - opts: A variadic slice of [Option] functions to customize the equality check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing the start and end
//     points of the line segments. If the absolute difference between the coordinates of
//     the points is less than epsilon, they are considered equal.
//
// Behavior:
//   - By default, the function performs an exact equality check, returning true only if
//     both the start and end points of l and other are identical.
//   - If [WithEpsilon] is provided, the function performs an approximate equality check,
//     considering the points equal if their coordinate differences are within the specified
//     epsilon threshold.
//
// Returns:
//   - bool: Returns true if both line segments have identical (or approximately equal with epsilon) start
//     and end points; otherwise, false.
//
// Notes:
//   - Approximate equality is useful when comparing line segments with floating-point coordinates,
//     where small precision errors might otherwise cause inequality.
//   - This function relies on the [Point.Eq] method, which supports epsilon adjustments.
func (l LineSegment[T]) Eq(other LineSegment[T], opts ...Option) bool {
	return l.start.Eq(other.start, opts...) && l.end.Eq(other.end, opts...)
}

// IntersectionPoint calculates the intersection [Point] between two line segments, if one exists.
//
// This method checks if the LineSegment l and LineSegment other intersect within their boundaries
// and, if so, calculates and returns the intersection point. If the segments do not intersect
// or are parallel, it returns a zero-value [Point] and false.
//
// It uses the parametric form of the line segments to solve for intersection parameters t and u.
// If t and u are both in the range [0, 1], the intersection point lies within the bounds of
// both segments.
//
// Parameters:
//   - other (LineSegment[T]): The second line segment to check for an intersection.
//
// Returns:
//   - [Point][T]: The intersection point.
//   - bool: If true, the first element is the intersection point. If false, there is
//     no intersection within the segments’ bounds or the segments are parallel.
func (l LineSegment[T]) IntersectionPoint(other LineSegment[T]) (Point[float64], bool) {
	// Define segment endpoints for AB (l) and CD (other)
	A, B := l.start.AsFloat(), l.end.AsFloat()
	C, D := other.start.AsFloat(), other.end.AsFloat()

	// Calculate the direction vectors

	dir1 := B.Translate(A.Negate())
	dir2 := D.Translate(C.Negate())

	// Calculate the determinants
	denominator := dir1.CrossProduct(dir2)

	// Check if the lines are parallel (no intersection)
	if denominator == 0 {
		return Point[float64]{}, false // No intersection
	}

	// Calculate parameters t and u
	AC := C.Translate(A.Negate())
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

// IntersectsLineSegment checks whether there is any intersection or overlap between LineSegment l and LineSegment other.
//
// This function returns true if segments l and other have an intersecting spatial relationship, such as intersection,
// overlap, containment, or endpoint coincidence. It leverages the [LineSegment.RelationshipToLineSegment] function to
// determine if the relationship value is greater than [RelationshipLineSegmentLineSegmentMiss], indicating that the segments are not fully
// disjoint.
//
// Parameters:
//   - other (LineSegment[T]): The line segment to compare with AB.
//
// Returns:
//   - bool: Returns true if l and other intersect, overlap, or share any form of intersecting relationship, and false if they are completely disjoint.
//
// Example usage:
//
//	segmentAB := NewLineSegment(NewPoint(0, 0), NewPoint(2, 2))
//	segmentCD := NewLineSegment(NewPoint(1, 1), NewPoint(3, 3))
//	intersects := segmentAB.IntersectsLineSegment(segmentCD)
//
// `intersects` will be `true` as there is an intersecting relationship between `AB` and `CD`.
// todo: may not be needed: RelationshipTo* methods can make IntersectsLineSegment and ContainsPoint redundant
func (l LineSegment[T]) IntersectsLineSegment(other LineSegment[T]) bool {
	if l.RelationshipToLineSegment(other) > RelationshipLineSegmentLineSegmentMiss {
		return true
	}
	return false
}

// Perimeter calculates the length of the line segment, optionally using an epsilon threshold
// to adjust the precision of the calculation.
//
// Parameters:
//   - opts: A variadic slice of [Option] functions to customize the behavior of the calculation.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for snapping small floating-point
//     deviations in the calculated length to cleaner values, improving robustness.
//
// Behavior:
//   - The function computes the Euclidean distance between the start and end points of the
//     line segment using the [LineSegment.DistanceToPoint] method.
//   - If [WithEpsilon] is provided, the resulting length is adjusted such that small deviations
//     due to floating-point precision errors are corrected.
//
// Returns:
//   - float64: The length of the line segment, optionally adjusted based on epsilon.
//
// Notes:
//   - This function relies on [LineSegment.DistanceToPoint], which supports epsilon adjustments for distance
//     calculations. Epsilon is particularly useful for floating-point coordinates where minor
//     imprecision might affect the result.
//   - Yes, we know "Perimeter" is an odd choice of name for a line segment (ideally it would be "Length").
//     Rest assured, this was done to meet the [Measurable] interface requirements and to standardize functions
//     across all geometric types.
func (l LineSegment[T]) Perimeter(opts ...Option) float64 {
	return l.start.DistanceToPoint(l.end, opts...)
}

// Points returns the two endpoints of the line segment as a slice of Points.
// The order of the points is [start, end].
//
// Returns:
//   - [][Point][T]: A slice containing the start and end points of the line segment.
func (l LineSegment[T]) Points() []Point[T] {
	return []Point[T]{l.start, l.end}
}

// Reflect reflects the line segment across the specified axis or custom line.
//
// Parameters:
//   - axis ([ReflectionAxis]): The axis or line to reflect across ([ReflectAcrossXAxis], [ReflectAcrossYAxis], or [ReflectAcrossCustomLine]).
//   - line (LineSegment[float64]): Optional. The line segment for [ReflectAcrossCustomLine] reflection.
//
// Returns:
//   - LineSegment[float64] - A new line segment where both endpoints are reflected accordingly.
func (l LineSegment[float64]) Reflect(axis ReflectionAxis, line ...LineSegment[float64]) LineSegment[float64] {
	var startReflected, endReflected Point[float64]
	switch axis {
	case ReflectAcrossXAxis:
		// Reflect across the x-axis
		startReflected = l.start.Reflect(ReflectAcrossXAxis)
		endReflected = l.end.Reflect(ReflectAcrossXAxis)
	case ReflectAcrossYAxis:
		// Reflect across the y-axis
		startReflected = l.start.Reflect(ReflectAcrossYAxis)
		endReflected = l.end.Reflect(ReflectAcrossYAxis)
	case ReflectAcrossCustomLine:
		// Reflect across a custom line if provided
		if len(line) > 0 {
			startReflected = l.start.Reflect(ReflectAcrossCustomLine, line[0])
			endReflected = l.end.Reflect(ReflectAcrossCustomLine, line[0])
		} else {
			// No custom line provided; return the original line segment unchanged
			return l
		}
	default:
		// Invalid axis, return the line segment unchanged
		return l
	}

	// Return a new line segment with reflected points
	return NewLineSegment(startReflected, endReflected)
}

// RelationshipToCircle determines the spatial relationship of the line segment
// to a [Circle]. It returns one of several possible relationships, such as whether
// the segment is inside, outside, tangent to, or intersects the [Circle].
//
// Parameters:
//   - c ([Circle][T]): The [Circle] to analyze.
//   - opts: A variadic slice of [Option] functions to customize the behavior of the relationship check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing distances to the circle's radius,
//     improving robustness against floating-point precision errors.
//
// Returns:
//   - [RelationshipLineSegmentCircle]: An enum value indicating the relationship.
//
// Notes:
//   - Epsilon adjustment is particularly useful for floating-point coordinates, where small precision
//     errors might otherwise cause incorrect classifications.
func (l LineSegment[T]) RelationshipToCircle(c Circle[T], opts ...Option) RelationshipLineSegmentCircle {
	return c.RelationshipToLineSegment(l, opts...)
}

// RelationshipToLineSegment determines the spatial relationship between two line segments, l and other.
//
//   - Let A = l.Start()
//   - Let B = l.End()
//   - Let C = other.Start()
//   - Let D = other.End()
//
// This function evaluates the relationship between two line segments, AB and CD, by checking for
// endpoint coincidences, intersections, collinear relationships, and containment. It returns a
// [RelationshipLineSegmentLineSegment] constant that describes the exact relationship between the segments, such
// as intersection, partial overlap, full containment, etc.
//
// Output constants may include references to A, B, C or D (for brevity).
//
// Parameters:
//   - other (LineSegment[T]): The line segment to compare with l.
//   - opts: A variadic slice of [Option] functions to customize the behavior of the relationship check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing points and collinearity calculations,
//     allowing for robust handling of floating-point precision errors.
//
// Behavior:
//   - The function first checks if the two line segments are exactly equal (or approximately equal if epsilon is provided).
//   - It then evaluates endpoint coincidences, collinearity, intersection, and containment using orientation tests,
//     point-on-segment checks, and direct comparisons.
//   - If [WithEpsilon] is provided, epsilon adjustments are applied to point comparisons, collinearity checks, and
//     point-on-segment tests to ensure robustness against floating-point imprecision.
//
// Returns:
//   - [RelationshipLineSegmentLineSegment]: A constant that describes the relationship between segments AB and CD.
//
// Notes:
//   - Epsilon adjustment is particularly useful when working with floating-point coordinates, where small
//     precision errors might otherwise cause incorrect results.
//
// todo: check wording with astra
func (l LineSegment[T]) RelationshipToLineSegment(other LineSegment[T], opts ...Option) RelationshipLineSegmentLineSegment {

	// Check if segments are exactly equal
	if (l.start.Eq(other.start, opts...) && l.end.Eq(other.end, opts...)) || (l.start.Eq(other.end, opts...) && l.end.Eq(other.start, opts...)) {
		return RelationshipLineSegmentLineSegmentCollinearEqual
	}

	switch {

	// Check if A and C coincide
	case l.start.Eq(other.start, opts...):
		return RelationshipLineSegmentLineSegmentAeqC

	// Check if A and D coincide
	case l.start.Eq(other.end, opts...):
		return RelationshipLineSegmentLineSegmentAeqD

	// Check if B and C coincide
	case l.end.Eq(other.start, opts...):
		return RelationshipLineSegmentLineSegmentBeqC

	// Check if B and D coincide
	case l.end.Eq(other.end, opts...):
		return RelationshipLineSegmentLineSegmentBeqD

	}

	// Determine orientations for intersection and collinearity checks
	o1 := Orientation(l.start, l.end, other.start)
	o2 := Orientation(l.start, l.end, other.end)
	o3 := Orientation(other.start, other.end, l.start)
	o4 := Orientation(other.start, other.end, l.end)

	// Non-collinear intersection cases
	if o1 != o2 && o3 != o4 {

		switch {

		// Check if A lies on CD
		case other.ContainsPoint(l.start) && !other.ContainsPoint(l.end):
			return RelationshipLineSegmentLineSegmentAonCD

		// Check if B lies on CD
		case !other.ContainsPoint(l.start) && other.ContainsPoint(l.end):
			return RelationshipLineSegmentLineSegmentBonCD

		// Check if C lies on l
		case l.ContainsPoint(other.start) && !l.ContainsPoint(other.end):
			return RelationshipLineSegmentLineSegmentConAB

		// Check if D lies on l
		case !l.ContainsPoint(other.start) && l.ContainsPoint(other.end):
			return RelationshipLineSegmentLineSegmentDonAB

		// Default case that lines intersect without any "edge cases"
		default:
			return RelationshipLineSegmentLineSegmentIntersects
		}
	}

	// PointsCollinear cases: All orientations are zero
	if o1 == 0 && o2 == 0 && o3 == 0 && o4 == 0 {
		// Check if segments are collinear and disjoint
		if !other.ContainsPoint(l.start) && !other.ContainsPoint(l.end) &&
			!l.ContainsPoint(other.start) && !l.ContainsPoint(other.end) {
			return RelationshipLineSegmentLineSegmentCollinearDisjoint
		}
		// Check if AB is fully contained within CD
		if other.ContainsPoint(l.start) && other.ContainsPoint(l.end) {
			return RelationshipLineSegmentLineSegmentCollinearABinCD
		}
		// Check if CD is fully contained within AB
		if l.ContainsPoint(other.start) && l.ContainsPoint(other.end) {
			return RelationshipLineSegmentLineSegmentCollinearCDinAB
		}
		// Check specific collinear partial overlaps
		if other.ContainsPoint(l.start) {
			return RelationshipLineSegmentLineSegmentCollinearAonCD
		}
		if other.ContainsPoint(l.end) {
			return RelationshipLineSegmentLineSegmentCollinearBonCD
		}
	}

	// If none of the conditions matched, the segments are disjoint
	return RelationshipLineSegmentLineSegmentMiss
}

// RelationshipToPoint determines the spatial relationship of a given [Point]
// to the line segment.
//
// This method evaluates the relationship between the line segment AB and the
// provided [Point] p by internally delegating the evaluation to the point's
// [Point.RelationshipToLineSegment] method. It checks whether the [Point] lies
// on the infinite line, on the segment, coincides with one of the endpoints,
// or is entirely disjoint.
//
// Parameters:
//   - p ([Point][T]): The [Point] to analyze relative to the line segment.
//   - opts: A variadic slice of [Option] functions to customize the behavior of the relationship check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing the point's location relative
//     to the line segment, improving robustness in floating-point calculations.
//
// Returns:
//   - [RelationshipPointLineSegment]: A constant that describes the relationship between the line segment
//     and the point.
//
// Notes:
//   - This function supports epsilon adjustments to account for floating-point precision errors.
//   - The function checks whether the point is entirely disjoint, coincides with one of the endpoints,
//     lies within the segment bounds, or lies on the infinite extension of the segment.
func (l LineSegment[T]) RelationshipToPoint(p Point[T], opts ...Option) RelationshipPointLineSegment {
	return p.RelationshipToLineSegment(l, opts...)
}

// RelationshipToPolyTree determines the spatial relationship of a line segment to a [PolyTree].
//
// The function evaluates whether the line segment:
//   - Intersects any boundary within the [PolyTree].
//   - Is entirely within a solid or hole polygon in the [PolyTree].
//   - Lies entirely outside the [PolyTree].
//
// Parameters:
//   - tree (*[PolyTree][T]): A pointer to the [PolyTree] to evaluate.
//   - opts: A variadic slice of [Option] functions to customize the behavior of the relationship check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing the point's location relative
//     to the line segment, improving robustness in floating-point calculations.
//
// Returns:
//   - [RelationshipLineSegmentPolyTree]: The relationship between the line segment and the [PolyTree].
//
// Behavior:
//   - If the line segment intersects any boundary (considering epsilon), the function immediately returns [PTLRIntersectsBoundary].
//   - If the segment's endpoints are entirely contained within the same polygon, the function returns [PTLRInsideSolid] or [PTLRInsideHole], depending on the polygon type.
//   - The function uses an epsilon tolerance when checking relationships between line segments and polygon edges.
//   - If no stronger relationship is found, the function returns [PTLRMiss], indicating the segment is entirely outside the [PolyTree].
func (l LineSegment[T]) RelationshipToPolyTree(tree *PolyTree[T], opts ...Option) RelationshipLineSegmentPolyTree {

	// as the points in a polytree contour are doubled, we need to also double the input line segment
	lineSegmentDoubled := l.Scale(NewPoint[T](0, 0), 2)

	highestRel := PTLRMiss // Default to outside

	// Iterate through each polygon in the tree
	for poly := range tree.iterPolys {
		// Check each edge of the polygon's contour
		for edge := range poly.contour.iterEdges {
			// Determine relationship between poly contour & line segment
			rel := edge.RelationshipToLineSegment(lineSegmentDoubled, opts...)
			// any intersection
			if rel > RelationshipLineSegmentLineSegmentMiss {
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

// RelationshipToRectangle determines the spatial relationship between a line segment and a [Rectangle].
//
// Parameters:
//   - r ([Rectangle][T]): The [Rectangle] to evaluate the relationship against.
//
// Returns:
//   - [RelationshipLineSegmentRectangle]: A constant describing the relationship between
//     the line segment and the [Rectangle].
//
// Behavior:
//   - This function internally calls the rectangle's [Rectangle.RelationshipToLineSegment]
//     method to determine the relationship.
//   - The relationship is computed based on the position of the line segment relative to
//     the rectangle's edges and vertices.
//
// Notes:
//   - The function relies on the rectangle's ability to evaluate the relationship
//     between itself and the line segment.
func (l LineSegment[T]) RelationshipToRectangle(r Rectangle[T]) RelationshipLineSegmentRectangle {
	return r.RelationshipToLineSegment(l)
}

// Rotate rotates the LineSegment around a given pivot [Point] by a specified angle in radians.
// Optionally, an epsilon threshold can be applied to adjust the precision of the resulting coordinates.
//
// Parameters:
//   - pivot ([Point][T]): The point around which to rotate the line segment.
//   - radians (float64): The rotation angle in radians.
//   - opts: A variadic slice of [Option] functions to customize the behavior of the rotation.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for snapping near-zero or near-integer
//     values in the resulting coordinates to cleaner values, improving robustness.
//
// Behavior:
//   - The function rotates the start and end points of the line segment around the given pivot
//     point by the specified angle using the [Point.Rotate] method.
//   - If [WithEpsilon] is provided, epsilon adjustments are applied to the rotated coordinates to
//     handle floating-point precision errors.
//
// Returns:
//   - LineSegment[float64]: A new line segment representing the rotated position, with floating-point coordinates.
//
// Notes:
//   - Epsilon adjustment is particularly useful when the rotation involves floating-point
//     calculations that could result in minor inaccuracies.
//   - The returned line segment always has float64 coordinates, ensuring precision regardless
//     of the coordinate type of the original line segment.
func (l LineSegment[T]) Rotate(pivot Point[T], radians float64, opts ...Option) LineSegment[float64] {
	newStart := l.start.Rotate(pivot, radians, opts...)
	newEnd := l.end.Rotate(pivot, radians, opts...)
	return NewLineSegment(newStart, newEnd)
}

// Scale scales the line segment by a given factor from a specified reference point.
//
// Parameters:
//   - ref ([Point][T]): The reference point from which the scaling is applied. Using the origin
//     point (0, 0) scales the segment relative to the coordinate system's origin, while specifying
//     a custom reference point scales the segment relative to that point.
//   - factor ([T]): The scaling factor, where a value greater than 1 expands the segment,
//     and a value between 0 and 1 shrinks it.
//
// Behavior:
//   - The function scales both endpoints of the line segment relative to the specified
//     reference point using the [Point.Scale] method.
//   - The scaling operation preserves the relative orientation of the segment.
//
// Returns:
//   - [LineSegment][T]: A new line segment, scaled relative to the specified reference point.
//
// Notes:
//   - Scaling by a factor of 1 will return a line segment identical to the original.
//   - Negative scaling factors will mirror the segment across the reference point
//     and scale its length accordingly.
//   - If the user wishes to shrink the segment (factor < 1), we recommend ensuring
//     the line segment's type is floating-point to avoid precision loss. Use the [LineSegment.AsFloat] method
//     to safely convert the segment to floating-point type before scaling.
func (l LineSegment[T]) Scale(ref Point[T], factor T) LineSegment[T] {
	return NewLineSegment(
		l.start.Scale(ref, factor),
		l.end.Scale(ref, factor),
	)
}

// Start returns the starting point of the line segment.
//
// This function provides access to the starting point of the LineSegment l, typically representing
// the beginning of the segment.
func (l LineSegment[T]) Start() Point[T] {
	return l.start
}

// String returns a formatted string representation of the line segment for debugging and logging purposes.
//
// The string representation includes the coordinates of the start and end points in the format:
// "LineSegment[(x1, y1) -> (x2, y2)]", where (x1, y1) are the coordinates of the start point,
// and (x2, y2) are the coordinates of the end point.
//
// Returns:
//   - string: A string representing the line segment's start and end coordinates.
func (l LineSegment[T]) String() string {
	return fmt.Sprintf("LineSegment[(%v, %v) -> (%v, %v)]", l.start.x, l.start.y, l.end.x, l.end.y)
}

// SubLineSegment subtracts the start and end points of another line segment from this one.
//
// This function performs an element-wise subtraction, where the start and end points
// of the other line segment are subtracted from the corresponding start and end points
// of the current line segment.
//
// Parameters:
//   - other (LineSegment[T]): The line segment to subtract from the current one.
//
// Returns:
//   - LineSegment[T] - A new line segment where each endpoint is the result of the element-wise subtraction.
func (l LineSegment[T]) SubLineSegment(other LineSegment[T]) LineSegment[T] {
	return NewLineSegment(
		l.start.Translate(other.start.Negate()),
		l.end.Translate(other.end.Negate()),
	)
}

// Translate moves the line segment by a specified vector.
//
// This method shifts the line segment's position in the 2D plane by translating
// both its start and end points by the given vector delta. The relative
// orientation and length of the line segment remain unchanged.
//
// Parameters:
//   - delta ([Point][T]): The vector by which to translate the line segment.
//
// Returns:
//   - [LineSegment][T]: A new line segment translated by the specified vector.
//
// Notes:
//   - Translating the line segment effectively adds the delta vector to both
//     the start and end points of the segment.
//   - This operation is equivalent to a uniform shift, maintaining the segment's
//     shape and size while moving it to a new position.
func (l LineSegment[T]) Translate(delta Point[T]) LineSegment[T] {
	return NewLineSegment(
		l.start.Translate(delta),
		l.end.Translate(delta),
	)
}
