package _old

import (
	"fmt"
	"github.com/mikenye/geom2d/types"
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
	return NewFromPoints(
		l.start.Translate(other.start),
		l.end.Translate(other.end),
	)
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
