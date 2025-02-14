// Package options provides configurable settings for geometric operations in the geom2d library.
//
// This package defines a functional options pattern, allowing users to modify the behavior
// of various geometric functions without changing their signatures. The options
// are applied using functional parameters that modify a GeometryOptions struct.
//
// # Key Features
//
//   - Floating-Point Precision Control: The Epsilon parameter allows users to define
//     a tolerance for numerical comparisons, mitigating precision issues in floating-point arithmetic.
//   - Functional Options Pattern: The GeometryOptionsFunc type provides a way to apply
//     optional configurations without requiring additional parameters in function signatures.
//
// # Functional Options
//
// The package provides the following functional options:
//
// - WithEpsilon(epsilon float64) GeometryOptionsFunc: Sets a small tolerance value for floating-point operations.
//
// These options are applied using ApplyGeometryOptions, which takes a default GeometryOptions struct
// and modifies it based on the provided options.
//
// This approach ensures a clean API while allowing flexible configuration for numerical stability
// in geometric computations.
package options
