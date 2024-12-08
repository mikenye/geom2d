package geom2d

import (
	"fmt"
	"image"
	"math"
	"slices"
)

// Point represents a point in two-dimensional space with x and y coordinates of a generic numeric type T.
// The Point struct provides methods for common vector operations such as addition, subtraction, and distance
// calculations, making it versatile for computational geometry and graphics applications.
//
// Type Parameter:
//   - T: The numeric type for the coordinates, constrained to signed number types by the [SignedNumber] interface.
type Point[T SignedNumber] struct {
	x T
	y T
}

// ConvexHull computes the [convex hull] of a finite set of points using the [Graham Scan] algorithm.
// The convex hull is the smallest convex polygon that encloses all points in the input set. This function is
// particularly useful in computational geometry applications where the outer boundary of a set of points is needed.
//
// This implementation finds the point with the lowest y-coordinate to serve as a reference for sorting points by their
// angle relative to this point. Starting from this reference, it orders the points counterclockwise, removing any
// points that cause a clockwise turn to ensure a convex boundary.
//
// Parameters:
//   - points ([Point][T]): A variable number of instances representing the set of points for which the
//     convex hull is to be computed.
//
// Returns:
//   - [][Point][T]: A slice of points representing the vertices of the convex hull in counterclockwise order.
//     The returned slice includes only the points that form the outer boundary of the input set.
//
// Note:
//   - If the points parameter is empty or has fewer than three points, the function returns the input points
//     unchanged, as a convex hull cannot be formed.
//
// [Graham Scan]: https://en.wikipedia.org/wiki/Graham_scan
// [convex hull]: https://en.wikipedia.org/wiki/Convex_hull
func ConvexHull[T SignedNumber](points []Point[T]) []Point[T] {

	var (
		pt0Index, pt1Index, pt2Index int
		pt0, pt1, pt2                Point[T]
	)

	// Copy points into a new slice, to prevent modifying input slice
	output := make([]Point[T], len(points))
	_ = copy(output, points)

	// Find the point with the lowest y-coordinate.
	// If the lowest y-coordinate exists in more than one point in the set,
	// the point with the lowest x-coordinate out of the candidates should be chosen.
	_, lowestPoint := findLowestLeftestPoint(output...)

	// Order the points by angle about the lowest point
	orderPointsByAngleAboutLowestPoint(lowestPoint, output)

	// Starting with the lowest point, work through points, popping off
	// any that cause a clockwise turn, to maintain convexity.
	for pt0Index = 0; pt0Index < len(output); pt0Index++ {
		pt1Index = (pt0Index + 1) % len(output)
		pt2Index = (pt1Index + 1) % len(output)
		pt0 = output[pt0Index]
		pt1 = output[pt1Index]
		pt2 = output[pt2Index]
		o := Orientation(pt0, pt1, pt2)
		if o == PointsClockwise {
			output = slices.Delete(output, pt1Index, pt1Index+1)
			pt0Index -= 3
			if pt0Index < 0 {
				pt0Index = 0
			}
		}
	}

	return output
}

// NewPoint creates a new Point with the specified x and y coordinates.
//
// This function is generic and requires the x and y values to satisfy the [SignedNumber] constraint.
//
// Parameters:
//   - x (T): The x-coordinate of the point.
//   - y (T): The y-coordinate of the point.
//
// Returns:
//   - Point[T]: A new Point instance with the given coordinates.
func NewPoint[T SignedNumber](x, y T) Point[T] {
	return Point[T]{
		x: x,
		y: y,
	}
}

// NewPointFromImagePoint creates and returns a new Point with integer x and y coordinates
// based on an [image.Point]. This function is useful for converting between graphics and
// computational geometry representations of points.
//
// Parameters:
//   - q ([image.Point]): An [image.Point] representing the source coordinates for the new point.
//
// Returns:
//   - Point[int]: A new Point with coordinates corresponding to the x and y values of the provided [image.Point].
func NewPointFromImagePoint(q image.Point) Point[int] {
	return Point[int]{
		x: q.X,
		y: q.Y,
	}
}

// AsFloat converts the Point's x and y coordinates to the float64 type, returning a new Point[float64].
// This method is useful when higher precision or floating-point arithmetic is needed on the coordinates.
//
// Returns:
//   - Point[float64]: A new Point with x and y coordinates converted to float64.
func (p Point[T]) AsFloat() Point[float64] {
	return Point[float64]{
		x: float64(p.x),
		y: float64(p.y),
	}
}

// AsInt converts the Point's x and y coordinates to the int type by truncating any decimal values.
// This method is useful when integer coordinates are needed for operations that require whole numbers,
// such as pixel-based calculations.
//
// Returns:
//   - Point[int]: A new Point with x and y coordinates converted to int by truncating any decimal portion.
func (p Point[T]) AsInt() Point[int] {
	return Point[int]{
		x: int(p.x),
		y: int(p.y),
	}
}

// AsIntRounded converts the Point's x and y coordinates to the int type by rounding to the nearest integer.
// This method is useful when integer coordinates are required and rounding provides more accurate results
// compared to truncation.
//
// Returns:
//   - Point[int]: A new Point with x and y coordinates converted to int by rounding to the nearest integer.
func (p Point[T]) AsIntRounded() Point[int] {
	return Point[int]{
		x: int(math.Round(float64(p.x))),
		y: int(math.Round(float64(p.y))),
	}
}

