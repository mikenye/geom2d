package geom2d

import (
	"fmt"
	"math"
)

// detailedLineSegmentRelationship defines the possible spatial relationships
// between two line segments in a 2D plane.
//
// This enumeration categorizes how two line segments relate to each other based on
// their positions, intersections, collinearity, etc.
type detailedLineSegmentRelationship int8

// Valid values for detailedLineSegmentRelationship
const (
	// lsrCollinearDisjoint indicates that the segments are collinear and do not intersect, overlap, or touch at any point.
	lsrCollinearDisjoint detailedLineSegmentRelationship = iota - 1

	// lsrMiss indicates that the segments are not collinear, disjoint, and do not intersect, overlap, or touch at any point.
	lsrMiss

	// lsrIntersects indicates that the segments intersect at a unique point that is not an endpoint.
	lsrIntersects

	// lsrAeqC indicates that point A of segment AB coincides with Point C of segment CD.
	lsrAeqC

	// lsrAeqD indicates that point A of segment AB coincides with Point D of segment CD.
	lsrAeqD

	// lsrBeqC indicates that the endpoint of segment AB coincides with Point C of segment CD.
	lsrBeqC

	// lsrBeqD indicates that the endpoint of segment AB coincides with Point D of segment CD.
	lsrBeqD

	// lsrAonCD indicates that point A lies on LineSegment CD.
	lsrAonCD

	// lsrBonCD indicates that the endpoint lies on LineSegment CD.
	lsrBonCD

	// lsrConAB indicates that point C lies on LineSegment AB.
	lsrConAB

	// lsrDonAB indicates that point D lies on LineSegment AB.
	lsrDonAB

	// lsrCollinearAonCD indicates that point A lies on LineSegment CD (partial overlap), and the line segments are collinear.
	lsrCollinearAonCD

	// lsrCollinearBonCD indicates that the endpoint lies on LineSegment CD (partial overlap), and the line segments are collinear.
	lsrCollinearBonCD

	// lsrCollinearABinCD indicates that segment AB is fully contained within segment CD.
	lsrCollinearABinCD

	// lsrCollinearCDinAB indicates that segment CD is fully contained within segment AB.
	lsrCollinearCDinAB

	// lsrCollinearEqual indicates that the segments AB and CD are exactly equal, sharing both endpoints in the same locations.
	lsrCollinearEqual
)

// String provides a string representation of a [detailedLineSegmentRelationship] value.
//
// This method converts the [detailedLineSegmentRelationship] enum value into a human-readable string,
// allowing for easier debugging and output interpretation. Each enum value maps to its corresponding
// string name.
//
// Returns:
//   - string: A string representation of the [detailedLineSegmentRelationship].
//
// Panics:
//   - If the enum value is not recognized, the method panics with an error indicating
//     an unsupported [detailedLineSegmentRelationship] value.
func (r detailedLineSegmentRelationship) String() string {
	switch r {
	case lsrCollinearDisjoint:
		return "lsrCollinearDisjoint"
	case lsrMiss:
		return "lsrMiss"
	case lsrIntersects:
		return "lsrIntersects"
	case lsrAeqC:
		return "lsrAeqC"
	case lsrAeqD:
		return "lsrAeqD"
	case lsrBeqC:
		return "lsrBeqC"
	case lsrBeqD:
		return "lsrBeqD"
	case lsrAonCD:
		return "lsrAonCD"
	case lsrBonCD:
		return "lsrBonCD"
	case lsrConAB:
		return "lsrConAB"
	case lsrDonAB:
		return "lsrDonAB"
	case lsrCollinearAonCD:
		return "lsrCollinearAonCD"
	case lsrCollinearBonCD:
		return "lsrCollinearBonCD"
	case lsrCollinearABinCD:
		return "lsrCollinearABinCD"
	case lsrCollinearCDinAB:
		return "lsrCollinearCDinAB"
	case lsrCollinearEqual:
		return "lsrCollinearEqual"
	default:
		panic(fmt.Errorf("unsupported detailedLineSegmentRelationship"))
	}
}

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
// This constructor function initializes a [LineSegment] with the specified starting and ending points.
// The generic type parameter "T" must satisfy the [SignedNumber] constraint, allowing various numeric types
// (such as int or float64) to be used for the segment’s coordinates.
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

