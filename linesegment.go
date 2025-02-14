package geom2d

import "math"

type LineSegment struct {
	pointUpper, pointLower *Point
}

// Endpoints returns the coordinates of the upper and lower endpoints of the line segment.
// The upper endpoint is the point with the greater Y-coordinate, or if Y-coordinates are equal,
// the one with the smaller X-coordinate. The lower endpoint is the other point.
//
// Returns:
//   - upper, lower: The upper and lower endpoints of the line segment.
func (l *LineSegment) Endpoints() (*Point, *Point) {
	return l.pointUpper, l.pointLower
}

// riseRunScaled returns the scaled rise (dy) and run (dx) of the line segment.
//
// This allows for integer-based slope comparisons, avoiding floating-point precision issues.
//
// Returns:
//   - dy: The scaled difference in Y-coordinates (rise).
//   - dx: The scaled difference in X-coordinates (run).
func (l *LineSegment) riseRunScaled() (int64, int64) {
	return l.pointLower.scaledY - l.pointUpper.scaledY, l.pointLower.scaledX - l.pointUpper.scaledX
}

// Slope computes the slope of the line segment.
//
// Special Cases:
// - If the segment is vertical (X-coordinates are the same), the function returns math.Inf(1).
// - If the segment is horizontal (Y-coordinates are the same), the function returns 0.
//
// Returns:
//   - The slope of the line segment as a float64.
func (l *LineSegment) Slope() float64 {
	// Handle vertical line case: return positive infinity
	if l.pointUpper.scaledX == l.pointLower.scaledX {
		return math.Inf(1)
	}

	// Compute slope (rise / run)
	rise, run := l.riseRunScaled()
	return float64(rise) /
		float64(run)
}

// SlopeCmp compares the slopes of two line segments using integer arithmetic.
//
// The comparison is based on cross-multiplication to avoid floating-point imprecision.
// It follows the same order as Go's `cmp.Compare(a, b)` function:
//
// - Returns `-1` if l is steeper than other (l < other).
// - Returns `1` if other is steeper than l (l > other).
// - Returns `0` if they have the same slope.
//
// Special Cases:
// - If both segments are vertical, they are considered equal (returns 0).
// - If one is vertical, the vertical segment is always considered steeper.
//
// Parameters:
//   - other: The `LineSegment` to compare against.
//
// Returns:
//   - An `int` indicating the comparison result (-1, 0, or 1).
func (l *LineSegment) SlopeCmp(other *LineSegment) int {
	// Get rise and run for both segments
	dy1, dx1 := l.riseRunScaled()
	dy2, dx2 := other.riseRunScaled()

	// Handle vertical lines explicitly
	if dx1 == 0 && dx2 == 0 {
		return 0 // Both are vertical, equal slopes
	}
	if dx1 == 0 {
		return -1 // l is vertical → steeper
	}
	if dx2 == 0 {
		return 1 // other is vertical → steeper
	}

	// Check if slopes are both positive, both negative, or opposing
	slope1Positive := dy1*dx1 > 0
	slope2Positive := dy2*dx2 > 0

	if slope1Positive != slope2Positive {
		// One is positive, one is negative → negative slopes are steeper
		if slope1Positive {
			return 1
		}
		return -1
	}

	// Both slopes are either positive or negative
	// Use cross-multiplication to compare without division:
	// dy1/dx1 < dy2/dx2  →  dy1 * dx2 < dy2 * dx1
	diff := dy1*dx2 - dy2*dx1

	// If the slopes are **negative**, reverse the comparison
	if !slope1Positive {
		diff = -diff
	}

	if diff < 0 {
		return 1
	} else if diff > 0 {
		return -1
	}
	return 0
}

