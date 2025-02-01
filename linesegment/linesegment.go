package linesegment

import (
	"fmt"
	"github.com/mikenye/geom2d/numeric"
	"github.com/mikenye/geom2d/options"
	"github.com/mikenye/geom2d/point"
	"github.com/mikenye/geom2d/types"
	"math"
)

// LineSegment represents a line segment in a 2D space, defined by two endpoints, the start [point.Point] and end [point.Point].
//
// The generic type parameter T must satisfy the [types.SignedNumber] constraint, allowing the segment
// to use various numeric types such as int or float64 for its coordinates.
type LineSegment[T types.SignedNumber] struct {
	start point.Point[T]
	end   point.Point[T]
}

// New creates a new LineSegment with the specified start and end x and y coordinates.
//
// This constructor function initializes a [LineSegment] with the specified starting and ending points.
// The generic type parameter "T" must satisfy the [types.SignedNumber] constraint, allowing various numeric types
// (such as int or float64) to be used for the segment’s coordinates.
//
// Parameters:
//   - x1,y1 (T): The starting point of the LineSegment.
//   - x2,y2 (T): The ending point of the LineSegment.
//
// Returns:
//   - LineSegment[T] - A new line segment defined by the start and end points.
func New[T types.SignedNumber](x1, y1, x2, y2 T) LineSegment[T] {
	return LineSegment[T]{
		start: point.New[T](x1, y1),
		end:   point.New[T](x2, y2),
	}
}

// NewFromPoints creates a new LineSegment from two endpoints, a start [point.Point] and end [point.Point].
//
// This constructor function initializes a [LineSegment] with the specified starting and ending points.
// The generic type parameter "T" must satisfy the [types.SignedNumber] constraint, allowing various numeric types
// (such as int or float64) to be used for the segment’s coordinates.
//
// Parameters:
//   - start ([point.Point][T]): The starting [point.Point] of the LineSegment.
//   - end ([point.Point][T]): The ending [point.Point] of the LineSegment.
//
// Returns:
//   - LineSegment[T] - A new line segment defined by the start and end points.
func NewFromPoints[T types.SignedNumber](start, end point.Point[T]) LineSegment[T] {
	return LineSegment[T]{
		start: start,
		end:   end,
	}
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
	return NewFromPoints(l.start.AsFloat32(), l.end.AsFloat32())
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
	return NewFromPoints(l.start.AsFloat64(), l.end.AsFloat64())
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
	return NewFromPoints(l.start.AsInt(), l.end.AsInt())
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
	return NewFromPoints(l.start.AsIntRounded(), l.end.AsIntRounded())
}

