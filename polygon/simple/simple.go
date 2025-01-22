package simple

import (
	"fmt"
	"github.com/mikenye/geom2d/linesegment"
	"github.com/mikenye/geom2d/options"
	"github.com/mikenye/geom2d/point"
	"github.com/mikenye/geom2d/types"
)

// IsWellFormed checks whether a given set of points defines a well-formed polygon.
// A polygon is considered well-formed if:
//
//  1. It has at least 3 points.
//  2. It has a non-zero area.
//  3. It does not contain any self-intersecting edges.
//
// Parameters:
//   - points ([]Point[T]): A slice of Point[T] representing the vertices of the polygon.
//
// Returns:
//   - A boolean indicating whether the polygon is well-formed.
//   - An error providing details if the polygon is not well-formed.
//
// Example Errors:
//   - "polygon must have at least 3 points"
//   - "polygon has zero area"
//   - "polygon has self-intersecting edges: edge1 and edge2"
//
// todo: should use sweep line for this
func IsWellFormed[T types.SignedNumber](points []point.Point[T], opts ...options.GeometryOptionsFunc) (bool, error) {
	// Check for minimum 3 points
	if len(points) < 3 {
		return false, fmt.Errorf("polygon must have at least 3 points")
	}

	// Check for non-zero area
	if Area2XSigned(points...) == 0 {
		return false, fmt.Errorf("polygon has zero area")
	}

	// Check for self-intersection
	segments := ToLineSegments(points...)

	// todo: options.GeometryOptionsFunc to allow user to change out linesegment.FindIntersectionsSlow (eg: sweepline)
	intersections := linesegment.FindIntersectionsSlow(segments, opts...)

	for _, intersection := range intersections {
		// Skip intersection if it is exactly at the endpoints of both line segments
		if intersection.IntersectionType == linesegment.IntersectionPoint {
			// Check if the intersection point is an endpoint of both segments
			pointOnSegment := func(p point.Point[float64], seg linesegment.LineSegment[T], opts ...options.GeometryOptionsFunc) bool {
				return p.Eq(seg.Start().AsFloat64(), opts...) || p.Eq(seg.End().AsFloat64(), opts...)
			}
			if pointOnSegment(intersection.IntersectionPoint, intersection.InputLineSegments[0], opts...) &&
				pointOnSegment(intersection.IntersectionPoint, intersection.InputLineSegments[1], opts...) {
				continue
			}
		}

		// If an intersection is found and not at a vertex, the polygon is invalid
		return false, fmt.Errorf("polygon has self-intersecting edges")
	}

	return true, nil // Polygon is well-formed
}

// ToLineSegments converts a set of points defining a polygon into a set of [linesegment.LineSegment]
// representing the edges of the polygon. Points are assumed to define a closed polygon,
// meaning the last point connects back to the first.
//
// Degenerate line segments (segments with zero length due to repeated points) are skipped.
//
// Parameters:
//   - points: A variadic slice of [point.Point][T] that defines the vertices of the polygon.
//
// Returns:
//   - []linesegment.LineSegment[T]: A slice of line segments representing the edges of the polygon.
//
// Behavior:
//   - If fewer than two points are provided, the function returns an empty slice.
//   - Degenerate line segments (zero-length segments) are excluded from the result.
func ToLineSegments[T types.SignedNumber](points ...point.Point[T]) []linesegment.LineSegment[T] {
	var segments []linesegment.LineSegment[T]
	n := len(points)

	if n < 2 {
		// Not enough points to form a line segment
		return segments
	}

	for i := 0; i < n; i++ {
		start := points[i]
		end := points[(i+1)%n] // Wrap around to close the polygon

		// Skip degenerate line segments
		if start.Eq(end) {
			continue
		}

		segments = append(segments, linesegment.NewFromPoints(start, end))
	}

	return segments
}
