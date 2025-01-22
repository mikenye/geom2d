package geom2d

import (
	"fmt"
	"github.com/mikenye/geom2d/types"
	"image"
)

// RelationshipToCircle determines the spatial relationship between a rectangle and a circle.
//
// This function evaluates whether the given [Circle] is:
//   - Disjoint from the rectangle ([RelationshipDisjoint])
//   - Intersecting the rectangle's boundary ([RelationshipIntersection])
//   - Fully contained within the rectangle ([RelationshipContains])
//
// The function delegates the relationship check to the circle's [Circle.RelationshipToRectangle] method and flips
// the containment perspective to represent the rectangle's relationship to the circle.
//
// Parameters:
//   - c ([Circle][T]): The circle to compare with the rectangle.
//   - opts: A variadic slice of [Option] functions to customize the equality check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing the center coordinates
//     and radii of the circles, allowing for robust handling of floating-point precision errors.
//
// Returns:
//   - [Relationship]: The relationship of the rectangle to the circle.
//
// Notes:
//   - The returned relationship reflects the rectangle's perspective.
//   - Use the [WithEpsilon] option to adjust for floating-point precision errors during the calculations.
func (r Rectangle[T]) RelationshipToCircle(c Circle[T], opts ...Option) Relationship {
	return c.RelationshipToRectangle(r, opts...).flipContainment()
}

// RelationshipToLineSegment determines the spatial relationship between a rectangle and a line segment.
//
// This function checks whether the given [LineSegment] is:
//   - Disjoint from the rectangle ([RelationshipDisjoint])
//   - Intersecting the rectangle's boundary ([RelationshipIntersection])
//   - Fully contained within the rectangle ([RelationshipContains])
//
// The relationship is determined by delegating the check to the line segment's [LineSegment.RelationshipToRectangle]
// method and then flipping the containment perspective to describe the rectangle's relationship to the line segment.
//
// Parameters:
//   - l ([LineSegment][T]): The line segment to compare with the rectangle.
//   - opts: A variadic slice of [Option] functions to customize the equality check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing the center coordinates
//     and radii of the circles, allowing for robust handling of floating-point precision errors.
//
// Returns:
//   - [Relationship]: The relationship of the rectangle to the line segment.
//
// Notes:
//   - The returned relationship is flipped to represent the rectangle's perspective.
//   - The behavior of the function can be customized using the [WithEpsilon] option to handle floating-point precision.
func (r Rectangle[T]) RelationshipToLineSegment(l LineSegment[T], opts ...Option) Relationship {
	return l.RelationshipToRectangle(r, opts...).flipContainment()
}

// RelationshipToPoint determines the spatial relationship between a rectangle and a point.
//
// This function checks whether the given [Point] is:
//   - Outside the rectangle ([RelationshipDisjoint])
//   - On the rectangle's edge or vertex ([RelationshipIntersection])
//   - Fully contained within the rectangle ([RelationshipContains])
//
// The relationship is determined by delegating the check to the point's [Point.RelationshipToRectangle] method
// and then flipping the containment perspective to describe the rectangle's relationship to the point.
//
// Parameters:
//   - p ([Point][T]): The point to compare with the rectangle.
//   - opts: A variadic slice of [Option] functions to customize the equality check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing the center coordinates
//     and radii of the circles, allowing for robust handling of floating-point precision errors.
//
// Returns:
//   - [Relationship]: The relationship of the rectangle to the point.
//
// Notes:
//   - The returned relationship is flipped to represent the rectangle's perspective.
//   - The behavior of the function can be customized using the [WithEpsilon] option to handle floating-point precision.
func (r Rectangle[T]) RelationshipToPoint(p Point[T], opts ...Option) Relationship {
	return p.RelationshipToRectangle(r, opts...).flipContainment()
}

