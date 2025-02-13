package _old

import (
	"fmt"
	"github.com/mikenye/geom2d/types"
	"image"
	"math"
	"slices"
)

// RelationshipToPolyTree determines the spatial relationship between the current Point and each polygon in a [PolyTree].
//
// This method returns a map, where the keys are pointers to the polygons in the [PolyTree], and the values are
// [Relationship] constants indicating the relationship of the point to each polygon.
//
// Relationships:
//   - [RelationshipContainedBy]: The point is inside the polygon but not on its boundary.
//   - [RelationshipIntersection]: The point lies on an edge or vertex of the polygon.
//   - [RelationshipDisjoint]: The point lies entirely outside the polygon.
//
// Parameters:
//   - pt (*[PolyTree][T]): The [PolyTree] to analyze.
//   - opts: A variadic slice of [Option] functions to customize the behavior of the calculation.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing the point's location relative
//     to the polygons, improving robustness in floating-point calculations.
//
// Returns:
//   - map[*[PolyTree][T]][Relationship]: A map where each polygon in the [PolyTree] is associated with its relationship to the point.
//
// Behavior:
//
// For each polygon in the [PolyTree], the function checks whether the point is:
//
//   - Contained within the polygon.
//   - On an edge or vertex of the polygon.
//   - Outside the polygon entirely.
//
// The relationship for each polygon is stored in the output map.
func (p Point[T]) RelationshipToPolyTree(pt *PolyTree[T], opts ...Option) map[*PolyTree[T]]Relationship {
	pDoubled := NewPoint[T](p.x*2, p.y*2)
	output := make(map[*PolyTree[T]]Relationship, pt.Len())
PointRelationshipToPolyTreeIterPolys:
	for poly := range pt.Nodes {

		// check if point on edge/vertex
		for edge := range poly.contour.iterEdges {
			if pDoubled.RelationshipToLineSegment(edge, opts...) == RelationshipIntersection {
				output[poly] = RelationshipIntersection
				continue PointRelationshipToPolyTreeIterPolys
			}
		}

		// check if point is contained in poly
		if poly.contour.isPointInside(pDoubled) {
			output[poly] = RelationshipContainedBy
			continue PointRelationshipToPolyTreeIterPolys
		}

		// else, no relationship
		output[poly] = RelationshipDisjoint
	}
	return output
}

// RoundPointToEpsilon rounds the coordinates of a point to the nearest multiple of a given epsilon value.
// This function is useful for reducing numerical precision issues in geometric computations.
//
// Parameters:
// - point (Point[float64]): The input Point[float64] whose coordinates need to be rounded.
// - epsilon (float64): The precision value to which the coordinates should be rounded.
//
// Returns:
// - A new Point[float64] where each coordinate is rounded to the nearest multiple of epsilon.
//
// Notes:
// - The epsilon value must be greater than zero or the function will panic.
func RoundPointToEpsilon(point Point[float64], epsilon float64) Point[float64] {
	return NewPoint[float64](
		math.Round(point.x/epsilon)*epsilon,
		math.Round(point.y/epsilon)*epsilon,
	)
}
