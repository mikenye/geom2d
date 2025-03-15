// Package circle provides a representation of circles in a two-dimensional space,
// along with methods for geometric operations, transformations, and relationships
// with other geometric primitives.
//
// Circles are fundamental in computational geometry, used in collision detection, spatial indexing, and graphics.
//
// # Overview
//
// The [Circle] type represents a circle defined by a center point and a radius. It supports
// operations such as computing the area and circumference, scaling, translating, and rotating
// the circle, as well as determining its relationship with points and other geometric entities.
//
// This package also includes Bresenham's circle algorithm for rasterization, enabling efficient
// integer-based rendering of circles.
//
// # Features
//
//   - Creation of circles from coordinates or points.
//   - Type conversion to different numeric representations.
//   - Relationship checks with points, including containment and intersection. // TODO: implement
//   - Support for geometric transformations such as translation, rotation, and scaling.
//   - Efficient rasterization using Bresenham's circle algorithm.
//
// This package is part of the geom2d library and integrates with other geometric primitives
// such as points, line segments, and polygons.
package circle

import (
	"encoding/json"
	"fmt"
	"github.com/mikenye/geom2d"
	"github.com/mikenye/geom2d/numeric"
	"github.com/mikenye/geom2d/point"
	"github.com/mikenye/geom2d/types"
	"math"
)

// Circle represents a circle in 2D space with a center point and a radius.
//
// The Circle type provides methods for calculating its circumference and area,
// determining if a point lies within the circle, checking if a line segment
// intersects the circle, and checking the relationship between the circle and other geometric shapes.
type Circle struct {
	center point.Point // The center point of the circle
	radius float64     // The radius of the circle
}

// New creates a new [Circle] with the specified center coordinates and radius.
//
// Parameters:
//   - x, y (T): The center coordinates of the [Circle].
//   - radius (T): The radius of the circle, of generic type T, which must satisfy the [types.SignedNumber] constraint.
//
// Returns:
//   - Circle[T]: A new Circle with the specified center and radius.
func New(x, y, radius float64) Circle {
	return Circle{
		center: point.New(x, y),
		radius: math.Abs(radius),
	}
}

// NewFromPoint creates a new [Circle] with the specified center [point.Point] and radius.
//
// Parameters:
//   - center ([point.Point][T]): The center [point.Point] of the [Circle].
//   - radius (T): The radius of the circle, of generic type T, which must satisfy the [types.SignedNumber] constraint.
//
// Returns:
//   - Circle[T]: A new Circle with the specified center and radius.
func NewFromPoint(center point.Point, radius float64) Circle {
	return Circle{
		center: center,
		radius: math.Abs(radius),
	}
}

// Area calculates the area of the circle.
//
// Returns:
//   - float64: The area of the circle, computed as π * radius².
func (c Circle) Area() float64 {
	return math.Pi * c.radius * c.radius
}

