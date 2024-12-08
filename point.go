package geom2d

import (
	"fmt"
	"image"
	"math"
	"slices"
)

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

// Point represents a point in two-dimensional space with x and y coordinates of a generic numeric type T.
// The Point struct provides methods for common vector operations such as addition, subtraction, and distance
// calculations, making it versatile for computational geometry and graphics applications.
//
// Type Parameter:
//   - T: The numeric type for the coordinates, constrained to signed number types by the SignedNumber interface.
//
// Usage:
//   - To create a new Point, use the NewPoint constructor: p := NewPoint(3, 4)
//   - To create from an image.Point, use NewPointFromImagePoint: p := NewPointFromImagePoint(imagePoint)
//
// Accessor Methods:
//   - p.X(): Returns the x-coordinate of the point.
//   - p.Y(): Returns the y-coordinate of the point.
type Point[T SignedNumber] struct {
	x T
	y T
}

// AsFloat converts the Point's x and y coordinates to the float64 type, returning a new Point[float64].
// This method is useful when higher precision or floating-point arithmetic is needed on the coordinates.
//
// Returns:
//   - Point[float64]: A new Point with x and y coordinates converted to float64.
//
// Example Usage:
//
//	p := NewPoint(3, 4)
//	floatPoint := p.AsFloat() // floatPoint is a Point[float64] with coordinates (3.0, 4.0)
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
//
// Example Usage:
//
//	p := NewPoint(3.7, 4.9)
//	intPoint := p.AsInt() // intPoint is a Point[int] with coordinates (3, 4)
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
//
// Example Usage:
//
//	p := NewPoint(3.7, 4.2)
//	roundedPoint := p.AsIntRounded() // roundedPoint is a Point[int] with coordinates (4, 4)
func (p Point[T]) AsIntRounded() Point[int] {
	return Point[int]{
		x: int(math.Round(float64(p.x))),
		y: int(math.Round(float64(p.y))),
	}
}

// CrossProduct calculates the cross product of the vector from the origin to Point p and the vector from the origin to Point q.
// The result is useful in determining the relative orientation of two vectors:
//   - A positive result indicates a counterclockwise turn (left turn),
//   - A negative result indicates a clockwise turn (right turn),
//   - A result of zero indicates that the points are collinear.
//
// Parameters:
//   - q: The Point with which to compute the cross product relative to the calling Point.
//
// Returns:
//   - T: The cross product of the vectors from the origin to p and q, indicating their relative orientation.
//
// Example Usage:
//
//	p := NewPoint(1, 2)
//	q := NewPoint(3, 4)
//	cross := p.CrossProduct(q) // cross > 0 means p and q form a counterclockwise turn
func (p Point[T]) CrossProduct(q Point[T]) T {
	return (p.x * q.y) - (p.y * q.x)
}

// DistanceToLineSegment calculates the orthogonal (shortest) distance from the current Point p
// to a specified LineSegment l. This distance is the length of the perpendicular line segment
// from p to the closest point on l.
//
// Parameters:
//   - l: The LineSegment to which the distance is calculated.
//   - opts: A variadic slice of Option functions to customize the calculation behavior.
//     WithEpsilon(epsilon float64): Adjusts the result by snapping small floating-point
//     deviations to cleaner values based on the specified epsilon threshold.
//
// Behavior:
//   - The function first computes the projection of p onto the given line segment l. This is
//     the closest point on l to p, whether it falls within the line segment or on one of its endpoints.
//   - The distance is then calculated as the Euclidean distance from p to the projected point,
//     using the DistanceToPoint method for precision.
//
// Returns:
//   - float64: The shortest distance between the point p and the line segment l, optionally
//     adjusted based on epsilon if provided.
//
// Example Usage:
//
//	p := NewPoint(3, 4)
//	l := NewLineSegment(NewPoint(0, 0), NewPoint(6, 8))
//
//	// Default behavior (no epsilon adjustment)
//	distance := p.DistanceToLineSegment(l)
//
//	// With epsilon adjustment
//	distanceWithEpsilon := p.DistanceToLineSegment(l, WithEpsilon(1e-4))
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
//   - q: The Point to which the distance is calculated.
//   - opts: A variadic slice of Option functions to customize the behavior of the calculation.
//     WithEpsilon(epsilon float64): Adjusts the result by snapping near-zero distances
//     (or distances close to clean values like integers) to a more precise value based
//     on the specified epsilon threshold.
//
// Behavior:
//   - The function computes the Euclidean distance using the formula:
//     distance = sqrt((p.x - q.x)^2 + (p.y - q.y)^2)
//   - If the `WithEpsilon` option is provided, the computed distance is adjusted such that
//     deviations within the epsilon range are snapped to a clean value.
//
// Returns:
//   - float64: The Euclidean distance between the two points, optionally adjusted based on epsilon.
//
// Example Usage:
//
//	p := NewPoint(3, 4)
//	q := NewPoint(6, 8)
//
//	// Default behavior (no epsilon adjustment)
//	distance := p.DistanceToPoint(q) // Straight-line distance between p and q
//
//	// With epsilon adjustment
//	distanceWithEpsilon := p.DistanceToPoint(q, WithEpsilon(1e-4))
//	// Adjusts small floating-point deviations based on epsilon
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

