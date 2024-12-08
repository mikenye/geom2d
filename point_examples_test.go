package geom2d_test

import (
	"fmt"
	"geom2d"
	"math"
)

func ExampleEnsureClockwise() {

	// Define points in counter-clockwise order
	points := []geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(4, 0),
		geom2d.NewPoint(4, 4),
	}

	// Re-order points in clockwise order
	geom2d.EnsureClockwise(points)

	// Print re-ordered points
	fmt.Println(points)
	// Output:
	// [Point[(4, 4)] Point[(4, 0)] Point[(0, 0)]]
}

func ExampleEnsureCounterClockwise() {

	// Define points in clockwise order
	points := []geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(4, 4),
		geom2d.NewPoint(4, 0),
	}

	// Re-order points in counter-clockwise order
	geom2d.EnsureCounterClockwise(points)

	// Print re-ordered points
	fmt.Println(points)
	// Output:
	// [Point[(4, 0)] Point[(4, 4)] Point[(0, 0)]]
}

func ExampleRelativeAngle() {
	A := geom2d.NewPoint[int](1, 0)
	B := geom2d.NewPoint[int](0, 1)
	radians := geom2d.RelativeAngle(A, B) // π/2 radians
	degrees := radians * 180 / math.Pi    // convert to degrees
	fmt.Println(degrees)
	// Output:
	// 90
}

func ExampleRelativeCosineOfAngle() {
	A := geom2d.NewPoint(1, 0)
	B := geom2d.NewPoint(0, 1)
	cosine := geom2d.RelativeCosineOfAngle(A, B) // cosine is 0 for a 90-degree angle
	fmt.Println(cosine)
	// Output:
	// 0
}

func ExampleSignedArea2X() {
	points := []geom2d.Point[int]{
		geom2d.NewPoint(0, 0),
		geom2d.NewPoint(4, 0),
		geom2d.NewPoint(4, 3),
	}
	signedArea := geom2d.SignedArea2X(points)
	fmt.Println(signedArea)
	// Output:
	// 12
}