// CrossProduct calculates the cross product of the vector from the origin to Point p and the vector from the origin
// to Point q. The result is useful in determining the relative orientation of two vectors:
//   - A positive result indicates a counterclockwise turn (left turn),
//   - A negative result indicates a clockwise turn (right turn),
//   - A result of zero indicates that the points are collinear.
//
// Parameters:
//   - q (Point[T]): The Point with which to compute the cross product relative to the calling Point.
//
// Returns:
//   - T: The cross product of the vectors from the origin to p and q, indicating their relative orientation.
func (p Point[T]) CrossProduct(q Point[T]) T {
	return (p.x * q.y) - (p.y * q.x)
}

// DistanceSquaredToPoint calculates the squared Euclidean distance between Point p and another Point q.
// This method returns the squared distance, which avoids the computational cost of a square root calculation
// and is useful in cases where only distance comparisons are needed.
//
// Parameters:
//   - q (Point[T]): The Point to which the squared distance is calculated from the calling Point.
//
// Returns:
//   - T: The squared Euclidean distance between p and q.
func (p Point[T]) DistanceSquaredToPoint(q Point[T]) T {
	return (q.x-p.x)*(q.x-p.x) + (q.y-p.y)*(q.y-p.y)
}

// DistanceToLineSegment calculates the orthogonal (shortest) distance from the current Point p
// to a specified [LineSegment] l. This distance is the length of the perpendicular line segment
// from p to the closest point on l.
//
// Parameters:
//   - l ([LineSegment][T]): The [LineSegment] to which the distance is calculated.
//   - opts: A variadic slice of [Option] functions to customize the calculation behavior.
//     [WithEpsilon](epsilon float64): Adjusts the result by snapping small floating-point
//     deviations to cleaner values based on the specified epsilon threshold.
//
// Behavior:
//   - The function first computes the projection of p onto the given [LineSegment] l. This is
//     the closest point on l to p, whether it falls within the line segment or on one of its endpoints.
//   - The distance is then calculated as the Euclidean distance from p to the projected point,
//     using the [Point.DistanceToPoint] method for precision.
//
// Returns:
//   - float64: The shortest distance between the point p and the line segment l, optionally
//     adjusted based on epsilon if provided.
//
// Notes:
//   - If the point p lies exactly on the line segment, the distance will be zero (or adjusted
//     to zero if within epsilon).
//   - This method ensures precision by converting points to float64 before performing calculations.
func (p Point[T]) DistanceToLineSegment(l LineSegment[T], opts ...Option) float64 {
	projectedPoint := p.ProjectOntoLineSegment(l)
	pf := p.AsFloat()
	return pf.DistanceToPoint(projectedPoint, opts...)
}

// DistanceToPoint calculates the Euclidean (straight-line) distance between the current Point p
// and another Point q. The result is returned as a float64, providing precise measurement
// of the straight-line separation between the two points.
//
// Parameters:
//   - q ([Point][T]): The Point to which the distance is calculated.
//   - opts: A variadic slice of [Option] functions to customize the behavior of the calculation.
//     [WithEpsilon](epsilon float64): Adjusts the result by snapping near-zero distances
//     (or distances close to clean values like integers) to a more precise value based
//     on the specified epsilon threshold.
//
// Behavior:
//   - The function computes the Euclidean distance using the formula:
//     distance = sqrt((p.x - q.x)^2 + (p.y - q.y)^2)
//   - If the [WithEpsilon] option is provided, the computed distance is adjusted such that
//     deviations within the epsilon range are snapped to a clean value.
//
// Returns:
//   - float64: The Euclidean distance between the two points, optionally adjusted based on epsilon.
//
// Notes:
//   - Epsilon adjustment is useful for snapping results to clean values, particularly when
//     small floating-point errors could propagate in further calculations.
func (p Point[T]) DistanceToPoint(q Point[T], opts ...Option) float64 {

	// Apply options with defaults
	options := applyOptions(geomOptions{epsilon: 0}, opts...)

	// Calculate distance
	distance := math.Sqrt(float64(p.DistanceSquaredToPoint(q)))

	// Apply epsilon if specified
	if options.epsilon > 0 {
		distance = applyEpsilon(distance, options.epsilon)
	}

	return distance
}

// DotProduct calculates the dot product of the vector represented by Point p with the vector represented by Point q.
// The dot product is defined as p.x*q.x + p.y*q.y and is widely used in geometry for angle calculations,
// projection operations, and determining the relationship between two vectors.
//
// Parameters:
//   - q (Point[T]): The Point with which to calculate the dot product relative to the calling Point.
//
// Returns:
//   - T: The dot product of the vectors represented by p and q.
func (p Point[T]) DotProduct(q Point[T]) T {
	return (p.x * q.x) + (p.y * q.y)
}