// DistanceSquaredToPoint calculates the squared Euclidean distance between Point p and another Point q.
// This method returns the squared distance, which avoids the computational cost of a square root calculation
// and is useful in cases where only distance comparisons are needed.
//
// Parameters:
//   - q: The Point to which the squared distance is calculated from the calling Point.
//
// Returns:
//   - T: The squared Euclidean distance between p and q.
//
// Example Usage:
//
//	p := NewPoint(3, 4)
//	q := NewPoint(6, 8)
//	distanceSquared := p.DistanceSquaredToPoint(q) // distanceSquared is the squared distance between p and q
func (p Point[T]) DistanceSquaredToPoint(q Point[T]) T {
	return (q.x-p.x)*(q.x-p.x) + (q.y-p.y)*(q.y-p.y)
}

// Div scales the calling Point p by dividing each of its coordinates by a scalar value k.
// The method returns a new Point of type Point[float64], ensuring that fractional values
// are preserved. This is particularly useful for applications where precise scaling is required.
//
// Parameters:
//   - k: The scalar value by which to divide the x and y coordinates of Point p.
//   - opts: A variadic slice of Option functions to customize the behavior of the division.
//     WithEpsilon(epsilon float64): Adjusts the result by snapping small floating-point
//     deviations to cleaner values based on the specified epsilon threshold.
//
// Behavior:
//   - Each coordinate of the Point is divided by k, resulting in fractional values where appropriate.
//   - If the `WithEpsilon` option is provided, the resulting coordinates are adjusted such that
//     deviations within the epsilon range are snapped to clean values (e.g., rounding near-zero values to zero).
//
// Returns:
//   - Point[float64]: A new Point with coordinates scaled by the division and optionally adjusted based on epsilon.
//
// Example Usage:
//
//	p := NewPoint(10, 20)
//
//	// Default behavior (no epsilon adjustment)
//	scaledPoint := p.Div(2) // scaledPoint is a Point[float64] with coordinates (5.0, 10.0)
//
//	// With epsilon adjustment
//	scaledPointWithEpsilon := p.Div(2, WithEpsilon(1e-4))
//	// Adjusts small floating-point deviations in the result based on epsilon
//
// Notes:
//   - If k is zero, this function will cause a runtime panic due to division by zero.
//     Ensure that k is non-zero before calling this method.
//   - Epsilon adjustment is particularly useful when dividing by values that might introduce
//     floating-point imprecisions in the resulting coordinates.
func (p Point[T]) Div(k T, opts ...Option) Point[float64] {

	// Apply options with defaults
	options := applyOptions(geomOptions{epsilon: 0}, opts...)

	divided := Point[float64]{
		x: float64(p.x) / float64(k),
		y: float64(p.y) / float64(k),
	}

	if options.epsilon > 0 {
		divided.x = applyEpsilon(divided.x, options.epsilon)
		divided.y = applyEpsilon(divided.y, options.epsilon)
	}

	return divided
}

