package point

import (
	"fmt"
	"github.com/mikenye/geom2d/options"
	"github.com/mikenye/geom2d/types"
	"image"
	"math"
)

// Point represents a point in two-dimensional space with x and y coordinates of a generic numeric type T.
// The Point struct provides methods for common vector operations such as addition, subtraction, and distance
// calculations, making it versatile for computational geometry and graphics applications.
//
// Type Parameter:
//   - T: The numeric type for the coordinates, constrained to signed number types by the [types.SignedNumber] interface.
type Point[T types.SignedNumber] struct {
	x T
	y T
}

// New creates a new Point with the specified x and y coordinates.
//
// This function is generic and requires the x and y values to satisfy the [types.SignedNumber] constraint.
//
// Parameters:
//   - x (T): The x-coordinate of the point.
//   - y (T): The y-coordinate of the point.
//
// Returns:
//   - Point[T]: A new Point instance with the given coordinates.
func New[T types.SignedNumber](x, y T) Point[T] {
	return Point[T]{
		x: x,
		y: y,
	}
}

// NewFromImagePoint creates and returns a new Point with integer x and y coordinates
// based on an [image.Point]. This function is useful for converting between graphics and
// computational geometry representations of points.
//
// Parameters:
//   - q ([image.Point]): An [image.Point] representing the source coordinates for the new point.
//
// Returns:
//   - Point[int]: A new Point with coordinates corresponding to the x and y values of the provided [image.Point].
func NewFromImagePoint(q image.Point) Point[int] {
	return Point[int]{
		x: q.X,
		y: q.Y,
	}
}

// AsFloat32 converts the Point's x and y coordinates to the float32 type, returning a new Point[float32].
// This method is useful when higher precision or floating-point arithmetic is needed on the coordinates.
//
// Returns:
//   - Point[float32]: A new Point with x and y coordinates converted to float32.
func (p Point[T]) AsFloat32() Point[float32] {
	return Point[float32]{
		x: float32(p.x),
		y: float32(p.y),
	}
}

// AsFloat64 converts the Point's x and y coordinates to the float64 type, returning a new Point[float64].
// This method is useful when higher precision or floating-point arithmetic is needed on the coordinates.
//
// Returns:
//   - Point[float64]: A new Point with x and y coordinates converted to float64.
func (p Point[T]) AsFloat64() Point[float64] {
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

// todo: doc comments, unit test, example func
func (p Point[T]) Coordinates() (x, y T) {
	return p.x, p.y
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
func (p Point[T]) Eq(q Point[T], opts ...options.GeometryOptionsFunc) bool {
	// Apply options with defaults
	geoOpts := options.ApplyGeometryOptions(options.GeometryOptions{Epsilon: 0}, opts...)

	if geoOpts.Epsilon > 0 {
		return math.Abs(float64(p.x)-float64(q.x)) < geoOpts.Epsilon &&
			math.Abs(float64(p.y)-float64(q.y)) < geoOpts.Epsilon
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
	return New[T](-p.x, -p.y)
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
func (p Point[T]) Rotate(pivot Point[T], radians float64, opts ...options.GeometryOptionsFunc) Point[float64] {
	// Apply options with defaults
	geoOpts := options.ApplyGeometryOptions(options.GeometryOptions{Epsilon: 0}, opts...)

	pf := p.AsFloat64()
	originf := pivot.AsFloat64()

	// Translate the point to the origin
	translatedX := pf.x - originf.x
	translatedY := pf.y - originf.y

	// Apply rotation
	rotatedX := translatedX*math.Cos(radians) - translatedY*math.Sin(radians)
	rotatedY := translatedX*math.Sin(radians) + translatedY*math.Cos(radians)

	// Apply epsilon if specified
	if geoOpts.Epsilon > 0 {
		rotatedX = types.SnapToEpsilon(rotatedX, geoOpts.Epsilon)
		rotatedY = types.SnapToEpsilon(rotatedY, geoOpts.Epsilon)
	}

	// Translate back
	newX := rotatedX + originf.x
	newY := rotatedY + originf.y

	return New(newX, newY)
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
	return New(
		ref.x+(p.x-ref.x)*k,
		ref.y+(p.y-ref.y)*k,
	)
}

// String returns a string representation of the Point p in the format "(x, y)".
// This provides a readable format for the point’s coordinates, useful for debugging
// and displaying points in logs or output.
//
// Returns:
//   - string: A string representation of the Point in the format "(x, y)".
func (p Point[T]) String() string {
	return fmt.Sprintf("(%v,%v)", p.x, p.y)
}

// Translate moves the Point by a given displacement vector.
//
// Parameters:
//   - delta (Point[T]): The displacement vector to apply.
//
// Returns:
//   - Point[T]: A new Point resulting from the translation.
func (p Point[T]) Translate(delta Point[T]) Point[T] {
	return New[T](p.x+delta.x, p.y+delta.y)
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
