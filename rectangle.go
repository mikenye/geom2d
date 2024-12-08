package geom2d

import (
	"fmt"
	"image"
	"slices"
)

// Rectangle represents an axis-aligned rectangle defined by its four corners.
type Rectangle[T SignedNumber] struct {
	topLeft     Point[T]
	topRight    Point[T]
	bottomLeft  Point[T]
	bottomRight Point[T]
}

// Area calculates the area of the rectangle.
//
// Returns:
//   - T: The area of the rectangle, calculated as Width * Height.
func (r Rectangle[T]) Area() T {
	return r.Width() * r.Height()
}

// AsFloat converts the Rectangle's corner points to the float64 type, useful for higher-precision operations.
//
// Returns:
//   - Rectangle[float64]: A new Rectangle with float64 coordinates.
func (r Rectangle[T]) AsFloat() Rectangle[float64] {
	return NewRectangleByOppositeCorners(
		r.topLeft.AsFloat(),
		r.bottomRight.AsFloat(),
	)
}

// AsInt converts the Rectangle's corner points to the int type by truncating any decimal values.
// This method is useful for operations requiring integer coordinates.
//
// Returns:
//   - Rectangle[int]: A new Rectangle with integer coordinates, truncated from the original values.
func (r Rectangle[T]) AsInt() Rectangle[int] {
	return NewRectangleByOppositeCorners(
		r.topLeft.AsInt(),
		r.bottomRight.AsInt(),
	)
}

// AsIntRounded converts the Rectangle's corner points to the int type by rounding to the nearest integer.
// This method is useful for operations requiring integer coordinates with rounding.
//
// Returns:
//   - Rectangle[int]: A new Rectangle with integer coordinates, rounded from the original values.
func (r Rectangle[T]) AsIntRounded() Rectangle[int] {
	return NewRectangleByOppositeCorners(
		r.topLeft.AsIntRounded(),
		r.bottomRight.AsIntRounded(),
	)
}

// ContainsPoint checks if a given point lies within the rectangle.
// This includes points on the rectangle's boundary.
//
// Parameters:
//   - p: The point to check.
//
// Returns:
//   - bool: True if the point lies within or on the boundary of the rectangle; otherwise, false.
func (r Rectangle[T]) ContainsPoint(p Point[T]) bool {
	return p.x >= r.topLeft.x &&
		p.x <= r.bottomRight.x &&
		p.y <= r.topLeft.y &&
		p.y >= r.bottomRight.y
}

func (r Rectangle[T]) Edges() []LineSegment[T] {
	return []LineSegment[T]{
		NewLineSegment(r.bottomLeft, r.bottomRight),
		NewLineSegment(r.bottomRight, r.topRight),
		NewLineSegment(r.topRight, r.topLeft),
		NewLineSegment(r.topLeft, r.bottomLeft),
	}
}

// Height calculates the height of the rectangle.
//
// Returns:
//   - T: The height of the rectangle, calculated as the absolute difference between the y-coordinates of the top-left and bottom-right corners.
func (r Rectangle[T]) Height() T {
	height := r.bottomRight.y - r.topLeft.y
	if height < 0 {
		return -height // Ensure height is always positive
	}
	return height
}

// isLineSegmentOnEdge checks if the given line segment lies entirely on one of the rectangle's edges.
//
// Parameters:
//   - segment: The line segment to check.
//
// Returns:
//   - bool: True if the segment lies on one of the rectangle's edges; otherwise, false.
func (r Rectangle[T]) isLineSegmentOnEdge(segment LineSegment[T]) bool {
	relToStart := r.RelationshipToPoint(segment.start)
	relToEnd := r.RelationshipToPoint(segment.end)
	return (relToStart == RelationshipPointRectanglePointOnVertex || relToStart == RelationshipPointRectanglePointOnEdge) && (relToEnd == RelationshipPointRectanglePointOnVertex || relToEnd == RelationshipPointRectanglePointOnEdge)
}

// isLineSegmentOnEdgeWithEndTouchingVertex checks if the given line segment lies on an edge of the rectangle
// and one or both of its endpoints touch a vertex.
//
// Parameters:
//   - segment: The line segment to check.
//
// Returns:
//   - bool: True if the segment lies on an edge and touches a vertex; otherwise, false.
func (r Rectangle[T]) isLineSegmentOnEdgeWithEndTouchingVertex(segment LineSegment[T]) bool {
	vertices := []Point[T]{
		r.topLeft,
		NewPoint(r.topLeft.x, r.bottomRight.y), // Bottom-left
		NewPoint(r.bottomRight.x, r.topLeft.y), // Top-right
		r.bottomRight,
	}
	return r.isLineSegmentOnEdge(segment) &&
		(slices.Contains(vertices, segment.start) || slices.Contains(vertices, segment.end))
}

