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

// AsFloat converts the Circle's center coordinates and radius to the float64 type, returning a new Circle[float64].
// This method is useful for cases where higher precision or floating-point arithmetic is required.
//
// Returns:
//   - Circle[float64]: A new Circle with the center coordinates and radius converted to float64.
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

// Eq determines whether the calling Circle is equal to another Circle, either exactly (default)
// or approximately using an epsilon threshold.
//
// Parameters
//   - c2: The Circle to compare with the calling Circle.
//   - opts: A variadic slice of Option functions to customize the equality check.
//   - WithEpsilon(epsilon float64): Specifies a tolerance for comparing the center coordinates
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

// RelationshipToCircle determines the relationship between the calling circle
// and another circle.
//
// This method calculates the distance between the centers of the two circles
// and compares it with their radii to determine their spatial relationship.
//
// Parameters:
//   - other (Circle[T]): The circle to compare with the calling circle.
//
// Returns:
//   - [RelationshipCircleCircle]: The relationship between the two circles.
func (c Circle[T]) RelationshipToCircle(other Circle[T], opts ...Option) RelationshipCircleCircle {
	options := applyOptions(geomOptions{epsilon: 0}, opts...)

	// Calculate the distance between the centers of the two circles
	centerDistance := c.center.DistanceToPoint(other.center, opts...)

	// Calculate the sum and absolute difference of the radii
	sumRadii := float64(c.radius) + float64(other.radius)
	diffRadii := math.Abs(float64(c.radius) - float64(other.radius))

	absCenterDistanceSubRadiiSub := math.Abs(centerDistance - sumRadii)
	absCenterDistanceSubDiffRadii := math.Abs(centerDistance - diffRadii)

	// Determine the relationship
	switch {
	case centerDistance == 0 && c.radius == other.radius:
		return RelationshipCircleCircleEqual // Circles are identical
	case absCenterDistanceSubRadiiSub < options.epsilon:
		return RelationshipCircleCircleExternallyTangent // Circles are externally tangent
	case absCenterDistanceSubDiffRadii < options.epsilon:
		return RelationshipCircleCircleInternallyTangent // Circles are internally tangent
	case centerDistance > sumRadii:
		return RelationshipCircleCircleMiss // Circles are disjoint
	case centerDistance < diffRadii:
		return RelationshipCircleCircleContained // One circle is fully contained within the other
	case centerDistance < sumRadii:
		return RelationshipCircleCircleIntersection // Circles overlap
	}

	// Fallback (should not happen)
	return RelationshipCircleCircleMiss
}