// AsFloat32 converts the line segment to a LineSegment[float32] type.
//
// This function converts both endpoints of the LineSegment l to [Point][float32]
// values, creating a new line segment with floating-point coordinates.
// It is useful for precise calculations where floating-point accuracy is needed.
//
// Returns:
//   - LineSegment[float32] - The line segment with both endpoints converted to float32.
func (l LineSegment[T]) AsFloat32() LineSegment[float32] {
	return NewLineSegment(l.start.AsFloat32(), l.end.AsFloat32())
}

// AsFloat64 converts the line segment to a LineSegment[float64] type.
//
// This function converts both endpoints of the LineSegment l to [Point][float64]
// values, creating a new line segment with floating-point coordinates.
// It is useful for precise calculations where floating-point accuracy is needed.
//
// Returns:
//   - LineSegment[float64] - The line segment with both endpoints converted to float64.
func (l LineSegment[T]) AsFloat64() LineSegment[float64] {
	return NewLineSegment(l.start.AsFloat64(), l.end.AsFloat64())
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

// BoundingBox computes the smallest axis-aligned [Rectangle] that fully contains the LineSegment.
//
// Returns:
//   - [Rectangle][T]: A [Rectangle] defined by the opposite corners of the LineSegment.
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

	start := l.start.AsFloat64()
	end := l.end.AsFloat64()

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

// detailedRelationshipToLineSegment determines the spatial relationship between two line segments, l and other.
//
//   - Let A = l.Start()
//   - Let B = l.End()
//   - Let C = other.Start()
//   - Let D = other.End()
//
// This function evaluates the relationship between two line segments, AB and CD, by checking for
// endpoint coincidences, intersections, collinear relationships, and containment. It returns a
// [detailedLineSegmentRelationship] constant that describes the exact relationship between the segments,
// such as intersection, partial overlap, or full containment.
//
// The output constants may reference A, B, C, or D for brevity.
//
// Parameters:
//   - other (LineSegment[T]): The line segment to compare with l.
//   - opts: A variadic slice of [Option] functions to customize the behavior of the relationship check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing points and collinearity calculations,
//     improving robustness against floating-point precision errors.
//
// Behavior:
//   - The function first checks if the two line segments are exactly equal (or approximately equal if an epsilon is provided).
//   - It evaluates endpoint coincidences, collinearity, intersection, and containment using orientation tests,
//     point-on-segment checks, and direct comparisons.
//   - If [WithEpsilon] is provided, epsilon adjustments are applied to point comparisons, collinearity checks, and
//     point-on-segment tests to handle floating-point imprecision.
//
// Returns:
//   - [detailedLineSegmentRelationship]: A constant describing the relationship between segments AB and CD.
//
// Notes:
//   - Epsilon adjustment is particularly useful when working with floating-point coordinates, where small
//     precision errors might lead to incorrect results.
func (l LineSegment[T]) detailedRelationshipToLineSegment(other LineSegment[T], opts ...Option) detailedLineSegmentRelationship {

	// Check if segments are exactly equal
	if (l.start.Eq(other.start, opts...) && l.end.Eq(other.end, opts...)) || (l.start.Eq(other.end, opts...) && l.end.Eq(other.start, opts...)) {
		return lsrCollinearEqual
	}

	switch {

	// Check if A and C coincide
	case l.start.Eq(other.start, opts...):
		return lsrAeqC

	// Check if A and D coincide
	case l.start.Eq(other.end, opts...):
		return lsrAeqD

	// Check if B and C coincide
	case l.end.Eq(other.start, opts...):
		return lsrBeqC

	// Check if B and D coincide
	case l.end.Eq(other.end, opts...):
		return lsrBeqD

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
			return lsrAonCD

		// Check if B lies on CD
		case !other.ContainsPoint(l.start) && other.ContainsPoint(l.end):
			return lsrBonCD

		// Check if C lies on l
		case l.ContainsPoint(other.start) && !l.ContainsPoint(other.end):
			return lsrConAB

		// Check if D lies on l
		case !l.ContainsPoint(other.start) && l.ContainsPoint(other.end):
			return lsrDonAB

		// Default case that lines intersect without any "edge cases"
		default:
			return lsrIntersects
		}
	}

	// PointsCollinear cases: All orientations are zero
	if o1 == 0 && o2 == 0 && o3 == 0 && o4 == 0 {
		// Check if segments are collinear and disjoint
		if !other.ContainsPoint(l.start) && !other.ContainsPoint(l.end) &&
			!l.ContainsPoint(other.start) && !l.ContainsPoint(other.end) {
			return lsrCollinearDisjoint
		}
		// Check if AB is fully contained within CD
		if other.ContainsPoint(l.start) && other.ContainsPoint(l.end) {
			return lsrCollinearABinCD
		}
		// Check if CD is fully contained within AB
		if l.ContainsPoint(other.start) && l.ContainsPoint(other.end) {
			return lsrCollinearCDinAB
		}
		// Check specific collinear partial overlaps
		if other.ContainsPoint(l.start) {
			return lsrCollinearAonCD
		}
		if other.ContainsPoint(l.end) {
			return lsrCollinearBonCD
		}
	}

	// If none of the conditions matched, the segments are disjoint
	return lsrMiss
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
	ABf, CDf := l.AsFloat64(), other.AsFloat64()

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
//   - other (LineSegment[T]): The line segment to compare with the current line segment.
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
//
// todo: update doc comments and reorder
func (l LineSegment[T]) IntersectionGeometry(other LineSegment[T], opts ...Option) LineSegmentIntersectionResult[float64] {

	// Apply options with defaults
	options := applyOptions(geomOptions{epsilon: 1e-10}, opts...)

	// Define segment endpoints for AB (l) and CD (other)
	A, B := l.start.AsFloat64(), l.end.AsFloat64()
	C, D := other.start.AsFloat64(), other.end.AsFloat64()

	// Calculate the direction vectors
	dir1 := B.Translate(A.Negate())
	dir2 := D.Translate(C.Negate())

	// Calculate the determinants
	denominator := dir1.CrossProduct(dir2)

	// Handle collinear case (denominator == 0)
	if denominator == 0 {
		// Check if the segments are collinear
		AC := C.Translate(A.Negate())
		if AC.CrossProduct(dir1) != 0 {
			return LineSegmentIntersectionResult[float64]{IntersectionType: IntersectionNone} // Parallel but not collinear
		}

		// Check overlap by projecting points onto the line
		tStart := (C.Translate(A.Negate())).DotProduct(dir1) / dir1.DotProduct(dir1)
		tEnd := (D.Translate(A.Negate())).DotProduct(dir1) / dir1.DotProduct(dir1)

		// Ensure tStart < tEnd for consistency
		if tStart > tEnd {
			tStart, tEnd = tEnd, tStart
		}

		// Check for overlap
		tOverlapStart := max(0.0, tStart)
		tOverlapEnd := min(1.0, tEnd)

		if tOverlapStart > tOverlapEnd {
			return LineSegmentIntersectionResult[float64]{IntersectionType: IntersectionNone} // No overlap
		}

		// Calculate the overlapping segment
		overlapStart := RoundPointToEpsilon(A.Translate(dir1.Scale(NewPoint[float64](0, 0), tOverlapStart)), options.epsilon)
		overlapEnd := RoundPointToEpsilon(A.Translate(dir1.Scale(NewPoint[float64](0, 0), tOverlapEnd)), options.epsilon)

		return LineSegmentIntersectionResult[float64]{
			IntersectionType:    IntersectionSegment,
			IntersectionSegment: NewLineSegment(overlapStart, overlapEnd),
		}
	}

	// Calculate parameters t and u for non-collinear case
	AC := C.Translate(A.Negate())
	tNumerator := AC.CrossProduct(dir2)
	uNumerator := AC.CrossProduct(dir1)

	t := tNumerator / denominator
	u := uNumerator / denominator

	// Check if intersection occurs within the segment bounds
	if t < 0 || t > 1 || u < 0 || u > 1 {
		return LineSegmentIntersectionResult[float64]{IntersectionType: IntersectionNone} // Intersection is outside the segments
	}

	// Calculate the intersection point
	intersection := RoundPointToEpsilon(A.Translate(dir1.Scale(NewPoint[float64](0, 0), t)), options.epsilon)
	return LineSegmentIntersectionResult[float64]{
		IntersectionType:  IntersectionPoint,
		IntersectionPoint: intersection,
	}
}

// IntersectsLineSegment checks whether there is any intersection or overlap between LineSegment l and LineSegment other.
//
// This function returns true if segments l and other have an intersecting spatial relationship, such as intersection,
// overlap, containment, or endpoint coincidence.
//
// Parameters:
//   - other (LineSegment[T]): The line segment to compare with l.
//
// Returns:
//   - bool: Returns true if l and other intersect, overlap, or share any form of intersecting relationship, and false if they are completely disjoint.
func (l LineSegment[T]) IntersectsLineSegment(other LineSegment[T]) bool {
	return l.detailedRelationshipToLineSegment(other) > lsrMiss
}

// Length calculates the length of the line segment, optionally using an epsilon threshold
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
func (l LineSegment[T]) Length(opts ...Option) float64 {
	return l.start.DistanceToPoint(l.end, opts...)
}

// Normalize returns a "normalized" version of the line segment where the start point
// is guaranteed to be the leftmost-lowest point. This ensures that:
// - The point with the smallest X-coordinate is the start point.
// - If the X-coordinates are equal, the point with the smallest Y-coordinate is the start point.
//
// This normalization is useful for algorithms that require consistent ordering of line segments,
// such as the Bentley-Ottmann sweep line algorithm or Boolean operations on polygons.
//
// Returns a new LineSegment with the points swapped if necessary.
func (l LineSegment[T]) Normalize() LineSegment[T] {
	// if start point is not leftest-lowest, swap points
	if l.start.x > l.end.x || (l.start.x == l.end.x && l.start.y > l.end.y) {
		return NewLineSegment[T](l.end, l.start)
	}
	// else, return original point
	return l
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
func (l LineSegment[T]) Reflect(axis ReflectionAxis, line ...LineSegment[float64]) LineSegment[float64] {
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
			return l.AsFloat64()
		}
	default:
		// Invalid axis, return the line segment unchanged
		return l.AsFloat64()
	}

	// Return a new line segment with reflected points
	return NewLineSegment(startReflected, endReflected)
}

