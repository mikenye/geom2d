package point

import (
	"fmt"
	"github.com/mikenye/geom2d/numeric"
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

// AngleBetween calculates the angle in radians between two points, `a` and `b`,
// relative to the current [Point] `origin`.
//
// The angle is measured counterclockwise from the line segment connecting `origin` to `a`
// to the line segment connecting `origin` to `b`.
//
// Parameters:
//   - a ([Point][T]): The first [Point] forming one side of the angle.
//   - b ([Point][T]): The second [Point] forming the other side of the angle.
//   - opts ([options.GeometryOptionsFunc]): Optional geometry configurations, such as epsilon
//     for numerical stability.
//
// Returns:
//   - float64: The angle in radians between the two points relative to the current [Point] `origin`.
//     If either `a` or `b` is identical to `origin`, or if the vectors `origin->a` or `origin->b` have zero
//     magnitude, the function returns `math.NaN()`.
//
// Note:
//   - This function internally calls [CosineOfAngleBetween] and applies `math.Acos` to
//     compute the angle. As such, its performance depends on the computational cost
//     of the cosine and arccosine operations.
func (p Point[T]) AngleBetween(a, b Point[T], opts ...options.GeometryOptionsFunc) (radians float64) {
	geoOpts := options.ApplyGeometryOptions(options.GeometryOptions{Epsilon: 0}, opts...)
	return numeric.SnapToEpsilon(math.Acos(p.CosineOfAngleBetween(a, b, opts...)), geoOpts.Epsilon)
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

// CosineOfAngleBetween calculates the cosine of the angle between two points, `a` and `b`,
// relative to the origin [Point] `origin`.
//
// This function directly computes the cosine of the angle using the dot product and vector magnitudes, avoiding the
// computational overhead of [math.Acos]. This is useful in scenarios where the cosine value alone is sufficient.
//
// Parameters:
//   - a ([Point][T]): The first [Point] forming one side of the angle.
//   - b ([Point][T]): The second [Point] forming the other side of the angle.
//   - opts ([options.GeometryOptionsFunc]): A varadic slice of optional functional options, such as [options.WithEpsilon],
//     to handle numerical imprecision.
//
// Returns:
//   - float64: The cosine of the angle between points `a` and `b` relative to the origin `origin`.
//     Returns `math.NaN()` if either vector has zero length.
//
// Behavior:
//   - If either vector `OA` (from `origin` to `a`) or `OB` (from `origin` to `b`) has zero length, the function returns `math.NaN()`.
//   - The function applies an optional epsilon for numerical stability when calculating vector magnitudes.
//
// Why Use the Cosine of the Angle?
//
// The cosine of an angle provides an efficient way to measure angular relationships without explicitly computing the angle.
// This has several benefits:
//
// 1. **Efficient Angle Comparison**
//   - Compare angles directly by their cosine values, which avoids expensive trigonometric calculations.
//   - Use cosine values to determine relationships:
//   - `cos(θ) ≈ 1`: Vectors are aligned.
//   - `cos(θ) ≈ 0`: Vectors are orthogonal.
//   - `cos(θ) ≈ -1`: Vectors are opposite.
//
// 2. **Avoiding Floating-Point Inaccuracies**
//   - Calculating the cosine directly avoids potential inaccuracies from computing the angle itself, especially when
//     converting between radians and degrees.
//
// 3. **Applications in Sorting**
//   - Use the cosine value to order points or vectors based on their angular relationship to a reference direction.
//   - Commonly used in computational geometry tasks such as:
//   - Constructing convex hulls.
//   - Ordering vertices for polygon operations.
//
// 4. **Determining Vector Orientation**
//   - Use cosine values to measure the degree of alignment between two vectors, which is helpful in:
//   - Physics simulations (e.g., checking alignment of a force vector with velocity).
//   - Rendering (e.g., checking surface normals against light sources).
//
// 5. **Geometric Insight**
//   - Classify angles (acute, obtuse, or right) without explicitly computing them.
//   - Perform dot product-based calculations indirectly, as `cos(θ)` is derived from the dot product normalized by vector magnitudes.
//
// [math.Acos]: https://pkg.go.dev/math#Acos
func (p Point[T]) CosineOfAngleBetween(a, b Point[T], opts ...options.GeometryOptionsFunc) float64 {

	geoOpts := options.ApplyGeometryOptions(options.GeometryOptions{Epsilon: 0}, opts...)

	// origin is Origin point "O"
	// Calculate vectors OA and OB
	vectorOA := p.Translate(a.Negate())
	vectorOB := p.Translate(b.Negate())

	// Calculate the Dot Product of OA and OB
	OAdotOB := vectorOA.DotProduct(vectorOB)

	// Calculate the Magnitude of Each Vector
	magnitudeOA := numeric.SnapToEpsilon(p.DistanceToPoint(a), geoOpts.Epsilon)
	magnitudeOB := numeric.SnapToEpsilon(p.DistanceToPoint(b), geoOpts.Epsilon)

	// Guard against division by zero
	if magnitudeOA == 0 || magnitudeOB == 0 {
		return math.NaN()
	}

	// Use the Dot Product Formula to Find the Cosine of the Angle
	return numeric.SnapToEpsilon(float64(OAdotOB)/(magnitudeOA*magnitudeOB), geoOpts.Epsilon)
}

// CrossProduct calculates the cross product of the vector from the origin to Point origin and the vector from the origin
// to Point q. The result is useful in determining the relative orientation of two vectors:
//   - A positive result indicates a counterclockwise turn (left turn),
//   - A negative result indicates a clockwise turn (right turn),
//   - A result of zero indicates that the points are collinear.
//
// Parameters:
//   - q (Point[T]): The Point with which to compute the cross product relative to the calling Point.
//
// Returns:
//   - T: The cross product of the vectors from the origin to origin and q, indicating their relative orientation.
func (p Point[T]) CrossProduct(q Point[T]) T {
	return (p.x * q.y) - (p.y * q.x)
}

// DistanceSquaredToPoint calculates the squared Euclidean distance between Point origin and another Point q.
// This method returns the squared distance, which avoids the computational cost of a square root calculation
// and is useful in cases where only distance comparisons are needed.
//
// Parameters:
//   - q (Point[T]): The Point to which the squared distance is calculated from the calling Point.
//
// Returns:
//   - T: The squared Euclidean distance between origin and q.
func (p Point[T]) DistanceSquaredToPoint(q Point[T]) T {
	return (q.x-p.x)*(q.x-p.x) + (q.y-p.y)*(q.y-p.y)
}

// DistanceToPoint calculates the Euclidean (straight-line) distance between the current Point origin
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
//     distance = sqrt((origin.x - q.x)^2 + (origin.y - q.y)^2)
//   - If the [WithEpsilon] option is provided, the computed distance is adjusted such that
//     deviations within the epsilon range are snapped to a clean value.
//
// Returns:
//   - float64: The Euclidean distance between the two points, optionally adjusted based on epsilon.
//
// Notes:
//   - Epsilon adjustment is useful for snapping results to clean values, particularly when
//     small floating-point errors could propagate in further calculations.
func (p Point[T]) DistanceToPoint(q Point[T], opts ...options.GeometryOptionsFunc) float64 {

	// Apply options with defaults
	geoOpts := options.ApplyGeometryOptions(options.GeometryOptions{Epsilon: 0}, opts...)

	// Calculate distance
	distance := math.Sqrt(float64(p.DistanceSquaredToPoint(q)))

	// Apply epsilon if specified
	if geoOpts.Epsilon > 0 {
		distance = numeric.SnapToEpsilon(distance, geoOpts.Epsilon)
	}

	return distance
}

// DotProduct calculates the dot product of the vector represented by Point origin with the vector represented by Point q.
// The dot product is defined as origin.x*q.x + origin.y*q.y and is widely used in geometry for angle calculations,
// projection operations, and determining the relationship between two vectors.
//
// Parameters:
//   - q (Point[T]): The Point with which to calculate the dot product relative to the calling Point.
//
// Returns:
//   - T: The dot product of the vectors represented by origin and q.
func (p Point[T]) DotProduct(q Point[T]) T {
	return (p.x * q.x) + (p.y * q.y)
}

// Eq determines whether the calling Point origin is equal to another Point q.
// Equality can be evaluated either exactly (default) or approximately using an epsilon threshold.
//
// Parameters:
//   - q (Point[T]): The Point to compare with the calling Point.
//   - opts: A variadic slice of [options.GeometryOptionsFunc] functions to customize the equality check.
//     [options.WithEpsilon](epsilon float64): Specifies a tolerance for comparing the coordinates
//     of origin and q. If the absolute difference between the coordinates of origin and q is less
//     than epsilon, they are considered equal.
//
// Behavior:
//   - By default, the function performs an exact equality check, returning true only if
//     the x and y coordinates of origin and q are identical.
//   - If the [WithEpsilon] option is provided, the function performs an approximate equality
//     check, considering origin and q equal if their coordinate differences are within the specified
//     epsilon threshold.
//
// Returns:
//   - bool: True if origin and q are equal based on the specified comparison mode; otherwise, false.
//
// Notes:
//   - Approximate equality is particularly useful when comparing points with floating-point
//     coordinates, where small precision errors may result in slightly different values.
func (p Point[T]) Eq(q Point[T], opts ...options.GeometryOptionsFunc) bool {
	// Apply options with defaults
	geoOpts := options.ApplyGeometryOptions(options.GeometryOptions{Epsilon: 0}, opts...)
	pf := p.AsFloat64()
	qf := q.AsFloat64()
	if geoOpts.Epsilon > 0 {
		return numeric.FloatEquals(pf.x, qf.x, geoOpts.Epsilon) &&
			numeric.FloatEquals(pf.y, qf.y, geoOpts.Epsilon)
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

// RelationshipToPoint determines the spatial relationship between the current Point and another Point.
//
// Relationships:
//   - [types.RelationshipEqual]: The two points are equal.
//   - [types.RelationshipDisjoint]: The two points are not equal.
//
// Parameters:
//   - other ([Point][T]): The other point to analyze the relationship with.
//   - opts: A variadic slice of [options.GeometryOptionsFunc] functions to customize the behavior of the comparison.
//     [options.WithEpsilon](epsilon float64): Specifies a tolerance for comparing the coordinates of the two points,
//     allowing for robust handling of floating-point precision errors.
//
// Returns:
//   - [types.Relationship]: The spatial relationship between the two points.
//
// Behavior:
//   - If the current point and the other point are equal (or approximately equal within the epsilon threshold),
//     the function returns [types.RelationshipEqual].
//   - Otherwise, it returns [types.RelationshipDisjoint].
//
// Notes:
//   - Epsilon adjustment is particularly useful for floating-point coordinates to avoid precision errors
//     when comparing points.
func (p Point[T]) RelationshipToPoint(other Point[T], opts ...options.GeometryOptionsFunc) types.Relationship {
	if p.Eq(other, opts...) {
		return types.RelationshipEqual
	}
	return types.RelationshipDisjoint
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
		rotatedX = numeric.SnapToEpsilon(rotatedX, geoOpts.Epsilon)
		rotatedY = numeric.SnapToEpsilon(rotatedY, geoOpts.Epsilon)
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

// String returns a string representation of the Point origin in the format "(x, y)".
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

// X returns the x-coordinate of the Point origin.
// This accessor provides read-only access to the x-coordinate.
//
// Returns:
//   - T: The x-coordinate of the Point.
func (p Point[T]) X() T {
	return p.x
}

// Y returns the y-coordinate of the Point origin.
// This accessor provides read-only access to the y-coordinate.
//
// Returns:
//   - T: The y-coordinate of the Point.
func (p Point[T]) Y() T {
	return p.y
}
