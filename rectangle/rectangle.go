package rectangle

import (
	"encoding/json"
	"fmt"
	"github.com/mikenye/geom2d/linesegment"
	"github.com/mikenye/geom2d/point"
	"github.com/mikenye/geom2d/types"
	"image"
)

// Rectangle represents an axis-aligned rectangle defined by its four corners.
type Rectangle struct {
	topLeft     point.Point
	topRight    point.Point
	bottomLeft  point.Point
	bottomRight point.Point
}

// New creates a rectangle given two opposite corners.
//
// This function determines the corners from the provided points,
// regardless of their order, and ensures a valid axis-aligned rectangle.
//
// Parameters:
//   - x1,y1 (float64): One corner of the rectangle.
//   - x2,y2 (float64): The opposite corner of the rectangle.
//
// Returns:
//   - Rectangle: A new rectangle defined by the given opposite corners.
func New(x1, y1, x2, y2 float64) Rectangle {
	return NewFromPoints(
		point.New(min(x1, x2), min(y1, y2)),
		point.New(min(x1, x2), max(y1, y2)),
		point.New(max(x1, x2), min(y1, y2)),
		point.New(max(x1, x2), max(y1, y2)),
	)
}

// NewFromImageRect creates a new Rectangle from an image.Rectangle.
//
// Parameters:
//   - r (image.Rectangle): The image.Rectangle to convert.
//
// Returns:
//   - Rectangle: A new rectangle matching the given [image.Rectangle].
//
// Behavior:
//   - The function maps the minimum point of the [image.Rectangle] to the top-left corner and the
//     maximum point to the bottom-right corner of the resulting rectangle.
func NewFromImageRect(r image.Rectangle) Rectangle {
	return NewFromPoints(
		point.New(float64(r.Min.X), float64(r.Min.Y)),
		point.New(float64(r.Max.X), float64(r.Max.Y)),
		point.New(float64(r.Min.X), float64(r.Max.Y)),
		point.New(float64(r.Max.X), float64(r.Min.Y)),
	)
}