// RelationshipToCircle determines the spatial relationship between the current LineSegment and a Circle.
//
// This function evaluates whether the LineSegment:
//   - Intersects the circle at any point.
//   - Lies entirely within the circle.
//   - Lies entirely outside the circle.
//
// Parameters:
//   - c ([Circle][T]): The circle to compare with the line segment.
//   - opts: A variadic slice of [Option] functions to customize the behavior of the relationship check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing distances, improving robustness against
//     floating-point precision errors.
//
// Behavior:
//   - The function calculates the distance between the circle's center and the endpoints of the line segment,
//     as well as the shortest distance between the circle's center and the line segment itself.
//   - If any of these distances are exactly equal to the circle's radius, the function determines that the
//     line segment intersects the circle.
//   - If both endpoints of the line segment lie within the circle, the function determines that the line segment
//     is contained within the circle.
//   - Otherwise, the function determines that the line segment lies entirely outside the circle.
//
// Returns:
//
// [Relationship]: A constant indicating the relationship of the line segment to the circle. Possible values are:
//   - [RelationshipDisjoint]: The line segment lies entirely outside the circle.
//   - [RelationshipIntersection]: The line segment intersects the circle at one or more points.
//   - [RelationshipContainedBy]: The line segment lies entirely within the circle.
//
// Notes:
//   - Epsilon adjustment is particularly useful when working with floating-point coordinates,
//     where minor precision errors could otherwise lead to incorrect results.
func (l LineSegment[T]) RelationshipToCircle(c Circle[T], opts ...Option) Relationship {
	circleFloat := c.AsFloat64()

	distanceCircleCenterToLineSegment := c.center.DistanceToLineSegment(l, opts...)
	distanceCircleCenterToLineSegmentStart := c.center.DistanceToPoint(l.start, opts...)
	distanceCircleCenterToLineSegmentEnd := c.center.DistanceToPoint(l.end, opts...)

	// check for intersection
	if distanceCircleCenterToLineSegmentStart == circleFloat.radius ||
		distanceCircleCenterToLineSegmentEnd == circleFloat.radius ||
		distanceCircleCenterToLineSegment == circleFloat.radius {
		return RelationshipIntersection
	}
	if (distanceCircleCenterToLineSegmentStart > circleFloat.radius || distanceCircleCenterToLineSegmentEnd > circleFloat.radius) &&
		distanceCircleCenterToLineSegment <= circleFloat.radius {
		return RelationshipIntersection
	}

	// check for containment
	if distanceCircleCenterToLineSegmentStart < circleFloat.radius &&
		distanceCircleCenterToLineSegmentEnd < circleFloat.radius {
		return RelationshipContainedBy
	}

	return RelationshipDisjoint
}

