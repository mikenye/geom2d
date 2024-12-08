package geom2d_test

import (
	"fmt"
	"geom2d"
)

func ExampleNewLineSegment() {
	segment := geom2d.NewLineSegment(geom2d.NewPoint(0, 0), geom2d.NewPoint(3, 4))
	fmt.Println(segment.String())
	// Output:
	// LineSegment[(0, 0) -> (3, 4)]
}

func ExampleLineSegment_AddLineSegment() {
	segment1 := geom2d.NewLineSegment(geom2d.NewPoint(1, 1), geom2d.NewPoint(4, 5))
	segment2 := geom2d.NewLineSegment(geom2d.NewPoint(2, 3), geom2d.NewPoint(1, 2))
	result := segment1.AddLineSegment(segment2)
	fmt.Println(result)
	// Output:
	// LineSegment[(3, 4) -> (5, 7)]
}

func ExampleLineSegment_Area() {
	segment := geom2d.NewLineSegment(geom2d.NewPoint(1, 2), geom2d.NewPoint(3, 4))
	fmt.Println(segment.Area())
	// Output:
	// 0
}

func ExampleLineSegment_AsFloat() {
	// Create a LineSegment with integer coordinates
	intSegment := geom2d.NewLineSegment(
		geom2d.NewPoint(1, 2), // Start point
		geom2d.NewPoint(3, 4), // End point
	)

	// Convert the LineSegment to float64
	floatSegment := intSegment.AsFloat()

	// Print the converted LineSegment
	fmt.Println("Integer LineSegment:", intSegment)
	fmt.Println("Float64 LineSegment:", floatSegment)
	// Output:
	// Integer LineSegment: LineSegment[(1, 2) -> (3, 4)]
	// Float64 LineSegment: LineSegment[(1, 2) -> (3, 4)]
}

func ExampleLineSegment_AsInt() {
	// Create a LineSegment with floating-point coordinates
	line := geom2d.NewLineSegment(
		geom2d.NewPoint(1.5, 2.7),
		geom2d.NewPoint(3.9, 4.2),
	)

	// Convert the LineSegment to integer coordinates
	intLine := line.AsInt()

	// Output the integer LineSegment
	fmt.Println(intLine)
	// Output:
	// LineSegment[(1, 2) -> (3, 4)]
}

func ExampleLineSegment_AsIntRounded() {
	// Create a LineSegment with floating-point coordinates
	line := geom2d.NewLineSegment(
		geom2d.NewPoint(1.5, 2.7),
		geom2d.NewPoint(3.9, 4.2),
	)

	// Convert the LineSegment to integer coordinates with rounding
	roundedIntLine := line.AsIntRounded()

	// Output the rounded integer LineSegment
	fmt.Println(roundedIntLine)
	// Output:
	// LineSegment[(2, 3) -> (4, 4)]
}

func ExampleLineSegment_BoundingBox() {
	// Create a LineSegment
	line := geom2d.NewLineSegment(
		geom2d.NewPoint(3, 1),
		geom2d.NewPoint(6, 4),
	)

	// Compute the bounding box of the LineSegment
	boundingBox := line.BoundingBox()

	// Output the bounding box
	fmt.Println(boundingBox)
	// Output:
	// Rectangle[(3, 1), (6, 1), (6, 4), (3, 4)]
}

func ExampleLineSegment_Center() {
	lineSegment := geom2d.NewLineSegment(geom2d.NewPoint(0, 0), geom2d.NewPoint(10, 10))
	center := lineSegment.Center()
	fmt.Printf("Center: %v\n", center)
	// Output:
	// Center: Point[(5, 5)]
}