// RelationshipToLineSegment determines the spatial relationship of a [LineSegment]
// to the circle. It returns one of several possible relationships, such as whether
// the segment is inside, outside, tangent to, or intersects the circle.
//
// Parameters:
//   - AB ([LineSegment][T]): The [LineSegment] to analyze.
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
func (c Circle[T]) RelationshipToLineSegment(AB LineSegment[T], opts ...Option) RelationshipLineSegmentCircle {
	// Apply options with defaults
	options := applyOptions(geomOptions{epsilon: 0}, opts...)

	// Calculate distances from the circle's center to the line segment's endpoints
	distStart := c.center.DistanceToPoint(AB.start, opts...)
	distEnd := c.center.DistanceToPoint(AB.end, opts...)

	// Check if both endpoints are within the circle's radius
	if distStart < float64(c.radius) && distEnd < float64(c.radius) {
		return RelationshipLineSegmentCircleContainedByCircle
	}

	// Check if both endpoints are exactly on the boundary
	if math.Abs(distStart-float64(c.radius)) < options.epsilon && math.Abs(distEnd-float64(c.radius)) < options.epsilon {
		return RelationshipLineSegmentCircleBothEndsOnCircumference
	}

	// Check if one endpoint is on the boundary
	if math.Abs(distStart-float64(c.radius)) < options.epsilon || math.Abs(distEnd-float64(c.radius)) < options.epsilon {
		if distStart < float64(c.radius) || distEnd < float64(c.radius) {
			// One endpoint is on the circumference, and the other is inside
			return RelationshipLineSegmentCircleEndOnCircumferenceInside
		} else {
			// One endpoint is on the circumference, and the other is outside
			return RelationshipLineSegmentCircleEndOnCircumferenceOutside
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
		segmentDirection := AB.end.Translate(AB.start.Negate())             // Vector along line segment
		radiusVector := closestPoint.Translate(c.center.AsFloat().Negate()) // Vector from center to closest point

		// Dot product should be zero for perpendicular vectors
		isPerpendicular := math.Abs(segmentDirection.AsFloat().DotProduct(radiusVector)) < options.epsilon

		if math.Abs(closestDistance-float64(c.radius)) < options.epsilon && isPerpendicular {
			return RelationshipLineSegmentCircleTangentToCircle
		}

		// Otherwise, it's intersecting
		return RelationshipLineSegmentCircleIntersecting
	}

	// If none of the conditions are met, the segment is outside the circle
	return RelationshipLineSegmentCircleMiss
}

// RelationshipToPoint determines the relationship of a given [Point] to the circle.
// It returns whether the point is Outside, OnCircumference, or Inside the circle.
//
// Parameters:
//   - p ([Point][T]): The [Point] to check.
//   - opts: A variadic slice of Option functions to customize the behavior of the relationship check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing the distance of the point
//     to the circle's radius, improving robustness against floating-point precision errors.
//
// Returns:
//   - [RelationshipPointCircle]: The relationship of the [Point] to the circle, indicating whether it
//     is outside, on the circumference, or inside the circle.
//
// Behavior:
//   - The function calculates the Euclidean distance between the [Point] p and the circle's center.
//   - It compares this distance to the circle's radius:
//   - [RelationshipPointCircleContainedByCircle]: The point lies inside the circle (distance < radius).
//   - [RelationshipPointCircleOnCircumference]: The point lies on the circumference of the circle (distance ≈ radius).
//   - [RelationshipPointCircleMiss]: The point lies outside the circle (distance > radius).
//   - If [WithEpsilon] is provided, epsilon adjustments are applied to the comparison with the radius.
//
// Notes:
//   - Epsilon adjustment is particularly useful when working with floating-point coordinates, where small
//     precision errors might otherwise cause incorrect classifications.
func (c Circle[T]) RelationshipToPoint(p Point[T], opts ...Option) RelationshipPointCircle {
	// Apply options with defaults
	options := applyOptions(geomOptions{epsilon: 0}, opts...)

	distance := c.center.DistanceToPoint(p, opts...)
	switch {
	case distance < float64(c.radius)-options.epsilon:
		return RelationshipPointCircleContainedByCircle
	case math.Abs(distance-float64(c.radius)) < options.epsilon:
		return RelationshipPointCircleOnCircumference
	default:
		return RelationshipPointCircleMiss
	}
}

// RelationshipToRectangle determines the spatial relationship between the circle
// and a [Rectangle].
//
// This method evaluates whether the circle and [Rectangle] are disjoint, tangent,
// intersecting, or whether one is fully contained within the other.
//
// Parameters:
//   - rect ([Rectangle][T]): The rectangle to compare with the circle.
//
// Returns:
//   - [RelationshipRectangleCircle]: The relationship between the circle and the [Rectangle].
func (c Circle[T]) RelationshipToRectangle(rect Rectangle[T], opts ...Option) RelationshipRectangleCircle {

	// Calculate the bounding box of the circle
	circleBounds := NewRectangle(
		[]Point[T]{
			NewPoint(c.center.x-T(c.radius), c.center.y-T(c.radius)),
			NewPoint(c.center.x+T(c.radius), c.center.y-T(c.radius)),
			NewPoint(c.center.x+T(c.radius), c.center.y+T(c.radius)),
			NewPoint(c.center.x-T(c.radius), c.center.y+T(c.radius)),
		},
	)

	// Check if bounding boxes are disjoint
	if circleBounds.RelationshipToRectangle(rect, opts...) == RelationshipRectangleRectangleMiss {
		return RelationshipRectangleCircleMiss
	}

	// Check if the circle is fully contained within the rectangle
	circlePoints := []Point[T]{
		NewPoint(c.center.x-T(c.radius), c.center.y-T(c.radius)),
		NewPoint(c.center.x+T(c.radius), c.center.y+T(c.radius)),
	}
	allCirclePointsInRect := true
	for _, p := range circlePoints {
		if rect.RelationshipToPoint(p) == RelationshipPointRectangleMiss {
			allCirclePointsInRect = false
			break
		}
	}
	if allCirclePointsInRect {
		return RelationshipRectangleCircleContainedByRectangle
	}

	// Check if the rectangle is fully contained within the circle
	rectVertices := []Point[T]{rect.topLeft, rect.topRight, rect.bottomLeft, rect.bottomRight}
	allVerticesInCircle := true
	for _, vertex := range rectVertices {
		if c.center.DistanceToPoint(vertex, opts...) > float64(c.radius) {
			allVerticesInCircle = false
			break
		}
	}
	if allVerticesInCircle {
		return RelationshipRectangleCircleContainedByCircle
	}

	// Check if the circle intersects the rectangle
	rectEdges := rect.Edges()
	for _, edge := range rectEdges {
		if c.RelationshipToLineSegment(edge, opts...) >= RelationshipLineSegmentCircleIntersecting {
			return RelationshipRectangleCircleIntersection
		}
	}

	return RelationshipRectangleCircleMiss
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
//
// todo: ensure other Rotate functions specify directionality.
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
