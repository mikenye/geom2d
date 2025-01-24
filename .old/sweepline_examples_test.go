package _old_test

import (
	"fmt"
	"github.com/mikenye/geom2d"
)

func ExampleSweepLine() {
	// Define a set of line segments.
	lineSegments := []geom2d.LineSegment[int]{
		geom2d.NewLineSegment(geom2d.NewPoint(0, 0), geom2d.NewPoint(10, 4)),
		geom2d.NewLineSegment(geom2d.NewPoint(0, 2), geom2d.NewPoint(10, 0)),
		geom2d.NewLineSegment(geom2d.NewPoint(0, 4), geom2d.NewPoint(10, 2)),
		geom2d.NewLineSegment(geom2d.NewPoint(1, -1), geom2d.NewPoint(10, 1)),
	}

	// Call SweepLine to find intersections.
	result := geom2d.SweepLine(lineSegments)

	// Print the intersection points.
	fmt.Println("Intersection Points:")
	for _, point := range result.IntersectionPoints {
		fmt.Printf("(%.4f, %.4f)\n", point.X(), point.Y())
	}

	// Output:
	// Intersection Points:
	// (3.3333, 1.3333)
	// (6.6667, 2.6667)
	// (7.6316, 0.4737)
}