// todo: function below commented out as redundant (RelationshipToLineSegment)
//// LineSegmentEntersAndExits checks if a LineSegment enters the rectangle through one edge
//// and exits through another, indicating that it crosses through the rectangle.
////
//// Parameters:
////   - segment: The line segment to check.
////
//// Returns:
////   - bool: True if the segment enters through one edge and exits through another; otherwise, false.
////
//// Explanation of Logic:
////   - The function iterates over each edge of the rectangle and checks the relationship
////     between the segment and each edge using the RelationshipToLineSegment function.
////   - If the segment strictly intersects an edge (RelationshipLineSegmentLineSegmentIntersects) or one end of the segment
////     lies on an edge without the entire segment being on that edge (RelationshipLineSegmentLineSegmentConAB), it is counted
////     as an entry or exit point. We only test RelationshipLineSegmentLineSegmentConAB (and not RelationshipLineSegmentLineSegmentDonAB) to avoid double counting.
////   - If there is more than one intersection or endpoint contact with the rectangle's edges,
////     the segment is considered to "enter and exit," returning true.
////   - This approach prevents double-counting cases where the segment might lie on or touch an
////     edge without fully crossing into the rectangle.
//func (r Rectangle[T]) LineSegmentEntersAndExits(segment LineSegment[T]) bool {
//	entryCount := 0
//
//	edges := []LineSegment[T]{
//		NewLineSegment(r.topLeft, NewPoint(r.bottomRight.x, r.topLeft.y)),     // Top edge
//		NewLineSegment(NewPoint(r.bottomRight.x, r.topLeft.y), r.bottomRight), // Right edge
//		NewLineSegment(r.bottomRight, NewPoint(r.topLeft.x, r.bottomRight.y)), // Bottom edge
//		NewLineSegment(NewPoint(r.topLeft.x, r.bottomRight.y), r.topLeft),     // Left edge
//	}
//
//	var rel RelationshipLineSegmentLineSegment
//	for _, edge := range edges {
//
//		// Check for intersections or an endpoint lying on an edge without full overlap.
//		// We only test RelationshipLineSegmentLineSegmentConAB (and not RelationshipLineSegmentLineSegmentDonAB) to avoid double counting.
//		rel = segment.RelationshipToLineSegment(edge)
//		if rel == RelationshipLineSegmentLineSegmentIntersects || rel == RelationshipLineSegmentLineSegmentConAB {
//			entryCount++
//		}
//		if entryCount > 1 {
//			return true // Enters and exits
//		}
//	}
//	return false
//}

// todo: function below commented out as redundant (RelationshipToLineSegment)
//// LineSegmentIntersectsEdges checks if the given line segment intersects one or more edges of the rectangle.
////
//// Parameters:
////   - segment: The line segment to check.
////
//// Returns:
////   - bool: True if the segment intersects any edge; otherwise, false.
//func (r Rectangle[T]) LineSegmentIntersectsEdges(segment LineSegment[T]) bool {
//	edges := []LineSegment[T]{
//		NewLineSegment(r.topLeft, NewPoint(r.bottomRight.x, r.topLeft.y)),     // Top edge
//		NewLineSegment(NewPoint(r.bottomRight.x, r.topLeft.y), r.bottomRight), // Right edge
//		NewLineSegment(r.bottomRight, NewPoint(r.topLeft.x, r.bottomRight.y)), // Bottom edge
//		NewLineSegment(NewPoint(r.topLeft.x, r.bottomRight.y), r.topLeft),     // Left edge
//	}
//
//	var rel RelationshipLineSegmentLineSegment
//	for _, edge := range edges {
//		rel = segment.RelationshipToLineSegment(edge)
//		switch rel {
//		case RelationshipLineSegmentLineSegmentIntersects, RelationshipLineSegmentLineSegmentConAB, RelationshipLineSegmentLineSegmentDonAB:
//			return true
//		default:
//		}
//	}
//	return false
//}

// Perimeter calculates the perimeter of the rectangle.
//
// Returns:
//   - T: The perimeter of the rectangle, calculated as 2 * (Width + Height).
func (r Rectangle[T]) Perimeter() T {
	return 2 * (r.Width() + r.Height())
}

// Points returns the four corner points of the rectangle in the following order:
// top-left, top-right, bottom-right, and bottom-left.
//
// Returns:
//   - []Point[T]: A slice containing the four corner points of the rectangle.
func (r Rectangle[T]) Points() []Point[T] {
	return []Point[T]{
		r.topLeft,
		r.topRight,
		r.bottomRight,
		r.bottomLeft,
	}
}

func (r Rectangle[T]) RelationshipToCircle(c Circle[T], opts ...Option) RelationshipRectangleCircle {
	return c.RelationshipToRectangle(r, opts...)
}