// Eq determines whether the calling Point p is equal to another Point q.
// Equality can be evaluated either exactly (default) or approximately using an epsilon threshold.
//
// Parameters:
//   - q (Point[T]): The Point to compare with the calling Point.
//   - opts: A variadic slice of [Option] functions to customize the equality check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing the coordinates
//     of p and q. If the absolute difference between the coordinates of p and q is less
//     than epsilon, they are considered equal.
//
// Behavior:
//   - By default, the function performs an exact equality check, returning true only if
//     the x and y coordinates of p and q are identical.
//   - If the [WithEpsilon] option is provided, the function performs an approximate equality
//     check, considering p and q equal if their coordinate differences are within the specified
//     epsilon threshold.
//
// Returns:
//   - bool: True if p and q are equal based on the specified comparison mode; otherwise, false.
//
// Notes:
//   - Approximate equality is particularly useful when comparing points with floating-point
//     coordinates, where small precision errors may result in slightly different values.
func (p Point[T]) Eq(q Point[T], opts ...Option) bool {
	// Apply options with defaults
	options := applyOptions(geomOptions{epsilon: 0}, opts...)

	if options.epsilon > 0 {
		return math.Abs(float64(p.x)-float64(q.x)) < options.epsilon &&
			math.Abs(float64(p.y)-float64(q.y)) < options.epsilon
	}

	// Exact equality for default behavior
	return p.x == q.x && p.y == q.y
}

// Negate returns a new Point with both x and y coordinates negated.
// This operation is equivalent to reflecting the point across the origin
// and is useful in reversing the direction of a vector or preparing
// a point for subtraction via translation.
//
// Returns:
//   - Point[T]: A new Point with negated x and y coordinates.
//
// Notes:
//   - The returned point has the same type as the calling point.
//   - This method does not modify the original point but returns a new instance.
func (p Point[T]) Negate() Point[T] {
	return NewPoint[T](-p.x, -p.y)
}

// ProjectOntoLineSegment projects the Point p onto a given [LineSegment] l.
//
// The function calculates the closest point on the LineSegment defined by
// the endpoints l.Start() and l.End() to the Point p. It utilizes vector
// mathematics to determine the projection of Point p onto the infinite line
// represented by the [LineSegment]. If the projected point falls beyond the ends of the
// [LineSegment], the function returns the closest endpoint of the segment.
//
// Parameters:
//   - l ([LineSegment][T]): The line segment onto which the point is projected. It consists of
//     two endpoints: Start() and End().
//
// Returns:
//   - A Point[float64] representing the coordinates of the projected point.
//     If the [LineSegment] is degenerate (both endpoints are the same),
//     the function returns the coordinates of the Start() Point of the LineSegment.
func (p Point[T]) ProjectOntoLineSegment(l LineSegment[T]) Point[float64] {

	// the direction vector of the line segment
	vecAB := l.end.Translate(l.start.Negate())

	// the vector from line segment start to point
	vecAP := p.Translate(l.start.Negate())

	// Calculate the dot products
	ABdotAB := vecAB.DotProduct(vecAB) // |vecAB|^2
	APdotAB := vecAP.DotProduct(vecAB) // vecAP • vecAB

	// Calculate the projection length as a fraction of the length of vecAB
	if ABdotAB == 0 { // Avoid division by zero; A and End are the same point
		return l.start.AsFloat()
	}
	projLen := float64(APdotAB) / float64(ABdotAB)

	// Clamp the projection length to the segment
	if projLen < 0 {
		return l.start.AsFloat() // Closest to line segment start
	} else if projLen > 1 {
		return l.end.AsFloat() // Closest to line segment end
	}

	// return the projection point
	return NewPoint(
		float64(l.start.x)+(projLen*float64(vecAB.x)),
		float64(l.start.y)+(projLen*float64(vecAB.y)),
	)
}

// Reflect reflects the point across the specified axis or custom line.
//
// Parameters:
//   - axis ([ReflectionAxis]): Axis - The axis of reflection ([ReflectAcrossXAxis], [ReflectAcrossYAxis], or [ReflectAcrossCustomLine]).
//   - line ([LineSegment][T]): Optional. The line for [ReflectAcrossCustomLine] reflection.
//
// Returns:
//   - Point[float64] - A new point representing the reflection of the original point.
func (p Point[float64]) Reflect(axis ReflectionAxis, line ...LineSegment[float64]) Point[float64] {
	switch axis {
	case ReflectAcrossXAxis:
		return NewPoint(p.x, -p.y)
	case ReflectAcrossYAxis:
		return NewPoint(-p.x, p.y)
	case ReflectAcrossCustomLine:
		if len(line) == 0 {
			// If no line is provided, return the point unchanged or handle the error
			return p
		}
		return p.reflectAcrossLine(line[0])
	default:
		return p // Return unchanged if axis is invalid
	}
}

// reflectAcrossLine reflects the point across an arbitrary line defined by a LineSegment.
//
// Parameters:
//   - line: LineSegment[float64] - The line for reflection.
//
// Returns:
//   - Point[float64] - The reflected point.
func (p Point[float64]) reflectAcrossLine(line LineSegment[float64]) Point[float64] {
	// Extract points from the line segment
	x1, y1 := line.start.x, line.start.y
	x2, y2 := line.end.x, line.end.y

	// Calculate the line's slope and intercept for projection
	dx, dy := x2-x1, y2-y1
	if dx == 0 && dy == 0 {
		return p // Degenerate line segment; return point unchanged
	}

	// Calculate the reflection using vector projection
	a := (dx*dx - dy*dy) / (dx*dx + dy*dy)
	b := 2 * dx * dy / (dx*dx + dy*dy)

	newX := a*(p.x-x1) + b*(p.y-y1) + x1
	newY := b*(p.x-x1) - a*(p.y-y1) + y1

	return NewPoint(newX, newY)
}

