package point

import (
	"fmt"
	"github.com/mikenye/geom2d"
	"math"
)

// todo: doc comments
type OrientationType uint8

// todo: doc comments
const (
	Collinear OrientationType = iota
	Counterclockwise
	Clockwise
)

// todo: doc comments
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

// Orientation checks turn direction using the cross product
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
