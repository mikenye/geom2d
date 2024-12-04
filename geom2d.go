// Package geom2d provides a comprehensive set of tools for computational geometry in two dimensions.
//
// This file serves as the entry point for the package and contains utility functions and foundational
// concepts shared across all geometric types. These include global constants, enums for relationships,
// and generic helper functions.
//
// The geom2d package is built around core types like Point, LineSegment, Circle, Rectangle, and PolyTree,
// supporting a wide range of operations including transformations, boolean geometry, and spatial relationships.
//
// Designed for both performance and clarity, geom2d leverages Go generics to handle various numeric types
// and provides intuitive APIs for working with 2D geometric data.

package geom2d

import (
	"fmt"
	"math"
)

// CircleLineSegmentRelationship (CLR) describes the spatial relationship between a line segment
// and a circle.
type CircleLineSegmentRelationship int

const (
	// CLROutside indicates that the line segment lies completely outside the circle,
	// with no intersection or tangency.
	CLROutside CircleLineSegmentRelationship = iota

	// CLRInside indicates that the line segment lies completely within the circle,
	// with both endpoints inside the circle's boundary.
	CLRInside

	// CLRIntersecting indicates that the line segment intersects the circle at exactly
	// two distinct points.
	CLRIntersecting

	// CLRTangent indicates that the line segment is tangent to the circle, touching it
	// at exactly one point with the tangent forming a 90-degree angle to the circle's radius.
	CLRTangent

	// CLROneEndOnCircumferenceOutside indicates that one endpoint of the line segment
	// lies on the circle's boundary, while the other endpoint lies outside the circle.
	CLROneEndOnCircumferenceOutside

	// CLROneEndOnCircumferenceInside indicates that one endpoint of the line segment
	// lies on the circle's boundary, while the other endpoint lies inside the circle.
	CLROneEndOnCircumferenceInside

	// CLRBothEndsOnCircumference indicates that both endpoints of the line segment lie
	// exactly on the circle's boundary. The line segment does not extend inside or outside
	// the circle.
	CLRBothEndsOnCircumference
)

func (r *CircleLineSegmentRelationship) String() string {
	switch *r {
	case CLROutside:
		return "CLROutside"
	case CLRInside:
		return "CLRInside"
	case CLRIntersecting:
		return "CLRIntersecting"
	case CLRTangent:
		return "CLRTangent"
	case CLROneEndOnCircumferenceOutside:
		return "CLROneEndOnCircumferenceOutside"
	case CLROneEndOnCircumferenceInside:
		return "CLROneEndOnCircumferenceInside"
	case CLRBothEndsOnCircumference:
		return "CLRBothEndsOnCircumference"
	default:
		panic(fmt.Errorf("unsupported CircleLineSegmentRelationship"))
	}
}

// PointCircleRelationship (PCR) represents the spatial relationship of a point to a circle.
type PointCircleRelationship int

const (
	PCROutside         PointCircleRelationship = iota // Point is outside the circle
	PCROnCircumference                                // Point lies exactly on the circle's circumference
	PCRInside                                         // Point is inside the circle
)

// PointRectangleRelationship (PRR) describes the relationship of a point to a rectangle.
type PointRectangleRelationship int

const (
	PRROutside  PointRectangleRelationship = iota // The point lies outside the rectangle.
	PRRInside                                     // The point lies strictly inside the rectangle.
	PRROnVertex                                   // The point lies on a vertex of the rectangle.
	PRROnEdge                                     // The point lies on an edge of the rectangle.
)

// RectangleLineSegmentRelationship (RLR) describes the relationship between a rectangle and a line segment.
type RectangleLineSegmentRelationship int

const (
	RLROutside                 RectangleLineSegmentRelationship = iota // The segment lies entirely outside the rectangle.
	RLROutsideEndTouchesEdge                                           // The segment lies outside the rectangle and one end touches an edge.
	RLROutsideEndTouchesVertex                                         // The segment lies outside the rectangle and one end touches a vertex.
	RLRInside                                                          // The segment lies entirely within the rectangle.
	RLRInsideEndTouchesEdge                                            // The segment lies within the rectangle and one end touches an edge without crossing the boundary.
	RLRInsideEndTouchesVertex                                          // The segment lies within the rectangle and one end touches a vertex without crossing the boundary.
	RLROnEdge                                                          // The segment lies entirely on one edge of the rectangle.
	RLROnEdgeEndTouchesVertex                                          // The segment lies entirely on one edge of the rectangle and one or both ends touch a vertex.
	RLRIntersects                                                      // The segment crosses one or more edges of the rectangle.
	RLREntersAndExits                                                  // The segment enters through one edge and exits through another.
)

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

func (r *PointCircleRelationship) String() string {
	switch *r {
	case PCROutside:
		return "PCROutside"
	case PCROnCircumference:
		return "PCROnCircumference"
	case PCRInside:
		return "PCRInside"
	default:
		panic(fmt.Errorf("unsupported PointCircleRelationship"))
	}
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

// Option is a functional option type used to configure optional parameters
// in geometric operations. Functions that accept an Option parameter allow
// users to customize behavior without modifying the primary function signature.
//
// Option functions take a pointer to an geomOptions struct and modify its fields
// to apply specific configurations.
//
// Example:
//
//	rotated := p.Rotate(pivot, math.Pi/2, WithEpsilon(1e-10))
//
// In this example, the WithEpsilon function returns an Option that sets the
// Epsilon field in the geomOptions struct, enabling numerical stability adjustments.
type Option func(*geomOptions)

// WithEpsilon returns an Option that sets the Epsilon value in the geomOptions struct.
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
// Example:
//
//	rotated := p.Rotate(pivot, math.Pi/2, WithEpsilon(1e-10))
//	// Configures the rotation to treat values within 1e-10 of zero as zero.
//
// Returns:
//
//	An Option function that modifies the Epsilon field in the geomOptions struct.
func WithEpsilon(epsilon float64) Option {
	return func(opts *geomOptions) {
		if epsilon < 0 {
			epsilon = 0 // Default to no adjustment
		}
		opts.epsilon = epsilon
	}
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