// RelationshipToCircle determines the spatial relationship of the Point p
// to a given [Circle] c.
//
// This method is a convenience wrapper around [Circle.RelationshipToPoint],
// reversing the parameter order to provide a more intuitive calling syntax
// from the perspective of the [Point].
//
// Parameters:
//   - c ([Circle][T]): The circle to which the relationship is evaluated.
//   - opts: A variadic slice of [Option] functions to customize the behavior of the relationship check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing distances, improving robustness
//     in floating-point calculations.
//
// Returns:
//   - [RelationshipPointCircle]: The relationship of the point to the circle, as defined by the
//     [RelationshipPointCircle] enum.
//
// Notes:
//   - This method directly calls [Circle.RelationshipToPoint], and the result is identical to calling
//     the circle's method with the point as the parameter.
//
// See [Circle.RelationshipToPoint] for examples and more information.
func (p Point[T]) RelationshipToCircle(c Circle[T], opts ...Option) RelationshipPointCircle {
	return c.RelationshipToPoint(p, opts...)
}

// RelationshipToLineSegment determines the spatial relationship of a [Point]
// to a given [LineSegment].
//
// This method evaluates the position of the calling [Point] relative to the
// specified [LineSegment]. The possible relationships are defined by the [RelationshipPointLineSegment].
//
// Parameters:
//   - seg ([LineSegment][T]): The line segment to which the relationship is evaluated.
//   - opts: A variadic slice of [Option] functions to customize the behavior of the relationship check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for floating-point comparisons,
//     improving robustness against precision errors.
//
// Returns:
//   - [RelationshipPointLineSegment]: The relationship of the point to the line segment, as defined by
//     [RelationshipPointLineSegment].
//
// Behavior:
//   - The function first checks if the point coincides with the start or end of the line segment.
//   - It then determines if the point lies on the infinite line defined by the segment using orientation tests.
//   - If the point is collinear with the segment, it checks whether the point lies within the bounds of the segment
//     (using the bounding box).
//   - If the point lies on the infinite line but is outside the bounds of the segment, it is classified as
//     [RelationshipPointLineSegmentCollinearDisjoint].
//
// Notes:
//   - This method is useful for spatial queries where precise relationships between points and line segments are needed.
//
// See [RelationshipPointLineSegment] for possible return values.
func (p Point[T]) RelationshipToLineSegment(seg LineSegment[T], opts ...Option) RelationshipPointLineSegment {
	options := applyOptions(geomOptions{epsilon: 0}, opts...)

	// Check if the point coincides with the segment's start or end
	if p.Eq(seg.start, opts...) {
		return RelationshipPointLineSegmentPointEqStart
	}
	if p.Eq(seg.end, opts...) {
		return RelationshipPointLineSegmentPointEqEnd
	}

	// Check if the point is collinear with the infinite line of the segment
	orientation := Orientation(seg.start, seg.end, p)
	if orientation != PointsCollinear {
		return RelationshipPointLineSegmentMiss
	}

	// Check if the point lies within the bounding box of the segment
	minX, maxX := min(seg.start.x, seg.end.x), max(seg.start.x, seg.end.x)
	minY, maxY := min(seg.start.y, seg.end.y), max(seg.start.y, seg.end.y)

	if float64(p.x) >= float64(minX)-options.epsilon && float64(p.x) <= float64(maxX)+options.epsilon &&
		float64(p.y) >= float64(minY)-options.epsilon && float64(p.y) <= float64(maxY)+options.epsilon {
		return RelationshipPointLineSegmentPointOnLineSegment
	}

	// If collinear but outside the segment's bounds
	return RelationshipPointLineSegmentCollinearDisjoint
}

