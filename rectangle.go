package geom2d

import (
	"image"
	"slices"
)

// Rectangle represents a 2D rectangle defined by its top-left and bottom-right corners.
type Rectangle[T SignedNumber] struct {
	topLeft     Point[T]
	bottomRight Point[T]
}

// Add adds a vector (Point[T]) or another Rectangle[T] to the current rectangle.
// If a Point is added, it shifts the rectangle. If a Rectangle is added, it adjusts both corners.
//
// Parameters:
//   - p: The Point[T] to add to both corners of the rectangle.
//
// Returns:
//   - Rectangle[T]: A new Rectangle with adjusted corners.
func (r Rectangle[T]) Add(p Point[T]) Rectangle[T] {
	return NewRectangleByPoints(
		r.topLeft.Add(p),
		r.bottomRight.Add(p),
	)
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
	return NewRectangleByPoints(
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
	return NewRectangleByPoints(
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
	return NewRectangleByPoints(
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
		p.y >= r.topLeft.y &&
		p.y <= r.bottomRight.y
}

// Div divides the rectangle’s dimensions by a scalar value.
//
// Parameters:
//   - divisor: The divisor by which to scale the rectangle.
//
// Returns:
//   - Rectangle[float64]: A new Rectangle with dimensions divided by the divisor.
func (r Rectangle[T]) Div(divisor float64) Rectangle[float64] {
	return r.Scale(1 / divisor)
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

// IsLineSegmentOnEdge checks if the given line segment lies entirely on one of the rectangle's edges.
//
// Parameters:
//   - segment: The line segment to check.
//
// Returns:
//   - bool: True if the segment lies on one of the rectangle's edges; otherwise, false.
func (r Rectangle[T]) IsLineSegmentOnEdge(segment LineSegment[T]) bool {
	return (segment.start.y == r.topLeft.y && segment.end.y == r.topLeft.y && segment.start.x >= r.topLeft.x && segment.end.x <= r.bottomRight.x) || // Top edge
		(segment.start.y == r.bottomRight.y && segment.end.y == r.bottomRight.y && segment.start.x >= r.topLeft.x && segment.end.x <= r.bottomRight.x) || // Bottom edge
		(segment.start.x == r.topLeft.x && segment.end.x == r.topLeft.x && segment.start.y >= r.topLeft.y && segment.end.y <= r.bottomRight.y) || // Left edge
		(segment.start.x == r.bottomRight.x && segment.end.x == r.bottomRight.x && segment.start.y >= r.topLeft.y && segment.end.y <= r.bottomRight.y) // Right edge
}

// IsLineSegmentOnEdgeWithEndTouchingVertex checks if the given line segment lies on an edge of the rectangle
// and one or both of its endpoints touch a vertex.
//
// Parameters:
//   - segment: The line segment to check.
//
// Returns:
//   - bool: True if the segment lies on an edge and touches a vertex; otherwise, false.
func (r Rectangle[T]) IsLineSegmentOnEdgeWithEndTouchingVertex(segment LineSegment[T]) bool {
	vertices := []Point[T]{
		r.topLeft,
		NewPoint(r.topLeft.x, r.bottomRight.y), // Bottom-left
		NewPoint(r.bottomRight.x, r.topLeft.y), // Top-right
		r.bottomRight,
	}
	return r.IsLineSegmentOnEdge(segment) &&
		(slices.Contains(vertices, segment.start) || slices.Contains(vertices, segment.end))
}

// LineSegmentEntersAndExits checks if a LineSegment enters the rectangle through one edge
// and exits through another, indicating that it crosses through the rectangle.
//
// Parameters:
//   - segment: The line segment to check.
//
// Returns:
//   - bool: True if the segment enters through one edge and exits through another; otherwise, false.
//
// Explanation of Logic:
//   - The function iterates over each edge of the rectangle and checks the relationship
//     between the segment and each edge using the RelationshipToLineSegment function.
//   - If the segment strictly intersects an edge (LSRIntersects) or one end of the segment
//     lies on an edge without the entire segment being on that edge (LSRConAB), it is counted
//     as an entry or exit point. We only test LSRConAB (and not LSRDonAB) to avoid double counting.
//   - If there is more than one intersection or endpoint contact with the rectangle's edges,
//     the segment is considered to "enter and exit," returning true.
//   - This approach prevents double-counting cases where the segment might lie on or touch an
//     edge without fully crossing into the rectangle.
func (r Rectangle[T]) LineSegmentEntersAndExits(segment LineSegment[T]) bool {
	entryCount := 0

	edges := []LineSegment[T]{
		NewLineSegment(r.topLeft, NewPoint(r.bottomRight.x, r.topLeft.y)),     // Top edge
		NewLineSegment(NewPoint(r.bottomRight.x, r.topLeft.y), r.bottomRight), // Right edge
		NewLineSegment(r.bottomRight, NewPoint(r.topLeft.x, r.bottomRight.y)), // Bottom edge
		NewLineSegment(NewPoint(r.topLeft.x, r.bottomRight.y), r.topLeft),     // Left edge
	}

	var rel TwoLinesRelationship
	for _, edge := range edges {

		// Check for intersections or an endpoint lying on an edge without full overlap.
		// We only test LSRConAB (and not LSRDonAB) to avoid double counting.
		rel = segment.RelationshipToLineSegment(edge)
		if rel == LSRIntersects || rel == LSRConAB {
			entryCount++
		}
		if entryCount > 1 {
			return true // Enters and exits
		}
	}
	return false
}

// LineSegmentIntersectsEdges checks if the given line segment intersects one or more edges of the rectangle.
//
// Parameters:
//   - segment: The line segment to check.
//
// Returns:
//   - bool: True if the segment intersects any edge; otherwise, false.
func (r Rectangle[T]) LineSegmentIntersectsEdges(segment LineSegment[T]) bool {
	edges := []LineSegment[T]{
		NewLineSegment(r.topLeft, NewPoint(r.bottomRight.x, r.topLeft.y)),     // Top edge
		NewLineSegment(NewPoint(r.bottomRight.x, r.topLeft.y), r.bottomRight), // Right edge
		NewLineSegment(r.bottomRight, NewPoint(r.topLeft.x, r.bottomRight.y)), // Bottom edge
		NewLineSegment(NewPoint(r.topLeft.x, r.bottomRight.y), r.topLeft),     // Left edge
	}

	var rel TwoLinesRelationship
	for _, edge := range edges {
		rel = segment.RelationshipToLineSegment(edge)
		switch rel {
		case LSRIntersects, LSRConAB, LSRDonAB:
			return true
		default:
		}
	}
	return false
}

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
		NewPoint(r.bottomRight.x, r.topLeft.y), // Top-right
		r.bottomRight,
		NewPoint(r.topLeft.x, r.bottomRight.y), // Bottom-left
	}
}

// RelationshipToLineSegment determines the relationship between a line segment and the rectangle.
// It returns one of several values indicating whether the segment is inside, outside, on an edge,
// touches vertices or edges, or intersects the rectangle in different ways.
//
// Returns:
//   - RectangleSegmentRelationship: Enum value describing the relationship of the segment to the rectangle.
func (r Rectangle[T]) RelationshipToLineSegment(segment LineSegment[T]) RectangleSegmentRelationship {
	// Determine relationships of each endpoint of the segment to the rectangle
	startRelationship := r.RelationshipToPoint(segment.start)
	endRelationship := r.RelationshipToPoint(segment.end)

	// Handle degenerate segment (start and end points are the same)
	if segment.start == segment.end {
		switch startRelationship {
		case PRRInside:
			return RSRInside
		case PRROnVertex:
			return RSROutsideEndTouchesVertex
		case PRROnEdge:
			return RSROutsideEndTouchesEdge
		default:
			return RSROutside
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
	edgeRelationships := make([]TwoLinesRelationship, len(edges))
	for i, edge := range edges {
		edgeRelationships[i] = edge.RelationshipToLineSegment(segment)
	}

	// Check if segment enters and exists
	if countOccurrences(edgeRelationships, LSRAonCD) >= 1 &&
		countOccurrences(edgeRelationships, LSRBonCD) >= 1 &&
		startRelationship == PRROutside && endRelationship == PRROutside {
		return RSREntersAndExits
	}

	// Check if segment fully inside
	if adjoiningEdges(edgeRelationships, LSRAeqC, LSRBeqD) &&
		adjoiningEdges(edgeRelationships, LSRAeqD, LSRBeqC) &&
		startRelationship == PRROnVertex && endRelationship == PRROnVertex {
		return RSRInsideEndTouchesVertex
	}

	// Check if segment is inside, with one end on an edge
	if countOccurrences(edgeRelationships, LSRMiss) == 3 &&
		(countOccurrences(edgeRelationships, LSRDonAB) == 1 ||
			countOccurrences(edgeRelationships, LSRConAB) == 1) &&
		((startRelationship == PRROnEdge && endRelationship == PRRInside) ||
			(startRelationship == PRRInside && endRelationship == PRROnEdge)) {
		return RSRInsideEndTouchesEdge
	}

	// Check if segment is fully outside
	if countOccurrences(edgeRelationships, LSRMiss) == len(edges) &&
		startRelationship == PRROutside && endRelationship == PRROutside {
		return RSROutside
	}

	// Check if the segment lies entirely on an edge
	if countOccurrences(edgeRelationships, LSRCollinearEqual) == 1 {
		if r.IsLineSegmentOnEdgeWithEndTouchingVertex(segment) {
			return RSROnEdgeEndTouchesVertex
		}
		return RSROnEdge
	}
	if countOccurrences(edgeRelationships, LSRCollinearCDinAB) == 1 &&
		countOccurrences(edgeRelationships, LSRMiss) == 3 &&
		startRelationship == PRROnEdge && endRelationship == PRROnEdge {
		return RSROnEdge
	}

	// Check if the segment intersects the rectangle through one or more edges
	intersectionCount := countOccurrences(edgeRelationships, LSRIntersects)
	if intersectionCount == 1 {
		return RSRIntersects
	} else if intersectionCount > 1 {
		return RSREntersAndExits
	}

	// Check if one endpoint is on an edge and the other is outside
	if countOccurrences(edgeRelationships, LSRConAB) == 1 && endRelationship == PRROutside {
		return RSROutsideEndTouchesEdge
	} else if countOccurrences(edgeRelationships, LSRDonAB) == 1 && startRelationship == PRROutside {
		return RSROutsideEndTouchesEdge
	}

	// Check if one endpoint is on a vertex and the other is inside or outside
	if countOccurrences(edgeRelationships, LSRAeqC) == 1 && adjoiningEdges(edgeRelationships, LSRAeqC, LSRBeqC) {
		if endRelationship == PRRInside {
			return RSRInsideEndTouchesVertex
		}
		return RSROutsideEndTouchesVertex
	} else if countOccurrences(edgeRelationships, LSRAeqD) == 1 && adjoiningEdges(edgeRelationships, LSRAeqD, LSRBeqD) {
		if startRelationship == PRRInside {
			return RSRInsideEndTouchesVertex
		}
		return RSROutsideEndTouchesVertex
	}
	if countOccurrences(edgeRelationships, LSRMiss) == 2 &&
		adjoiningEdges(edgeRelationships, LSRAonCD, LSRBonCD) &&
		((startRelationship == PRROutside && endRelationship == PRRInside) ||
			(startRelationship == PRRInside && endRelationship == PRROutside)) {
		return RSRIntersects
	}

	if countOccurrences(edgeRelationships, LSRMiss) == 2 &&
		(adjoiningEdges(edgeRelationships, LSRAonCD, LSRCollinearBonCD) ||
			adjoiningEdges(edgeRelationships, LSRBonCD, LSRCollinearAonCD)) &&
		((startRelationship == PRROnEdge && endRelationship == PRROutside) ||
			(startRelationship == PRROutside && endRelationship == PRROnEdge)) {
		return RSROutsideEndTouchesEdge
	}

	// If both endpoints are inside
	if startRelationship == PRRInside && endRelationship == PRRInside &&
		countOccurrences(edgeRelationships, LSRMiss) == len(edges) {
		return RSRInside
	}

	return RSROutside
}

// RelationshipToPoint determines the relationship of a point to the rectangle.
// The relationship can be Inside, Outside, On a Vertex, or On an Edge.
//
// Parameters:
//   - p: The point to check.
//
// Returns:
//   - PointRectangleRelationship: The relationship of the point to the rectangle.
func (r Rectangle[T]) RelationshipToPoint(p Point[T]) PointRectangleRelationship {
	// Check if the point is strictly inside
	if p.x > r.topLeft.x && p.x < r.bottomRight.x &&
		p.y > r.topLeft.y && p.y < r.bottomRight.y {
		return PRRInside
	}

	// Check if the point is on a vertex
	if (p == r.topLeft) ||
		(p == r.bottomRight) ||
		(p == NewPoint(r.topLeft.x, r.bottomRight.y)) || // Bottom-left vertex
		(p == NewPoint(r.bottomRight.x, r.topLeft.y)) { // Top-right vertex
		return PRROnVertex
	}

	// Check if the point is on an edge
	if (p.x == r.topLeft.x || p.x == r.bottomRight.x) && (p.y >= r.topLeft.y && p.y <= r.bottomRight.y) ||
		(p.y == r.topLeft.y || p.y == r.bottomRight.y) && (p.x >= r.topLeft.x && p.x <= r.bottomRight.x) {
		return PRROnEdge
	}

	// Otherwise, the point is outside
	return PRROutside
}

// Scale scales the rectangle by a scalar value from the top-left corner.
//
// Parameters:
//   - factor: The scaling factor. A value > 1 enlarges the rectangle; < 1 shrinks it.
//
// Returns:
//   - Rectangle[float64]: A new Rectangle scaled by the given factor.
func (r Rectangle[T]) Scale(factor float64) Rectangle[float64] {
	topLeft := r.topLeft.AsFloat()
	bottomRight := r.bottomRight.AsFloat()
	return NewRectangleByPoints(
		topLeft,
		NewPoint(
			topLeft.x+(bottomRight.x-topLeft.x)*factor,
			topLeft.y+(bottomRight.y-topLeft.y)*factor,
		),
	)
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
	return NewRectangleByPoints(
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
	return NewRectangleByPoints(
		topLeft,
		NewPoint(topLeft.x+width, r.bottomRight.AsFloat().y),
	)
}

// Sub subtracts a vector (Point[T]) or another Rectangle[T] from the current rectangle.
//
// Parameters:
//   - p: The Point[T] to subtract from both corners of the rectangle.
//
// Returns:
//   - Rectangle[T]: A new Rectangle with adjusted corners.
func (r Rectangle[T]) Sub(p Point[T]) Rectangle[T] {
	return NewRectangleByPoints(
		r.topLeft.Sub(p),
		r.bottomRight.Sub(p),
	)
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

// NewRectangleByDimensions creates a rectangle given the top-left point, width, and height.
//
// Parameters:
//   - topLeft: The top-left corner of the rectangle.
//   - width: The width of the rectangle.
//   - height: The height of the rectangle.
//
// Returns:
//   - Rectangle[T]: A new rectangle defined by the given dimensions.
func NewRectangleByDimensions[T SignedNumber](topLeft Point[T], width, height T) Rectangle[T] {
	return Rectangle[T]{
		topLeft:     topLeft,
		bottomRight: NewPoint(topLeft.x+width, topLeft.y+height),
	}
}

// NewRectangleByPoints creates a rectangle given the top-left and bottom-right points.
//
// Parameters:
//   - topLeft: The top-left corner of the rectangle.
//   - bottomRight: The bottom-right corner of the rectangle.
//
// Returns:
//   - Rectangle[T]: A new rectangle defined by the given points.
func NewRectangleByPoints[T SignedNumber](topLeft, bottomRight Point[T]) Rectangle[T] {
	return Rectangle[T]{topLeft: topLeft, bottomRight: bottomRight}
}

// NewRectangleFromImageRect creates a new Rectangle[T] from an image.Rectangle.
//
// Parameters:
//   - r: The image.Rectangle to convert.
//
// Returns:
//   - Rectangle[int]: A new Rectangle with integer coordinates matching the given image.Rectangle.
//
// Example Usage:
//
//	imgRect := image.Rect(0, 0, 100, 200)
//	rect := NewRectangleFromImageRect(imgRect) // Rectangle with same coordinates as imgRect.
func NewRectangleFromImageRect(r image.Rectangle) Rectangle[int] {
	return NewRectangleByPoints(
		NewPoint(r.Min.X, r.Min.Y),
		NewPoint(r.Max.X, r.Max.Y),
	)
}