// RelationshipToLineSegment determines the spatial relationship between a line segment and the rectangle.
// It evaluates how the segment interacts with the rectangle, considering endpoints, edges, and intersections.
//
// Parameters:
//   - segment: LineSegment[T] - The line segment to analyze.
//
// Returns:
//   - RelationshipLineSegmentRectangle: An enum value describing the relationship of the segment
//     to the rectangle. Possible relationships include:
//   - RelationshipLineSegmentRectangleMiss: The segment lies entirely outside the rectangle.
//   - RelationshipLineSegmentRectangleContainedByRectangle: The segment lies entirely within the rectangle.
//   - RelationshipLineSegmentRectangleEdgeCollinear: The segment lies entirely on one of the rectangle’s edges.
//   - RelationshipLineSegmentRectangleEdgeCollinearTouchingVertex: The segment lies on an edge, and one or both endpoints touch a vertex.
//   - RelationshipLineSegmentRectangleEndTouchesEdgeInternally: One endpoint of the segment lies on an edge, and the other is inside the rectangle.
//   - RelationshipLineSegmentRectangleEndTouchesVertexInternally: One endpoint of the segment lies on a vertex, and the other is inside the rectangle.
//   - RelationshipLineSegmentRectangleEndTouchesEdgeExternally: One endpoint lies on an edge, and the other is outside the rectangle.
//   - RelationshipLineSegmentRectangleEndTouchesVertexExternally: One endpoint lies on a vertex, and the other is outside the rectangle.
//   - RelationshipLineSegmentRectangleIntersects: The segment intersects the rectangle through one or more edges.
//   - RelationshipLineSegmentRectangleEntersAndExits: The segment passes through the rectangle, entering and exiting via different edges.
//   - RelationshipLineSegmentRectangleMiss: The segment does not interact with the rectangle in any of the above ways.
//
// Behavior:
//   - The function first evaluates the relationship of each endpoint of the segment to the rectangle using
//     `RelationshipToPoint`.
//   - If the segment is degenerate (both endpoints are the same), it directly determines the relationship based on
//     the point’s position.
//   - For non-degenerate segments, it checks interactions with the rectangle’s edges, using a combination of
//     `LineSegment` relationships and endpoint relationships.
//   - The function examines conditions such as:
//   - Full containment
//   - Touching or lying on edges
//   - Intersection at one or more points
//   - Entry and exit through different edges
//
// Example Usage:
//
//	rect := NewRectangleByPoints(NewPoint(0, 10), NewPoint(10, 0))
//	segment := NewLineSegment(NewPoint(5, 15), NewPoint(5, -5))
//
//	relationship := rect.RelationshipToLineSegment(segment)
//	// Returns RelationshipLineSegmentRectangleEntersAndExits because the segment passes through the rectangle.
//
// Notes:
//   - The rectangle is treated as axis-aligned.
//   - Edge and vertex relationships are evaluated based on geometric precision.
//   - Epsilon adjustments are not currently applied, so the function relies on exact evaluations of segment and point relationships.
func (r Rectangle[T]) RelationshipToLineSegment(segment LineSegment[T]) RelationshipLineSegmentRectangle {
	// Determine relationships of each endpoint of the segment to the rectangle
	startRelationship := r.RelationshipToPoint(segment.start)
	endRelationship := r.RelationshipToPoint(segment.end)

	// Handle degenerate segment (start and end points are the same)
	if segment.start == segment.end {
		switch startRelationship {
		case RelationshipPointRectangleContainedByRectangle:
			return RelationshipLineSegmentRectangleContainedByRectangle
		case RelationshipPointRectanglePointOnVertex:
			return RelationshipLineSegmentRectangleEndTouchesVertexExternally
		case RelationshipPointRectanglePointOnEdge:
			return RelationshipLineSegmentRectangleEndTouchesEdgeExternally
		default:
			return RelationshipLineSegmentRectangleMiss
		}
	}

	// Define rectangle edges
	edges := []LineSegment[T]{
		NewLineSegment(r.topLeft, NewPoint(r.bottomRight.x, r.topLeft.y)),     // Top edge
		NewLineSegment(NewPoint(r.bottomRight.x, r.topLeft.y), r.bottomRight), // Right edge
		NewLineSegment(r.bottomRight, NewPoint(r.topLeft.x, r.bottomRight.y)), // Bottom edge
		NewLineSegment(NewPoint(r.topLeft.x, r.bottomRight.y), r.topLeft),     // Left edge
	}

	// Identify relationships with each edge
	edgeRelationships := make([]RelationshipLineSegmentLineSegment, len(edges))
	for i, edge := range edges {
		edgeRelationships[i] = edge.RelationshipToLineSegment(segment)
	}

	// Check if segment enters and exists
	if countOccurrencesInSlice(edgeRelationships, RelationshipLineSegmentLineSegmentAonCD) >= 1 &&
		countOccurrencesInSlice(edgeRelationships, RelationshipLineSegmentLineSegmentBonCD) >= 1 &&
		startRelationship == RelationshipPointRectangleMiss && endRelationship == RelationshipPointRectangleMiss {
		return RelationshipLineSegmentRectangleEntersAndExits
	}

	// Check if segment fully inside
	if adjacentInSlice(edgeRelationships, RelationshipLineSegmentLineSegmentAeqC, RelationshipLineSegmentLineSegmentBeqD) &&
		adjacentInSlice(edgeRelationships, RelationshipLineSegmentLineSegmentAeqD, RelationshipLineSegmentLineSegmentBeqC) &&
		startRelationship == RelationshipPointRectanglePointOnVertex && endRelationship == RelationshipPointRectanglePointOnVertex {
		return RelationshipLineSegmentRectangleEndTouchesVertexInternally
	}

	// Check if segment is inside, with one end on an edge
	if countOccurrencesInSlice(edgeRelationships, RelationshipLineSegmentLineSegmentMiss) == 3 &&
		(countOccurrencesInSlice(edgeRelationships, RelationshipLineSegmentLineSegmentDonAB) == 1 ||
			countOccurrencesInSlice(edgeRelationships, RelationshipLineSegmentLineSegmentConAB) == 1) &&
		((startRelationship == RelationshipPointRectanglePointOnEdge && endRelationship == RelationshipPointRectangleContainedByRectangle) ||
			(startRelationship == RelationshipPointRectangleContainedByRectangle && endRelationship == RelationshipPointRectanglePointOnEdge)) {
		return RelationshipLineSegmentRectangleEndTouchesEdgeInternally
	}

	// Check if segment is fully outside
	if countOccurrencesInSlice(edgeRelationships, RelationshipLineSegmentLineSegmentMiss) == len(edges) &&
		startRelationship == RelationshipPointRectangleMiss && endRelationship == RelationshipPointRectangleMiss {
		return RelationshipLineSegmentRectangleMiss
	}

	// Check if the segment lies entirely on an edge
	if countOccurrencesInSlice(edgeRelationships, RelationshipLineSegmentLineSegmentCollinearEqual) == 1 {
		if r.isLineSegmentOnEdgeWithEndTouchingVertex(segment) {
			return RelationshipLineSegmentRectangleEdgeCollinearTouchingVertex
		}
		return RelationshipLineSegmentRectangleEdgeCollinear
	}
	if countOccurrencesInSlice(edgeRelationships, RelationshipLineSegmentLineSegmentCollinearCDinAB) == 1 &&
		countOccurrencesInSlice(edgeRelationships, RelationshipLineSegmentLineSegmentMiss) == 3 &&
		startRelationship == RelationshipPointRectanglePointOnEdge && endRelationship == RelationshipPointRectanglePointOnEdge {
		return RelationshipLineSegmentRectangleEdgeCollinear
	}

	// Check if the segment intersects the rectangle through one or more edges
	intersectionCount := countOccurrencesInSlice(edgeRelationships, RelationshipLineSegmentLineSegmentIntersects)
	if intersectionCount == 1 {
		return RelationshipLineSegmentRectangleIntersects
	} else if intersectionCount > 1 {
		return RelationshipLineSegmentRectangleEntersAndExits
	}

	// Check if one endpoint is on an edge and the other is outside
	if countOccurrencesInSlice(edgeRelationships, RelationshipLineSegmentLineSegmentConAB) == 1 && endRelationship == RelationshipPointRectangleMiss {
		return RelationshipLineSegmentRectangleEndTouchesEdgeExternally
	} else if countOccurrencesInSlice(edgeRelationships, RelationshipLineSegmentLineSegmentDonAB) == 1 && startRelationship == RelationshipPointRectangleMiss {
		return RelationshipLineSegmentRectangleEndTouchesEdgeExternally
	}

	// Check if one endpoint is on a vertex and the other is inside or outside
	if countOccurrencesInSlice(edgeRelationships, RelationshipLineSegmentLineSegmentAeqC) == 1 && adjacentInSlice(edgeRelationships, RelationshipLineSegmentLineSegmentAeqC, RelationshipLineSegmentLineSegmentBeqC) {
		if endRelationship == RelationshipPointRectangleContainedByRectangle {
			return RelationshipLineSegmentRectangleEndTouchesVertexInternally
		}
		return RelationshipLineSegmentRectangleEndTouchesVertexExternally
	} else if countOccurrencesInSlice(edgeRelationships, RelationshipLineSegmentLineSegmentAeqD) == 1 && adjacentInSlice(edgeRelationships, RelationshipLineSegmentLineSegmentAeqD, RelationshipLineSegmentLineSegmentBeqD) {
		if startRelationship == RelationshipPointRectangleContainedByRectangle {
			return RelationshipLineSegmentRectangleEndTouchesVertexInternally
		}
		return RelationshipLineSegmentRectangleEndTouchesVertexExternally
	}
	if countOccurrencesInSlice(edgeRelationships, RelationshipLineSegmentLineSegmentMiss) == 2 &&
		adjacentInSlice(edgeRelationships, RelationshipLineSegmentLineSegmentAonCD, RelationshipLineSegmentLineSegmentBonCD) &&
		((startRelationship == RelationshipPointRectangleMiss && endRelationship == RelationshipPointRectangleContainedByRectangle) ||
			(startRelationship == RelationshipPointRectangleContainedByRectangle && endRelationship == RelationshipPointRectangleMiss)) {
		return RelationshipLineSegmentRectangleIntersects
	}

	if countOccurrencesInSlice(edgeRelationships, RelationshipLineSegmentLineSegmentMiss) == 2 &&
		(adjacentInSlice(edgeRelationships, RelationshipLineSegmentLineSegmentAonCD, RelationshipLineSegmentLineSegmentCollinearBonCD) ||
			adjacentInSlice(edgeRelationships, RelationshipLineSegmentLineSegmentBonCD, RelationshipLineSegmentLineSegmentCollinearAonCD)) &&
		((startRelationship == RelationshipPointRectanglePointOnEdge && endRelationship == RelationshipPointRectangleMiss) ||
			(startRelationship == RelationshipPointRectangleMiss && endRelationship == RelationshipPointRectanglePointOnEdge)) {
		return RelationshipLineSegmentRectangleEndTouchesEdgeExternally
	}

	// If both endpoints are inside
	if startRelationship == RelationshipPointRectangleContainedByRectangle && endRelationship == RelationshipPointRectangleContainedByRectangle &&
		countOccurrencesInSlice(edgeRelationships, RelationshipLineSegmentLineSegmentMiss) == len(edges) {
		return RelationshipLineSegmentRectangleContainedByRectangle
	}

	return RelationshipLineSegmentRectangleMiss
}

