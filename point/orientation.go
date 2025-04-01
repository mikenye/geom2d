package point

import (
	"fmt"
	"github.com/mikenye/geom2d"
	"math"
)

// OrientationType represents the orientation relationship between three points in a 2D plane.
//
// This type categorizes how three points are oriented relative to each other:
// - Collinear: The points lie on a straight line
// - Clockwise: The points form a clockwise turn
// - Counterclockwise: The points form a counterclockwise turn
//
// The orientation is determined by evaluating the cross product of vectors formed by these points,
// and is a fundamental concept in computational geometry algorithms such as convex hull construction,
// line segment intersection detection, and polygon operations.
type OrientationType uint8

// Orientation constants define the possible orientation relationships between three points.
const (
	// Collinear indicates that three points lie on a straight line.
	Collinear OrientationType = iota

	// Counterclockwise indicates that three points form a counterclockwise turn.
	Counterclockwise

	// Clockwise indicates that three points form a clockwise turn.
	Clockwise
)

// String returns a human-readable string representation of the orientation type.
//
// This method allows OrientationType values to be easily displayed and logged.
//
// Returns:
//   - string: The name of the orientation: "Collinear", "Counterclockwise", or "Clockwise".
//
// Panics:
//   - If the OrientationType value is not one of the defined constants.
func (o OrientationType) String() string {
	switch o {
	case Collinear:
		return "Collinear"
	case Counterclockwise:
		return "Counterclockwise"
	case Clockwise:
		return "Clockwise"
	default:
		panic(fmt.Errorf("unsupported point orientation: %d", o))
	}
}

// Orientation determines the relative orientation of three points in a 2D plane.
//
// This function calculates whether three points p, q, and r make a clockwise turn,
// a counterclockwise turn, or are collinear, using the cross product of the vectors (q-p) and (r-p).
//
// Parameters:
//   - p, q, r (Point): The three points to determine orientation for
//
// Returns:
//   - OrientationType: The orientation relationship:
//   - Collinear: The points lie on a straight line
//   - Clockwise: The points make a clockwise turn
//   - Counterclockwise: The points make a counterclockwise turn
//
// Behavior:
//   - Uses an adaptive epsilon based on the distance between points to handle floating-point precision
//   - Relies on the sign of the cross product:
//   - Positive → Counterclockwise
//   - Negative → Clockwise
//   - Near zero (within epsilon) → Collinear
//
// Note:
//   - This is a fundamental operation for many computational geometry algorithms including
//     convex hull construction, polygon triangulation, and line segment intersection detection.
func Orientation(p, q, r Point) OrientationType {
	val := (q.Sub(p)).CrossProduct(r.Sub(p))

	// Compute adaptive epsilon based on segment lengths
	epsilon := geom2d.GetEpsilon() * (p.DistanceToPoint(q) + p.DistanceToPoint(r))

	if math.Abs(val) < epsilon {
		return Collinear // Collinear
	}
	if val > 0 {
		return Counterclockwise // Counterclockwise
	}
	return Clockwise // Clockwise
}
