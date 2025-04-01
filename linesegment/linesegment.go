// Package linesegment provides fundamental geometric operations on line segments,
// including intersection detection, transformations, and computational geometry algorithms.
//
// # Overview
//
// This package defines the [LineSegment] type, which represents a finite straight segment
// between two points in a 2D plane. It supports various operations,
// such as computing intersections, transformations (scaling, rotation, translation),
// and checking geometric relationships.
//
// # Features
//
//   - Basic Operations: Methods for retrieving endpoints, length, midpoint, and orientation.
//   - Geometric Relationships: Functions to determine whether a point lies on the segment,
//     whether two segments intersect, and whether a segment is collinear with another.
//   - Transformations: Functions to translate, rotate, and scale line segments.
//   - Intersection detection via FindIntersectionsSlow:
//     A naive brute-force approach that compares all segment pairs.
//   - Intersection detection via FindIntersectionsFast:
//     A more efficient method using the sweep line algorithm from
//     [Computational Geometry: Algorithms and Applications], suitable for larger datasets.
//
// # Line Segment Intersection Algorithms
//
// There are two methods for determining intersections between a set of line segments:
//   - Naive Method (FindIntersectionsSlow)
//   - Sweep Line Algorithm (FindIntersectionsFast)
//
// The naive method FindIntersectionsSlow iterates over all pairs of line segments and directly checks whether they
// intersect using the [Intersection] method. This has O(n²) time complexity, making it
// inefficient for large datasets but useful as a reference for correctness. In fact,
// the testing/fuzzing of FindIntersectionsFast compares results to FindIntersectionsSlow as a reference.
//
// The sweep line method FindIntersectionsFast is implemented to more efficiently find all intersections
// among a set of line segments. This algorithm sweeps a vertical line from Y-max to Y-min across
// the plane, maintaining an active set of segments that intersect the sweep line.
// This method is outlined in Section 2.1 of [Computational Geometry: Algorithms and Applications].
//
// [Computational Geometry: Algorithms and Applications]: https://www.springer.com/gp/book/9783540779735
package linesegment

import (
	"encoding/json"
	"fmt"
	"github.com/mikenye/geom2d"
	"github.com/mikenye/geom2d/numeric"
	"github.com/mikenye/geom2d/point"
	"github.com/mikenye/geom2d/types"
	"math"
)

// LineSegment represents a line segment in a 2D space, defined by two endpoints.
// Internally, the points are stored as "upper" and "lower" points rather than "start" and "end",
// where the upper point has a higher Y-coordinate (or leftmost X if Y-coordinates are equal)
// and the lower point has a lower Y-coordinate (or rightmost X if Y-coordinates are equal).
// This ordering is enforced by the constructors to provide consistent behavior.
type LineSegment struct {
	upper point.Point
	lower point.Point
}

// New creates a new LineSegment with the specified start and end x and y coordinates.
//
// This constructor function initializes a [LineSegment] with the specified starting and ending points.
//
// Parameters:
//   - x1,y1 (float64): The starting point of the LineSegment.
//   - x2,y2 (float64): The ending point of the LineSegment.
//
// Returns:
//   - LineSegment: A new line segment defined by the two endpoints.
func New(x1, y1, x2, y2 float64) LineSegment {
	p1 := point.New(x1, y1)
	p2 := point.New(x2, y2)

	return NewFromPoints(p1, p2)
}

