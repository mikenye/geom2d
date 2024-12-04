// This file defines the Circle type, which represents a circle specified by a center point and a radius.
// The Circle type includes methods for geometric computations such as circumference, area, and checking
// relationships like containment and intersection with other shapes.
//
// The Circle type integrates seamlessly with the geom2d package, allowing it to be used in combination
// with other geometric primitives, including points, line segments, and polygons.

package geom2d

import (
	"fmt"
	"math"
)

// Circle represents a circle in 2D space with a center point and a radius.
//
// The Circle type provides methods for calculating its circumference and area,
// determining if a point lies within the circle, and checking if a line segment
// intersects the circle.
type Circle[T SignedNumber] struct {
	center Point[T] // The center point of the circle
	radius T        // The radius of the circle
}

// Add translates the circle's center by adding a vector.
//
// Parameters:
//   - v: Point[T] - The vector to add to the circle's center.
//
// Returns:
//   - Circle[T]: A new circle with the center moved by the specified vector.
func (c Circle[T]) Add(v Point[T]) Circle[T] {
	return Circle[T]{center: c.center.Add(v), radius: c.radius}
}

// Area calculates the area of the circle.
//
// Returns:
//   - float64: The area of the circle, computed as π * radius^2.
func (c Circle[T]) Area() float64 {
	return math.Pi * float64(c.radius) * float64(c.radius)
}

// AsFloat converts the Circle's center coordinates and radius to the float64 type, returning a new Circle[float64].
// This method is useful for cases where higher precision or floating-point arithmetic is required.
//
// Returns:
//   - Circle[float64]: A new Circle with the center coordinates and radius converted to float64.
//
// Example Usage:
//
//	c := NewCircle(NewPoint(3, 4), 5)
//	floatCircle := c.AsFloat() // floatCircle is a Circle[float64] with center (3.0, 4.0) and radius 5.0
func (c Circle[T]) AsFloat() Circle[float64] {
	return Circle[float64]{
		center: c.center.AsFloat(),
		radius: float64(c.radius),
	}
}

// AsInt converts the Circle's center coordinates and radius to the int type by truncating any decimal values.
// This method is useful when integer values are needed, such as for pixel-based or grid-based calculations.
//
// Returns:
//   - Circle[int]: A new Circle with center coordinates and radius converted to int by truncating any decimal portion.
//
// Example Usage:
//
//	c := NewCircle(NewPoint(3.7, 4.9), 5.6)
//	intCircle := c.AsInt() // intCircle is a Circle[int] with center (3, 4) and radius 5
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
//
// Example Usage:
//
//	c := NewCircle(NewPoint(3.7, 4.2), 5.6)
//	roundedCircle := c.AsIntRounded() // roundedCircle is a Circle[int] with center (4, 4) and radius 6
func (c Circle[T]) AsIntRounded() Circle[int] {
	return Circle[int]{
		center: c.center.AsIntRounded(),
		radius: int(math.Round(float64(c.radius))),
	}
}

// Center returns the center point of the Circle.
//
// Returns:
//   - Point[T]: The center point of the Circle.
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

// Div scales down the radius of the circle by dividing by a scalar factor.
// Optionally, an epsilon threshold can be applied to adjust the precision of the resulting radius.
//
// Parameters:
//   - divisor: T - The factor by which to divide the radius.
//   - opts: A variadic slice of Option functions to customize the behavior of the division.
//     WithEpsilon(epsilon float64): Specifies a tolerance for snapping the resulting radius
//     to cleaner values, improving robustness in floating-point calculations.
//
// Returns:
//   - Circle[float64]: A new circle with the radius scaled down by the specified factor.
func (c Circle[T]) Div(divisor T, opts ...Option) Circle[float64] {
	// Apply options with defaults
	options := applyOptions(geomOptions{epsilon: 0}, opts...)

	// Scale the radius
	scaledRadius := float64(c.radius) / float64(divisor)

	// Apply epsilon if specified
	if options.epsilon > 0 {
		scaledRadius = applyEpsilon(scaledRadius, options.epsilon)
	}

	return Circle[float64]{center: c.center.AsFloat(), radius: scaledRadius}
}

