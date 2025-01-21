package geom2d

import (
	"fmt"
	"github.com/mikenye/geom2d/types"
	"math"
)

// BoundingBox calculates the axis-aligned bounding box (AABB) of the circle.
//
// The bounding box is the smallest rectangle, aligned with the coordinate axes, that completely encloses the circle.
// This is useful for collision detection, spatial partitioning, and other geometric operations.
//
// Returns:
//   - [Rectangle][T]: The axis-aligned bounding box that encloses the circle.
//
// Notes:
//   - The bounding box is a rectangle defined by the four corner points derived from the circle's center and radius.
func (c Circle[T]) BoundingBox() Rectangle[T] {
	return NewRectangle[T]([]Point[T]{
		NewPoint(c.center.x-c.radius, c.center.y-c.radius),
		NewPoint(c.center.x+c.radius, c.center.y-c.radius),
		NewPoint(c.center.x+c.radius, c.center.y+c.radius),
		NewPoint(c.center.x-c.radius, c.center.y+c.radius),
	})
}

// RelationshipToCircle determines the spatial relationship between two circles.
//
// This function evaluates the relationship between the current circle and another
// circle by comparing their center points and radii. The possible relationships include:
//   - [RelationshipEqual]: The circles are identical.
//   - [RelationshipContainedBy]: The current circle is completely contained within the other circle.
//   - [RelationshipContains]: The current circle completely contains the other circle.
//   - [RelationshipIntersection]: The circles overlap, including tangency.
//   - [RelationshipDisjoint]: The circles do not overlap.
//
// Parameters:
//   - other (Circle[T]): The circle to compare against the current circle.
//   - opts: A variadic slice of [Option] functions to customize the equality check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing the center coordinates
//     and radii of the circles, allowing for robust handling of floating-point precision errors.
//
// Returns:
//   - [Relationship]: A constant representing the relationship between the circles.
//
// Behavior:
//   - The function first checks for equality by comparing center points and radii.
//   - It then checks for containment by comparing the distance between centers and radii.
//   - Intersection is detected if the distance between centers is less than or equal to the sum of the radii.
//   - If no other relationship is found, the circles are considered disjoint.
func (c Circle[T]) RelationshipToCircle(other Circle[T], opts ...Option) Relationship {
	distanceBetweenCenters := c.center.DistanceToPoint(other.center, opts...)
	cFloat := c.AsFloat64()
	otherFloat := other.AsFloat64()

	// check for equality
	if c.Eq(other) {
		return RelationshipEqual
	}

	// check for c contained by other
	if distanceBetweenCenters+cFloat.radius < otherFloat.radius {
		return RelationshipContainedBy
	}

	// check for c contains other
	if distanceBetweenCenters+otherFloat.radius < cFloat.radius {
		return RelationshipContains
	}

	// check for intersection
	if distanceBetweenCenters <= cFloat.radius+otherFloat.radius {
		return RelationshipIntersection
	}

	return RelationshipDisjoint

}

// RelationshipToLineSegment determines the spatial relationship between the current Circle and a
// given [LineSegment].
//
// This function evaluates the relationship between the circle and the line segment,
// which can be one of the following:
//   - [RelationshipDisjoint]: The line segment lies entirely outside the circle.
//   - [RelationshipIntersection]: The line segment intersects the circle's boundary.
//   - [RelationshipContains]: The line segment is fully contained within the circle.
//
// This method internally calls [LineSegment.RelationshipToCircle], flipping the containment
// direction to align with the perspective of the circle.
//
// Parameters:
//   - l ([LineSegment][T]): The line segment to compare against the circle.
//   - opts: A variadic slice of [Option] functions to customize the equality check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing the center coordinates
//     and radii of the circles, allowing for robust handling of floating-point precision errors.
//
// Returns:
//   - [Relationship]: A constant representing the relationship between the circle and the line segment.
func (c Circle[T]) RelationshipToLineSegment(l LineSegment[T], opts ...Option) Relationship {
	return l.RelationshipToCircle(c, opts...).flipContainment()
}

// RelationshipToPoint determines the spatial relationship between the current Circle and a given [Point].
//
// This function evaluates the relationship between the circle and the point,
// which can be one of the following:
//   - [RelationshipDisjoint]: The point lies outside the circle.
//   - [RelationshipIntersection]: The point lies exactly on the circle's boundary.
//   - [RelationshipContains]: The point lies inside the circle.
//
// This method internally calls [Point.RelationshipToCircle], flipping the containment
// direction to align with the perspective of the circle.
//
// Parameters:
//   - p ([Point][T]): The point to compare against the circle.
//   - opts: A variadic slice of [Option] functions to customize the equality check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing the center coordinates
//     and radii of the circles, allowing for robust handling of floating-point precision errors.
//
// Returns:
//   - [Relationship]: A constant representing the relationship between the circle and the point.
func (c Circle[T]) RelationshipToPoint(p Point[T], opts ...Option) Relationship {
	return p.RelationshipToCircle(c, opts...).flipContainment()
}