// RelationshipToLineSegment determines the high-level spatial relationship between two line segments.
//
// It categorizes the relationship as:
//   - Disjoint (no intersection or overlap).
//   - Intersection (the segments intersect at one or more points).
//   - Equal (both segments are collinear and share the same endpoints).
//
// Parameters:
//   - other ([LineSegment][T]): The line segment to compare against the current one.
//   - opts: A variadic slice of [Option] functions to customize the behavior of the relationship check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing points and collinearity calculations,
//     allowing for robust handling of floating-point precision errors.
//
// Behavior:
//
// The detailed relationship mapped to a [Relationship] constant:
//   - [RelationshipDisjoint]: If the segments are collinear but disjoint, or if they miss entirely.
//   - [RelationshipIntersection]: The segments intersect at one or more points.
//   - [RelationshipEqual]: If the segments are collinear and exactly equal (share the same endpoints).
//
// Returns:
//   - [Relationship]: A constant describing the high-level relationship between the two line segments.
//
// Notes:
//   - The use of epsilon adjustments ensures robustness against floating-point imprecision.
func (l LineSegment[T]) RelationshipToLineSegment(other LineSegment[T], opts ...Option) Relationship {
	rel := l.detailedRelationshipToLineSegment(other, opts...)
	switch rel {
	case lsrCollinearDisjoint, lsrMiss:
		return RelationshipDisjoint
	case lsrCollinearEqual:
		return RelationshipEqual
	default:
		return RelationshipIntersection
	}
}