// DotProduct calculates the dot product of the vector represented by Point p with the vector represented by Point q.
// The dot product is defined as p.x*q.x + p.y*q.y and is widely used in geometry for angle calculations,
// projection operations, and determining the relationship between two vectors.
//
// Parameters:
//   - q: The Point with which to calculate the dot product relative to the calling Point.
//
// Returns:
//   - T: The dot product of the vectors represented by p and q.
//
// Example Usage:
//
//	p := NewPoint(1, 2)
//	q := NewPoint(3, 4)
//	dot := p.DotProduct(q) // dot is the dot product of p and q
func (p Point[T]) DotProduct(q Point[T]) T {
	return (p.x * q.x) + (p.y * q.y)
}

// Eq determines whether the calling Point p is equal to another Point q.
// Equality can be evaluated either exactly (default) or approximately using an epsilon threshold.
//
// Parameters:
//   - q: The Point to compare with the calling Point.
//   - opts: A variadic slice of Option functions to customize the equality check.
//     WithEpsilon(epsilon float64): Specifies a tolerance for comparing the coordinates
//     of p and q. If the absolute difference between the coordinates of p and q is less
//     than epsilon, they are considered equal.
//
// Behavior:
//   - By default, the function performs an exact equality check, returning true only if
//     the x and y coordinates of p and q are identical.
//   - If the `WithEpsilon` option is provided, the function performs an approximate equality
//     check, considering p and q equal if their coordinate differences are within the specified
//     epsilon threshold.
//
// Returns:
//   - bool: True if p and q are equal based on the specified comparison mode; otherwise, false.
//
// Example Usage:
//
//	p := NewPoint(3.0, 4.0)
//	q := NewPoint(3.0, 4.0)
//
//	// Exact equality
//	isEqual := p.Eq(q) // isEqual is true
//
//	// Approximate equality with epsilon
//	r := NewPoint(3.000001, 4.000001)
//	isApproximatelyEqual := p.Eq(r, WithEpsilon(1e-5)) // isApproximatelyEqual is true
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

