package simple

import (
	"github.com/mikenye/geom2d/point"
	"github.com/mikenye/geom2d/types"
	"slices"
)

// MakeClockwise ensures that a slice of points representing a polygon is ordered in a clockwise direction.
//
// This function checks the orientation of the points based on the signed area of the polygon.
// If the signed area is positive, indicating a counterclockwise orientation, the function reverses
// the order of the points (in-place) to make them clockwise. If the points are already clockwise, no changes are made.
//
// Parameters:
//   - points ([][Point][T]): A slice of points representing the vertices of a polygon. The points are assumed
//     to form a closed loop or define a valid polygon.
//
// Behavior:
//   - Calculates the signed area of the polygon using [SignedArea2X].
//   - If the signed area is positive (counterclockwise orientation), reverses the order of the points.
//   - If the signed area is negative or zero (already clockwise or degenerate polygon), does nothing.
//
// Notes:
//   - This function modifies the input slice of points in place.
//   - A zero area polygon is considered already clockwise and is left unchanged.
//
// todo: update doc comments, unit test, example func
func MakeClockwise[T types.SignedNumber](points ...point.Point[T]) []point.Point[T] {
	if Orientation(points...) == types.PointsClockwise {
		return points
	}
	slices.Reverse(points)
	return points
}

// MakeCounterClockwise ensures that a slice of points representing a polygon is ordered in a counterclockwise direction.
//
// This function checks the orientation of the points based on the signed area of the polygon.
// If the signed area is negative, indicating a clockwise orientation, the function reverses
// the order of the points (in-place) to make them counterclockwise. If the points are already counterclockwise,
// no changes are made.
//
// Parameters:
//
//   - points ([][Point][T]): A slice of points representing the vertices of a polygon. The points are assumed
//     to form a closed loop or define a valid polygon.
//
// Behavior:
//
//   - Calculates the signed area of the polygon using [SignedArea2X].
//   - If the signed area is negative (clockwise orientation), reverses the order of the points.
//   - If the signed area is positive or zero (already counterclockwise or degenerate polygon), does nothing.
//
// Notes:
//
//   - This function modifies the input slice of points in place.
//   - A zero area polygon is considered already counterclockwise and is left unchanged.
//
// todo: update doc comments, unit test, example func
func MakeCounterClockwise[T types.SignedNumber](points ...point.Point[T]) []point.Point[T] {
	if Orientation(points...) == types.PointsCounterClockwise {
		return points
	}
	slices.Reverse(points)
	return points
}

// Orientation determines the winding (orientation) of a simple polygon defined
// by a sequence of points ([point.Point]).
//
// It calculates the signed area of the polygon using the [Shoelace Formula]
// to determine if the points are arranged in a counterclockwise (CCW) order,
// clockwise (CW) order, or are collinear (zero area). This function is widely
// used in computational geometry to classify the orientation of polygons
// and validate input for algorithms like convex hull or boolean operations.
//
// Parameters:
//   - points: A variadic slice of [point.Point] representing the vertices of
//     the polygon in sequential order. The last point is assumed to connect
//     back to the first to close the polygon.
//
// Returns a [types.PointOrientation]: A constant indicating the orientation of the polygon. This can be one of:
//   - [types.PointsCounterClockwise] if the polygon points form a CCW turn,
//   - [types.PointsClockwise] if the polygon points form a CW turn,
//   - [types.PointsCollinear] if the points are collinear (zero area).
//
// Notes:
//   - The input is assumed to form a closed polygon; the algorithm treats the
//     last point as connecting to the first.
//   - This function relies on [Area2XSigned] to compute the signed area
//     of the polygon.
//
// [Shoelace Formula]: https://en.wikipedia.org/wiki/Shoelace_formula
func Orientation[T types.SignedNumber](points ...point.Point[T]) types.PointOrientation {
	area2x := Area2XSigned(points...)
	switch {
	case area2x < 0:
		return types.PointsClockwise
	case area2x > 0:
		return types.PointsCounterClockwise
	default: // area2x == 0
		return types.PointsCollinear
	}
}