// Bresenham generates all points on the perimeter of a circle using Bresenham's circle-drawing algorithm.
//
// This method is typically used for rasterized circle rendering.
//
// The function is designed to be used with a for-loop, and thus takes a callback yield that processes each point.
// If the callback returns false at any point (if the calling for-loop is terminated, for example), the function
// halts further generation.
//
// This algorithm utilizes integer arithmetic to efficiently calculate the points on the circle,
// making it suitable for rendering or other grid-based operations.
//
// This algorithm requires circles using integer coordinates because Bresenham's circle algorithm relies
// on integer arithmetic to avoid floating-point precision errors.
//
// Parameters:
//   - yield (func([point.Point][int]) bool): A function that processes each generated point.
//     Returning false will stop further point generation.
func (c Circle) Bresenham(yield func(point.Point) bool) {
	var xc, yc, r, x, y, p float64

	xc = c.center.X()
	yc = c.center.Y()
	r = c.radius

	// Starting at the top of the circle
	x = 0
	y = r
	p = 1 - r // Initial decision parameter

	// Yield the initial points for all octants
	for _, pt := range reflectAcrossCircleOctants(xc, yc, x, y) {
		if !yield(pt) {
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
		for _, pt := range reflectAcrossCircleOctants(xc, yc, x, y) {
			if !yield(pt) {
				return
			}
		}
	}
}

// Center returns the center [Point] of the Circle.
//
// Returns:
//   - Point[T]: The center [Point] of the Circle.
func (c Circle) Center() point.Point {
	return c.center
}

// Circumference calculates the circumference of the circle.
//
// Returns:
//   - float64: The circumference of the circle, computed as 2 * π * radius.
func (c Circle) Circumference() float64 {
	return 2 * math.Pi * c.radius
}

// RelationshipToPoint determines the spatial relationship between the Circle and a [point.Point].
//
// This function evaluates whether the point lies outside, on the boundary of, or inside the given circle.
// The possible relationships are:
//   - [types.RelationshipDisjoint]: The point lies outside the circle.
//   - [types.RelationshipIntersection]: The point lies exactly on the circle's boundary.
//   - [types.RelationshipContainedBy]: The point is inside the circle.
//
// Parameters:
//   - p ([point.Point][T]): The point to compare with the current Circle
//   - opts: A variadic slice of [options.GeometryOptionsFunc] functions to customize the equality check.
//     [options.WithEpsilon](epsilon float64): Specifies a tolerance for comparing distances to handle floating-point
//     precision errors.
//
// Returns:
//   - [types.Relationship]: The relationship of the point to the circle, indicating whether the point is disjoint from,
//     on the boundary of, or contained within the circle.
//
// Behavior:
//   - The function computes the [Euclidean distance] between the point and the circle's center.
//   - It compares this distance to the circle's radius (converted to float64 for precision).
//   - If the distance equals the radius, the relationship is [types.RelationshipIntersection].
//   - If the distance is less than the radius, the relationship is [types.RelationshipContainedBy].
//   - Otherwise, the relationship is [types.RelationshipDisjoint].
//
// Notes:
//   - Epsilon adjustments can be used to account for floating-point precision issues when comparing the distance
//     to the circle's radius.
//
// [Euclidean distance]: https://en.wikipedia.org/wiki/Euclidean_distance
func (c Circle) RelationshipToPoint(p point.Point) types.Relationship {
	distancePointToCircleCenter := p.DistanceToPoint(c.center)
	switch {
	case distancePointToCircleCenter == c.radius:
		return types.RelationshipIntersection
	case distancePointToCircleCenter < c.radius:
		return types.RelationshipContainedBy
	default:
		return types.RelationshipDisjoint
	}
}

// Eq determines whether the calling Circle (c) is equal to another Circle (other), either exactly (default)
// or approximately using an epsilon threshold.
//
// Parameters
//   - other: The Circle to compare with the calling Circle.
//   - opts: A variadic slice of [options.GeometryOptionsFunc] functions to customize the equality check.
//     [options.WithEpsilon](epsilon float64): Specifies a numerical tolerance (epsilon) for comparing the
//     center coordinates and radii of the circles, to avoid false negatives due to floating-point imprecision.
//
// Behavior
//   - The function checks whether the center points of the two circles are equal (using the [point.Point.Eq] method
//     of the [point.Point] type) and whether their radii are equal.
//   - If [options.WithEpsilon] is provided, the comparison allows for small differences in the radius values
//     and center coordinates within the epsilon threshold.
//
// Returns
//   - bool: true if the center coordinates and radius of the two circles are equal (or approximately equal
//     within epsilon); otherwise, false.
//
// Notes:
//   - Approximate equality is particularly useful when comparing circles with floating-point
//     coordinates or radii, where small precision errors might otherwise cause inequality.
func (c Circle) Eq(other Circle) bool {
	// Check equality for the center points
	centersEqual := c.center.Eq(other.center)

	// Check equality for the radii with epsilon adjustment
	radiiEqual := numeric.FloatEquals(c.radius, other.radius, geom2d.GetEpsilon())

	// Equal if centers and radii are equal
	return centersEqual && radiiEqual
}

// MarshalJSON serializes Circle as JSON while preserving its original type.
func (c Circle) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Center point.Point `json:"center"`
		Radius float64     `json:"radius"`
	}{
		Center: c.center,
		Radius: c.radius,
	})
}

// Radius returns the radius of the Circle.
//
// Returns:
//   - T: The radius of the Circle.
func (c Circle) Radius() float64 {
	return c.radius
}