// RelationshipToPolyTree determines the spatial relationship of a [Point] to a given [PolyTree].
//
// This method evaluates the position of the calling [Point] relative to the specified [PolyTree],
// which can include nested polygons such as solid regions, holes, and islands. The possible relationships
// are defined by [RelationshipPointPolyTree].
//
// Parameters:
//   - tree ([PolyTree][T]): A pointer to the [PolyTree] to analyze.
//   - opts: A variadic slice of [Option] functions to customize the behavior of the relationship check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for floating-point comparisons,
//     improving robustness against precision errors.
//
// Returns:
//   - [RelationshipPointPolyTree]: The relationship of the point to the [PolyTree], as defined by
//     [RelationshipPointPolyTree].
//
// Behavior:
//   - It first checks if the point is on an edge of any polygon in the [PolyTree], returning
//     [PPTRPointOnVertex] or [PPTRPointOnEdge] if a match is found.
//   - It then determines if the point is inside a polygon and evaluates the type of polygon:
//   - [PPTRPointInside] if the point is inside the root polygon.
//   - [PPTRPointInHole] if the point is inside a hole.
//   - [PPTRPointInsideIsland] if the point is inside a nested solid region (island).
//   - If no stronger relationship is found, the function returns [PPTRPointOutside], indicating
//     the point is entirely outside the [PolyTree].
//
// Notes:
//   - This method is particularly useful for complex geometry queries involving hierarchical polygon structures.
//   - The method leverages the efficient iteration of polygons and edges provided by the [PolyTree] structure.
//
// See [RelationshipPointPolyTree] for possible return values.
func (p Point[T]) RelationshipToPolyTree(tree *PolyTree[T], opts ...Option) RelationshipPointPolyTree {
	highestRel := PPTRPointOutside // Default to outside

	// as the points in a polytree contour are doubled, we need to also double the input point
	pDoubled := p.Scale(NewPoint[T](0, 0), 2)

	for poly := range tree.iterPolys {
		// Check if the point is on an edge
		for edge := range poly.contour.iterEdges {
			rel := edge.RelationshipToPoint(pDoubled, opts...)
			switch rel {
			case RelationshipPointLineSegmentPointEqStart, RelationshipPointLineSegmentPointEqEnd:
				return PPTRPointOnVertex // Early return for vertex relationship
			case RelationshipPointLineSegmentPointOnLineSegment:
				return PPTRPointOnEdge // Early return for edge relationship
			}
		}

		// Check if the point is inside the polygon
		if poly.contour.isPointInside(pDoubled) {
			switch {
			case poly.parent == nil:
				highestRel = PPTRPointInside
			case poly.polygonType == PTHole:
				highestRel = PPTRPointInHole
			case poly.polygonType == PTSolid:
				highestRel = PPTRPointInsideIsland
			}
		}
	}

	return highestRel
}

// RelationshipToRectangle determines the spatial relationship of the calling [Point] to a given [Rectangle].
//
// This method evaluates the position of the point relative to the specified rectangle and returns
// a [RelationshipPointRectangle] value indicating whether the point is inside, outside, on an edge,
// or on a vertex of the rectangle.
//
// Parameters:
//   - r ([Rectangle][T]): The rectangle to analyze.
//
// Returns:
//   - [RelationshipPointRectangle]: The relationship of the point to the rectangle, as defined by
//     [RelationshipPointRectangle].
//
// Behavior:
//   - This method delegates the relationship check to [Rectangle.RelationshipToPoint].
//   - The evaluation considers the position of the point relative to the rectangle's edges and vertices.
//
// Notes:
//   - For detailed behavior, examples, and further information, refer to the documentation for
//     [Rectangle.RelationshipToPoint].
//   - This method provides a convenient shorthand for evaluating the relationship between a point and a rectangle
//     without needing to directly call [Rectangle.RelationshipToPoint].
func (p Point[T]) RelationshipToRectangle(r Rectangle[T]) RelationshipPointRectangle {
	return r.RelationshipToPoint(p)
}

// Rotate rotates the point by a specified angle (in radians), counter-clockwise, around a given pivot point.
//
// Parameters:
//   - pivot ([Point][T]): The point around which the rotation is performed.
//   - radians (float64): The angle of rotation, specified in radians.
//   - opts: A variadic slice of [Option] functions to customize the behavior of the relationship check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for floating-point comparisons,
//     improving robustness against precision errors.
//
// Behavior:
//   - The method first translates the point to the origin (relative to the pivot),
//     applies the rotation matrix, and then translates the point back to its original position.
//   - If the [WithEpsilon] option is provided, small numerical deviations in the rotated coordinates
//     will be adjusted based on the specified epsilon.
//
// Returns:
//   - Point[float64]: A new point representing the rotated position.
//
// Notes:
//   - If no options are provided, the default behavior applies, and no epsilon adjustment is made.
//   - The return type is always Point[float64] to ensure precision in the rotated result.
func (p Point[T]) Rotate(pivot Point[T], radians float64, opts ...Option) Point[float64] {
	// Apply options with defaults
	options := applyOptions(geomOptions{epsilon: 0}, opts...)

	pf := p.AsFloat()
	originf := pivot.AsFloat()

	// Translate the point to the origin
	translatedX := pf.x - originf.x
	translatedY := pf.y - originf.y

	// Apply rotation
	rotatedX := translatedX*math.Cos(radians) - translatedY*math.Sin(radians)
	rotatedY := translatedX*math.Sin(radians) + translatedY*math.Cos(radians)

	// Apply epsilon if specified
	if options.epsilon > 0 {
		rotatedX = applyEpsilon(rotatedX, options.epsilon)
		rotatedY = applyEpsilon(rotatedY, options.epsilon)
	}

	// Translate back
	newX := rotatedX + originf.x
	newY := rotatedY + originf.y

	return NewPoint(newX, newY)
}

// Scale scales the point by a factor k relative to a reference point ref.
//
// Parameters:
//   - ref (Point[float64]): The reference point from which scaling originates.
//   - k (float64): The scaling factor.
//
// Returns:
//   - Point[float64] - A new point scaled relative to the reference point.
func (p Point[T]) Scale(ref Point[T], k T) Point[T] {
	return NewPoint(
		ref.x+(p.x-ref.x)*k,
		ref.y+(p.y-ref.y)*k,
	)
}

