// Package point defines the foundational geometric primitive in the geom2d library, the Point type.
// All other geometric types—such as line segments, rectangles, polygons, etc. are built upon this type.
//
// # Overview
//
// The Point type represents a two-dimensional point with floating-point coordinates. It provides
// fundamental geometric operations such as translation, distance measurement, vector arithmetic, and angle
// calculations. Points are essential building blocks in computational geometry, enabling higher-level
// constructs such as line segments and polygons.
//
// # Key Features
//
// Creation & Type Conversion
//   - Points can be created using New and NewFromImagePoint.
//   - Conversion methods (AsFloat32, AsFloat64, AsInt, AsIntRounded) allow working with different numeric types.
//
// Vector Operations
//   - Basic operations like Translate and Negate enable geometric transformations.
//   - Scale allows uniform scaling around a reference point.
//
// Distance & Angle Measurements
//   - DistanceToPoint and DistanceSquaredToPoint provide Euclidean distance calculations.
//   - AngleBetween and CosineOfAngleBetween help determine angular relationships between points.
//   - CrossProduct and DotProduct support vector orientation and projection calculations.
//
// Equality & Relationships
//   - Eq checks exact or approximate equality (with epsilon-based tolerance).
//   - RelationshipToPoint determines if two points are equal or disjoint.
//
// # Notes
//
//   - Floating-point operations may introduce precision errors. Most comparison operations use the global
//     epsilon value for approximate comparisons, which can be adjusted using [geom2d.SetEpsilon].
//   - This package is optimized for computational geometry applications, balancing precision, performance, and ease of use.
//
// The [Point] type serves as the core building block for all geometric structures in geom2d.
package point

import (
	"encoding/json"
	"fmt"
	"github.com/mikenye/geom2d"
	"github.com/mikenye/geom2d/numeric"
	"github.com/mikenye/geom2d/types"
	"image"
	"math"
)

var origin Point

func init() {
	origin = New(0, 0)
}

// Origin returns the origin point (0,0) in the 2D coordinate system.
//
// This function provides efficient access to a pre-initialized point at the
// coordinate system origin. This is useful for operations that reference the
// origin, such as vector calculations and coordinate transformations.
//
// Returns:
//   - Point: A Point instance representing the origin (0,0).
//
// Note:
//   - This is more efficient than repeatedly creating new origin points with New(0,0).
//   - The returned point is immutable, but a copy is returned, so it can be safely
//     used in any context without affecting the stored origin.
func Origin() Point {
	return origin
}

// Point represents a point in two-dimensional space with x and y coordinates of type float64.
// The Point struct provides methods for common vector operations such as addition, subtraction, and distance
// calculations, making it versatile for computational geometry and graphics applications.
type Point struct {
	x float64
	y float64
}