// RelationshipToPoint determines the spatial relationship of a point relative to the rectangle.
// The relationship can be one of the following:
//   - RelationshipPointRectangleContainedByRectangle: The point lies strictly within the rectangle.
//   - RelationshipPointRectangleMiss: The point lies outside the rectangle.
//   - RelationshipPointRectanglePointOnVertex: The point coincides with one of the rectangle’s vertices.
//   - RelationshipPointRectanglePointOnEdge: The point lies on one of the rectangle’s edges but not on a vertex.
//
// Parameters:
//   - p: Point[T] - The point to analyze.
//
// Returns:
//   - RelationshipPointRectangle: An enum value indicating the relationship of the point to the rectangle.
//
// Behavior:
//   - The function first checks if the point lies strictly inside the rectangle by comparing its coordinates
//     to the bounds defined by the `topLeft` and `bottomRight` corners.
//   - If not inside, it evaluates whether the point coincides with one of the rectangle’s vertices.
//   - If the point is neither inside nor on a vertex, it checks if the point lies on any of the rectangle’s edges.
//   - If none of the above conditions are met, the point is classified as outside the rectangle.
//
// Example Usage:
//
//	rect := NewRectangleByPoints(NewPoint(0, 10), NewPoint(10, 0))
//
//	// Point inside the rectangle
//	pointInside := NewPoint(5, 5)
//	relationshipInside := rect.RelationshipToPoint(pointInside) // Returns RelationshipPointRectangleContainedByRectangle
//
//	// Point on the edge of the rectangle
//	pointOnEdge := NewPoint(0, 5)
//	relationshipOnEdge := rect.RelationshipToPoint(pointOnEdge) // Returns RelationshipPointRectanglePointOnEdge
//
//	// Point on a vertex of the rectangle
//	pointOnVertex := NewPoint(0, 10)
//	relationshipOnVertex := rect.RelationshipToPoint(pointOnVertex) // Returns RelationshipPointRectanglePointOnVertex
//
//	// Point outside the rectangle
//	pointOutside := NewPoint(-5, 15)
//	relationshipOutside := rect.RelationshipToPoint(pointOutside) // Returns RelationshipPointRectangleMiss
//
// Notes:
//   - This function assumes that the rectangle is axis-aligned, meaning its sides are parallel to the coordinate axes.
//   - Precision issues are not accounted for; all comparisons are performed using exact equality.
func (r Rectangle[T]) RelationshipToPoint(p Point[T]) RelationshipPointRectangle {
	// Check if the point is strictly inside
	if p.x > r.topLeft.x && p.x < r.bottomRight.x &&
		p.y < r.topLeft.y && p.y > r.bottomRight.y {
		return RelationshipPointRectangleContainedByRectangle
	}

	// Check if the point is on a vertex
	if (p == r.topLeft) ||
		(p == r.bottomRight) ||
		(p == NewPoint(r.topLeft.x, r.bottomRight.y)) || // Bottom-left vertex
		(p == NewPoint(r.bottomRight.x, r.topLeft.y)) { // Top-right vertex
		return RelationshipPointRectanglePointOnVertex
	}

	// Check if the point is on an edge
	onLeft := p.x == r.topLeft.x && p.y < r.topLeft.y && p.y > r.bottomRight.y
	onRight := p.x == r.bottomRight.x && p.y < r.topLeft.y && p.y > r.bottomRight.y
	onTop := p.y == r.topLeft.y && p.x > r.topLeft.x && p.x < r.bottomRight.x
	onBottom := p.y == r.bottomRight.y && p.x > r.topLeft.x && p.x < r.bottomRight.x
	if onLeft || onRight || onTop || onBottom {
		return RelationshipPointRectanglePointOnEdge
	}

	// Otherwise, the point is outside
	return RelationshipPointRectangleMiss
}

