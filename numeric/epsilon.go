package numeric

import "math"

// SnapToEpsilon adjusts a floating-point value to eliminate small numerical imprecisions
// by snapping it to the nearest whole number if the difference is within a specified epsilon.
//
// Parameters:
//   - value: The floating-point number to adjust.
//   - epsilon: A small positive threshold. If the absolute difference between `value` and
//     the nearest whole number is less than `epsilon`, the value is snapped to that whole number.
//
// Returns:
//   - A floating-point number adjusted based on the specified epsilon, or the original value
//     if no adjustment is needed.
func SnapToEpsilon(value, epsilon float64) float64 {
	rounded := math.Round(value)
	if math.Abs(value-rounded) < epsilon {
		return rounded
	}
	return value
}
