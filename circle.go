package geom2d

import (
	"fmt"
	"math"
)

// Circle represents a circle in 2D space with a center point and a radius.
//
// The Circle type provides methods for calculating its circumference and area,
// determining if a point lies within the circle, checking if a line segment
// intersects the circle, and checking the relationship between the circle and other geometric shapes.
type Circle[T SignedNumber] struct {
	center Point[T] // The center point of the circle
	radius T        // The radius of the circle
}

// NewCircle creates a new [Circle] with the specified center point and radius.
//
// Parameters:
//   - center ([Point][T]): The center [Point] of the [Circle].
//   - radius (T): The radius of the circle, of generic type T, which must satisfy the [SignedNumber] constraint.
//
// Returns:
//   - Circle[T]: A new Circle with the specified center and radius.
func NewCircle[T SignedNumber](center Point[T], radius T) Circle[T] {
	return Circle[T]{
		center: center,
		radius: radius,
	}
}

// Area calculates the area of the circle.
//
// Returns:
//   - float64: The area of the circle, computed as π * radius^2.
func (c Circle[T]) Area() float64 {
	return math.Pi * float64(c.radius) * float64(c.radius)
}

// AsFloat32 converts the Circle's center coordinates and radius to the float32 type, returning a new Circle[float32].
// This method is useful for cases where higher precision or floating-point arithmetic is required.
//
// Returns:
//   - Circle[float32]: A new Circle with the center coordinates and radius converted to float64.
func (c Circle[T]) AsFloat32() Circle[float32] {
	return Circle[float32]{
		center: c.center.AsFloat32(),
		radius: float32(c.radius),
	}
}

// AsFloat64 converts the Circle's center coordinates and radius to the float64 type, returning a new Circle[float64].
// This method is useful for cases where higher precision or floating-point arithmetic is required.
//
// Returns:
//   - Circle[float64]: A new Circle with the center coordinates and radius converted to float64.
func (c Circle[T]) AsFloat64() Circle[float64] {
	return Circle[float64]{
		center: c.center.AsFloat64(),
		radius: float64(c.radius),
	}
}

// AsInt converts the Circle's center coordinates and radius to the int type by truncating any decimal values.
// This method is useful when integer values are needed, such as for pixel-based or grid-based calculations.
//
// Returns:
//   - Circle[int]: A new Circle with center coordinates and radius converted to int by truncating any decimal portion.
func (c Circle[T]) AsInt() Circle[int] {
	return Circle[int]{
		center: c.center.AsInt(),
		radius: int(c.radius),
	}
}

// AsIntRounded converts the Circle's center coordinates and radius to the int type by rounding to the nearest integer.
// This method is useful when integer values are needed and rounding provides a more accurate representation
// compared to truncation.
//
// Returns:
//   - Circle[int]: A new Circle with center coordinates and radius converted to int by rounding to the nearest integer.
func (c Circle[T]) AsIntRounded() Circle[int] {
	return Circle[int]{
		center: c.center.AsIntRounded(),
		radius: int(math.Round(float64(c.radius))),
	}
}

// BoundingBox calculates the axis-aligned bounding box (AABB) of the circle.
//
// The bounding box is the smallest rectangle, aligned with the coordinate axes, that completely encloses the circle.
// This is useful for collision detection, spatial partitioning, and other geometric operations.
//
// Returns:
//   - [Rectangle][T]: The axis-aligned bounding box that encloses the circle.
//
// Notes:
//   - The bounding box is a rectangle defined by the four corner points derived from the circle's center and radius.
func (c Circle[T]) BoundingBox() Rectangle[T] {
	return NewRectangle[T]([]Point[T]{
		NewPoint(c.center.x-c.radius, c.center.y-c.radius),
		NewPoint(c.center.x+c.radius, c.center.y-c.radius),
		NewPoint(c.center.x+c.radius, c.center.y+c.radius),
		NewPoint(c.center.x-c.radius, c.center.y+c.radius),
	})
}

