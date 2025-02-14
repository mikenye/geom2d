package options

// GeometryOptionsFunc is a functional option type used to configure optional parameters
// in geometric operations. Functions that accept a GeometryOptionsFunc parameter allow
// users to customize behavior without modifying the primary function signature.
//
// GeometryOptionsFunc functions take a pointer to a GeometryOptions struct and modify its fields
// to apply specific configurations.
type GeometryOptionsFunc func(*GeometryOptions)

// GeometryOptions defines a set of configurable parameters for geometric operations.
// These options allow users to customize the behavior of functions in the library,
// such as applying numerical stability adjustments or other optional features.
type GeometryOptions struct {
	// Epsilon is a small positive value used to adjust for floating-point precision errors.
	// When set, values within the range [-Epsilon, Epsilon] are treated as zero in
	// calculations to improve numerical stability. A value of 0 disables this adjustment.
	//
	// For example:
	//   - Epsilon > 0: Small deviations caused by floating-point arithmetic are corrected.
	//   - Epsilon = 0: No adjustment is applied, leaving results as-is.
	//
	// Default: 0 (no epsilon adjustment)
	Epsilon float64
}

// ApplyGeometryOptions applies a set of functional options to a given options struct,
// starting with a set of default values.
//
// Parameters:
//   - defaults (GeometryOptions): The initial geomOptions struct containing default values.
//   - opts: A variadic slice of GeometryOptionsFunc functions that modify the geomOptions struct.
//
// Behavior:
//   - Each GeometryOptionsFunc function in the opts slice is applied in the order it is provided.
//   - The defaults parameter serves as a base configuration, which can be
//     overridden by the provided geomOptions.
//
// Returns:
//
// A new GeometryOptions struct that reflects the default values combined with any
// modifications made by the GeometryOptionsFunc functions.
//
// This function is used internally to provide a consistent way to handle
// optional parameters across multiple functions.
func ApplyGeometryOptions(defaults GeometryOptions, opts ...GeometryOptionsFunc) GeometryOptions {
	for _, opt := range opts {
		opt(&defaults)
	}
	return defaults
}
