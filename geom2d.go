// Package geom2d provides a comprehensive set of tools for computational geometry in two dimensions.
//
// The geom2d package is built around core types like [Point], [LineSegment], [Circle], [Rectangle], and [PolyTree],
// supporting a wide range of operations including transformations, boolean geometry, and spatial relationships.
//
// Designed for both performance and clarity, geom2d leverages Go generics to handle various numeric types
// and provides intuitive APIs for working with 2D geometric data.
//
// # Coordinate System
//
// This library assumes a standard Cartesian coordinate system where the x-axis increases to the right and the y-axis
// increases upward. This system is commonly referred to as a right-handed Cartesian coordinate system.
// All geometric operations and relationships (e.g., clockwise or counterclockwise points) are based on this convention.
//
// # Core Geometric Types
//
// The library includes support for the following 2D geometric types:
//
//   - [Point]: Represents a single coordinate in 2D space.
//   - [LineSegment]: Represents a straight line segment defined by two endpoints.
//   - [Rectangle]: Represents an axis-aligned rectangle, defined by its corners.
//   - [Circle]: Represents a circle defined by a center point and radius.
//   - [PolyTree]: Represents a hierarchical structure of polygons, supporting sibling polygons,
//     nested holes and nested islands.
//
// # Support for Generics
//
// geom2d leverages Go’s generics, allowing you to use the library with different numeric types
// (int, float32, float64, etc.). This flexibility ensures the library can adapt to various applications,
// from integer-based grids to floating-point precision computations.
//
// # Precision Control with Epsilon
//
// geom2d incorporates an epsilon parameter in many of its relationship methods to handle floating-point
// precision issues. This allows you to control the tolerance for comparisons, making the library robust
// for real-world applications where precision errors can occur.
//
// # Relationships Between Types
//
// The library provides methods to compute relationships between geometric types, including:
//
//   - [Point]-to-Geometry Relationships: Determine whether a [Point] is inside, outside, or on the edge of a geometry.
//   - [LineSegment]-to-Geometry Relationships: Check if a [LineSegment] intersects, touches, or is entirely inside or outside another geometry.
//   - Geometry-to-Geometry Relationships: Explore relationships between shapes, such as overlap, containment, tangency, and intersection.
//
// # Acknowledgments
//
// geom2d builds upon the work of others and is grateful for the foundations they have laid. Specifically:
//
//   - Martínez et al.: Their paper on Boolean operations on polygons has been instrumental in the implementation of
//     the Martínez algorithm in this library. See [A simple algorithm for Boolean operations on polygons].
//   - Tom Wright: The inspiration for starting this library came from Tom Wright’s repository
//     [Provably Correct Polygon Algorithms] and his accompanying paper. While geom2d follows its own approach,
//     certain ideas have been influenced by his work.
//   - This project is a collaborative effort, with assistance from [OpenAI's Assistant] for brainstorming, debugging,
//     and refining implementations.
//
// [A simple algorithm for Boolean operations on polygons]: https://web.archive.org/web/20230514184409/https://www.sciencedirect.com/science/article/abs/pii/S0925772199000124
// [Provably Correct Polygon Algorithms]: https://github.com/TooOldCoder/Provably-Correct-Polygon-Algorithms
// [OpenAI's Assistant]: https://openai.com/
package geom2d

import (
	"math"
)

// Option is a functional option type used to configure optional parameters
// in geometric operations. Functions that accept an Option parameter allow
// users to customize behavior without modifying the primary function signature.
//
// Option functions take a pointer to an geomOptions struct and modify its fields
// to apply specific configurations.
type Option func(*geomOptions)

