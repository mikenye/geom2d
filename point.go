package geom2d

import "errors"

var (
	errPointOutOfBounds = errors.New("point is out of bounds")
)

type Point struct {
	scaledX, scaledY int64
	plane            *Plane
}

func (pt *Point) Eq(other *Point) bool {
	return pt.scaledX == other.scaledX && pt.scaledY == other.scaledY && pt.plane.scaleFactor == other.plane.scaleFactor
}

func (pt *Point) X() float64 {
	return float64(pt.scaledX) / pt.plane.scaleFactor
}

func (pt *Point) Y() float64 {
	return float64(pt.scaledY) / pt.plane.scaleFactor
}
