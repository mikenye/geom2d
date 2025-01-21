package linesegment

import (
	"fmt"
	"github.com/mikenye/geom2d/numeric"
	"github.com/mikenye/geom2d/options"
	"github.com/mikenye/geom2d/point"
	"github.com/mikenye/geom2d/types"
)

// LineSegment represents a line segment in a 2D space, defined by two endpoints, the start [point.Point] and end [point.Point].
//
// The generic type parameter T must satisfy the [types.SignedNumber] constraint, allowing the segment
// to use various numeric types such as int or float64 for its coordinates.
type LineSegment[T types.SignedNumber] struct {
	start point.Point[T]
	end   point.Point[T]
}

// New creates a new LineSegment with the specified start and end x and y coordinates.
//
// This constructor function initializes a [LineSegment] with the specified starting and ending points.
// The generic type parameter "T" must satisfy the [types.SignedNumber] constraint, allowing various numeric types
// (such as int or float64) to be used for the segment’s coordinates.
//
// Parameters:
//   - x1,y1 (T): The starting point of the LineSegment.
//   - x2,y2 (T): The ending point of the LineSegment.
//
// Returns:
//   - LineSegment[T] - A new line segment defined by the start and end points.
func New[T types.SignedNumber](x1, y1, x2, y2 T) LineSegment[T] {
	return LineSegment[T]{
		start: point.New[T](x1, y1),
		end:   point.New[T](x2, y2),
	}
}

// NewFromPoints creates a new LineSegment from two endpoints, a start [point.Point] and end [point.Point].
//
// This constructor function initializes a [LineSegment] with the specified starting and ending points.
// The generic type parameter "T" must satisfy the [types.SignedNumber] constraint, allowing various numeric types
// (such as int or float64) to be used for the segment’s coordinates.
//
// Parameters:
//   - start ([point.Point][T]): The starting [point.Point] of the LineSegment.
//   - end ([point.Point][T]): The ending [point.Point] of the LineSegment.
//
// Returns:
//   - LineSegment[T] - A new line segment defined by the start and end points.
func NewFromPoints[T types.SignedNumber](start, end point.Point[T]) LineSegment[T] {
	return LineSegment[T]{
		start: start,
		end:   end,
	}
}

// AsFloat32 converts the line segment to a LineSegment[float32] type.
//
// This function converts both endpoints of the LineSegment l to [Point][float32]
// values, creating a new line segment with floating-point coordinates.
// It is useful for precise calculations where floating-point accuracy is needed.
//
// Returns:
//   - LineSegment[float32] - The line segment with both endpoints converted to float32.
func (l LineSegment[T]) AsFloat32() LineSegment[float32] {
	return NewFromPoints(l.start.AsFloat32(), l.end.AsFloat32())
}

// AsFloat64 converts the line segment to a LineSegment[float64] type.
//
// This function converts both endpoints of the LineSegment l to [Point][float64]
// values, creating a new line segment with floating-point coordinates.
// It is useful for precise calculations where floating-point accuracy is needed.
//
// Returns:
//   - LineSegment[float64] - The line segment with both endpoints converted to float64.
func (l LineSegment[T]) AsFloat64() LineSegment[float64] {
	return NewFromPoints(l.start.AsFloat64(), l.end.AsFloat64())
}

// AsInt converts the line segment to a LineSegment[int] type.
//
// This function converts both endpoints of the line segment l to [Point][int]
// by truncating any decimal places. It is useful for converting a floating-point
// line segment to integer coordinates without rounding.
//
// Returns:
//   - LineSegment[int] - The line segment with both endpoints converted to integer coordinates by truncation.
func (l LineSegment[T]) AsInt() LineSegment[int] {
	return NewFromPoints(l.start.AsInt(), l.end.AsInt())
}

// AsIntRounded converts the line segment to a LineSegment[int] type with rounded coordinates.
//
// This function converts both endpoints of the line segment l to [Point][int]
// by rounding each coordinate to the nearest integer. It is useful when you need to
// approximate the segment’s position with integer coordinates while minimizing the
// rounding error.
//
// Returns:
//   - LineSegment[int] - The line segment with both endpoints converted to integer coordinates by rounding.
func (l LineSegment[T]) AsIntRounded() LineSegment[int] {
	return NewFromPoints(l.start.AsIntRounded(), l.end.AsIntRounded())
}

