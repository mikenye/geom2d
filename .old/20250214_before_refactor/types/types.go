// Package types defines core type constraints and relationships used across the geom2d library.
//
// This package provides foundational types such as SignedNumber, which restricts generic
// operations to signed numeric types, and Relationship, which describes spatial relationships
// between geometric entities.
//
// # Key Features
//
//   - SignedNumber Interface: Defines a type set that includes all signed integer and floating-point types,
//     ensuring that geometric operations remain compatible with various numeric representations.
//   - Relationship Enum: Encapsulates possible geometric relationships between shapes, such as containment,
//     intersection, or equality, allowing for standardized comparisons between geometric objects.
//
// # Usage
//
// This package is primarily used internally within the geom2d library to enable type safety and consistency
// in geometric operations. Functions and structures throughout the library rely on these types to enforce
// correct input parameters and return meaningful results.
//
// See the documentation for each type for more details.
package types