// Bresenham generates all the integer points along the LineSegment using
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
//   - yield (func([point.Point][int]) bool): A function that processes each generated point.
//     Returning false will stop further point generation.
//
// Note: This method requires integer-type coordinates for the line segment.
func (l LineSegment[int]) Bresenham(yield func(point.Point[int]) bool) {

	var x1, x2, y1, y2, dx, dy, sx, sy int

	x1 = l.start.X()
	x2 = l.end.X()
	y1 = l.start.Y()
	y2 = l.end.Y()

	// Calculate absolute deltas
	dx = numeric.Abs(x2 - x1)
	dy = numeric.Abs(y2 - y1)

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
		if !yield(point.New(x1, y1)) {
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
func (l LineSegment[T]) Center(opts ...options.GeometryOptionsFunc) point.Point[float64] {
	// Apply geomOptions with defaults
	geoOpts := options.ApplyGeometryOptions(options.GeometryOptions{Epsilon: 0}, opts...)

	start := l.start.AsFloat64()
	end := l.end.AsFloat64()

	midX := (start.X() + end.X()) / 2
	midY := (start.Y() + end.Y()) / 2

	// Apply epsilon if specified
	if geoOpts.Epsilon > 0 {
		midX = numeric.SnapToEpsilon(midX, geoOpts.Epsilon)
		midY = numeric.SnapToEpsilon(midY, geoOpts.Epsilon)
	}

	return point.New[float64](midX, midY)
}

// ContainsPoint determines whether the given [point.Point] lies on the LineSegment.
//
// This method calculates the shortest distance from the given point to the LineSegment
// using the DistanceToPoint method. If the distance is zero (considering an epsilon threshold),
// the point is determined to be on the segment.
//
// Parameters:
//   - p ([point.Point][T]): The [Point] to test against the LineSegment.
//   - opts: A variadic slice of [options.GeometryOptionsFunc] to customize the epsilon for
//     floating-point precision adjustments.
//
// Returns:
//   - bool: true if the [point.Point] lies on the LineSegment, false otherwise.
//
// Notes:
//   - This function uses the DistanceToPoint method to compute the distance.
//   - Floating-point precision issues are handled using the epsilon parameter if provided in opts.
//   - The [point.Point] must also be within the bounding box defined by the segment endpoints to return true.
func (l LineSegment[T]) ContainsPoint(p point.Point[T], opts ...options.GeometryOptionsFunc) bool {

	geoOpts := options.ApplyGeometryOptions(options.GeometryOptions{Epsilon: 0}, opts...)

	lf := l.AsFloat64()
	pf := p.AsFloat64()

	// Early bounding box check
	if pf.X() < min(lf.start.X(), lf.end.X())-geoOpts.Epsilon ||
		pf.X() > max(lf.start.X(), lf.end.X())+geoOpts.Epsilon ||
		pf.Y() < min(lf.start.Y(), lf.end.Y())-geoOpts.Epsilon ||
		pf.Y() > max(lf.start.Y(), lf.end.Y())+geoOpts.Epsilon {
		return false
	}

	// Check distance to the segment
	d := numeric.SnapToEpsilon(l.DistanceToPoint(p, opts...), geoOpts.Epsilon)

	// if d == 0 then point is on the segment
	return d == 0
}

// DistanceToLineSegment calculates the minimum distance between two line segments, l and other.
//
// If the segments intersect or touch at any point, the function immediately returns 0, as the distance is effectively zero.
// Otherwise, it calculates the shortest distance by considering:
//  1. The distances between endpoints of one segment and the other segment.
//  2. The distances from endpoints of one segment to the perpendicular projections onto the other segment.
//
// Parameters:
//   - other (LineSegment[T]): The line segment to compare with l.
//   - opts: A variadic slice of [options.GeometryOptionsFunc] functions to customize the calculation behavior.
//     Common options include:
//   - [options.WithEpsilon](epsilon float64): Specifies a tolerance to handle small floating-point deviations
//     in distance calculations and ensure near-zero results are treated as zero.
//
// Returns:
//   - float64: The minimum distance between the two line segments. If the segments intersect or touch, this value is 0.
//
// Behavior:
//
// First, the function checks whether the two segments intersect or touch using the Intersection method. If so, the distance is 0.
//
// For non-intersecting segments, the function calculates distances using:
//  1. Endpoints of one segment to endpoints of the other.
//  2. Endpoints of one segment to the perpendicular projections on the other segment.
//
// The smallest of these distances is returned.
//
// Notes:
//   - This function converts the line segments to float64 precision for robust calculations.
//   - This is a comprehensive calculation suitable for exact and approximate (epsilon-adjusted) distance checks.
func (l LineSegment[T]) DistanceToLineSegment(other LineSegment[T], opts ...options.GeometryOptionsFunc) float64 {
	// If line segments intersect, the distance is zero.
	intersection := l.Intersection(other, opts...)
	if intersection.IntersectionType != IntersectionNone {
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
	updateMinDist(ABf.start.DistanceToPoint(CDf.ProjectPoint(ABf.start), opts...))
	updateMinDist(ABf.end.DistanceToPoint(CDf.ProjectPoint(ABf.end), opts...))
	updateMinDist(CDf.start.DistanceToPoint(ABf.ProjectPoint(CDf.start), opts...))
	updateMinDist(CDf.end.DistanceToPoint(ABf.ProjectPoint(CDf.end), opts...))

	return minDist
}

// DistanceToPoint calculates the orthogonal (shortest) distance from the [LineSegment] l to the [point.Point] p.
// This distance is the length of the perpendicular line from p to the closest point on l.
//
// Parameters:
//   - p ([point.Point][T]): The [point.Point] to which the distance is calculated from [LineSegment] l.
//   - opts: A variadic slice of [options.GeometryOptionsFunc] functions to customize the calculation behavior.
//     [options.WithEpsilon](epsilon float64): Adjusts the result by snapping small floating-point
//     deviations to cleaner values based on the specified epsilon threshold.
//
// Behavior:
//   - The function first computes the projection of p onto the given [LineSegment] l. This is
//     the closest point on l to p, whether it falls within the line segment or on one of its endpoints.
//   - The distance is then calculated as the Euclidean distance from p to the projected point,
//     using the [point.Point.DistanceToPoint] method for precision.
//
// Returns:
//   - float64: The shortest distance between the point p and the line segment l, optionally
//     adjusted based on epsilon if provided.
//
// Notes:
//   - If the point p lies exactly on the line segment, the distance will be zero (or adjusted
//     to zero if within epsilon).
//   - This method ensures precision by converting points to float64 before performing calculations.
func (l LineSegment[T]) DistanceToPoint(p point.Point[T], opts ...options.GeometryOptionsFunc) float64 {
	projectedPoint := l.ProjectPoint(p)
	pf := p.AsFloat64()
	return pf.DistanceToPoint(projectedPoint, opts...)
}

// End returns the ending [point.Point] of the LineSegment.
//
// This function provides access to the ending [point.Point] of the LineSegment l, typically representing
// the endpoint of the segment.
func (l LineSegment[T]) End() point.Point[T] {
	return l.end
}

// Eq checks if two line segments are equal by comparing their start and end points.
// Equality can be evaluated either exactly (default) or approximately using an epsilon threshold.
//
// Parameters:
//   - other (LineSegment[T]): The line segment to compare with the current line segment.
//   - opts: A variadic slice of [options.GeometryOptionsFunc] functions to customize the equality check.
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
func (l LineSegment[T]) Eq(other LineSegment[T], opts ...options.GeometryOptionsFunc) bool {
	return l.start.Eq(other.start, opts...) && l.end.Eq(other.end, opts...)
}

// todo: doc comments, unit test, example
func (l LineSegment[T]) Flip() LineSegment[T] {
	return NewFromPoints(l.End(), l.Start())
}

// Length calculates the Euclidean distance (length) between the start and end points of the line segment.
//
// Parameters:
//   - opts: A variadic slice of [options.GeometryOptionsFunc] functions to customize the calculation behavior.
//     Common options include:
//   - [options.WithEpsilon](epsilon float64): Specifies a tolerance for snapping small floating-point
//     deviations to cleaner values, ensuring robustness in length calculations.
//
// Returns:
//   - float64: The length of the line segment, optionally adjusted based on epsilon.
//
// Behavior:
//   - The function computes the Euclidean distance between the start and end points of the line segment
//     using [point.Point.DistanceToPoint].
//   - If an epsilon threshold is provided via [options.WithEpsilon], the resulting length is adjusted to
//     correct small deviations caused by floating-point precision.
//
// Notes:
//   - This function is particularly useful for geometric computations where precise segment length is required.
//   - For integer coordinates, the epsilon adjustment has no effect and is ignored.
//
// Complexity:
//   - O(1): The calculation involves a single distance computation.
func (l LineSegment[T]) Length(opts ...options.GeometryOptionsFunc) float64 {
	return l.start.DistanceToPoint(l.end, opts...)
}

// normalize ensures that the line segment's coordinates match the expected ordering for the sweep line algorithm in FindIntersections.
// makes l.Start() the "upper" point, and seg.End() the "lower" point
// todo: doc comment, unit test, example
func (l LineSegment[T]) normalize() LineSegment[T] {

	// if start Y is smaller than end Y, then flip (top-to-bottom)
	if l.Start().Y() < l.End().Y() {
		return l.Flip()
	}

	// if Ys are equal, then order by X, with smallest X first (left-to-right
	if l.Start().Y() == l.End().Y() && l.Start().X() > l.End().X() {
		return l.Flip()
	}

	return l
}

// Points returns the start [point.Point] and end [point.Point] of the LineSegment.
//
// Returns:
//   - start ([Point][T]): The start [point.Point] of the LineSegment.
//   - end ([Point][T]): The end [point.Point] of the LineSegment.
func (l LineSegment[T]) Points() (start, end point.Point[T]) {
	return l.start, l.end
}

// ProjectPoint projects the [point.Point] p onto a given [LineSegment] l.
//
// The function calculates the closest point on the LineSegment to the Point p.
// It utilizes vector mathematics to determine the projection of Point p onto the infinite line
// represented by the [LineSegment]. If the projected point falls beyond the ends of the
// [LineSegment], the function returns the closest endpoint of the segment.
//
// Parameters:
//   - p ([point.Point][T]): The point to be projected onto line segment l
//
// Returns:
//   - A [point.Point][float64] representing the coordinates of the projected point.
//     If the [LineSegment] is degenerate (both endpoints are the same),
//     the function returns the coordinates of the Start() Point of the LineSegment.
func (l LineSegment[T]) ProjectPoint(p point.Point[T]) point.Point[float64] {

	// the direction vector of the line segment
	vecAB := l.end.Translate(l.start.Negate())

	// the vector from line segment start to point
	vecAP := p.Translate(l.start.Negate())

	// Calculate the dot products
	ABdotAB := vecAB.DotProduct(vecAB) // |vecAB|^2
	APdotAB := vecAP.DotProduct(vecAB) // vecAP • vecAB

	// Calculate the projection length as a fraction of the length of vecAB
	if ABdotAB == 0 { // Avoid division by zero; A and End are the same point
		return l.start.AsFloat64()
	}
	projLen := float64(APdotAB) / float64(ABdotAB)

	// Clamp the projection length to the segment
	if projLen < 0 {
		return l.start.AsFloat64() // Closest to line segment start
	} else if projLen > 1 {
		return l.end.AsFloat64() // Closest to line segment end
	}

	// return the projection point
	return point.New(
		float64(l.start.X())+(projLen*float64(vecAB.X())),
		float64(l.start.Y())+(projLen*float64(vecAB.Y())),
	)
}

// ReflectLineSegment reflects a given [LineSegment] `other` across the current line segment.
//
// This function calculates the reflection of each endpoint of the `other` line segment across
// the current line segment and returns a new reflected line segment.
//
// Parameters:
//   - other (LineSegment[T]): The line segment to be reflected.
//
// Returns:
//   - LineSegment[float64]: A new line segment whose endpoints are the reflections of
//     the endpoints of the `other` line segment across the current line segment.
//
// Behavior:
//   - The function uses the [LineSegment.ReflectPoint] method to compute the reflection of each endpoint
//     of `other` across the current line segment.
//   - The resulting line segment has endpoints represented as [float64], as reflections
//     may involve non-integer coordinates.
//
// Notes:
//   - This function assumes that the current line segment (`l`) is not degenerate (i.e., has non-zero length).
//   - If the `other` line segment coincides with the current line segment, the result is a line segment that
//     mirrors the original.
func (l LineSegment[T]) ReflectLineSegment(other LineSegment[T]) LineSegment[float64] {
	return NewFromPoints(l.ReflectPoint(other.Start()), l.ReflectPoint(other.End()))
}

// ReflectPoint reflects the [point.Point] across the axis defined by LineSegment l.
//
// Parameters:
//   - p ([point.Point][T]): The [point.Point] to be reflected about LineSegment l.
//
// Returns:
//   - Point[float64] - A new point representing the reflection of the original point.
func (l LineSegment[T]) ReflectPoint(p point.Point[T]) point.Point[float64] {
	pFloat := p.AsFloat64()
	lFloat := l.AsFloat64()

	// Extract points from the line segment
	x1, y1 := lFloat.start.X(), lFloat.start.Y()
	x2, y2 := lFloat.end.X(), lFloat.end.Y()

	// Calculate the line's slope and intercept for projection
	dx, dy := x2-x1, y2-y1
	if dx == 0 && dy == 0 {
		return pFloat // Degenerate line segment; return point unchanged
	}

	// Calculate the reflection using vector projection
	a := (dx*dx - dy*dy) / (dx*dx + dy*dy)
	b := 2 * dx * dy / (dx*dx + dy*dy)

	newX := a*(pFloat.X()-x1) + b*(pFloat.Y()-y1) + x1
	newY := b*(pFloat.X()-x1) - a*(pFloat.Y()-y1) + y1

	return point.New(newX, newY)
}

// RelationshipToPoint determines the spatial relationship of the current Point to a given [LineSegment].
//
// The function calculates the orthogonal (shortest) distance from the point to the line segment
// and determines the relationship based on this distance.
//
// Relationships:
//   - [types.RelationshipIntersection]: The point lies on the line segment.
//   - [types.RelationshipDisjoint]: The point does not lie on the line segment.
//
// Parameters:
//   - p ([point.Point][T]): The point to analyse the relationship with.
//   - opts: A variadic slice of [options.GeometryOptionsFunc] functions to customize the calculation.
//     [options.WithEpsilon](epsilon float64): Adjusts the precision for distance comparisons, enabling robust handling of floating-point errors.
//
// Returns:
//   - [types.Relationship]: The spatial relationship of the point to the line segment.
//
// Behavior:
//   - If the shortest distance between the point and the line segment is zero (or within the epsilon threshold),
//     the function returns [types.RelationshipIntersection].
//   - Otherwise, it returns [types.RelationshipDisjoint].
//
// Notes:
//   - This method is useful for determining if a point lies on a line segment, including endpoints and interior points.
//   - Epsilon adjustment is particularly useful for floating-point coordinates to avoid precision errors.
func (l LineSegment[T]) RelationshipToPoint(p point.Point[T], opts ...options.GeometryOptionsFunc) types.Relationship {
	distancePointToLineSegment := l.DistanceToPoint(p, opts...)
	if distancePointToLineSegment == 0 {
		return types.RelationshipIntersection
	}
	return types.RelationshipDisjoint
}

// Rotate rotates the LineSegment around a given pivot [point.Point] by a specified angle in radians counterclockwise.
// Optionally, an epsilon threshold can be applied to adjust the precision of the resulting coordinates.
//
// Parameters:
//   - pivot ([point.Point][T]): The point around which to rotate the line segment.
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
func (l LineSegment[T]) Rotate(pivot point.Point[T], radians float64, opts ...options.GeometryOptionsFunc) LineSegment[float64] {
	newStart := l.start.Rotate(pivot, radians, opts...)
	newEnd := l.end.Rotate(pivot, radians, opts...)
	return NewFromPoints(newStart, newEnd)
}

// Scale scales the line segment by a given factor from a specified reference point.
//
// Parameters:
//   - ref ([point.Point][T]): The reference point from which the scaling is applied. Using the origin
//     point (0, 0) scales the segment relative to the coordinate system's origin, while specifying
//     a custom reference point scales the segment relative to that point.
//   - factor ([T]): The scaling factor, where a value greater than 1 expands the segment,
//     and a value between 0 and 1 shrinks it.
//
// Behavior:
//   - The function scales both endpoints of the line segment relative to the specified
//     reference point using the [point.Point.Scale] method.
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
func (l LineSegment[T]) Scale(ref point.Point[T], factor T) LineSegment[T] {
	return NewFromPoints(
		l.start.Scale(ref, factor),
		l.end.Scale(ref, factor),
	)
}

// Slope calculates the slope of the line segment.
//
// The slope is calculated as the change in y-coordinates (dy) divided by
// the change in x-coordinates (dx) of the line segment. If the line segment
// is vertical (dx = 0), the slope is undefined, and the function returns math.NaN().
//
// Returns:
//   - float64: The calculated slope of the line segment. Returns math.NaN() if the slope is undefined.
//
// Notes:
//   - Vertical lines (dx = 0) are identified as having an undefined slope.
//   - Use math.IsNaN() to check if the slope is undefined.
func (l LineSegment[T]) Slope() float64 {
	dx := float64(l.end.X() - l.start.X())
	dy := float64(l.end.Y() - l.start.Y())

	if dx == 0 {
		return math.NaN() // Vertical line, slope undefined
	}
	return dy / dx
}

// Start returns the starting point of the line segment.
//
// This function provides access to the starting point of the LineSegment l, typically representing
// the beginning of the segment.
func (l LineSegment[T]) Start() point.Point[T] {
	return l.start
}

// String returns a formatted string representation of the line segment for debugging and logging purposes.
//
// The string representation includes the coordinates of the start and end points in the format:
// "(x1, y1)(x2, y2)", where (x1, y1) are the coordinates of the start point,
// and (x2, y2) are the coordinates of the end point.
//
// Returns:
//   - string: A string representing the line segment's start and end coordinates.
func (l LineSegment[T]) String() string {
	return fmt.Sprintf("(%v,%v)(%v,%v)", l.start.X(), l.start.Y(), l.end.X(), l.end.Y())
}

// Translate moves the LineSegment by a specified vector.
//
// This method shifts the LineSegment's position in the 2D plane by translating
// both its start and end points by the given vector delta. The relative
// orientation and length of the LineSegment remain unchanged.
//
// Parameters:
//   - delta ([point.Point][T]): The vector by which to translate the line segment.
//
// Returns:
//   - [LineSegment][T]: A new LineSegment translated by the specified vector.
//
// Notes:
//   - Translating the line segment effectively adds the delta vector to both
//     the start and end points of the segment.
//   - This operation is equivalent to a uniform shift, maintaining the segment's
//     shape and size while moving it to a new position.
func (l LineSegment[T]) Translate(delta point.Point[T]) LineSegment[T] {
	return NewFromPoints(
		l.start.Translate(delta),
		l.end.Translate(delta),
	)
}

// XAtY calculates the x-coordinate on the line segment at a given y-coordinate.
//
// Parameters:
//   - y (T): The y-coordinate at which to find the corresponding x-coordinate.
//
// Returns:
//   - float64: The x-coordinate at the given y-coordinate on the line segment, or math.NaN()
//     if the line is horizontal, vertical, or if the y-coordinate is outside the bounds of the segment.
//
// Behavior:
//   - If the line segment is vertical (undefined slope), the function returns math.NaN().
//   - If the provided y-coordinate is outside the range of the segment's y-coordinates, the function returns math.NaN().
//   - If the line segment is horizontal, the x-coordinate is constant and will match the start or end x-coordinate.
//
// Example:
//   - For a line segment from (1, 2) to (4, 6), calling XAtY(4) will return 2.5.
func (l LineSegment[T]) XAtY(y float64) float64 {
	A, B := l.Start().AsFloat64(), l.End().AsFloat64()

	// Ensure y is within bounds
	if (y < A.Y() && y < B.Y()) || (y > A.Y() && y > B.Y()) {
		return math.NaN()
	}

	// Handle vertical line case: x is constant for all y values in range
	if A.X() == B.X() {
		return A.X() // Valid as long as y is within bounds
	}

	// Compute x using interpolation
	return A.X() + (y-A.Y())*(B.X()-A.X())/(B.Y()-A.Y())
}

// YAtX calculates the y-coordinate on the line segment at a given x-coordinate.
//
// Parameters:
//   - x (T): The x-coordinate at which to find the corresponding y-coordinate.
//
// Returns:
//   - float64: The y-coordinate at the given x-coordinate on the line segment, or math.NaN()
//     if the line is horizontal, vertical, or if the x-coordinate is outside the bounds of the segment.
//
// Behavior:
//   - If the line segment is horizontal (zero slope), the function returns math.NaN().
//   - If the provided x-coordinate is outside the range of the segment's x-coordinates, the function returns math.NaN().
//   - If the line segment is vertical, the y-coordinate is constant and will match the start or end y-coordinate.
//
// Example:
//   - For a line segment from (1, 2) to (4, 6), calling YAtX(2.5) will return 3.5.
func (l LineSegment[T]) YAtX(x float64) float64 {
	A, B := l.Start().AsFloat64(), l.End().AsFloat64()

	// Ensure x is within bounds
	if (x < A.X() && x < B.X()) || (x > A.X() && x > B.X()) {
		return math.NaN()
	}

	// Handle horizontal line case: y is constant for all x values in range
	if A.Y() == B.Y() {
		return A.Y()
	}

	// Compute y using interpolation
	return A.Y() + (x-A.X())*(B.Y()-A.Y())/(B.X()-A.X())
}