// String returns a string representation of the Point p in the format "Point[(x, y)]".
// This provides a readable format for the point’s coordinates, useful for debugging
// and displaying points in logs or output.
//
// Returns:
//   - string: A string representation of the Point in the format "Point[(x, y)]".
func (p Point[T]) String() string {
	return fmt.Sprintf("Point[(%v, %v)]", p.x, p.y)
}

// Translate moves the Point by a given displacement vector.
//
// Parameters:
//   - delta (Point[T]): The displacement vector to apply.
//
// Returns:
//   - Point[T]: A new Point resulting from the translation.
func (p Point[T]) Translate(delta Point[T]) Point[T] {
	return NewPoint[T](p.x+delta.x, p.y+delta.y)
}

// X returns the x-coordinate of the Point p.
// This accessor provides read-only access to the x-coordinate.
//
// Returns:
//   - T: The x-coordinate of the Point.
func (p Point[T]) X() T {
	return p.x
}

// Y returns the y-coordinate of the Point p.
// This accessor provides read-only access to the y-coordinate.
//
// Returns:
//   - T: The y-coordinate of the Point.
func (p Point[T]) Y() T {
	return p.y
}

// PointOrientation represents the relative orientation of three points in a two-dimensional plane.
// It describes whether the points are collinear, form a clockwise turn, or form a counterclockwise turn.
// This type is commonly used in computational geometry algorithms to determine the spatial relationship
// between points in relation to each other.
type PointOrientation uint8

// Valid values for PointOrientation.
const (
	// PointsCollinear indicates that the points are collinear, meaning they lie on a single straight line.
	PointsCollinear PointOrientation = iota

	// PointsClockwise indicates that the points are arranged in a clockwise orientation.
	PointsClockwise

	// PointsCounterClockwise indicates that the points are arranged in a counterclockwise orientation.
	PointsCounterClockwise
)

// findLowestLeftestPoint identifies the point with the lowest y-coordinate from a given set of points.
// If multiple points share the lowest y-coordinate, it selects the point with the lowest x-coordinate among them.
//
// Parameters:
//   - points: A variadic list of Point[T] instances from which the lowest leftmost point is determined.
//
// Returns:
//   - int: The index of the lowest leftmost point within the provided points.
//   - Point[T]: The Point with the lowest y-coordinate (and lowest x-coordinate in case of ties).
//
// Example Usage:
//
//	points := []Point[int]{{3, 4}, {1, 5}, {1, 4}}
//	index, lowestPoint := findLowestLeftestPoint(points...)
//	// lowestPoint is Point[int]{1, 4}, and index is 2
func findLowestLeftestPoint[T SignedNumber](points ...Point[T]) (int, Point[T]) {
	lowestPointIndex := 0
	for i := range points {
		switch {
		case points[i].y < points[lowestPointIndex].y:
			lowestPointIndex = i
		case points[i].y == points[lowestPointIndex].y:
			if points[i].x < points[lowestPointIndex].x {
				lowestPointIndex = i
			}
		}
	}
	return lowestPointIndex, points[lowestPointIndex]
}

// orderPointsByAngleAboutLowestPoint sorts a slice of points by their angular order around a reference point, lowestPoint.
// This sorting is used in computational geometry algorithms, such as the Graham scan, to arrange points in a counterclockwise
// order around a pivot point. Collinear points are ordered by increasing distance from the lowestPoint.
//
// Parameters:
//   - lowestPoint: The reference Point from which angles are calculated for sorting. This point is usually the starting point in a convex hull algorithm.
//   - points: A slice of points to be sorted by their angle relative to the lowestPoint.
//
// Sorting Logic:
//   - The function uses the cross product of vectors from lowestPoint to each point to determine the angular order:
//   - If the cross product is positive, a is counterclockwise to b.
//   - If the cross product is negative, a is clockwise to b.
//   - If the cross product is zero, the points are collinear, so they are sorted by their distance to lowestPoint.
//
// Example Usage:
//
//	points := []Point[int]{{3, 4}, {1, 5}, {2, 2}}
//	lowestPoint := NewPoint(1, 2)
//	orderPointsByAngleAboutLowestPoint(lowestPoint, points)
//	// points are now sorted counterclockwise around lowestPoint, with collinear points ordered by distance.
func orderPointsByAngleAboutLowestPoint[T SignedNumber](lowestPoint Point[T], points []Point[T]) {
	slices.SortStableFunc(points, func(a Point[T], b Point[T]) int {

		// Ensure lowestPoint is always the first point
		switch {
		case a.Eq(lowestPoint):
			return -1
		case b.Eq(lowestPoint):
			return 1
		}

		// Calculate relative vectors from lowestPoint to start and end
		relativeA := a.Translate(lowestPoint.Negate())
		relativeB := b.Translate(lowestPoint.Negate())
		crossProduct := relativeA.CrossProduct(relativeB)

		// Use cross product to determine angular order
		switch {
		case crossProduct > 0:
			return -1 // start is counterclockwise to end
		case crossProduct < 0:
			return 1 // start is clockwise to end
		}

		// If cross product is zero, points are collinear; order by distance to lowestPoint
		distAtoLP := lowestPoint.DistanceSquaredToPoint(a)
		distBtoLP := lowestPoint.DistanceSquaredToPoint(b)

		// Sort closer points first
		switch {
		case distAtoLP < distBtoLP:
			return -1
		case distAtoLP > distBtoLP:
			return 1
		default:
			return 0
		}
	})
}

