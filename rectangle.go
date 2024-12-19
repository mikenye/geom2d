package geom2d

import (
	"fmt"
	"image"
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

// AsFloat32 converts the Rectangle's corner points to the float32 type, useful for higher-precision operations.
//
// Returns:
//   - Rectangle[float32]: A new Rectangle with float32 coordinates.
func (r Rectangle[T]) AsFloat32() Rectangle[float32] {
	return newRectangleByOppositeCorners(
		r.topLeft.AsFloat32(),
		r.bottomRight.AsFloat32(),
	)
}

// AsFloat64 converts the Rectangle's corner points to the float64 type, useful for higher-precision operations.
//
// Returns:
//   - Rectangle[float64]: A new Rectangle with float64 coordinates.
func (r Rectangle[T]) AsFloat64() Rectangle[float64] {
	return newRectangleByOppositeCorners(
		r.topLeft.AsFloat64(),
		r.bottomRight.AsFloat64(),
	)
}

// AsInt converts the Rectangle's corner points to the int type by truncating any decimal values.
// This method is useful for operations requiring integer coordinates.
//
// Returns:
//   - Rectangle[int]: A new Rectangle with integer coordinates, truncated from the original values.
func (r Rectangle[T]) AsInt() Rectangle[int] {
	return newRectangleByOppositeCorners(
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
	return newRectangleByOppositeCorners(
		r.topLeft.AsIntRounded(),
		r.bottomRight.AsIntRounded(),
	)
}

// ContainsPoint checks if a given [Point] lies within or on the boundary of the [Rectangle].
//
// Parameters:
//   - p: The [Point] to check.
//
// Returns:
//   - bool: Returns true if the point lies inside or on the boundary of the rectangle, false otherwise.
//
// Behavior:
//   - A point is considered contained if its x-coordinate is between the left and right edges of the [Rectangle],
//     and its y-coordinate is between the top and bottom edges of the rectangle.
//   - The rectangle's boundary is inclusive for both x and y coordinates.
func (r Rectangle[T]) ContainsPoint(p Point[T]) bool {
	return p.x >= r.topLeft.x &&
		p.x <= r.bottomRight.x &&
		p.y <= r.topLeft.y &&
		p.y >= r.bottomRight.y
}

// Edges returns the edges of the rectangle as a slice of [LineSegment][T].
// Each edge is represented as a line segment connecting two adjacent corners of the rectangle.
//
// Returns:
//   - [][LineSegment][T]: A slice of line segments representing the edges of the rectangle.
func (r Rectangle[T]) Edges() []LineSegment[T] {
	return []LineSegment[T]{
		NewLineSegment(r.bottomLeft, r.bottomRight),
		NewLineSegment(r.bottomRight, r.topRight),
		NewLineSegment(r.topRight, r.topLeft),
		NewLineSegment(r.topLeft, r.bottomLeft),
	}
}

// Eq checks if two [Rectangle] instances are equal.
//
// Parameters:
//   - other (Rectangle[T]): The [Rectangle] to compare against the current [Rectangle].
//
// Returns:
//   - bool: Returns true if the two rectangles have identical corner points
//     (bottom-left, bottom-right, top-right, and top-left), false otherwise.
//
// Behavior:
//   - The comparison is based on the exact equality of the corner points.
//   - Both rectangles must have the same coordinates for all four corners to be considered equal.
func (r Rectangle[T]) Eq(other Rectangle[T]) bool {
	if r.bottomLeft == other.bottomLeft &&
		r.bottomRight == other.bottomRight &&
		r.topRight == other.topRight &&
		r.topLeft == other.topLeft {
		return true
	}
	return false
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

// Perimeter calculates the perimeter of the rectangle.
//
// Returns:
//   - T: The perimeter of the rectangle, calculated as 2 * (Width + Height).
func (r Rectangle[T]) Perimeter() T {
	return 2 * (r.Width() + r.Height())
}

// Contour returns the four corner points of the rectangle in the following order:
// top-left, top-right, bottom-right, and bottom-left.
//
// Returns:
//   - [][Point][T]: A slice containing the four corner points of the rectangle.
func (r Rectangle[T]) Contour() []Point[T] {
	return []Point[T]{
		r.topLeft,
		r.topRight,
		r.bottomRight,
		r.bottomLeft,
	}
}

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

// Scale scales the [Rectangle] relative to a specified reference [Point] by a given scalar factor.
//
// Each corner of the rectangle is scaled relative to the reference point using the provided factor.
// The resulting rectangle maintains its axis-aligned orientation.
//
// Parameters:
//   - ref ([Point][T]): The reference point relative to which the rectangle is scaled.
//   - k (T): The scaling factor. A value > 1 enlarges the rectangle; < 1 shrinks it.
//
// Returns:
//   - [Rectangle][T]: A new rectangle with corners scaled relative to the reference point.
//
// Notes:
//   - The function delegates the scaling of each corner to the [Point.Scale] method.
//   - The rectangle remains axis-aligned after scaling.
//   - If the scaling factor k is 1, the rectangle remains unchanged.
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
//   - factor (float64): The scaling factor for the height. A value > 1 enlarges the height; < 1 shrinks it.
//
// Returns:
//   - Rectangle[float64]: A new Rectangle with the height scaled by the given factor.
func (r Rectangle[T]) ScaleHeight(factor float64) Rectangle[float64] {
	topLeft := r.topLeft.AsFloat64()
	height := float64(r.Height()) * factor
	return newRectangleByOppositeCorners(
		topLeft,
		NewPoint(r.bottomRight.AsFloat64().x, topLeft.y+height),
	)
}

// ScaleWidth scales the width of the rectangle by a scalar value.
//
// Parameters:
//   - factor (float64): The scaling factor for the width. A value > 1 enlarges the width; < 1 shrinks it.
//
// Returns:
//   - Rectangle[float64]: A new Rectangle with the width scaled by the given factor.
func (r Rectangle[T]) ScaleWidth(factor float64) Rectangle[float64] {
	topLeft := r.topLeft.AsFloat64()
	width := float64(r.Width()) * factor
	return newRectangleByOppositeCorners(
		topLeft,
		NewPoint(topLeft.x+width, r.bottomRight.AsFloat64().y),
	)
}

// String returns a string representation of the rectangle.
// The representation includes the coordinates of the rectangle's corners in counter-clockwise order,
// in the format:
// "Rectangle[(bottomLeft), (bottomRight), (topRight), (topLeft)]".
//
// This is primarily useful for debugging and logging.
//
// Returns:
//   - string: A formatted string showing the coordinates of the rectangle's corners.
func (r Rectangle[T]) String() string {
	return fmt.Sprintf("Rectangle[(%v, %v), (%v, %v), (%v, %v), (%v, %v)]", r.bottomLeft.x, r.bottomLeft.y, r.bottomRight.x, r.bottomRight.y, r.topRight.x, r.topRight.y, r.topLeft.x, r.topLeft.y)
}

// ToImageRect converts the [Rectangle][int] to an [image.Rectangle].
// This method is only available for [Rectangle][int] as [image.Rectangle] requires integer coordinates.
//
// Returns:
//   - [image.Rectangle]: A new [image.Rectangle] with coordinates matching the [Rectangle].
func (r Rectangle[int]) ToImageRect() image.Rectangle {
	topLeft := r.topLeft.AsInt()
	bottomRight := r.bottomRight.AsInt()
	return image.Rect(topLeft.x, topLeft.y, bottomRight.x, bottomRight.y)
}

// Translate moves the rectangle by a specified vector.
//
// This method shifts the rectangle's position in the 2D plane by translating
// both its top-left and bottom-right corners by the given vector p. The
// dimensions of the rectangle remain unchanged.
//
// Parameters:
//   - p ([Point][T]): The vector by which to translate the rectangle.
//
// Returns:
//   - [Rectangle][T]: A new [Rectangle] translated by the specified vector.
func (r Rectangle[T]) Translate(p Point[T]) Rectangle[T] {
	return newRectangleByOppositeCorners(
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
//   - points ([][Point][T]): A slice of four points.
//
// Returns:
//   - [Rectangle][T]: A new [Rectangle] initialized with the correct corner points.
//
// Panics:
//   - If the provided points do not form an axis-aligned rectangle, the function panics.
func NewRectangle[T SignedNumber](points []Point[T]) Rectangle[T] {

	if len(points) != 4 {
		panic("NewRectangle requires exactly four points")
	}

	// Find min and max x and y coordinates
	minX, maxX := points[0].x, points[0].x
	minY, maxY := points[0].y, points[0].y

	for _, p := range points[1:] {
		minX = min(minX, p.x)
		minY = min(minY, p.y)
		maxX = max(maxX, p.x)
		maxY = max(maxY, p.y)
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

// newRectangleByOppositeCorners creates a rectangle given two opposite corners.
//
// This function determines the top-left and bottom-right corners from the provided points,
// regardless of their order, and ensures a valid axis-aligned rectangle.
//
// Parameters:
//   - corner ([Point][T]): One corner of the rectangle.
//   - oppositeCorner ([Point][T]): The opposite corner of the rectangle.
//
// Returns:
//   - [Rectangle][T]: A new rectangle defined by the determined top-left and bottom-right corners.
func newRectangleByOppositeCorners[T SignedNumber](corner, oppositeCorner Point[T]) Rectangle[T] {
	return NewRectangle([]Point[T]{
		NewPoint(min(corner.x, oppositeCorner.x), min(corner.y, oppositeCorner.y)),
		NewPoint(min(corner.x, oppositeCorner.x), max(corner.y, oppositeCorner.y)),
		NewPoint(max(corner.x, oppositeCorner.x), min(corner.y, oppositeCorner.y)),
		NewPoint(max(corner.x, oppositeCorner.x), max(corner.y, oppositeCorner.y)),
	})
}

// NewRectangleFromImageRect creates a new [Rectangle][T] from an [image.Rectangle].
//
// Parameters:
//   - r [image.Rectangle]: The [image.Rectangle] to convert.
//
// Returns:
//   - [Rectangle][int]: A new rectangle with integer coordinates matching the given [image.Rectangle].
//
// Behavior:
//   - The function maps the minimum point of the [image.Rectangle] to the top-left corner and the
//     maximum point to the bottom-right corner of the resulting rectangle.
func NewRectangleFromImageRect(r image.Rectangle) Rectangle[int] {
	return NewRectangle([]Point[int]{
		NewPoint(r.Min.X, r.Min.Y),
		NewPoint(r.Max.X, r.Max.Y),
		NewPoint(r.Min.X, r.Max.Y),
		NewPoint(r.Max.X, r.Min.Y),
	})
}