// NewFromPoints creates a new Rectangle from four points.
// The points can be provided in any order, but they must form an axis-aligned rectangle.
//
// Parameters:
//   - pt1,pt2,pt3,pt4 (point.Point): Points forming an axis-aligned rectangle.
//
// Returns:
//   - Rectangle: A new Rectangle initialized with the four given points.
//
// Panics:
//   - If the provided points do not form an axis-aligned rectangle, the function panics.
func NewFromPoints(pt1, pt2, pt3, pt4 point.Point) Rectangle {

	points := []point.Point{pt1, pt2, pt3, pt4}

	// Find min and max x and y coordinates
	minX, maxX := points[0].X(), points[0].X()
	minY, maxY := points[0].Y(), points[0].Y()

	for _, p := range points[1:] {
		minX = min(minX, p.X())
		minY = min(minY, p.Y())
		maxX = max(maxX, p.X())
		maxY = max(maxY, p.Y())
	}

	// Validate that the points form an axis-aligned rectangle
	corners := map[point.Point]bool{
		point.New(minX, maxY): false, // top-left
		point.New(maxX, maxY): false, // top-right
		point.New(minX, minY): false, // bottom-left
		point.New(maxX, minY): false, // bottom-right
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
	return Rectangle{
		topLeft:     point.New(minX, maxY),
		topRight:    point.New(maxX, maxY),
		bottomLeft:  point.New(minX, minY),
		bottomRight: point.New(maxX, minY),
	}
}

// Area calculates the area of the rectangle.
//
// Returns:
//   - float64: The area of the rectangle, calculated as Width * Height.
func (r Rectangle) Area() float64 {
	return r.Width() * r.Height()
}

// ContainsPoint checks if a given point lies within or on the boundary of the Rectangle.
//
// Parameters:
//   - p: The [point.Point] to check.
//
// Returns:
//   - bool: Returns true if the point lies inside or on the boundary of the rectangle, false otherwise.
//
// Behavior:
//   - A point is considered contained if its x-coordinate is between the left and right edges of the [Rectangle],
//     and its y-coordinate is between the top and bottom edges of the rectangle.
//   - The rectangle's boundary is inclusive for both x and y coordinates.
func (r Rectangle) ContainsPoint(p point.Point) bool {
	return p.X() >= r.topLeft.X() &&
		p.X() <= r.bottomRight.X() &&
		p.Y() <= r.topLeft.Y() &&
		p.Y() >= r.bottomRight.Y()
}

// Contour returns the four corner points of the rectangle in the following order:
// top-left, top-right, bottom-right, and bottom-left.
//
// Returns:
//   - bottomLeft, bottomRight, topRight, topLeft (point.Point): The four corner points of the rectangle.
func (r Rectangle) Contour() (bottomLeft, bottomRight, topRight, topLeft point.Point) {
	return r.bottomLeft,
		r.bottomRight,
		r.topRight,
		r.topLeft
}

// Edges returns the edges of the rectangle as line segments.
// Each edge is represented as a line segment connecting two adjacent corners of the rectangle.
//
// Returns:
//   - bottom, right, top, left (linesegment.LineSegment): line segments representing the edges of the rectangle.
func (r Rectangle) Edges() (bottom, right, top, left linesegment.LineSegment) {
	return linesegment.NewFromPoints(r.bottomLeft, r.bottomRight),
		linesegment.NewFromPoints(r.bottomRight, r.topRight),
		linesegment.NewFromPoints(r.topRight, r.topLeft),
		linesegment.NewFromPoints(r.topLeft, r.bottomLeft)
}

// EdgesIter iterates over the edges of the rectangle in counter-clockwise order,
// yielding each edge as a [linesegment.LineSegment].
//
// It is a ["range-over"] function, designed for use with a for-loop.
//
// Example usage:
//
//		rect := NewRectangle(NewPoint(0, 0), NewPoint(4, 3))
//	 fmt.Println("Edges in rect:")
//		for seg := range rect.EdgesIter {
//		    fmt.Printf(" - %s\n", seg)
//		}
//
// The edges are yielded in the following order:
//  1. Bottom edge (left to right)
//  2. Right edge (bottom to top)
//  3. Top edge (right to left)
//  4. Left edge (top to bottom)
//
// If the loop body returns false (due to the for-loop being broken early), iteration stops early.
//
// ["range-over"]: https://go.dev/blog/range-functions
func (r Rectangle) EdgesIter(yield func(segment linesegment.LineSegment) bool) {
	if !yield(linesegment.NewFromPoints(r.bottomLeft, r.bottomRight)) {
		return
	}
	if !yield(linesegment.NewFromPoints(r.bottomRight, r.topRight)) {
		return
	}
	if !yield(linesegment.NewFromPoints(r.topRight, r.topLeft)) {
		return
	}
	if !yield(linesegment.NewFromPoints(r.topLeft, r.bottomLeft)) {
		return
	}
}

// Eq checks if two [Rectangle] instances are equal.
//
// Parameters:
//   - other (Rectangle): The Rectangle to compare against the current Rectangle.
//
// Returns:
//   - bool: Returns true if the two rectangles have identical corner points
//     (bottom-left, bottom-right, top-right, and top-left), false otherwise.
//
// Behavior:
//   - The comparison is based on the exact equality of the corner points.
//   - Both rectangles must have the same coordinates for all four corners to be considered equal.
func (r Rectangle) Eq(other Rectangle) bool {
	return r.bottomLeft.Eq(other.bottomLeft) &&
		r.bottomRight.Eq(other.bottomRight) &&
		r.topRight.Eq(other.topRight) &&
		r.topLeft.Eq(other.topLeft)
}

// Height calculates the height of the rectangle.
//
// Returns:
//   - float64: The height of the rectangle, calculated as the absolute difference between the y-coordinates of the top-left and bottom-right corners.
func (r Rectangle) Height() float64 {
	height := r.bottomRight.Y() - r.topLeft.Y()
	if height < 0 {
		return -height // Ensure height is always positive
	}
	return height
}

// MarshalJSON serializes Rectangle as JSON while preserving its original type.
func (r Rectangle) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		TopLeft     point.Point `json:"top_left"`
		TopRight    point.Point `json:"top_right"`
		BottomLeft  point.Point `json:"bottom_left"`
		BottomRight point.Point `json:"bottom_right"`
	}{
		TopLeft:     r.topLeft,
		TopRight:    r.topRight,
		BottomLeft:  r.bottomLeft,
		BottomRight: r.bottomRight,
	})
}

// Perimeter calculates the perimeter of the rectangle.
//
// Returns:
//   - float64: The perimeter of the rectangle, calculated as 2 * (Width + Height).
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width() + r.Height())
}

