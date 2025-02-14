package numeric

import "github.com/mikenye/geom2d/types"

// Abs computes the absolute value of a signed number.
//
// This function is generic and works for any type that satisfies the
// [SignedNumber] constraint (e.g., int, int32, int64, float32, float64).
//
// Parameters:
//   - n (T): The signed number whose absolute value is to be computed.
//
// Returns:
//   - The absolute value of the input number.
func Abs[T types.SignedNumber](n T) T {
	if n < 0 {
		return -n
	}
	return n
}
