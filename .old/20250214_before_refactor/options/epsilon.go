package options

// WithEpsilon returns an [GeometryOptionsFunc] that sets the Epsilon value for functions that support it.
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
//   - An [GeometryOptionsFunc] function that modifies the Epsilon field in the geomOptions struct.
func WithEpsilon(epsilon float64) GeometryOptionsFunc {
	return func(opts *GeometryOptions) {
		if epsilon < 0 {
			epsilon = 0 // Default to no adjustment
		}
		opts.Epsilon = epsilon
	}
}
