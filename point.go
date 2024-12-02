// This file contains the implementation of the Point type, which represents a point in 2D space.
// It includes methods for basic arithmetic operations, transformations, relationships with other types,
// and utilities like equality checks and conversions between numeric types.
//
// The Point type is designed to be flexible and efficient, leveraging generics to support various numeric
// types (e.g., int, float64). It serves as a building block for more complex geometric types like LineSegment,
// Circle, and Rectangle.

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
	PointsCollinear = PointOrientation(iota)

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

// Add returns a new Point that represents the sum of the calling Point p and another Point q.
// Each coordinate of the result is the sum of the corresponding coordinates of p and q.
//
// Parameters:
//   - q: The Point to add to the calling Point.
//
// Returns:
//   - Point[T]: A new Point where the x and y coordinates are the sums of the x and y coordinates of p and q.
//
// Example Usage:
//
//	p1 := NewPoint(3, 4)
//	p2 := NewPoint(1, 2)
//	result := p1.Add(p2) // result is a Point with coordinates (4, 6)
func (p Point[T]) Add(q Point[T]) Point[T] {
	return Point[T]{
		x: p.x + q.x,
		y: p.y + q.y,
	}
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

// DistanceToLineSegment calculates the orthogonal distance from the Point p
// to the specified LineSegment l. This distance is determined by projecting
// the Point p onto the LineSegment and measuring the distance from p to
// the projected Point.
//
// Returns the distance as a float64.
//
// The function first computes the projection of p onto the line segment defined
// by the endpoints of l. It then converts the original point p to a
// Point[float64] to ensure accurate distance calculation, as the
// DistanceToPoint function operates on float64 points.
func (p Point[T]) DistanceToLineSegment(l LineSegment[T]) float64 {
	projectedPoint := p.ProjectOntoLineSegment(l)
	pf := p.AsFloat()
	return pf.DistanceToPoint(projectedPoint)
}

// DistanceToPoint calculates the Euclidean (straight-line) distance between Point p and another Point q.
// This method returns the distance as a float64, allowing for precise measurement of the straight-line
// separation between two points.
//
// Parameters:
//   - q: The Point to which the distance is calculated from the calling Point.
//
// Returns:
//   - float64: The Euclidean distance between p and q.
//
// Example Usage:
//
//	p := NewPoint(3, 4)
//	q := NewPoint(6, 8)
//	distance := p.DistanceToPoint(q) // distance is the straight-line distance between p and q
func (p Point[T]) DistanceToPoint(q Point[T]) float64 {
	return math.Sqrt(float64(p.DistanceSquaredToPoint(q)))
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

// Div returns a new Point that scales the calling Point p by dividing each of its coordinates by a scalar value k.
// This method performs division and returns a Point of type Point[float64] to preserve any fractional values,
// making it suitable for applications where precise scaling is required.
//
// Parameters:
//   - k: The scalar value by which to divide the x and y coordinates of Point p.
//
// Returns:
//   - Point[float64]: A new Point with x and y coordinates scaled by the division, as floating-point values.
//
// Example Usage:
//
//	p := NewPoint(10, 20)
//	scaledPoint := p.Div(2) // scaledPoint is a Point[float64] with coordinates (5.0, 10.0)
func (p Point[T]) Div(k T) Point[float64] {
	return Point[float64]{
		x: float64(p.x) / float64(k),
		y: float64(p.y) / float64(k),
	}
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

// Eq reports whether the calling Point p is equal to another Point q.
// Two points are considered equal if both their x and y coordinates are identical.
//
// Parameters:
//   - q: The Point to compare with the calling Point.
//
// Returns:
//   - bool: True if the x and y coordinates of p and q are equal; otherwise, false.
//
// Example Usage:
//
//	p := NewPoint(3, 4)
//	q := NewPoint(3, 4)
//	isEqual := p.Eq(q) // isEqual is true because p and q have the same coordinates
func (p Point[T]) Eq(q Point[T]) bool {
	return p.x == q.x && p.y == q.y
}

// IsOnLineSegment determines whether a Point lies on a given LineSegment.
//
// This method first checks if the point is collinear with the endpoints of the
// line segment using an orientation check. If the point is not collinear, it
// cannot be on the segment. If the point is collinear, the function then verifies
// if the point lies within the bounding box defined by the segment's endpoints.
//
// Parameters:
//   - l: LineSegment[T], the line segment against which the Point's position is tested.
//
// Returns:
//   - bool: true if the Point lies on the LineSegment, false otherwise.
//
// Example usage:
//
//	p := NewPoint[float64](1, 1)
//	segment := NewLineSegment(NewPoint[float64](0, 0), NewPoint[float64](2, 2))
//	onSegment := p.IsOnLineSegment(segment)  // true as p lies on segment
func (p Point[T]) IsOnLineSegment(l LineSegment[T]) bool {

	// Check collinearity first; if not collinear, point is not on the line segment
	if Orientation(p, l.start, l.end) != PointsCollinear {
		return false
	}

	// Check if the point lies within the bounding box defined by A and End
	return p.x >= min(l.start.x, l.end.x) && p.x <= max(l.start.x, l.end.x) &&
		p.y >= min(l.start.y, l.end.y) && p.y <= max(l.start.y, l.end.y)
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

// Rotate rotates the point by a specified angle (in radians) around a given origin point.
//
// Parameters:
//   - origin: Point[T] - The origin point around which to rotate.
//   - radians: float64 - The rotation angle in radians.
//
// Returns:
//   - Point[float64] - A new point representing the rotated position.
//
// Example usage:
//
//	p := NewPoint[float64](3.0, 4.0)
//	origin := NewPoint[float64](1.0, 1.0)
//	rotated := p.Rotate(origin, math.Pi / 2) // Rotates p 90 degrees around origin
func (p Point[T]) Rotate(pivot Point[T], radians float64) Point[float64] {
	pf := p.AsFloat()
	originf := pivot.AsFloat()

	// Translate the point to the origin
	translatedX := pf.x - originf.x
	translatedY := pf.y - originf.y

	// Apply rotation
	rotatedX := translatedX*math.Cos(radians) - translatedY*math.Sin(radians)
	rotatedY := translatedX*math.Sin(radians) + translatedY*math.Cos(radians)

	// Translate back
	newX := rotatedX + originf.x
	newY := rotatedY + originf.y

	return NewPoint(newX, newY)
}

// Scale returns a new Point that scales the calling Point p by a scalar value k.
// Both the x and y coordinates of p are multiplied by the scalar, producing a Point
// that is scaled proportionally in both dimensions.
//
// Parameters:
//   - k: The scalar value by which to multiply the x and y coordinates of Point p.
//
// Returns:
//   - Point[T]: A new Point where each coordinate is the result of scaling by k.
//
// Example Usage:
//
//	p := NewPoint(3, 4)
//	scaledPoint := p.Scale(2) // scaledPoint is a Point with coordinates (6, 8)
func (p Point[T]) Scale(k T) Point[T] {
	return p.ScaleFrom(NewPoint[T](0, 0), k)
}

// ScaleFrom scales the point by a factor k relative to a reference point ref.
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
//	scaled := p.ScaleFrom(ref, 2) // scaled is now (5, 7), relative to ref.
func (p Point[float64]) ScaleFrom(ref Point[float64], k float64) Point[float64] {
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
//   - points []Point[T]: A slice of points representing the vertices of a polygon. The points are assumed
//     to form a closed loop or define a valid polygon.
//
// Behavior:
//   - Calculates the signed area of the polygon using SignedArea2X.
//   - If the signed area is positive (counterclockwise orientation), reverses the order of the points.
//   - If the signed area is negative or zero (already clockwise or degenerate polygon), does nothing.
//
// Notes:
//   - This function modifies the input slice of points in place.
//   - A zero area polygon is considered already clockwise and is left unchanged.
//
// Example:
//
//	points := []Point[float64]{
//		NewPoint(0, 0),
//		NewPoint(4, 0),
//		NewPoint(4, 4),
//	}
//	EnsureClockwise(points)
//	// points is now ordered in a clockwise direction.
//
// Dependencies:
//   - This function uses SignedArea2X to compute the signed area of the polygon.
//   - The slices.Reverse function is used to reverse the order of the points.
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
//   - points []Point[T]: A slice of points representing the vertices of a polygon. The points are assumed
//     to form a closed loop or define a valid polygon.
//
// Behavior:
//   - Calculates the signed area of the polygon using SignedArea2X.
//   - If the signed area is negative (clockwise orientation), reverses the order of the points.
//   - If the signed area is positive or zero (already counterclockwise or degenerate polygon), does nothing.
//
// Notes:
//   - This function modifies the input slice of points in place.
//   - A zero area polygon is considered already counterclockwise and is left unchanged.
//
// Example:
//
//		points := []Point[float64]{
//			NewPoint(0, 0),
//			NewPoint(4, 4),
//	     	NewPoint(4, 0),
//		}
//		EnsureCounterClockwise(points)
//		// points is now ordered in a counterclockwise direction.
//
// Dependencies:
//   - This function uses SignedArea2X to compute the signed area of the polygon.
//   - The slices.Reverse function is used to reverse the order of the points.
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

// RelativeAngle calculates the angle in radians between two points, A and B, relative to an optional origin Point O.
// This angle is measured from the origin to the line segments connecting A and B.
// If no origin is provided, the function defaults to using the point (0, 0) as the origin.
//
// Parameters:
//   - A: The first point forming one side of the angle.
//   - B: The second point forming the other side of the angle.
//   - O: An optional origin Point. If not provided, the origin defaults to (0, 0).
//
// Returns:
//   - radians (float64): The angle between points A and B relative to the origin, in radians.
//
// Example Usage:
//
//	A := NewPoint(1, 0)
//	B := NewPoint(0, 1)
//	radians := RelativeAngle(A, B) // radians is π/2 (90 degrees) for a right angle
func RelativeAngle[T SignedNumber](A, B Point[T], O ...Point[T]) (radians float64) {
	return math.Acos(RelativeCosineOfAngle(A, B, O...))
}

// RelativeCosineOfAngle calculates the cosine of the angle between two points, A and B, relative to an optional origin Point O.
// This function returns the cosine of the angle directly, avoiding the costly acos calculation, which can improve performance
// in applications where only the cosine value is needed.
//
// If no origin is provided, the function defaults to using the point (0, 0) as the origin.
//
// Parameters:
//   - A: The first point forming one side of the angle.
//   - B: The second point forming the other side of the angle.
//   - O: An optional origin Point. If not provided, the origin defaults to (0, 0).
//
// Returns:
//   - float64: The cosine of the angle between points A and B relative to the origin.
//
// Example Usage:
//
//	A := NewPoint(1, 0)
//	B := NewPoint(0, 1)
//	cosine := RelativeCosineOfAngle(A, B) // cosine is 0 for a 90-degree angle
//
// Note:
//   - This function does not currently handle division by zero errors. If either vector OA or OB has zero length,
//     a division by zero could occur. Consider adding validation if needed in such cases.
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
// The function calculates the signed area of the polygon using the Shoelace formula,
// adapted to sum the areas of triangles formed by consecutive points. The result is
// twice the actual signed area, which avoids introducing fractional values and simplifies
// calculations with integer-based coordinate types.
//
// A positive signed area indicates that the points are ordered counterclockwise,
// while a negative signed area indicates clockwise order. This function is commonly
// used to determine the orientation of a polygon or to compute its area efficiently.
//
// Parameters:
//   - points: A slice of Point values representing the vertices of the polygon in order.
//     The polygon is assumed to be closed, meaning the first point connects
//     to the last point.
//
// Returns:
//   - The signed area multiplied by 2 (hence "2X").
//     Returns 0 if the number of points is fewer than 3, as a valid polygon cannot be formed.
//
// Example Usage:
//
//	points := []Point[int]{
//		NewPoint(0, 0),
//		NewPoint(4, 0),
//		NewPoint(4, 3),
//	}
//	signedArea := SignedArea2X(points)
//	fmt.Println(signedArea) // Output: 12 (twice the signed area)
//
// Notes:
//   - The function assumes the input points form a simple polygon (no self-intersections).
//   - If the points are not in order, the result may not represent the true orientation
//     or area of the intended polygon.
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