// Bresenham generates all the integer points along the LineSegment using
// Bresenham's line algorithm. It is an efficient way to rasterize a line
// in a grid or pixel-based system.
//
// The function is designed to be used with a for-loop, and thus takes a callback yield that processes each point.
// If the callback returns false at any point (if the calling for-loop is terminated, for example), the function
// halts further generation.
//
// Example use cases include:
// - Rendering lines in graphics applications.
// - Generating grid points for pathfinding.
//
// Parameters:
//   - yield (func([point.Point][int]) bool): A function that processes each generated point.
//     Returning false will stop further point generation.
//
// Note: This method requires integer-type coordinates for the line segment.
func (l LineSegment[int]) Bresenham(yield func(point.Point[int]) bool) {

	var x1, x2, y1, y2, dx, dy, sx, sy int

	x1 = l.start.X()
	x2 = l.end.X()
	y1 = l.start.Y()
	y2 = l.end.Y()

	// Calculate absolute deltas
	dx = numeric.Abs(x2 - x1)
	dy = numeric.Abs(y2 - y1)

	// Determine the direction of the increments
	sx = 1
	if x1 > x2 {
		sx = -1
	}
	sy = 1
	if y1 > y2 {
		sy = -1
	}

	// Bresenham's algorithm
	err := dx - dy
	for {
		if !yield(point.New(x1, y1)) {
			return
		}

		// Break the loop if we've reached the end point
		if x1 == x2 && y1 == y2 {
			return
		}

		// Calculate the error
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

// Center calculates the midpoint of the line segment, optionally applying an epsilon
// threshold to adjust the precision of the result.
//
// Parameters:
//   - opts: A variadic slice of [Option] functions to customize the behavior of the calculation.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for snapping near-integer or
//     near-zero results to cleaner values, improving robustness in floating-point calculations.
//
// Behavior:
//   - The midpoint is calculated by averaging the x and y coordinates of the start and end
//     points of the line segment.
//   - If [WithEpsilon] is provided, the resulting midpoint coordinates are adjusted such that
//     small deviations due to floating-point precision errors are corrected.
//
// Returns:
//   - [Point][float64]: The midpoint of the line segment as a point with floating-point coordinates,
//     optionally adjusted based on epsilon.
//
// Notes:
//   - Epsilon adjustment is particularly useful when working with floating-point coordinates
//     where minor imprecision could affect the midpoint calculation.
//   - The midpoint is always returned as [Point][float64], ensuring precision regardless of the
//     coordinate type of the original line segment.
func (l LineSegment[T]) Center(opts ...options.GeometryOptionsFunc) point.Point[float64] {
	// Apply geomOptions with defaults
	geoOpts := options.ApplyGeometryOptions(options.GeometryOptions{Epsilon: 0}, opts...)

	start := l.start.AsFloat64()
	end := l.end.AsFloat64()

	midX := (start.X() + end.X()) / 2
	midY := (start.Y() + end.Y()) / 2

	// Apply epsilon if specified
	if geoOpts.Epsilon > 0 {
		midX = types.SnapToEpsilon(midX, geoOpts.Epsilon)
		midY = types.SnapToEpsilon(midY, geoOpts.Epsilon)
	}

	return point.New[float64](midX, midY)
}

// End returns the ending [point.Point] of the LineSegment.
//
// This function provides access to the ending [point.Point] of the LineSegment l, typically representing
// the endpoint of the segment.
func (l LineSegment[T]) End() point.Point[T] {
	return l.end
}

// Eq checks if two line segments are equal by comparing their start and end points.
// Equality can be evaluated either exactly (default) or approximately using an epsilon threshold.
//
// Parameters:
//   - other (LineSegment[T]): The line segment to compare with the current line segment.
//   - opts: A variadic slice of [options.GeometryOptionsFunc] functions to customize the equality check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing the start and end
//     points of the line segments. If the absolute difference between the coordinates of
//     the points is less than epsilon, they are considered equal.
//
// Behavior:
//   - By default, the function performs an exact equality check, returning true only if
//     both the start and end points of l and other are identical.
//   - If [WithEpsilon] is provided, the function performs an approximate equality check,
//     considering the points equal if their coordinate differences are within the specified
//     epsilon threshold.
//
// Returns:
//   - bool: Returns true if both line segments have identical (or approximately equal with epsilon) start
//     and end points; otherwise, false.
//
// Notes:
//   - Approximate equality is useful when comparing line segments with floating-point coordinates,
//     where small precision errors might otherwise cause inequality.
//   - This function relies on the [Point.Eq] method, which supports epsilon adjustments.
func (l LineSegment[T]) Eq(other LineSegment[T], opts ...options.GeometryOptionsFunc) bool {
	return l.start.Eq(other.start, opts...) && l.end.Eq(other.end, opts...)
}

// Points returns the start [point.Point] and end [point.Point] of the LineSegment.
//
// Returns:
//   - start ([Point][T]): The start [point.Point] of the LineSegment.
//   - end ([Point][T]): The end [point.Point] of the LineSegment.
func (l LineSegment[T]) Points() (start, end point.Point[T]) {
	return l.start, l.end
}

// Rotate rotates the LineSegment around a given pivot [point.Point] by a specified angle in radians counterclockwise.
// Optionally, an epsilon threshold can be applied to adjust the precision of the resulting coordinates.
//
// Parameters:
//   - pivot ([point.Point][T]): The point around which to rotate the line segment.
//   - radians (float64): The rotation angle in radians.
//   - opts: A variadic slice of [Option] functions to customize the behavior of the rotation.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for snapping near-zero or near-integer
//     values in the resulting coordinates to cleaner values, improving robustness.
//
// Behavior:
//   - The function rotates the start and end points of the line segment around the given pivot
//     point by the specified angle using the [Point.Rotate] method.
//   - If [WithEpsilon] is provided, epsilon adjustments are applied to the rotated coordinates to
//     handle floating-point precision errors.
//
// Returns:
//   - LineSegment[float64]: A new line segment representing the rotated position, with floating-point coordinates.
//
// Notes:
//   - Epsilon adjustment is particularly useful when the rotation involves floating-point
//     calculations that could result in minor inaccuracies.
//   - The returned line segment always has float64 coordinates, ensuring precision regardless
//     of the coordinate type of the original line segment.
func (l LineSegment[T]) Rotate(pivot point.Point[T], radians float64, opts ...options.GeometryOptionsFunc) LineSegment[float64] {
	newStart := l.start.Rotate(pivot, radians, opts...)
	newEnd := l.end.Rotate(pivot, radians, opts...)
	return NewFromPoints(newStart, newEnd)
}

// Scale scales the line segment by a given factor from a specified reference point.
//
// Parameters:
//   - ref ([point.Point][T]): The reference point from which the scaling is applied. Using the origin
//     point (0, 0) scales the segment relative to the coordinate system's origin, while specifying
//     a custom reference point scales the segment relative to that point.
//   - factor ([T]): The scaling factor, where a value greater than 1 expands the segment,
//     and a value between 0 and 1 shrinks it.
//
// Behavior:
//   - The function scales both endpoints of the line segment relative to the specified
//     reference point using the [point.Point.Scale] method.
//   - The scaling operation preserves the relative orientation of the segment.
//
// Returns:
//   - [LineSegment][T]: A new line segment, scaled relative to the specified reference point.
//
// Notes:
//   - Scaling by a factor of 1 will return a line segment identical to the original.
//   - Negative scaling factors will mirror the segment across the reference point
//     and scale its length accordingly.
//   - If the user wishes to shrink the segment (factor < 1), we recommend ensuring
//     the line segment's type is floating-point to avoid precision loss. Use the [LineSegment.AsFloat64] method
//     to safely convert the segment to floating-point type before scaling.
func (l LineSegment[T]) Scale(ref point.Point[T], factor T) LineSegment[T] {
	return NewFromPoints(
		l.start.Scale(ref, factor),
		l.end.Scale(ref, factor),
	)
}

// Slope calculates the slope of the line segment.
//
// The slope is calculated as the change in y-coordinates (dy) divided by
// the change in x-coordinates (dx) of the line segment. This function
// returns the slope as a float64 and a boolean indicating whether the
// slope is defined.
//
// Returns:
//   - (float64, true): The calculated slope if the line segment is not vertical.
//   - (0, false): Indicates the slope is undefined (the line segment is vertical).
func (l LineSegment[T]) Slope() (float64, bool) {
	dx := float64(l.end.X() - l.start.X())
	dy := float64(l.end.Y() - l.start.Y())

	if dx == 0 {
		return 0, false // Vertical line, slope undefined
	}
	return dy / dx, true
}

// Start returns the starting point of the line segment.
//
// This function provides access to the starting point of the LineSegment l, typically representing
// the beginning of the segment.
func (l LineSegment[T]) Start() point.Point[T] {
	return l.start
}

// String returns a formatted string representation of the line segment for debugging and logging purposes.
//
// The string representation includes the coordinates of the start and end points in the format:
// "(x1, y1)(x2, y2)", where (x1, y1) are the coordinates of the start point,
// and (x2, y2) are the coordinates of the end point.
//
// Returns:
//   - string: A string representing the line segment's start and end coordinates.
func (l LineSegment[T]) String() string {
	return fmt.Sprintf("(%v,%v)(%v,%v)", l.start.X(), l.start.Y(), l.end.X(), l.end.Y())
}

// Translate moves the LineSegment by a specified vector.
//
// This method shifts the LineSegment's position in the 2D plane by translating
// both its start and end points by the given vector delta. The relative
// orientation and length of the LineSegment remain unchanged.
//
// Parameters:
//   - delta ([point.Point][T]): The vector by which to translate the line segment.
//
// Returns:
//   - [LineSegment][T]: A new LineSegment translated by the specified vector.
//
// Notes:
//   - Translating the line segment effectively adds the delta vector to both
//     the start and end points of the segment.
//   - This operation is equivalent to a uniform shift, maintaining the segment's
//     shape and size while moving it to a new position.
func (l LineSegment[T]) Translate(delta point.Point[T]) LineSegment[T] {
	return NewFromPoints(
		l.start.Translate(delta),
		l.end.Translate(delta),
	)
}