// RelationshipToPoint determines the spatial relationship between the current [Rectangle] and a [point.Point].
//
// Relationships:
//   - [types.RelationshipIntersection]: The point lies on one of the rectangle's edges.
//   - [types.RelationshipContainedBy]: The point is inside the rectangle but not on its boundary.
//   - [types.RelationshipDisjoint]: The point lies entirely outside the rectangle.
//
// Parameters:
//   - p (point.Point): The point to analyze the relationship with.
//
// Returns:
//   - [types.Relationship]: The spatial relationship between the point and the rectangle.
//
// Behavior:
//   - The function checks if the point lies on any of the rectangle's edges. If so, it returns [types.RelationshipIntersection].
//   - If the point is not on an edge but is inside the rectangle, it returns [types.RelationshipContainedBy].
//   - If the point is neither on an edge nor inside the rectangle, it returns [types.RelationshipDisjoint].
func (r Rectangle) RelationshipToPoint(p point.Point) types.Relationship {
	for edge := range r.EdgesIter {
		if edge.RelationshipToPoint(p) == types.RelationshipIntersection {
			return types.RelationshipIntersection
		}
	}
	if r.ContainsPoint(p) {
		return types.RelationshipContainedBy
	}
	return types.RelationshipDisjoint
}

// Scale scales the Rectangle relative to a specified reference point by a given scalar factor.
//
// Each corner of the rectangle is scaled relative to the reference point using the provided factor.
// The resulting rectangle maintains its axis-aligned orientation.
//
// Parameters:
//   - ref (point.Point): The reference point relative to which the rectangle is scaled.
//   - k (float64): The scaling factor. A value > 1 enlarges the rectangle; < 1 shrinks it.
//
// Returns:
//   - Rectangle: A new rectangle with corners scaled relative to the reference point.
//
// Notes:
//   - The function delegates the scaling of each corner to the [point.Point.Scale] method.
//   - The rectangle remains axis-aligned after scaling.
//   - If the scaling factor k is 1, the rectangle remains unchanged.
func (r Rectangle) Scale(ref point.Point, k float64) Rectangle {
	return NewFromPoints(
		r.topLeft.Scale(ref, k),
		r.topRight.Scale(ref, k),
		r.bottomLeft.Scale(ref, k),
		r.bottomRight.Scale(ref, k),
	)
}

// ScaleHeight scales the height of the rectangle by the given factor,
// keeping the bottom edge fixed and adjusting the top edge proportionally.
//
// Parameters:
//   - factor (float64): The scaling factor to apply to the height. A value of 1
//     keeps the height unchanged, a value greater than 1 increases the height,
//     and a value between 0 and 1 decreases the height. Negative values may
//     result in unexpected behavior.
//
// Returns:
//   - Rectangle: A new rectangle with the scaled height, maintaining the
//     same width and bottom edge.
//
// Behavior:
//   - The bottom edge of the rectangle remains fixed in place.
//   - The top edge is adjusted vertically by scaling the height of the
//     rectangle by the specified factor.
//
// Notes:
//   - If the scaling factor is negative, the top edge will be positioned below
//     the bottom edge, potentially creating an inverted rectangle.
func (r Rectangle) ScaleHeight(factor float64) Rectangle {
	newTopY := r.bottomLeft.Y() + (r.Height() * factor)
	return NewFromPoints(
		point.New(r.topLeft.X(), newTopY),
		point.New(r.topRight.X(), newTopY),
		r.bottomLeft,
		r.bottomRight,
	)
}

// ScaleWidth scales the width of the rectangle by the given factor,
// keeping the left edge fixed and adjusting the right edge proportionally.
//
// Parameters:
//   - factor (float64): The scaling factor to apply to the width. A value of 1
//     keeps the width unchanged, a value greater than 1 increases the width,
//     and a value between 0 and 1 decreases the width. Negative values may
//     result in unexpected behavior.
//
// Returns:
//   - Rectangle: A new rectangle with the scaled width, maintaining the
//     same height and left edge.
//
// Behavior:
//   - The left edge of the rectangle remains fixed in place.
//   - The right edge is adjusted horizontally by scaling the width of the
//     rectangle by the specified factor.
//
// Notes:
//   - If the scaling factor is negative, the right edge will be positioned to
//     the left of the left edge, potentially creating an inverted rectangle.
func (r Rectangle) ScaleWidth(factor float64) Rectangle {
	newRightX := r.bottomLeft.X() + (r.Width() * factor)
	return NewFromPoints(
		point.New(newRightX, r.bottomRight.Y()),
		point.New(newRightX, r.topRight.Y()),
		r.bottomLeft,
		r.topLeft,
	)
}