// Center returns the center [Point] of the Circle.
//
// Returns:
//   - Point[T]: The center [Point] of the Circle.
func (c Circle[T]) Center() Point[T] {
	return c.center
}

// Circumference calculates the circumference of the circle.
//
// Returns:
//   - float64: The circumference of the circle, computed as 2 * π * radius.
func (c Circle[T]) Circumference() float64 {
	return 2 * math.Pi * float64(c.radius)
}

// Eq determines whether the calling Circle (c) is equal to another Circle (other), either exactly (default)
// or approximately using an epsilon threshold.
//
// Parameters
//   - other: The Circle to compare with the calling Circle.
//   - opts: A variadic slice of [Option] functions to customize the equality check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing the center coordinates
//     and radii of the circles, allowing for robust handling of floating-point precision errors.
//
// Behavior
//   - The function checks whether the center points of the two circles are equal (using the [Point.Eq] method
//     of the [Point] type) and whether their radii are equal.
//   - If [WithEpsilon] is provided, the comparison allows for small differences in the radius values
//     and center coordinates within the epsilon threshold.
//
// Returns
//   - bool: true if the center coordinates and radius of the two circles are equal (or approximately equal
//     within epsilon); otherwise, false.
//
// Notes:
//   - Approximate equality is particularly useful when comparing circles with floating-point
//     coordinates or radii, where small precision errors might otherwise cause inequality.
func (c Circle[T]) Eq(other Circle[T], opts ...Option) bool {
	// Apply options with defaults
	options := applyOptions(geomOptions{epsilon: 0}, opts...)

	// Check equality for the center points
	centersEqual := c.center.Eq(other.center, opts...)

	// Check equality for the radii with epsilon adjustment
	radiiEqual := c.radius == other.radius
	if options.epsilon > 0 {
		radiiEqual = math.Abs(float64(c.radius)-float64(other.radius)) < options.epsilon
	}

	return centersEqual && radiiEqual
}

// Radius returns the radius of the Circle.
//
// Returns:
//   - T: The radius of the Circle.
func (c Circle[T]) Radius() T {
	return c.radius
}

// RelationshipToCircle determines the spatial relationship between two circles.
//
// This function evaluates the relationship between the current circle and another
// circle by comparing their center points and radii. The possible relationships include:
//   - [RelationshipEqual]: The circles are identical.
//   - [RelationshipContainedBy]: The current circle is completely contained within the other circle.
//   - [RelationshipContains]: The current circle completely contains the other circle.
//   - [RelationshipIntersection]: The circles overlap, including tangency.
//   - [RelationshipDisjoint]: The circles do not overlap.
//
// Parameters:
//   - other (Circle[T]): The circle to compare against the current circle.
//   - opts: A variadic slice of [Option] functions to customize the equality check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing the center coordinates
//     and radii of the circles, allowing for robust handling of floating-point precision errors.
//
// Returns:
//   - [Relationship]: A constant representing the relationship between the circles.
//
// Behavior:
//   - The function first checks for equality by comparing center points and radii.
//   - It then checks for containment by comparing the distance between centers and radii.
//   - Intersection is detected if the distance between centers is less than or equal to the sum of the radii.
//   - If no other relationship is found, the circles are considered disjoint.
func (c Circle[T]) RelationshipToCircle(other Circle[T], opts ...Option) Relationship {
	distanceBetweenCenters := c.center.DistanceToPoint(other.center, opts...)
	cFloat := c.AsFloat64()
	otherFloat := other.AsFloat64()

	// check for equality
	if c.Eq(other) {
		return RelationshipEqual
	}

	// check for c contained by other
	if distanceBetweenCenters+cFloat.radius < otherFloat.radius {
		return RelationshipContainedBy
	}

	// check for c contains other
	if distanceBetweenCenters+otherFloat.radius < cFloat.radius {
		return RelationshipContains
	}

	// check for intersection
	if distanceBetweenCenters <= cFloat.radius+otherFloat.radius {
		return RelationshipIntersection
	}

	return RelationshipDisjoint

}