// NewFromPoints creates a new LineSegment from two endpoints.
//
// This constructor function initializes a LineSegment with the specified points, automatically
// determining which will be the upper and lower points according to the LineSegment ordering rules.
//
// Parameters:
//   - p1 (point.Point): The first endpoint of the LineSegment.
//   - p2 (point.Point): The second endpoint of the LineSegment.
//
// Returns:
//   - LineSegment: A new line segment defined by the two endpoints.
func NewFromPoints(p1, p2 point.Point) LineSegment {

	// Ensure p1 is the "upper" point (higher Y first, or rightmost X if tied)
	if p2.Y() > p1.Y() || (p2.Y() == p1.Y() && p2.X() < p1.X()) {
		p1, p2 = p2, p1 // Swap to maintain order
	}

	return LineSegment{
		upper: p1, // Always uppermost point first
		lower: p2,
	}
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
//   - Rendering lines in graphics applications.
//   - Generating grid points for pathfinding.
//
// Parameters:
//   - yield (func([point.Point]) bool): A function that processes each generated point.
//     Returning false will stop further point generation.
//
// Note: This method requires integer-type coordinates for the line segment.
func (l LineSegment) Bresenham(yield func(point.Point) bool) {

	var x1, x2, y1, y2, dx, dy, sx, sy float64

	x1 = l.upper.X()
	x2 = l.lower.X()
	y1 = l.upper.Y()
	y2 = l.lower.Y()

	// Calculate absolute deltas
	dx = math.Abs(x2 - x1)
	dy = math.Abs(y2 - y1)

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

// Center calculates the midpoint of the line segment.
//
// Behavior:
//   - The midpoint is calculated by averaging the x and y coordinates of the upper and lower
//     points of the line segment.
//
// Returns:
//   - point.Point: The midpoint of the line segment.
func (l LineSegment) Center() point.Point {

	midX := (l.upper.X() + l.lower.X()) / 2
	midY := (l.upper.Y() + l.lower.Y()) / 2

	return point.New(midX, midY)
}

// ContainsPoint determines whether the given point lies on the LineSegment.
//
// This method calculates the shortest distance from the given point to the LineSegment
// using the DistanceToPoint method. If the distance is zero (considering epsilon threshold),
// the point is determined to be on the segment.
//
// Parameters:
//   - p (point.Point): The point to test against the LineSegment.
//
// Returns:
//   - bool: true if the point lies on the LineSegment, false otherwise.
//
// Notes:
//   - This function uses the DistanceToPoint method to compute the distance.
//   - Floating-point precision issues are handled using the global epsilon value.
//   - The point must also be within the bounding box defined by the segment endpoints to return true.
func (l LineSegment) ContainsPoint(p point.Point) bool {

	epsilon := geom2d.GetEpsilon()

	// Compute vectors AP and AB
	ap := p.Sub(l.upper)
	ab := l.lower.Sub(l.upper)

	// Dynamically adjust epsilon based on the segment length
	segmentLength := ab.DistanceToPoint(point.Origin())
	adaptiveEpsilon := epsilon * segmentLength

	// Check if cross product is within adaptive epsilon (collinearity test)
	cross := math.Abs(ap.CrossProduct(ab))
	if cross > adaptiveEpsilon {
		return false // P is not on the line
	}

	// Check if P is within the bounding box of the segment
	xMin, xMax := math.Min(l.upper.X(), l.lower.X()), math.Max(l.upper.X(), l.lower.X())
	yMin, yMax := math.Min(l.upper.Y(), l.lower.Y()), math.Max(l.upper.Y(), l.lower.Y())

	return (p.X() >= xMin-adaptiveEpsilon && p.X() <= xMax+adaptiveEpsilon) &&
		(p.Y() >= yMin-adaptiveEpsilon && p.Y() <= yMax+adaptiveEpsilon)
}

// DistanceToLineSegment calculates the minimum distance between two line segments, l and other.
//
// If the segments intersect or touch at any point, the function immediately returns 0, as the distance is effectively zero.
// Otherwise, it calculates the shortest distance by considering:
//  1. The distances between endpoints of one segment and the other segment.
//  2. The distances from endpoints of one segment to the perpendicular projections onto the other segment.
//
// Parameters:
//   - other (LineSegment): The line segment to compare with l.
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
//
// DistanceToLineSegment returns the shortest distance between two line segments.
func (l LineSegment) DistanceToLineSegment(other LineSegment) float64 {
	// If segments intersect, the distance is zero.
	if l.Intersects(other) {
		return 0
	}

	// Compute distances from endpoints to the other segment
	d1 := l.DistanceToPoint(other.Upper()) // Distance from A to C-D
	d2 := l.DistanceToPoint(other.Lower()) // Distance from B to C-D
	d3 := other.DistanceToPoint(l.Upper()) // Distance from C to A-B
	d4 := other.DistanceToPoint(l.Lower()) // Distance from D to A-B

	// Return the minimum distance
	return math.Min(math.Min(d1, d2), math.Min(d3, d4))
}

// DistanceToPoint calculates the orthogonal (shortest) distance from the LineSegment l to the point p.
// This distance is the length of the perpendicular line from p to the closest point on l.
//
// Parameters:
//   - p (point.Point): The point to which the distance is calculated from LineSegment l.
//
// Behavior:
//   - The function first computes the projection of p onto the given LineSegment l. This is
//     the closest point on l to p, whether it falls within the line segment or on one of its endpoints.
//   - The distance is then calculated as the Euclidean distance from p to the projected point,
//     using the point.Point.DistanceToPoint method for precision.
//
// Returns:
//   - float64: The shortest distance between the point p and the line segment l.
//
// Notes:
//   - If the point p lies exactly on the line segment, the distance will be zero.
//   - Note that ContainsPoint (not this method) uses epsilon-based comparison to determine
//     if a point is considered to be on the line segment.
func (l LineSegment) DistanceToPoint(p point.Point) float64 {
	projectedPoint := l.ProjectPoint(p)
	return p.DistanceToPoint(projectedPoint)
}

// Eq checks if two line segments are equal by comparing their upper and lower points.
// Equality is evaluated using an epsilon-based comparison with the global epsilon value.
//
// Parameters:
//   - other (LineSegment): The line segment to compare with the current line segment.
//
// Behavior:
//   - The function performs an approximate equality check using the global epsilon value,
//     considering the points equal if their coordinate differences are within this threshold.
//
// Returns:
//   - bool: Returns true if both line segments have approximately equal upper and lower points;
//     otherwise, false.
//
// Notes:
//   - Approximate equality is useful when comparing line segments with floating-point coordinates,
//     where small precision errors might otherwise cause inequality.
//   - This function relies on the point.Point.Eq method, which uses the global epsilon value.
//
// Note: In this implementation, equality is defined by matching upper and lower points using epsilon-based
// comparison. NewFromPoints enforces a consistent internal ordering of points, so segments with the
// same endpoints will compare as equal regardless of the order the points were provided.
func (l LineSegment) Eq(other LineSegment) bool {
	return l.upper.Eq(other.upper) && l.lower.Eq(other.lower)
}

// IntersectionPoints calculates the points of intersection between two line segments.
//
// This function determines if and where two line segments intersect. It handles different cases:
// 1. Regular intersection at a single point
// 2. Collinear overlapping segments (returns both endpoints of the overlap)
// 3. No intersection (returns empty slice)
//
// Parameters:
//   - other (LineSegment): The line segment to check for intersection with this segment.
//
// Returns:
//   - []point.Point: A slice containing intersection points. For regular intersections, this contains
//     a single point. For collinear overlapping segments, it contains the two endpoints of the overlap.
//   - bool: True if there is at least one intersection point, false otherwise.
//
// Behavior:
//   - For non-parallel segments, calculates the intersection point using line equations and checks
//     if this point lies on both segments.
//   - For collinear segments, determines if they overlap and returns both endpoints of the
//     overlapping section if they do.
//   - Returns an empty slice and false if segments don't intersect.
func (l LineSegment) IntersectionPoints(other LineSegment) ([]point.Point, bool) {

	// Line AB represented as a1x + b1y = c1
	a1 := l.lower.Y() - l.upper.Y()
	b1 := l.upper.X() - l.lower.X()
	c1 := a1*(l.upper.X()) + b1*(l.upper.Y())

	// Line CD represented as a2x + b2y = c2
	a2 := other.lower.Y() - other.upper.Y()
	b2 := other.upper.X() - other.lower.X()
	c2 := a2*(other.upper.X()) + b2*(other.upper.Y())

	determinant := a1*b2 - a2*b1

	if numeric.FloatEquals(determinant, 0, geom2d.GetEpsilon()) {
		// Lines are either parallel or collinear
		if l.ContainsPoint(other.lower) || l.ContainsPoint(other.upper) ||
			other.ContainsPoint(l.lower) || other.ContainsPoint(l.upper) {

			// Determine the overlapping segment
			overlapStart := point.New(math.Max(l.lower.X(), other.lower.X()), math.Max(l.lower.Y(), other.lower.Y()))
			overlapEnd := point.New(math.Min(l.upper.X(), other.upper.X()), math.Min(l.upper.Y(), other.upper.Y()))

			if overlapStart.X() > overlapEnd.X() || overlapStart.Y() > overlapEnd.Y() {
				return []point.Point{}, false // No overlap
			}

			return []point.Point{overlapStart, overlapEnd}, true
		}
		// Parallel but not collinear
		return []point.Point{}, false
	} else {
		// Compute intersection point
		x := (b2*c1 - b1*c2) / determinant
		y := (a1*c2 - a2*c1) / determinant
		intersection := point.New(x, y)

		// Check if the intersection is within both line segments
		if l.ContainsPoint(intersection) && other.ContainsPoint(intersection) {
			return []point.Point{intersection}, true
		}
		return []point.Point{}, false
	}
}

// Intersects determines if two line segments intersect.
//
// This function checks if this line segment intersects with another line segment.
// It uses orientation tests to efficiently determine intersection without computing
// the actual intersection points.
//
// Parameters:
//   - other (LineSegment): The line segment to check for intersection with this segment.
//
// Returns:
//   - bool: True if the segments intersect, false otherwise.
//
// Behavior:
//   - For segments to intersect, they must cross each other or share a point.
//   - The function first checks if the segments straddle each other (orientation test).
//   - It also handles special cases like collinear segments that may overlap.
//   - Two segments with coincident endpoints are considered to intersect.
func (l LineSegment) Intersects(other LineSegment) bool {
	a, b := l.upper, l.lower
	c, d := other.upper, other.lower

	// Compute orientation values
	o1 := point.Orientation(a, b, c)
	o2 := point.Orientation(a, b, d)
	o3 := point.Orientation(c, d, a)
	o4 := point.Orientation(c, d, b)

	// General case: If the two segments straddle each other
	if o1 != o2 && o3 != o4 {
		return true
	}

	// Special case: Check for collinear overlap
	if o1 == 0 && NewFromPoints(a, c).ContainsPoint(b) {
		return true
	}
	if o2 == 0 && NewFromPoints(a, d).ContainsPoint(b) {
		return true
	}
	if o3 == 0 && NewFromPoints(c, a).ContainsPoint(d) {
		return true
	}
	if o4 == 0 && NewFromPoints(c, b).ContainsPoint(d) {
		return true
	}

	return false
}

// Length calculates the Euclidean distance (length) between the start and end points of the line segment.
//
// Returns:
//   - float64: The length of the line segment.
//
// Behavior:
//   - The function computes the Euclidean distance between the start and end points of the line segment
//     using [point.Point.DistanceToPoint].
func (l LineSegment) Length() float64 {
	return l.upper.DistanceToPoint(l.lower)
}

// Lower returns the lower [point.Point] of the LineSegment.
func (l LineSegment) Lower() point.Point {
	return l.lower
}

// MarshalJSON serializes LineSegment as JSON while preserving its original type.
func (l LineSegment) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Upper point.Point `json:"upper"`
		Lower point.Point `json:"lower"`
	}{
		Upper: l.Upper(),
		Lower: l.Lower(),
	})
}