// String returns a string representation of the rectangle.
// The representation includes the coordinates of the rectangle's corners in counter-clockwise order,
// in the format: "[(bottomLeft),(topRight)]".
//
// This is primarily useful for debugging and logging.
//
// Returns:
//   - string: A formatted string showing the coordinates of the rectangle's corners.
func (r Rectangle) String() string {
	return fmt.Sprintf("[(%v,%v),(%v,%v)]", r.bottomLeft.X(), r.bottomLeft.Y(), r.topRight.X(), r.topRight.Y())
}

// ToImageRect converts the Rectangle to an image.Rectangle.
// This method converts the rectangle coordinates to integers by truncating the decimal portion
// (via int type casting) as image.Rectangle requires integer coordinates.
//
// Returns:
//   - image.Rectangle: A new image.Rectangle with coordinates matching the Rectangle.
func (r Rectangle) ToImageRect() image.Rectangle {
	return image.Rect(
		int(r.topLeft.X()),
		int(r.topLeft.Y()),
		int(r.bottomRight.X()),
		int(r.bottomRight.Y()),
	)
}

// Translate moves the rectangle by a specified vector.
//
// This method shifts the rectangle's position in the 2D plane by translating
// both its corners by the given vector p. The
// dimensions of the rectangle remain unchanged.
//
// Parameters:
//   - p (point.Point): The vector by which to translate the rectangle.
//
// Returns:
//   - Rectangle: A new Rectangle translated by the specified vector.
func (r Rectangle) Translate(p point.Point) Rectangle {
	return NewFromPoints(
		r.topLeft.Translate(p),
		r.topRight.Translate(p),
		r.bottomLeft.Translate(p),
		r.bottomRight.Translate(p),
	)
}

// UnmarshalJSON deserializes JSON into a Rectangle while keeping the exact original type.
func (r *Rectangle) UnmarshalJSON(data []byte) error {
	var temp struct {
		TopLeft     point.Point `json:"top_left"`
		TopRight    point.Point `json:"top_right"`
		BottomLeft  point.Point `json:"bottom_left"`
		BottomRight point.Point `json:"bottom_right"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	r.topLeft = temp.TopLeft
	r.topRight = temp.TopRight
	r.bottomLeft = temp.BottomLeft
	r.bottomRight = temp.BottomRight
	return r.validate()
}

// validate checks if the rectangle is axis-aligned and correctly ordered.
func (r Rectangle) validate() error {
	if r.topLeft.Y() != r.topRight.Y() {
		return fmt.Errorf("topLeft (%v) and topRight (%v) must have the same y-coordinate", r.topLeft, r.topRight)
	}
	if r.bottomLeft.Y() != r.bottomRight.Y() {
		return fmt.Errorf("bottomLeft (%v) and bottomRight (%v) must have the same y-coordinate", r.bottomLeft, r.bottomRight)
	}
	if r.topLeft.X() != r.bottomLeft.X() {
		return fmt.Errorf("topLeft (%v) and bottomLeft (%v) must have the same x-coordinate", r.topLeft, r.bottomLeft)
	}
	if r.topRight.X() != r.bottomRight.X() {
		return fmt.Errorf("topRight (%v) and bottomRight (%v) must have the same x-coordinate", r.topRight, r.bottomRight)
	}
	if r.topLeft.Y() <= r.bottomLeft.Y() {
		return fmt.Errorf("topLeft (%v) must be above bottomLeft (%v)", r.topLeft, r.bottomLeft)
	}
	if r.topRight.Y() <= r.bottomRight.Y() {
		return fmt.Errorf("topRight (%v) must be above bottomRight (%v)", r.topRight, r.bottomRight)
	}
	if r.topLeft.X() >= r.topRight.X() {
		return fmt.Errorf("topLeft (%v) must be to the left of topRight (%v)", r.topLeft, r.topRight)
	}
	if r.bottomLeft.X() >= r.bottomRight.X() {
		return fmt.Errorf("bottomLeft (%v) must be to the left of bottomRight (%v)", r.bottomLeft, r.bottomRight)
	}
	return nil // Rectangle is valid
}

// Width calculates the width of the rectangle.
//
// Returns:
//   - float64: The width of the rectangle, calculated as the absolute difference between the x-coordinates of the top-left and bottom-right corners.
func (r Rectangle) Width() float64 {
	width := r.bottomRight.X() - r.topLeft.X()
	if width < 0 {
		return -width // Ensure width is always positive
	}
	return width
}