// Rotate rotates the Circle's center around a specified pivot [point.Point] by a given angle in radians
// counter-clockwise, while keeping the radius unchanged. Optionally, an epsilon threshold can be applied
// to adjust the precision of the resulting coordinates.
//
// Parameters:
//   - pivot ([point.Point][T]): The [point.Point] around which to rotate the circle's center.
//   - radians: The rotation angle in radians (counter-clockwise).
//   - opts: A variadic slice of [options.GeometryOptionsFunc] functions to customize the behavior of the rotation.
//     [options.WithEpsilon](epsilon float64): Specifies a tolerance for snapping the resulting center coordinates
//     to cleaner values, improving robustness in floating-point calculations.
//
// Returns:
//   - [Circle][float64]: A new [Circle] with the center rotated around the pivot [point.Point] by the specified angle,
//     and with the radius unchanged.
//
// Behavior:
//   - The function rotates the circle's center point around the given pivot by the specified angle using
//     the [point.Point.Rotate] method.
//   - The rotation is performed in a counter-clockwise direction relative to the pivot point.
//   - The radius remains unchanged in the resulting [Circle].
//   - If [options.WithEpsilon] is provided, epsilon adjustments are applied to the rotated coordinates to handle
//     floating-point precision errors.
//
// Notes:
//   - Epsilon adjustment is particularly useful when the rotation involves floating-point calculations
//     that could result in minor inaccuracies.
//   - The returned [Circle] always has a center with float64 coordinates, ensuring precision regardless
//     of the coordinate type of the original [Circle].
func (c Circle) Rotate(pivot point.Point, radians float64) Circle {
	return NewFromPoint(
		c.center.Rotate(pivot, radians),
		c.radius,
	)
}

// Scale scales the radius of the circle by a scalar factor.
//
// Parameters:
//   - factor (T): The factor by which to scale the radius.
//
// Returns:
//   - Circle[T]: A new circle with the radius scaled by the specified factor.
//
// todo: update doc comment, examples after adding numeric.Abs to radius
func (c Circle) Scale(factor float64) Circle {
	return Circle{center: c.center, radius: math.Abs(c.radius * factor)}
}

// String returns a string representation of the Circle, including its center coordinates and radius.
// This is useful for debugging and logging.
//
// Returns a string representation of the Circle in the format "(h, k, r)", where:
//   - h: x-coordinate of the center.
//   - k: y-coordinate of the center.
//   - r: radius
func (c Circle) String() string {
	return fmt.Sprintf("(%f,%f; r=%f)", c.center.X(), c.center.Y(), c.radius)
}

// Translate moves the circle by a specified vector (given as a [point.Point]).
//
// This method shifts the circle's center by the given vector v, effectively
// translating the circle's position in the 2D plane. The radius of the circle
// remains unchanged.
//
// Parameters:
//   - v ([point.Point][T]): The vector by which to translate the circle's center.
//
// Returns:
//   - Circle[T]: A new Circle translated by the specified vector.
func (c Circle) Translate(v point.Point) Circle {
	return Circle{center: c.center.Translate(v), radius: c.radius}
}

// UnmarshalJSON deserializes JSON into a Circle while keeping the exact original type.
func (c *Circle) UnmarshalJSON(data []byte) error {
	var temp struct {
		Center point.Point `json:"center"`
		Radius float64     `json:"radius"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Validate radius (ensure it's non-negative)
	if temp.Radius < 0 {
		return fmt.Errorf("invalid radius: must be non-negative, got %v", temp.Radius)
	}

	c.center = temp.Center
	c.radius = temp.Radius
	return nil
}

// reflectAcrossCircleOctants generates a slice of points that represent the reflection
// of a given point (x, y) across all eight octants of a circle centered at (xc, yc).
//
// The function is typically used in circle-drawing algorithms, such as Bresenham's Circle Algorithm,
// to exploit the symmetry of circles for efficient computation.
//
// Parameters:
//   - xc, yc: The coordinates of the circle's center.
//   - x, y: The coordinates of the point to reflect across the octants.
//
// Returns:
//   - A slice of Point[T] containing the reflected points in the following order:
//     1. Octant 1: (xc + x, yc + y)
//     2. Octant 2: (xc - x, yc + y)
//     3. Octant 8: (xc + x, yc - y)
//     4. Octant 7: (xc - x, yc - y)
//     5. Octant 3: (xc + y, yc + x)
//     6. Octant 4: (xc - y, yc + x)
//     7. Octant 6: (xc + y, yc - x)
//     8. Octant 5: (xc - y, yc - x)
func reflectAcrossCircleOctants(xc, yc, x, y float64) []point.Point {
	return []point.Point{
		point.New(xc+x, yc+y), // Octant 1
		point.New(xc-x, yc+y), // Octant 2
		point.New(xc+x, yc-y), // Octant 8
		point.New(xc-x, yc-y), // Octant 7
		point.New(xc+y, yc+x), // Octant 3
		point.New(xc-y, yc+x), // Octant 4
		point.New(xc+y, yc-x), // Octant 6
		point.New(xc-y, yc-x), // Octant 5
	}
}