// ProjectOntoLineSegment projects the Point p onto a given LineSegment l.
//
// The function calculates the closest point on the LineSegment defined by
// the endpoints l.Start() and l.End() to the Point p. It utilizes vector
// mathematics to determine the projection of Point p onto the infinite line
// represented by the LineSegment. If the projected point falls beyond the ends of the
// LineSegment, the function returns the closest endpoint of the segment.
//
// Parameters:
//   - l: The line segment onto which the point is projected. It consists of
//     two endpoints: Start() and End().
//
// Returns:
//   - A Point[float64] representing the coordinates of the projected point.
//     If the LineSegment is degenerate (both endpoints are the same),
//     the function returns the coordinates of the Start() Point of the LineSegment.
func (p Point[T]) ProjectOntoLineSegment(l LineSegment[T]) Point[float64] {

	// the direction vector of the line segment
	vecAB := l.end.Sub(l.start)

	// the vector from line segment start to point
	vecAP := p.Sub(l.start)

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
//   - axis: Axis - The axis of reflection (ReflectAcrossXAxis, ReflectAcrossYAxis, or ReflectAcrossCustomLine).
//   - line: Optional LineSegment[T] - The line for ReflectAcrossCustomLine reflection.
//
// Returns:
//   - Point[float64] - A new point representing the reflection of the original point.
//
// Example usage:
//
//	p := NewPoint(3, 4)
//	reflectedX := p.Reflect(ReflectAcrossXAxis)               // Reflect across x-axis: (3, -4)
//	customLine := NewLineSegment(NewPoint(0, 0), NewPoint(1, 1))
//	reflectedLine := p.Reflect(ReflectAcrossCustomLine, customLine) // Reflect across y = x
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

func (p Point[T]) RelationshipToCircle(c Circle[T], opts ...Option) PointCircleRelationship {
	return c.RelationshipToPoint(p, opts...)
}

func (p Point[T]) RelationshipToLineSegment(seg LineSegment[T], opts ...Option) PointLineSegmentRelationship {
	options := applyOptions(geomOptions{epsilon: 0}, opts...)

	// Check if the point coincides with the segment's start or end
	if p.Eq(seg.start, opts...) {
		return PLRPointEqStart
	}
	if p.Eq(seg.end, opts...) {
		return PLRPointEqEnd
	}

	// Check if the point is collinear with the infinite line of the segment
	orientation := Orientation(seg.start, seg.end, p)
	if orientation != PointsCollinear {
		return PLRMiss
	}

	// Check if the point lies within the bounding box of the segment
	minX, maxX := min(seg.start.x, seg.end.x), max(seg.start.x, seg.end.x)
	minY, maxY := min(seg.start.y, seg.end.y), max(seg.start.y, seg.end.y)

	if float64(p.x) >= float64(minX)-options.epsilon && float64(p.x) <= float64(maxX)+options.epsilon &&
		float64(p.y) >= float64(minY)-options.epsilon && float64(p.y) <= float64(maxY)+options.epsilon {
		return PLRPointOnLineSegment
	}

	// If collinear but outside the segment's bounds
	return PLRPointOnLine
}

func (p Point[T]) RelationshipToPolyTree(tree *PolyTree[T], opts ...Option) PointPolyTreeRelationship {
	highestRel := PPTRPointOutside // Default to outside

	// as the points in a polytree contour are doubled, we need to also double the input point
	pDoubled := p.Scale(NewPoint[T](0, 0), 2)

	for poly := range tree.iterPolys {
		// Check if the point is on an edge
		for edge := range poly.contour.iterEdges {
			rel := edge.RelationshipToPoint(pDoubled, opts...)
			switch rel {
			case PLRPointEqStart, PLRPointEqEnd:
				return PPTRPointOnVertex // Early return for vertex relationship
			case PLRPointOnLineSegment:
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

func (p Point[T]) RelationshipToRectangle(r Rectangle[T]) PointRectangleRelationship {
	return r.RelationshipToPoint(p)
}

// Rotate rotates the point by a specified angle (in radians) around a given pivot point.
//
// Parameters:
//   - pivot: Point[T] - The point around which the rotation is performed.
//   - radians: float64 - The angle of rotation, specified in radians.
//   - opts: ...Option - Functional options for customizing the rotation behavior. WithEpsilon(epsilon float64): Specifies a tolerance for rounding small floating-point deviations to cleaner values (e.g., snapping near-zero values to zero).
//
// Behavior:
//   - The method first translates the point to the origin (relative to the pivot),
//     applies the rotation matrix, and then translates the point back to its original position.
//   - If the `WithEpsilon` option is provided, small numerical deviations in the rotated coordinates
//     (e.g., -0.9999999999999998 instead of -1) will be adjusted based on the specified epsilon.
//
// Returns:
//   - Point[float64] - A new point representing the rotated position.
//
// Example Usage:
//
//	pivot := geom2d.NewPoint(1.0, 1.0)
//	circle := geom2d.NewCircle(geom2d.NewPoint(3.0, 3.0), 5.0)
//
//	// Rotates the circle 90 degrees around (1.0, 1.0)
//	rotatedCircle := circle.Rotate(pivot, math.Pi/2, geom2d.WithEpsilon(1e-10))
//
//	// rotatedCircle = Circle[center=(-1, 3), radius=5]
//
// Notes:
//   - If no options are provided, the default behavior applies, and no epsilon adjustment is made.
//   - The return type is always `Point[float64]` to ensure precision in the rotated result.
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
//   - ref: Point[float64] - The reference point from which scaling originates.
//   - k: float64 - The scaling factor.
//
// Returns:
//   - Point[float64] - A new point scaled relative to the reference point.
//
// Example:
//
//	p := NewPoint(3, 4)
//	ref := NewPoint(1, 1)
//	scaled := p.Scale(ref, 2) // scaled is now (5, 7), relative to ref.
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
//
// Example Usage:
//
//	p := NewPoint(3, 4)
//	str := p.String() // str is "Point[(3, 4)]"
func (p Point[T]) String() string {
	return fmt.Sprintf("Point[(%v, %v)]", p.x, p.y)
}

// Sub returns a new Point representing the vector from the calling Point p to another Point q.
// This is equivalent to subtracting the coordinates of q from those of p, resulting in a vector
// (as a Point) that points from q to p.
//
// Parameters:
//   - q: The Point to subtract from the calling Point.
//
// Returns:
//   - Point[T]: A new Point representing the vector from p to q, with coordinates equal to p - q.
//
// Example Usage:
//
//	p := NewPoint(5, 7)
//	q := NewPoint(3, 2)
//	vector := p.Sub(q) // vector is a Point with coordinates (2, 5), representing the vector from q to p
func (p Point[T]) Sub(q Point[T]) Point[T] {
	return Point[T]{
		x: p.x - q.x,
		y: p.y - q.y,
	}
}

// Translate moves the Point by a given displacement vector.
//
// Parameters:
//   - delta: Point[T] - The displacement vector to apply.
//
// Returns:
//   - Point[T]: A new Point resulting from the translation.
//
// Example Usage:
//
//	p := NewPoint(3, 4)
//	delta := NewPoint(2, -1)
//	translated := p.Translate(delta) // translated is a Point with coordinates (5, 3)
func (p Point[T]) Translate(delta Point[T]) Point[T] {
	return NewPoint(p.x+delta.x, p.y+delta.y)
}

// X returns the x-coordinate of the Point p.
// This accessor provides read-only access to the x-coordinate.
//
// Returns:
//   - T: The x-coordinate of the Point.
//
// Example Usage:
//
//	p := NewPoint(3, 4)
//	x := p.X() // x is 3
func (p Point[T]) X() T {
	return p.x
}

// Y returns the y-coordinate of the Point p.
// This accessor provides read-only access to the y-coordinate.
//
// Returns:
//   - T: The y-coordinate of the Point.
//
// Example Usage:
//
//	p := NewPoint(3, 4)
//	y := p.Y() // y is 4
func (p Point[T]) Y() T {
	return p.y
}

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
		relativeA := a.Sub(lowestPoint)
		relativeB := b.Sub(lowestPoint)
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

// ConvexHull computes the convex hull of a finite set of points using the Graham scan algorithm.
// The convex hull is the smallest convex polygon that encloses all points in the input set. This function is particularly useful in computational
// geometry applications where the outer boundary of a set of points is needed.
//
// This implementation finds the point with the lowest y-coordinate to serve as a reference for sorting points by their angle relative to this point.
// Starting from this reference, it orders the points counterclockwise, removing any points that cause a clockwise turn to ensure a convex boundary.
//
// Parameters:
//   - points: A variable number of `Point[T]` instances representing the set of points for which the convex hull is to be computed.
//
// Returns:
//   - []Point[T]: A slice of points representing the vertices of the convex hull in counterclockwise order.
//     The returned slice includes only the points that form the outer boundary of the input set.
//
// Note:
//   - If the `points` parameter is empty or has fewer than three points, the function returns the input points unchanged, as a convex hull cannot be formed.
//
// See https://en.wikipedia.org/wiki/Graham_scan for more information on the Graham scan algorithm.
//
// See https://en.wikipedia.org/wiki/Convex_hull for more information on convex hulls.
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

// NewPoint creates and returns a new Point with the specified x and y coordinates.
// This function is a generic constructor for Point, allowing the x and y values to be of any signed numeric type.
//
// Parameters:
//   - x: The x-coordinate of the new point.
//   - y: The y-coordinate of the new point.
//
// Returns:
//   - Point[T]: A new Point instance with the specified coordinates.
func NewPoint[T SignedNumber](x, y T) Point[T] {
	return Point[T]{
		x: x,
		y: y,
	}
}

// NewPointFromImagePoint creates and returns a new Point with integer x and y coordinates
// based on an image.Point. This function is useful for converting between graphics and
// computational geometry representations of points.
//
// Parameters:
//   - q: An image.Point representing the source coordinates for the new point.
//
// Returns:
//   - Point[int]: A new Point with coordinates corresponding to the x and y values of the provided image.Point.
func NewPointFromImagePoint(q image.Point) Point[int] {
	return Point[int]{
		x: q.X,
		y: q.Y,
	}
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
	vectorOA := origin.Sub(A)
	vectorOB := origin.Sub(B)

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