// RelationshipToRectangle determines the spatial relationship between two rectangles.
//
// Parameters:
//   - other: Rectangle[T] - The rectangle to compare with.
//
// Returns:
//   - RelationshipRectangleRectangle: The relationship between the two rectangles.
//
// Possible Relationships:
//   - RelationshipRectangleRectangleMiss: Rectangles are disjoint, with no overlap or touching.
//   - RelationshipRectangleRectangleSharedEdge: Rectangles share a complete edge but do not overlap.
//   - RelationshipRectangleRectangleSharedVertex: Rectangles share a single vertex but do not overlap.
//   - RelationshipRectangleRectangleIntersection: Rectangles overlap but neither is fully contained in the other.
//   - RelationshipRectangleRectangleContained: One rectangle is fully contained within the other without touching edges.
//   - RelationshipRectangleRectangleContainedTouching: One rectangle is fully contained within the other and touches edges.
//   - RelationshipRectangleRectangleEqual: Rectangles are identical in position and size.
func (r Rectangle[T]) RelationshipToRectangle(other Rectangle[T], opts ...Option) RelationshipRectangleRectangle {
	// Apply options if necessary
	// options := applyOptions(geomOptions{epsilon: 0}, opts...)

	// Check if the rectangles are disjoint
	disjointHorizontally := r.bottomRight.x < other.topLeft.x || other.bottomRight.x < r.topLeft.x
	disjointVertically := r.topLeft.y < other.bottomRight.y || other.topLeft.y < r.bottomRight.y

	if disjointHorizontally && disjointVertically {
		return RelationshipRectangleRectangleMiss
	}

	// Check if the rectangles are equal
	if r.topLeft.Eq(other.topLeft, opts...) && r.bottomRight.Eq(other.bottomRight, opts...) {
		return RelationshipRectangleRectangleEqual
	}

	// Check for containment (strictly inside or touching edges)
	rectCorners := []Point[T]{r.topLeft, r.topRight, r.bottomLeft, r.bottomRight}
	otherCorners := []Point[T]{other.topLeft, other.topRight, other.bottomLeft, other.bottomRight}

	allOtherCornersInside := true
	for _, corner := range otherCorners {
		if r.RelationshipToPoint(corner) == RelationshipPointRectangleMiss {
			allOtherCornersInside = false
			break
		}
	}

	if allOtherCornersInside {
		allStrictlyInside := true
		for _, corner := range otherCorners {
			if r.RelationshipToPoint(corner) != RelationshipPointRectangleContainedByRectangle {
				allStrictlyInside = false
				break
			}
		}
		if allStrictlyInside {
			return RelationshipRectangleRectangleContained
		}
		return RelationshipRectangleRectangleContainedTouching
	}

	allRectCornersInside := true
	for _, corner := range rectCorners {
		if other.RelationshipToPoint(corner) == RelationshipPointRectangleMiss {
			allRectCornersInside = false
			break
		}
	}

	if allRectCornersInside {
		allStrictlyInside := true
		for _, corner := range rectCorners {
			if other.RelationshipToPoint(corner) != RelationshipPointRectangleContainedByRectangle {
				allStrictlyInside = false
				break
			}
		}
		if allStrictlyInside {
			return RelationshipRectangleRectangleContained
		}
		return RelationshipRectangleRectangleContainedTouching
	}

	// Check for edge touching
	rectEdges := r.Edges()
	otherEdges := other.Edges()
	for _, edge := range rectEdges {
		for _, otherEdge := range otherEdges {
			if edge.RelationshipToLineSegment(otherEdge) == RelationshipLineSegmentLineSegmentCollinearEqual {
				return RelationshipRectangleRectangleSharedEdge
			}
		}
	}

	// Check for vertex touching
	vertexTouch := false
	for _, rectVertex := range rectCorners {
		for _, otherVertex := range otherCorners {
			if rectVertex.Eq(otherVertex, opts...) {
				vertexTouch = true
				break
			}
		}
	}
	if vertexTouch {
		return RelationshipRectangleRectangleSharedVertex
	}

	// If none of the above, the rectangles intersect
	return RelationshipRectangleRectangleIntersection
}

