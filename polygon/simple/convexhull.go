package simple

import (
	"github.com/mikenye/geom2d/point"
	"github.com/mikenye/geom2d/types"
	"slices"
)

// ConvexHull computes the [convex hull] of a finite set of points using the [Graham Scan] algorithm.
// The convex hull is the smallest convex polygon that encloses all points in the input set.
//
// Parameters:
//   - points ([Point][T]): A variable number of points for which the convex hull is to be computed.
//
// Returns:
//   - [][Point][T]: A slice of points representing the vertices of the convex hull in counterclockwise order.
//
// Behavior:
//   - If the input points have fewer than 3 points, the function returns the input unchanged.
//   - The function handles duplicate points gracefully, ensuring the output contains only the outer boundary.
//
// Note:
//   - Points are assumed to form a "closed" polygon, where the last point connects to the first.
//   - Collinear points along the boundary are included if they are farthest from the reference point.
//
// [Graham Scan]: https://en.wikipedia.org/wiki/Graham_scan
// [convex hull]: https://en.wikipedia.org/wiki/Convex_hull
func ConvexHull[T types.SignedNumber](points ...point.Point[T]) []point.Point[T] {
	// Handle edge cases where a convex hull cannot be formed.
	if len(points) < 3 {
		return points
	}

	// Step 1: Find the point with the lowest Y-coordinate (break ties by X-coordinate).
	_, lowestPoint := findLowestLeftestPoint(points...)

	// Step 2: Sort points by angle about the lowest point, using a stable sort.
	sortedPoints := make([]point.Point[T], len(points))
	copy(sortedPoints, points)
	orderPointsByAngleAboutLowestPoint(lowestPoint, sortedPoints)

	// Step 3: Initialize the convex hull stack.
	hull := make([]point.Point[T], 0, len(sortedPoints))
	hull = append(hull, sortedPoints[0], sortedPoints[1]) // Add the first two points.

	// Step 4: Process the sorted points.
	for i := 2; i < len(sortedPoints); i++ {
		for len(hull) > 1 {
			top := hull[len(hull)-1]
			nextToTop := hull[len(hull)-2]

			// Check the orientation of the top two points on the stack and the current point.
			if Orientation(nextToTop, top, sortedPoints[i]) != types.PointsClockwise {
				break // Stop popping if the turn is counterclockwise or collinear.
			}

			// Remove the top point as it causes a clockwise turn.
			hull = hull[:len(hull)-1]
		}
		hull = append(hull, sortedPoints[i]) // Add the current point to the hull.
	}

	return hull
}

// findLowestLeftestPoint identifies the point with the lowest y-coordinate from a given set of points.
// If multiple points share the lowest y-coordinate, it selects the point with the lowest x-coordinate among them.
//
// Parameters:
//   - points: A variadic list of [point.Point][T] instances from which the lowest leftmost point is determined.
//
// Returns:
//   - int: The index of the lowest leftmost point within the provided points.
//   - [point.Point][T]: The Point with the lowest y-coordinate (and lowest x-coordinate in case of ties).
//
// Example Usage:
//
//	points := []Point[int]{{3, 4}, {1, 5}, {1, 4}}
//	index, lowestPoint := findLowestLeftestPoint(points...)
//	// lowestPoint is Point[int]{1, 4}, and index is 2
func findLowestLeftestPoint[T types.SignedNumber](points ...point.Point[T]) (int, point.Point[T]) {

	lowestIndex := 0
	lowestPoint := points[0]

	for i := 1; i < len(points); i++ {
		current := points[i]
		if current.Y() < lowestPoint.Y() || (current.Y() == lowestPoint.Y() && current.X() < lowestPoint.X()) {
			lowestIndex = i
			lowestPoint = current
		}
	}
	return lowestIndex, lowestPoint
}

// orderPointsByAngleAboutLowestPoint sorts a slice of points by their angular order around a reference point, lowestPoint.
// This sorting is used in computational geometry algorithms, such as the Graham scan, to arrange points in a counterclockwise
// order around a pivot point. Collinear points are ordered by increasing distance from the lowestPoint.
//
// Parameters:
//   - lowestPoint: The reference Point from which angles are calculated for sorting. This point is usually the starting point in a convex hull algorithm.
//   - points: A slice of points to be sorted by their angle relative to the lowestPoint.
//
// Sorting Logic:
//   - The function uses the cross product of vectors from lowestPoint to each point to determine the angular order:
//   - If the cross product is positive, a is counterclockwise to b.
//   - If the cross product is negative, a is clockwise to b.
//   - If the cross product is zero, the points are collinear, so they are sorted by their distance to lowestPoint.
//
// Example Usage:
//
//	points := []Point[int]{{3, 4}, {1, 5}, {2, 2}}
//	lowestPoint := NewPoint(1, 2)
//	orderPointsByAngleAboutLowestPoint(lowestPoint, points)
//	// points are now sorted counterclockwise around lowestPoint, with collinear points ordered by distance.
func orderPointsByAngleAboutLowestPoint[T types.SignedNumber](lowestPoint point.Point[T], points []point.Point[T]) {
	slices.SortStableFunc(points, func(a point.Point[T], b point.Point[T]) int {

		// Ensure lowestPoint is always the first point
		switch {
		case a.Eq(lowestPoint):
			return -1
		case b.Eq(lowestPoint):
			return 1
		}

		// Calculate relative vectors from lowestPoint to start and end
		relativeA := a.Translate(lowestPoint.Negate())
		relativeB := b.Translate(lowestPoint.Negate())
		crossProduct := relativeA.CrossProduct(relativeB)

		// Use cross product to determine angular order
		switch {
		case crossProduct > 0:
			return -1 // start is counterclockwise to end
		case crossProduct < 0:
			return 1 // start is clockwise to end
		}

		// If cross product is zero, points are collinear; order by distance to lowestPoint
		distAtoLP := lowestPoint.DistanceSquaredToPoint(a)
		distBtoLP := lowestPoint.DistanceSquaredToPoint(b)

		// Sort closer points first
		switch {
		case distAtoLP < distBtoLP:
			return -1
		case distAtoLP > distBtoLP:
			return 1
		default:
			return 0
		}
	})
}
