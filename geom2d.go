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
// # Relationships Between Geometric Types
//
// This library provides methods to compute relationships between geometric types using a standardized set of relationships:
// [RelationshipDisjoint], [RelationshipIntersection], [RelationshipContainedBy], [RelationshipContains], and [RelationshipEqual].
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
	"fmt"
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

// Relationship defines the spatial or geometric relationship between two shapes or entities
// in 2D space. This type is used across various functions to represent how one geometric
// entity relates to another.
//
// The possible relationships include:
//   - Disjoint: The entities do not intersect, overlap, or touch in any way.
//   - Intersection: The entities share some boundary or overlap.
//   - Contained: One entity is fully within the boundary of another.
//   - ContainedBy: One entity is fully enclosed by another.
//   - Equal: The entities are identical in shape, size, and position.
//
// Each relationship type is represented as a constant value of the Relationship type.
// Functions that evaluate relationships between geometric entities typically return one
// of these constants to describe the spatial relationship between them.
//
// See the individual constant definitions for more details.
type Relationship uint8

// Valid values for Relationship:
const (
	// RelationshipDisjoint indicates that the two entities are completely separate,
	// with no overlap, touching, or intersection.
	RelationshipDisjoint Relationship = iota

	// RelationshipIntersection indicates that the two entities overlap or share
	// a boundary. This includes cases where the entities partially intersect
	// or where they touch at one or more points.
	RelationshipIntersection

	// RelationshipContainedBy indicates that the first entity is fully enclosed
	// within the second entity. The boundary of the first entity does not extend
	// outside the boundary of the second entity.
	RelationshipContainedBy

	// RelationshipContains indicates that the first entity fully encloses the
	// second entity. The boundary of the second entity does not extend outside
	// the boundary of the first entity.
	RelationshipContains

	// RelationshipEqual indicates that the two entities are identical in shape,
	// size, and position. This includes cases where their boundaries coincide exactly.
	RelationshipEqual
)

// flipContainment reverses containment relationships for a [Relationship].
//
// This method is used to swap the roles of containment when interpreting
// relationships. Specifically:
//   - If the [Relationship] is RelationshipContainedBy, it is flipped to RelationshipContains.
//   - If the [Relationship] is RelationshipContains, it is flipped to RelationshipContainedBy.
//   - All other [Relationship] values are returned unchanged.
//
// Returns:
//   - [Relationship]: The flipped or unchanged [Relationship].
//
// Example:
//
//	rel := RelationshipContainedBy
//	flipped := rel.flipContainment()
//	fmt.Println(flipped) // Output: RelationshipContains
func (r Relationship) flipContainment() Relationship {
	switch r {
	case RelationshipContainedBy:
		return RelationshipContains
	case RelationshipContains:
		return RelationshipContainedBy
	default:
		return r
	}
}

// String converts a [Relationship] value to its string representation.
//
// This method provides a human-readable string corresponding to the [Relationship]
// constant, such as RelationshipDisjoint or RelationshipContainedBy. It is useful
// for debugging and logging purposes.
//
// Supported [Relationship] values:
//   - [RelationshipDisjoint]: The objects are disjoint and do not touch or intersect.
//   - [RelationshipIntersection]: The objects intersect or overlap at some point.
//   - [RelationshipContainedBy]: The object is fully contained within another object.
//   - [RelationshipContains]: The object fully contains another object.
//   - [RelationshipEqual]: The objects are identical in size, shape, and position.
//
// Returns:
//   - string: The string representation of the [Relationship].
//
// Panics:
//   - If the [Relationship] value is not supported, this method panics with an error message.
func (r Relationship) String() string {
	switch r {
	case RelationshipDisjoint:
		return "RelationshipDisjoint"
	case RelationshipIntersection:
		return "RelationshipIntersection"
	case RelationshipContainedBy:
		return "RelationshipContainedBy"
	case RelationshipContains:
		return "RelationshipContains"
	case RelationshipEqual:
		return "RelationshipEqual"
	default:
		panic(fmt.Errorf("unsupported Relationship: %d", r))
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

// abs computes the absolute value of a signed number.
//
// This function is generic and works for any type that satisfies the
// [SignedNumber] constraint (e.g., int, int32, int64, float32, float64).
//
// Parameters:
//   - n (T): The signed number whose absolute value is to be computed.
//
// Returns:
//   - The absolute value of the input number.
func abs[T SignedNumber](n T) T {
	if n < 0 {
		return -n
	}
	return n
}