// RelationshipToPoint determines the high-level spatial relationship between a line segment and a point.
//
// This function evaluates how a line segment relates to a point by delegating the computation to the
// [Point.RelationshipToLineSegment] method. The relationship is determined based on whether the point
// lies on the segment, coincides with an endpoint, or is disjoint from the segment.
//
// Parameters:
//   - p ([Point][T]): The point to compare against the line segment.
//   - opts: A variadic slice of [Option] functions to customize the behavior of the relationship check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing distances and collinearity calculations,
//     ensuring robust handling of floating-point precision errors.
//
// Behavior:
//   - The function calls the [Point.RelationshipToLineSegment] method for computation.
//   - The returned relationship constant describes whether the point is on the line segment ([RelationshipIntersection]), or
//     disjoint from the line segment ([RelationshipDisjoint]).
//
// Returns:
//   - [Relationship]: A constant describing the relationship between the line segment and the point.
//
// Notes:
//   - Epsilon adjustment ensures robustness against floating-point imprecision.
func (l LineSegment[T]) RelationshipToPoint(p Point[T], opts ...Option) Relationship {
	return p.RelationshipToLineSegment(l, opts...)
}

// RelationshipToPolyTree determines the relationship between a line segment and each polygon in a [PolyTree].
//
// This function evaluates how a line segment relates to the polygons in the tree, computing whether the segment
// is disjoint, intersects any edge, or is fully contained within a polygon. The result is returned as a map,
// where each key is a pointer to a polygon in the [PolyTree], and the value is the relationship between the
// line segment and that polygon.
//
// Parameters:
//   - pt (*[PolyTree][T]): The [PolyTree] to analyze for relationships.
//   - opts: A variadic slice of [Option] functions to customize the behavior of the relationship check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing distances and collinearity calculations,
//     ensuring robust handling of floating-point precision errors.
//
// Behavior:
//
// The function first checks if the line segment intersects any edge of each polygon.
//   - If an intersection or equality is found, the relationship is set to [RelationshipIntersection].
//   - If the segment's start and end points are both contained within a polygon, the relationship is set to [RelationshipContainedBy].
//   - If neither of the above conditions is satisfied, the relationship defaults to [RelationshipDisjoint].
//
// Returns:
//   - map[*PolyTree[T]]Relationship: A map where the keys are polygons in the [PolyTree] and the values are their
//     relationships with the line segment.
//
// Notes:
//   - Epsilon adjustment ensures robustness against floating-point imprecision.
//   - The function assumes polygons in the [PolyTree] have doubled points for accurate containment checks.
func (l LineSegment[T]) RelationshipToPolyTree(pt *PolyTree[T], opts ...Option) map[*PolyTree[T]]Relationship {
	lDoubled := NewLineSegment[T](NewPoint[T](l.start.x*2, l.start.y*2), NewPoint[T](l.end.x*2, l.end.y*2))
	output := make(map[*PolyTree[T]]Relationship, pt.Len())
RelationshipToPolyTreeIterPolys:
	for poly := range pt.Nodes {

		// check for intersection
		for edge := range poly.contour.iterEdges {
			rel := lDoubled.RelationshipToLineSegment(edge, opts...)
			if rel == RelationshipIntersection || rel == RelationshipEqual {
				output[poly] = RelationshipIntersection
				continue RelationshipToPolyTreeIterPolys
			}
		}

		// check for containment
		lineStartInPoly := poly.contour.isPointInside(lDoubled.start)
		lineEndInPoly := poly.contour.isPointInside(lDoubled.end)
		if lineStartInPoly && lineEndInPoly {
			output[poly] = RelationshipContainedBy
			continue RelationshipToPolyTreeIterPolys
		}

		// otherwise, disjoint
		output[poly] = RelationshipDisjoint
	}
	return output
}