// RelationshipToLineSegment determines the spatial relationship between the current Circle and a
// given [LineSegment].
//
// This function evaluates the relationship between the circle and the line segment,
// which can be one of the following:
//   - [RelationshipDisjoint]: The line segment lies entirely outside the circle.
//   - [RelationshipIntersection]: The line segment intersects the circle's boundary.
//   - [RelationshipContains]: The line segment is fully contained within the circle.
//
// This method internally calls [LineSegment.RelationshipToCircle], flipping the containment
// direction to align with the perspective of the circle.
//
// Parameters:
//   - l ([LineSegment][T]): The line segment to compare against the circle.
//   - opts: A variadic slice of [Option] functions to customize the equality check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing the center coordinates
//     and radii of the circles, allowing for robust handling of floating-point precision errors.
//
// Returns:
//   - [Relationship]: A constant representing the relationship between the circle and the line segment.
func (c Circle[T]) RelationshipToLineSegment(l LineSegment[T], opts ...Option) Relationship {
	return l.RelationshipToCircle(c, opts...).flipContainment()
}

// RelationshipToPoint determines the spatial relationship between the current Circle and a given [Point].
//
// This function evaluates the relationship between the circle and the point,
// which can be one of the following:
//   - [RelationshipDisjoint]: The point lies outside the circle.
//   - [RelationshipIntersection]: The point lies exactly on the circle's boundary.
//   - [RelationshipContains]: The point lies inside the circle.
//
// This method internally calls [Point.RelationshipToCircle], flipping the containment
// direction to align with the perspective of the circle.
//
// Parameters:
//   - p ([Point][T]): The point to compare against the circle.
//   - opts: A variadic slice of [Option] functions to customize the equality check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing the center coordinates
//     and radii of the circles, allowing for robust handling of floating-point precision errors.
//
// Returns:
//   - [Relationship]: A constant representing the relationship between the circle and the point.
func (c Circle[T]) RelationshipToPoint(p Point[T], opts ...Option) Relationship {
	return p.RelationshipToCircle(c, opts...).flipContainment()
}

// RelationshipToPolyTree determines the spatial relationships between the Circle and the polygons in the given [PolyTree].
//
// This method evaluates whether the circle intersects, contains, or is contained by each polygon in the [PolyTree].
// It uses a doubled representation of the Circle to align with the doubled points in the [PolyTree] for robust computations.
//
// Parameters:
//   - pt (*[PolyTree][T]): The [PolyTree] whose polygons will be compared to the Circle.
//   - opts: A variadic slice of [Option] functions to customize the equality check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing the center coordinates
//     and radii of the circles, allowing for robust handling of floating-point precision errors.
//
// Behavior:
//   - For each polygon in the [PolyTree], the function iterates through its edges and checks the relationship to the Circle.
//   - Intersection is checked first; if any edge of the polygon intersects the circle, the relationship is marked as [RelationshipIntersection].
//   - If all edges of the polygon lie within the circle's radius, the relationship is marked as [RelationshipContains].
//   - If the circle's center lies within the polygon and its minimum distance to any edge is greater than its radius, the relationship is marked as [RelationshipContainedBy].
//   - If none of the above conditions are satisfied, the relationship is marked as [RelationshipDisjoint].
//
// Returns:
//   - map[*[PolyTree][T]][Relationship]: A map where each key is a pointer to a polygon in the [PolyTree], and the value is the relationship between the Circle and that polygon.
func (c Circle[T]) RelationshipToPolyTree(pt *PolyTree[T], opts ...Option) map[*PolyTree[T]]Relationship {
	output := make(map[*PolyTree[T]]Relationship, pt.Len())
	cDoubled := NewCircle(NewPoint(c.center.x*2, c.center.y*2), c.radius*2)
	cFloatDoubled := cDoubled.AsFloat64()

RelationshipToPolyTreeIterPolys:
	for poly := range pt.Nodes {
		minDistCircleCenterToEdge := math.MaxFloat64
		allEdgesWithinCircle := true

		for edge := range poly.contour.iterEdges {
			rel := cDoubled.RelationshipToLineSegment(edge, opts...)

			// Check for intersection
			if rel == RelationshipIntersection {
				output[poly] = RelationshipIntersection
				continue RelationshipToPolyTreeIterPolys
			}

			// Check if all edges are within the circle's radius
			distanceToEdge := cDoubled.center.DistanceToLineSegment(edge, opts...)
			minDistCircleCenterToEdge = min(minDistCircleCenterToEdge, distanceToEdge)
			if distanceToEdge > cFloatDoubled.radius {
				allEdgesWithinCircle = false
			}
		}

		// Check for containment: circle fully contains the polygon
		if allEdgesWithinCircle {
			output[poly] = RelationshipContains
			continue RelationshipToPolyTreeIterPolys
		}

		// Check for containment: polygon fully contains the circle
		if poly.contour.isPointInside(cDoubled.center) && minDistCircleCenterToEdge > cFloatDoubled.radius {
			output[poly] = RelationshipContainedBy
			continue RelationshipToPolyTreeIterPolys
		}

		// Default: no relationship found
		output[poly] = RelationshipDisjoint
	}

	return output
}