// triangleAreaX2Signed calculates twice the signed area of the triangle formed by points p0, p1, and p2.
// The result is positive if the points are in counterclockwise order, negative if clockwise, and zero if collinear.
// This method is useful in computational geometry for determining point orientation and triangle area without
// computing floating-point values.
//
// Parameters:
//   - p0, p1, p2: The three points that define the triangle.
//
// Returns:
//   - T: Twice the signed area of the triangle. This value helps determine the relative orientation of the points.
//
// The formula here uses the 2D cross product of vectors (p1 - p0) and (p2 - p0)
// to compute twice the signed area of the triangle formed by p0, p1, p2.
func triangleAreaX2Signed[T SignedNumber](p0, p1, p2 Point[T]) T {
	return (p1.x-p0.x)*(p2.y-p0.y) - (p2.x-p0.x)*(p1.y-p0.y)
}

// EnsureClockwise ensures that a slice of points representing a polygon is ordered in a clockwise direction.
//
// This function checks the orientation of the points based on the signed area of the polygon.
// If the signed area is positive, indicating a counterclockwise orientation, the function reverses
// the order of the points (in-place) to make them clockwise. If the points are already clockwise, no changes are made.
//
// Parameters:
//   - points ([][Point][T]): A slice of points representing the vertices of a polygon. The points are assumed
//     to form a closed loop or define a valid polygon.
//
// Behavior:
//   - Calculates the signed area of the polygon using [SignedArea2X].
//   - If the signed area is positive (counterclockwise orientation), reverses the order of the points.
//   - If the signed area is negative or zero (already clockwise or degenerate polygon), does nothing.
//
// Notes:
//   - This function modifies the input slice of points in place.
//   - A zero area polygon is considered already clockwise and is left unchanged.
func EnsureClockwise[T SignedNumber](points []Point[T]) {
	if SignedArea2X(points) < 0 {
		return // Already clockwise
	}
	slices.Reverse(points)
}

// EnsureCounterClockwise ensures that a slice of points representing a polygon is ordered in a counterclockwise direction.
//
// This function checks the orientation of the points based on the signed area of the polygon.
// If the signed area is negative, indicating a clockwise orientation, the function reverses
// the order of the points (in-place) to make them counterclockwise. If the points are already counterclockwise,
// no changes are made.
//
// Parameters:
//
//   - points ([][Point][T]): A slice of points representing the vertices of a polygon. The points are assumed
//     to form a closed loop or define a valid polygon.
//
// Behavior:
//
//   - Calculates the signed area of the polygon using [SignedArea2X].
//   - If the signed area is negative (clockwise orientation), reverses the order of the points.
//   - If the signed area is positive or zero (already counterclockwise or degenerate polygon), does nothing.
//
// Notes:
//
//   - This function modifies the input slice of points in place.
//   - A zero area polygon is considered already counterclockwise and is left unchanged.
func EnsureCounterClockwise[T SignedNumber](points []Point[T]) {
	if SignedArea2X(points) > 0 {
		return // Already counterclockwise
	}
	slices.Reverse(points)
}

// Orientation determines the relative orientation of three points: p0, p1, and p2.
// It calculates the signed area of the triangle formed by these points to determine if the
// points make a counterclockwise turn, a clockwise turn, or are collinear. This method is
// widely used in computational geometry to classify point arrangements in polygon and convex hull algorithms.
//
// Parameters:
//   - p0, p1, p2: The three points for which the orientation is determined.
//
// Returns:
//
// PointOrientation: A constant indicating the relative orientation of the points:
//   - PointsCounterClockwise if the points form a counterclockwise turn,
//   - PointsClockwise if the points form a clockwise turn,
//   - PointsCollinear if the points are collinear (lie on a single line).
//
// Example Usage:
//
//	p0 := NewPoint(0, 0)
//	p1 := NewPoint(1, 1)
//	p2 := NewPoint(2, 0)
//	orientation := Orientation(p0, p1, p2) // orientation will be PointsClockwise in this case
func Orientation[T SignedNumber](p0, p1, p2 Point[T]) PointOrientation {
	area2x := triangleAreaX2Signed(p0, p1, p2)
	switch {
	case area2x < 0:
		return PointsClockwise
	case area2x > 0:
		return PointsCounterClockwise
	default: // area2x == 0
		return PointsCollinear
	}
}

// RelativeAngle calculates the angle in radians between two points, A and B, relative to an optional origin [Point] O.
// This angle is measured from the origin to the line segments connecting A and B.
// If no origin is provided, the function defaults to using the point (0, 0) as the origin.
//
// Parameters:
//
//   - A: The first [Point] forming one side of the angle.
//   - B: The second [Point] forming the other side of the angle.
//   - O: An optional origin [Point]. If not provided, the origin defaults to (0, 0).
//
// Returns:
//
//   - radians (float64): The angle between points A and B relative to the origin, in radians.
func RelativeAngle[T SignedNumber](A, B Point[T], O ...Point[T]) (radians float64) {
	return math.Acos(RelativeCosineOfAngle(A, B, O...))
}