// RelationshipToRectangle determines the high-level spatial relationship between a line segment and a rectangle.
//
// This function evaluates how a line segment relates to a rectangle, considering possible intersections,
// containment, or disjoint relationships. The function iterates through the rectangle's edges to check for
// intersections with the line segment and verifies whether the line segment is entirely contained within
// the rectangle.
//
// Parameters:
//   - r ([Rectangle][T]): The rectangle to compare against the line segment.
//   - opts: A variadic slice of [Option] functions to customize the behavior of the relationship check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing points and collinearity calculations,
//     ensuring robust handling of floating-point precision errors.
//
// Behavior:
//
// The function checks each edge of the rectangle against the line segment:
//   - If any edge intersects or is equal to the line segment, the function returns [RelationshipIntersection].
//   - If both endpoints of the line segment are contained within the rectangle, the function returns [RelationshipContainedBy].
//   - If no intersection or containment is found, the function returns [RelationshipDisjoint].
//
// Returns:
//   - [Relationship]: A constant describing the relationship between the line segment and the rectangle.
//
// Notes:
//   - Epsilon adjustment ensures robustness against floating-point imprecision.
func (l LineSegment[T]) RelationshipToRectangle(r Rectangle[T], opts ...Option) Relationship {
	for _, edge := range r.Edges() {
		rel := edge.RelationshipToLineSegment(l, opts...)
		if rel == RelationshipIntersection || rel == RelationshipEqual {
			return RelationshipIntersection
		}
	}
	if r.ContainsPoint(l.start) && r.ContainsPoint(l.end) {
		return RelationshipContainedBy
	}
	return RelationshipDisjoint
}