// Points returns the upper and lower points of the LineSegment.
//
// Returns:
//   - upper (point.Point): The upper point of the LineSegment.
//   - lower (point.Point): The lower point of the LineSegment.
func (l LineSegment) Points() (upper, lower point.Point) {
	return l.upper, l.lower
}

// ProjectPoint projects the point p onto a given LineSegment l.
//
// The function calculates the closest point on the LineSegment to the point p.
// It utilizes vector mathematics to determine the projection of point p onto the infinite line
// represented by the LineSegment. If the projected point falls beyond the ends of the
// LineSegment, the function returns the closest endpoint of the segment.
//
// Parameters:
//   - p (point.Point): The point to be projected onto line segment l
//
// Returns:
//   - point.Point: The coordinates of the projected point.
//     If the LineSegment is degenerate (both endpoints are the same),
//     the function returns the coordinates of the lower point of the LineSegment
//     (the point with lower Y-coordinate, or if Y-coordinates are equal, the point with higher X-coordinate).
func (l LineSegment) ProjectPoint(p point.Point) point.Point {
	// Compute the direction vector of the line segment
	vecAB := l.upper.Sub(l.lower) // Ensure this is (upper - lower)

	// Compute the vector from segment start to the point
	vecAP := p.Sub(l.lower)

	// Compute the dot products
	ABdotAB := vecAB.DotProduct(vecAB) // |AB|^2
	APdotAB := vecAP.DotProduct(vecAB) // AP • AB

	// If segment has zero length, return the lower point
	if ABdotAB == 0 {
		return l.lower
	}

	// Compute projection scalar t and clamp it to [0,1]
	t := math.Max(0, math.Min(1, APdotAB/ABdotAB))

	// Compute the projected point as lower + t * vecAB
	return l.lower.Add(vecAB.Scale(point.New(0, 0), t))
}