// RelationshipToRectangle determines the spatial relationship between the circle and the rectangle.
//
// This function evaluates whether the circle is:
//   - Disjoint from the rectangle (no overlap or touching),
//   - Intersects with the rectangle (crosses its boundary),
//   - Fully contains the rectangle (encloses it entirely),
//   - Fully contained by the rectangle (is completely inside the rectangle).
//
// Parameters:
//   - r ([Rectangle][T]): The rectangle to compare with the circle.
//   - opts: A variadic slice of [Option] functions to customize the behavior of the relationship check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for floating-point precision.
//
// Behavior:
//   - The function checks each edge of the rectangle for potential intersections with the circle.
//   - If all edges of the rectangle are fully contained within the circle, it returns [RelationshipContains].
//   - If the rectangle fully contains the circle (the circle’s center is inside the rectangle, and the circle
//     does not extend beyond any edge), it returns [RelationshipContainedBy].
//   - If none of these conditions are met, it determines whether the circle and rectangle are disjoint.
//
// Returns:
//
// [Relationship]: One of the following constants:
//   - [RelationshipDisjoint]: The circle and rectangle are entirely separate.
//   - [RelationshipIntersection]: The circle intersects with one or more edges of the rectangle.
//   - [RelationshipContains]: The circle completely encloses the rectangle.
//   - [RelationshipContainedBy]: The circle is fully contained within the rectangle.
func (c Circle[T]) RelationshipToRectangle(r Rectangle[T], opts ...Option) Relationship {
	cContainsR := true
	cFloat := c.AsFloat64()
	minDistCircleCenterToEdge := math.MaxFloat64
	for _, edge := range r.Edges() {
		rel := edge.RelationshipToCircle(c, opts...)

		// check for intersection
		if rel == RelationshipIntersection {
			return RelationshipIntersection
		}

		// check for containment
		if rel != RelationshipContainedBy {
			cContainsR = false
		}

		edgeFloat := edge.AsFloat64()
		minDistCircleCenterToEdge = min(minDistCircleCenterToEdge, cFloat.center.DistanceToLineSegment(edgeFloat, opts...))
	}

	// check c contain r
	if cContainsR {
		return RelationshipContains
	}

	// check r contains c
	if r.ContainsPoint(c.center) && minDistCircleCenterToEdge > cFloat.radius {
		return RelationshipContainedBy
	}

	return RelationshipDisjoint
}

