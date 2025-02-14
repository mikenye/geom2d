// Package numeric provides utility functions for numerical computations,
// particularly focused on handling floating-point precision issues and
// operations on signed numbers.
//
// # Overview
//
// The numeric package contains a set of helper functions designed for
// common numerical operations that arise in computational geometry and
// other domains where precision is important. This includes absolute
// value computation, floating-point comparisons with epsilon tolerance,
// and precision adjustments.
//
// # Features
//
//   - Absolute Value Calculation: The Abs function computes the
//     absolute value of any signed number, supporting both integer and
//     floating-point types.
//
//   - Floating-Point Comparisons: Functions such as FloatEquals,
//     FloatGreaterThan, FloatLessThan, and their variants provide
//     robust comparisons between floating-point numbers using an epsilon
//     threshold to mitigate precision errors.
//
//   - Precision Adjustment: The SnapToEpsilon function allows
//     floating-point numbers to be snapped to the nearest whole number if
//     they are within an acceptable tolerance, reducing small precision
//     artifacts.
//
// # Usage
//
// This package is particularly useful in scenarios where direct equality
// checks for floating-point numbers are unreliable due to the inherent
// imprecision of floating-point arithmetic.
package numeric