// RelativeCosineOfAngle calculates the cosine of the angle between two points, A and B, relative to an optional origin [Point] O.
// This function returns the cosine of the angle directly, avoiding the costly [math.Acos] calculation, which can improve performance
// in applications where only the cosine value is needed.
//
// If no origin is provided, the function defaults to using the point (0, 0) as the origin.
//
// Parameters:
//   - A ([Point][T]): The first [Point] forming one side of the angle.
//   - B ([Point][T]): The second [Point] forming the other side of the angle.
//   - O ([Point][T]): An optional origin [Point]. If not provided, the origin defaults to (0, 0).
//
// Returns:
//   - float64: The cosine of the angle between points A and B relative to the origin.
//
// Note:
//   - This function does not currently handle division by zero errors. If either vector OA or OB has zero length,
//     a division by zero could occur. Consider adding validation if needed in such cases.
//
// Why Would Anyone Just Need The Cosine?
//
// By using the cosine of the angle, you can determine not just the relative angle but also its directionality and
// magnitude in terms of alignment. Here's why someone might want this:
//
// 1. Efficient Angle Comparison
//
// Instead of calculating the actual angle using trigonometric functions (which are computationally expensive), you can
// compare the cosine of angles directly. Since the cosine function is monotonic between 0 and π, you can use the cosine
// values to determine:
//   - If the vectors are pointing in the same direction (cos(θ) ≈ 1).
//   - If the vectors are orthogonal (cos(θ) ≈ 0).
//   - If the vectors are pointing in opposite directions (cos(θ) ≈ -1).
//
// 2. Avoiding Floating-Point Inaccuracies
//
// Computing the cosine of the angle avoids potential inaccuracies associated with computing the angle itself
// (e.g., due to limited precision when converting radians to degrees or vice versa).
//
// 3. Applications in Sorting
//
// If you are sorting points or vectors relative to a reference direction, you can use [RelativeCosineOfAngle] to order
// them based on their angular relationship. This is particularly useful in computational geometry tasks like:
//   - Constructing a convex hull.
//   - Ordering vertices for polygon operations.
//
// 4. Determining Vector Orientation
//
// You can use the cosine value to determine the degree of alignment between two vectors, which is helpful in:
//   - Physics simulations (e.g., checking if a force vector aligns with a velocity vector).
//   - Rendering graphics (e.g., checking the alignment of a surface normal with a light source).
//
// 5. Geometric Insight
//
// In geometry, understanding the relative cosine helps in:
//   - Classifying angles (acute, obtuse, or right) without explicitly calculating them.
//   - Performing dot product-based calculations indirectly, as cos(θ) is derived from the dot product normalized by the vectors' magnitudes.
//
// [math.Acos]: https://pkg.go.dev/math#Acos
func RelativeCosineOfAngle[T SignedNumber](A, B Point[T], O ...Point[T]) float64 {
	// Set origin point to (0,0) if not provided
	origin := NewPoint[T](0, 0)
	if len(O) > 0 {
		origin = O[0]
	}

	// Calculate vectors OA and OB
	vectorOA := origin.Translate(A.Negate())
	vectorOB := origin.Translate(B.Negate())

	// Calculate the Dot Product of OA and OB
	OAdotOB := vectorOA.DotProduct(vectorOB)

	// Calculate the Magnitude of Each Vector
	magnitudeOA := origin.DistanceToPoint(A)
	magnitudeOB := origin.DistanceToPoint(B)

	// todo: check for divide by zero errors & handle (return error?)

	// Use the Dot Product Formula to Find the Cosine of the Angle
	return float64(OAdotOB) / (magnitudeOA * magnitudeOB)
}

// SignedArea2X computes twice the signed area of a polygon defined by the given points.
//
// The function calculates the signed area of the polygon using the [Shoelace Formula],
// adapted to sum the areas of triangles formed by consecutive points. The result is
// twice the actual signed area, which avoids introducing fractional values and simplifies
// calculations with integer-based coordinate types.
//
// A positive signed area indicates that the points are ordered counterclockwise,
// while a negative signed area indicates clockwise order. This function is commonly
// used to determine the orientation of a polygon or to compute its area efficiently.
//
// Parameters:
//   - points ([][Point][T]): A slice of [Point] values representing the vertices of the polygon in order.
//     The polygon is assumed to be closed, meaning the first [Point] connects to the last [Point].
//
// Returns:
//   - The signed area multiplied by 2 (hence "2X").
//     Returns 0 if the number of points is fewer than 3, as a valid polygon cannot be formed.
//
// Notes:
//   - The function assumes the input points form a simple polygon (no self-intersections).
//   - If the points are not in order, the result may not represent the true orientation
//     or area of the intended polygon.
//
// [Shoelace Formula]: https://en.wikipedia.org/wiki/Shoelace_formula
func SignedArea2X[T SignedNumber](points []Point[T]) T {
	var area T
	n := len(points)
	if n < 3 {
		return 0 // Not a polygon
	}
	for i := 1; i < n-1; i++ {
		area += triangleAreaX2Signed(points[0], points[i], points[i+1])
	}
	return area
}