// ReflectLineSegment reflects a given [LineSegment] `other` across the current line segment.
//
// This function calculates the reflection of each endpoint of the `other` line segment across
// the current line segment and returns a new reflected line segment.
//
// Parameters:
//   - other (LineSegment): The line segment to be reflected.
//
// Returns:
//   - LineSegment: A new line segment whose endpoints are the reflections of
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
func (l LineSegment) ReflectLineSegment(other LineSegment) LineSegment {
	return NewFromPoints(l.ReflectPoint(other.upper), l.ReflectPoint(other.lower))
}

// ReflectPoint reflects the [point.Point] across the axis defined by LineSegment l.
//
// Parameters:
//   - p (point.Point): The point to be reflected about LineSegment l.
//
// Returns:
//   - point.Point: A new point representing the reflection of the original point.
func (l LineSegment) ReflectPoint(p point.Point) point.Point {

	// Extract points from the line segment
	x1, y1 := l.upper.X(), l.upper.Y()
	x2, y2 := l.lower.X(), l.lower.Y()

	// Calculate the line's slope and intercept for projection
	dx, dy := x2-x1, y2-y1
	if dx == 0 && dy == 0 {
		return p // Degenerate line segment; return point unchanged
	}

	// Calculate the reflection using vector projection
	a := (dx*dx - dy*dy) / (dx*dx + dy*dy)
	b := 2 * dx * dy / (dx*dx + dy*dy)

	newX := a*(p.X()-x1) + b*(p.Y()-y1) + x1
	newY := b*(p.X()-x1) - a*(p.Y()-y1) + y1

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
//   - p (point.Point): The point to analyze the relationship with.
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
func (l LineSegment) RelationshipToPoint(p point.Point) types.Relationship {
	distancePointToLineSegment := l.DistanceToPoint(p)
	if distancePointToLineSegment == 0 {
		return types.RelationshipIntersection
	}
	return types.RelationshipDisjoint
}

// Rotate rotates the LineSegment around a given pivot [point.Point] by a specified angle in radians counterclockwise.
// Optionally, an epsilon threshold can be applied to adjust the precision of the resulting coordinates.
//
// Parameters:
//   - pivot (point.Point): The point around which to rotate the line segment.
//   - radians (float64): The rotation angle in radians.
//
// Behavior:
//   - The function rotates the upper and lower points of the line segment around the given pivot
//     point by the specified angle using the [point.Point.Rotate] method.
//
// Returns:
//   - LineSegment: A new line segment representing the rotated position.
func (l LineSegment) Rotate(pivot point.Point, radians float64) LineSegment {
	newStart := l.upper.Rotate(pivot, radians)
	newEnd := l.lower.Rotate(pivot, radians)
	return NewFromPoints(newStart, newEnd)
}

// Scale scales the line segment by a given factor from a specified reference point.
//
// Parameters:
//   - ref (point.Point): The reference point from which the scaling is applied. Using the origin
//     point (0, 0) scales the segment relative to the coordinate system's origin, while specifying
//     a custom reference point scales the segment relative to that point.
//   - factor (float64): The scaling factor, where a value greater than 1 expands the segment,
//     and a value between 0 and 1 shrinks it.
//
// Behavior:
//   - The function scales both endpoints of the line segment relative to the specified
//     reference point using the [point.Point.Scale] method.
//   - The scaling operation preserves the relative orientation of the segment.
//
// Returns:
//   - LineSegment: A new line segment, scaled relative to the specified reference point.
//
// Notes:
//   - Scaling by a factor of 1 will return a line segment identical to the original.
//   - Negative scaling factors will mirror the segment across the reference point
//     and scale its length accordingly.
func (l LineSegment) Scale(ref point.Point, factor float64) LineSegment {
	return NewFromPoints(
		l.upper.Scale(ref, factor),
		l.lower.Scale(ref, factor),
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
//   - Use math.IsNaN() to check if the slope is undefined (and thus vertical).
func (l LineSegment) Slope() float64 {
	dx := float64(l.lower.X() - l.upper.X())
	dy := float64(l.lower.Y() - l.upper.Y())

	if dx == 0 {
		return math.NaN() // Vertical line, slope undefined
	}
	return dy / dx
}

// String returns a formatted string representation of the line segment for debugging and logging purposes.
//
// The string representation includes the coordinates of the start and end points in the format:
// "(x1, y1)(x2, y2)", where (x1, y1) are the coordinates of the start point,
// and (x2, y2) are the coordinates of the end point.
//
// Returns:
//   - string: A string representing the line segment's start and end coordinates.
func (l LineSegment) String() string {
	return fmt.Sprintf("(%v,%v)(%v,%v)", l.upper.X(), l.upper.Y(), l.lower.X(), l.lower.Y())
}

// Translate moves the LineSegment by a specified vector.
//
// This method shifts the LineSegment's position in the 2D plane by translating
// both its start and end points by the given vector delta. The relative
// orientation and length of the LineSegment remain unchanged.
//
// Parameters:
//   - delta (point.Point): The vector by which to translate the line segment.
//
// Returns:
//   - LineSegment: A new LineSegment translated by the specified vector.
//
// Notes:
//   - Translating the line segment effectively adds the delta vector to both
//     the start and end points of the segment.
//   - This operation is equivalent to a uniform shift, maintaining the segment's
//     shape and size while moving it to a new position.
func (l LineSegment) Translate(delta point.Point) LineSegment {
	return NewFromPoints(
		l.upper.Translate(delta),
		l.lower.Translate(delta),
	)
}

// UnmarshalJSON deserializes JSON into a LineSegment while keeping the exact original type.
func (l *LineSegment) UnmarshalJSON(data []byte) error {
	var temp struct {
		Upper point.Point `json:"upper"`
		Lower point.Point `json:"lower"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	l.upper = temp.Upper
	l.lower = temp.Lower
	return nil
}

// Upper returns the upper point of the line segment.
func (l LineSegment) Upper() point.Point {
	return l.upper
}

// XAtY calculates the x-coordinate on the line segment at a given y-coordinate.
//
// Parameters:
//   - y (float64): The y-coordinate at which to find the corresponding x-coordinate.
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
// Notes:
//   - Users should check for `NaN` with `math.IsNaN()` before using the result.
//
// Example:
//   - For a line segment from (1, 2) to (4, 6), calling XAtY(4) will return 2.5.
func (l LineSegment) XAtY(y float64) float64 {
	A, B := l.upper, l.lower

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
//   - x (float64): The x-coordinate at which to find the corresponding y-coordinate.
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
// Notes:
//   - Users should check for `NaN` with `math.IsNaN()` before using the result.
//
// Example:
//   - For a line segment from (1, 2) to (4, 6), calling YAtX(2.5) will return 3.5.
func (l LineSegment) YAtX(x float64) float64 {
	A, B := l.upper, l.lower

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

// Commented out for future implementation:
//// mergeSegments combines two slices of line segments, removing duplicates.
//func mergeSegments(a, b []LineSegment) []LineSegment {
//	input := append(a, b...)
//	output := make([]LineSegment, 0, len(a)+len(b))
//	for _, seg := range input {
//		if !slices.ContainsFunc(output, func(l LineSegment) bool {
//			return l.Eq(seg)
//		}) {
//			output = append(output, seg)
//		}
//	}
//	return output
//}

// Left returns the endpoint with the smaller X coordinate.
//
// This method determines which of the line segment's endpoints is leftmost
// (has the smaller X coordinate). If both endpoints have the same X-coordinate,
// it returns the endpoint with the smaller Y-coordinate (the one lower on the coordinate plane).
//
// Returns:
//   - point.Point: The leftmost endpoint of the line segment.
//
// Behavior:
//   - Compares the X coordinates of both endpoints first.
//   - If X coordinates are equal, uses Y coordinates as a tiebreaker.
//   - For vertical segments, returns the endpoint with smaller Y.
func (l LineSegment) Left() point.Point {
	if l.upper.X() < l.lower.X() || (l.upper.X() == l.lower.X() && l.upper.Y() < l.lower.Y()) {
		return l.upper
	}
	return l.lower
}

// Right returns the endpoint with the greater X coordinate.
//
// This method determines which of the line segment's endpoints is rightmost
// (has the greater X coordinate). If both endpoints have the same X-coordinate,
// it returns the endpoint with the greater Y-coordinate (the one higher on the coordinate plane).
//
// Returns:
//   - point.Point: The rightmost endpoint of the line segment.
//
// Behavior:
//   - Compares the X coordinates of both endpoints first.
//   - If X coordinates are equal, uses Y coordinates as a tiebreaker.
//   - For vertical segments, returns the endpoint with greater Y.
func (l LineSegment) Right() point.Point {
	if l.upper.X() > l.lower.X() || (l.upper.X() == l.lower.X() && l.upper.Y() > l.lower.Y()) {
		return l.upper
	}
	return l.lower
}