// Eq determines whether the calling Circle is equal to another Circle, either exactly (default)
// or approximately using an epsilon threshold.
//
// Parameters:
//   - c2: The Circle to compare with the calling Circle.
//   - opts: A variadic slice of Option functions to customize the equality check.
//   - WithEpsilon(epsilon float64): Specifies a tolerance for comparing the center coordinates
//     and radii of the circles, allowing for robust handling of floating-point precision errors.
//
// Behavior:
//   - The function checks whether the center points of the two circles are equal (using the `Eq` method
//     of the `Point` type) and whether their radii are equal.
//   - If `WithEpsilon` is provided, the comparison allows for small differences in the radius values
//     and center coordinates within the epsilon threshold.
//
// Returns:
//   - bool: `true` if the center coordinates and radius of the two circles are equal (or approximately equal
//     within epsilon); otherwise, `false`.
//
// Example Usage:
//
//	c1 := NewCircle(NewPoint(3, 4), 5)
//	c2 := NewCircle(NewPoint(3, 4), 5)
//
//	// Default behavior (exact equality)
//	isEqual := c1.Eq(c2) // isEqual is true
//
//	// Approximate equality with epsilon
//	c3 := NewCircle(NewPoint(3.0001, 4.0001), 5.0001)
//	isApproximatelyEqual := c1.Eq(c3, WithEpsilon(1e-3)) // isApproximatelyEqual is true
//
// Notes:
//   - Approximate equality is particularly useful when comparing circles with floating-point
//     coordinates or radii, where small precision errors might otherwise cause inequality.
func (c Circle[T]) Eq(c2 Circle[T], opts ...Option) bool {
	// Apply options with defaults
	options := applyOptions(geomOptions{epsilon: 0}, opts...)

	// Check equality for the center points
	centersEqual := c.center.Eq(c2.center, opts...)

	// Check equality for the radii with epsilon adjustment
	radiiEqual := c.radius == c2.radius
	if options.epsilon > 0 {
		radiiEqual = math.Abs(float64(c.radius)-float64(c2.radius)) < options.epsilon
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

// RelationshipToLineSegment determines the spatial relationship of a line segment
// to the circle. It returns one of several possible relationships, such as whether
// the segment is inside, outside, tangent to, or intersects the circle.
//
// Parameters:
//   - AB: The line segment to analyze.
//   - opts: A variadic slice of Option functions to customize the behavior of the relationship check.
//   - WithEpsilon(epsilon float64): Specifies a tolerance for comparing distances to the circle's radius,
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
//	relationship := c.RelationshipToLineSegment(segment)
//
//	// With epsilon adjustment
//	relationshipWithEpsilon := c.RelationshipToLineSegment(segment, WithEpsilon(1e-4))
//
// Notes:
//   - Epsilon adjustment is particularly useful for floating-point coordinates, where small precision
//     errors might otherwise cause incorrect classifications.
func (c Circle[T]) RelationshipToLineSegment(AB LineSegment[T], opts ...Option) CircleLineSegmentRelationship {
	// Apply options with defaults
	options := applyOptions(geomOptions{epsilon: 0}, opts...)

	// Calculate distances from the circle's center to the line segment's endpoints
	distStart := c.center.DistanceToPoint(AB.start, opts...)
	distEnd := c.center.DistanceToPoint(AB.end, opts...)

	// Check if both endpoints are within the circle's radius
	if distStart < float64(c.radius) && distEnd < float64(c.radius) {
		return CLRInside
	}

	// Check if both endpoints are exactly on the boundary
	if math.Abs(distStart-float64(c.radius)) < options.epsilon && math.Abs(distEnd-float64(c.radius)) < options.epsilon {
		return CLRBothEndsOnCircumference
	}

	// Check if one endpoint is on the boundary
	if math.Abs(distStart-float64(c.radius)) < options.epsilon || math.Abs(distEnd-float64(c.radius)) < options.epsilon {
		if distStart < float64(c.radius) || distEnd < float64(c.radius) {
			// One endpoint is on the circumference, and the other is inside
			return CLROneEndOnCircumferenceInside
		} else {
			// One endpoint is on the circumference, and the other is outside
			return CLROneEndOnCircumferenceOutside
		}
	}

	// Calculate the closest point on the segment to the circle's center
	closestPoint := c.center.ProjectOntoLineSegment(AB)
	closestDistance := closestPoint.DistanceToPoint(c.center.AsFloat(), opts...)

	// Check if the closest point is on the circle's boundary or inside the circle
	if closestDistance <= float64(c.radius)+options.epsilon {

		// True tangent check:
		// Confirm perpendicularity (right angle)
		// Calculate direction vector for line segment and radius vector to closest point
		segmentDirection := AB.end.Sub(AB.start)             // Vector along line segment
		radiusVector := closestPoint.Sub(c.center.AsFloat()) // Vector from center to closest point

		// Dot product should be zero for perpendicular vectors
		isPerpendicular := math.Abs(segmentDirection.AsFloat().DotProduct(radiusVector)) < options.epsilon

		if math.Abs(closestDistance-float64(c.radius)) < options.epsilon && isPerpendicular {
			return CLRTangent
		}

		// Otherwise, it's intersecting
		return CLRIntersecting
	}

	// If none of the conditions are met, the segment is outside the circle
	return CLROutside
}

// RelationshipToPoint determines the relationship of a given point to the circle.
// It returns whether the point is Outside, OnCircumference, or Inside the circle.
//
// Parameters:
//   - p: The point to check, of type Point[T].
//   - opts: A variadic slice of Option functions to customize the behavior of the relationship check.
//   - WithEpsilon(epsilon float64): Specifies a tolerance for comparing the distance of the point
//     to the circle's radius, improving robustness against floating-point precision errors.
//
// Returns:
//   - PointCircleRelationship: The relationship of the point to the circle, indicating whether it
//     is outside, on the circumference, or inside the circle.
//
// Behavior:
//   - The function calculates the Euclidean distance between the point `p` and the circle's center.
//   - It compares this distance to the circle's radius:
//   - `PCRInside`: The point lies inside the circle (distance < radius).
//   - `PCROnCircumference`: The point lies on the circumference of the circle (distance ≈ radius).
//   - `PCROutside`: The point lies outside the circle (distance > radius).
//   - If `WithEpsilon` is provided, epsilon adjustments are applied to the comparison with the radius.
//
// Example Usage:
//
//	c := NewCircle(NewPoint(0.0, 0.0), 5.0)
//	point := NewPoint(3.0, 4.0)
//
//	// Default behavior (no epsilon adjustment)
//	relationship := c.RelationshipToPoint(point) // Returns PCRInside since (3, 4) is within radius 5
//
//	// With epsilon adjustment
//	point2 := NewPoint(5.000001, 0.0)
//	relationshipWithEpsilon := c.RelationshipToPoint(point2, WithEpsilon(1e-4))
//	// Returns PCROnCircumference because (5.000001, 0.0) is close enough to the circumference.
//
// Notes:
//   - Epsilon adjustment is particularly useful when working with floating-point coordinates, where small
//     precision errors might otherwise cause incorrect classifications.
func (c Circle[T]) RelationshipToPoint(p Point[T], opts ...Option) PointCircleRelationship {
	// Apply options with defaults
	options := applyOptions(geomOptions{epsilon: 0}, opts...)

	distance := c.center.DistanceToPoint(p, opts...)
	switch {
	case distance < float64(c.radius)-options.epsilon:
		return PCRInside
	case math.Abs(distance-float64(c.radius)) < options.epsilon:
		return PCROnCircumference
	default:
		return PCROutside
	}
}

// Rotate rotates the Circle's center around a specified pivot point by a given angle in radians,
// while keeping the radius unchanged. Optionally, an epsilon threshold can be applied to adjust
// the precision of the resulting coordinates.
//
// Parameters:
//   - pivot: The point around which to rotate the circle's center.
//   - radians: The rotation angle in radians.
//   - opts: A variadic slice of Option functions to customize the behavior of the rotation.
//     WithEpsilon(epsilon float64): Specifies a tolerance for snapping the resulting center coordinates
//     to cleaner values, improving robustness in floating-point calculations.
//
// Returns:
//   - Circle[float64]: A new Circle with the center rotated around the pivot point by the specified angle,
//     and with the radius unchanged.
//
// Behavior:
//   - The function rotates the circle's center point around the given pivot by the specified angle using
//     the `Point.Rotate` method.
//   - The radius remains unchanged in the resulting Circle.
//   - If `WithEpsilon` is provided, epsilon adjustments are applied to the rotated coordinates to handle
//     floating-point precision errors.
//
// Example Usage:
//
//	pivot := NewPoint(1.0, 1.0)
//	circle := NewCircle(NewPoint(3.0, 3.0), 5.0)
//	angle := math.Pi / 2 // Rotate 90 degrees
//
//	// Default behavior (no epsilon adjustment)
//	rotatedCircle := circle.Rotate(pivot, angle)
//
//	// With epsilon adjustment
//	rotatedCircleWithEpsilon := circle.Rotate(pivot, angle, WithEpsilon(1e-10))
//
// Notes:
//   - Epsilon adjustment is particularly useful when the rotation involves floating-point calculations
//     that could result in minor inaccuracies.
//   - The returned Circle always has a center with `float64` coordinates, ensuring precision regardless
//     of the coordinate type of the original Circle.
func (c Circle[T]) Rotate(pivot Point[T], radians float64, opts ...Option) Circle[float64] {
	return NewCircle[float64](
		c.center.Rotate(pivot, radians, opts...),
		float64(c.radius),
	)
}

// Scale scales the radius of the circle by a scalar factor.
//
// Parameters:
//   - factor: T - The factor by which to scale the radius.
//
// Returns:
//   - Circle[float64]: A new circle with the radius scaled by the specified factor.
func (c Circle[T]) Scale(factor T) Circle[float64] {
	return Circle[float64]{center: c.center.AsFloat(), radius: float64(c.radius) * float64(factor)}
}

// String returns a string representation of the Circle, including its center coordinates and radius.
// This is useful for debugging and logging.
//
// Returns:
//   - string: A string representation of the Circle in the format "Circle[center=(x, y), radius=r]".
//
// Example Usage:
//
//	c := NewCircle(NewPoint(3, 4), 5)
//	fmt.Println(c.String()) // Output: "Circle[center=(3, 4), radius=5]"
func (c Circle[T]) String() string {
	return fmt.Sprintf("Circle[center=(%v, %v), radius=%v]", c.center.x, c.center.y, c.radius)
}

// Sub translates the circle's center by subtracting a vector.
//
// Parameters:
//   - v: Point[T] - The vector to subtract from the circle's center.
//
// Returns:
//   - Circle[T]: A new circle with the center moved by the specified vector.
func (c Circle[T]) Sub(v Point[T]) Circle[T] {
	return Circle[T]{center: c.center.Sub(v), radius: c.radius}
}

// NewCircle creates a new Circle with the specified center point and radius.
//
// Parameters:
//   - center: The center point of the circle, of type Point[T].
//   - radius: The radius of the circle, of generic type T, which must satisfy the SignedNumber constraint.
//
// Returns:
//   - Circle[T]: A new Circle with the specified center and radius.
//
// Example Usage:
//
//	center := NewPoint(3, 4)
//	radius := 5
//	circle := NewCircle(center, radius) // Creates a Circle with center (3, 4) and radius 5
func NewCircle[T SignedNumber](center Point[T], radius T) Circle[T] {
	return Circle[T]{
		center: center,
		radius: radius,
	}
}