// RelationshipToPolyTree determines the spatial relationships between the Circle and the polygons in the given [PolyTree].
//
// This method evaluates whether the circle intersects, contains, or is contained by each polygon in the [PolyTree].
// It uses a doubled representation of the Circle to align with the doubled points in the [PolyTree] for robust computations.
//
// Parameters:
//   - pt (*[PolyTree][T]): The [PolyTree] whose polygons will be compared to the Circle.
//   - opts: A variadic slice of [Option] functions to customize the equality check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for comparing the center coordinates
//     and radii of the circles, allowing for robust handling of floating-point precision errors.
//
// Behavior:
//   - For each polygon in the [PolyTree], the function iterates through its edges and checks the relationship to the Circle.
//   - Intersection is checked first; if any edge of the polygon intersects the circle, the relationship is marked as [RelationshipIntersection].
//   - If all edges of the polygon lie within the circle's radius, the relationship is marked as [RelationshipContains].
//   - If the circle's center lies within the polygon and its minimum distance to any edge is greater than its radius, the relationship is marked as [RelationshipContainedBy].
//   - If none of the above conditions are satisfied, the relationship is marked as [RelationshipDisjoint].
//
// Returns:
//   - map[*[PolyTree][T]][Relationship]: A map where each key is a pointer to a polygon in the [PolyTree], and the value is the relationship between the Circle and that polygon.
func (c Circle[T]) RelationshipToPolyTree(pt *PolyTree[T], opts ...Option) map[*PolyTree[T]]Relationship {
	output := make(map[*PolyTree[T]]Relationship, pt.Len())
	cDoubled := NewCircle(NewPoint(c.center.x*2, c.center.y*2), c.radius*2)
	cFloatDoubled := cDoubled.AsFloat64()

RelationshipToPolyTreeIterPolys:
	for poly := range pt.Nodes {
		minDistCircleCenterToEdge := math.MaxFloat64
		allEdgesWithinCircle := true

		for edge := range poly.contour.iterEdges {
			rel := cDoubled.RelationshipToLineSegment(edge, opts...)

			// Check for intersection
			if rel == RelationshipIntersection {
				output[poly] = RelationshipIntersection
				continue RelationshipToPolyTreeIterPolys
			}

			// Check if all edges are within the circle's radius
			distanceToEdge := cDoubled.center.DistanceToLineSegment(edge, opts...)
			minDistCircleCenterToEdge = min(minDistCircleCenterToEdge, distanceToEdge)
			if distanceToEdge > cFloatDoubled.radius {
				allEdgesWithinCircle = false
			}
		}

		// Check for containment: circle fully contains the polygon
		if allEdgesWithinCircle {
			output[poly] = RelationshipContains
			continue RelationshipToPolyTreeIterPolys
		}

		// Check for containment: polygon fully contains the circle
		if poly.contour.isPointInside(cDoubled.center) && minDistCircleCenterToEdge > cFloatDoubled.radius {
			output[poly] = RelationshipContainedBy
			continue RelationshipToPolyTreeIterPolys
		}

		// Default: no relationship found
		output[poly] = RelationshipDisjoint
	}

	return output
}

// RelationshipToRectangle determines the spatial relationship between the circle and the rectangle.
//
// This function evaluates whether the circle is:
//   - Disjoint from the rectangle (no overlap or touching),
//   - Intersects with the rectangle (crosses its boundary),
//   - Fully contains the rectangle (encloses it entirely),
//   - Fully contained by the rectangle (is completely inside the rectangle).
//
// Parameters:
//   - r ([Rectangle][T]): The rectangle to compare with the circle.
//   - opts: A variadic slice of [Option] functions to customize the behavior of the relationship check.
//     [WithEpsilon](epsilon float64): Specifies a tolerance for floating-point precision.
//
// Behavior:
//   - The function checks each edge of the rectangle for potential intersections with the circle.
//   - If all edges of the rectangle are fully contained within the circle, it returns [RelationshipContains].
//   - If the rectangle fully contains the circle (the circle’s center is inside the rectangle, and the circle
//     does not extend beyond any edge), it returns [RelationshipContainedBy].
//   - If none of these conditions are met, it determines whether the circle and rectangle are disjoint.
//
// Returns:
//
// [Relationship]: One of the following constants:
//   - [RelationshipDisjoint]: The circle and rectangle are entirely separate.
//   - [RelationshipIntersection]: The circle intersects with one or more edges of the rectangle.
//   - [RelationshipContains]: The circle completely encloses the rectangle.
//   - [RelationshipContainedBy]: The circle is fully contained within the rectangle.
func (c Circle[T]) RelationshipToRectangle(r Rectangle[T], opts ...Option) Relationship {
	cContainsR := true
	cFloat := c.AsFloat64()
	minDistCircleCenterToEdge := math.MaxFloat64
	for _, edge := range r.Edges() {
		rel := edge.RelationshipToCircle(c, opts...)

		// check for intersection
		if rel == RelationshipIntersection {
			return RelationshipIntersection
		}

		// check for containment
		if rel != RelationshipContainedBy {
			cContainsR = false
		}

		edgeFloat := edge.AsFloat64()
		minDistCircleCenterToEdge = min(minDistCircleCenterToEdge, cFloat.center.DistanceToLineSegment(edgeFloat, opts...))
	}

	// check c contain r
	if cContainsR {
		return RelationshipContains
	}

	// check r contains c
	if r.ContainsPoint(c.center) && minDistCircleCenterToEdge > cFloat.radius {
		return RelationshipContainedBy
	}

	return RelationshipDisjoint
}