// New creates a new Point with the specified x and y coordinates.
//
// Parameters:
//   - x (float64): The x-coordinate of the point.
//   - y (float64): The y-coordinate of the point.
//
// Returns:
//   - Point: A new Point instance with the given coordinates.
func New(x, y float64) Point {
	return Point{
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
//   - Point: A new Point with coordinates corresponding to the x and y values of the provided [image.Point].
func NewFromImagePoint(q image.Point) Point {
	return Point{
		x: float64(q.X),
		y: float64(q.Y),
	}
}

// Add returns the sum of two points as if they were vectors.
// It performs component-wise addition:
//
//	(p.X + q.X, p.Y + q.Y)
func (p Point) Add(q Point) Point {
	return Point{
		x: p.X() + q.X(),
		y: p.Y() + q.Y(),
	}
}

// AngleBetween calculates the angle in radians between two points, a and b,
// relative to the current [Point] origin.
//
// The angle is measured counterclockwise from the line segment connecting origin to a
// to the line segment connecting origin to b.
//
// Parameters:
//   - a (Point): The first Point forming one side of the angle.
//   - b (Point): The second Point forming the other side of the angle.
//
// Returns:
//   - float64: The angle in radians between the two points relative to the current Point origin.
//     If either a or b is identical to origin, or if the vectors origin->a or origin->b have zero
//     magnitude, the function returns math.NaN().
//
// Note:
//   - This function internally calls CosineOfAngleBetween and applies math.Acos to
//     compute the angle.
//   - Due to floating-point precision limitations with math.Acos, this function
//     requires a smaller epsilon (~1e-6) than other functions. The issue is that
//     when the cosine value is close to -1 or 1 (angles near 0° or 180°), even tiny
//     floating-point errors are magnified by math.Acos due to the nonlinear nature of
//     the inverse cosine function near these bounds.
func (p Point) AngleBetween(a, b Point) (radians float64) {
	return math.Acos(p.CosineOfAngleBetween(a, b))
}

// Coordinates returns the X and Y coordinates of the Point as separate values.
// This function allows convenient access to the individual components of a Point.
//
// Returns:
//   - x (float64): The X-coordinate of the point.
//   - y (float64): The Y-coordinate of the point.
func (p Point) Coordinates() (x, y float64) {
	return p.x, p.y
}

// CosineOfAngleBetween calculates the cosine of the angle between two points, a and b,
// relative to the origin [Point] origin.
//
// This function directly computes the cosine of the angle using the dot product and vector magnitudes, avoiding the
// computational overhead of math.Acos. This is useful in scenarios where the cosine value alone is sufficient.
//
// Parameters:
//   - a (Point): The first Point forming one side of the angle.
//   - b (Point): The second Point forming the other side of the angle.
//
// Returns:
//   - float64: The cosine of the angle between points a and b relative to the origin.
//     Returns math.NaN() if either vector has zero length.
//
// Behavior:
//   - If either vector OA (from origin to a) or OB (from origin to b) has zero length, the function returns math.NaN().
//   - The function ensures numerical stability through normalization and clamping of the result to [-1,1].
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
func (p Point) CosineOfAngleBetween(a, b Point) float64 {

	// origin is Origin point "O"
	// Calculate vectors OA and OB
	vectorOA := p.Translate(a.Negate())
	vectorOB := p.Translate(b.Negate())

	// Calculate the Dot Product of OA and OB
	OAdotOB := vectorOA.DotProduct(vectorOB)

	// Calculate the Magnitude of Each Vector
	magnitudeOA := p.DistanceToPoint(a)
	magnitudeOB := p.DistanceToPoint(b)

	// Guard against division by zero
	if magnitudeOA == 0 || magnitudeOB == 0 {
		return math.NaN()
	}

	// Use the Dot Product Formula to Find the Cosine of the Angle
	cosTheta := OAdotOB / (magnitudeOA * magnitudeOB)

	// Clamp to [-1,1]
	return math.Max(-1, math.Min(1, cosTheta))
}

// CrossProduct returns the 2D cross product (determinant) of two vectors:
//
//	a × b = a.x * b.y - a.y * b.x
//
// This function is useful in computational geometry for determining relative orientation:
//   - A positive result indicates a counterclockwise turn (left turn),
//   - A negative result indicates a clockwise turn (right turn),
//   - A result of zero indicates that the points are collinear.
//
// Returns:
//   - float64: The signed cross product value.
func (a Point) CrossProduct(b Point) float64 {
	return a.x*b.y - a.y*b.x
}

// DistanceSquaredToPoint calculates the squared Euclidean distance between Point origin and another Point q.
// This method returns the squared distance, which avoids the computational cost of a square root calculation
// and is useful in cases where only distance comparisons are needed.
//
// Parameters:
//   - q (Point): The Point to which the squared distance is calculated from the calling Point.
//
// Returns:
//   - float64: The squared Euclidean distance between origin and q.
func (p Point) DistanceSquaredToPoint(q Point) float64 {
	return (q.x-p.x)*(q.x-p.x) + (q.y-p.y)*(q.y-p.y)
}

// DistanceToPoint calculates the Euclidean (straight-line) distance between the current Point origin
// and another Point q. The result is returned as a float64, providing precise measurement
// of the straight-line separation between the two points.
//
// Parameters:
//   - q (Point): The Point to which the distance is calculated.
//
// Behavior:
//   - The function computes the Euclidean distance using the formula:
//     distance = sqrt((origin.x - q.x)^2 + (origin.y - q.y)^2)
//
// Returns:
//   - float64: The Euclidean distance between the two points.
func (p Point) DistanceToPoint(q Point) float64 {
	// Calculate distance
	return math.Sqrt(p.DistanceSquaredToPoint(q))
}

// DotProduct calculates the dot product of the vector represented by Point origin with the vector represented by Point q.
// The dot product is defined as origin.x*q.x + origin.y*q.y and is widely used in geometry for angle calculations,
// projection operations, and determining the relationship between two vectors.
//
// Parameters:
//   - q (Point): The Point with which to calculate the dot product relative to the calling Point.
//
// Returns:
//   - float64: The dot product of the vectors represented by origin and q.
func (p Point) DotProduct(q Point) float64 {
	return (p.x * q.x) + (p.y * q.y)
}

// Eq determines whether the calling Point origin is equal to another Point q using
// the global epsilon value to account for floating-point precision.
//
// Parameters:
//   - q (Point): The Point to compare with the calling Point.
//
// Behavior:
//   - The function performs an approximate equality check using the global epsilon value,
//     comparing the x and y coordinates of origin and q.
//
// Returns:
//   - bool: True if origin and q are approximately equal within the epsilon threshold; otherwise, false.
//
// Notes:
//   - Approximate equality is particularly useful when comparing points with floating-point
//     coordinates, where small precision errors may result in slightly different values.
func (p Point) Eq(q Point) bool {
	return numeric.FloatEquals(p.x, q.x, geom2d.GetEpsilon()) && numeric.FloatEquals(p.y, q.y, geom2d.GetEpsilon())
}

// MarshalJSON serializes Point as JSON.
func (p Point) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}{
		X: p.X(),
		Y: p.Y(),
	})
}