// Scale scales the rectangle relative to a specified reference point by a given scalar factor.
//
// Each corner of the rectangle is scaled relative to the reference point using the provided factor.
// The resulting rectangle maintains its axis-aligned orientation.
//
// Parameters:
//   - ref: Point[T] - The reference point relative to which the rectangle is scaled.
//   - k: T - The scaling factor. A value > 1 enlarges the rectangle; < 1 shrinks it.
//
// Returns:
//   - Rectangle[T]: A new rectangle with corners scaled relative to the reference point.
//
// Example Usage:
//
//	rect := NewRectangleByPoints(NewPoint(0, 10), NewPoint(10, 0))
//	ref := NewPoint(5, 5) // Center of the rectangle
//
//	scaledRect := rect.ScaleFrom(ref, 2)
//	// The rectangle is scaled relative to (5, 5), doubling its size.
//
// Notes:
//   - The function delegates the scaling of each corner to the `ScaleFrom` method of the `Point` type.
//   - The rectangle remains axis-aligned after scaling.
//   - If the scaling factor `k` is 1, the rectangle remains unchanged.
func (r Rectangle[T]) Scale(ref Point[T], k T) Rectangle[T] {
	return NewRectangle([]Point[T]{
		r.topLeft.Scale(ref, k),
		r.topRight.Scale(ref, k),
		r.bottomLeft.Scale(ref, k),
		r.bottomRight.Scale(ref, k),
	})
}