// WithEpsilon returns an [Option] that sets the Epsilon value for functions that support it.
// Epsilon is a small positive value used to adjust for floating-point precision errors,
// ensuring numerical stability in geometric calculations.
//
// Parameters:
//   - epsilon: A small positive value specifying the tolerance range. Values within
//     [-epsilon, epsilon] are treated as zero.
//
// Behavior:
//   - When this option is applied, functions will use the specified Epsilon value
//     to handle near-zero results caused by floating-point arithmetic.
//   - If a negative epsilon is provided, it will default to 0 (no adjustment).
//   - If not set (default), Epsilon remains 0, and no adjustment is applied.
//
// Returns:
//   - An [Option] function that modifies the Epsilon field in the geomOptions struct.
func WithEpsilon(epsilon float64) Option {
	return func(opts *geomOptions) {
		if epsilon < 0 {
			epsilon = 0 // Default to no adjustment
		}
		opts.epsilon = epsilon
	}
}

// SignedNumber is a generic interface representing signed numeric types supported by this package.
// This interface allows functions and structs to operate generically on various numeric types,
// including integer and floating-point types, while restricting to signed values only.
//
// Supported types:
//   - int
//   - int32
//   - int64
//   - float32
//   - float64
//
// By using SignedNumber, functions can handle multiple numeric types without needing to be rewritten
// for each specific type, enabling flexible and type-safe operations across different numeric data.
type SignedNumber interface {
	int | int32 | int64 | float32 | float64
}

// GeometricType is a type constraint for generic programming, representing all geometric types
// supported in this library. It allows functions and methods to operate on any of the included
// geometric types in a type-safe manner.
//
// Supported Types:
//   - [Point]: Represents a 2D point with x and y coordinates.
//   - [LineSegment]: Represents a line segment defined by two endpoints.
//   - [Rectangle]: Represents an axis-aligned rectangle defined by its corner points.
//   - [Circle]: Represents a circle defined by its center point and radius.
//   - [PolyTree]: Represents a hierarchical polygon structure, supporting holes and islands.
//
// Type Parameter:
//   - T: A type that satisfies the [SignedNumber] constraint, representing the numeric type
//     used for the coordinates and dimensions (e.g., int, float64).
//
// Notes:
//   - This interface is used purely as a type constraint for generic functions and methods.
//   - It does not define methods that the types must implement, but rather ensures that
//     only valid geometric types are used.
//
// todo: consider making private & finish doco
type GeometricType[T SignedNumber] interface {
	Point[T] | LineSegment[T] | Rectangle[T] | Circle[T] | PolyTree[T]
}

// todo: consider making private & finish doco
type Measurable[T SignedNumber] interface {
	Area() float64
	Center(opts ...Option) Point[T]
	Perimeter(opts ...Option) float64
}

// todo: consider making private & finish doco
type Spatial[T SignedNumber] interface {
	BoundingBox() Rectangle[T]
	ContainsPoint(p Point[T]) bool
}

// todo: consider making private & finish doco
type Transformable[T SignedNumber, R GeometricType[T]] interface {
	Rotate(pivot Point[T], radians float64, opts ...Option) R
	Scale(ref Point[T], k T) R
	Translate(delta Point[T]) R
}

// todo: consider making private & finish doco
type Geometry[T SignedNumber, R GeometricType[T]] interface {
	Measurable[T]
	Spatial[T]
	Transformable[T, R]
	Eq(other R, opts ...Option) bool
	Intersects(g R) bool
}

// ReflectionAxis specifies the axis or line across which a point or line segment should be reflected.
//
// This type defines the possible axes for reflection, including the standard x-axis and y-axis,
// as well as an arbitrary line defined by a custom line segment.
type ReflectionAxis int

const (
	// ReflectAcrossXAxis reflects a point or line segment across the x-axis, flipping the y-coordinate.
	ReflectAcrossXAxis ReflectionAxis = iota

	// ReflectAcrossYAxis reflects a point or line segment across the y-axis, flipping the x-coordinate.
	ReflectAcrossYAxis

	// ReflectAcrossCustomLine reflects a point or line segment across an arbitrary line defined by a LineSegment.
	// This line segment can be specified as an additional argument to the Reflect method.
	ReflectAcrossCustomLine
)

