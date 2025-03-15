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
func Orientation(p, q, r Point) int {
	val := (q.Sub(p)).CrossProduct(r.Sub(p))
	if math.Abs(val) < geom2d.GetEpsilon() {
		return 0 // Collinear
	}
	if val > 0 {
		return 1 // Counterclockwise
	}
	return -1 // Clockwise
}