// ScaleHeight scales the height of the rectangle by a scalar value.
//
// Parameters:
//   - factor: The scaling factor for the height. A value > 1 enlarges the height; < 1 shrinks it.
//
// Returns:
//   - Rectangle[float64]: A new Rectangle with the height scaled by the given factor.
func (r Rectangle[T]) ScaleHeight(factor float64) Rectangle[float64] {
	topLeft := r.topLeft.AsFloat()
	height := float64(r.Height()) * factor
	return NewRectangleByOppositeCorners(
		topLeft,
		NewPoint(r.bottomRight.AsFloat().x, topLeft.y+height),
	)
}

// ScaleWidth scales the width of the rectangle by a scalar value.
//
// Parameters:
//   - factor: The scaling factor for the width. A value > 1 enlarges the width; < 1 shrinks it.
//
// Returns:
//   - Rectangle[float64]: A new Rectangle with the width scaled by the given factor.
func (r Rectangle[T]) ScaleWidth(factor float64) Rectangle[float64] {
	topLeft := r.topLeft.AsFloat()
	width := float64(r.Width()) * factor
	return NewRectangleByOppositeCorners(
		topLeft,
		NewPoint(topLeft.x+width, r.bottomRight.AsFloat().y),
	)
}

func (r Rectangle[T]) String() string {
	return fmt.Sprintf("Rectangle[(%v, %v), (%v, %v), (%v, %v), (%v, %v)]", r.bottomLeft.x, r.bottomLeft.y, r.bottomRight.x, r.bottomRight.y, r.topRight.x, r.topRight.y, r.topLeft.x, r.topLeft.y)
}

// ToImageRect converts the Rectangle[int] to an image.Rectangle.
// This method is only available for Rectangle[int] as image.Rectangle requires integer coordinates.
//
// Returns:
//   - image.Rectangle: A new image.Rectangle with coordinates matching the Rectangle.
//
// Example Usage:
//
//	rect := NewRectangleByPoints(NewPoint(0, 0), NewPoint(100, 200))
//	imgRect := rect.ToImageRect() // image.Rect(0, 0, 100, 200)
func (r Rectangle[int]) ToImageRect() image.Rectangle {
	topLeft := r.topLeft.AsInt()
	bottomRight := r.bottomRight.AsInt()
	return image.Rect(topLeft.x, topLeft.y, bottomRight.x, bottomRight.y)
}

// Translate moves the rectangle by a specified vector.
//
// This method shifts the rectangle's position in the 2D plane by translating
// both its top-left and bottom-right corners by the given vector `p`. The
// dimensions of the rectangle remain unchanged.
//
// Parameters:
//   - p: Point[T] - The vector by which to translate the rectangle.
//
// Returns:
//   - Rectangle[T]: A new Rectangle translated by the specified vector.
//
// Example Usage:
//
//	rectangle := NewRectangleByOppositeCorners(NewPoint(1, 1), NewPoint(4, 4))
//	translationVector := NewPoint(2, 3)
//	translatedRectangle := rectangle.Translate(translationVector)
//	// translatedRectangle has its top-left corner at (3, 4)
//	// and bottom-right corner at (6, 7), preserving its dimensions.
func (r Rectangle[T]) Translate(p Point[T]) Rectangle[T] {
	return NewRectangleByOppositeCorners(
		r.topLeft.Translate(p),
		r.bottomRight.Translate(p),
	)
}