// adjacentInSlice checks if two values are adjacent in the slice.
func adjacentInSlice[T comparable](s []T, a, b T) bool {
	var i, i2 int
	for i = 0; i < len(s); i++ {
		i2 = i + 1 // todo: (i+1)%len(s) and remove all the %len(s) below?
		if s[i%len(s)] == a && s[i2%len(s)] == b || s[i%len(s)] == b && s[i2%len(s)] == a {
			return true
		}
	}
	return false
}

// countOccurrencesInSlice returns the number of occurrences of an element in a slice.
func countOccurrencesInSlice[T comparable](s []T, element T) int {
	count := 0
	for _, v := range s {
		if v == element {
			count++
		}
	}
	return count
}

// inOrder returns true if b lies between a and c
func inOrder[T SignedNumber](a, b, c T) bool {
	return (a-b)*(b-c) > 0
}

// geomOptions defines a set of configurable parameters for geometric operations.
// These options allow users to customize the behavior of functions in the library,
// such as applying numerical stability adjustments or other optional features.
type geomOptions struct {
	// epsilon is a small positive value used to adjust for floating-point precision errors.
	// When set, values within the range [-Epsilon, Epsilon] are treated as zero in
	// calculations to improve numerical stability. A value of 0 disables this adjustment.
	//
	// For example:
	//   - epsilon > 0: Small deviations caused by floating-point arithmetic are corrected.
	//   - epsilon = 0: No adjustment is applied, leaving results as-is.
	//
	// Default: 0 (no epsilon adjustment)
	epsilon float64
}

// applyOptions applies a set of functional options to a given options struct,
// starting with a set of default values.
//
// Parameters:
//   - defaults: The initial geomOptions struct containing default values.
//   - opts: A variadic slice of Option functions that modify the geomOptions struct.
//
// Behavior:
//   - Each Option function in the opts slice is applied in the order it is provided.
//   - The defaults parameter serves as a base configuration, which can be
//     overridden by the provided geomOptions.
//
// Returns:
//
// A new geomOptions struct that reflects the default values combined with any
// modifications made by the Option functions.
//
// Example:
//
//	defaults := geomOptions{Epsilon: 0}
//	geomOptions := applyOptions(defaults, WithEpsilon(1e-10))
//	fmt.Println(geomOptions.Epsilon) // Output: 1e-10
//
// This function is used internally to provide a consistent way to handle
// optional parameters across multiple functions.
func applyOptions(defaults geomOptions, opts ...Option) geomOptions {
	for _, opt := range opts {
		opt(&defaults)
	}
	return defaults
}

// applyEpsilon adjusts a floating-point value to eliminate small numerical imprecisions
// by snapping it to the nearest whole number if the difference is within a specified epsilon.
//
// Parameters:
//   - value: The floating-point number to adjust.
//   - epsilon: A small positive threshold. If the absolute difference between `value` and
//     the nearest whole number is less than `epsilon`, the value is snapped to that whole number.
//
// Behavior:
//   - If `math.Abs(value - math.Round(value)) < epsilon`, the function returns the rounded value.
//   - Otherwise, it returns the original value unchanged.
//
// Example:
//
//	result := applyEpsilon(-0.9999999999999998, 1e-4)
//	// Output: -1.0 (snapped to the nearest whole number)
//
//	result := applyEpsilon(1.0001, 1e-4)
//	// Output: 1.0001 (unchanged, as it's outside the epsilon threshold)
//
// Returns:
//
//	A floating-point number adjusted based on the specified epsilon, or the original value
//	if no adjustment is needed.
//
// Notes:
//   - This function is commonly used to address floating-point imprecision in geometric
//     computations where small deviations can accumulate and affect results.
func applyEpsilon(value, epsilon float64) float64 {
	// Round to the nearest whole number if within epsilon
	rounded := math.Round(value)
	if math.Abs(value-rounded) < epsilon {
		return rounded
	}
	// Otherwise, return the original value
	return value
}