// Negate returns a new Point with both x and y coordinates negated.
// This operation is equivalent to reflecting the point across the origin
// and is useful in reversing the direction of a vector or preparing
// a point for subtraction via translation.
//
// Returns:
//   - Point: A new Point with negated x and y coordinates.
//
// Notes:
//   - The returned point has the same type as the calling point.
//   - This method does not modify the original point but returns a new instance.
func (p Point) Negate() Point {
	return New(-p.x, -p.y)
}

// RelationshipToPoint determines the spatial relationship between the current Point and another Point.
//
// Relationships:
//   - [types.RelationshipEqual]: The two points are equal.
//   - [types.RelationshipDisjoint]: The two points are not equal.
//
// Parameters:
//   - other (Point): The other point to analyze the relationship with.
//
// Returns:
//   - [types.Relationship]: The spatial relationship between the two points.
//
// Behavior:
//   - If the current point and the other point are equal (or approximately equal within the epsilon threshold),
//     the function returns [types.RelationshipEqual].
//   - Otherwise, it returns [types.RelationshipDisjoint].
func (p Point) RelationshipToPoint(other Point) types.Relationship {
	if p.Eq(other) {
		return types.RelationshipEqual
	}
	return types.RelationshipDisjoint
}

// Rotate rotates the point by a specified angle (in radians), counter-clockwise, around a given pivot point.
//
// Parameters:
//   - pivot (Point): The point around which the rotation is performed.
//   - radians (float64): The angle of rotation, specified in radians.
//
// Behavior:
//   - The method first translates the point to the origin (relative to the pivot),
//     applies the rotation matrix, and then translates the point back to its original position.
//
// Returns:
//   - Point: A new point representing the rotated position.
//
// Notes:
//   - The return type is Point with float64 precision to ensure accurate rotation calculations.
func (p Point) Rotate(pivot Point, radians float64) Point {

	// Translate the point to the origin (pivot)
	translatedX := p.x - pivot.x
	translatedY := p.y - pivot.y

	// Apply rotation
	rotatedX := translatedX*math.Cos(radians) - translatedY*math.Sin(radians)
	rotatedY := translatedX*math.Sin(radians) + translatedY*math.Cos(radians)

	// Translate back
	newX := rotatedX + pivot.x
	newY := rotatedY + pivot.y

	return New(newX, newY)
}

// Scale scales the point by a factor k relative to a reference point ref.
//
// Parameters:
//   - ref (Point): The reference point from which scaling originates.
//   - k (float64): The scaling factor.
//
// Returns:
//   - Point - A new point scaled relative to the reference point.
func (p Point) Scale(ref Point, k float64) Point {
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
func (p Point) String() string {
	return fmt.Sprintf("(%f,%f)", p.x, p.y)
}

// Sub returns the vector from this point to another point.
func (p Point) Sub(q Point) Point {
	return New(p.x-q.x, p.y-q.y)
}

// Translate moves the Point by a given displacement vector.
//
// Parameters:
//   - delta (Point): The displacement vector to apply.
//
// Returns:
//   - Point: A new Point resulting from the translation.
func (p Point) Translate(delta Point) Point {
	return New(p.x+delta.x, p.y+delta.y)
}

// UnmarshalJSON deserializes JSON into a Point.
func (p *Point) UnmarshalJSON(data []byte) error {
	var temp struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	p.x = temp.X
	p.y = temp.Y
	return nil
}

// X returns the x-coordinate of the Point origin.
// This accessor provides read-only access to the x-coordinate.
//
// Returns:
//   - float64: The x-coordinate of the Point.
func (p Point) X() float64 {
	return p.x
}

// Y returns the y-coordinate of the Point origin.
// This accessor provides read-only access to the y-coordinate.
//
// Returns:
//   - float64: The y-coordinate of the Point.
func (p Point) Y() float64 {
	return p.y
}