// Width calculates the width of the rectangle.
//
// Returns:
//   - T: The width of the rectangle, calculated as the absolute difference between the x-coordinates of the top-left and bottom-right corners.
func (r Rectangle[T]) Width() T {
	width := r.bottomRight.x - r.topLeft.x
	if width < 0 {
		return -width // Ensure width is always positive
	}
	return width
}

// NewRectangle creates a new Rectangle from a slice of four points.
// The points can be provided in any order, but they must form an axis-aligned rectangle.
//
// Parameters:
//   - points: []Point[T] - A slice of four points.
//
// Returns:
//   - Rectangle[T]: A new Rectangle initialized with the correct corner points.
//
// Panics:
//   - If the provided points do not form an axis-aligned rectangle, the function panics.
//
// Example Usage:
//
//	rect := NewRectangle([]Point[int]{
//	    NewPoint(0, 10), NewPoint(10, 10), NewPoint(0, 0), NewPoint(10, 0),
//	})
func NewRectangle[T SignedNumber](points []Point[T]) Rectangle[T] {

	if len(points) != 4 {
		panic("NewRectangle requires exactly four points")
	}

	// Find min and max x and y coordinates
	minX, maxX := points[0].x, points[0].x
	minY, maxY := points[0].y, points[0].y

	for _, p := range points[1:] {
		if p.x < minX {
			minX = p.x
		}
		if p.x > maxX {
			maxX = p.x
		}
		if p.y < minY {
			minY = p.y
		}
		if p.y > maxY {
			maxY = p.y
		}
	}

	//fmt.Printf("minX: %v, maxX: %v, minY: %v, maxY: %v\n", minX, maxX, minY, maxY)

	// Validate that the points form an axis-aligned rectangle
	corners := map[Point[T]]bool{
		NewPoint(minX, maxY): false, // top-left
		NewPoint(maxX, maxY): false, // top-right
		NewPoint(minX, minY): false, // bottom-left
		NewPoint(maxX, minY): false, // bottom-right
	}

	for _, p := range points {
		if _, ok := corners[p]; ok {
			corners[p] = true
		} else {
			panic("Points do not form an axis-aligned rectangle")
		}
	}

	for _, found := range corners {
		if !found {
			panic("Points do not form an axis-aligned rectangle")
		}
	}

	// Assign points to the correct fields
	return Rectangle[T]{
		topLeft:     NewPoint(minX, maxY),
		topRight:    NewPoint(maxX, maxY),
		bottomLeft:  NewPoint(minX, minY),
		bottomRight: NewPoint(maxX, minY),
	}
}

// NewRectangleByOppositeCorners creates a rectangle given two opposite corners.
//
// This function determines the top-left and bottom-right corners from the provided points,
// regardless of their order, and ensures a valid axis-aligned rectangle.
//
// Parameters:
//   - corner: Point[T] - One corner of the rectangle.
//   - oppositeCorner: Point[T] - The opposite corner of the rectangle.
//
// Returns:
//   - Rectangle[T]: A new rectangle defined by the determined top-left and bottom-right corners.
func NewRectangleByOppositeCorners[T SignedNumber](corner, oppositeCorner Point[T]) Rectangle[T] {
	return NewRectangle([]Point[T]{
		NewPoint(min(corner.x, oppositeCorner.x), min(corner.y, oppositeCorner.y)),
		NewPoint(min(corner.x, oppositeCorner.x), max(corner.y, oppositeCorner.y)),
		NewPoint(max(corner.x, oppositeCorner.x), min(corner.y, oppositeCorner.y)),
		NewPoint(max(corner.x, oppositeCorner.x), max(corner.y, oppositeCorner.y)),
	})
}

// NewRectangleFromImageRect creates a new Rectangle[T] from an image.Rectangle.
//
// Parameters:
//   - r: image.Rectangle - The image.Rectangle to convert.
//
// Returns:
//   - Rectangle[int]: A new rectangle with integer coordinates matching the given image.Rectangle.
//
// Behavior:
//   - The function maps the minimum point of the `image.Rectangle` to the top-left corner and the
//     maximum point to the bottom-right corner of the resulting rectangle.
//
// Example Usage:
//
//	imgRect := image.Rect(0, 0, 100, 200)
//	rect := NewRectangleFromImageRect(imgRect)
//	// rect represents a rectangle with top-left (0, 0) and bottom-right (100, 200).
func NewRectangleFromImageRect(r image.Rectangle) Rectangle[int] {
	return NewRectangle([]Point[int]{
		NewPoint(r.Min.X, r.Min.Y),
		NewPoint(r.Max.X, r.Max.Y),
		NewPoint(r.Min.X, r.Max.Y),
		NewPoint(r.Max.X, r.Min.Y),
	})
}