// Rotate rotates the Circle's center around a specified pivot [Point] by a given angle in radians
// counterclockwise, while keeping the radius unchanged. Optionally, an epsilon threshold can be applied
// to adjust the precision of the resulting coordinates.
//
// Parameters:
//   - pivot ([Point][T]): The [Point] around which to rotate the circle's center.
//   - radians: The rotation angle in radians (counterclockwise).
//   - opts: A variadic slice of [Option] functions to customize the behavior of the rotation.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for snapping the resulting center coordinates
//     to cleaner values, improving robustness in floating-point calculations.
//
// Returns:
//   - [Circle][float64]: A new [Circle] with the center rotated around the pivot [Point] by the specified angle,
//     and with the radius unchanged.
//
// Behavior:
//   - The function rotates the circle's center point around the given pivot by the specified angle using
//     the [Point.Rotate] method.
//   - The rotation is performed in a counterclockwise direction relative to the pivot point.
//   - The radius remains unchanged in the resulting [Circle].
//   - If [WithEpsilon] is provided, epsilon adjustments are applied to the rotated coordinates to handle
//     floating-point precision errors.
//
// Notes:
//   - Epsilon adjustment is particularly useful when the rotation involves floating-point calculations
//     that could result in minor inaccuracies.
//   - The returned [Circle] always has a center with float64 coordinates, ensuring precision regardless
//     of the coordinate type of the original [Circle].
func (c Circle[T]) Rotate(pivot Point[T], radians float64, opts ...Option) Circle[float64] {
	return NewCircle[float64](
		c.center.Rotate(pivot, radians, opts...),
		float64(c.radius),
	)
}

// Scale scales the radius of the circle by a scalar factor.
//
// Parameters:
//   - factor (T): The factor by which to scale the radius.
//
// Returns:
//   - Circle[T]: A new circle with the radius scaled by the specified factor.
func (c Circle[T]) Scale(factor T) Circle[T] {
	return Circle[T]{center: c.center, radius: c.radius * factor}
}

// String returns a string representation of the Circle, including its center coordinates and radius.
// This is useful for debugging and logging.
//
// Returns:
//   - string: A string representation of the Circle in the format "Circle[center=(x, y), radius=r]".
func (c Circle[T]) String() string {
	return fmt.Sprintf("Circle[center=(%v, %v), radius=%v]", c.center.x, c.center.y, c.radius)
}

// Translate moves the circle by a specified vector (given as a [Point]).
//
// This method shifts the circle's center by the given vector v, effectively
// translating the circle's position in the 2D plane. The radius of the circle
// remains unchanged.
//
// Parameters:
//   - v ([Point][T]): The vector by which to translate the circle's center.
//
// Returns:
//   - Circle[T]: A new Circle translated by the specified vector.
func (c Circle[T]) Translate(v Point[T]) Circle[T] {
	return Circle[T]{center: c.center.Translate(v), radius: c.radius}
}

// todo: doc comments, unit tests & example function
func (c Circle[int]) Bresenham(yield func(Point[int]) bool) {
	var xc, yc, r, x, y, p int

	xc = c.center.x
	yc = c.center.y
	r = c.radius

	// Starting at the top of the circle
	x = 0
	y = r
	p = 1 - r // Initial decision parameter

	// Yield the initial points for all octants
	for _, point := range reflectAcrossCircleOctants(xc, yc, x, y) {
		if !yield(point) {
			return
		}
	}

	// Loop until x meets y
	for x < y {
		x++
		if p < 0 {
			// Midpoint is inside the circle
			p += 2*x + 1
		} else {
			// Midpoint is outside or on the circle
			y--
			p += 2*(x-y) + 1
		}

		// Yield the points for the current x, y
		for _, point := range reflectAcrossCircleOctants(xc, yc, x, y) {
			if !yield(point) {
				return
			}
		}
	}
}

// Reflect across all octants
// todo: doc comments, unit tests
func reflectAcrossCircleOctants[T SignedNumber](xc, yc, x, y T) []Point[T] {
	return []Point[T]{
		NewPoint[T](xc+x, yc+y), // Octant 1
		NewPoint[T](xc-x, yc+y), // Octant 2
		NewPoint[T](xc+x, yc-y), // Octant 8
		NewPoint[T](xc-x, yc-y), // Octant 7
		NewPoint[T](xc+y, yc+x), // Octant 3
		NewPoint[T](xc-y, yc+x), // Octant 4
		NewPoint[T](xc+y, yc-x), // Octant 6
		NewPoint[T](xc-y, yc-x), // Octant 5
	}
}
