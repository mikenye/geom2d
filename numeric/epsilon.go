package numeric

import "math"

// FloatEquals returns true if a and b are equal within the given epsilon tolerance.
//
// This is useful for comparing floating-point numbers where direct equality checks may be unreliable
// due to floating point imprecision.
func FloatEquals(a, b, epsilon float64) bool {
	return math.Abs(a-b) <= epsilon
}

// FloatGreaterThan returns true if a is greater than b within the given epsilon tolerance.
//
// This is useful for comparing floating-point numbers where direct equality checks may be unreliable
// due to floating point imprecision.
func FloatGreaterThan(a, b, epsilon float64) bool {
	return a > b && !FloatEquals(a, b, epsilon)
}

// FloatGreaterThanOrEqualTo returns true if a is greater than or equal to b within the given epsilon tolerance.
//
// This is useful for comparing floating-point numbers where direct equality checks may be unreliable
// due to floating point imprecision.
func FloatGreaterThanOrEqualTo(a, b, epsilon float64) bool {
	return a > b || FloatEquals(a, b, epsilon)
}

// FloatLessThan returns true if a is less than b within the given epsilon tolerance.
//
// This is useful for comparing floating-point numbers where direct equality checks may be unreliable
// due to floating point imprecision.
func FloatLessThan(a, b, epsilon float64) bool {
	return a < b && !FloatEquals(a, b, epsilon)
}

// FloatLessThanOrEqualTo checks if 'a' is less than or equal to 'b', using epsilon for approximate equality.
//
// This function returns true if 'a' is strictly less than 'b' or if 'a' and 'b' are approximately
// equal within the given epsilon tolerance.
func FloatLessThanOrEqualTo(a, b, epsilon float64) bool {
	return a < b || FloatEquals(a, b, epsilon)
}

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