// Rotate rotates the LineSegment around a given pivot [Point] by a specified angle in radians counterclockwise.
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

// RoundToEpsilon returns a new LineSegment where the coordinates of both the start
// and end points are rounded to the nearest multiple of the given epsilon.
//
// Parameters:
//   - epsilon: The value to which the coordinates should be rounded.
//
// Returns:
//
//	A new LineSegment with rounded coordinates.
//
// Notes:
//   - The epsilon value must be greater than 0. If it is 0 or negative,
//     the function will panic.
func (l LineSegment[T]) RoundToEpsilon(epsilon float64) LineSegment[float64] {
	return NewLineSegment(
		RoundPointToEpsilon(l.start.AsFloat64(), epsilon),
		RoundPointToEpsilon(l.end.AsFloat64(), epsilon),
	)
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
//     the line segment's type is floating-point to avoid precision loss. Use the [LineSegment.AsFloat64] method
//     to safely convert the segment to floating-point type before scaling.
func (l LineSegment[T]) Scale(ref Point[T], factor T) LineSegment[T] {
	return NewLineSegment(
		l.start.Scale(ref, factor),
		l.end.Scale(ref, factor),
	)
}

// Slope calculates the slope of the line segment.
//
// The slope is calculated as the change in y-coordinates (dy) divided by
// the change in x-coordinates (dx) of the line segment. This function
// returns the slope as a float64 and a boolean indicating whether the
// slope is defined.
//
// Returns:
//   - (float64, true): The calculated slope if the line segment is not vertical.
//   - (0, false): Indicates the slope is undefined (the line segment is vertical).
func (l LineSegment[T]) Slope() (float64, bool) {
	dx := float64(l.end.x - l.start.x)
	dy := float64(l.end.y - l.start.y)

	if dx == 0 {
		return 0, false // Vertical line, slope undefined
	}
	return dy / dx, true
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

// XAtY calculates the X-coordinate of the line segment at the given Y-coordinate.
//
// The function determines the X-coordinate of the intersection point between the line
// segment and the horizontal line at the specified Y. It handles vertical, horizontal,
// and diagonal line segments.
//
// Parameters:
//   - y (float64): The Y-coordinate for which the corresponding X-coordinate is to be calculated.
//
// Returns:
//   - float64: The X-coordinate corresponding to the given Y, if it lies within the bounds of the segment.
//   - bool: A boolean indicating whether the given Y is within the vertical range of the segment.
//
// Behavior:
//   - If the line segment is vertical (constant X), the function returns the constant X-coordinate
//     if the Y-coordinate is within the segment's range, and false otherwise.
//   - If the line segment is horizontal (constant Y), the function returns false unless the Y-coordinate
//     matches the segment's Y-coordinate.
//   - For diagonal line segments, the function computes the X-coordinate using the line equation
//     and returns true if the calculated X lies within the segment's bounds.
func (l LineSegment[T]) XAtY(y float64) (float64, bool) {
	// Handle vertical line segment (undefined slope)
	lf := l.AsFloat64()
	if lf.start.x == lf.end.x {
		if y >= lf.start.y && y <= lf.end.y || y >= lf.end.y && y <= lf.start.y {
			return lf.start.x, true
		}
		return 0, false // Y is out of bounds
	}

	// Handle horizontal line segment
	if lf.start.y == lf.end.y {
		if y == lf.start.y {
			return lf.start.x, true // Any X on the segment is valid
		}
		return 0, false // Y is out of bounds
	}

	// Calculate the X value using the line equation
	slope := (lf.end.y - lf.start.y) / (lf.end.x - lf.start.x)
	intercept := lf.start.y - slope*lf.start.x // y = mx + b -> b = y - mx

	// Rearrange y = mx + b to solve for x: x = (y - b) / m
	x := (y - intercept) / slope

	// Check if the calculated X is within the segment bounds
	if (x >= lf.start.x && x <= lf.end.x) || (x >= lf.end.x && x <= lf.start.x) {
		return x, true
	}
	return 0, false // X is out of bounds
}

// YAtX calculates the Y-coordinate of the line segment at the given X-coordinate.
//
// The function determines the Y-coordinate of the intersection point between the line
// segment and the vertical line at the specified X. It handles vertical, horizontal,
// and diagonal line segments.
//
// Parameters:
//   - x (float64): The X-coordinate for which the corresponding Y-coordinate is to be calculated.
//
// Returns:
//   - float64: The Y-coordinate corresponding to the given X, if it lies within the bounds of the segment.
//   - bool: A boolean indicating whether the given X is within the horizontal range of the segment.
//
// Behavior:
//   - If the line segment is vertical (constant X), the function returns the Y-coordinate
//     of the line segment if the X-coordinate matches the segment's X-coordinate, and false otherwise.
//   - If the line segment is horizontal (constant Y), the function returns the constant Y-coordinate
//     if the X-coordinate is within the segment's range, and false otherwise.
//   - For diagonal line segments, the function computes the Y-coordinate using the line equation
//     and returns true if the calculated Y lies within the segment's bounds.
func (l LineSegment[T]) YAtX(x float64) (float64, bool) {
	lf := l.AsFloat64()
	// Handle vertical line segment
	if lf.start.x == lf.end.x {
		if x == lf.start.x {
			return lf.start.y, true // Any Y within the segment's range is valid
		}
		return 0, false // X is out of bounds
	}

	// Handle horizontal line segment
	if lf.start.y == lf.end.y {
		if x >= lf.start.x && x <= lf.end.x || x >= lf.end.x && x <= lf.start.x {
			return lf.start.y, true
		}
		return 0, false // X is out of bounds
	}

	// Calculate Y value using the line equation: y = mx + b
	slope := (lf.end.y - lf.start.y) / (lf.end.x - lf.start.x)
	intercept := lf.start.y - slope*lf.start.x // y = mx + b -> b = y - mx

	// Calculate y = mx + b
	y := slope*x + intercept

	// Check if X is within the segment bounds
	if (x >= lf.start.x && x <= lf.end.x) || (x >= lf.end.x && x <= lf.start.x) {
		return y, true
	}
	return 0, false // X is out of bounds
}

// Bresenham generates all the integer points along the line segment using
// Bresenham's line algorithm. It is an efficient way to rasterize a line
// in a grid or pixel-based system.
//
// The function is designed to be used with a for-loop, and thus takes a callback yield that processes each point.
// If the callback returns false at any point (if the calling for-loop is terminated, for example), the function
// halts further generation.
//
// Example use cases include:
// - Rendering lines in graphics applications.
// - Generating grid points for pathfinding.
//
// Parameters:
//   - yield (func([Point][int]) bool): A function that processes each generated point.
//     Returning false will stop further point generation.
//
// Note: This method requires integer-type coordinates for the line segment.
func (l LineSegment[int]) Bresenham(yield func(Point[int]) bool) {

	var x1, x2, y1, y2, dx, dy, sx, sy int

	x1 = l.start.x
	x2 = l.end.x
	y1 = l.start.y
	y2 = l.end.y

	// Calculate absolute deltas
	dx = abs(x2 - x1)
	dy = abs(y2 - y1)

	// Determine the direction of the increments
	sx = 1
	if x1 > x2 {
		sx = -1
	}
	sy = 1
	if y1 > y2 {
		sy = -1
	}

	// Bresenham's algorithm
	err := dx - dy
	for {
		if !yield(NewPoint(x1, y1)) {
			return
		}

		// Break the loop if we've reached the end point
		if x1 == x2 && y1 == y2 {
			return
		}

		// Calculate the error
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

type LineSegmentIntersectionType uint8

const (
	IntersectionNone LineSegmentIntersectionType = iota
	IntersectionPoint
	IntersectionSegment
)

type LineSegmentIntersectionResult[T SignedNumber] struct {
	IntersectionType    LineSegmentIntersectionType // Type of intersection
	IntersectionPoint   Point[T]                    // Valid if Type == IntersectionPoint
	IntersectionSegment LineSegment[T]              // Valid if Type == IntersectionSegment
}