// XAtY computes the X-coordinate of the line segment at a given Y-coordinate.
// If the given Y-coordinate is out of the segment's bounds, the function returns math.NaN().
//
// Special Cases:
// - If the line is vertical, X remains constant for all Y values in range.
// - If Y is outside the range of the segment, the function returns math.NaN().
//
// Parameters:
//   - y: The Y-coordinate where the X-coordinate is needed.
//
// Returns:
//   - The X-coordinate at the given Y value, or NaN if Y is out of bounds.
func (l *LineSegment) XAtY(y float64) float64 {
	x, ok := l.xAtYScaled(int64(y * l.pointUpper.plane.scaleFactor))
	if !ok {
		return math.NaN()
	}
	return float64(x) / l.pointUpper.plane.scaleFactor
}

// xAtYScaled computes the integer-scaled X-coordinate at a given integer-scaled Y-coordinate.
// It ensures the given Y-value is within the segment's range before computing the result.
//
// Special Cases:
// - If the line is vertical, it returns a constant X value.
// - If Y is outside the valid range, the function returns (0, false).
//
// Parameters:
//   - y: The integer-scaled Y-coordinate where the X-coordinate is needed.
//
// Returns:
//   - The integer-scaled X-coordinate at the given Y value.
//   - A boolean indicating whether the computation was successful.
func (l *LineSegment) xAtYScaled(y int64) (int64, bool) {

	// Ensure y is within bounds
	if y < min(l.pointUpper.scaledY, l.pointLower.scaledY) || y > max(l.pointUpper.scaledY, l.pointLower.scaledY) {
		return 0, false
	}

	// Handle vertical line case: x is constant for all y values in range
	if l.pointUpper.scaledX == l.pointLower.scaledX {
		return l.pointUpper.scaledX, true
	}

	// Compute x using interpolation
	// Compute x using interpolation
	return int64(math.FMA(
		float64(y-l.pointUpper.scaledY),
		float64(l.pointLower.scaledX-l.pointUpper.scaledX)/float64(l.pointLower.scaledY-l.pointUpper.scaledY),
		float64(l.pointUpper.scaledX),
	)), true
}

// YAtX computes the Y-coordinate of the line segment at a given X-coordinate.
// If the given X-coordinate is out of the segment's bounds, the function returns math.NaN().
//
// Special Cases:
// - If the line is horizontal, Y remains constant for all X values in range.
// - If X is outside the range of the segment, the function returns math.NaN().
//
// Parameters:
//   - x: The X-coordinate where the Y-coordinate is needed.
//
// Returns:
//   - The Y-coordinate at the given X value, or NaN if X is out of bounds.
func (l *LineSegment) YAtX(x float64) float64 {
	y, ok := l.yAtXScaled(int64(x * l.pointUpper.plane.scaleFactor))
	if !ok {
		return math.NaN()
	}
	return float64(y) / l.pointUpper.plane.scaleFactor
}

// yAtXScaled computes the integer-scaled Y-coordinate at a given integer-scaled X-coordinate.
// It ensures the given X-value is within the segment's range before computing the result.
//
// Special Cases:
// - If the line is horizontal, it returns a constant Y value.
// - If X is outside the valid range, the function returns (0, false).
//
// Parameters:
//   - x: The integer-scaled X-coordinate where the Y-coordinate is needed.
//
// Returns:
//   - The integer-scaled Y-coordinate at the given X value.
//   - A boolean indicating whether the computation was successful.
func (l *LineSegment) yAtXScaled(x int64) (int64, bool) {

	// Ensure x is within bounds
	if x < min(l.pointUpper.scaledX, l.pointLower.scaledX) || x > max(l.pointUpper.scaledX, l.pointLower.scaledX) {
		return 0, false
	}

	// Handle horizontal line case: y is constant for all x values in range
	if l.pointUpper.scaledY == l.pointLower.scaledY {
		return l.pointUpper.scaledY, true
	}

	// Compute y using interpolation
	return int64(math.FMA(
		float64(x-l.pointUpper.scaledX),
		float64(l.pointLower.scaledY-l.pointUpper.scaledY)/float64(l.pointLower.scaledX-l.pointUpper.scaledX),
		float64(l.pointUpper.scaledY),
	)), true
}
