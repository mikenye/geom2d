package simple

import (
	"github.com/mikenye/geom2d/point"
	"github.com/mikenye/geom2d/types"
)

// Area2XSigned calculates twice the signed area of a simple polygon defined by a series
// of points. The function uses the [Shoelace Formula] (also known as Gauss's area formula)
// to compute the signed area efficiently.
//
// The input points are assumed to form a closed polygon, where the last
// point connects back to the first. If the input slice does not explicitly
// include the closing point, the algorithm still assumes the connection
// between the last and first points.
//
// The "signed" area means that the result is positive if the points are ordered in a
// counterclockwise direction (CCW) and negative if they are ordered in a clockwise
// direction (CW). This property is useful for determining the orientation of the polygon.
//
// Parameters:
//   - points ([]point.Point[T]): A varadic slice of [point.Point] instances defining the vertices of the polygon.
//     The points must represent a simple polygon (no self-intersections) and should be
//     ordered either clockwise or counterclockwise.
//
// Returns:
//   - T: Twice the signed area of the polygon. The value is positive if the points are
//     ordered counterclockwise, negative if clockwise, and zero if the polygon is degenerate
//     (e.g., collinear points or less than 3 vertices).
//
// Behavior:
//   - The function expects at least three points to form a valid polygon.
//   - If fewer than three points are provided, the function returns 0 (not a polygon).
//
// Notes:
//   - The points forming the polygon are assumed
//   - For concave polygons, the function still produces the correct signed area.
//   - The result is undefined for self-intersecting polygons (e.g., bow-tie shapes).
//   - This function is commonly used in computational geometry for tasks like checking
//     polygon orientation or validating polygon geometry.
//
// [Shoelace Formula]: https://en.wikipedia.org/wiki/Shoelace_formula
func Area2XSigned[T types.SignedNumber](points ...point.Point[T]) T {
	n := len(points)
	if n < 3 {
		// Not enough points to form a polygon
		return 0
	}

	var area T
	for i := 0; i < n; i++ {
		// Current point and the next point, wrapping around
		p1 := points[i]
		p2 := points[(i+1)%n]

		// Shoelace formula: x1*y2 - x2*y1
		area += (p1.X() * p2.Y()) - (p2.X() * p1.Y())
	}

	return area
}