// RelationshipToPolyTree determines the spatial relationship between a rectangle and a [PolyTree].
//
// This method evaluates how the calling [Rectangle] (r) relates to each polygon in the given [PolyTree] (pt).
// The relationships include intersection, containment, and disjoint.
//
// Parameters:
//   - pt: A pointer to the [PolyTree] to compare with the calling rectangle.
//   - opts: A variadic slice of [Option] functions to customize the behavior of the relationship check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing points and collinearity calculations,
//     allowing for robust handling of floating-point precision errors.
//
// Behavior:
//   - For each polygon in the [PolyTree], the function determines whether the rectangle intersects, contains,
//     is contained by, or is disjoint from the polygon.
//   - Intersection is determined by checking if any edge of the rectangle intersects any edge of the polygon.
//   - Containment is determined by checking whether all edges of one shape lie within the other.
//
// Returns:
//   - map[*[PolyTree][T]][Relationship]: A map where each key is a polygon in the [PolyTree], and the value is the
//     relationship between the rectangle and that polygon.
//
// Notes:
//   - The function assumes that both the rectangle and the polygons in the [PolyTree] are valid (e.g., non-degenerate).
//   - Epsilon adjustment is useful for floating-point coordinates, where small precision errors might otherwise
//     cause incorrect results.
func (r Rectangle[T]) RelationshipToPolyTree(pt *PolyTree[T], opts ...Option) map[*PolyTree[T]]Relationship {
	output := make(map[*PolyTree[T]]Relationship, pt.Len())

RelationshipToPolyTreeIterPolys:
	for poly := range pt.Nodes {
		rectangleContainsPoly := true
		polyContainsRectangle := true

		for edge := range poly.contour.iterEdges {

			edgeHalved := NewLineSegment[T](
				NewPoint[T](edge.start.x/2, edge.start.y/2),
				NewPoint[T](edge.end.x/2, edge.end.y/2),
			)

			rel := r.RelationshipToLineSegment(edgeHalved, opts...)

			// Check for intersection
			if rel == RelationshipIntersection {
				output[poly] = RelationshipIntersection
				continue RelationshipToPolyTreeIterPolys
			}

			// Check containment of poly by rectangle
			if !r.ContainsPoint(edgeHalved.start) || !r.ContainsPoint(edgeHalved.end) {
				rectangleContainsPoly = false
			}

			// Check containment of rectangle by poly
			for _, rectVertex := range r.Contour() {
				rectVertexDoubled := NewPoint[T](rectVertex.x*2, rectVertex.y*2)
				if !poly.contour.isPointInside(rectVertexDoubled) {
					polyContainsRectangle = false
				}
			}
		}

		// Determine containment relationships
		if rectangleContainsPoly {
			output[poly] = RelationshipContains
			continue RelationshipToPolyTreeIterPolys
		}
		if polyContainsRectangle {
			output[poly] = RelationshipContainedBy
			continue RelationshipToPolyTreeIterPolys
		}

		// If no stronger relationship is found, disjoint
		output[poly] = RelationshipDisjoint
	}

	return output
}

// RelationshipToRectangle determines the spatial relationship between two rectangles.
//
// This method evaluates the relationship between the calling [Rectangle] (r) and another [Rectangle] (other).
// It checks for equality, intersections, containment, and disjoint relationships. The function considers
// edge and vertex overlap to ensure accurate results.
//
// Parameters:
//   - other: The [Rectangle] to compare with the calling [Rectangle].
//   - opts: A variadic slice of [Option] functions to customize the behavior of the relationship check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing points and collinearity calculations,
//     allowing for robust handling of floating-point precision errors.
//
// Behavior:
//   - The function first checks if the two rectangles are equal.
//   - It then evaluates whether the rectangles intersect by checking all edge pairs.
//   - If no intersections are found, the function checks if one rectangle is fully contained within the other.
//   - If neither intersects nor is contained, the rectangles are considered disjoint.
//
// Returns:
//
// [Relationship]: A constant indicating the relationship between the two rectangles, which can be:
//   - [RelationshipEqual]: The rectangles are identical in position and size.
//   - [RelationshipIntersection]: The rectangles overlap but are not fully contained.
//   - [RelationshipContains]: The calling rectangle fully contains the other rectangle.
//   - [RelationshipContainedBy]: The calling rectangle is fully contained within the other rectangle.
//   - [RelationshipDisjoint]: The rectangles do not overlap or touch.
//
// Notes:
//   - The function assumes the input rectangles are valid (e.g., non-degenerate).
//   - Epsilon adjustment is useful for floating-point coordinates, where small precision errors might otherwise cause incorrect results.
func (r Rectangle[T]) RelationshipToRectangle(other Rectangle[T], opts ...Option) Relationship {

	// check for equality
	if r.Eq(other) {
		return RelationshipEqual
	}

	rInsideOther := true
	otherInsideR := true

	for _, rEdge := range r.Edges() {
		for _, otherEdge := range other.Edges() {

			// check for intersection
			rel := rEdge.RelationshipToLineSegment(otherEdge, opts...)
			if rel == RelationshipIntersection || rel == RelationshipEqual {
				return RelationshipIntersection
			}

			// check for containment
			if !(r.ContainsPoint(otherEdge.start) && r.ContainsPoint(otherEdge.end)) {
				otherInsideR = false
			}
			if !(other.ContainsPoint(rEdge.start) && other.ContainsPoint(rEdge.end)) {
				rInsideOther = false
			}
		}
	}

	// containment
	if otherInsideR {
		return RelationshipContains
	}
	if rInsideOther {
		return RelationshipContainedBy
	}

	return RelationshipDisjoint
}
